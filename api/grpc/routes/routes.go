package routes

import (
	pbBar "github.com/GarotoCowboy/vttProject/api/grpc/proto/bar/pb"
	"github.com/GarotoCowboy/vttProject/api/grpc/proto/character/pb"
	pbChat "github.com/GarotoCowboy/vttProject/api/grpc/proto/chat/pb"
	pbImage "github.com/GarotoCowboy/vttProject/api/grpc/proto/imageLibrary/pb"
	pbToken "github.com/GarotoCowboy/vttProject/api/grpc/proto/token/pb"
	"github.com/GarotoCowboy/vttProject/api/grpc/service/bar"
	"github.com/GarotoCowboy/vttProject/api/grpc/service/character"
	"github.com/GarotoCowboy/vttProject/api/grpc/service/chat"
	imageLibrary "github.com/GarotoCowboy/vttProject/api/grpc/service/imageLibrary"
	"github.com/GarotoCowboy/vttProject/api/grpc/service/token"
	"github.com/GarotoCowboy/vttProject/config"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

// Routes code for GRPC
func Routes(r *grpc.Server, db *gorm.DB, logger *config.Logger) {
	characterService := character.NewCharacterService(db, logger)
	chatService := chat.NewChatService(db, logger)
	tokenService := token.NewTokenService(db, logger)
	barService := bar.NewBarService(db, logger)
	imageService := imageLibrary.NewImageLibraryService(db, logger)
	//Implements the router for characterServiceGRPC

	pb.RegisterCharacterServiceServer(r, characterService)

	//Implements the router for chatServiceGRPC
	pbChat.RegisterChatServer(r, chatService)

	//Implements the router for tokenServiceGRPC
	pbToken.RegisterTokenServiceServer(r, tokenService)

	//Implements the router for barServiceGRPC
	pbBar.RegisterBarServiceServer(r, barService)

	//Implements the router for ImageLibrary
	pbImage.RegisterImageLibraryServiceServer(r, imageService)

}
