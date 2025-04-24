package utils

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// hash password with bcrypt
func HashPassword(password string) ([]byte, error) {

	if password == "" {
		return nil, fmt.Errorf("password cannot be empty")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}
	return hash, nil
}

func VerifyPassword(password string, hash []byte) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(password))
}
