package events

import (
	"github.com/GarotoCowboy/vttProject/api/grpc/pb/chat"
	"github.com/GarotoCowboy/vttProject/api/grpc/pb/sync"
)

func NewSendChatMessage(c *chat.ChatMessageResponse) *sync.SyncResponse {

	return &sync.SyncResponse{
		SceneId: 0,
		TableId: c.TableId,
		Action: &sync.SyncResponse_MessageSended{
			MessageSended: &chat.ChatMessageSent{
				Message: c,
			},
		},
	}
}

func NewUpdateChatMessage(c *chat.ChatMessageResponse) *sync.SyncResponse {

	return &sync.SyncResponse{
		SceneId: 0,
		TableId: c.TableId,
		Action: &sync.SyncResponse_MessageUpdated{
			MessageUpdated: &chat.ChatMessageUpdated{
				Message: c,
			},
		},
	}
}

func NewRemoveChatMessage(messageUUID string, tableID uint64) *sync.SyncResponse {

	return &sync.SyncResponse{
		SceneId: 0,
		TableId: tableID,
		Action: &sync.SyncResponse_MessageDeleted{
			MessageDeleted: &chat.ChatMessageDeleted{
				TableId:     tableID,
				MessageUuid: messageUUID,
			},
		},
	}
}
