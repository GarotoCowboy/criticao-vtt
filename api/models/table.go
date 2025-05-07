package models

import "gorm.io/gorm"

// todo: TA ERRADO TEM Q ARRUMAR >:(
type Table struct {
	gorm.Model

	Name string `json:"name" gorm:"not null"`

	OwnerID uint `json:"owner_id"`
	Owner   User `json:"owner" gorm:"foreignKey:OwnerID"`

	Members []TableUser `gorm:"foreignKey:TableID"`

	InviteLink string `json:"inviteLink" gorm:"unique"`
	Password   string `json:"password"`
	//ActionLog []string `json:"actionLog"`
}
