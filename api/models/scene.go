package models

import "gorm.io/gorm"

type Scene struct {
	gorm.Model
	Name    string `json:"name"`
	Width   uint   `json:"width" gorm:"not null"`
	Height  uint   `json:"height" gorm:"not null"`
	TableID uint   `json:"table_id" gorm:"not null"`

	PlacedTokens []*PlacedToken `json:"placedTokens" gorm:"foreignKey:SceneID"`
	PlacedImages []*PlacedImage `json:"placedImages" gorm:"foreignKey:SceneID"`
}
