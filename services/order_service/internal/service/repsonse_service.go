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
	ResponseToOrder(model *models.NewResponseModel, UserData models.UserData) (string, error)
	GetResponseById(ResponseID string, UserData models.UserData) (*models.ResponseDB, error)
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

func (s *ResponseService) GetResponseById(ResponseID string, UserData models.UserData) (*models.ResponseDB, error) {
	//UserInfo, err := s.GRPCClient.GetUserByTelegramID(context.Background(), UserData.TelegramID)

	//if err != nil {
	//	return nil, err
	//}

	return s.orderRepository.GetResponseById(ResponseID, UserData.UserID)
}

func (s *ResponseService) ResponseToOrder(Response *models.NewResponseModel, UserData models.UserData) (string, error) {
	//if UserData.Role ==

	TutorInfo, err := s.GRPCClient.GetUser(context.Background(), UserData.UserID)
	log.Println("ok")
	studentID, err := s.orderRepository.GetUserByOrder(Response.OrderId)
	log.Println(studentID)
	log.Println(err)
	StudentInfo, err := s.GRPCClient.GetStudent(context.Background(), studentID)
	log.Println(err)

	if err != nil {
		return "", custom_errors.ErrorGetUser
	}

	log.Println(StudentInfo)

	//StudentID, err := s.orderRepository.GetUserByOrder(Response.OrderId)
	//// тут ошибка
	//if err != nil || StudentID == nil {
	//	return "", custom_errors.ErrStudentByOrderNotFound
	//}

	responseID, err := s.orderRepository.CreateResponse(Response, TutorInfo)

	if err != nil {
		if errors.Is(custom_errors.ErrResponseAlredyExist, err) {
			return responseID, nil
		}
		return "", err
	}

	var ResponseToBroker models.ResponseToBrokerModel

	ResponseToBroker = models.ResponseToBrokerModel{
		UserId:  StudentInfo.TelegramID,
		OrderId: Response.OrderId,
		ChatId:  StudentInfo.TelegramID, // тут типо chatID
	}

	err = s.ProducerBroker.Publish("order_response", ResponseToBroker)

	return responseID, nil
}
