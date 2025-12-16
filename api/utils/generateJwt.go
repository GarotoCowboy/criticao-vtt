package utils

import (
	"time"

	"github.com/GarotoCowboy/vttProject/config"
	"github.com/golang-jwt/jwt/v5"
)

// Method generate a new JWT with userID
func GenerateJWT(userID uint) (string, time.Time, error) {

	expiresAt := time.Now().Add(time.Hour * 24)

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expiresAt.Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(config.JWT_SECRET)
	if err != nil {
		return "", time.Time{}, err
	}
	return tokenString, expiresAt, nil
}
