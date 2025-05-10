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
	"github.com/randnull/Lessons/internal/utils"
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
	response, err := s.orderRepository.GetResponseById(ResponseID)

	if err != nil {
		if !errors.Is(err, custom_errors.ErrorNotFound) {
			logger.Error("[ResponseService] GetResponseById Error GetResponseById: " + err.Error())
			return nil, custom_errors.ErrorServiceError
		}
		return nil, custom_errors.ErrorNotFound
	}

	if UserData.UserID != response.TutorID {
		isUserRequest, err := s.orderRepository.CheckOrderByStudentID(response.OrderID, UserData.UserID)

		if err != nil {
			logger.Error("[ResponseService] GetResponseById Error CheckOrderByStudentID: " + err.Error())
			return nil, custom_errors.ErrorServiceError
		}

		if !isUserRequest {
			logger.Info("[ResponseService] GetResponseById Not Allowed CheckOrderByStudentID")
			return nil, custom_errors.ErrNotAllowed
		}
	}

	return response, nil
}

func (s *ResponseService) ResponseToOrder(orderID string, newResponse *models.NewResponseModel, UserData models.UserData) (string, error) {
	if UserData.Role != models.TutorType {
		return "", custom_errors.ErrNotAllowed
	}

	if utils.ContainsBadWords(newResponse.Greetings) {
		return "", custom_errors.ErrorBanWords
	}

	if s.orderRepository.CheckResponseExist(UserData.UserID, orderID) {
		return "", custom_errors.ErrResponseAlredyExist
	}

	TutorInfoRaw, err := s.GRPCClient.GetTutorInfoById(context.Background(), UserData.UserID, true)

	if err != nil {
		return "", custom_errors.ErrorGetUser
	}

	if TutorInfoRaw.ResponseCount <= 0 {
		return "", custom_errors.ErrNotAllowed
	}

	TutorInfo := &models.Tutor{
		Id:         TutorInfoRaw.Tutor.Id,
		TelegramID: TutorInfoRaw.Tutor.TelegramID,
		Name:       TutorInfoRaw.Tutor.Name,
		Bio:        TutorInfoRaw.Bio,
		Tags:       TutorInfoRaw.Tags,
	}

	isAvailable, err := s.GRPCClient.CreateNewResponse(context.Background(), UserData.UserID)

	if err != nil {
		logger.Error("[ResponseService] ResponseToOrder error CreateNewResponse: " + err.Error())
		return "", err
	}

	if !isAvailable {
		return "", custom_errors.ErrorServiceError
	}

	Order, err := s.orderRepository.GetOrderByID(orderID)

	if err != nil {
		logger.Error("[ResponseService] ResponseToOrder error GetOrderByID: " + err.Error())
		return "", err
	}

	if Order.Status != models.StatusNew {
		logger.Info("[ResponseService] ResponseToOrder GetOrderByID InActive Order: ")
		return "", custom_errors.ErrorBadStatus
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
		logger.Error("[ResponseService] ResponseToOrder error CreateResponse: " + err.Error())
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

	err = s.ProducerBroker.Publish(models.QueueNewResponse, ResponseToBroker)

	if err != nil {
		logger.Error("[OrderService] ResponseToOrder Error publishing order: " + err.Error())
		return "", err
	}

	return responseID, nil
}
