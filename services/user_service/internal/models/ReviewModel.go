package models

import "time"

type Review struct {
	ID        string    `json:"id" db:"id"`
	TutorID   string    `json:"tutor_id" db:"tutor_id"`
	OrderID   string    `json:"order_id" db:"order_id"`
	Rating    int       `json:"rating" db:"rating"`
	Comment   string    `json:"comment" db:"comment"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
