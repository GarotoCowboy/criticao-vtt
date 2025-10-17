package models

import (
	"github.com/GarotoCowboy/vttProject/api/models/consts"
	"gorm.io/gorm"
)

type PlacedToken struct {
	gorm.Model
	TokenID uint  `json:"token_id" gorm:"not null"`
	Token   Token `gorm:"foreignKey:TokenID"`

	SceneID uint  `json:"scene_id" gorm:"not null"`
	Scene   Scene `gorm:"foreignKey:SceneID"`

	PosX int32 `json:"posX" gorm:"not null;default:0"`
	PosY int32 `json:"posY" gorm:"not null;default:0"`

	Size int32 `json:"size" gorm:"not null;default:1"`

	LayerType consts.LayerType `json:"layer_type" gorm:"not null;default:1"`

	Rotation int `json:"rotation" gorm:"not null;default:0"`

	CanBeViewedBy   consts.PermissionLevel `json:"can_view_by" gorm:"default:4"`
	CanBeModifiedBy consts.PermissionLevel `json:"can_be_modified_by" gorm:"default:2"`

	Owners []GameObjectOwner `json:"owners" gorm:"foreignKey:PlacedTokenID;constraint:OnDelete:CASCADE"`
}
