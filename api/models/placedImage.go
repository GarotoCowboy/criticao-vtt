package models

import (
	"github.com/GarotoCowboy/vttProject/api/models/consts"
	"gorm.io/gorm"
)

type PlacedImage struct {
	gorm.Model
	PosX uint `json:"posX" gorm:"not null"`
	PosY uint `json:"posY" gorm:"not null"`

	ImageID uint   `json:"imageID" gorm:"not null"`
	Image   *Image `json:"image" gorm:"not null"`

	SceneID uint   `json:"sceneID" gorm:"constraint:OnUpdate:CASCADE"`
	Scene   *Scene `json:"scene" gorm:"constraint:OnUpdate:CASCADE"`

	LayerType consts.LayerType `json:"layer_type"`
}
