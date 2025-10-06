package broker

import (
	"sync"

	syncBroker "github.com/GarotoCowboy/vttProject/api/grpc/pb/sync"
)

type Broker struct {
	mu          sync.RWMutex
	subscribers map[string]map[chan *syncBroker.SyncResponse]struct{}
}

func NewBroker() *Broker {
	return &Broker{
		subscribers: make(map[string]map[chan *syncBroker.SyncResponse]struct{}),
	}
}
