package rabbitmq

import "github.com/randnull/Lessons/internal/models"

type RabbitMQInterface interface {
	Publish(queueName string, OrderInfo models.OrderToBrokerWithID) error
	Close()
}
