package models

import (
	"github.com/GarotoCowboy/vttProject/api/models/consts"
	"gorm.io/gorm"
)

type Scene struct {
	gorm.Model
	Name                string          `json:"name" gorm:"not null"`
	Width               uint            `json:"width" gorm:"not null"`
	Height              uint            `json:"height" gorm:"not null"`
	BackgroundImagePath string          `json:"background_image_path"`
	TableID             uint            `json:"table_id" gorm:"not null"`
	IsVisible           bool            `json:"visibilty" gorm:"not null,default:false"`
	BackGroundColor     string          `json:"back_ground_color" gorm:"not null"`
	GridCellDistance    uint            `json:"distance_grid_cells"`
	GridType            consts.GridType `json:"grid_type" gorm:"not null"`
	PlacedTokens        []*PlacedToken  `json:"placedTokens" gorm:"foreignKey:SceneID"`
	PlacedImages        []*PlacedImage  `json:"placedImages" gorm:"foreignKey:SceneID"`
	Drawings            []*Drawing      `json:"drawing" gorm:"foreignKey:SceneID"`
}
