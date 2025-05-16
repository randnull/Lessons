package models

import "github.com/golang-jwt/jwt/v5"

type UserData struct {
	UserID     string `json:"user_id"`
	TelegramID int64  `json:"telegram_id"`
	Username   string `json:"username"`
	Role       string `json:"role"`
}

type Claims struct {
	UserID     string `json:"user_id"`
	Username   string `json:"username"`
	TelegramID int64  `json:"telegram_id"`
	Role       string `json:"role"`
	jwt.RegisteredClaims
}
