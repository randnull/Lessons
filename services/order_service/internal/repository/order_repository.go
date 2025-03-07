package repository

import (
	"github.com/randnull/Lessons/internal/models"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

type OrderRepository interface {
	CreateOrder(order *models.NewOrder, InitData initdata.InitData) (*models.OrderToBrokerWithID, error)
	GetByID(id string, InitData initdata.InitData) (*models.OrderDetails, error)
	GetOrderByIdTutor(id string, InitData initdata.InitData) (*models.OrderDetailsTutor, error)
	GetAllOrders(InitData initdata.InitData) ([]*models.Order, error)
	UpdateOrder(orderID string, order *models.UpdateOrder, InitData initdata.InitData) error
	GetUserByOrder(orderID string) (*int64, error)
	GetAllUsersOrders(InitData initdata.InitData) ([]*models.Order, error)
	DeleteOrder(id string, InitData initdata.InitData) error
	CreateResponse(response *models.NewResponseModel, Tutor *models.User) (string, error)
}
