package logs

import "time"

type KickLog struct {
	Timestamp    time.Time `json:"timestamp"`
	Reason       string    `json:reason`
	TableId      uint      `json:table_id`
	MasterId     uint      `json:master_id`
	UserKickedId uint      `json:user_kick_id`
}
