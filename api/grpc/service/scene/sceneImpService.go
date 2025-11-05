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
	"github.com/GarotoCowboy/vttProject/api/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"
)

func (s *SceneService) CreateScene(ctx context.Context, req *scene.CreateSceneRequest) (*scene.CreateSceneResponse, error) {
	s.Logger.InfoF("GRPC Requisition to Create Scene started...")
	// Validate the incoming request payload
	if err := Validate(req); err != nil {
		s.Logger.ErrorF("Invalid create scene request: %v", err)
		return &scene.CreateSceneResponse{}, err
	}

	var tableModel models.Table

	s.Logger.InfoF("Searching if table exists with ID: %d", req.GetTableId())

	// 1. Fetch the table
	if err := s.DB.WithContext(ctx).Where("id = ?", req.GetTableId()).First(&tableModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.ErrorF("Table not found: %d", req.GetTableId())
			return nil, status.Errorf(codes.NotFound, "tableModel not found")
		}
		s.Logger.ErrorF("Error fetching table %d: %v", req.GetTableId(), err)
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	// 2. MASTER-LEVEL AUTHORIZATION
	// The interceptor already verified the user is a MEMBER.
	// Now we verify they are a MASTER, which is required to CREATE a scene.
	err := utils.CheckUserIsMaster(ctx, s.DB, tableModel.ID)
	if err != nil {
		s.Logger.ErrorF("User does not have master permissions for table %d: %v", tableModel.ID, err)
		return nil, err // Return the permission error
	}

	s.Logger.InfoF("Creating a model and saving Scene...")
	var sceneModel = models.Scene{
		TableID:             tableModel.ID,
		Name:                req.Name,
		GridType:            consts.GridType(req.GridType),
		BackGroundColor:     req.BackgroundColor,
		BackgroundImagePath: req.BackgroundImagePath,
		IsVisible:           false, // Scenes always start invisible
		GridCellDistance:    uint(req.GridCellDistance),
		Width:               uint(req.Width),
		Height:              uint(req.Height),
	}

	// Create the scene record in the database
	if err := s.DB.WithContext(ctx).Create(&sceneModel).Error; err != nil {
		s.Logger.ErrorF("Failed to create Scene on DB: %v", err.Error)
		return &scene.CreateSceneResponse{}, status.Errorf(codes.Internal, "cannot create scene")
	}

	// Build the gRPC response
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
		BackgroundColor:     sceneModel.BackGroundColor, // Corrected (was BackgroundImagePath)
	}

	s.Logger.InfoF("Synchronizing new CreateScene event for table %d", sceneModel.TableID)
	// Publish the event to the table's sync channel
	event := events.NewCreateSceneEvent(response)
	s.Broker.Publish(pubSubSyncConst.TableSync, uint64(sceneModel.TableID), event)

	s.Logger.InfoF("GRPC Requisition to Create Scene finished successfully.")

	return &scene.CreateSceneResponse{Scene: response}, nil
}

func (s SceneService) EditScene(ctx context.Context, req *scene.EditSceneRequest) (*scene.EditSceneResponse, error) {
	s.Logger.InfoF("GRPC Requisition to EditScene started...")

	// 1. Validate the field mask to know which fields to update
	updatesMap, err := ValidateAndBuildUpdateMap(req)
	if err != nil {
		s.Logger.ErrorF("Invalid edit scene request or field mask: %v", err)
		return nil, err
	}

	var sceneModel models.Scene

	s.Logger.InfoF("Searching on DB for scene ID: %d", req.GetScene().GetSceneId())

	if err := s.DB.WithContext(ctx).Where("id = ?", req.GetScene().GetSceneId()).First(&sceneModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.ErrorF("Scene not found: %d", req.GetScene().GetSceneId())
			return nil, status.Errorf(codes.NotFound, "scene id not found")
		}
		s.Logger.ErrorF("Error fetching scene %d: %v", req.GetScene().GetSceneId(), err)
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}

	if sceneModel.TableID != uint(req.GetScene().GetTableId()) {
		s.Logger.ErrorF("Scene %d does not belong to table %d (it belongs to %d)", sceneModel.ID, req.GetScene().GetTableId(), sceneModel.TableID)
		return nil, status.Errorf(codes.NotFound, "scene not found in the specified table")
	}

	// 4. MASTER-LEVEL AUTHORIZATION
	// The interceptor verified the user is a MEMBER.
	// Now we verify they are a MASTER of the table this scene *actually* belongs to.
	err = utils.CheckUserIsMaster(ctx, s.DB, sceneModel.TableID) // Using sceneModel.TableID
	if err != nil {
		s.Logger.ErrorF("User does not have master permissions for table %d: %v", sceneModel.TableID, err)
		return nil, err
	}
	// --- END SECURITY FIX ---

	// 5. Apply the updates
	s.Logger.InfoF("Applying updates to scene %d: %v", sceneModel.ID, updatesMap)
	if err := s.DB.WithContext(ctx).Model(&sceneModel).Updates(updatesMap).Error; err != nil {
		s.Logger.ErrorF("Failed to update scene %d: %v", sceneModel.ID, err)
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}

	s.Logger.Info("GRPC Requisition to EditScene finished...")

	// Build the response
	responseScene := &scene.Scene{
		TableId:             uint64(sceneModel.TableID),
		SceneId:             uint64(sceneModel.ID),
		Name:                sceneModel.Name,
		Width:               uint32(sceneModel.Width),
		Height:              uint32(sceneModel.Height),
		CreatedAt:           timestamppb.New(sceneModel.CreatedAt),
		UpdatedAt:           timestamppb.New(sceneModel.UpdatedAt), // GORM updates this
		IsVisible:           wrapperspb.Bool(sceneModel.IsVisible),
		GridCellDistance:    uint64(sceneModel.GridCellDistance),
		GridType:            scene.GridType(sceneModel.GridType),
		BackgroundImagePath: sceneModel.BackgroundImagePath,
		BackgroundColor:     sceneModel.BackGroundColor, // --- BUG FIX (Typo) ---
	}

	s.Logger.InfoF("Synchronizing new UpdateScene event for table %d", sceneModel.TableID)
	// Publish the update event
	event := events.NewUpdateSceneEvent(responseScene)
	s.Broker.Publish(pubSubSyncConst.TableSync, uint64(sceneModel.TableID), event)

	return &scene.EditSceneResponse{
		Scene: responseScene,
	}, nil
}

func (s SceneService) DeleteScene(ctx context.Context, req *scene.DeleteSceneRequest) (*emptypb.Empty, error) {
	s.Logger.InfoF("GRPC Requisition to Delete Scene started...")

	// 1. Validation
	if req.SceneId == 0 || req.TableId == 0 { // Simplified
		s.Logger.ErrorF("Invalid request: SceneId and TableId are required")
		return nil, status.Errorf(codes.InvalidArgument, "scene_id and table_id cannot be null")
	}

	// 2. MASTER-LEVEL AUTHORIZATION
	// The interceptor verified the user is a MEMBER.
	// Now we verify they are a MASTER, which is required to DELETE.
	err := utils.CheckUserIsMaster(ctx, s.DB, uint(req.TableId))
	if err != nil {
		s.Logger.ErrorF("User does not have master permissions for table %d: %v", req.TableId, err)
		return nil, err
	}

	s.Logger.InfoF("Deleting scene %d from table %d...", req.SceneId, req.TableId)

	// 3. Execution (Safe and Atomic)
	// This WHERE clause ensures the user only deletes the scene if it belongs
	// to the table they were authorized for.
	result := s.DB.WithContext(ctx).Where("id = ? AND table_id = ?", req.SceneId, req.TableId).Delete(&models.Scene{})

	if result.Error != nil {
		s.Logger.ErrorF("Error deleting scene: %v", result.Error)
		return nil, status.Errorf(codes.Internal, "cannot delete scene: %v", result.Error)
	}

	// If RowsAffected is 0, the scene wasn't found (or didn't match the table_id)
	if result.RowsAffected == 0 {
		s.Logger.WarningF("No scene found to delete with id %d and table_id %d", req.SceneId, req.TableId)
		return nil, status.Errorf(codes.NotFound, "scene not found or does not belong to the specified table")
	}

	s.Logger.InfoF("Synchronizing new DeleteScene event for table %d", req.TableId)
	// Publish the delete event
	event := events.NewDeleteSceneEvent(req.TableId, req.SceneId)
	s.Broker.Publish(pubSubSyncConst.TableSync, req.TableId, event)

	s.Logger.InfoF("GRPC Requisition to Delete Scene finished...")

	return &emptypb.Empty{}, nil
}

func (s SceneService) ListAllScenesForTable(ctx context.Context, req *scene.ListAllScenesRequest) (*scene.ListAllScenesResponse, error) {
	s.Logger.InfoF("GRPC Requisition to ListAllScenes on table %d", req.TableId)

	// 1. Validation
	if req.TableId == 0 {
		s.Logger.ErrorF("Invalid request: Table id must be greater than zero")
		return nil, status.Errorf(codes.InvalidArgument, "Table id must be greater than zero")
	}

	// 2. MEMBER-LEVEL AUTHORIZATION
	// No code needed here. The `GrpcTableMemberInterceptor` already
	// handled this. If the code reached this point, the user is a member.

	// 3. Pagination Setup
	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 25 // Default page size
	}
	if pageSize > 100 {
		pageSize = 100 // Max page size
	}

	s.Logger.InfoF("Searching DB for all scene models on table %d...", req.TableId)

	// 4. Query Build
	query := s.DB.Where("table_id = ?", req.TableId).Order("id ASC")
	// Keyset pagination: if a page token (last ID) is provided, fetch items after it
	if req.PageToken != "" {
		lastID, err := strconv.ParseInt(req.PageToken, 10, 64)
		if err != nil {
			s.Logger.ErrorF("Invalid page token: %s", req.PageToken)
			return nil, status.Errorf(codes.InvalidArgument, "invalid page token")
		}
		query = query.Where("id > ?", lastID)
	}

	var scenesModels []*models.Scene
	// Fetch (pageSize + 1) items to check if a next page exists
	if err := query.Limit(pageSize + 1).Find(&scenesModels).Error; err != nil {
		s.Logger.ErrorF("Failed to list Scenes from DB: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to list scenesModels from table")
	}

	var nextPageToken string

	// If we fetched more items than the page size, a next page exists
	if len(scenesModels) > pageSize {
		// The token is the ID of the last item *on this page* (index pageSize-1)
		nextPageToken = fmt.Sprintf("%d", scenesModels[pageSize-1].ID)
		// Trim the slice to only return items for *this page*
		scenesModels = scenesModels[:pageSize]
	}

	var protoScenes []*scene.Scene
	s.Logger.InfoF("Listing %d scenes and creating response...", len(scenesModels))

	// Convert DB models to protobuf messages
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
			BackgroundColor:     dbScene.BackGroundColor, // --- BUG FIX (Typo) ---
			CreatedAt:           timestamppb.New(dbScene.CreatedAt),
			UpdatedAt:           timestamppb.New(dbScene.UpdatedAt),
		})
	}

	// Build the final response
	response := &scene.ListAllScenesResponse{
		Scenes:        protoScenes,
		NextPageToken: nextPageToken,
	}

	s.Logger.InfoF("GRPC Requisition to ListAllScenes finished successfully.")
	return response, nil
}
