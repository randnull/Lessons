package models

import "time"

type NewResponseModel struct {
	Greetings string `json:"greetings"`
}

type Response struct {
	ID        string    `json:"id" db:"id"`
	TutorID   string    `json:"tutor_id" db:"tutor_id"`
	Name      string    `json:"name" db:"name"`
	IsFinal   bool      `json:"is_final" db:"is_final"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type ResponseDB struct {
	ID            string    `json:"id" db:"id"`
	OrderID       string    `json:"order_id" db:"order_id"`
	TutorID       string    `json:"tutor_id" db:"tutor_id"`
	TutorUsername string    `json:"tutor_username" db:"tutor_username"`
	Name          string    `json:"name" db:"name"`
	Greetings     string    `json:"greetings" db:"greetings"`
	IsFinal       bool      `json:"is_final" db:"is_final"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

type ResponseToBrokerModel struct {
	UserId  int64  `json:"user_id"`
	OrderId string `json:"order_id"`
	ChatId  int64  `json:"chat_id"`
}
