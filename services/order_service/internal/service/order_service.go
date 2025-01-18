package service

import (
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/internal/rabbitmq"
	"github.com/randnull/Lessons/internal/repository"
	initdata "github.com/telegram-mini-apps/init-data-golang"
	"log"
)

type OrderServiceInt interface {
	CreateOrder(order *models.NewOrder, InitData initdata.InitData) (string, error)
	GetOrderById(id string, InitData initdata.InitData) (*models.Order, error)
	GetAllOrders(InitData initdata.InitData) ([]*models.Order, error)
	UpdateOrder(orderID string, order *models.NewOrder, InitData initdata.InitData) error
	DeleteOrder(orderID string, InitData initdata.InitData) error
}

type OrderService struct {
	orderRepository repository.OrderRepository
	ProducerBroker  rabbitmq.RabbitMQInterface
}

func NewOrderService(orderRepo repository.OrderRepository, producerBroker rabbitmq.RabbitMQInterface) OrderServiceInt {
	return &OrderService{
		orderRepository: orderRepo,
		ProducerBroker:  producerBroker,
	}
}

func (orderServ *OrderService) CreateOrder(order *models.NewOrder, InitData initdata.InitData) (string, error) {
	createdOrder, err := orderServ.orderRepository.CreateOrder(order, InitData)

	if err != nil {
		log.Printf("Error creating order: %v", err)
	}

	err = orderServ.ProducerBroker.Publish("my_queue", createdOrder)
	if err != nil {
		log.Printf("Error publishing order: %v", err)
		// нужно что-то придумать .
	}

	return createdOrder.ID, nil
}

func (orderServ *OrderService) GetOrderById(id string, InitData initdata.InitData) (*models.Order, error) {
	return orderServ.orderRepository.GetByID(id, InitData)
}

func (orderServ *OrderService) GetAllOrders(InitData initdata.InitData) ([]*models.Order, error) {
	return orderServ.orderRepository.GetAllOrders(InitData)
}

func (orderServ *OrderService) UpdateOrder(orderID string, order *models.NewOrder, InitData initdata.InitData) error {
	return orderServ.orderRepository.UpdateOrder(orderID, order, InitData)
}

func (orderServ *OrderService) DeleteOrder(orderID string, InitData initdata.InitData) error {
	return orderServ.orderRepository.DeleteOrder(orderID, InitData)
}
