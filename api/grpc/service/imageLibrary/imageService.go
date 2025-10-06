package libraryImg

import (
	image_library "github.com/GarotoCowboy/vttProject/api/grpc/pb/imageLibrary"
	"github.com/GarotoCowboy/vttProject/api/grpc/service/sync/broker"
	"github.com/GarotoCowboy/vttProject/config"
	"gorm.io/gorm"
)

type ImageLibraryService struct {
	image_library.UnimplementedImageLibraryServiceServer
	DB     *gorm.DB
	Logger *config.Logger
	Broker *broker.Broker
}

func NewImageLibraryService(db *gorm.DB, Logger *config.Logger, broker *broker.Broker) *ImageLibraryService {
	return &ImageLibraryService{
		DB:     db,
		Logger: Logger,
		Broker: broker,
	}
}
