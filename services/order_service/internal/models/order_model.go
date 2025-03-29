package models

import (
	"github.com/lib/pq"
	"time"
)

type NewOrder struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Grade       string   `json:"grade"`
	MinPrice    int      `json:"min_price"`
	MaxPrice    int      `json:"max_price"`
	Tags        []string `json:"tags"`
}

type UpdateOrder struct {
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Grade       string   `json:"grade,omitempty"`
	MinPrice    int      `json:"min_price,omitempty"`
	MaxPrice    int      `json:"max_price,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type Order struct {
	ID            string         `json:"id"`
	StudentID     string         `json:"student_id"`
	Title         string         `json:"title"`
	Description   string         `json:"description"`
	Grade         string         `json:"grade"`
	MinPrice      int            `json:"min_price"`
	MaxPrice      int            `json:"max_price"`
	Tags          pq.StringArray `json:"tags"`
	Status        string         `json:"status"`
	ResponseCount int            `json:"response_count"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type OrderPagination struct {
	Orders []*Order
	Pages  int
}

type OrderDetailsTutor struct {
	ID            string         `json:"id" db:"id"`
	Title         string         `json:"title" db:"title"`
	Description   string         `json:"description" db:"description"`
	Grade         string         `json:"grade" db:"grade"`
	MinPrice      int            `json:"min_price" db:"min_price"`
	MaxPrice      int            `json:"max_price" db:"max_price"`
	Tags          pq.StringArray `json:"tags" db:"tags"`
	Status        string         `json:"status" db:"status"`
	IsResponsed   bool           `json:"is_responsed"`
	ResponseCount int            `json:"response_count" db:"response_count"`
	CreatedAt     time.Time      `json:"created_at" db:"created_at"`
}

type OrderDetails struct {
	ID            string         `json:"id"`
	StudentID     string         `json:"student_id"`
	Title         string         `json:"title"`
	Description   string         `json:"description"`
	Grade         string         `json:"grade"`
	MinPrice      int            `json:"min_price"`
	MaxPrice      int            `json:"max_price"`
	Tags          pq.StringArray `json:"tags"`
	Status        string         `json:"status"`
	ResponseCount int            `json:"response_count"`
	Responses     []Response     `json:"responses"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type OrderToBrokerWithID struct {
	ID          string   `json:"id"`
	StudentID   int64    `json:"student_id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	MinPrice    int      `json:"min_price"`
	MaxPrice    int      `json:"max_price"`
	Tags        []string `json:"tags"`
	ChatID      int64    `json:"chat_id"`
	Status      string   `json:"status"`
}
