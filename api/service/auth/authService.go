package auth

import (
	"errors"
	"time"

	"github.com/GarotoCowboy/vttProject/api/models"
	"github.com/GarotoCowboy/vttProject/api/utils"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func LoginService(db *gorm.DB, req LoginRequest) (string, time.Time, error) {
	var user models.User

	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", time.Time{}, errors.New("invalid credentials")
		}
		return "", time.Time{}, err
	}

	if err := utils.VerifyPassword(req.Password, []byte(user.Password)); err != nil {
		return "", time.Time{}, errors.New("invalid credentials")
	}

	token, expiresAt, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return "", time.Time{}, errors.New("failed to generate JWT")
	}
	return token, expiresAt, nil
}
