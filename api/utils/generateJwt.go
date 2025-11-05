package utils

import (
	"time"

	"github.com/GarotoCowboy/vttProject/config"
	"github.com/golang-jwt/jwt/v5"
)

// Method generate a new JWT with userID
func GenerateJWT(userID uint) (string, error) {

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(config.JWT_SECRET)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
