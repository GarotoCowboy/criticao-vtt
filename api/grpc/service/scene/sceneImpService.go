package scene

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/GarotoCowboy/vttProject/api/grpc/proto/scene/pb"
	"github.com/GarotoCowboy/vttProject/api/models"
	"github.com/GarotoCowboy/vttProject/api/models/consts"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"
)

func (s *SceneService) CreateScene(ctx context.Context, req *pb.CreateSceneRequest) (*pb.CreateSceneResponse, error) {
	s.Logger.InfoF("GRPC Requisition to Create Scene started...")
	if err := Validate(req); err != nil {
		return &pb.CreateSceneResponse{}, err
	}

	var table models.Table

	s.Logger.InfoF("Searching if table exists...")

	if err := s.DB.Where("id = ?", req.GetTableId()).First(&table).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "table not found")
		}
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	s.Logger.InfoF("creating a model and save Scene...")
	var scene = models.Scene{
		TableID:             table.ID,
		Name:                req.Name,
		GridType:            consts.GridType(req.GridType),
		BackGroundColor:     req.BackgroundColor,
		BackgroundImagePath: req.BackgroundImagePath,
		IsVisible:           false,
		GridCellDistance:    uint(req.GridCellDistance),
		Width:               uint(req.Width),
		Height:              uint(req.Height),
	}

	if err := s.DB.Create(&scene); err != nil {
		return &pb.CreateSceneResponse{}, status.Errorf(codes.Canceled, "cannot create scene")
	}

	response := &pb.Scene{
		SceneId:             uint64(scene.ID),
		TableId:             uint64(scene.TableID),
		Name:                scene.Name,
		Width:               uint32(scene.Width),
		Height:              uint32(scene.Height),
		IsVisible:           wrapperspb.Bool(scene.IsVisible),
		GridCellDistance:    uint64(scene.GridCellDistance),
		GridType:            pb.GridType(scene.GridType),
		BackgroundImagePath: scene.BackgroundImagePath,
		BackgroundColor:     scene.BackgroundImagePath,
	}

	s.Logger.InfoF("GRPC Requisition to Create Scene finished...")

	return &pb.CreateSceneResponse{Scene: response}, nil
}
func (s SceneService) EditScene(ctx context.Context, req *pb.EditSceneRequest) (*pb.EditSceneResponse, error) {

	s.Logger.InfoF("GRPC Requisition to EditScene started...")
	//validate the mask
	updatesMap, err := ValidateAndBuildUpdateMap(req)
	if err != nil {
		return nil, err
	}

	//create models
	var scene models.Scene
	var table models.Table

	s.Logger.InfoF("Searching on DB...")

	//search if table exists on DB
	if err := s.DB.Where("id = ?", req.GetScene().GetTableId()).First(&table).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "invalid table id: %v", "table id not found")
		}
		return nil, status.Errorf(codes.NotFound, "internal error: %v", err.Error())
	}
	//search if scene exits on DB
	if err := s.DB.Where("id = ?", req.GetScene().GetSceneId()).First(&scene).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "invalid image id: %v", "image id not found")
		}
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}

	if err := s.DB.Model(&scene).Updates(updatesMap).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}

	s.Logger.Info("GRPC Requisition to EditScene finished...")

	responseScene := &pb.Scene{
		TableId:             uint64(scene.TableID),
		SceneId:             uint64(scene.ID),
		Name:                scene.Name,
		Width:               uint32(scene.Width),
		Height:              uint32(scene.Height),
		CreatedAt:           timestamppb.New(scene.CreatedAt),
		UpdatedAt:           timestamppb.New(scene.UpdatedAt),
		IsVisible:           wrapperspb.Bool(scene.IsVisible),
		GridCellDistance:    uint64(scene.GridCellDistance),
		GridType:            pb.GridType(scene.GridType),
		BackgroundImagePath: scene.BackgroundImagePath,
		BackgroundColor:     scene.BackgroundImagePath,
	}

	return &pb.EditSceneResponse{
		Scene: responseScene,
	}, nil
}

func (s SceneService) DeleteScene(ctx context.Context, req *pb.DeleteSceneRequest) (*emptypb.Empty, error) {
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

	var table models.Table

	if err := s.DB.Where("id = ?", req.TableId).First(&table).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "table not found")
		}
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	s.Logger.InfoF("Deleting scene...")

	if err := s.DB.Where("id = ? AND table_id = ?", req.SceneId, table.ID).Delete(&models.Scene{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "scene not found")
		}
		return nil, status.Errorf(codes.Internal, "cannot delete scene")
	}

	s.Logger.InfoF("GRPC Requisition to Delete Scene finished...")

	return &emptypb.Empty{}, status.Errorf(codes.Unimplemented, "method DeleteScene not implemented")
}
func (s SceneService) ListAllScenesForTable(ctx context.Context, req *pb.ListAllScenesRequest) (*pb.ListAllScenesResponse, error) {
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

	s.Logger.InfoF("searching on DB all scenes...")

	query := s.DB.Where("table_id = ?", req.TableId).Order("id ASC")

	if req.PageToken != "" {
		lastID, err := strconv.ParseInt(req.PageToken, 10, 64)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid page token")
		}
		query = query.Where("id > ?", lastID)
	}

	var scenes []*models.Scene
	if err := query.Limit(pageSize + 1).Find(&scenes).Error; err != nil {
		s.Logger.ErrorF("failed to list Scenes : %v", err)
		return nil, status.Errorf(codes.Internal, "failed to list scenes from table")
	}

	var nextPageToken string

	if len(scenes) > pageSize {
		nextPageToken = fmt.Sprintf("%d", scenes[len(scenes)-1].ID)
		scenes = scenes[:pageSize]
	}

	var protoScenes []*pb.Scene

	s.Logger.InfoF("Listing all scenes and creating response...")

	for _, dbScene := range scenes {
		protoScenes = append(protoScenes, &pb.Scene{
			SceneId:             uint64(dbScene.ID),
			TableId:             uint64(dbScene.TableID),
			Name:                dbScene.Name,
			Width:               uint32(dbScene.Width),
			Height:              uint32(dbScene.Height),
			IsVisible:           wrapperspb.Bool(dbScene.IsVisible),
			GridCellDistance:    uint64(dbScene.GridCellDistance),
			GridType:            pb.GridType(dbScene.GridType),
			BackgroundImagePath: dbScene.BackgroundImagePath,
			BackgroundColor:     dbScene.BackgroundImagePath,
			CreatedAt:           timestamppb.New(dbScene.CreatedAt),
			UpdatedAt:           timestamppb.New(dbScene.UpdatedAt),
		})
	}

	response := &pb.ListAllScenesResponse{
		Scene:         protoScenes,
		NextPageToken: nextPageToken,
	}

	return response, nil
}
