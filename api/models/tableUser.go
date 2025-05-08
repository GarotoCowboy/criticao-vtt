package models

import "gorm.io/gorm"

type TableUser struct {
	gorm.Model
	TableID uint `gorm:"not null;uniqueIndex:idx_table_user_user"`
	UserID  uint `gorm:"not null;uniqueIndex:idx_table_user_user"`

	Role Role

	Table Table
	User  User
}
