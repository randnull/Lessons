package models

type User struct {
	Id         string `json:"id"`
	TelegramID int64  `json:"telegram_id"`
	Name       string `json:"name"`
	IsBanned   bool   `json:"is_banned"`
	Role       string `json:"role"`
}

type NewUser struct {
	TelegramID int64  `json:"telegram_id"`
	Username   string `json:"username"`
	Name       string `json:"name"`
	Role       string `json:"role"`
}
