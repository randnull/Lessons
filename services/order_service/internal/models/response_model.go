package models

import "time"

type NewResponseModel struct {
	Greetings string `json:"greetings"`
}

type Response struct {
	ID        string    `json:"id" db:"id"`
	OrderID   string    `json:"order_id" db:"order_id"`
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
	ResponseID string `json:"response_id"`
	TutorID    int64  `json:"tutor_id"`
	StudentID  int64  `json:"student_id"`
	OrderID    string `json:"order_id"`
	Title      string `json:"order_name"`
}

type SelectedResponseToBroker struct {
	ResponseID string `json:"response_id"`
}
