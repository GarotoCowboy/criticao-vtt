package models

import (
	"github.com/GarotoCowboy/vttProject/api/models/consts"
	"gorm.io/gorm"
)

type PlacedImage struct {
	gorm.Model
	PosX int32 `json:"posX" gorm:"not null"`
	PosY int32 `json:"posY" gorm:"not null"`

	ImageID uint   `json:"imageID" gorm:"not null"`
	Image   *Image `json:"image" gorm:"not null"`

	SceneID uint   `json:"sceneID" gorm:"constraint:OnUpdate:CASCADE"`
	Scene   *Scene `json:"scene" gorm:"constraint:OnUpdate:CASCADE"`

	LayerType consts.LayerType `json:"layer_type" gorm:"not null;default:1"`
	Rotation  int              `json:"rotation" gorm:"not null;default:0"`

	Width  uint `json:"width" gorm:"not null"`
	Height uint `json:"height" gorm:"not null"`

	CanBeViewedBy   consts.PermissionLevel `json:"can_view_by" gorm:"default:4"`
	CanBeModifiedBy consts.PermissionLevel `json:"can_be_modified_by" gorm:"default:2"`

	Owners []GameObjectOwner `json:"owners" gorm:"foreignKey:PlacedImageID;constraint:OnDelete:CASCADE"`
}
