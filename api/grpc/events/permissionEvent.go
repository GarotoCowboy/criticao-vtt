package events

import (
	"github.com/GarotoCowboy/vttProject/api/grpc/pb/permission"
	"github.com/GarotoCowboy/vttProject/api/grpc/pb/sync"
)

func NewPlacedObjectAccessEvent(tableID uint64, eventData *permission.PlacedObjectAccessUpdated) *sync.SyncResponse {

	return &sync.SyncResponse{
		SceneId: 0,
		TableId: tableID,
		Action: &sync.SyncResponse_PlacedObjectAccessUpdated{
			PlacedObjectAccessUpdated: eventData,
		},
	}
}

func NewLibraryObjectVisibilityUpdatedEvent(tableID uint64, eventData *permission.LibraryObjectVisibilityUpdated) *sync.SyncResponse {

	return &sync.SyncResponse{
		SceneId: 0,
		TableId: tableID,
		Action: &sync.SyncResponse_LibraryObjectVisiblityUpdated{
			LibraryObjectVisiblityUpdated: eventData,
		},
	}
}
