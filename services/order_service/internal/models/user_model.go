package models

type CreateUser struct {
	Name       string `json:"name"`
	TelegramId int64  `json:"telegram_id"`
}

type User struct {
	Id         string `json:"id"`
	TelegramID int64  `json:"telegram_id"`
	Name       string `json:"name"`
	//CreatedAt  time.Time `json:"created_at"`
}
