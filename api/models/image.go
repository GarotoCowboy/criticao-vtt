package models

import "gorm.io/gorm"

type Image struct {
	gorm.Model
	ImageURL string `json:"imageUrl" gorm:"not null"`
	Name     string `json:"name" gorm:"not null"`
	Width    uint   `json:"width" gorm:"not null"`
	Height   uint   `json:"height" gorm:"not null"`
}
