package auth

import (
	"crypto/rand"
	"encoding/base64"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		slog.Error("error hashing password", "error", err.Error())
	}

	return string(hashedPassword)
}

func ComparePassword(password string, hashedPassword string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		slog.Error("passwords do not match", "error", err.Error())
		return false
	}
	return true
}

func GenerateToken(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		slog.Error("error creating bcrypt bytes", "error", err.Error())
	}

	return base64.URLEncoding.EncodeToString(bytes)
}
