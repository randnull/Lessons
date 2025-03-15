package service

import (
	"github.com/randnull/Lessons/internal/gRPC_client"
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/internal/rabbitmq"
	"github.com/randnull/Lessons/internal/repository"
	"log"
)

type OrderServiceInt interface {
	CreateOrder(order *models.NewOrder, UserData models.UserData) (string, error)
	GetOrderById(id string, UserData models.UserData) (*models.OrderDetails, error)
	GetOrderByIdTutor(id string, UserData models.UserData) (*models.OrderDetailsTutor, error)
	GetAllOrders(UserData models.UserData) ([]*models.Order, error)
	GetAllUsersOrders(UserData models.UserData) ([]*models.Order, error)
	UpdateOrder(orderID string, order *models.UpdateOrder, UserData models.UserData) error
	DeleteOrder(orderID string, UserData models.UserData) error
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

func (orderServ *OrderService) CreateOrder(order *models.NewOrder, UserData models.UserData) (string, error) {
	//UserInfo, err := orderServ.GRPCClient.GetUserByTelegramID(context.Background(), UserData.TelegramID)

	//if err != nil {
	//	return "", err
	//}

	createdOrder, err := orderServ.orderRepository.CreateOrder(order, UserData.UserID, UserData.TelegramID)

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

func (orderServ *OrderService) GetOrderById(id string, UserData models.UserData) (*models.OrderDetails, error) {
	//UserInfo, err := orderServ.GRPCClient.GetUserByTelegramID(context.Background(), UserData.TelegramID)
	//
	//if err != nil {
	//	return nil, err
	//}

	return orderServ.orderRepository.GetByID(id, UserData.UserID)
}

func (orderServ *OrderService) GetOrderByIdTutor(id string, UserData models.UserData) (*models.OrderDetailsTutor, error) {
	//UserInfo, err := orderServ.GRPCClient.GetUserByTelegramID(context.Background(), UserData.TelegramID)
	//
	//if err != nil {
	//	return nil, err
	//}

	return orderServ.orderRepository.GetOrderByIdTutor(id, UserData.UserID)
}

func (orderServ *OrderService) GetAllOrders(UserData models.UserData) ([]*models.Order, error) {
	//UserInfo, err := orderServ.GRPCClient.GetUserByTelegramID(context.Background(), UserData.TelegramID)
	//
	//if err != nil {
	//	return nil, err
	//}

	return orderServ.orderRepository.GetAllOrders(UserData.UserID)
}

func (orderServ *OrderService) UpdateOrder(orderID string, order *models.UpdateOrder, UserData models.UserData) error {
	//UserInfo, err := orderServ.GRPCClient.GetUserByTelegramID(context.Background(), UserData.TelegramID)
	//
	//if err != nil {
	//	return err
	//}

	return orderServ.orderRepository.UpdateOrder(orderID, order, UserData.UserID)
}

func (orderServ *OrderService) DeleteOrder(orderID string, UserData models.UserData) error {
	//UserInfo, err := orderServ.GRPCClient.GetUserByTelegramID(context.Background(), UserData.TelegramID)
	//
	//if err != nil {
	//	return err
	//}

	return orderServ.orderRepository.DeleteOrder(orderID, UserData.UserID)
}

func (orderServ *OrderService) GetAllUsersOrders(UserData models.UserData) ([]*models.Order, error) {
	//UserInfo, err := orderServ.GRPCClient.GetUserByTelegramID(context.Background(), UserData.TelegramID)
	//
	//if err != nil {
	//	return nil, err
	//}

	return orderServ.orderRepository.GetAllUsersOrders(UserData.UserID)
}
