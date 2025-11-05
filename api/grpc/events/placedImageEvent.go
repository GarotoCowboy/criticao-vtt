package events

import (
	placed_image "github.com/GarotoCowboy/vttProject/api/grpc/pb/placedImage"
	"github.com/GarotoCowboy/vttProject/api/grpc/pb/sync"
)

func NewPlacedImageCreatedEvent(pi *placed_image.PlacedImage) *sync.SyncResponse {

	return &sync.SyncResponse{
		SceneId: pi.SceneId,
		Action: &sync.SyncResponse_PlacedImageCreated{
			PlacedImageCreated: &placed_image.PlacedImageCreated{
				PlacedImage: pi,
			},
		},
	}
}

func NewPlacedImageUpdatedEvent(pi *placed_image.PlacedImage) *sync.SyncResponse {

	return &sync.SyncResponse{
		SceneId: pi.SceneId,
		Action: &sync.SyncResponse_PlacedImageUpdated{
			PlacedImageUpdated: &placed_image.PlacedImageUpdated{
				PlacedImage: pi,
			},
		},
	}
}

func NewPlacedImageDeletedEvent(sceneID, placedImageID uint) *sync.SyncResponse {

	return &sync.SyncResponse{
		SceneId: uint64(sceneID),
		Action: &sync.SyncResponse_PlacedImageDeleted{
			PlacedImageDeleted: &placed_image.PlacedImageDeleted{
				PlacedImageId: uint64(placedImageID),
				SceneId:       uint64(sceneID),
			},
		},
	}
}

func NewPlacedImageMovedEvent(sceneId, placedImageID uint64, posX, posY int32) *sync.SyncResponse {

	return &sync.SyncResponse{
		SceneId: sceneId,
		Action: &sync.SyncResponse_PlacedImageMoved{
			PlacedImageMoved: &placed_image.PlacedImageMoved{
				PlacedImageId: placedImageID,
				SceneId:       sceneId,
				PosY:          posY,
				PosX:          posX,
			},
		},
	}
}
