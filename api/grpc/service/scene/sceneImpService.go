package scene

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/GarotoCowboy/vttProject/api/grpc/events"
	"github.com/GarotoCowboy/vttProject/api/grpc/pb/scene"
	"github.com/GarotoCowboy/vttProject/api/models"
	"github.com/GarotoCowboy/vttProject/api/models/consts"
	"github.com/GarotoCowboy/vttProject/api/models/consts/pubSubSyncConst"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"
)

func (s *SceneService) CreateScene(ctx context.Context, req *scene.CreateSceneRequest) (*scene.CreateSceneResponse, error) {
	s.Logger.InfoF("GRPC Requisition to Create Scene started...")
	if err := Validate(req); err != nil {
		return &scene.CreateSceneResponse{}, err
	}

	var tableModel models.Table

	s.Logger.InfoF("Searching if tableModel exists...")

	if err := s.DB.WithContext(ctx).Where("id = ?", req.GetTableId()).First(&tableModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "tableModel not found")
		}
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	s.Logger.InfoF("creating a model and save Scene...")
	var sceneModel = models.Scene{
		TableID:             tableModel.ID,
		Name:                req.Name,
		GridType:            consts.GridType(req.GridType),
		BackGroundColor:     req.BackgroundColor,
		BackgroundImagePath: req.BackgroundImagePath,
		IsVisible:           false,
		GridCellDistance:    uint(req.GridCellDistance),
		Width:               uint(req.Width),
		Height:              uint(req.Height),
	}

	if err := s.DB.WithContext(ctx).Create(&sceneModel).Error; err != nil {
		s.Logger.ErrorF("Failed to create Scene on DB: %v", err.Error)
		return &scene.CreateSceneResponse{}, status.Errorf(codes.Internal, "cannot create scene")
	}

	response := &scene.Scene{
		SceneId:             uint64(sceneModel.ID),
		TableId:             uint64(sceneModel.TableID),
		Name:                sceneModel.Name,
		Width:               uint32(sceneModel.Width),
		Height:              uint32(sceneModel.Height),
		IsVisible:           wrapperspb.Bool(sceneModel.IsVisible),
		GridCellDistance:    uint64(sceneModel.GridCellDistance),
		GridType:            scene.GridType(sceneModel.GridType),
		BackgroundImagePath: sceneModel.BackgroundImagePath,
		BackgroundColor:     sceneModel.BackGroundColor,
	}

	s.Logger.InfoF("Sycronizing this new Event")
	event := events.NewCreateSceneEvent(response)

	s.Broker.Publish(pubSubSyncConst.TableSync, uint64(sceneModel.TableID), event)

	s.Logger.InfoF("GRPC Requisition to Create Scene finished...")

	return &scene.CreateSceneResponse{Scene: response}, nil
}
func (s SceneService) EditScene(ctx context.Context, req *scene.EditSceneRequest) (*scene.EditSceneResponse, error) {

	s.Logger.InfoF("GRPC Requisition to EditScene started...")
	//validate the mask
	updatesMap, err := ValidateAndBuildUpdateMap(req)
	if err != nil {
		return nil, err
	}

	//create models
	var sceneModel models.Scene
	var tableModel models.Table

	s.Logger.InfoF("Searching on DB...")

	//search if table exists on DB
	if err := s.DB.WithContext(ctx).Where("id = ?", req.GetScene().GetTableId()).First(&tableModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "invalid table id: %v", "table id not found")
		}
		return nil, status.Errorf(codes.NotFound, "internal error: %v", err.Error())
	}
	//search if scene exits on DB
	if err := s.DB.WithContext(ctx).Where("id = ?", req.GetScene().GetSceneId()).First(&sceneModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "invalid image id: %v", "image id not found")
		}
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}

	if err := s.DB.WithContext(ctx).Model(&sceneModel).Updates(updatesMap).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}

	s.Logger.Info("GRPC Requisition to EditScene finished...")

	responseScene := &scene.Scene{
		TableId:             uint64(sceneModel.TableID),
		SceneId:             uint64(sceneModel.ID),
		Name:                sceneModel.Name,
		Width:               uint32(sceneModel.Width),
		Height:              uint32(sceneModel.Height),
		CreatedAt:           timestamppb.New(sceneModel.CreatedAt),
		UpdatedAt:           timestamppb.New(sceneModel.UpdatedAt),
		IsVisible:           wrapperspb.Bool(sceneModel.IsVisible),
		GridCellDistance:    uint64(sceneModel.GridCellDistance),
		GridType:            scene.GridType(sceneModel.GridType),
		BackgroundImagePath: sceneModel.BackgroundImagePath,
		BackgroundColor:     sceneModel.BackgroundImagePath,
	}

	s.Logger.InfoF("Sycronizing this new Event")
	event := events.NewUpdateSceneEvent(responseScene)

	s.Broker.Publish(pubSubSyncConst.TableSync, uint64(sceneModel.TableID), event)

	return &scene.EditSceneResponse{
		Scene: responseScene,
	}, nil
}

func (s SceneService) DeleteScene(ctx context.Context, req *scene.DeleteSceneRequest) (*emptypb.Empty, error) {
	s.Logger.InfoF("GRPC Requisition to Delete Scene started...")

	if req.SceneId == 0 && req.TableId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "fields cannot be null")
	}

	if req.SceneId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "scene_id cannot be null")
	}

	if req.TableId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "table_id cannot be null")
	}

	s.Logger.InfoF("Deleting scene...")

	result := s.DB.WithContext(ctx).Where("id = ? AND table_id = ?", req.SceneId, req.TableId).Delete(&models.Scene{})

	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "cannot delete scene: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		s.Logger.WarningF("No scene found to delete with id %d and table_id %d", req.SceneId, req.TableId)
		return nil, status.Errorf(codes.NotFound, "scene not found or does not belong to the specified table")
	}

	s.Logger.InfoF("Sycronizing this new Event")
	event := events.NewDeleteSceneEvent(req.TableId, req.SceneId)
	s.Broker.Publish(pubSubSyncConst.TableSync, req.TableId, event)

	s.Logger.InfoF("GRPC Requisition to Delete Scene finished...")

	return &emptypb.Empty{}, nil
}
func (s SceneService) ListAllScenesForTable(ctx context.Context, req *scene.ListAllScenesRequest) (*scene.ListAllScenesResponse, error) {
	s.Logger.InfoF("GRPC Requisiton to ListAllScenes on table %d", req.TableId)

	if req.TableId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Table id must be greater than zero")
	}

	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 25
	}

	if pageSize > 100 {
		pageSize = 100
	}

	s.Logger.InfoF("searching on DB all scenesModels...")

	query := s.DB.Where("table_id = ?", req.TableId).Order("id ASC")

	if req.PageToken != "" {
		lastID, err := strconv.ParseInt(req.PageToken, 10, 64)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid page token")
		}
		query = query.Where("id > ?", lastID)
	}

	var scenesModels []*models.Scene
	if err := query.Limit(pageSize + 1).Find(&scenesModels).Error; err != nil {
		s.Logger.ErrorF("failed to list Scenes : %v", err)
		return nil, status.Errorf(codes.Internal, "failed to list scenesModels from table")
	}

	var nextPageToken string

	if len(scenesModels) > pageSize {
		nextPageToken = fmt.Sprintf("%d", scenesModels[len(scenesModels)-1].ID)
		scenesModels = scenesModels[:pageSize]
	}

	var protoScenes []*scene.Scene

	s.Logger.InfoF("Listing all scenesModels and creating response...")

	for _, dbScene := range scenesModels {
		protoScenes = append(protoScenes, &scene.Scene{
			SceneId:             uint64(dbScene.ID),
			TableId:             uint64(dbScene.TableID),
			Name:                dbScene.Name,
			Width:               uint32(dbScene.Width),
			Height:              uint32(dbScene.Height),
			IsVisible:           wrapperspb.Bool(dbScene.IsVisible),
			GridCellDistance:    uint64(dbScene.GridCellDistance),
			GridType:            scene.GridType(dbScene.GridType),
			BackgroundImagePath: dbScene.BackgroundImagePath,
			BackgroundColor:     dbScene.BackgroundImagePath,
			CreatedAt:           timestamppb.New(dbScene.CreatedAt),
			UpdatedAt:           timestamppb.New(dbScene.UpdatedAt),
		})
	}

	response := &scene.ListAllScenesResponse{
		Scene:         protoScenes,
		NextPageToken: nextPageToken,
	}

	return response, nil
}
