package events

import (
	"github.com/GarotoCowboy/vttProject/api/grpc/pb/imageLibrary"
	"github.com/GarotoCowboy/vttProject/api/grpc/pb/sync"
)

func NewImageLibraryUploadedEvent(il *imageLibrary.Image) *sync.SyncResponse {

	return &sync.SyncResponse{
		SceneId: 0,
		TableId: il.TableId,
		Action: &sync.SyncResponse_ImageUploaded{
			ImageUploaded: &imageLibrary.ImageUploaded{
				Image: il,
			},
		},
	}
}

func NewImageLibraryUpdatedEvent(il *imageLibrary.Image) *sync.SyncResponse {
	return &sync.SyncResponse{
		SceneId: 0,
		TableId: il.TableId,
		Action: &sync.SyncResponse_ImageUpdated{
			ImageUpdated: &imageLibrary.ImageUpdated{
				Image: il,
			},
		},
	}
}

func NewImageLibraryDeletedEvent(imageID, tableID uint64) *sync.SyncResponse {
	return &sync.SyncResponse{
		SceneId: 0,
		TableId: tableID,
		Action: &sync.SyncResponse_ImageDeleted{
			ImageDeleted: &imageLibrary.ImageDeleted{
				TableId: tableID,
				ImageId: imageID,
			},
		},
	}
}
