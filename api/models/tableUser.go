package models

import "gorm.io/gorm"

// TableUser representa um relacionamento de associação entre usuários e tabelas
// @description Relacionamento entre usuários e tabelas que eles pertencem.
// @type TableUser
type TableUser struct {
	gorm.Model
	TableID uint `gorm:"not null;uniqueIndex:idx_table_user_user"`
	UserID  uint `gorm:"not null;uniqueIndex:idx_table_user_user"`

	Role Role

	Table Table
	User  User
}
