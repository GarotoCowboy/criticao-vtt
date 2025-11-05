package utils

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

const RoleGM = 2

func CheckUserIsMaster(ctx context.Context, db *gorm.DB, tableID uint) error {

	userId, err := PickUserIdJWT(ctx)
	if err != nil {
		return err
	}

	if tableID == 0 {
		return status.Errorf(codes.InvalidArgument, "invalid table id")
	}

	var exists bool

	err = db.WithContext(ctx).Raw(
		`
			SELECT EXISTS(
				SELECT 1
				FROM table_users
				WHERE user_id = ? AND table_id = ? AND role = ? AND deleted_at IS NULL
					)
			`, userId, tableID, RoleGM).
		Scan(&exists).Error
	if err != nil {
		return status.Errorf(codes.Internal, "database error")
	}
	if !exists {
		return status.Error(codes.PermissionDenied, "user is not the game master")
	}
	return nil
}
