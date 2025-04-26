package models

import (
	"github.com/lib/pq"
	"time"
)

type NewOrder struct {
	Title       string   `json:"title"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Grade       string   `json:"grade"`
	MinPrice    int      `json:"min_price"`
	MaxPrice    int      `json:"max_price"`
	Tags        []string `json:"tags"`
}

type CreateOrder struct {
	StudentID string
	Order     *NewOrder
}

type UpdateOrder struct {
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Grade       string   `json:"grade,omitempty"`
	MinPrice    int      `json:"min_price,omitempty"`
	MaxPrice    int      `json:"max_price,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type OrderPagination struct {
	Orders []*Order
	Pages  int
}

type Order struct {
	ID            string         `json:"id" db:"id"`
	Name          string         `json:"name" db:"name"`
	StudentID     string         `json:"student_id" db:"student_id"`
	Title         string         `json:"title" db:"title"`
	Description   string         `json:"description" db:"description"`
	Grade         string         `json:"grade" db:"grade"`
	MinPrice      int            `json:"min_price" db:"min_price"`
	MaxPrice      int            `json:"max_price" db:"max_price"`
	Tags          pq.StringArray `json:"tags" db:"tags"`
	Status        string         `json:"status" db:"status"`
	ResponseCount int            `json:"response_count" db:"response_count"`
	CreatedAt     time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at" db:"updated_at"`
}

type OrderDetailsTutor struct {
	Order
	IsResponded bool `json:"is_responsed"`
}

type OrderDetails struct {
	Order
	Responses []Response `json:"responses"`
}

type OrderToBroker struct {
	ID        string   `json:"order_id"`
	StudentID int64    `json:"student_id"`
	Title     string   `json:"order_name"`
	Tags      []string `json:"tags"`
	Status    string   `json:"status"`
}

type SuggestOrder struct {
	ID              string `json:"order_id"`
	TutorTelegramID int64  `json:"tutor_telegram_id"`
	Title           string `json:"order_name"`
	Description     string `json:"description"`
	MinPrice        int    `json:"min_price"`
	MaxPrice        int    `json:"max_price"`
}
