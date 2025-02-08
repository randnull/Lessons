package service

import (
	"github.com/randnull/Lessons/internal/repository"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

type OrderServiceInt interface {
}

type OrderService struct {
	orderRepository repository.UserRepository
}

func NewOrderService(orderRepo repository.UserRepository) OrderServiceInt {
	return &OrderService{
		orderRepository: orderRepo,
	}
}
