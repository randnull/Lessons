package service

import (
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/internal/repository"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

type OrderServiceInt interface {
	CreateOrder(order *models.NewOrder, InitData initdata.InitData) (string, error)
	GetOrderById(id string) (*models.Order, error)
	GetAllOrders() ([]*models.Order, error)
	UpdateOrder(order *models.Order) error
	DeleteOrder(id string) error
}

type OrderService struct {
	orderRepository repository.OrderRepository
}

func NewOrderService(orderRepo repository.OrderRepository) OrderServiceInt {
	return &OrderService{
		orderRepository: orderRepo,
	}
}

func (orderServ *OrderService) CreateOrder(order *models.NewOrder, InitData initdata.InitData) (string, error) {
	return orderServ.orderRepository.CreateOrder(order, InitData)
}

func (orderServ *OrderService) GetOrderById(id string) (*models.Order, error) {
	return orderServ.orderRepository.GetByID(id)
}

func (orderServ *OrderService) GetAllOrders() ([]*models.Order, error) {
	return orderServ.orderRepository.GetAllOrders()
}

func (orderServ *OrderService) UpdateOrder(order *models.Order) error {
	return orderServ.orderRepository.UpdateOrder(order)
}

func (orderServ *OrderService) DeleteOrder(id string) error {
	return orderServ.orderRepository.DeleteOrder(id)
}
