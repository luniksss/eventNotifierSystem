package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type NotificationPublisher struct {
	channel  *amqp.Channel
	exchange string
}

func NewNotificationPublisher(conn *amqp.Connection, exchange string) (*NotificationPublisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &NotificationPublisher{channel: ch, exchange: exchange}, nil
}

func (p *NotificationPublisher) Publish(ctx context.Context, notification interface{}) error {
	body, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	err = p.channel.PublishWithContext(ctx,
		p.exchange,
		"notification.send",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("publish failed: %v", err)
		return fmt.Errorf("failed to publish notification: %w", err)
	}
	log.Println("publish succeeded")
	return nil
}

func (p *NotificationPublisher) Close() error {
	return p.channel.Close()
}
