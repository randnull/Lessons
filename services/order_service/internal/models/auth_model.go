package models

type UserData struct {
	UserID     string `json:"user_id"`
	TelegramID int64  `json:"telegram_id"`
	Username   string `json:"username"`
	Role       string `json:"role"`
}
