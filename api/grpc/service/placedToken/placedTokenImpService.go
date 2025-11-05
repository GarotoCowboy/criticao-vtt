package placedToken

import (
	"context"
	"errors"

	"github.com/GarotoCowboy/vttProject/api/grpc/events"
	"github.com/GarotoCowboy/vttProject/api/grpc/pb/placedToken"
	"github.com/GarotoCowboy/vttProject/api/models"
	"github.com/GarotoCowboy/vttProject/api/models/consts/pubSubSyncConst"
	"github.com/GarotoCowboy/vttProject/api/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func (s *PlacedTokenService) CreatePlacedToken(ctx context.Context, req *placedToken.CreatePlacedTokenRequest) (*placedToken.CreatePlacedTokenResponse, error) {

	s.Logger.InfoF("GRPC Requisition to Create placedToken started...") // Log when the request starts

	// Validate the request body payload
	if err := Validate(req); err != nil {
		s.Logger.ErrorF("Invalid request body: %v", err) // Log if validation fails
		return &placedToken.CreatePlacedTokenResponse{}, err
	}

	// Declare models for database operations
	var sceneModel models.Scene
	var tokenModel models.Token

	s.Logger.InfoF("Searching if scene exists...") // Log to search for the scene

	// Start a database transaction
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		// Check if the scene exists using the provided SceneId
		if err := tx.WithContext(ctx).Where("id = ?", req.SceneId).First(&sceneModel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Logger.ErrorF("Scene with SceneId %d not found", req.SceneId) // Log if scene is not found
				return status.Errorf(codes.NotFound, "scene not found")
			}
			s.Logger.ErrorF("Error fetching scene with SceneId %d: %v", req.SceneId, err) // Log other errors during scene fetch
			return status.Errorf(codes.Internal, "internal error")
		}

		s.Logger.InfoF("Searching if token exists...") // Log to search for the token

		// Check if the token exists using the provided TokenId
		if err := tx.WithContext(ctx).Where("id = ?", req.TokenId).First(&tokenModel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Logger.ErrorF("Token with TokenId %d not found", req.TokenId) // Log if token is not found
				return status.Errorf(codes.NotFound, "token not found")
			}
			s.Logger.ErrorF("Error fetching token with TokenId %d: %v", req.TokenId, err) // Log other errors during token fetch
			return status.Errorf(codes.Internal, "internal error")
		}

		// Check if the user making the request has master permissions for the scene's table
		err := utils.CheckUserIsMaster(ctx, tx, sceneModel.TableID)
		if err != nil {
			s.Logger.ErrorF("User does not have master permissions for SceneId %d", req.SceneId) // Log if user doesn't have master permissions
			return err                                                                           // Return permission error
		}

		// If all checks pass, commit the transaction
		return nil
	})

	// If the transaction failed, log and return an error
	if err != nil {
		s.Logger.WarningF("Cannot create placed token due to transaction error: %v", err) // Log if the creation fails
		return nil, status.Errorf(codes.Internal, "cannot create a placed token")
	}

	s.Logger.InfoF("Creating a model and saving placedToken...") // Log before saving the placed token

	// Create the placed token model struct to be inserted
	var placedTokenModel = models.PlacedToken{
		SceneID: sceneModel.ID, // Set SceneID from the found scene
		TokenID: tokenModel.ID, // Set TokenID from the found token
		PosY:    req.PosY,      // Set position Y from request
		PosX:    req.PosX,      // Set position X from request
	}

	// Save the new placed token record in the database
	if err := s.DB.Create(&placedTokenModel).Error; err != nil {
		s.Logger.ErrorF("Error saving placedToken: %v", err) // Log if save fails
		return &placedToken.CreatePlacedTokenResponse{}, status.Errorf(codes.Canceled, "cannot create placedToken")
	}

	// Prepare the gRPC response protobuf message
	response := &placedToken.PlacedToken{
		SceneId: uint64(sceneModel.ID),
		TokenId: uint64(tokenModel.ID),
		PosY:    req.PosY,
		PosX:    req.PosX,
	}

	s.Logger.InfoF("GRPC Requisition to Create Scene finished...") // Log when the request finishes

	// Publish an event indicating a new placed token was created
	event := events.NewPlacedTokenCreatedEvent(response)
	s.Broker.Publish(pubSubSyncConst.SceneSync, req.SceneId, event)

	// Return the successful response
	return &placedToken.CreatePlacedTokenResponse{
		PlacedToken: response,
	}, nil
}

func (s *PlacedTokenService) EditPlacedToken(ctx context.Context, req *placedToken.EditPlacedTokenRequest) (*placedToken.EditPlacedTokenResponse, error) {
	s.Logger.InfoF("GRPC Requisition to Edit placedToken started...") // Log when the request starts

	// Validate the request and build a map of fields to update
	updatesMap, err := ValidateAndBuildUpdateMap(req)
	if err != nil {
		s.Logger.ErrorF("Error validating and building update map: %v", err) // Log if update map creation fails
		return nil, status.Errorf(codes.InvalidArgument, "invalid request body: %v", err.Error())
	}

	// Declare models for database operations
	var placedTokenModel models.PlacedToken
	var sceneModel models.Scene

	s.Logger.InfoF("Searching if scene and placedToken exists...") // Log to search for the scene and placed token

	// Start a database transaction
	err = s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Check if the scene exists
		s.Logger.InfoF("Checking if scene exists with SceneId: %d", req.PlacedToken.SceneId)
		if err := tx.Where("id = ?", req.PlacedToken.SceneId).First(&sceneModel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Logger.ErrorF("Scene with SceneId %d not found", req.PlacedToken.SceneId) // Log if scene is not found
				return status.Errorf(codes.NotFound, "scene not found")
			}
			s.Logger.ErrorF("Error fetching scene with SceneId %d: %v", req.PlacedToken.SceneId, err) // Log other errors during scene fetch
			return status.Errorf(codes.Internal, "internal error")
		}

		// Check if the placed token exists
		s.Logger.InfoF("Checking if placedToken exists with PlacedTokenId: %d", req.PlacedToken.PlacedTokenId)
		if err := tx.Where("id = ?", req.PlacedToken.PlacedTokenId).First(&placedTokenModel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Logger.ErrorF("PlacedToken with PlacedTokenId %d not found", req.PlacedToken.PlacedTokenId) // Log if placed token is not found
				return status.Errorf(codes.NotFound, "placedtoken not found")
			}
			s.Logger.ErrorF("Error fetching placedToken with PlacedTokenId %d: %v", req.PlacedToken.PlacedTokenId, err) // Log other errors during placed token fetch
			return status.Errorf(codes.Internal, "internal error")
		}

		// Verify if user is a master OR has permission to edit the specific token object
		s.Logger.InfoF("Verifying if user is a master for SceneId: %d", sceneModel.TableID)
		errMaster := utils.CheckUserIsMaster(ctx, tx, sceneModel.TableID)

		// Check if the master check failed with a permission error
		if errMaster != nil {
			st, ok := status.FromError(errMaster)
			if !ok || st.Code() != codes.PermissionDenied {
				// If it's not a permission error, it's an internal error
				s.Logger.ErrorF("Error checking master permissions for SceneId %d: %v", sceneModel.TableID, errMaster)
				return errMaster // Return the internal error
			}

			// If it IS a permission error (user is not master), check for specific token edit rights
			s.Logger.InfoF("User is not master, verifying if user can edit PlacedTokenId: %d", req.PlacedToken.PlacedTokenId)
			errEdit := utils.CheckUserCanEditTokenObject(ctx, tx, uint(req.PlacedToken.PlacedTokenId))
			if errEdit != nil {
				s.Logger.ErrorF("User does not have permission to edit PlacedTokenId: %d", req.PlacedToken.PlacedTokenId) // Log if user doesn't have permission to edit
				return errEdit                                                                                            // Return the specific edit permission error
			}
		}
		// If errMaster was nil, user is master and check passes

		// Update placedToken with the new values from the map
		s.Logger.InfoF("Updating PlacedToken with new values: %v", updatesMap)
		if err := tx.Model(&placedTokenModel).Updates(updatesMap).Error; err != nil {
			s.Logger.ErrorF("Error updating PlacedToken with PlacedTokenId %d: %v", placedTokenModel.ID, err) // Log if update fails
			return status.Errorf(codes.Internal, "internal error %v", err.Error())
		}

		// Commit the transaction
		return nil
	})

	// If the transaction failed, log and return error
	if err != nil {
		s.Logger.ErrorF("Transaction failed: %v", err) // Log if transaction fails
		return nil, status.Errorf(codes.Internal, "cannot update scene or placedToken")
	}

	// Prepare the gRPC response
	response := placedToken.PlacedToken{
		SceneId:       uint64(sceneModel.ID),
		TokenId:       req.PlacedToken.TokenId,
		PlacedTokenId: uint64(placedTokenModel.ID),
		Size:          req.PlacedToken.Size,
		Layer:         req.PlacedToken.Layer,
		Rotation:      req.PlacedToken.Rotation,
		UpdatedAt:     timestamppb.New(placedTokenModel.UpdatedAt), // Send updated timestamp
	}

	// Publish the event after update
	event := events.NewPlacedTokenUpdatedEvent(&response)
	s.Broker.Publish(pubSubSyncConst.SceneSync, uint64(sceneModel.ID), event)

	s.Logger.InfoF("PlacedToken updated successfully with PlacedTokenId: %d", placedTokenModel.ID) // Log when update is successful
	return &placedToken.EditPlacedTokenResponse{
		PlacedToken: &response,
	}, nil
}

func (s *PlacedTokenService) DeletePlacedToken(ctx context.Context, req *placedToken.DeletePlacedTokenRequest) (*placedToken.DeletePlacedTokenResponse, error) {
	s.Logger.InfoF("GRPC Requisition to Delete placedToken started...")

	// Validate required IDs
	if req.GetSceneId() == 0 && req.GetPlacedTokenId() == 0 {
		s.Logger.ErrorF("Invalid request: SceneId and PlacedTokenId are required")
		return &placedToken.DeletePlacedTokenResponse{}, status.Error(codes.InvalidArgument, "scene id and placed_token id are required")
	}

	// Declare models
	var placedTokenModel models.PlacedToken
	var sceneModel models.Scene

	// Start database transaction
	err := s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		s.Logger.InfoF("Searching if scene exists with SceneId: %d", req.SceneId)
		// Check if scene exists
		if err := tx.Where("id = ?", req.SceneId).First(&sceneModel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Logger.ErrorF("Scene not found with SceneId: %d", req.SceneId)
				return status.Errorf(codes.NotFound, "scene not found")
			}
			s.Logger.ErrorF("Error fetching scene: %v", err)
			return status.Errorf(codes.Internal, "internal error")
		}

		s.Logger.InfoF("Searching if placed token exists with PlacedTokenId: %d", req.PlacedTokenId)
		// Check if placed token exists
		if err := tx.Where("id = ?", req.PlacedTokenId).First(&placedTokenModel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Logger.ErrorF("PlacedToken not found with PlacedTokenId: %d", req.PlacedTokenId)
				return status.Errorf(codes.NotFound, "placed token not found")
			}
			s.Logger.ErrorF("Error fetching placed token: %v", err)
			return status.Errorf(codes.Internal, "internal error")
		}

		// Verify if user is a master OR has permission to edit the object
		s.Logger.InfoF("Verifying if user is a master for SceneId: %d", sceneModel.TableID)
		errMaster := utils.CheckUserIsMaster(ctx, tx, sceneModel.TableID)

		if errMaster != nil { // If master check failed
			st, ok := status.FromError(errMaster)
			if !ok || st.Code() != codes.PermissionDenied {
				s.Logger.ErrorF("Error checking master permissions for SceneId %d: %v", sceneModel.TableID, errMaster)
				return errMaster // Return internal error
			}

			// User is not master, check specific edit rights
			s.Logger.InfoF("User is not master, verifying if user can edit PlacedTokenId: %d", req.PlacedTokenId)
			errEdit := utils.CheckUserCanEditTokenObject(ctx, tx, uint(req.PlacedTokenId))
			if errEdit != nil {
				s.Logger.ErrorF("User does not have permission to edit/delete PlacedTokenId: %d", req.PlacedTokenId)
				return errEdit // Return permission error
			}
		}
		// User is master or has specific rights, proceed with delete

		s.Logger.InfoF("Deleting PlacedTokenId: %d", placedTokenModel.ID)
		// Delete the placed token
		if err := tx.Delete(&placedTokenModel).Error; err != nil {
			s.Logger.ErrorF("Error deleting placed token: %v", err)
			return status.Errorf(codes.Internal, "cannot delete placedToken")
		}

		// Commit transaction
		return nil
	})

	// Check transaction error
	if err != nil {
		s.Logger.ErrorF("Transaction failed: %v", err)
		return &placedToken.DeletePlacedTokenResponse{}, err
	}

	// Publish delete event
	event := events.NewPlacedTokenDeletedEvent(sceneModel.ID, placedTokenModel.ID)
	s.Broker.Publish(pubSubSyncConst.SceneSync, req.SceneId, event)

	s.Logger.InfoF("GRPC Requisition to Delete Scene finished successfully.")
	return &placedToken.DeletePlacedTokenResponse{
		Empty: &emptypb.Empty{},
	}, nil
}
func (s *PlacedTokenService) ListAllTokensOnScene(ctx context.Context, req *placedToken.ListAllTokensRequest) (*placedToken.ListAllTokensResponse, error) {
	s.Logger.InfoF("GRPC Requisition to ListAllTokens started...")

	// Validate input
	if req.SceneId == 0 {
		s.Logger.ErrorF("Invalid request: SceneId is required")
		return nil, status.Errorf(codes.InvalidArgument, "scene id is required")
	}

	var placedTokenModels []models.PlacedToken

	// Fetch all placed tokens for the given scene, preloading associated Token data
	if err := s.DB.WithContext(ctx).Preload("Token").Where("scene_id = ?", req.SceneId).
		Find(&placedTokenModels).Error; err != nil {
		s.Logger.ErrorF("Error fetching placed tokens for the scene %d: %v", req.SceneId, err)
		return nil, status.Errorf(codes.Internal, "failed to fetch tokens")
	}

	s.Logger.InfoF("Found %d tokens for scene %d", len(placedTokenModels), req.SceneId)

	// Prepare response slice
	responseTokens := make([]*placedToken.PlacedToken, 0, len(placedTokenModels))

	// Iterate over models and convert to protobuf response format
	for _, model := range placedTokenModels {
		responseTokens = append(responseTokens, &placedToken.PlacedToken{
			SceneId:       uint64(model.SceneID),
			PlacedTokenId: uint64(model.ID),
			TokenId:       uint64(model.TokenID),
			Size:          model.Size,
			Layer:         placedToken.LayerType(model.LayerType),
			PosX:          model.PosX,
			PosY:          model.PosY,
			Rotation:      int32(model.Rotation),
			CreatedAt:     timestamppb.New(model.CreatedAt),
			UpdatedAt:     timestamppb.New(model.UpdatedAt),
			// Note: Preloaded Token data (model.Token) isn't being mapped here,
			// you might want to add fields from model.Token to the PlacedToken proto if needed.
		})
	}

	s.Logger.InfoF("Successfully listed all tokens for scene %d.", req.SceneId)
	// Return the list of tokens
	return &placedToken.ListAllTokensResponse{
		PlacedTokens: responseTokens,
	}, nil
}
func (s *PlacedTokenService) MoveToken(ctx context.Context, req *placedToken.MoveTokenRequest) (*placedToken.MoveTokenResponse, error) {
	s.Logger.InfoF("GRPC Requisition to moveToken started...")

	var sceneModel models.Scene
	var placedTokenModel models.PlacedToken

	// Start database transaction
	err := s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		s.Logger.InfoF("Searching if scene exists with SceneId: %d", req.SceneId)
		// Check if scene exists
		if err := tx.First(&sceneModel, req.SceneId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Logger.ErrorF("Scene not found: %d", req.SceneId)
				return status.Errorf(codes.NotFound, "scene not found")
			}
			s.Logger.ErrorF("Failed to check scene: %v", err)
			return status.Errorf(codes.Internal, "failed to check scene")
		}

		s.Logger.InfoF("searching if placed Token exists with PlacedTokenId: %d", req.PlacedTokenId)
		// Check if placed token exists
		if err := tx.First(&placedTokenModel, req.PlacedTokenId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Logger.ErrorF("PlacedToken not found: %d", req.PlacedTokenId)
				return status.Errorf(codes.NotFound, "token not found")
			}
			s.Logger.ErrorF("Failed to check token: %v", err)
			return status.Errorf(codes.Internal, "failed to check token")
		}

		// Verify if user is a master OR has permission to edit the object
		s.Logger.InfoF("Verifying if user is a master for SceneId: %d (TableID: %d)", req.SceneId, sceneModel.TableID)
		errMaster := utils.CheckUserIsMaster(ctx, tx, sceneModel.TableID)

		if errMaster != nil { // If master check failed
			st, ok := status.FromError(errMaster)
			if !ok || st.Code() != codes.PermissionDenied {
				s.Logger.ErrorF("Error checking master permissions for SceneId %d: %v", req.SceneId, errMaster)
				return errMaster // Return internal error
			}

			// User is not master, check specific edit rights
			s.Logger.InfoF("User is not master, verifying if user can edit PlacedTokenId: %d", req.PlacedTokenId)
			errEdit := utils.CheckUserCanEditTokenObject(ctx, tx, uint(req.PlacedTokenId))
			if errEdit != nil {
				s.Logger.ErrorF("User does not have permission to edit/move PlacedTokenId: %d", req.PlacedTokenId)
				return errEdit // Return permission error
			}
		}
		// User is master or has specific rights, proceed with move

		// Prepare update data with new positions
		updateData := models.PlacedToken{
			PosX: int32(req.PosX),
			PosY: int32(req.PosY),
		}

		s.Logger.InfoF("updating token position for PlacedTokenId: %d...", req.PlacedTokenId)
		// Update the token's position
		if err := tx.Model(&placedTokenModel).Updates(updateData).Error; err != nil {
			s.Logger.ErrorF("Failed to move token in database: %v", err)
			return status.Errorf(codes.Internal, "failed to move token")
		}

		// Commit transaction
		return nil
	})

	// Check transaction error
	if err != nil {
		s.Logger.ErrorF("Transaction failed: %v", err)
		return nil, err
	}

	s.Logger.InfoF("publishing placedToken moved event for scene %d", req.SceneId)

	// Publish move event
	event := events.NewPlacedTokenMovedEvent(req.SceneId, req.PlacedTokenId, int32(req.PosX), int32(req.PosY))
	s.Broker.Publish(pubSubSyncConst.SceneSync, req.SceneId, event)

	// Prepare successful response
	responseProto := &placedToken.MoveTokenResponse{
		Success: true,
	}

	s.Logger.InfoF("GRPC Requisition to moveToken finished successfully.")
	return responseProto, nil
}
