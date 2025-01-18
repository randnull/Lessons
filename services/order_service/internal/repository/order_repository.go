package repository

import (
	"github.com/randnull/Lessons/internal/models"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

type OrderRepository interface {
	CreateOrder(order *models.NewOrder, InitData initdata.InitData) (models.OrderToBrokerWithID, error)
	GetByID(id string, InitData initdata.InitData) (*models.Order, error)
	GetAllOrders(InitData initdata.InitData) ([]*models.Order, error)
	UpdateOrder(orderID string, order *models.NewOrder, InitData initdata.InitData) error
	DeleteOrder(id string, InitData initdata.InitData) error
}
