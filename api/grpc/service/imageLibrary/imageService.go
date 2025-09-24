package libraryImg

import (
	"github.com/GarotoCowboy/vttProject/api/grpc/proto/imageLibrary/pb"
	"github.com/GarotoCowboy/vttProject/config"
	"gorm.io/gorm"
)

type ImageLibraryService struct {
	pb.UnimplementedImageLibraryServiceServer
	DB     *gorm.DB
	Logger *config.Logger
}

func NewImageLibraryService(db *gorm.DB, Logger *config.Logger) *ImageLibraryService {
	return &ImageLibraryService{
		DB:     db,
		Logger: Logger,
	}
}
