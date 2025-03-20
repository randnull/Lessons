package service

import (
	"context"
	"fmt"
	"github.com/randnull/Lessons/internal/custom_errors"
	"github.com/randnull/Lessons/internal/gRPC_client"
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/internal/rabbitmq"
	"github.com/randnull/Lessons/internal/redis"
	"github.com/randnull/Lessons/internal/repository"
	"log"
	"time"
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
	RedisClient     redis.RedisInterface
}

func NewOrderService(orderRepo repository.OrderRepository, producerBroker rabbitmq.RabbitMQInterface, grpcClient gRPC_client.GRPCClientInt, redisClient redis.RedisInterface) OrderServiceInt {
	return &OrderService{
		orderRepository: orderRepo,
		ProducerBroker:  producerBroker,
		GRPCClient:      grpcClient,
		RedisClient:     redisClient,
	}
}

func (orderServ *OrderService) CreateOrder(order *models.NewOrder, UserData models.UserData) (string, error) {
	createdOrder, err := orderServ.orderRepository.CreateOrder(order, UserData.UserID, UserData.TelegramID)

	if err != nil {
		log.Printf("Error creating order: %v", err)
		return "", err
	}

	err = orderServ.RedisClient.SaveNewOrder(createdOrder)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	orderToBroker := &models.OrderToBrokerWithID{
		ID:          createdOrder.ID,
		StudentID:   UserData.TelegramID,
		Title:       createdOrder.Title,
		Description: createdOrder.Description,
		Tags:        createdOrder.Tags,
		MinPrice:    createdOrder.MinPrice,
		MaxPrice:    createdOrder.MaxPrice,
		ChatID:      UserData.TelegramID,
	}

	err = orderServ.ProducerBroker.Publish("my_queue", orderToBroker)
	if err != nil {
		log.Printf("Error publishing order: %v", err)
		// нужно что-то придумать .
	}

	return createdOrder.ID, nil
}

func (orderServ *OrderService) GetOrderById(id string, UserData models.UserData) (*models.OrderDetails, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	orderCached, errOrder := orderServ.RedisClient.GetOrderHSET(ctx, id)
	responses, errResponse := orderServ.RedisClient.GetAllResponses(ctx, id)

	if errOrder != nil || errResponse != nil || orderCached == nil || orderCached.ID != id || (responses == nil && orderCached.ResponseCount != 0) {
		log.Println(errOrder, errResponse)
		order, err := orderServ.orderRepository.GetByID(id)

		if err != nil {
			return nil, err
		}

		err = orderServ.RedisClient.AddOrder(order)

		if err != nil {
			return nil, err
		}

		if order.StudentID != UserData.UserID {
			return nil, custom_errors.ErrNotAllowed
		}

		return order, nil
	}

	fmt.Println("CACHED GetOrderById")

	orderDetail := &models.OrderDetails{
		ID:            orderCached.ID,
		StudentID:     orderCached.StudentID,
		Title:         orderCached.Title,
		Description:   orderCached.Description,
		Grade:         orderCached.Grade,
		MinPrice:      orderCached.MinPrice,
		MaxPrice:      orderCached.MaxPrice,
		Tags:          orderCached.Tags,
		Status:        orderCached.Status,
		ResponseCount: orderCached.ResponseCount,
		Responses:     responses,
		CreatedAt:     orderCached.CreatedAt,
		UpdatedAt:     orderCached.UpdatedAt,
	}

	log.Println("order", orderCached)
	return orderDetail, nil
}

func (orderServ *OrderService) GetOrderByIdTutor(id string, UserData models.UserData) (*models.OrderDetailsTutor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	orderCached, err := orderServ.RedisClient.GetOrderHSET(ctx, id)

	if err != nil || orderCached == nil || orderCached.ID != id {
		order, err := orderServ.orderRepository.GetByID(id)

		if err != nil {
			log.Println(err)
			return nil, err
		}

		return &models.OrderDetailsTutor{
			ID:            order.ID,
			Title:         order.Title,
			Description:   order.Description,
			Grade:         order.Grade,
			MinPrice:      order.MinPrice,
			MaxPrice:      order.MaxPrice,
			Tags:          order.Tags,
			Status:        order.Status,
			ResponseCount: order.ResponseCount,
			CreatedAt:     order.CreatedAt,
		}, nil
	}

	log.Println("CACHED GetOrderByIdTutor")

	return &models.OrderDetailsTutor{
		ID:            orderCached.ID,
		Title:         orderCached.Title,
		Description:   orderCached.Description,
		Grade:         orderCached.Grade,
		MinPrice:      orderCached.MinPrice,
		MaxPrice:      orderCached.MaxPrice,
		Tags:          orderCached.Tags,
		Status:        orderCached.Status,
		ResponseCount: orderCached.ResponseCount,
		CreatedAt:     orderCached.CreatedAt,
	}, nil
}

func (orderServ *OrderService) GetAllOrders(UserData models.UserData) ([]*models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	orders, err := orderServ.RedisClient.GetAllOrders(ctx)

	log.Println(orders, err)

	if err != nil || orders == nil {
		log.Println(err)
		return orderServ.orderRepository.GetAllOrders(UserData.UserID)
	}

	log.Println("CACHED GetAllOrders")

	return orders, nil
}

func (orderServ *OrderService) UpdateOrder(orderID string, order *models.UpdateOrder, UserData models.UserData) error {
	err := orderServ.orderRepository.UpdateOrder(orderID, order, UserData.UserID)

	if err != nil {
		return err
	}

	err = orderServ.RedisClient.UpdateOrderHSET(orderID, order)

	if err != nil {
		log.Println(err)
	}

	return nil
}

func (orderServ *OrderService) DeleteOrder(orderID string, UserData models.UserData) error {
	err := orderServ.orderRepository.DeleteOrder(orderID, UserData.UserID)

	if err != nil {
		return err
	}

	err = orderServ.RedisClient.DeleteOrderHSET(orderID)

	if err != nil {
		log.Println(err)
	}

	return nil
}

func (orderServ *OrderService) GetAllUsersOrders(UserData models.UserData) ([]*models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	orders, err := orderServ.RedisClient.GetAllUsersOrders(ctx, UserData.UserID)

	if err != nil || orders == nil {
		return orderServ.orderRepository.GetAllUsersOrders(UserData.UserID)
	}

	fmt.Println("CACHED GetAllUsersOrders")

	return orders, nil
}
