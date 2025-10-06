package sync

import (
	syncBroker "github.com/GarotoCowboy/vttProject/api/grpc/pb/sync"
	"github.com/GarotoCowboy/vttProject/api/grpc/service/sync/broker"
	"github.com/GarotoCowboy/vttProject/config"
)

type SyncServer struct {
	syncBroker.UnimplementedSyncServiceServer
	Broker *broker.Broker
	Logger *config.Logger
}

func NewSyncServer(broker *broker.Broker, Logger *config.Logger) *SyncServer {
	return &SyncServer{
		Broker: broker,
		Logger: Logger,
	}
}
