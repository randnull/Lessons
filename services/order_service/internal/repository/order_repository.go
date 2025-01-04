package repository

import (
	"github.com/randnull/Lessons/internal/models"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

type OrderRepository interface {
	CreateOrder(order *models.NewOrder, InitData initdata.InitData) (string, error)
	GetByID(id string) (*models.Order, error)
	GetAllOrders() ([]*models.Order, error)
	UpdateOrder(order *models.Order) error
	DeleteOrder(id string) error
}
