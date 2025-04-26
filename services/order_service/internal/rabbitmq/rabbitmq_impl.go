package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/randnull/Lessons/internal/config"
	"github.com/randnull/Lessons/internal/logger"
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

	log.Printf("[RabbitMQ] connecting to %v...", connectionLink)

	conn, err := amqp.Dial(connectionLink)

	if err != nil {
		if conn != nil {
			conn.Close()
		}
		log.Fatal("[RabbitMQ] failed to connect to RabbitMQ:" + err.Error())
	} else {
		log.Println("[RabbitMQ] Connected to RabbitMQ")
	}

	channel, err := conn.Channel()

	if err != nil {
		if channel != nil {
			channel.Close()
		}
		log.Fatal("[RabbitMQ] failed to open a channel:" + err.Error())
	}

	log.Println("[RabbitMQ] Channel opened")

	return &RabbitMQ{
		conn:    conn,
		channel: channel,
	}
}

func (r *RabbitMQ) Publish(queueName string, messageData interface{}) error {
	message, err := json.Marshal(messageData)
	if err != nil {
		logger.Error("[RabbitMQ] error marshal message" + err.Error())
	}

	_, err = r.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		logger.Error("[RabbitMQ] error init queue" + err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = r.channel.PublishWithContext(ctx,
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)

	if err != nil {
		logger.Error("[RabbitMQ] failed to publish a message" + err.Error())
		return err
	}

	logger.Info("[RabbitMQ] message pushed")

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
