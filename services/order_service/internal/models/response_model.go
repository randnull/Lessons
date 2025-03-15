package models

import "time"

type NewResponseModel struct {
	OrderId string `json:"order_id"`
}

type Response struct {
	ID        string    `json:"id" db:"id"`
	TutorID   string    `json:"tutor_id" db:"tutor_id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type ResponseDB struct {
	ID        string    `json:"id" db:"id"`
	OrderID   string    `json:"order_id" db:"order_id"`
	TutorID   string    `json:"tutor_id" db:"tutor_id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type ResponseToBrokerModel struct {
	UserId  int64  `json:"user_id"`
	OrderId string `json:"order_id"`
	ChatId  int64  `json:"chat_id"`
}
