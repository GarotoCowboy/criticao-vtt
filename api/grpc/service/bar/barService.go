package bar

import (
	"github.com/GarotoCowboy/vttProject/api/grpc/pb/bar"
	"github.com/GarotoCowboy/vttProject/config"
	"gorm.io/gorm"
)

type BarService struct {
	bar.UnimplementedBarServiceServer
	DB     *gorm.DB
	Logger *config.Logger
}

func NewBarService(db *gorm.DB, Logger *config.Logger) *BarService {
	return &BarService{
		DB:     db,
		Logger: Logger,
	}
}
