package service

import (
	"context"
	"errors"
	"github.com/randnull/Lessons/internal/custom_errors"
	"github.com/randnull/Lessons/internal/gRPC_client"
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/internal/rabbitmq"
	"github.com/randnull/Lessons/internal/repository"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

type ResponseServiceInt interface {
	ResponseToOrder(model *models.NewResponseModel, InitData initdata.InitData) (string, error)
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

func (s *ResponseService) ResponseToOrder(Response *models.NewResponseModel, InitData initdata.InitData) (string, error) {
	TutorInfo, err := s.GRPCClient.GetUser(context.Background(), InitData.User.ID)

	if err != nil {
		return "", custom_errors.ErrorGetUser
	}

	StudentID, err := s.orderRepository.GetUserByOrder(Response.OrderId)

	if err != nil || StudentID == nil {
		return "", custom_errors.ErrStudentByOrderNotFound
	}

	responseID, err := s.orderRepository.CreateResponse(Response, TutorInfo)

	if err != nil {
		if errors.Is(custom_errors.ErrResponseAlredyExist, err) {
			return responseID, nil
		}
		return "", err
	}

	var ResponseToBroker models.ResponseToBrokerModel

	ResponseToBroker = models.ResponseToBrokerModel{
		UserId:  *StudentID,
		OrderId: Response.OrderId,
		ChatId:  InitData.Chat.ID,
	}

	err = s.ProducerBroker.Publish("order_response", ResponseToBroker)

	return responseID, nil
}
