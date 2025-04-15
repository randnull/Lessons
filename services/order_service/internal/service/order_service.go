package service

import (
	"context"
	"errors"
	"github.com/randnull/Lessons/internal/custom_errors"
	"github.com/randnull/Lessons/internal/gRPC_client"
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/internal/rabbitmq"
	"github.com/randnull/Lessons/internal/repository"
	"log"
)

type OrderServiceInt interface {
	CreateOrder(order *models.NewOrder, UserData models.UserData) (string, error)
	GetOrderById(id string, UserData models.UserData) (*models.OrderDetails, error)
	GetStudentOrdersWithPagination(page int, size int, UserData models.UserData) (*models.OrderPagination, error)
	GetOrdersWithPagination(page int, size int, tag string, UserData models.UserData) (*models.OrderPagination, error)
	GetOrderByIdTutor(id string, UserData models.UserData) (*models.OrderDetailsTutor, error)
	GetAllOrders(UserData models.UserData) ([]*models.Order, error)
	GetAllUsersOrders(UserData models.UserData) ([]*models.Order, error)
	UpdateOrder(orderID string, order *models.UpdateOrder, UserData models.UserData) error
	DeleteOrder(orderID string, UserData models.UserData) error
	SelectTutor(responseID string, UserData models.UserData) error
	ApproveSelectedOrderByTutor(responseID string, UserData models.UserData) error
	SetActiveOrderStatus(orderID string, IsActive bool, UserData models.UserData) error
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
	_, err := orderServ.GRPCClient.GetStudent(context.Background(), UserData.UserID)

	if err != nil {
		return "", custom_errors.ErrorGetUser
	}

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
	_, err := orderServ.orderRepository.CheckOrderByStudentID(id, UserData.UserID)
	if err != nil {
		return nil, err
	}

	return orderServ.orderRepository.GetOrderByID(id)
}

func (orderServ *OrderService) GetOrderByIdTutor(id string, UserData models.UserData) (*models.OrderDetailsTutor, error) {
	return orderServ.orderRepository.GetOrderByIdTutor(id, UserData.UserID)
}

func (orderServ *OrderService) GetOrdersWithPagination(page int, size int, tag string, UserData models.UserData) (*models.OrderPagination, error) {
	limit := size
	offset := (page - 1) * size

	orders, count, err := orderServ.orderRepository.GetOrdersPagination(limit, offset, tag)
	if err != nil {
		return nil, err
	}

	addPage := 0

	if count%size != 0 {
		addPage += 1
	}

	return &models.OrderPagination{
		Orders: orders,
		Pages:  count/size + addPage,
	}, nil
}

func (orderServ *OrderService) GetStudentOrdersWithPagination(page int, size int, UserData models.UserData) (*models.OrderPagination, error) {
	limit := size
	offset := (page - 1) * size

	orders, count, err := orderServ.orderRepository.GetStudentOrdersPagination(limit, offset, UserData.UserID)

	if err != nil {
		return nil, err
	}

	addPage := 0

	if count%size != 0 {
		addPage += 1
	}

	return &models.OrderPagination{
		Orders: orders,
		Pages:  count/size + addPage,
	}, nil
}

func (orderServ *OrderService) GetAllOrders(UserData models.UserData) ([]*models.Order, error) {
	return orderServ.orderRepository.GetOrders()
}

func (orderServ *OrderService) UpdateOrder(orderID string, order *models.UpdateOrder, UserData models.UserData) error {
	isExist, err := orderServ.orderRepository.CheckOrderByStudentID(orderID, UserData.UserID)

	if !isExist || err != nil {
		return custom_errors.ErrNotAllowed
	}

	return orderServ.orderRepository.UpdateOrder(orderID, order, UserData.UserID)
}

func (orderServ *OrderService) DeleteOrder(orderID string, UserData models.UserData) error {
	isExist, err := orderServ.orderRepository.CheckOrderByStudentID(orderID, UserData.UserID)

	if !isExist || err != nil {
		return custom_errors.ErrNotAllowed
	}

	return orderServ.orderRepository.DeleteOrder(orderID)
}

func (orderServ *OrderService) GetAllUsersOrders(UserData models.UserData) ([]*models.Order, error) {
	return orderServ.orderRepository.GetStudentOrders(UserData.UserID)
}

func (orderServ *OrderService) SelectTutor(responseID string, UserData models.UserData) error {
	response, err := orderServ.orderRepository.GetResponseById(responseID)

	isUserRequest, err := orderServ.orderRepository.CheckOrderByStudentID(responseID, UserData.UserID)

	if err != nil {
		return err
	}

	if !isUserRequest {
		return custom_errors.ErrNotAllowed
	}

	log.Println(response, err)

	if err != nil || response == nil {
		return custom_errors.ErrNotAllowed
	}

	isExist, err := orderServ.orderRepository.CheckOrderByStudentID(response.OrderID, UserData.UserID)

	log.Println(isExist, err)

	if err != nil || !isExist {
		return custom_errors.ErrNotAllowed
	}

	err = orderServ.orderRepository.SetTutorToOrder(response, UserData)

	if err != nil {
		return errors.New("error with select tutor")
	}

	err = orderServ.ProducerBroker.Publish("selected_orders", models.SelectedResponseToBroker{
		ResponseID: responseID,
	})

	if err != nil {
		log.Printf("error with push response selected to broker")
	}

	return nil
}

func (orderServ *OrderService) ApproveSelectedOrderByTutor(responseID string, UserData models.UserData) error {
	response, err := orderServ.orderRepository.GetResponseById(responseID)

	if err != nil {
		return errors.New("error with get response")
	}

	if response.TutorID != UserData.UserID {
		return errors.New("error not allowed")
	}

	if !response.IsFinal {
		return errors.New("error not final") // тут проверяем, что репетитора правда выбрали
	}

	err = orderServ.orderRepository.SetOrderStatus("Approved", response.OrderID)

	if err != nil {
		return errors.New("error with approved order")
	}

	err = orderServ.ProducerBroker.Publish("approved_orders", models.SelectedResponseToBroker{
		ResponseID: responseID,
	})

	if err != nil {
		log.Println("error with publish to broker")
	}

	return nil
}

func (orderServ *OrderService) SetActiveOrderStatus(orderID string, IsActive bool, UserData models.UserData) error {
	order, err := orderServ.orderRepository.GetOrderByID(orderID)

	if err != nil {
		return errors.New("error with get order")
	}

	if order.StudentID != UserData.UserID {
		return errors.New("error not allowed")
	}

	if IsActive {
		if order.Status != "Inactive" {
			return errors.New("error not Inactive state")
		}

		err = orderServ.orderRepository.SetOrderStatus("New", orderID)

		if err != nil {
			return errors.New("error with NEW state order")
		}
	} else {
		if order.Status != "New" {
			return errors.New("error not NEW state")
		}

		err = orderServ.orderRepository.SetOrderStatus("Inactive", orderID)

		if err != nil {
			return errors.New("error with inactive order")
		}
	}

	return nil
}
