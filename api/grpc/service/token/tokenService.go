package token

import (
	"github.com/GarotoCowboy/vttProject/api/grpc/proto/token/pb"
	"github.com/GarotoCowboy/vttProject/config"
	"gorm.io/gorm"
)

type TokenService struct {
	pb.UnimplementedTokenServiceServer
	DB     *gorm.DB
	Logger *config.Logger
}

func NewTokenService(db *gorm.DB, logger *config.Logger) *TokenService {
	return &TokenService{
		DB:     db,
		Logger: logger,
	}
}
