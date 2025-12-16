package libraryImg

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/GarotoCowboy/vttProject/api/grpc/events"
	"github.com/GarotoCowboy/vttProject/api/grpc/pb/imageLibrary"
	"github.com/GarotoCowboy/vttProject/api/models"
	"github.com/GarotoCowboy/vttProject/api/models/consts/pubSubSyncConst"
	"github.com/GarotoCowboy/vttProject/api/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// maxImageSize defines the maximum allowed image size (10 MB)
const maxImageSize = 10 << 20 // 10 * 2^20 bytes = 10 MiB

func (s *ImageLibraryService) UploadImage(stream imageLibrary.ImageLibraryService_UploadImageServer) error {
	s.Logger.InfoF("GRPC Requisition to Upload Image started... ")

	// --- 1. Receive Initial Metadata ---
	// The first message from the client must be the 'init' message containing metadata
	req, err := stream.Recv()
	if err != nil {
		s.Logger.ErrorF("failed to receive first message from stream: %v ", err)
		return status.Error(codes.Internal, "cannot possible receive data from image")
	}

	// Get the 'init' payload from the first message
	initMsg := req.GetInit()
	// Validate the initial metadata (e.g., check for TableId, Name)
	if err := ValidateUploadInit(initMsg); err != nil {
		s.Logger.ErrorF("failed to validate initial upload message:  %v", err)
		return status.Errorf(codes.InvalidArgument, "invalid initial upload data: %v", err)
	}

	s.Logger.InfoF("Receiving image '%s' for table %d", initMsg.Name, initMsg.TableId)

	// --- 2. Receive Image Chunks ---
	imageData := bytes.Buffer{}
	imageSize := 0

	// Loop to receive all subsequent messages, which must be chunks
	for {
		// Receive the next message in the stream
		req, err := stream.Recv()
		// io.EOF signals the client has finished sending chunks
		if err == io.EOF {
			s.Logger.Info("End of received image stream (EOF)")
			break
		}
		if err != nil {
			s.Logger.ErrorF("failed to receive chunk from stream: %v ", err)
			return status.Error(codes.Internal, "cannot possible receive chunk from stream")
		}

		// Get the chunk data
		chunk := req.GetChunk()
		if chunk == nil {
			s.Logger.ErrorF("Invalid stream: received non-init message that was not a chunk")
			return status.Errorf(codes.InvalidArgument, "after initialization, all messages must be chunks")
		}

		// --- 3. Validate Size and Write Chunk ---
		imageSize += len(chunk.Data)
		// Enforce the maximum image size limit
		if imageSize > maxImageSize {
			s.Logger.ErrorF("Image size limit exceeded: %d > %d", imageSize, maxImageSize)
			return status.Errorf(codes.InvalidArgument, "image size is too big, max allowed is %dMB", maxImageSize>>20)
		}

		// Write the received chunk data to the in-memory buffer
		_, err = imageData.Write(chunk.Data)
		if err != nil {
			s.Logger.ErrorF("Could not write image chunk to buffer: %v", err)
			return status.Errorf(codes.Internal, "could not write image chunk")
		}
	}

	s.Logger.InfoF("Total image size received: %d bytes. Validating image content...", imageSize)

	// --- 4. Validate Image Content ---
	imageBytes := imageData.Bytes()
	// Decode image config to validate format and get dimensions
	imgConfig, _, err := validateImageContent(imageBytes)
	if err != nil {
		s.Logger.ErrorF("Image content validation failed: %v", err)
		return status.Errorf(codes.InvalidArgument, err.Error())
	}

	// Detect and validate the MIME type
	contentType := http.DetectContentType(imageBytes)
	if contentType != "image/jpeg" && contentType != "image/png" { // Adjusted check
		s.Logger.ErrorF("Invalid content type: %s. Only JPEG or PNG allowed.", contentType)
		return status.Errorf(codes.InvalidArgument, "file type not allowed, use PNG or JPEG")
	}
	s.Logger.InfoF("Image validated: %s, %dx%d", contentType, imgConfig.Width, imgConfig.Height)

	// --- 5. Save Image to Disk ---
	s.Logger.InfoF("Saving content to disk...")
	// Generate checksum to prevent duplicates and create a unique filename
	checksum := fmt.Sprintf("%x", sha256.Sum256(imageBytes))
	fileName := fmt.Sprintf("%s%s", checksum, filepath.Ext(initMsg.Name)) // Use original extension
	tableIdStr := fmt.Sprintf("%d", initMsg.TableId)

	// Define the directory path: vttData/[TableID]/img
	dirPath := filepath.Join("vttData", tableIdStr, "img")

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		s.Logger.ErrorF("Failed to create table directory (%s): %v", dirPath, err)
		return status.Error(codes.Internal, "failed to prepare image storage")
	}

	// Define the full file path
	filePath := filepath.Join(dirPath, fileName)

	// Write the image bytes to the file
	err = os.WriteFile(filePath, imageBytes, 0664)
	if err != nil {
		s.Logger.ErrorF("failed to save image to disk: %v", err)
		return status.Error(codes.Internal, "failed to save image")
	}
	s.Logger.InfoF("Saved image to disk: %s", filePath)

	// --- 6. Save Metadata to Database ---
	imageModel := &models.Image{
		Name:        initMsg.Name, // Original uploaded name
		TableID:     uint(initMsg.TableId),
		ContentType: contentType,
		ImagePath:   filePath, // Path in the filesystem
		CheckSum:    checksum,
		Width:       uint(imgConfig.Width),
		Height:      uint(imgConfig.Height),
		// CanBeViewedBy is likely defaulting to a value, e.g., "All" or "GM"
	}

	if err := s.DB.Create(imageModel).Error; err != nil {
		s.Logger.ErrorF("Failed to save image metadata to database: %v", err)
		// Attempt to clean up the saved file if DB insert fails
		os.Remove(filePath)
		return status.Error(codes.Internal, "failed to save image to database")
	}
	s.Logger.InfoF("Saved image metadata to DB with ID: %d", imageModel.ID)

	// --- 7. Prepare Response and Publish Event ---
	responseImage := &imageLibrary.Image{
		TableId:     uint64(imageModel.TableID),
		ImageId:     uint64(imageModel.ID),
		Name:        imageModel.Name,
		Width:       uint32(imgConfig.Width),
		Height:      uint32(imgConfig.Height),
		ContentType: imageModel.ContentType,
		ImageUrl:    filePath, // Note: This is the internal path, client might need a public URL
		CreatedAt:   timestamppb.New(imageModel.CreatedAt),
		UpdatedAt:   timestamppb.New(imageModel.UpdatedAt),
		Checksum:    checksum,
	}
	res := &imageLibrary.UploadImageResponse{
		Image: responseImage,
	}

	s.Logger.InfoF("Publishing ImageLibraryUploaded event to TableSync for table: %d", imageModel.TableID)
	event := events.NewImageLibraryUploadedEvent(responseImage)
	// Publish to the table channel so all users in that table get the update
	s.Broker.Publish(pubSubSyncConst.TableSync, uint64(imageModel.TableID), event)

	s.Logger.InfoF("Image uploaded successfully")
	// Send the final response and close the stream
	return stream.SendAndClose(res)
}

func (s *ImageLibraryService) EditImage(ctx context.Context, req *imageLibrary.EditImageRequest) (*imageLibrary.Image, error) {
	s.Logger.InfoF("GRPC Requisition to EditImage started...")

	// Validate the request and field mask, build a map of fields to update
	updatesMap, err := ValidadeAndBuildUpdateMap(req)
	if err != nil {
		s.Logger.ErrorF("Invalid edit image request or field mask: %v", err)
		return nil, err
	}

	// Note: Authentication (PickUserIdJWT) is missing here, but should be present.
	// Also, a CheckUserIsMaster check is likely needed.

	// create models
	var imageModel models.Image
	//var tableModel models.Table // This is fetched but not used

	s.Logger.InfoF("Searching on DB for image %d on table %d", req.GetImage().GetImageId(), req.GetImage().GetTableId())

	// Start a transaction for the update
	err = s.DB.Transaction(func(tx *gorm.DB) error {
		// Check if the table exists
		// This check is redundant if we only check the image and its TableID matches
		if err := tx.WithContext(ctx).Where("id = ?", req.GetImage().GetTableId()).First(&models.Table{}).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Logger.ErrorF("Table not found: %d", req.GetImage().GetTableId())
				return status.Errorf(codes.NotFound, "invalid tableModel id: %v", "tableModel id not found")
			}
			s.Logger.ErrorF("Error fetching table: %v", err)
			return status.Errorf(codes.Internal, "internal error: %v", err.Error())
		}

		// Find the image to be updated
		if err := tx.WithContext(ctx).Where("id = ?", req.GetImage().GetImageId()).First(&imageModel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Logger.ErrorF("Image not found: %d", req.GetImage().GetImageId())
				return status.Errorf(codes.NotFound, "invalid imageModel id: %v", "imageModel id not found")
			}
			s.Logger.ErrorF("Error fetching image: %v", err)
			return status.Errorf(codes.Internal, "internal error: %v", err.Error())
		}

		// Security check: Ensure the found image belongs to the table specified in the request
		if imageModel.TableID != uint(req.GetImage().GetTableId()) {
			s.Logger.ErrorF("Image %d does not belong to table %d (belongs to %d)", imageModel.ID, req.GetImage().GetTableId(), imageModel.TableID)
			return status.Errorf(codes.NotFound, "image not found in the specified table")
		}

		if err = utils.CheckUserIsMaster(ctx, tx, imageModel.TableID); err != nil {
			s.Logger.ErrorF("%v", err)
			return err
		}

		s.Logger.InfoF("Updating image %d with data: %v", imageModel.ID, updatesMap)
		// Apply the updates from the map
		if err := tx.WithContext(ctx).Model(&imageModel).Updates(updatesMap).Error; err != nil {
			s.Logger.ErrorF("Failed to update image in DB: %v", err)
			return status.Errorf(codes.Internal, "internal error: %v", err.Error())
		}

		// Commit transaction
		return nil
	})

	// Check for transaction error
	if err != nil {
		s.Logger.ErrorF("Transaction failed: %v", err)
		// Don't wrap internal errors, just return them
		return nil, err
	}

	// Prepare the response protobuf
	responseImage := &imageLibrary.Image{
		TableId:     uint64(imageModel.TableID),
		ImageId:     uint64(imageModel.ID),
		Name:        imageModel.Name,
		Width:       uint32(imageModel.Width),
		Height:      uint32(imageModel.Height),
		ContentType: imageModel.ContentType,
		CreatedAt:   timestamppb.New(imageModel.CreatedAt),
		UpdatedAt:   timestamppb.New(imageModel.UpdatedAt), // Will be updated by GORM
		// ImageUrl is missing, may be needed by client
	}

	s.Logger.InfoF("Publishing ImageLibraryUpdated event to TableSync on table: %d ", imageModel.TableID)
	// Create and publish the update event
	event := events.NewImageLibraryUpdatedEvent(responseImage)
	s.Broker.Publish(pubSubSyncConst.TableSync, uint64(imageModel.TableID), event)

	s.Logger.Info("GRPC Requisition to EditImage finished...")
	return responseImage, nil
}

func (s *ImageLibraryService) DeleteImage(ctx context.Context, req *imageLibrary.DeleteImageRequest) (*emptypb.Empty, error) {
	s.Logger.InfoF("Received request to delete image: image_id=%d, table_id=%d", req.ImageId, req.TableId)

	// TODO: Add Authentication (PickUserIdJWT) and Authorization (CheckUserIsMaster)

	// --- Validation ---
	if req.GetImageId() == 0 || req.GetTableId() == 0 {
		s.Logger.ErrorF("DeleteImage requires ImageId and TableId fields")
		return nil, status.Error(codes.InvalidArgument, "invalid request: ImageId and TableId are required")
	}

	var imageModel models.Image

	// --- 1. Find Image in DB ---
	// Find the image matching both ID and TableID
	if err := s.DB.Where("id = ? AND table_id = ?", req.ImageId, req.TableId).First(&imageModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.ErrorF("Image not found: image_id=%d, table_id=%d", req.ImageId, req.TableId)
			return nil, status.Error(codes.NotFound, "image not found") // Corrected from Internal
		}
		s.Logger.ErrorF("Error fetching image from DB: %v", err)
		return nil, status.Error(codes.Internal, "failed to fetch image")
	}

	// --- 2. Delete Image from Disk ---
	if imageModel.ImagePath != "" { // Only attempt delete if path exists
		if err := os.Remove(imageModel.ImagePath); err != nil {
			// Log as an error but don't stop. The DB record is the source of truth.
			s.Logger.ErrorF("Failed to delete image from disk (%s): %v. Continuing to delete DB record.", imageModel.ImagePath, err)
		} else {
			s.Logger.InfoF("Successfully deleted image from disk: %s", imageModel.ImagePath)
		}
	} else {
		s.Logger.WarningF("Image %d had no ImagePath set, skipping disk delete.", imageModel.ID)
	}

	if err := utils.CheckUserIsMaster(ctx, s.DB, imageModel.TableID); err != nil {
		return nil, err
	}

	// --- 3. Delete Image from DB ---
	// This will also cascade delete placed images if constraints are set correctly
	if err := s.DB.Delete(&imageModel).Error; err != nil {
		s.Logger.ErrorF("Failed to delete image from DB: %v", err)
		return nil, status.Error(codes.Internal, "error to delete image")
	}

	s.Logger.InfoF("Deleted image metadata from DB: ID %d", imageModel.ID)

	// --- 4. Publish Event ---
	event := events.NewImageLibraryDeletedEvent(req.ImageId, req.TableId)
	s.Broker.Publish(pubSubSyncConst.TableSync, req.TableId, event)
	s.Logger.InfoF("Published ImageLibraryDeleted event to TableSync for table %d", req.TableId)

	return &emptypb.Empty{}, nil
}

func (s *ImageLibraryService) ListImages(ctx context.Context, req *imageLibrary.ListImagesRequest) (*imageLibrary.ListImagesResponse, error) {
	s.Logger.InfoF("GRPC Requisition to ListImages on table %d", req.TableId)

	// TODO: Add Authentication (PickUserIdJWT).
	// Authorization (CheckUserIsMember) might be needed if library isn't public to all table members.

	// --- 1. Validate Parameters ---
	if req.GetTableId() <= 0 {
		s.Logger.ErrorF("TableId must be greater than zero")
		return nil, status.Errorf(codes.InvalidArgument, "Table id must be greater than zero")
	}

	// Define page size with defaults and limits
	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 25 // Default page size
	}
	if pageSize > 100 {
		pageSize = 100 // Max page size
	}

	// --- 2. Build Query ---
	// Base query for the specified table, ordered by ID (keyset pagination)
	query := s.DB.Where("table_id = ?", req.TableId).Order("id ASC")

	// If PageToken (which is the last ID from the previous page) is provided,
	// fetch items *after* that ID.
	if req.PageToken != "" {
		lastID, err := strconv.ParseInt(req.PageToken, 10, 64)
		if err != nil {
			s.Logger.ErrorF("Invalid page token format: %s", req.PageToken)
			return nil, status.Errorf(codes.InvalidArgument, "invalid page token")
		}
		query = query.Where("id > ?", lastID)
	}

	var images []*models.Image

	// --- 3. Execute Query ---
	// Fetch (pageSize + 1) items to check if there's a next page
	if err := query.Limit(pageSize + 1).Find(&images).Error; err != nil {
		s.Logger.ErrorF("Failed to list images from DB: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to list images from library")
	}

	// --- 4. Determine Next Page Token ---
	var nextPageToken string
	// If we fetched more items than the page size, a next page exists
	if len(images) > pageSize {
		// The next page token is the ID of the last item *on this page*
		nextPageToken = fmt.Sprintf("%d", images[pageSize-1].ID)
		// Trim the slice to only include items for *this page*
		images = images[:pageSize]
	}

	// --- 5. Format Response ---
	var protoImages []*imageLibrary.Image
	for _, dbImg := range images {
		// TODO: This URL should point to a static file server, not the internal file path.
		// Example: "http://my-cdn.com/vttData/1/img/checksum.png"
		// Using a hardcoded localhost URL for now.
		publicURL := fmt.Sprintf("http://localhost:8081/%s", filepath.ToSlash(dbImg.ImagePath))

		protoImages = append(protoImages, &imageLibrary.Image{
			TableId:     uint64(dbImg.TableID),
			ImageId:     uint64(dbImg.ID),
			Name:        dbImg.Name,
			Width:       uint32(dbImg.Width),
			Height:      uint32(dbImg.Height),
			ContentType: dbImg.ContentType,
			CreatedAt:   timestamppb.New(dbImg.CreatedAt),
			UpdatedAt:   timestamppb.New(dbImg.UpdatedAt),
			Checksum:    dbImg.CheckSum,
			ImageUrl:    publicURL, // Use the publicly accessible URL
			// CanBeViewedBy is missing, might be needed by client to filter
		})
	}

	response := &imageLibrary.ListImagesResponse{
		Images:        protoImages,
		NextPageToken: nextPageToken,
	}

	s.Logger.InfoF("Successfully listed %d images for table %d", len(protoImages), req.TableId)
	return response, nil
}
