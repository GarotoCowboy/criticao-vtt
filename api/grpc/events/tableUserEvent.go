package events

import (
	"github.com/GarotoCowboy/vttProject/api/grpc/pb/sync"
	"github.com/GarotoCowboy/vttProject/api/grpc/pb/tableUser"
)

func NewTableUserPromotedOrDemotedEvent(t *tableUser.TableUser) *sync.SyncResponse {

	return &sync.SyncResponse{
		SceneId: 0,
		TableId: t.TableId,
		Action: &sync.SyncResponse_UserPromotedDemoted{
			UserPromotedDemoted: &tableUser.PromotedOrDemotedUserEvent{
				TableUser: t,
			},
		},
	}
}
