package broker

import (
	syncBroker "github.com/GarotoCowboy/vttProject/api/grpc/pb/sync"
	"github.com/GarotoCowboy/vttProject/api/models/consts/pubSubSyncConst"
)

func (b *Broker) SubscribeToTopic(topicType pubSubSyncConst.PubSubSyncType, id uint64, ch chan *syncBroker.SyncResponse) {

	b.mu.Lock()
	defer b.mu.Unlock()

	topic := getTopic(topicType, id)

	if _, ok := b.subscribers[topic]; !ok {
		b.subscribers[topic] = make(map[chan *syncBroker.SyncResponse]struct{})
	}

	b.subscribers[topic][ch] = struct{}{}
}

func (b *Broker) UnsubscribeToTopic(topicType pubSubSyncConst.PubSubSyncType, id uint64, ch chan *syncBroker.SyncResponse) {

	b.mu.Lock()
	defer b.mu.Unlock()

	topic := getTopic(topicType, id)

	if subs, ok := b.subscribers[topic]; ok {
		delete(subs, ch)
	}

}

func (b *Broker) Publish(topicType pubSubSyncConst.PubSubSyncType, id uint64, msg *syncBroker.SyncResponse) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	topic := getTopic(topicType, id)
	if subs, ok := b.subscribers[topic]; ok {
		for ch := range subs {
			select {
			case ch <- msg:
			default:

			}
		}
	}
}
