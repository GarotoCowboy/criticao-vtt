package routes

import (
	barProto "github.com/GarotoCowboy/vttProject/api/grpc/pb/bar"
	characterProto "github.com/GarotoCowboy/vttProject/api/grpc/pb/character"
	chatProto "github.com/GarotoCowboy/vttProject/api/grpc/pb/chat"
	imageLibraryProto "github.com/GarotoCowboy/vttProject/api/grpc/pb/imageLibrary"
	permissionProto "github.com/GarotoCowboy/vttProject/api/grpc/pb/permission"
	placedImageProto "github.com/GarotoCowboy/vttProject/api/grpc/pb/placedImage"
	placedTokenProto "github.com/GarotoCowboy/vttProject/api/grpc/pb/placedToken"
	sceneProto "github.com/GarotoCowboy/vttProject/api/grpc/pb/scene"
	syncProto "github.com/GarotoCowboy/vttProject/api/grpc/pb/sync"
	tableUserProto "github.com/GarotoCowboy/vttProject/api/grpc/pb/tableUser"
	tokenProto "github.com/GarotoCowboy/vttProject/api/grpc/pb/token"
	"github.com/GarotoCowboy/vttProject/api/grpc/service/placedImage"
	"github.com/GarotoCowboy/vttProject/api/grpc/service/tableUser"

	"github.com/GarotoCowboy/vttProject/api/grpc/service/bar"
	characterNewService "github.com/GarotoCowboy/vttProject/api/grpc/service/character"
	"github.com/GarotoCowboy/vttProject/api/grpc/service/chat"
	imageLibraryS "github.com/GarotoCowboy/vttProject/api/grpc/service/imageLibrary"
	"github.com/GarotoCowboy/vttProject/api/grpc/service/permission"
	"github.com/GarotoCowboy/vttProject/api/grpc/service/placedToken"
	"github.com/GarotoCowboy/vttProject/api/grpc/service/scene"
	"github.com/GarotoCowboy/vttProject/api/grpc/service/sync"
	"github.com/GarotoCowboy/vttProject/api/grpc/service/sync/broker"
	"github.com/GarotoCowboy/vttProject/api/grpc/service/token"
	"github.com/GarotoCowboy/vttProject/config"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

// Routes code for GRPC
func Routes(r *grpc.Server, db *gorm.DB, logger *config.Logger, broker *broker.Broker) {
	characterService := characterNewService.NewCharacterService(db, logger)
	chatService := chat.NewChatService(db, logger, broker)
	tokenService := token.NewTokenService(db, logger, broker)
	barService := bar.NewBarService(db, logger)
	imageService := imageLibraryS.NewImageLibraryService(db, logger, broker)
	sceneService := scene.NewSceneService(logger, db, broker)
	placedTokenService := placedToken.NewPlacedTokenService(db, logger, broker)
	permissionService := permission.NewPermissionService(db, logger, broker)
	syncService := sync.NewSyncServer(broker, logger)
	tableUserService := tableUser.NewTableUserService(db, logger, broker)
	placedImageService := placedImage.NewPlacedImageService(db, logger, broker)
	//Implements the router for characterServiceGRPC

	characterProto.RegisterCharacterServiceServer(r, characterService)

	//Implements the router for chatServiceGRPC
	chatProto.RegisterChatServer(r, chatService)

	//Implements the router for tokenServiceGRPC
	tokenProto.RegisterTokenServiceServer(r, tokenService)

	//Implements the router for barServiceGRPC
	barProto.RegisterBarServiceServer(r, barService)

	//Implements the router for ImageLibrary
	imageLibraryProto.RegisterImageLibraryServiceServer(r, imageService)

	//Implements the router for scene
	sceneProto.RegisterSceneServiceServer(r, sceneService)

	//Implements the router for PlacedToken
	placedTokenProto.RegisterPlacedTokenServiceServer(r, placedTokenService)

	//Implements the router for sync
	syncProto.RegisterSyncServiceServer(r, syncService)

	//Implements the router for permissions
	permissionProto.RegisterPermissionServiceServer(r, permissionService)

	//Implements the router for tableUser
	tableUserProto.RegisterTableUserServiceServer(r, tableUserService)

	//Implements the router for placedImage
	placedImageProto.RegisterPlacedImageServiceServer(r, placedImageService)
}
