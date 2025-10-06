package events

import (
	"github.com/GarotoCowboy/vttProject/api/grpc/pb/sync"
	"github.com/GarotoCowboy/vttProject/api/grpc/pb/token"
)

func NewCreateTokenEvent(t *token.Token) *sync.SyncResponse {
	return &sync.SyncResponse{
		SceneId: 0,
		TableId: t.TableId,
		Action: &sync.SyncResponse_TokenCreated{
			TokenCreated: &token.TokenCreated{
				Token: t,
			},
		},
	}
}

func NewUpdatedTokenEvent(t *token.Token) *sync.SyncResponse {
	return &sync.SyncResponse{
		SceneId: 0,
		TableId: t.TableId,
		Action: &sync.SyncResponse_TokenUpdated{
			TokenUpdated: &token.TokenUpdated{
				Token: t,
			},
		},
	}
}
func NewDeleteTokenEvent(tableID, tokenID uint64) *sync.SyncResponse {
	return &sync.SyncResponse{
		SceneId: 0,
		TableId: tableID,
		Action: &sync.SyncResponse_TokenDeleted{
			TokenDeleted: &token.TokenDeleted{
				TokenId: tokenID,
				TableId: tableID,
			},
		},
	}
}
