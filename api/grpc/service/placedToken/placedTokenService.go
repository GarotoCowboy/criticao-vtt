package placedToken

import (
	"github.com/GarotoCowboy/vttProject/api/grpc/pb/placedToken"
	"github.com/GarotoCowboy/vttProject/api/grpc/service/sync/broker"

	"github.com/GarotoCowboy/vttProject/config"
	"gorm.io/gorm"
)

type PlacedTokenService struct {
	placedToken.UnimplementedPlacedTokenServiceServer
	Logger *config.Logger
	DB     *gorm.DB
	Broker *broker.Broker
}

func NewPlacedTokenService(db *gorm.DB, logger *config.Logger, broker *broker.Broker) *PlacedTokenService {
	return &PlacedTokenService{
		DB:     db,
		Logger: logger,
		Broker: broker,
	}
}
