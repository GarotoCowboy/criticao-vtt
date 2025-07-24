package models

import "gorm.io/gorm"

// create an table
type Table struct {
	gorm.Model

	Name string `json:"name" gorm:"not null"`

	OwnerID uint `json:"owner_id"`
	Owner   User `json:"owner" gorm:"foreignKey:OwnerID;constraint:OnUpdate:CASCADE;constraint:OnDelete:CASCADE"`

	Members []TableUser `gorm:"foreignKey:TableID;constraint:OnDelete:CASCADE"`

	InviteLink string `json:"inviteLink" gorm:"unique"`
	Password   string `json:"password"`
	//ActionLog []string `json:"actionLog"`
}
