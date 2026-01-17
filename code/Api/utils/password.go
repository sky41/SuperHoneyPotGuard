package utils

import (
	"golang.org/x/crypto/bcrypt"
	"superhoneypotguard/config"
)

func HashPassword(password string) (string, error) {
	cost := config.AppConfig.BCryptCost
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}

func ComparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
