package models

import "github.com/GarotoCowboy/vttProject/api/models/consts"

type PlacedToken struct {
	TokenID uint  `json:"token_id" gorm:"not null"`
	Token   Token `gorm:"foreignKey:TokenID"`

	SceneID uint  `json:"scene_id" gorm:"not null"`
	Scene   Scene `gorm:"foreignKey:SceneID"`

	PosX uint `json:"posX" gorm:"not null"`
	PosY uint `json:"posY" gorm:"not null"`

	Width  uint `json:"width" gorm:"not null"`
	Height uint `json:"height" gorm:"not null"`

	LayerType consts.LayerType `json:"layer_type"`
}
