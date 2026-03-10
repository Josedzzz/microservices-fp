// Package messaging provides functionality for producing and consuming messages from a message broker
package messaging

import (
	"context"
	"encoding/json"

	"github.com/rabbitmq/amqp091-go"
)

// Producer is responsible for publishing messages to RabbitMQ
type Producer struct {
	channel *amqp091.Channel
}

// NewProducer creates a new Producer instance with the provided RabbitMQ channel
func NewProducer(ch *amqp091.Channel) *Producer {
	return &Producer{channel: ch}
}

// RecoverPasswordPayload represents the data for the password recovery event
type RecoverPasswordPayload struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

// PublishRecoverPasswordEvent publishes a recover password event to the "auth.events" exchange
func (p *Producer) PublishRecoverPasswordEvent(ctx context.Context, email, token string) error {
	payload := RecoverPasswordPayload{
		Email: email,
		Token: token,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// We'll declare the exchange here to ensure it exists
	err = p.channel.ExchangeDeclare(
		"auth.events", // name
		"topic",       // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		return err
	}

	return p.channel.PublishWithContext(
		ctx,
		"auth.events",      // exchange
		"password.recover", // routing key
		false,              // mandatory
		false,              // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
