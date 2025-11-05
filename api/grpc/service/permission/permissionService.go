package permission

import (
	"github.com/GarotoCowboy/vttProject/api/grpc/pb/permission"
	"github.com/GarotoCowboy/vttProject/api/grpc/service/sync/broker"
	"github.com/GarotoCowboy/vttProject/config"
	"gorm.io/gorm"
)

type PermissionService struct {
	permission.UnimplementedPermissionServiceServer
	DB     *gorm.DB
	Logger *config.Logger
	Broker *broker.Broker
}

func NewPermissionService(db *gorm.DB, logger *config.Logger, broker *broker.Broker) *PermissionService {
	return &PermissionService{
		DB:     db,
		Logger: logger,
		Broker: broker,
	}
}
