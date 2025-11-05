package placedImage

import (
	"context"
	"errors"

	"github.com/GarotoCowboy/vttProject/api/grpc/events"
	"github.com/GarotoCowboy/vttProject/api/grpc/pb/placedImage"
	"github.com/GarotoCowboy/vttProject/api/models"
	"github.com/GarotoCowboy/vttProject/api/models/consts/pubSubSyncConst"
	"github.com/GarotoCowboy/vttProject/api/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func (s *PlacedImageService) CreatePlacedImage(ctx context.Context, req *placedImage.CreatePlacedImageRequest) (*placedImage.CreatePlacedImageResponse, error) {
	s.Logger.InfoF("GRPC Requisition to Create placedImage started...") // Log when the request starts

	// Validate the request body
	if err := Validate(req); err != nil {
		s.Logger.ErrorF("Invalid request: %v", err) // Log if validation fails
		return &placedImage.CreatePlacedImageResponse{}, err
	}

	var sceneModel models.Scene
	var imageModel models.Image

	s.Logger.InfoF("Searching if scene exists...") // Log to search for the scene

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		// Check if the scene exists
		if err := tx.WithContext(ctx).Where("id = ?", req.SceneId).First(&sceneModel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Logger.ErrorF("Scene not found with SceneId: %d", req.SceneId) // Log if the scene is not found
				return status.Errorf(codes.NotFound, "scene not found")
			}
			s.Logger.ErrorF("Error fetching scene with SceneId: %d, error: %v", req.SceneId, err) // Log other errors during scene fetch
			return status.Errorf(codes.Internal, "internal error")
		}

		// Verify if the user is a master for the scene
		s.Logger.InfoF("Checking if user is a master for SceneId: %d", req.SceneId)
		err := utils.CheckUserIsMaster(ctx, tx, sceneModel.TableID)
		if err != nil {
			s.Logger.ErrorF("User does not have master permissions for SceneId: %d", req.SceneId) // Log if user doesn't have master permissions
			return err
		}

		s.Logger.InfoF("Searching if image exists...") // Log to search for the image

		// Check if the image exists
		if err := tx.WithContext(ctx).Where("id = ?", req.ImageId).First(&imageModel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Logger.ErrorF("Image not found with ImageId: %d", req.ImageId) // Log if image is not found
				return status.Errorf(codes.NotFound, "image not found")
			}
			s.Logger.ErrorF("Error fetching image with ImageId: %d, error: %v", req.ImageId, err) // Log other errors during image fetch
			return status.Errorf(codes.Internal, "internal error")
		}

		// Create and save placedImage model
		s.Logger.InfoF("Creating placedImage model and saving...") // Log before saving placed image
		placedImageModel := models.PlacedImage{
			SceneID: sceneModel.ID,
			ImageID: imageModel.ID,
			PosY:    req.PosY,
			PosX:    req.PosX,
		}

		if err := s.DB.Create(&placedImageModel).Error; err != nil {
			s.Logger.ErrorF("Error creating placedImage: %v", err) // Log error if creation fails
			return status.Errorf(codes.Canceled, "cannot create placedImage")
		}

		return nil
	})

	if err != nil {
		s.Logger.WarningF("Cannot create placedImage: %v", err) // Log if transaction fails
		return nil, status.Errorf(codes.Internal, "cannot create placedImage")
	}

	// Prepare the response
	response := &placedImage.PlacedImage{
		SceneId: uint64(sceneModel.ID),
		ImageId: uint64(imageModel.ID),
		PosY:    req.PosY,
		PosX:    req.PosX,
	}

	s.Logger.InfoF("GRPC Requisition to Create placedImage finished...") // Log when the request finishes

	// Publish event after creation
	event := events.NewPlacedImageCreatedEvent(response)
	s.Broker.Publish(pubSubSyncConst.SceneSync, req.SceneId, event)

	return &placedImage.CreatePlacedImageResponse{
		PlacedImage: response,
	}, nil
}

func (s *PlacedImageService) EditPlacedImage(ctx context.Context, req *placedImage.EditPlacedImageRequest) (*placedImage.EditPlacedImageResponse, error) {
	s.Logger.InfoF("GRPC Requisition to Edit placedImage started...") // Log when the request starts

	updatesMap, err := ValidateAndBuildUpdateMap(req)
	if err != nil {
		s.Logger.ErrorF("Invalid request body: %v", err.Error()) // Log if validation fails
		return nil, status.Errorf(codes.InvalidArgument, "invalid request body: %v", err.Error())
	}

	var placedImageModel models.PlacedImage
	var sceneModel models.Scene

	s.Logger.InfoF("Searching if scene and placedImage exist...") // Log to search for the scene and placed image

	err = s.DB.Transaction(func(tx *gorm.DB) error {
		// Check if the scene exists
		if err := tx.WithContext(ctx).Where("id = ?", req.PlacedImage.SceneId).First(&sceneModel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Logger.ErrorF("Scene not found with SceneId: %d", req.PlacedImage.SceneId) // Log if scene is not found
				return status.Errorf(codes.NotFound, "scene not found")
			}
			s.Logger.ErrorF("Error fetching scene with SceneId: %d, error: %v", req.PlacedImage.SceneId, err) // Log other errors during scene fetch
			return status.Errorf(codes.Internal, "internal error")
		}

		// Verify if the user is a master or has permission to edit the object
		s.Logger.InfoF("Verifying if user is a master for SceneId: %d", sceneModel.TableID)
		errMaster := utils.CheckUserIsMaster(ctx, tx, sceneModel.TableID)
		if errMaster == nil {
			s.Logger.InfoF("User is a master, skipping permission check") // Log if user is a master
		} else {
			s.Logger.InfoF("Verifying if user can edit PlacedImageId: %d", req.PlacedImage.PlacedImageId)
			errEdit := utils.CheckUserCanEditTokenObject(ctx, tx, uint(req.PlacedImage.PlacedImageId))
			if errEdit != nil {
				s.Logger.ErrorF("User does not have permission to edit PlacedImageId: %d", req.PlacedImage.PlacedImageId) // Log if user doesn't have permission to edit
				return errEdit
			}
		}

		// Check if the placed image exists
		if err := tx.WithContext(ctx).Where("id = ?", req.PlacedImage.PlacedImageId).First(&placedImageModel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Logger.ErrorF("PlacedImage not found with PlacedImageId: %d", req.PlacedImage.PlacedImageId) // Log if placed image is not found
				return status.Errorf(codes.NotFound, "placed image not found")
			}
			s.Logger.ErrorF("Error fetching placedImage with PlacedImageId: %d, error: %v", req.PlacedImage.PlacedImageId, err) // Log other errors during placed image fetch
			return status.Errorf(codes.Internal, "internal error")
		}

		// Update placed image
		if err := tx.Model(&placedImageModel).Updates(updatesMap).Error; err != nil {
			s.Logger.ErrorF("Error updating placedImage with PlacedImageId: %d, error: %v", placedImageModel.ID, err) // Log if update fails
			return status.Errorf(codes.Internal, "internal error %v", err.Error())
		}
		return nil
	})

	if err != nil {
		s.Logger.ErrorF("Cannot update placedImage: %v", err) // Log if transaction fails
		return nil, status.Errorf(codes.Internal, "cannot update placedImage")
	}

	// Prepare the response
	response := placedImage.PlacedImage{
		SceneId:       uint64(sceneModel.ID),
		ImageId:       req.PlacedImage.ImageId,
		PlacedImageId: uint64(placedImageModel.ID),
		Width:         req.PlacedImage.Width,
		Height:        req.PlacedImage.Height,
		Layer:         req.PlacedImage.Layer,
		Rotation:      req.PlacedImage.Rotation,
		UpdatedAt:     timestamppb.New(placedImageModel.UpdatedAt),
	}

	// Publish event after update
	event := events.NewPlacedImageUpdatedEvent(&response)
	s.Broker.Publish(pubSubSyncConst.SceneSync, uint64(sceneModel.ID), event)

	s.Logger.InfoF("GRPC Requisition to Edit placedImage finished...") // Log when the request finishes

	return &placedImage.EditPlacedImageResponse{
		PlacedImage: &response,
	}, nil
}
func (s *PlacedImageService) DeletePlacedImage(ctx context.Context, req *placedImage.DeletePlacedImageRequest) (*emptypb.Empty, error) {
	s.Logger.InfoF("GRPC Requisition to Delete placedImage started...") // Log when the request starts

	// Check if the required parameters are present in the request
	if req.GetSceneId() == 0 && req.GetPlacedImageId() == 0 {
		s.Logger.ErrorF("Scene ID and PlacedImage ID are required") // Log if IDs are missing
		return &emptypb.Empty{}, status.Error(codes.InvalidArgument, "scene id and placed_image id are required")
	}

	var placedImageModel models.PlacedImage
	var sceneModel models.Scene

	// Begin the database transaction
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		s.Logger.InfoF("Searching if scene exists...") // Log to search for the scene

		// Check if the scene exists
		if err := tx.WithContext(ctx).Where("id = ?", req.SceneId).First(&sceneModel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Logger.ErrorF("Scene not found with SceneId: %d", req.SceneId) // Log if the scene is not found
				return status.Errorf(codes.NotFound, "scene not found")
			}
			s.Logger.ErrorF("Error fetching scene with SceneId: %d, error: %v", req.SceneId, err) // Log other errors during scene fetch
			return status.Errorf(codes.Internal, "internal error")
		}

		// Verify if the user is a master or has permission to edit the object
		s.Logger.InfoF("Verifying if user is a master for SceneId: %d", sceneModel.TableID)
		errMaster := utils.CheckUserIsMaster(ctx, tx, sceneModel.TableID)
		if errMaster == nil {
			st, ok := status.FromError(errMaster)
			if !ok || st.Code() != codes.PermissionDenied {
				s.Logger.ErrorF("User does not have master permissions for SceneId: %d", sceneModel.TableID) // Log if user doesn't have master permissions
				return errMaster
			}
		} else {
			s.Logger.InfoF("Verifying if user can edit PlacedImageId: %d", req.PlacedImageId)
			errEdit := utils.CheckUserCanEditTokenObject(ctx, tx, uint(req.PlacedImageId))
			if errEdit != nil {
				s.Logger.ErrorF("User does not have permission to edit PlacedImageId: %d", req.PlacedImageId) // Log if user doesn't have permission to edit
				return errEdit
			}
		}

		// Check if the placed image exists
		s.Logger.InfoF("Searching if placed image exists...") // Log to search for the placed image
		if err := tx.Where("id = ?", req.PlacedImageId).First(&placedImageModel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Logger.ErrorF("PlacedImage not found with PlacedImageId: %d", req.PlacedImageId) // Log if placed image is not found
				return status.Errorf(codes.NotFound, "placed image not found")
			}
			s.Logger.ErrorF("Error fetching placed image with PlacedImageId: %d, error: %v", req.PlacedImageId, err) // Log other errors during placed image fetch
			return status.Errorf(codes.Internal, "internal error")
		}

		// Attempt to delete the placed image
		if err := tx.Delete(&placedImageModel).Error; err != nil {
			s.Logger.ErrorF("Error deleting placed image with PlacedImageId: %d, error: %v", req.PlacedImageId, err) // Log if deletion fails
			return status.Errorf(codes.Internal, "cannot delete placed image")
		}

		return nil
	})

	if err != nil {
		s.Logger.ErrorF("Error during transaction: %v", err) // Log if the transaction fails
		return &emptypb.Empty{}, err
	}

	// Publish the event after the image is deleted
	event := events.NewPlacedImageDeletedEvent(sceneModel.ID, placedImageModel.ID)
	s.Broker.Publish(pubSubSyncConst.SceneSync, req.SceneId, event)

	s.Logger.InfoF("GRPC Requisition to Delete Scene finished...") // Log when the request finishes
	return &emptypb.Empty{}, nil
}
func (s *PlacedImageService) ListAllImagesOnScene(ctx context.Context, req *placedImage.ListAllImagesRequest) (*placedImage.ListAllImagesResponse, error) {
	s.Logger.InfoF("GRPC Requisition to ListAllImages started...") // Log when the request starts

	// Validate if the scene ID is provided
	if req.SceneId == 0 {
		s.Logger.ErrorF("Scene ID is required") // Log if scene ID is missing
		return nil, status.Errorf(codes.InvalidArgument, "scene id is required")
	}

	var placedImageModels []models.PlacedImage

	// Fetch placed images for the scene
	if err := s.DB.WithContext(ctx).Preload("Image").Where("scene_id = ?", req.SceneId).
		Find(&placedImageModels).Error; err != nil {
		s.Logger.ErrorF("Error fetching placed images for the scene: %v", err) // Log if fetching images fails
	}

	s.Logger.InfoF("Found %d images for scene %d", len(placedImageModels), req.SceneId) // Log the number of images found

	responseTokens := make([]*placedImage.PlacedImage, 0, len(placedImageModels))

	// Prepare the response for each placed image
	for _, model := range placedImageModels {
		responseTokens = append(responseTokens, &placedImage.PlacedImage{
			SceneId:       uint64(model.SceneID),
			PlacedImageId: uint64(model.ID),
			ImageId:       uint64(model.ImageID),
			Width:         uint64(model.Width),
			Height:        uint64(model.Height),
			Layer:         placedImage.LayerType(model.LayerType),
			PosX:          model.PosX,
			PosY:          model.PosY,
			Rotation:      int32(model.Rotation),
			CreatedAt:     timestamppb.New(model.CreatedAt),
			UpdatedAt:     timestamppb.New(model.UpdatedAt),
		})
	}

	return &placedImage.ListAllImagesResponse{
		PlacedImage: responseTokens,
	}, nil
}

func (s *PlacedImageService) MoveImage(ctx context.Context, req *placedImage.MoveImageRequest) (*placedImage.MoveImageResponse, error) {
	s.Logger.InfoF("GRPC Requisition to moveImage started...") // Log when the request starts

	// Attempt to pick the user ID from JWT

	// Validate if the request is valid
	if err := MoveTokenValidate(req); err != nil {
		s.Logger.ErrorF("MoveTokenValidate failed: %v", err) // Log if validation fails
		return &placedImage.MoveImageResponse{}, err
	}

	var sceneModel models.Scene
	var placedImageModel models.PlacedImage

	// Start the database transaction
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		s.Logger.InfoF("Searching if scene exists...") // Log to search for the scene

		// Check if the scene exists
		if err := tx.First(&sceneModel, req.SceneId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Logger.ErrorF("Scene not found with SceneId: %d", req.SceneId) // Log if scene is not found
				return status.Errorf(codes.NotFound, "scene not found")
			}
			s.Logger.ErrorF("Error fetching scene with SceneId: %d, error: %v", req.SceneId, err) // Log other errors during scene fetch
			return status.Errorf(codes.Internal, "failed to check scene")
		}

		s.Logger.InfoF("Searching if placed image exists...") // Log to search for the placed image

		// Check if the placed image exists
		if err := tx.First(&placedImageModel, req.PlacedImageId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Logger.ErrorF("Placed image not found with PlacedImageId: %d", req.PlacedImageId) // Log if placed image is not found
				return status.Errorf(codes.NotFound, "image not found")
			}
			s.Logger.ErrorF("Error fetching placed image with PlacedImageId: %d, error: %v", req.PlacedImageId, err) // Log other errors during placed image fetch
			return status.Errorf(codes.Internal, "failed to check token")
		}

		// Verify if the user is a master or has permission to edit the object
		s.Logger.InfoF("Verifying if user is a master for SceneId: %d", sceneModel.TableID)
		errMaster := utils.CheckUserIsMaster(ctx, tx, sceneModel.TableID)
		if errMaster == nil {
			st, ok := status.FromError(errMaster)
			if !ok || st.Code() != codes.PermissionDenied {
				s.Logger.ErrorF("User does not have master permissions for SceneId: %d", sceneModel.TableID) // Log if user doesn't have master permissions
				return errMaster
			}
		} else {
			s.Logger.InfoF("Verifying if user can edit PlacedImageId: %d", req.PlacedImageId)
			errEdit := utils.CheckUserCanEditTokenObject(ctx, tx, uint(req.PlacedImageId))
			if errEdit != nil {
				s.Logger.ErrorF("User does not have permission to edit PlacedImageId: %d", req.PlacedImageId) // Log if user doesn't have permission to edit
				return errEdit
			}
		}

		// Update the placed image position
		updateData := models.PlacedImage{
			PosX: int32(req.PosX),
			PosY: int32(req.PosY),
		}

		// Perform the update operation
		if err := tx.Model(&placedImageModel).Updates(updateData).Error; err != nil {
			s.Logger.ErrorF("Error updating placed image position: %v", err) // Log if update fails
			return status.Errorf(codes.Internal, "failed to move image")
		}

		s.Logger.InfoF("Updating image position...") // Log after updating the position

		return nil
	})
	if err != nil {
		s.Logger.ErrorF("Error during transaction: %v", err) // Log if transaction fails
		return nil, err
	}

	// Publish the event after moving the image
	s.Logger.InfoF("Publishing placedImage move event for SceneId: %d", req.SceneId)
	event := events.NewPlacedImageMovedEvent(req.SceneId, req.PlacedImageId, int32(req.PosX), int32(req.PosY))
	s.Broker.Publish(pubSubSyncConst.SceneSync, req.SceneId, event)

	// Prepare the response
	responseProto := &placedImage.MoveImageResponse{
		Success: true,
	}

	return responseProto, nil
}
