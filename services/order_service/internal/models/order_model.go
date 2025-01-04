package models

import (
	"github.com/lib/pq"
	"time"
)

type NewOrder struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	MinPrice    int      `json:"min_price"`
	MaxPrice    int      `json:"max_price"`
	Tags        []string `json:"tags"`
}

type Order struct {
	ID        string `json:"id"`
	StudentID int    `json:"student_id"`
	//TutorID     string    	`json:"tutor_id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	MinPrice    int            `json:"min_price"`
	MaxPrice    int            `json:"max_price"`
	Tags        pq.StringArray `json:"tags"`
	Status      string         `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
