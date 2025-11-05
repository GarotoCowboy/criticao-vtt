package utils

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func CheckUserCanEditTokenObject(ctx context.Context, db *gorm.DB, placedTokenID uint) error {

	userId, err := PickUserIdJWT(ctx)
	if err != nil {
		return err
	}

	if placedTokenID == 0 {
		return status.Errorf(codes.InvalidArgument, "invalid token id")
	}

	var exists bool

	err = db.WithContext(ctx).Raw(
		`
			SELECT EXISTS(
				SELECT 1
				FROM game_object_owners
				WHERE user_id = ? AND placed_token_id = ?  AND deleted_at IS NULL
					)
			`, userId, placedTokenID).Scan(&exists).Error
	if err != nil {
		return status.Errorf(codes.Internal, "database error")
	}

	if !exists {
		return status.Error(codes.PermissionDenied, "user cannot have permission to edit the object")
	}
	return nil

}

func CheckUserCanEditImageObject(ctx context.Context, db *gorm.DB, placedImageId uint) error {

	userId, err := PickUserIdJWT(ctx)
	if err != nil {
		return err
	}

	if placedImageId == 0 {
		return status.Errorf(codes.InvalidArgument, "invalid token id")
	}

	var exists bool

	err = db.WithContext(ctx).Raw(
		`
			SELECT EXISTS(
				SELECT 1
				FROM game_object_owners
				WHERE user_id = ? AND placed_image_id = ?  AND deleted_at IS NULL
					)
			`, userId, placedImageId).Scan(&exists).Error
	if err != nil {
		return status.Errorf(codes.Internal, "database error")
	}

	if !exists {
		return status.Error(codes.PermissionDenied, "user cannot have permission to edit the object")
	}
	return nil

}
