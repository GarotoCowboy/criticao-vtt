package sync

import (
	"sync"

	pbBar "github.com/GarotoCowboy/vttProject/api/grpc/pb/bar"
	pbScene "github.com/GarotoCowboy/vttProject/api/grpc/pb/scene"
	pbSync "github.com/GarotoCowboy/vttProject/api/grpc/pb/sync"
	pbToken "github.com/GarotoCowboy/vttProject/api/grpc/pb/token"
	"github.com/GarotoCowboy/vttProject/config"
	"gorm.io/gorm"
)

type ConnectionHub struct {
	mu sync.Mutex

	scenes        map[uint64][]pbSync.SyncService_SyncServer
	broadCastChan chan *pbSync.SyncResponse
}

type Server struct {
	pbSync.UnimplementedSyncServiceServer
	pbScene.UnimplementedSceneServiceServer
	pbToken.UnimplementedTokenServiceServer
	pbBar.UnimplementedBarServiceServer

	db     *gorm.DB
	logger *config.Logger
	hub    *ConnectionHub
}
