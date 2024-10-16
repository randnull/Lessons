package repository

import (
	"github.com/randnull/Lessons/internal/models"
)

type OrderRepository interface {
	CreateOrder(order *models.Order) error
	GetByID(id string) (*models.Order, error)
}
