package scene

import (
	"github.com/GarotoCowboy/vttProject/api/grpc/proto/scene/pb"
	"github.com/GarotoCowboy/vttProject/config"
	"gorm.io/gorm"
)

type SceneService struct {
	pb.UnimplementedSceneServiceServer
	Logger *config.Logger
	DB     *gorm.DB
}

func NewSceneService(logger *config.Logger, db *gorm.DB) *SceneService {
	return &SceneService{
		Logger: logger,
		DB:     db,
	}
}
