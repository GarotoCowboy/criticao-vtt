package permission

import (
	"context"
	"errors"

	"github.com/GarotoCowboy/vttProject/api/grpc/events"
	"github.com/GarotoCowboy/vttProject/api/grpc/pb/permission"
	syncBroker "github.com/GarotoCowboy/vttProject/api/grpc/pb/sync"
	"github.com/GarotoCowboy/vttProject/api/models"
	"github.com/GarotoCowboy/vttProject/api/models/consts"
	"github.com/GarotoCowboy/vttProject/api/models/consts/pubSubSyncConst"
	"github.com/GarotoCowboy/vttProject/api/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

func (s *PermissionService) UpdatePlacedObjectAccess(ctx context.Context, req *permission.UpdatePlacedObjectAccessRequest) (*emptypb.Empty, error) {
	s.Logger.InfoF("GRPC Requisition to Update a PlacedObject access permissions started...")

	// We need a valid object ID to update
	if req.ObjectId == 0 {
		s.Logger.ErrorF("Invalid request: ObjectId must be greater than 0")
		return nil, status.Error(codes.InvalidArgument, "objectId must be higher than 0")
	}

	// --- Variables to store context retrieved from the object ---
	var sceneID uint
	var tableID uint // To check master permissions

	// --- Object models ---
	var placedTokenModel models.PlacedToken
	var placedImageModel models.PlacedImage

	s.Logger.InfoF("Selecting the object Type to fetch data...")
	// Switch on the object type to find the object, its scene, and its table
	switch req.GetObjectType() {

	case permission.GameObjectType_TOKEN:
		s.Logger.InfoF("Object Type: Token Selected...")

		// Fetch the placed token and preload its associated Scene
		if err := s.DB.Preload("Scene").Where("id = ?", req.ObjectId).First(&placedTokenModel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Logger.ErrorF("Placed token with ID %d not found", req.ObjectId)
				return nil, status.Errorf(codes.NotFound, "token not found")
			}
			s.Logger.ErrorF("Failed to fetch placed token %d: %v", req.ObjectId, err)
			return nil, status.Errorf(codes.Internal, "failed to fetch token: %v", err)
		}

		// Data integrity check
		if placedTokenModel.Scene.ID == 0 {
			s.Logger.ErrorF("Token %d data is inconsistent: missing scene link", placedTokenModel.ID)
			return nil, status.Errorf(codes.Internal, "object data is inconsistent")
		}

		// Store context
		sceneID = placedTokenModel.Scene.ID
		tableID = placedTokenModel.Scene.TableID

		s.Logger.InfoF("Verifying user is master for TableID: %d", tableID)
		// Check if the user has rights to change permissions (must be GM)
		err := utils.CheckUserIsMaster(ctx, s.DB, tableID)
		if err != nil {
			s.Logger.ErrorF("User does not have master permissions for this object: %v", err)
			return nil, err
		}

	case permission.GameObjectType_IMAGE:
		s.Logger.InfoF("Object Type: Image Selected...")

		// Fetch the placed image and preload its associated Scene
		if err := s.DB.Preload("Scene").Where("id = ?", req.ObjectId).First(&placedImageModel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Logger.ErrorF("Placed image with ID %d not found", req.ObjectId)
				return nil, status.Errorf(codes.NotFound, "image not found")
			}
			s.Logger.ErrorF("Failed to fetch placed image %d: %v", req.ObjectId, err)
			return nil, status.Errorf(codes.Internal, "failed to fetch image: %v", err)
		}

		// Data integrity check
		if placedImageModel.Scene.ID == 0 {
			s.Logger.ErrorF("Image %d data is inconsistent: missing scene link", placedImageModel.ID)
			return nil, status.Errorf(codes.Internal, "object data is inconsistent")
		}

		// Store context
		sceneID = placedImageModel.Scene.ID
		tableID = placedImageModel.Scene.TableID

		s.Logger.InfoF("Verifying user is master for TableID: %d", tableID)
		// Check if the user has rights to change permissions (must be GM)
		err := utils.CheckUserIsMaster(ctx, s.DB, tableID)
		if err != nil {
			s.Logger.ErrorF("User does not have master permissions for this object: %v", err)
			return nil, err
		}

	default:
		s.Logger.ErrorF("Invalid object_type provided: %s", req.GetObjectType().String())
		return nil, status.Errorf(codes.InvalidArgument, "invalid object_type")
	}

	s.Logger.InfoF("Starting transaction to update owners for object %d", req.ObjectId)

	// Start a transaction to ensure atomicity
	tx := s.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		s.Logger.ErrorF("Failed to start transaction: %v", tx.Error)
		return nil, status.Errorf(codes.Internal, "failed to start transaction: %v", tx.Error)
	}

	// Defer rollback in case of any error; tx.Commit() will nullify this if successful
	defer tx.Rollback()

	// Helper function to remove duplicate user IDs and 0s from the request
	dedupe := func(xs []uint64) []uint64 {
		if len(xs) == 0 {
			return xs
		}
		m := make(map[uint64]struct{}, len(xs))
		out := make([]uint64, 0, len(xs))
		for _, v := range xs {
			if v == 0 { // Ignore invalid user ID 0
				continue
			}
			if _, ok := m[v]; !ok {
				m[v] = struct{}{}
				out = append(out, v)
			}
		}
		return out
	}

	// Clean the incoming list of owners
	cleanOwners := dedupe(req.OwnerUserIds)
	s.Logger.InfoF("Updating owners list (new count: %d)", len(cleanOwners))

	// Switch again to perform the actual database updates
	switch req.GetObjectType() {

	case permission.GameObjectType_TOKEN:
		s.Logger.InfoF("Removing old owners from token object %d...", req.ObjectId)
		// Clear all existing owners for this token
		if err := tx.Where("placed_token_id = ?", req.ObjectId).Delete(&models.GameObjectOwner{}).Error; err != nil {
			s.Logger.ErrorF("Failed to delete old owners for token %d: %v", placedTokenModel.ID, err)
			return nil, status.Error(codes.Internal, "failed to clear owners")
		}

		// If a new list of owners was provided, create them
		if len(cleanOwners) > 0 {
			s.Logger.InfoF("Creating %d new owner entries for token %d", len(cleanOwners), req.ObjectId)
			rows := make([]models.GameObjectOwner, 0, len(cleanOwners))
			for _, ids := range cleanOwners {
				rows = append(rows, models.GameObjectOwner{
					UserId:        uint(ids),
					PlacedTokenID: &placedTokenModel.ID,
				})
			}
			// Batch create the new owner entries
			if err := tx.Create(&rows).Error; err != nil {
				s.Logger.ErrorF("Failed to create new owners for token %d: %v", placedTokenModel.ID, err)
				return nil, status.Error(codes.Internal, "failed to set owners")
			}
		}

		// Update permission levels if they were provided in the request
		if req.CanBeViewedBy != nil {
			placedTokenModel.CanBeViewedBy = consts.PermissionLevel(req.GetCanBeViewedBy())
		}
		if req.CanBeModifiedBy != nil {
			placedTokenModel.CanBeModifiedBy = consts.PermissionLevel(req.GetCanBeModifiedBy())
		}
		s.Logger.InfoF("Saving new permission levels for token object %d...", req.ObjectId)
		// Save the updated permission levels (View/Modify)
		if err := tx.Model(&placedTokenModel).Select("can_be_viewed_by", "can_be_modified_by").
			Updates(&placedTokenModel).Error; err != nil {
			s.Logger.ErrorF("Failed to update permissions for token %d: %v", placedTokenModel.ID, err)
			return nil, status.Error(codes.Internal, "failed to update permissions")
		}

	case permission.GameObjectType_IMAGE:
		s.Logger.InfoF("Removing old owners from image object %d...", req.ObjectId)
		// Clear all existing owners for this image
		if err := tx.Where("placed_image_id = ?", req.ObjectId).
			Delete(&models.GameObjectOwner{}).Error; err != nil {
			s.Logger.ErrorF("Failed to delete old owners for image %d: %v", placedImageModel.ID, err)
			return nil, status.Error(codes.Internal, "failed to clear owners")
		}

		// If a new list of owners was provided, create them
		if len(cleanOwners) > 0 {
			s.Logger.InfoF("Creating %d new owner entries for image %d", len(cleanOwners), req.ObjectId)
			rows := make([]models.GameObjectOwner, 0, len(cleanOwners))
			for _, uid := range cleanOwners {
				rows = append(rows, models.GameObjectOwner{
					UserId:        uint(uid),
					PlacedImageID: &placedImageModel.ID,
				})
			}
			// Batch create the new owner entries
			if err := tx.Create(&rows).Error; err != nil {
				s.Logger.ErrorF("Failed to create new owners for image %d: %v", placedImageModel.ID, err)
				return nil, status.Error(codes.Internal, "failed to set owners")
			}
		}

		// Update permission levels if they were provided in the request
		if req.CanBeViewedBy != nil {
			placedImageModel.CanBeViewedBy = consts.PermissionLevel(req.GetCanBeViewedBy())
		}
		if req.CanBeModifiedBy != nil {
			placedImageModel.CanBeModifiedBy = consts.PermissionLevel(req.GetCanBeModifiedBy())
		}
		s.Logger.InfoF("Saving new permission levels for image object %d...", req.ObjectId)
		// Save the updated permission levels (View/Modify)
		if err := tx.Model(&placedImageModel).Select("can_be_viewed_by", "can_be_modified_by").
			Updates(&placedImageModel).Error; err != nil {
			s.Logger.ErrorF("Failed to update permissions for image %d: %v", placedImageModel.ID, err)
			return nil, status.Error(codes.Internal, "failed to update permissions")
		}
	}

	// If all steps succeeded, commit the transaction
	if err := tx.Commit().Error; err != nil {
		s.Logger.ErrorF("Failed to commit transaction for object %d: %v", req.ObjectId, err)
		return nil, status.Error(codes.Internal, "failed to commit changes")
	}

	s.Logger.InfoF("Transaction committed. Reloading owners from DB for event payload.")

	// After successful commit, reload the definitive list of owners from DB for the event
	var owners []uint64
	switch req.GetObjectType() {
	case permission.GameObjectType_TOKEN:
		var rows []models.GameObjectOwner
		// Read the owners for the token
		if err := s.DB.WithContext(ctx).Model(&models.GameObjectOwner{}).
			Select("user_id").Where("placed_token_id = ?", req.ObjectId).Find(&rows).Error; err == nil {
			for _, r := range rows {
				owners = append(owners, uint64(r.UserId))
			}
		} else {
			s.Logger.WarningF("Failed to reload owners for event (token %d): %v", req.ObjectId, err)
		}
	case permission.GameObjectType_IMAGE:
		var rows []models.GameObjectOwner
		// Read the owners for the image
		if err := s.DB.WithContext(ctx).Model(&models.GameObjectOwner{}).
			Select("user_id").Where("placed_image_id = ?", req.ObjectId).Find(&rows).Error; err == nil {
			for _, r := range rows {
				owners = append(owners, uint64(r.UserId))
			}
		} else {
			s.Logger.WarningF("Failed to reload owners for event (image %d): %v", req.ObjectId, err)
		}
	}

	// Set final permission levels for the event payload
	var finalModifiedBy, finalViewedBy consts.PermissionLevel
	if req.ObjectType == permission.GameObjectType_TOKEN {
		finalModifiedBy = placedTokenModel.CanBeModifiedBy
		finalViewedBy = placedTokenModel.CanBeViewedBy
	} else {
		finalModifiedBy = placedImageModel.CanBeModifiedBy
		finalViewedBy = placedImageModel.CanBeModifiedBy
	}

	s.Logger.InfoF("Creating the response for requisition...")

	// Create the synchronization event payload
	syncPermission := syncBroker.SyncResponse{
		SceneId: uint64(sceneID), // Publish to the specific scene
		Action: &syncBroker.SyncResponse_PlacedObjectAccessUpdated{
			PlacedObjectAccessUpdated: &permission.PlacedObjectAccessUpdated{
				ObjectId:        req.ObjectId,
				CanBeModifiedBy: permission.Permissionlevel(finalModifiedBy), // Final value from DB
				CanBeViewedBy:   permission.Permissionlevel(finalViewedBy),   // Final value from DB
				ObjectType:      req.ObjectType,
				OwnerUserIds:    owners, // The actual list of owners now in the DB
			},
		},
	}
	s.Logger.InfoF("Publishing new PlacedObjectAccessUpdated event to SceneSync...")

	// Publish the event
	s.Broker.Publish(pubSubSyncConst.SceneSync, uint64(sceneID), &syncPermission)

	s.Logger.InfoF("Event published. Requisition finished successfully.")

	return &emptypb.Empty{}, nil
}

func (s *PermissionService) UpdateLibraryObjectVisibility(ctx context.Context, req *permission.UpdateLibraryObjectVisibilityRequest) (*emptypb.Empty, error) {
	s.Logger.InfoF("GRPC Requisition to Update a Library Object visibility started...")

	var tableID uint // To store the table ID for permission checks and event publishing

	// --- Basic Request Validation ---
	if req.GetObjectId() == 0 {
		s.Logger.ErrorF("Invalid request: objectId must be greater than 0")
		return nil, status.Error(codes.InvalidArgument, "objectId is must be higher than 0")
	}

	switch req.GetObjectType() {
	case permission.GameObjectType_TOKEN, permission.GameObjectType_IMAGE:
		// Valid types
	default:
		s.Logger.ErrorF("Invalid request: invalid objectType %s", req.GetObjectType().String())
		return nil, status.Error(codes.InvalidArgument, "invalid objectType")
	}

	if req.GetVisibility() == permission.Permissionlevel_PERMISSION_LEVEL_UNSPECIFIED {
		s.Logger.ErrorF("Invalid request: invalid visibility level (UNSPECIFIED)")
		return nil, status.Error(codes.InvalidArgument, "invalid visibility")
	}
	// --- End Validation ---

	s.Logger.InfoF("Selecting the object Type...")
	switch req.GetObjectType() {

	case permission.GameObjectType_TOKEN:
		s.Logger.InfoF("Object Type: Token Selected...")
		var tokenModel models.Token
		s.Logger.InfoF("Searching the library token object on DB...")

		// Find the token in the library (models.Token)
		if err := s.DB.WithContext(ctx).Where("id = ?", req.ObjectId).First(&tokenModel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Logger.ErrorF("Library token not found with ID %d", req.ObjectId)
				return nil, status.Error(codes.NotFound, "token not found")
			}
			s.Logger.ErrorF("Failed to fetch library token %d: %v", req.ObjectId, err)
			return nil, status.Error(codes.Internal, "failed to fetch token")
		}

		// Get TableID for permission check and event publishing
		tableID = tokenModel.TableID
		s.Logger.InfoF("Verifying user is master for TableID %d", tableID)
		if err := utils.CheckUserIsMaster(ctx, s.DB, tableID); err != nil {
			s.Logger.ErrorF("User is not master for table %d: %v", tableID, err)
			return nil, err
		}

		s.Logger.InfoF("Updating the visibility on DB for token %d...", tokenModel.ID)

		// Update the visibility level
		newVis := consts.PermissionLevel(req.GetVisibility())
		if err := s.DB.WithContext(ctx).
			Model(&tokenModel).Update("can_be_viewed_by", newVis).Error; err != nil {
			s.Logger.ErrorF("Failed to update token visibility %d: %v", tokenModel.ID, err)
			return nil, status.Error(codes.Internal, "failed to update token visibility")
		}

	case permission.GameObjectType_IMAGE:
		s.Logger.InfoF("Object Type: Image Selected...")
		var imageModel models.Image
		s.Logger.InfoF("Searching the library image object on DB")

		// Find the image in the library (models.Image)
		if err := s.DB.WithContext(ctx).Where("id = ?", req.ObjectId).First(&imageModel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.Logger.ErrorF("Library image not found with ID %d", req.ObjectId)
				return nil, status.Error(codes.NotFound, "image not found")
			}
			s.Logger.ErrorF("Failed to fetch library image %d: %v", req.ObjectId, err)
			return nil, status.Error(codes.Internal, "failed to fetch image")
		}

		// Get TableID for permission check and event publishing
		tableID = imageModel.TableID
		s.Logger.InfoF("Verifying user is master for TableID %d", tableID)
		if err := utils.CheckUserIsMaster(ctx, s.DB, tableID); err != nil {
			s.Logger.ErrorF("User is not master for table %d: %v", tableID, err)
			return nil, err
		}

		s.Logger.InfoF("Updating the visibility on DB for image %d...", imageModel.ID)

		// Update the visibility level
		newVis := consts.PermissionLevel(req.GetVisibility())
		if err := s.DB.WithContext(ctx).
			Model(&imageModel).Update("can_be_viewed_by", newVis).Error; err != nil {
			s.Logger.ErrorF("Failed to update image visibility %d: %v", imageModel.ID, err)
			return nil, status.Error(codes.Internal, "failed to update image visibility")
		}
	}

	s.Logger.InfoF("Creating event payload for library visibility update...")

	// Prepare the event payload
	resp := &permission.LibraryObjectVisibilityUpdated{
		ObjectId:   req.GetObjectId(),
		ObjectType: req.GetObjectType(),
		Visibility: req.GetVisibility(),
	}

	// Create the event
	event := events.NewLibraryObjectVisibilityUpdatedEvent(uint64(tableID), resp)

	s.Logger.InfoF("Publishing new event to TableSync for TableID %d...", tableID)

	// Publish the event to the whole table
	s.Broker.Publish(pubSubSyncConst.TableSync, uint64(tableID), event)
	s.Logger.InfoF("Event published. Requisition finished successfully.")

	return &emptypb.Empty{}, nil
}
