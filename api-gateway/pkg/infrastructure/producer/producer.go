package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type TaskEventProducer struct {
	channel    *amqp.Channel
	exchange   string
	routingKey string
}

func NewTaskEventProducer(conn *amqp.Connection, exchange, routingKey string) (*TaskEventProducer, error) {
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

	return &TaskEventProducer{
		channel:    ch,
		exchange:   exchange,
		routingKey: routingKey,
	}, nil
}

func (tep *TaskEventProducer) PublishTaskCreated(ctx context.Context, data map[string]interface{}) error {
	event := map[string]interface{}{
		"type": "TaskCreated",
		"data": data,
		"time": time.Now(),
	}
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = tep.channel.PublishWithContext(
		ctx,
		tep.exchange,
		tep.routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}
	return nil
}

func (tep *TaskEventProducer) Close() error {
	return tep.channel.Close()
}
