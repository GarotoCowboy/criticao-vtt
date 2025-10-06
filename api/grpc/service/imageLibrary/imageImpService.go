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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

const maxImageSize = 10 << 20

func (s *ImageLibraryService) UploadImage(stream imageLibrary.ImageLibraryService_UploadImageServer) error {
	s.Logger.InfoF("GRPC Requisition to Upload Image started... ")

	//Receive  first metadata
	req, err := stream.Recv()

	if err != nil {
		s.Logger.ErrorF("failed to receive first message from stream: %v ", err)
		return status.Error(codes.Internal, "cannot possible receive data from image")
	}

	initMsg := req.GetInit()

	if err := ValidateUploadInit(initMsg); err != nil {
		s.Logger.ErrorF("failed to receive first message from stream:  %v", err)
		return status.Errorf(codes.Internal, "cannot possible receive data from image")
	}

	s.Logger.InfoF("Receiving image '%s' for table %d", initMsg.Name, initMsg.TableId)

	imageData := bytes.Buffer{}
	imageSize := 0

	// receive chunks from image
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			s.Logger.Info("End of received image stream")
			break
		}
		if err != nil {
			s.Logger.ErrorF("failed to receive chunk from stream: %v ", err)
			return status.Error(codes.Internal, "cannot possible receive chunk from stream")
		}
		chunk := req.GetChunk()
		if chunk == nil {
			return status.Errorf(codes.InvalidArgument, "after initialization, all messages must be chunks")
		}

		imageSize += len(chunk.Data)
		if imageSize > maxImageSize {
			return status.Errorf(codes.InvalidArgument, "image size is too big, max allowed is %dMB", maxImageSize>>20)
		}
		_, err = imageData.Write(chunk.Data)
		if err != nil {
			return status.Errorf(codes.Internal, "could not write image chunk")
		}
	}

	s.Logger.InfoF("validating image content")
	//validate image content
	imageBytes := imageData.Bytes()
	imgConfig, _, err := validateImageContent(imageBytes)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}

	contentType := http.DetectContentType(imageBytes)
	if contentType != "image/jpg" && contentType != "image/jpeg" && contentType != "image/png" {
		return status.Errorf(codes.InvalidArgument, "file type not allowed, use PNG,JPEG or JPG")
	}

	s.Logger.InfoF("saving content on db")
	//Save Image and register on DB
	checksum := fmt.Sprintf("%x", sha256.Sum256(imageBytes))
	fileName := fmt.Sprintf("%s%s", checksum, filepath.Ext(initMsg.Name))
	tableIdStr := fmt.Sprintf("%d", initMsg.TableId)
	dirPath := filepath.Join("vttData", tableIdStr, "img")

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		s.Logger.ErrorF("Failed to create table directory (%s): %v", dirPath, err)
		return status.Error(codes.Internal, "failed to prepare image storage")
	}
	filePath := filepath.Join(dirPath, fileName)

	err = os.WriteFile(filePath, imageBytes, 0664)
	if err != nil {
		s.Logger.ErrorF("failed to save image to disk: %v", err)
		return status.Error(codes.Internal, "failed to save image")
	}

	s.Logger.InfoF("Saved image to disk: %s", filePath)

	imageModel := &models.Image{
		Name:        initMsg.Name,
		TableID:     uint(initMsg.TableId),
		ContentType: contentType,
		ImagePath:   filePath,
		CheckSum:    checksum,
		Width:       uint(imgConfig.Width),
		Height:      uint(imgConfig.Height),
	}

	if err := s.DB.Create(imageModel).Error; err != nil {
		s.Logger.ErrorF("Failed to save image to database: %v", err)
		return status.Error(codes.Internal, "failed to save image to database")
	}

	responseImage := &imageLibrary.Image{
		TableId:     uint64(imageModel.TableID),
		ImageId:     uint64(imageModel.ID),
		Name:        imageModel.Name,
		Width:       uint32(imgConfig.Width),
		Height:      uint32(imgConfig.Height),
		ContentType: imageModel.ContentType,
		ImageUrl:    filePath,
		CreatedAt:   timestamppb.New(imageModel.CreatedAt),
		UpdatedAt:   timestamppb.New(imageModel.UpdatedAt),
		Checksum:    checksum,
	}
	res := &imageLibrary.UploadImageResponse{
		Image: responseImage,
	}

	s.Logger.InfoF("publishing event on table: %d", imageModel.TableID)
	event := events.NewImageLibraryUploadedEvent(responseImage)
	s.Broker.Publish(pubSubSyncConst.TableSync, uint64(imageModel.ID), event)

	s.Logger.InfoF("image uploaded with sucess")
	return stream.SendAndClose(res)
}
func (s *ImageLibraryService) EditImage(ctx context.Context, req *imageLibrary.EditImageRequest) (*imageLibrary.Image, error) {
	s.Logger.InfoF("GRPC Requisition to EditImage started...")
	//validate the mask
	updatesMap, err := ValidadeAndBuildUpdateMap(req)
	if err != nil {
		return nil, err
	}

	//create models
	var imageModel models.Image
	var tableModel models.Table

	s.Logger.InfoF("Searching on DB...")

	err = s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Where("id = ?", req.GetImage().GetTableId()).First(&tableModel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return status.Errorf(codes.NotFound, "invalid tableModel id: %v", "tableModel id not found")
			}
			return status.Errorf(codes.NotFound, "internal error: %v", err.Error())
		}

		if err := tx.WithContext(ctx).Where("id = ?", req.GetImage().GetImageId()).First(&imageModel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return status.Errorf(codes.NotFound, "invalid imageModel id: %v", "imageModel id not found")
			}
			return status.Errorf(codes.Internal, "internal error: %v", err.Error())
		}

		if err := tx.WithContext(ctx).Model(&imageModel).Updates(updatesMap).Error; err != nil {
			return status.Errorf(codes.Internal, "internal error: %v", err.Error())
		}
		return nil
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}

	responseImage := &imageLibrary.Image{
		TableId:     uint64(imageModel.TableID),
		ImageId:     uint64(imageModel.ID),
		Name:        imageModel.Name,
		Width:       uint32(imageModel.Width),
		Height:      uint32(imageModel.Height),
		ContentType: imageModel.ContentType,
		CreatedAt:   timestamppb.New(imageModel.CreatedAt),
		UpdatedAt:   timestamppb.New(imageModel.UpdatedAt),
	}

	s.Logger.InfoF("publishing update event on tableModel: %d ", imageModel.TableID)
	event := events.NewImageLibraryUpdatedEvent(responseImage)
	s.Broker.Publish(pubSubSyncConst.TableSync, uint64(imageModel.TableID), event)

	s.Logger.Info("GRPC Requisition to EditImage finished...")
	return responseImage, nil
}
func (s *ImageLibraryService) DeleteImage(ctx context.Context, req *imageLibrary.DeleteImageRequest) (*emptypb.Empty, error) {
	s.Logger.InfoF("Received request to delete imageModel: image_id=%s, table_id=%s", req.ImageId, req.TableId)

	//VALIDATION SIMPLE, I WILL CREATE AN VALIDADE ARCHIVE
	if req.GetImageId() == 0 || req.GetTableId() == 0 {
		s.Logger.ErrorF("DeleteImage requires imageId and tableId field")
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var imageModel models.Image

	//SEARCH IMAGE ON DB
	if err := s.DB.Where("id = ? AND table_id = ?", req.ImageId, req.TableId).First(&imageModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.ErrorF("Image not found: image_id=%d, table_id=%d", req.ImageId, req.TableId)
			return nil, status.Error(codes.Internal, "imageModel not found")
		}
	}
	//DELETE IMAGE FROM DISK
	if err := os.Remove(imageModel.ImagePath); err != nil {
		s.Logger.ErrorF("Failed to delete imageModel from disk: %v, but continuing to delete record from DB: %v", imageModel.ImagePath, err)
	} else {
		s.Logger.ErrorF("Error to delete imageModel from disk: %s", imageModel.ImagePath)
	}

	if err := s.DB.Delete(&imageModel).Error; err != nil {
		s.Logger.ErrorF("Failed to delete imageModel from DB: %v", err)
		return nil, status.Error(codes.Internal, "error to delete imageModel")
	}

	s.Logger.InfoF("Deleted imageModel from DB: %s", imageModel.ImagePath)

	event := events.NewImageLibraryDeletedEvent(uint64(imageModel.ID), req.TableId)
	s.Broker.Publish(pubSubSyncConst.TableSync, req.TableId, event)

	return &emptypb.Empty{}, nil
}
func (s *ImageLibraryService) ListImages(ctx context.Context, req *imageLibrary.ListImagesRequest) (*imageLibrary.ListImagesResponse, error) {
	s.Logger.InfoF("GRPC Requisiton to ListImages on table %d", req.TableId)

	if req.GetTableId() <= 0 {
		s.Logger.ErrorF("TableId must be greater than zero")
		return nil, status.Errorf(codes.InvalidArgument, "Table id must be greater than zero")
	}

	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 25
	}

	if pageSize > 100 {
		pageSize = 100
	}

	query := s.DB.Where("table_id = ?", req.TableId).Order("id ASC")

	if req.PageToken != "" {
		lastID, err := strconv.ParseInt(req.PageToken, 10, 64)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid page token")
		}
		query = query.Where("id > ?", lastID)
	}

	var images []*models.Image

	if err := query.Limit(pageSize + 1).Find(&images).Error; err != nil {
		s.Logger.ErrorF("Failed to list images: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to list images from library")
	}

	var nextPageToken string

	if len(images) > pageSize {
		nextPageToken = fmt.Sprintf("%d", images[pageSize-1].ID)
		images = images[:pageSize]
	}

	var protoImages []*imageLibrary.Image
	for _, dbImg := range images {
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
			ImageUrl:    publicURL,
		})
	}
	response := &imageLibrary.ListImagesResponse{
		Images:        protoImages,
		NextPageToken: nextPageToken,
	}
	s.Logger.InfoF("GRPC Requisit to ListImages on table %d", req.TableId)
	return response, nil
}
