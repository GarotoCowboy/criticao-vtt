package events

import (
	placed_token "github.com/GarotoCowboy/vttProject/api/grpc/pb/placedToken"
	"github.com/GarotoCowboy/vttProject/api/grpc/pb/sync"
)

func NewPlacedTokenCreatedEvent(pt *placed_token.PlacedToken) *sync.SyncResponse {

	return &sync.SyncResponse{
		SceneId: pt.SceneId,
		Action: &sync.SyncResponse_PlacedTokenCreated{
			PlacedTokenCreated: &placed_token.PlacedTokenCreated{
				PlacedToken: pt,
			},
		},
	}
}

func NewPlacedTokenUpdatedEvent(pt *placed_token.PlacedToken) *sync.SyncResponse {

	return &sync.SyncResponse{
		SceneId: pt.SceneId,
		Action: &sync.SyncResponse_PlacedTokenUpdated{
			PlacedTokenUpdated: &placed_token.PlacedTokenUpdated{
				PlacedToken: pt,
			},
		},
	}
}

func NewPlacedTokenDeletedEvent(sceneID, placedTokenID uint) *sync.SyncResponse {

	return &sync.SyncResponse{
		SceneId: uint64(sceneID),
		Action: &sync.SyncResponse_PlacedTokenDeleted{
			PlacedTokenDeleted: &placed_token.PlacedTokenDeleted{
				PlacedTokenId: uint64(placedTokenID),
				SceneId:       uint64(sceneID),
			},
		},
	}
}

func NewPlacedTokenMovedEvent(sceneId, placedTokenId uint64, posX, posY int32) *sync.SyncResponse {

	return &sync.SyncResponse{
		SceneId: sceneId,
		Action: &sync.SyncResponse_PlacedTokenMoved{
			PlacedTokenMoved: &placed_token.PlacedTokenMoved{
				PlacedTokenId: placedTokenId,
				SceneId:       sceneId,
				PosY:          posY,
				PosX:          posX,
			},
		},
	}
}
