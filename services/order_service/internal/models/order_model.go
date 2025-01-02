package models

import "time"

type NewOrder struct {
	Subject     string `json:"subject"`
	Description string `json:"description"`
	MinPrice    int    `json:"min_price"`
	MaxPrice    int    `json:"max_price"`
}

type Order struct {
	ID          string    `json:"id"`
	StudentID   string    `json:"student_id"`
	TutorID     string    `json:"tutor_id"`
	Subject     string    `json:"subject"`
	Description string    `json:"description"`
	MinPrice    int       `json:"min_price"`
	MaxPrice    int       `json:"max_price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
