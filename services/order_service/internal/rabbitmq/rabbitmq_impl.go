package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/randnull/Lessons/internal/config"
	"github.com/randnull/Lessons/internal/models"
	"log"
	"time"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQ(cfg config.MQConfig) *RabbitMQ {
	MqUser := cfg.User
	MqPassword := cfg.Pass
	MqHost := cfg.Host
	MqPort := cfg.Port

	connectionLink := fmt.Sprintf("amqp://%v:%v@%v:%v/", MqUser, MqPassword, MqHost, MqPort)

	log.Printf("connecting to %v...", connectionLink)

	conn, err := amqp.Dial(connectionLink)

	if err != nil {
		if conn != nil {
			conn.Close()
		}
		//log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}

	//log.Println("Connected to RabbitMQ")

	channel, err := conn.Channel()

	if err != nil {
		if channel != nil {
			channel.Close()
		}
		//log.Fatalf("failed to open a channel: %v", err)
	}

	//log.Println("Channel opened")

	return &RabbitMQ{
		conn:    conn,
		channel: channel,
	}
}

func (r *RabbitMQ) Publish(queueName string, OrderInfo models.OrderToBrokerWithID) error {
	_, err := r.channel.QueueDeclare(
		queueName, // Имя очереди
		true,      // Долговечность (durable)
		false,     // Автоудаление (auto-delete)
		false,     // Эксклюзивная (exclusive)
		false,     // Без ожидания (no-wait)
		nil,       // Доп. аргументы
	)

	message, err := json.Marshal(OrderInfo)

	if err != nil {
		log.Fatalf("failed to marshal order info: %v", err)
	}

	if err != nil {
		log.Printf("failed to declare a queue: %v", err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = r.channel.PublishWithContext(ctx,
		"",
		queueName,
		false, // Mandatory
		false, // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)

	if err != nil {
		log.Printf("failed to publish a message: %v", err)
		return err
	}

	log.Printf("Successfully published a message %v", message)

	return nil
}

func (r *RabbitMQ) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}
