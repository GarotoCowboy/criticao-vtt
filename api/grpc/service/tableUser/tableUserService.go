package tableUser

import (
	"github.com/GarotoCowboy/vttProject/api/grpc/pb/tableUser"
	"github.com/GarotoCowboy/vttProject/api/grpc/service/sync/broker"
	"github.com/GarotoCowboy/vttProject/config"
	"gorm.io/gorm"
)

type TableUserService struct {
	tableUser.UnimplementedTableUserServiceServer
	DB     *gorm.DB
	Logger *config.Logger
	Broker *broker.Broker
}

func NewTableUserService(db *gorm.DB, logger *config.Logger, broker *broker.Broker) *TableUserService {
	return &TableUserService{
		DB:     db,
		Logger: logger,
		Broker: broker,
	}
}
