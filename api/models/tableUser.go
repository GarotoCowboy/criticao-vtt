package models

import (
	"github.com/GarotoCowboy/vttProject/api/models/consts"
	"gorm.io/gorm"
)

// TableUser representa um relacionamento de associação entre usuários e tabelas
// @description Relacionamento entre usuários e tabelas que eles pertencem.
// @type TableUser
type TableUser struct {
	gorm.Model
	TableID uint `gorm:"not null;uniqueIndex:idx_table_user_user"`
	UserID  uint `gorm:"not null;uniqueIndex:idx_table_user_user"`

	Role consts.Role

	Table Table `gorm:"constraint:OnDelete:CASCADE"`
	User  User  `gorm:"constraint:OnDelete:CASCADE"`
}
