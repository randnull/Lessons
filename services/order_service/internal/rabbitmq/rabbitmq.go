package rabbitmq

type RabbitMQInterface interface {
	Publish(queueName string, messageData interface{}) error
	Close()
}
