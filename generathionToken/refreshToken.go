package generathionToken

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"main.go/db"
)

// GenerateRefreshToken создаёт уникальный refresh token и сохраняет его в БД
func GenerateRefreshToken(phone string) (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	token := hex.EncodeToString(bytes)

	_, err := db.DB.Exec(`INSERT INTO refresh_tokens (phone, token, created_at) VALUES ($1, $2, $3)`, phone, token, time.Now())
	if err != nil {
		return "", err
	}
	return token, nil
}

// ValidateRefreshToken проверяет, что токен существует
func ValidateRefreshToken(token string) (string, error) {
	var phone string
	err := db.DB.QueryRow(`SELECT phone FROM refresh_tokens WHERE token = $1`, token).Scan(&phone)
	if err != nil {
		return "", err
	}
	return phone, nil
}

// DeleteRefreshToken — используется при logout
func DeleteRefreshToken(token string) {
	db.DB.Exec(`DELETE FROM refresh_tokens WHERE token = $1`, token)
}
