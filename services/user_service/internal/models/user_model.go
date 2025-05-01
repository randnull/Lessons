package models

import "time"

type User struct {
	Id   string `json:"user_id" db:"user_id"`
	Name string `json:"name" db:"name"`
}

type TutorWithResponse struct {
	Id            string `json:"user_id" db:"user_id"`
	Name          string `json:"name" db:"name"`
	TelegramID    int64  `json:"telegram_id" db:"telegram_id"`
	ResponseCount int32  `json:"response_count" db:"response_count"`
}

type UserDB struct {
	Id         string    `json:"id" db:"id"`
	TelegramID int64     `json:"telegram_id" db:"telegram_id"`
	Name       string    `json:"name" db:"name"`
	Role       string    `json:"role" db:"role"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

type TutorDB struct {
	Id             string    `db:"id"`
	TelegramID     int64     `db:"telegram_id"`
	Name           string    `db:"name"`
	Role           string    `db:"role"`
	CreatedAt      time.Time `db:"created_at"`
	Bio            string    `db:"bio"`
	ResponseCount  int32     `db:"response_count"`
	Tags           []string  `db:"tags"`
	IsActive       bool      `db:"is_active"`
	TutorCreatedAt time.Time `db:"tutor_created_at"`
}

type CreateUser struct {
	Name       string `json:"name" db:"name"`
	TelegramId int64  `json:"telegram_id" db:"telegram_id"`
	Role       string `json:"role" db:"role"`
}

type TutorDetails struct {
	Tutor         TutorDB
	ResponseCount int32
	Reviews       []Review
	Tags          []string
}
