package routes

import (
	"github.com/GarotoCowboy/vttProject/api/grpc/proto/character/pb"
	pbChat "github.com/GarotoCowboy/vttProject/api/grpc/proto/chat/pb"
	"github.com/GarotoCowboy/vttProject/api/grpc/service/character"
	"github.com/GarotoCowboy/vttProject/api/grpc/service/chat"
	"github.com/GarotoCowboy/vttProject/config"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

// Routes code for GRPC
func Routes(r *grpc.Server, db *gorm.DB, logger *config.Logger) {
	characterService := character.NewCharacterService(db, logger)
	chatService := chat.NewChatService(db, logger)
	//Implements the router for characterServiceGRPC
	pb.RegisterCharacterServiceServer(r, characterService)

	pbChat.RegisterChatServer(r, chatService)

}
