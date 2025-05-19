package service

import (
	"context"
	"errors"
	"github.com/randnull/Lessons/internal/gRPC_client"
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/internal/repository"
	"github.com/randnull/Lessons/internal/utils"
	"github.com/randnull/Lessons/pkg/custom_errors"
	"github.com/randnull/Lessons/pkg/logger"
	"github.com/randnull/Lessons/pkg/rabbitmq"
)

type OrderServiceInt interface {
	CreateOrder(order *models.NewOrder, UserData models.UserData) (string, error)
	UpdateOrder(orderID string, order *models.UpdateOrder, UserData models.UserData) error
	DeleteOrder(orderID string, UserData models.UserData) error

	GetOrderById(id string, UserData models.UserData) (*models.OrderDetails, error)
	GetStudentOrdersWithPagination(page int, size int, UserData models.UserData) (*models.OrderPagination, error)
	GetOrdersWithPagination(page int, size int, tag string, UserData models.UserData) (*models.OrderPagination, error)
	GetOrderByIdTutor(id string, UserData models.UserData) (*models.OrderDetailsTutor, error)
	GetAllOrders(UserData models.UserData) ([]*models.Order, error)
	GetAllUsersOrders(UserData models.UserData) ([]*models.Order, error)

	SelectTutor(responseID string, UserData models.UserData) error
	SuggestOrderToTutor(orderID, tutorID string, UserData models.UserData) error

	SetActiveOrderStatus(orderID string, IsActive bool, UserData models.UserData) error
	SetBanOrderStatus(orderID string, UserData models.UserData) error
	SetApprovedOrderStatus(orderID string, UserData models.UserData) error
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
	if order.MinPrice > order.MaxPrice {
		return "", custom_errors.ErrorParams
	}

	if utils.ContainsBadWords(order.Title) || utils.ContainsBadWords(order.Description) {
		return "", custom_errors.ErrorBanWords
	}

	_, err := orderServ.GRPCClient.GetStudent(context.Background(), UserData.UserID)

	if err != nil {
		logger.Error("[OrderService] CreateOrder error GetStudent: " + err.Error())
		return "", custom_errors.ErrorServiceError
	}

	OrderToCreate := &models.CreateOrder{
		Order:     order,
		StudentID: UserData.UserID,
	}

	OrderID, err := orderServ.orderRepository.CreateOrder(OrderToCreate)

	if err != nil {
		logger.Error("[OrderService] CreateOrder error CreateOrder: " + err.Error())
		return "", err
	}

	OrderToBroker := models.OrderToBroker{
		ID:        OrderID,
		StudentID: UserData.TelegramID,
		Title:     order.Title,
		Tags:      order.Tags,
		Status:    models.StatusWaiting,
	}

	err = orderServ.ProducerBroker.Publish(models.QueueNewOrder, OrderToBroker)

	if err != nil {
		logger.Error("[OrderService] CreateOrder Error publishing order: " + err.Error())
	}

	return OrderID, nil
}

func (orderServ *OrderService) GetOrderById(id string, UserData models.UserData) (*models.OrderDetails, error) {
	order, err := orderServ.orderRepository.GetOrderByID(id)

	if err != nil {
		return nil, err
	}

	if order.StudentID != UserData.UserID {
		logger.Info("[OrderService] GetOrderById not allowed. User: " + UserData.UserID + " User-Order: " + order.StudentID)
		return nil, custom_errors.ErrNotAllowed
	}

	responses, err := orderServ.orderRepository.GetResponsesByOrderID(id)

	if err != nil {
		logger.Error("[OrderService] GetOrderById Error : " + err.Error())
		return nil, err
	}

	orderDetails := &models.OrderDetails{
		Order:     *order,
		Responses: responses,
	}

	return orderDetails, nil
}

func (orderServ *OrderService) GetOrderByIdTutor(id string, UserData models.UserData) (*models.OrderDetailsTutor, error) {
	order, err := orderServ.orderRepository.GetOrderByID(id)

	if err != nil {
		logger.Error("[OrderService] GetOrderByIdTutor Error GetOrderByID: " + err.Error())
		return nil, err
	}

	isResponded, err := orderServ.orderRepository.GetTutorIsRespond(id, UserData.UserID)

	if err != nil {
		logger.Error("[OrderService] GetOrderByIdTutor Error GetTutorIsRespond: " + err.Error())
	}

	orderDetails := &models.OrderDetailsTutor{
		Order:       *order,
		IsResponded: isResponded,
	}

	return orderDetails, nil
}

func (orderServ *OrderService) GetOrdersWithPagination(page int, size int, tag string, UserData models.UserData) (*models.OrderPagination, error) {
	limit := size
	offset := (page - 1) * size

	orders, count, err := orderServ.orderRepository.GetOrdersPagination(limit, offset, tag)

	if err != nil {
		logger.Error("[OrderService] GetOrdersWithPagination Error GetOrdersPagination: " + err.Error())
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
		logger.Error("[OrderService] GetStudentOrdersWithPagination Error GetStudentOrdersPagination: " + err.Error())
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
	if utils.ContainsBadWords(order.Title) || utils.ContainsBadWords(order.Description) {
		return custom_errors.ErrorBanWords
	}

	isExist, err := orderServ.orderRepository.CheckOrderByStudentID(orderID, UserData.UserID)

	if err != nil {
		logger.Error("[OrderService] UpdateOrder Error CheckOrderByStudentID: " + err.Error())
		return custom_errors.ErrorServiceError
	}

	if !isExist {
		logger.Error("[OrderService] UpdateOrder Not allowed. User: " + UserData.UserID + " Order: " + orderID)
		return custom_errors.ErrNotAllowed
	}

	return orderServ.orderRepository.UpdateOrder(orderID, order)
}

func (orderServ *OrderService) DeleteOrder(orderID string, UserData models.UserData) error {
	isExist, err := orderServ.orderRepository.CheckOrderByStudentID(orderID, UserData.UserID)

	if err != nil {
		logger.Error("[OrderService] UpdateOrder Error CheckOrderByStudentID: " + err.Error())
		return custom_errors.ErrorServiceError
	}

	if !isExist {
		logger.Error("[OrderService] UpdateOrder Not allowed. User: " + UserData.UserID + " Order: " + orderID)
		return custom_errors.ErrNotAllowed
	}

	return orderServ.orderRepository.DeleteOrder(orderID)
}

func (orderServ *OrderService) GetAllUsersOrders(UserData models.UserData) ([]*models.Order, error) {
	return orderServ.orderRepository.GetStudentOrders(UserData.UserID)
}

func (orderServ *OrderService) SelectTutor(responseID string, UserData models.UserData) error {
	response, err := orderServ.orderRepository.GetResponseById(responseID)

	if err != nil {
		logger.Error("[OrderService] SelectTutor Error GetResponseById: " + err.Error())
		return err
	}

	order, err := orderServ.orderRepository.GetOrderByID(response.OrderID)

	if err != nil {
		logger.Error("[OrderService] SelectTutor Error GetOrderByID: " + err.Error())
		return err
	}

	if order.StudentID != UserData.UserID {
		logger.Info("[OrderService] SelectTutor Not allowed. User: " + UserData.UserID + " Order: " + response.OrderID)
		return custom_errors.ErrNotAllowed
	}

	if order.Status != models.StatusNew {
		logger.Info("[OrderService] SelectTutor Order Not New. User: " + UserData.UserID + " Order: " + response.OrderID)
		return custom_errors.ErrorBadStatus
	}

	tutor, err := orderServ.GRPCClient.GetTutor(context.Background(), response.TutorID)

	if err != nil {
		logger.Error("[OrderService] SelectTutor GetTutor Error + " + err.Error() + "  Order: " + response.OrderID)
		return custom_errors.ErrorNotFound
	}

	student, err := orderServ.GRPCClient.GetStudent(context.Background(), order.StudentID)

	if err != nil {
		logger.Error("[OrderService] SelectTutor GetStudent Error + " + err.Error() + "  Order: " + response.OrderID)
	}

	err = orderServ.orderRepository.SetTutorToOrder(response, UserData)

	if err != nil {
		logger.Error("[OrderService] SelectTutor Error SetTutorToOrder: " + err.Error())
		return custom_errors.ErrorSelectTutor
	}

	err = orderServ.ProducerBroker.Publish("selected_orders", models.SelectedResponseToBroker{
		OrderName:  order.Title,
		OrderID:    order.ID,
		ResponseID: responseID,
		StudentID:  student.TelegramID,
		TutorID:    tutor.TelegramID,
	})

	if err != nil {
		logger.Error("[OrderService] SelectTutor error with push response selected to broker. Error: " + err.Error())
		return nil
	}

	return nil
}

func (orderServ *OrderService) SetActiveOrderStatus(orderID string, IsActive bool, UserData models.UserData) error {
	order, err := orderServ.orderRepository.GetOrderByID(orderID)

	if err != nil {
		logger.Error("[OrderService] SetActiveOrderStatus Error GetOrderByID: " + err.Error())
		return custom_errors.ErrGetOrder
	}

	if order.StudentID != UserData.UserID {
		logger.Info("[OrderService] SetActiveOrderStatus Not allowed. User: " + UserData.UserID + " User-Order: " + order.StudentID)
		return custom_errors.ErrNotAllowed
	}

	if IsActive {
		if order.Status != models.StatusInactive {
			logger.Info("[OrderService] SetActiveOrderStatus Not Inactive state. OrderID:" + orderID)
			return custom_errors.ErrorBadStatus
		}

		err = orderServ.orderRepository.SetOrderStatus(models.StatusNew, orderID)

		if err != nil {
			logger.Error("[OrderService] SetActiveOrderStatus Error SetOrderStatus: " + err.Error())
			return err
		}
	} else {
		if order.Status != models.StatusNew {
			logger.Info("[OrderService] SetActiveOrderStatus Not Active state. OrderID:" + orderID)
			return custom_errors.ErrorBadStatus
		}

		err = orderServ.orderRepository.SetOrderStatus(models.StatusInactive, orderID)

		if err != nil {
			logger.Error("[OrderService] SetActiveOrderStatus Error SetOrderStatus: " + err.Error())
			return err
		}
	}

	return nil
}

func (orderServ *OrderService) SetBanOrderStatus(orderID string, UserData models.UserData) error {
	_, err := orderServ.orderRepository.GetOrderByID(orderID)

	if err != nil {
		logger.Error("[OrderService] SetBanOrderStatus Error GetOrderByID: " + err.Error())
		return err
	}

	err = orderServ.orderRepository.SetOrderStatus(models.StatusBan, orderID)

	if err != nil {
		logger.Error("[OrderService] SetBanOrderStatus Error SetOrderStatus: " + err.Error())
		return err
	}

	return nil
}

func (orderServ *OrderService) SuggestOrderToTutor(orderID, tutorID string, UserData models.UserData) error {
	user, err := orderServ.GRPCClient.GetTutor(context.Background(), tutorID)

	if err != nil {
		logger.Error("[OrderService] SuggestOrderToTutor Error GRPCClient.GetTutor: " + err.Error())
		return err
	}

	orderInfo, err := orderServ.orderRepository.GetOrderByID(orderID)

	if err != nil {
		logger.Error("[OrderService] SuggestOrderToTutor Error GetOrderByID: " + err.Error())
		return err
	}

	if orderInfo.Status != models.StatusNew {
		logger.Info("[OrderService] SuggestOrderToTutor Not New state. OrderID:" + orderID)
		return errors.New("order not new state")
	}

	if orderInfo.StudentID != UserData.UserID {
		logger.Info("[OrderService] SuggestOrderToTutor Not allowed. User: " + UserData.UserID + " User-Order: " + orderInfo.StudentID)

		return custom_errors.ErrNotAllowed
	}

	suggestOrderModel := models.SuggestOrder{
		ID:              orderInfo.ID,
		TutorTelegramID: user.TelegramID,
		Title:           orderInfo.Title,
		Description:     orderInfo.Description,
		MinPrice:        orderInfo.MinPrice,
		MaxPrice:        orderInfo.MaxPrice,
	}

	err = orderServ.ProducerBroker.Publish(models.QueueSuggestOrder, suggestOrderModel)
	if err != nil {
		logger.Error("[OrderService] CreateOrder Error publishing order: " + err.Error())
		return err
	}

	return nil
}

func (orderServ *OrderService) SetApprovedOrderStatus(orderID string, UserData models.UserData) error {
	order, err := orderServ.orderRepository.GetOrderByID(orderID)

	if err != nil {
		logger.Error("[OrderService] SetApprovedOrderStatus Error GetOrderByID: " + err.Error())
		return err
	}

	if order.Status != models.StatusWaiting {
		return custom_errors.ErrorBadStatus
	}

	err = orderServ.orderRepository.SetOrderStatus(models.StatusNew, orderID)

	if err != nil {
		logger.Error("[OrderService] SetApprovedOrderStatus Error SetOrderStatus: " + err.Error())
		return err
	}

	student, err := orderServ.GRPCClient.GetStudent(context.Background(), order.StudentID)

	if err != nil {
		logger.Error("[OrderService] SetApprovedOrderStatus Error GetStudent: " + err.Error())
		return err
	}

	OrderToBroker := models.OrderToBroker{
		ID:        orderID,
		StudentID: student.TelegramID,
		Title:     order.Title,
		Tags:      order.Tags,
		Status:    models.StatusNew,
	}

	err = orderServ.ProducerBroker.Publish(models.QueueNewOrder, OrderToBroker)

	if err != nil {
		logger.Error("[OrderService] SetApprovedOrderStatus Error publishing order: " + err.Error())
	}

	return nil
}
