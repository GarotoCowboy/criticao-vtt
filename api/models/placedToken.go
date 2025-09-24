package models

import "github.com/GarotoCowboy/vttProject/api/models/consts"

type PlacedToken struct {
	TokenID uint  `json:"token_id" gorm:"not null"`
	Token   Token `gorm:"foreignKey:TokenID"`

	SceneID uint  `json:"scene_id" gorm:"not null"`
	Scene   Scene `gorm:"foreignKey:SceneID"`

	PosX uint `json:"posX" gorm:"not null;default:0"`
	PosY uint `json:"posY" gorm:"not null;default:0"`

	Size int `json:"size" gorm:"not null;default:1"`

	LayerType consts.LayerType `json:"layer_type" gorm:"not null;default:1"`

	Rotation int `json:"rotation" gorm:"not null;default:0"`
}
