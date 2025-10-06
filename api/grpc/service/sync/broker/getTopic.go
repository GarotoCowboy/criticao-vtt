package broker

import (
	"fmt"

	"github.com/GarotoCowboy/vttProject/api/models/consts/pubSubSyncConst"
)

func getTopic(topicType pubSubSyncConst.PubSubSyncType, id uint64) string {
	switch topicType {
	case pubSubSyncConst.TableSync:
		return fmt.Sprintf("table:%d", id)
	case pubSubSyncConst.SceneSync:
		return fmt.Sprintf("scene:%d", id)
	default:
		return "unknown:0"
	}

}
