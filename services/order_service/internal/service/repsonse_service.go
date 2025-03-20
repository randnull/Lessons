package service

import (
	"errors"
	"github.com/randnull/Lessons/internal/custom_errors"
	"github.com/randnull/Lessons/internal/gRPC_client"
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/internal/rabbitmq"
	"github.com/randnull/Lessons/internal/redis"
	"github.com/randnull/Lessons/internal/repository"
	"log"
	"time"
)

type ResponseServiceInt interface {
	ResponseToOrder(model *models.NewResponseModel, UserData models.UserData) (string, error)
	GetResponseById(ResponseID string, UserData models.UserData) (*models.ResponseDB, error)
}

type ResponseService struct {
	orderRepository repository.OrderRepository
	ProducerBroker  rabbitmq.RabbitMQInterface
	GRPCClient      gRPC_client.GRPCClientInt
	RedisClient     redis.RedisInterface
}

func NewResponseService(orderRepo repository.OrderRepository, producerBroker rabbitmq.RabbitMQInterface, grpcClient gRPC_client.GRPCClientInt, redisClient redis.RedisInterface) ResponseServiceInt {
	return &ResponseService{
		orderRepository: orderRepo,
		ProducerBroker:  producerBroker,
		GRPCClient:      grpcClient,
		RedisClient:     redisClient,
	}
}

func (s *ResponseService) GetResponseById(ResponseID string, UserData models.UserData) (*models.ResponseDB, error) {
	//ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	//defer cancel()
	//
	//log.Println("Started")
	//response, err := s.RedisClient.GetResponseById(ctx, ResponseID)
	//log.Println("Stopped")
	//
	//log.Println(response, err)
	//
	//if err != nil || response == nil {
	//	return s.orderRepository.GetResponseById(ResponseID, UserData.UserID)
	//}
	//
	//fmt.Println("CACHED GetResponseById")
	//return response, nil

	return s.orderRepository.GetResponseById(ResponseID, UserData.UserID)
}

func (s *ResponseService) ResponseToOrder(Response *models.NewResponseModel, UserData models.UserData) (string, error) {
	//log.Println("check if tutor")
	if UserData.Role != "Tutor" {
		return "", custom_errors.ErrNotAllowed
	}
	//log.Println("tutor!")
	//log.Println("get info tutor")

	TutorInfo, err := s.GRPCClient.GetUser(UserData.UserID)
	//log.Println("info", TutorInfo)

	if err != nil {
		log.Println(err)
		return "", err
	}

	//log.Println("get info student")

	studentID, err := s.orderRepository.GetUserByOrder(Response.OrderId)
	//log.Println("info", studentID)

	if err != nil {
		return "", err
	}
	//log.Println("get info student 2")

	StudentInfo, err := s.GRPCClient.GetStudent(studentID)
	//log.Println("info", StudentInfo)

	if err != nil {
		log.Println(err)
		return "", custom_errors.ErrorGetUser
	}

	log.Println(StudentInfo)

	//fmt.Println(Response, TutorInfo, UserData.Username)
	responseID, err := s.orderRepository.CreateResponse(Response, TutorInfo, UserData.Username)

	if err != nil {
		if errors.Is(custom_errors.ErrResponseAlredyExist, err) {
			return responseID, nil
		}
		return "", err
	}

	log.Println("RESPONSE ID", Response.OrderId)

	err = s.RedisClient.SaveNewResponse(&models.ResponseDB{
		ID:            responseID,
		OrderID:       Response.OrderId,
		TutorID:       TutorInfo.Id,
		TutorUsername: UserData.Username,
		Name:          TutorInfo.Name,
		CreatedAt:     time.Now(), // Тут что-то придумать, так как время не совпадает с добавленным в базу
	})

	if err != nil {
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
