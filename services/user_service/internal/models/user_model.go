package models

import "time"

type User struct {
	Id   string `json:"user_id"`
	Name string `json:"name"`
}

type UserDB struct {
	Id         string    `json:"id"`
	TelegramID int64     `json:"telegram_id"`
	Name       string    `json:"name"`
	Role       string    `json:"role"`
	CreatedAt  time.Time `json:"created_at"`
}

type TutorDB struct {
	Id             string
	TelegramID     int64
	Name           string
	Role           string
	CreatedAt      time.Time
	Bio            string
	ResponseCount  int32
	Tags           []string
	IsActive       bool
	TutorCreatedAt time.Time
}

type CreateUser struct {
	Name       string `json:"name"`
	TelegramId int64  `json:"telegram_id"`
	Role       string `json:"role"`
}

type TutorDetails struct {
	Tutor         TutorDB
	ResponseCount int32
	Reviews       []Review
	Tags          []string
}
