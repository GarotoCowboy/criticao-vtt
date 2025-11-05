package placedImage

import (
	"github.com/GarotoCowboy/vttProject/api/grpc/pb/placedImage"
	"github.com/GarotoCowboy/vttProject/api/grpc/service/sync/broker"
	"github.com/GarotoCowboy/vttProject/config"
	"gorm.io/gorm"
)

type PlacedImageService struct {
	placedImage.UnimplementedPlacedImageServiceServer
	Logger *config.Logger
	DB     *gorm.DB
	Broker *broker.Broker
}

func NewPlacedImageService(db *gorm.DB, logger *config.Logger, broker *broker.Broker) *PlacedImageService {
	return &PlacedImageService{
		DB:     db,
		Logger: logger,
		Broker: broker,
	}
}
