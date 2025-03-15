package models

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID     string `json:"user_id"`
	TelegramID int64  `json:"telegram_id"`
	Role       string `json:"role"`
	jwt.RegisteredClaims
}
