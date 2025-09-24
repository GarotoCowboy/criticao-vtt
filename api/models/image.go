package models

import "gorm.io/gorm"

type Image struct {
	gorm.Model
	ImagePath   string `json:"imageUrl" gorm:"not null;unique"`
	Name        string `json:"name" gorm:"not null"`
	Width       uint   `json:"width" gorm:"not null"`
	Height      uint   `json:"height" gorm:"not null"`
	TableID     uint   `json:"tableId" gorm:"not null"`
	CheckSum    string `json:"-" gorm:"not null;unique"`
	ContentType string `json:"contentType" gorm:"not null"`

	Table Table `gorm:"constraint:OnDelete:CASCADE"`
}
