package models

import (
	"encoding/json"
	"github.com/GarotoCowboy/vttProject/api/models/consts"
	"gorm.io/gorm"
)

type Character struct {
	gorm.Model
	TableUserID uint `json:"table_user_id" gorm:"not null"`
	TableUser TableUser `gorm:"foreignKey:TableUserID"`
	PlayerName string `json:"player_name,omitempty"`
	Name string `json:"character_name" gorm:"not null"`
	SystemKey consts.SystemKey `json:"system_key" gorm:"not null"`
	SheetData json.RawMessage `json:"sheet_data" gorm:"type:jsonb"`
}
