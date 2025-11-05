package tableUser

import (
	"context"
	"errors"

	"github.com/GarotoCowboy/vttProject/api/grpc/events"
	"github.com/GarotoCowboy/vttProject/api/grpc/pb/tableUser"
	"github.com/GarotoCowboy/vttProject/api/models"
	"github.com/GarotoCowboy/vttProject/api/models/consts"
	"github.com/GarotoCowboy/vttProject/api/models/consts/pubSubSyncConst"
	"github.com/GarotoCowboy/vttProject/api/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

func (s *TableUserService) PromoteOrDemoteUser(ctx context.Context, req *tableUser.PromoteOrDemoteUserRequest) (*emptypb.Empty, error) {

	err := utils.CheckUserIsMaster(ctx, s.DB, uint(req.TableId))
	if err != nil {
		return nil, err
	}

	var tableUserModel models.TableUser
	var tableModel models.Table

	if err := s.DB.Where("table_id = ? and user_id = ? AND deleted_at IS NULL", req.GetTableId(), req.GetUserId()).First(&tableUserModel).Preload("Table", &tableModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.ErrorF("user %d not found in table %s", req.GetUserId(), req.GetTableId())
			return nil, status.Error(codes.NotFound, "user not found in this table")
		}
		s.Logger.ErrorF("db error when loading table_user (table=%d user=%d): %v", req.GetTableId(), req.GetUserId(), err)
		return nil, status.Error(codes.Internal, "db error")
	}

	if tableUserModel.Role == consts.Role(tableUser.Role_MASTER) {
		return &emptypb.Empty{}, nil
	}

	if tableModel.Owner.ID == uint(req.UserId) && req.GetRole() == tableUser.Role_PLAYER {
		s.Logger.WarningF("Owner cannot be demoted a Player")
		return nil, status.Error(codes.PermissionDenied, "Owner cannot be demoted a Player")
	}

	if err := s.DB.WithContext(ctx).Model(&models.TableUser{}).Where("id = ? AND deleted_at IS NULL", tableUserModel.ID).
		Update("role", req.GetRole()).Error; err != nil {

		s.Logger.ErrorF("failed to promote or demote user %d in table %d: %v", req.GetUserId(), req.GetTableId(), err)
		return nil, status.Error(codes.Internal, "failed to promote or demote user")
	}

	resp := tableUser.TableUser{
		TableId: req.GetTableId(),
		UserId:  req.GetUserId(),
		Role:    req.GetRole(),
	}

	event := events.NewTableUserPromotedOrDemotedEvent(&resp)

	s.Broker.Publish(pubSubSyncConst.TableSync, req.GetTableId(), event)

	return &emptypb.Empty{}, nil
}
