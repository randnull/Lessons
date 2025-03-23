package models

import "time"

type User struct {
	Id   string `json:"user_id"`
	Name string `json:"name"`
	//Univer  string `json:"univer"`
	//AboutMe string `json:"info"`
}

type UserDB struct {
	Id         string    `json:"id"`
	TelegramID int64     `json:"telegram_id"`
	Name       string    `json:"name"`
	Role       string    `json:"role"`
	CreatedAt  time.Time `json:"created_at"`
}

type TutorDB struct {
	Id         string    `json:"id"`
	TelegramID int64     `json:"telegram_id"`
	Name       string    `json:"name"`
	Role       string    `json:"role"`
	Bio        string    `json:"bio"`
	CreatedAt  time.Time `json:"created_at"`
}

type CreateUser struct {
	Name       string `json:"name"`
	TelegramId int64  `json:"telegram_id"`
	Role       string `json:"role"`
}
