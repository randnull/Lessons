package models

type User struct {
	Id         string   `json:"id"`
	TelegramID int64    `json:"telegram_id"`
	Name       string   `json:"name"`
	Role       RoleType `json:"role"`
}

type NewUser struct {
	TelegramID int64    `json:"telegram_id"`
	Name       string   `json:"name"`
	Role       RoleType `json:"role"`
}
