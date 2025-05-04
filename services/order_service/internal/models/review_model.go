package models

import "time"

type Review struct {
	ID        string    `json:"id"`
	TutorID   string    `json:"tutor_id"`
	OrderID   string    `json:"order_id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type ReviewRequest struct {
	ResponseID string `json:"response_id" validate:"required,uuid"`
	Comment    string `json:"comment" validate:"required,min=5,max=500"`
	Rating     int    `json:"rating" validate:"required,gte=1,lte=5"`
}

type ReviewToBroker struct {
	ReviewID        string `json:"review_id"`
	ResponseID      string `json:"response_id"`
	OrderID         string `json:"order_id"`
	OrderName       string `json:"order_name"`
	TutorTelegramID int64  `json:"tutor_telegram_id"`
}

type ReviewActive struct {
	ReviewID string `json:"review_id" validate:"required,uuid"`
}
