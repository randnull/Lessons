package service

import (
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/internal/rabbitmq"
	"github.com/randnull/Lessons/internal/repository"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

type ResponseServiceInt interface {
	ResponseToOrder(model *models.NewResponseModel, InitData initdata.InitData) error
}

type ResponseService struct {
	orderRepository repository.OrderRepository
	ProducerBroker  rabbitmq.RabbitMQInterface
}

func NewResponseService(orderRepo repository.OrderRepository, producerBroker rabbitmq.RabbitMQInterface) ResponseServiceInt {
	return &ResponseService{
		orderRepository: orderRepo,
		ProducerBroker:  producerBroker,
	}
}

func (s *ResponseService) ResponseToOrder(Response *models.NewResponseModel, InitData initdata.InitData) error {
	err := s.orderRepository.CreateResponse(Response, InitData)

	err = s.ProducerBroker.Publish("order_response", Response)

	if err != nil {
		return err
	}

	return nil
}
