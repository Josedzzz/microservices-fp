// Package messaging provides functionality for consuming messages from a message broker
package messaging

import (
	"context"
	"encoding/json"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

// AuthService defines the interface for user-related business logic
type AuthService interface {
	CreateUser(ctx context.Context, email, role, authStatus string) error
	DisableUser(ctx context.Context, email string) error
}

// Consumer is responsible for consuming messages from a RabbitMQ queue
type Consumer struct {
	service AuthService
	channel *amqp091.Channel
}

// NewConsumer creates a new Consumer instance with the provided RabbitMQ channel
func NewConsumer(ch *amqp091.Channel, service AuthService) *Consumer {
	return &Consumer{service: service, channel: ch}
}

// Start begins consuming messages from the "employee.created" and "employee.deleted" queues
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

	// 1. Setup consumer for employee.created
	createQueue, err := c.channel.QueueDeclare(
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
		createQueue.Name,
		"employee.created",
		"employees.events",
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// 2. Setup consumer for employee.deleted
	deleteQueue, err := c.channel.QueueDeclare(
		"auth.employee.deleted",
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
		deleteQueue.Name,
		"employee.deleted",
		"employees.events",
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Start consuming from both (shared handler logic or separate)
	// We'll use a single consumer on a merged queue or separate.
	// For simplicity, we can just declare a single queue bound to both keys if we want,
	// but separate queues are safer. Let's use separate for clarity.

	if err := c.consume(createQueue.Name, "created"); err != nil {
		return err
	}
	if err := c.consume(deleteQueue.Name, "deleted"); err != nil {
		return err
	}

	return nil
}

func (c *Consumer) consume(queueName, action string) error {
	msgs, err := c.channel.Consume(
		queueName,
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
				Email  string `json:"email"`
				Role   string `json:"role"`
				Status string `json:"status"`
			}

			err := json.Unmarshal(msg.Body, &payload)
			if err != nil {
				log.Println("invalid message payload:", err)
				continue
			}

			if action == "created" {
				authStatus := payload.Status
				if authStatus == "" {
					authStatus = "ACTIVE"
				}
				err = c.service.CreateUser(context.Background(), payload.Email, payload.Role, authStatus)
			} else {
				err = c.service.DisableUser(context.Background(), payload.Email)
			}

			if err != nil {
				log.Printf("error processing %s event: %v", action, err)
			}
		}
	}()

	log.Printf("Auth service listening for employee.%s events", action)
	return nil
}
