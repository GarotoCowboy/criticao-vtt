package scene

import (
	"github.com/GarotoCowboy/vttProject/api/grpc/pb/scene"
	"github.com/GarotoCowboy/vttProject/api/grpc/service/sync/broker"
	"github.com/GarotoCowboy/vttProject/config"
	"gorm.io/gorm"
)

type SceneService struct {
	scene.UnimplementedSceneServiceServer
	Logger *config.Logger
	DB     *gorm.DB
	Broker *broker.Broker
}

func NewSceneService(logger *config.Logger, db *gorm.DB, broker *broker.Broker) *SceneService {
	return &SceneService{
		Logger: logger,
		DB:     db,
		Broker: broker,
	}
}
