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

type ResponseServiceInt interface {
	ResponseToOrder(orderID string, newResponse *models.NewResponseModel, UserData models.UserData) (string, error)
	GetResponseById(ResponseID string, UserData models.UserData) (*models.ResponseDB, error)
	GetTutorsResponses(UserData models.UserData) ([]models.Response, error)
}

type ResponseService struct {
	orderRepository repository.OrderRepository
	ProducerBroker  rabbitmq.RabbitMQInterface
	GRPCClient      gRPC_client.GRPCClientInt
}

func NewResponseService(orderRepo repository.OrderRepository, producerBroker rabbitmq.RabbitMQInterface, grpcClient gRPC_client.GRPCClientInt) ResponseServiceInt {
	return &ResponseService{
		orderRepository: orderRepo,
		ProducerBroker:  producerBroker,
		GRPCClient:      grpcClient,
	}
}

func (s *ResponseService) GetTutorsResponses(UserData models.UserData) ([]models.Response, error) {
	return s.orderRepository.GetTutorsResponses(UserData.UserID)
}

func (s *ResponseService) GetResponseById(ResponseID string, UserData models.UserData) (*models.ResponseDB, error) {
	return s.orderRepository.GetResponseById(ResponseID, UserData.UserID)
}

func (s *ResponseService) ResponseToOrder(orderID string, newResponse *models.NewResponseModel, UserData models.UserData) (string, error) {
	if UserData.Role != "Tutor" {
		return "", custom_errors.ErrNotAllowed
	}

	TutorInfo, err := s.GRPCClient.GetTutor(context.Background(), UserData.UserID)

	if err != nil {
		return "", err
	}

	log.Println(UserData.UserID)

	isAvaliable, err := s.GRPCClient.CreateNewResponse(context.Background(), UserData.UserID)

	if err != nil {
		log.Println(err)
		return "", err
	}

	if !isAvaliable {
		log.Println("no resposes")
		return "", errors.New("No responses")
	}

	Order, err := s.orderRepository.GetOrderByID(orderID)

	if err != nil {
		return "", err
	}

	StudentInfo, err := s.GRPCClient.GetStudent(context.Background(), Order.StudentID)

	if err != nil {
		return "", custom_errors.ErrorGetUser
	}

	responseID, err := s.orderRepository.CreateResponse(orderID, newResponse, TutorInfo, UserData.Username)

	if err != nil {
		if errors.Is(custom_errors.ErrResponseAlredyExist, err) {
			return responseID, nil
		}
		return "", err
	}

	var ResponseToBroker models.ResponseToBrokerModel

	ResponseToBroker = models.ResponseToBrokerModel{
		ResponseID: responseID,
		TutorID:    TutorInfo.TelegramID,
		StudentID:  StudentInfo.TelegramID,
		OrderID:    orderID,
		Title:      Order.Title,
	}

	err = s.ProducerBroker.Publish("order_response", ResponseToBroker)

	return responseID, nil
}
