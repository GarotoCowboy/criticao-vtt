package models

import (
	"gorm.io/gorm"
)

type Token struct {
	gorm.Model
	Name     string `json:"name" gorm:"not null"`
	ImageURL string `json:"image_url"`
	Bars     []*Bar `json:"bars" gorm:"foreignkey:TokenID"`

	TableID uint  `json:"table_id" gorm:"not null"`
	Table   Table `json:"table" gorm:"foreignkey:TableID"`
}
