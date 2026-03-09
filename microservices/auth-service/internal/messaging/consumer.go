// Package messaging provides functionality for consuming messages from a message broker
package messaging

import (
	"context"
	"encoding/json"
	"log"

	"github.com/rabbitmq/amqp091-go"

	"auth-service/internal/service"
)

// Consumer is responsible for consuming messages from a RabbitMQ queue
type Consumer struct {
	service *service.AuthService
	channel *amqp091.Channel
}

// NewConsumer creates a new Consumer instance with the provided RabbitMQ channel
func NewConsumer(ch *amqp091.Channel, service *service.AuthService) *Consumer {
	return &Consumer{service: service, channel: ch}
}

// Start begins consuming messages from the "employee.created" queue
func (c *Consumer) Start() error {
	err := c.channel.ExchangeDeclare(
		"employees.events",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	queue, err := c.channel.QueueDeclare(
		"auth.employee.created",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = c.channel.QueueBind(
		queue.Name,
		"employee.created",
		"employees.events",
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// 4️⃣ consume
	msgs, err := c.channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {

			var payload struct {
				Email string `json:"email"`
			}

			err := json.Unmarshal(msg.Body, &payload)
			if err != nil {
				log.Println("invalid message")
				continue
			}

			err = c.service.CreateUser(context.Background(), payload.Email, "EMPLOYEE")
			if err != nil {
				log.Println(err)
			}
		}
	}()

	log.Println("Auth service listening for employee.created events")

	return nil
}
