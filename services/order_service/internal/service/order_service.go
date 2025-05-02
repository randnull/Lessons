package service

import (
	"context"
	"errors"
	"github.com/randnull/Lessons/internal/custom_errors"
	"github.com/randnull/Lessons/internal/gRPC_client"
	"github.com/randnull/Lessons/internal/logger"
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/internal/rabbitmq"
	"github.com/randnull/Lessons/internal/repository"
)

type OrderServiceInt interface {
	// order edit
	CreateOrder(order *models.NewOrder, UserData models.UserData) (string, error)
	UpdateOrder(orderID string, order *models.UpdateOrder, UserData models.UserData) error
	DeleteOrder(orderID string, UserData models.UserData) error
	SelectTutor(responseID string, UserData models.UserData) error
	SuggestOrderToTutor(orderID, tutorID string, UserData models.UserData) error
	SetActiveOrderStatus(orderID string, IsActive bool, UserData models.UserData) error

	// order getting
	GetOrderById(id string, UserData models.UserData) (*models.OrderDetails, error)
	GetStudentOrdersWithPagination(page int, size int, UserData models.UserData) (*models.OrderPagination, error)
	GetOrdersWithPagination(page int, size int, tag string, UserData models.UserData) (*models.OrderPagination, error)
	GetOrderByIdTutor(id string, UserData models.UserData) (*models.OrderDetailsTutor, error)
	GetAllOrders(UserData models.UserData) ([]*models.Order, error)
	GetAllUsersOrders(UserData models.UserData) ([]*models.Order, error)
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
		return "", errors.New("max price less than min")
	}

	if order.Grade == "" {
		return "", errors.New("grade is null")
	}

	if len(order.Title) < 5 || len(order.Title) > 150 {
		return "", errors.New("error size")
	}

	if len(order.Description) < 5 || len(order.Description) > 500 {
		return "", errors.New("error size")
	}

	_, err := orderServ.GRPCClient.GetStudent(context.Background(), UserData.UserID)

	if err != nil {
		logger.Error("[OrderService] CreateOrder error get student: " + err.Error())
		return "", custom_errors.ErrorGetUser
	}

	OrderToCreate := &models.CreateOrder{
		Order:     order,
		StudentID: UserData.UserID,
	}

	OrderID, err := orderServ.orderRepository.CreateOrder(OrderToCreate)

	if err != nil {
		logger.Error("[OrderService] CreateOrder error create order: " + err.Error())
		return "", err
	}

	OrderToBroker := models.OrderToBroker{
		ID:        OrderID,
		StudentID: UserData.TelegramID,
		Title:     order.Title,
		Tags:      order.Tags,
		Status:    "New",
	}

	err = orderServ.ProducerBroker.Publish("new_orders", OrderToBroker)
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
	isExist, err := orderServ.orderRepository.CheckOrderByStudentID(orderID, UserData.UserID)

	if !isExist || err != nil {
		if err != nil {
			logger.Error("[OrderService] UpdateOrder Error CheckOrderByStudentID: " + err.Error())
			return custom_errors.ErrorServiceError
		}
		logger.Info("[OrderService] UpdateOrder Not allowed. User: " + UserData.UserID + " Order: " + orderID)
		return custom_errors.ErrNotAllowed
	}

	return orderServ.orderRepository.UpdateOrder(orderID, order)
}

func (orderServ *OrderService) DeleteOrder(orderID string, UserData models.UserData) error {
	isExist, err := orderServ.orderRepository.CheckOrderByStudentID(orderID, UserData.UserID)

	if !isExist || err != nil {
		if err != nil {
			logger.Error("[OrderService] UpdateOrder Error CheckOrderByStudentID: " + err.Error())
			return custom_errors.ErrorServiceError
		}
		logger.Info("[OrderService] UpdateOrder Not allowed. User: " + UserData.UserID + " Order: " + orderID)
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

	//err = orderServ.orderRepository.SetOrderStatus(models.StatusWaiting, order.ID)
	//
	//if err != nil {
	//	logger.Error("[UserService] SelectTutor error SetOrderStatus: " + err.Error())
	//}

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
			logger.Info("[OrderService] SetActiveOrderStatus Not Invative state. OrderID:" + orderID)
			return errors.New("error not Inactive state")
		}

		err = orderServ.orderRepository.SetOrderStatus(models.StatusNew, orderID)

		if err != nil {
			logger.Error("[OrderService] SetActiveOrderStatus Error SetOrderStatus: " + err.Error())
			return custom_errors.ErrorSetStatus
		}
	} else {
		if order.Status != models.StatusNew {
			logger.Info("[OrderService] SetActiveOrderStatus Not Active state. OrderID:" + orderID)
		}

		err = orderServ.orderRepository.SetOrderStatus(models.StatusInactive, orderID)

		if err != nil {
			logger.Error("[OrderService] SetActiveOrderStatus Error SetOrderStatus: " + err.Error())
			return custom_errors.ErrorSetStatus
		}
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
		return errors.New("order not NEW state")
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

	err = orderServ.ProducerBroker.Publish("suggest_order", suggestOrderModel)
	if err != nil {
		logger.Error("[OrderService] CreateOrder Error publishing order: " + err.Error())
		return err
	}

	return nil
}
