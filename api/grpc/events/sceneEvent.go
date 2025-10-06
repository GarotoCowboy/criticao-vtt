package events

import (
	"github.com/GarotoCowboy/vttProject/api/grpc/pb/scene"
	"github.com/GarotoCowboy/vttProject/api/grpc/pb/sync"
)

func NewCreateSceneEvent(s *scene.Scene) *sync.SyncResponse {

	return &sync.SyncResponse{
		SceneId: 0,
		TableId: s.TableId,
		Action: &sync.SyncResponse_SceneCreated{
			SceneCreated: &scene.SceneCreated{
				Scene: s,
			},
		},
	}
}

func NewUpdateSceneEvent(s *scene.Scene) *sync.SyncResponse {

	return &sync.SyncResponse{
		SceneId: 0,
		TableId: s.TableId,
		Action: &sync.SyncResponse_SceneUpdated{
			SceneUpdated: &scene.SceneUpdated{
				Scene: s,
			},
		},
	}
}

func NewDeleteSceneEvent(tableId, sceneId uint64) *sync.SyncResponse {

	return &sync.SyncResponse{
		SceneId: 0,
		TableId: tableId,
		Action: &sync.SyncResponse_SceneDeleted{
			SceneDeleted: &scene.SceneDeleted{
				SceneId: sceneId,
			},
		},
	}
}
