package models

import (
	"github.com/lib/pq"
	"time"
)

type NewOrder struct {
	Title       string   `json:"title" validate:"required,min=5,max=100"`
	Name        string   `json:"name" validate:"required,min=2,max=50"`
	Description string   `json:"description" validate:"required,min=10,max=1000"`
	Grade       string   `json:"grade" validate:"required,min=1,max=20"`
	MinPrice    int      `json:"min_price" validate:"gte=0"`
	MaxPrice    int      `json:"max_price" validate:"gte=0"`
	Tags        []string `json:"tags" validate:"required,min=1,max=10,dive,required"`
}

type CreateOrder struct {
	StudentID string
	Order     *NewOrder
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
type UpdateOrder struct {
	Title       string `json:"title,omitempty" validate:"omitempty,min=5,max=100"`
	Description string `json:"description,omitempty" validate:"omitempty,min=10,max=1000"`
	Grade       string `json:"grade,omitempty" validate:"omitempty,min=1,max=20"`
}

type OrderPagination struct {
	Orders []*Order
	Pages  int
}

type OrderDetailsTutor struct {
	Order
	IsResponded bool `json:"is_responsed"`
}

type OrderDetails struct {
	Order
	Responses []Response `json:"responses"`
}

type ChangeActive struct {
	IsActive bool `json:"is_active"`
}
