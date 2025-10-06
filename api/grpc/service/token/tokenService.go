package token

import (
	"github.com/GarotoCowboy/vttProject/api/grpc/pb/token"
	"github.com/GarotoCowboy/vttProject/api/grpc/service/sync/broker"
	"github.com/GarotoCowboy/vttProject/config"
	"gorm.io/gorm"
)

type TokenService struct {
	token.UnimplementedTokenServiceServer
	DB     *gorm.DB
	Logger *config.Logger
	Broker *broker.Broker
}

func NewTokenService(db *gorm.DB, logger *config.Logger, broker *broker.Broker) *TokenService {
	return &TokenService{
		DB:     db,
		Logger: logger,
		Broker: broker,
	}
}
