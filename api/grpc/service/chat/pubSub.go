package chat

import (
	"github.com/GarotoCowboy/vttProject/api/grpc/proto/chat/pb"
	"github.com/google/uuid"
)

func (s *ChatService) publicChatSubscribe(id uint, stream pb.Chat_SendMessageServer) string {
	subID := uuid.NewString()
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.publicSubscribers[id] == nil {
		s.publicSubscribers[id] = make(map[string]pb.Chat_SendMessageServer)
	}
	s.publicSubscribers[id][subID] = stream
	return subID
}

func (s *ChatService) publicChatUnsubscribe(id uint, subID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if subs, ok := s.publicSubscribers[id]; ok {
		delete(subs, subID)
		if len(subs) == 0 {
			delete(s.publicSubscribers, id)
		}
	}
}

func (s *ChatService) privateChatSubscribe(fromUserId, toUserId uint, stream pb.Chat_SendPrivateMessageServer) string {
	subId := uuid.NewString()
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.privateSubscribers[fromUserId] == nil {
		s.privateSubscribers[fromUserId] = make(map[uint]map[string]pb.Chat_SendPrivateMessageServer)
	}

	if s.privateSubscribers[fromUserId][toUserId] == nil {
		s.privateSubscribers[fromUserId][toUserId] = make(map[string]pb.Chat_SendPrivateMessageServer)
	}
	s.privateSubscribers[fromUserId][toUserId][subId] = stream
	return subId
}

func (s *ChatService) privateChatUnsubscribe(fromUserId, toUserId uint, subID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if subs, ok := s.privateSubscribers[fromUserId][toUserId]; ok {
		delete(subs, subID)
		if len(subs) == 0 {
			delete(s.privateSubscribers[fromUserId], toUserId)
		}
	}
}
