package models

import "gorm.io/gorm"

type Bar struct {
	gorm.Model
	Name     string `json:"name"`
	Value    int    `json:"value"`
	MaxValue int    `json:"max_value"`
	Color    string `json:"color"`

	TokenID uint  `json:"token_id" gorm:"not null"`
	Token   Token `json:"token" gorm:"foreignKey:TokenID"`
}
