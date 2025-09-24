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

	"github.com/GarotoCowboy/vttProject/api/grpc/proto/imageLibrary/pb"
	"github.com/GarotoCowboy/vttProject/api/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

const maxImageSize = 10 << 20

func (s *ImageLibraryService) UploadImage(stream pb.ImageLibraryService_UploadImageServer) error {
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
	res := &pb.UploadImageResponse{
		Image: &pb.Image{
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
		},
	}

	return stream.SendAndClose(res)
}
func (s *ImageLibraryService) EditImage(ctx context.Context, req *pb.EditImageRequest) (*pb.Image, error) {
	s.Logger.InfoF("GRPC Requisition to EditImage started...")
	//validate the mask
	updatesMap, err := ValidadeAndBuildUpdateMap(req)
	if err != nil {
		return nil, err
	}

	//create models
	var image models.Image
	var table models.Table

	s.Logger.InfoF("Searching on DB...")

	//search if table exists on DB
	if err := s.DB.Where("id = ?", req.GetImage().GetTableId()).First(&table).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "invalid table id: %v", "table id not found")
		}
		return nil, status.Errorf(codes.NotFound, "internal error: %v", err.Error())
	}

	if err := s.DB.Where("id = ?", req.GetImage().GetImageId()).First(&image).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "invalid image id: %v", "image id not found")
		}
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}

	if err := s.DB.Model(&image).Updates(updatesMap).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}

	s.Logger.Info("GRPC Requisition to EditImage finished...")

	responseImage := &pb.Image{
		TableId:     uint64(image.TableID),
		ImageId:     uint64(image.ID),
		Name:        image.Name,
		Width:       uint32(image.Width),
		Height:      uint32(image.Height),
		ContentType: image.ContentType,
		CreatedAt:   timestamppb.New(image.CreatedAt),
		UpdatedAt:   timestamppb.New(image.UpdatedAt),
	}

	return responseImage, nil
}
func (s *ImageLibraryService) DeleteImage(ctx context.Context, req *pb.DeleteImageRequest) (*emptypb.Empty, error) {
	s.Logger.InfoF("Received request to delete image: image_id=%s, table_id=%s", req.ImageId, req.TableId)

	//VALIDATION SIMPLE, I WILL CREATE AN VALIDADE ARCHIVE
	if req.GetImageId() == 0 || req.GetTableId() == 0 {
		s.Logger.ErrorF("DeleteImage requires imageId and tableId field")
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var image models.Image

	//SEARCH IMAGE ON DB
	if err := s.DB.Where("id = ? AND table_id = ?", req.ImageId, req.TableId).First(&image).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.ErrorF("Image not found: image_id=%d, table_id=%d", req.ImageId, req.TableId)
			return nil, status.Error(codes.Internal, "image not found")
		}
	}
	//DELETE IMAGE FROM DISK
	if err := os.Remove(image.ImagePath); err != nil {
		s.Logger.ErrorF("Failed to delete image from disk: %v, but continuing to delete record from DB: %v", image.ImagePath, err)
	} else {
		s.Logger.ErrorF("Error to delete image from disk: %s", image.ImagePath)
	}

	if err := s.DB.Delete(&image).Error; err != nil {
		s.Logger.ErrorF("Failed to delete image from DB: %v", err)
		return nil, status.Error(codes.Internal, "error to delete image")
	}

	s.Logger.InfoF("Deleted image from DB: %s", image.ImagePath)

	return &emptypb.Empty{}, nil
}
func (s *ImageLibraryService) ListImages(ctx context.Context, req *pb.ListImagesRequest) (*pb.ListImagesResponse, error) {
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

	var protoImages []*pb.Image
	for _, dbImg := range images {
		publicURL := fmt.Sprintf("http://localhost:8081/%s", filepath.ToSlash(dbImg.ImagePath))
		protoImages = append(protoImages, &pb.Image{
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
	response := &pb.ListImagesResponse{
		Images:        protoImages,
		NextPageToken: nextPageToken,
	}
	s.Logger.InfoF("GRPC Requisit to ListImages on table %d", req.TableId)
	return response, nil
}
