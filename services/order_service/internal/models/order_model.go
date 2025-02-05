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

type UpdateOrder struct {
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	MinPrice    int      `json:"min_price,omitempty"`
	MaxPrice    int      `json:"max_price,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type Order struct {
	ID          string         `json:"id"`
	StudentID   int            `json:"student_id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	MinPrice    int            `json:"min_price"`
	MaxPrice    int            `json:"max_price"`
	Tags        pq.StringArray `json:"tags"`
	Status      string         `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type OrderToBrokerWithID struct {
	ID          string   `json:"id"`
	StudentID   int      `json:"student_id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	MinPrice    int      `json:"min_price"`
	MaxPrice    int      `json:"max_price"`
	Tags        []string `json:"tags"`
	ChatID      int64    `json:"chat_id"`
	Status      string   `json:"status"`
}
