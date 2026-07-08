package chat

import (
	"github.com/GarotoCowboy/vttProject/api/grpc/pb/chat"
	"github.com/GarotoCowboy/vttProject/api/grpc/service/sync/broker"
	"github.com/GarotoCowboy/vttProject/config"
	"gorm.io/gorm"
)

type ChatService struct {
	chat.UnimplementedChatServer
	Db                 *gorm.DB
	Logger             *config.Logger
	Broker             *broker.Broker
}

func NewChatService(db *gorm.DB, logger *config.Logger, broker *broker.Broker) *ChatService {
	return &ChatService{
		Db:                 db,
		Logger:             logger,
		Broker:             broker,
	}
}
