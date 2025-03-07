package service

import (
	"github.com/randnull/Lessons/internal/gRPC_client"
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/internal/rabbitmq"
	"github.com/randnull/Lessons/internal/repository"
	initdata "github.com/telegram-mini-apps/init-data-golang"
	"log"
)

type OrderServiceInt interface {
	CreateOrder(order *models.NewOrder, InitData initdata.InitData) (string, error)
	GetOrderById(id string, InitData initdata.InitData) (*models.OrderDetails, error)
	GetOrderByIdTutor(id string, InitData initdata.InitData) (*models.OrderDetailsTutor, error)
	GetAllOrders(InitData initdata.InitData) ([]*models.Order, error)
	GetAllUsersOrders(InitData initdata.InitData) ([]*models.Order, error)
	UpdateOrder(orderID string, order *models.UpdateOrder, InitData initdata.InitData) error
	DeleteOrder(orderID string, InitData initdata.InitData) error
}

type OrderService struct {
	orderRepository repository.OrderRepository
	ProducerBroker  rabbitmq.RabbitMQInterface
	GRPCClient      gRPC_client.GRPCClientInt
}

func NewOrderService(orderRepo repository.OrderRepository, producerBroker rabbitmq.RabbitMQInterface, grpcClient gRPC_client.GRPCClientInt) OrderServiceInt {
	return &OrderService{
		orderRepository: orderRepo,
		ProducerBroker:  producerBroker,
		GRPCClient:      grpcClient,
	}
}

func (orderServ *OrderService) CreateOrder(order *models.NewOrder, InitData initdata.InitData) (string, error) {
	createdOrder, err := orderServ.orderRepository.CreateOrder(order, InitData)

	if err != nil {
		log.Printf("Error creating order: %v", err)
		return "", err
	}

	err = orderServ.ProducerBroker.Publish("my_queue", createdOrder)
	if err != nil {
		log.Printf("Error publishing order: %v", err)
		return createdOrder.ID, nil
		// нужно что-то придумать .
	}

	return createdOrder.ID, nil
}

func (orderServ *OrderService) GetOrderById(id string, InitData initdata.InitData) (*models.OrderDetails, error) {
	return orderServ.orderRepository.GetByID(id, InitData)
}

func (orderServ *OrderService) GetOrderByIdTutor(id string, InitData initdata.InitData) (*models.OrderDetailsTutor, error) {
	return orderServ.orderRepository.GetOrderByIdTutor(id, InitData)
}

func (orderServ *OrderService) GetAllOrders(InitData initdata.InitData) ([]*models.Order, error) {
	return orderServ.orderRepository.GetAllOrders(InitData)
}

func (orderServ *OrderService) UpdateOrder(orderID string, order *models.UpdateOrder, InitData initdata.InitData) error {
	return orderServ.orderRepository.UpdateOrder(orderID, order, InitData)
}

func (orderServ *OrderService) DeleteOrder(orderID string, InitData initdata.InitData) error {
	return orderServ.orderRepository.DeleteOrder(orderID, InitData)
}

func (orderServ *OrderService) GetAllUsersOrders(InitData initdata.InitData) ([]*models.Order, error) {
	return orderServ.orderRepository.GetAllUsersOrders(InitData)
}
