package models

import "gorm.io/gorm"

type Bar struct {
	gorm.Model
	Name     string `json:"name" gorm:"not null"`
	Value    int32  `json:"value"`
	MaxValue int32  `json:"max_value"`
	Color    string `json:"color"`
	TokenID  uint   `json:"token_id" gorm:"not null"`
	Token    Token  `json:"token" gorm:"foreignKey:TokenID"`
}
