package service

import (
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/internal/repository"
)

type OrderServiceInt interface {
	CreateOrder(order *models.NewOrder) (string, error)
	GetOrderById(id string) (*models.Order, error)
	GetAllOrders() ([]*models.Order, error)
}

type OrderService struct {
	orderRepository repository.OrderRepository
}

func NewOrderService(orderRepo repository.OrderRepository) OrderServiceInt {
	return &OrderService{
		orderRepository: orderRepo,
	}
}

func (orderServ *OrderService) CreateOrder(order *models.NewOrder) (string, error) {
	return orderServ.orderRepository.CreateOrder(order)
}

func (orderServ *OrderService) GetOrderById(id string) (*models.Order, error) {
	return orderServ.orderRepository.GetByID(id)
}

func (orderServ *OrderService) GetAllOrders() ([]*models.Order, error) {
	return orderServ.orderRepository.GetAllOrders()
}
