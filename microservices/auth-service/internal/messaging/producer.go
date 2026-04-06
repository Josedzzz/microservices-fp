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

// UserEventPayload represents the data for auth-related user events
type UserEventPayload struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

// PublishUserCreatedEvent publishes a user.created event when a new user is onboarded
func (p *Producer) PublishUserCreatedEvent(ctx context.Context, email, token string) error {
	return p.publishAuthEvent(ctx, "user.created", email, token)
}

// PublishUserRecoveryEvent publishes a user.recovery event when password recovery is requested
func (p *Producer) PublishUserRecoveryEvent(ctx context.Context, email, token string) error {
	return p.publishAuthEvent(ctx, "user.recovery", email, token)
}

// publishAuthEvent is a helper to publish events to the "auth.events" exchange
func (p *Producer) publishAuthEvent(ctx context.Context, routingKey, email, token string) error {
	payload := UserEventPayload{
		Email: email,
		Token: token,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Ensure the exchange exists
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
		"auth.events", // exchange
		routingKey,    // routing key
		false,         // mandatory
		false,         // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

// PublishRecoverPasswordEvent (DEPRECATED) kept for compatibility during transition if needed
func (p *Producer) PublishRecoverPasswordEvent(ctx context.Context, email, token string) error {
	return p.publishAuthEvent(ctx, "password.recover", email, token)
}
