package middleware

import (
	"context"
	"errors"

	"github.com/GarotoCowboy/vttProject/api/grpc/interfaces"
	"github.com/GarotoCowboy/vttProject/api/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func GrpcTableMemberInterceptor(db *gorm.DB) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		userID, ok := ctx.Value("user_id").(uint)
		if !ok {
			return nil, status.Error(codes.Internal, "user id is not in context")
		}
		tableReq, ok := req.(interfaces.TableRequest)
		if !ok {
			return handler(ctx, req)
		}

		tableID := tableReq.GetTableId()
		var tableUser models.TableUser

		if err := db.Where("user_id = ? AND table_id = ?", userID, tableID).First(&tableUser).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, status.Errorf(codes.NotFound, "user does not have permission to access table %d", tableID)
			}
			return nil, status.Errorf(codes.Internal, "error checking table permissions")
		}
		return handler(ctx, req)
	}
}
