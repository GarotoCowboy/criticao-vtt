package models

import "gorm.io/gorm"

type GameObjectOwner struct {
	gorm.Model
	UserId uint `json:"user_id" gorm:"not null;uniqueIndex:idx_user_ptoken;uniqueIndex:idx_user_pimage"`
	User   User `gorm:"constraint:OnUpdate:CASCADE;constraint:OnDelete:CASCADE"`

	PlacedTokenID *uint `json:"placed_token_id" gorm:"uniqueIndex:idx_user_ptoken"`
	PlacedImageID *uint `json:"placed_image_id" gorm:"uniqueIndex:idx_user_pimage"`
}
