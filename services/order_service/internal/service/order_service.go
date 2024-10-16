package service

import (
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/internal/repository"
)

type OrderServiceInt interface {
	CreateOrder(order *models.Order) error
	GetOrderById(id string) (*models.Order, error)
}

type OrderService struct {
	orderRepository repository.OrderRepository
}

func NewOrderService(orderRepo repository.OrderRepository) OrderServiceInt {
	return &OrderService{
		orderRepository: orderRepo,
	}
}

func (orderServ *OrderService) CreateOrder(order *models.Order) error {
	return orderServ.orderRepository.CreateOrder(order)
}

func (orderServ *OrderService) GetOrderById(id string) (*models.Order, error) {
	return orderServ.orderRepository.GetByID(id)
}
