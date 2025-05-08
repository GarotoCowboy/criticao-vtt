package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string      `json:"username"  gorm:"not null;unique"`
	Password  string      `json:"password"    gorm:"not null"`
	Email     string      `json:"email"   gorm:"not null;unique"`
	Firstname string      `json:"firstname"   gorm:"not null"`
	Lastname  string      `json:"lastname"    gorm:"not null"`
	ImageLink string      `json:"imageLink"`
	Tables    []TableUser `json:"joinedTable" gorm:"foreignKey:UserID"`
}

type GameMaster interface {
	BanUser(u *User)
	SetRole(set Role, u *User) bool
	KickUser(u User)
	sendAnonymousMessage(u User, message string) string
	//iniciativeTracker(sheet Sheet)Sheet
	//createSheet(sheet Sheet)Sheet
}
