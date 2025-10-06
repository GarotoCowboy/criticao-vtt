package placedToken

import (
	"context"
	"errors"

	"github.com/GarotoCowboy/vttProject/api/grpc/events"
	"github.com/GarotoCowboy/vttProject/api/grpc/pb/placedToken"
	"github.com/GarotoCowboy/vttProject/api/models"
	"github.com/GarotoCowboy/vttProject/api/models/consts/pubSubSyncConst"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func (s *PlacedTokenService) CreatePlacedToken(ctx context.Context, req *placedToken.CreatePlacedTokenRequest) (*placedToken.CreatePlacedTokenResponse, error) {

	s.Logger.InfoF("GRPC Requisition to Create placedToken started...")
	if err := Validate(req); err != nil {
		return &placedToken.CreatePlacedTokenResponse{}, err
	}

	var sceneModel models.Scene
	var tokenModel models.Token

	s.Logger.InfoF("Searching if scene exists...")

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Where("id = ?", req.SceneId).First(&sceneModel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return status.Errorf(codes.NotFound, "scene not found")
			}
			return status.Errorf(codes.Internal, "internal error")
		}

		s.Logger.InfoF("Searching if token exists...")

		if err := tx.WithContext(ctx).Where("id = ?", req.TokenId).First(&tokenModel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return status.Errorf(codes.NotFound, "token not found")
			}
			return status.Errorf(codes.Internal, "internal error")
		}
		return nil
	})
	if err != nil {
		s.Logger.WarningF("cannot create placed token")
		return nil, status.Errorf(codes.Internal, "cannot create a placed token")
	}

	s.Logger.InfoF("creating a model and save placedToken...")
	var placedTokenModel = models.PlacedToken{
		SceneID: sceneModel.ID,
		TokenID: tokenModel.ID,
		PosY:    req.PosY,
		PosX:    req.PosX,
	}

	if err := s.DB.Create(&placedTokenModel).Error; err != nil {
		return &placedToken.CreatePlacedTokenResponse{}, status.Errorf(codes.Canceled, "cannot create placedToken")
	}

	response := &placedToken.PlacedToken{
		SceneId: uint64(sceneModel.ID),
		TokenId: uint64(tokenModel.ID),
		PosY:    req.PosY,
		PosX:    req.PosX,
	}

	s.Logger.InfoF("GRPC Requisition to Create Scene finished...")

	event := events.NewPlacedTokenCreatedEvent(response)

	s.Broker.Publish(pubSubSyncConst.SceneSync, req.SceneId, event)

	return &placedToken.CreatePlacedTokenResponse{
		PlacedToken: response,
	}, nil

}
func (s *PlacedTokenService) EditPlacedToken(ctx context.Context, req *placedToken.EditPlacedTokenRequest) (*placedToken.EditPlacedTokenResponse, error) {
	s.Logger.InfoF("GRPC Requisition to Edit placedToken started...")
	updatesMap, err := ValidateAndBuildUpdateMap(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request body: %v", err.Error())
	}
	var placedTokenModel models.PlacedToken
	var sceneModel models.Scene

	s.Logger.InfoF("Searching if scene and placedToken exists...")

	err = s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", req.PlacedToken.SceneId).First(&sceneModel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return status.Errorf(codes.NotFound, "scene not found")
			}
			return status.Errorf(codes.Internal, "internal error")
		}
		if err := tx.Where("id = ?", req.PlacedToken.PlacedTokenId).First(&placedTokenModel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return status.Errorf(codes.NotFound, "placedtoken not found")
			}
			return status.Errorf(codes.Internal, "internal error")
		}

		if err := tx.Model(&placedTokenModel).Updates(updatesMap).Error; err != nil {
			return status.Errorf(codes.Internal, "internal error %v", err.Error())
		}
		return nil
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot update scene or placedToken")
	}

	response := placedToken.PlacedToken{
		SceneId:       uint64(sceneModel.ID),
		TokenId:       req.PlacedToken.TokenId,
		PlacedTokenId: uint64(placedTokenModel.ID),
		Size:          req.PlacedToken.Size,
		Layer:         req.PlacedToken.Layer,
		Rotation:      req.PlacedToken.Rotation,
		UpdatedAt:     timestamppb.New(placedTokenModel.UpdatedAt),
	}

	event := events.NewPlacedTokenUpdatedEvent(&response)
	s.Broker.Publish(pubSubSyncConst.SceneSync, uint64(sceneModel.ID), event)

	return &placedToken.EditPlacedTokenResponse{
		PlacedToken: &response,
	}, nil
}
func (s *PlacedTokenService) DeletePlacedToken(ctx context.Context, req *placedToken.DeletePlacedTokenRequest) (*placedToken.DeletePlacedTokenResponse, error) {
	s.Logger.InfoF("GRPC Requisition to Delete placedToken started...")

	if req.GetSceneId() == 0 && req.GetPlacedTokenId() == 0 {
		return &placedToken.DeletePlacedTokenResponse{}, status.Error(codes.InvalidArgument, "scene id and placed_token id are required")
	}

	var placedTokenModel models.PlacedToken
	var sceneModel models.Scene

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		s.Logger.InfoF("Searching if scene exists...")

		if err := tx.Where("id = ?", req.SceneId).First(&sceneModel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return status.Errorf(codes.NotFound, "scene not found")
			}
			return status.Errorf(codes.Internal, "internal error")
		}
		s.Logger.InfoF("Searching if placed token exists...")
		if err := tx.Where("id = ?", req.PlacedTokenId).First(&placedTokenModel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return status.Errorf(codes.NotFound, "placed token not found")
			}
			return status.Errorf(codes.Internal, "internal error")
		}
		if err := tx.Delete(&placedTokenModel).Error; err != nil {
			return status.Errorf(codes.Internal, "cannot delete placedToken")
		}
		return nil
	})
	if err != nil {
		return &placedToken.DeletePlacedTokenResponse{}, err
	}

	event := events.NewPlacedTokenDeletedEvent(sceneModel.ID, placedTokenModel.ID)
	s.Broker.Publish(pubSubSyncConst.SceneSync, req.SceneId, event)
	s.Logger.InfoF("GRPC Requisition to Delete Scene finished...")
	return &placedToken.DeletePlacedTokenResponse{
		Empty: &emptypb.Empty{},
	}, nil
}
func (s *PlacedTokenService) ListAllTokensOnScene(ctx context.Context, req *placedToken.ListAllTokensRequest) (*placedToken.ListAllTokensResponse, error) {
	s.Logger.InfoF("GRPC Requisition to ListAllTokens started...")

	if req.SceneId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "scene id is required")
	}

	var placedTokenModels []models.PlacedToken

	if err := s.DB.WithContext(ctx).Preload("Token").Where("scene_id = ?", req.SceneId).
		Find(&placedTokenModels).Error; err != nil {
		s.Logger.InfoF("Error fetching placed tokens for the scene")
	}
	s.Logger.InfoF("Found %d tokens for scene %d", len(placedTokenModels), req.SceneId)

	responseTokens := make([]*placedToken.PlacedToken, 0, len(placedTokenModels))

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
		})
	}

	return &placedToken.ListAllTokensResponse{
		PlacedToken: responseTokens,
	}, nil
}
func (s *PlacedTokenService) MoveToken(ctx context.Context, req *placedToken.MoveTokenRequest) (*placedToken.MoveTokenResponse, error) {
	s.Logger.InfoF("GRPC Requisition to moveToken started...")

	if err := MoveTokenValidate(req); err != nil {
		return &placedToken.MoveTokenResponse{}, err
	}

	var sceneModel models.Scene
	var placedTokenModel models.PlacedToken

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		s.Logger.InfoF("Searching if scene exists...")

		if err := tx.First(&sceneModel, req.SceneId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return status.Errorf(codes.NotFound, "scene not found")
			}
			return status.Errorf(codes.Internal, "failed to check scene")
		}
		s.Logger.InfoF("searching if placed Token exists...")
		if err := tx.First(&placedTokenModel, req.PlacedTokenId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return status.Errorf(codes.NotFound, "token not found")
			}
			return status.Errorf(codes.Internal, "failed to check token")
		}

		updateData := models.PlacedToken{
			PosX: int32(req.PosX),
			PosY: int32(req.PosY),
		}

		if err := tx.Model(&placedTokenModel).Updates(updateData).Error; err != nil {
			return status.Errorf(codes.Internal, "failed to move token")
		}

		s.Logger.InfoF("updating token position...")

		return nil
	})
	if err != nil {
		return nil, err
	}

	s.Logger.InfoF("publishing placedToken event for scene %d", req.SceneId)

	event := events.NewPlacedTokenMovedEvent(req.SceneId, req.PlacedTokenId, int32(req.PosX), int32(req.PosY))

	s.Broker.Publish(pubSubSyncConst.SceneSync, req.SceneId, event)
	responseProto := &placedToken.MoveTokenResponse{
		Sucess: true,
	}

	return responseProto, nil
}
