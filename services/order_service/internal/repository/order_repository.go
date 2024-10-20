package repository

import (
	"github.com/randnull/Lessons/internal/models"
)

type OrderRepository interface {
	CreateOrder(order *models.NewOrder) (string, error)
	GetByID(id string) (*models.Order, error)
	GetAllOrders() ([]*models.Order, error)
}
