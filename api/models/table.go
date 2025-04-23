package models

import "gorm.io/gorm"


//todo: TA ERRADO TEM Q ARRUMAR >:(
type Table struct{
	gorm.Model
	Name string `json:"name" gorm:"not null"`
	PlayerList []User
	GameMasters []User
	InviteLink string



}
