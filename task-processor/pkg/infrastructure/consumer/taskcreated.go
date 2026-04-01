package consumer

import (
	"context"
	"log"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"

	model2 "taskprocessor/pkg/app/model"
	"taskprocessor/pkg/infrastructure/publisher"
	"taskprocessor/pkg/infrastructure/validator"
)

type TaskCreatedConsumer struct {
	channel   *amqp.Channel
	queueName string
	publisher *publisher.NotificationPublisher
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

func NewTaskCreatedConsumer(conn *amqp.Connection, queueName, bindingKey string, pub *publisher.NotificationPublisher) (*TaskCreatedConsumer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(
		q.Name,
		bindingKey,
		"tasks",
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &TaskCreatedConsumer{
		channel:   ch,
		queueName: q.Name,
		publisher: pub,
		ctx:       ctx,
		cancel:    cancel,
	}, nil
}

func (c *TaskCreatedConsumer) Start() error {
	deliveries, err := c.channel.Consume(
		c.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			select {
			case <-c.ctx.Done():
				log.Println("context cancelled, stopping goroutine")
				return
			case d, ok := <-deliveries:
				if !ok {
					log.Println("deliveries channel closed")
					return
				}
				log.Printf("received message, delivery tag %d, redelivered=%v", d.DeliveryTag, d.Redelivered)
				c.handleMessage(d)
			}
		}
	}()
	return nil
}

func (c *TaskCreatedConsumer) Stop() error {
	c.cancel()
	return nil
}

func (c *TaskCreatedConsumer) Wait() {
	c.wg.Wait()
}

func (c *TaskCreatedConsumer) Close() error {
	if c.channel != nil {
		return c.channel.Close()
	}
	return nil
}

func (c *TaskCreatedConsumer) handleMessage(d amqp.Delivery) {
	var event model2.TaskCreatedEvent
	if err := event.FromJSON(d.Body); err != nil {
		log.Printf("failed to parse event: %v", err)
		d.Nack(false, false)
		return
	}
	log.Printf("event: TaskID=%s, Email=%s, Phone=%s, Title=%s",
		event.Data.TaskID, event.Data.Email, event.Data.Phone, event.Data.Title)

	if !validator.ValidateEvent(&event) {
		log.Printf("invalid event: missing task ID or contact info")
		d.Nack(false, false)
		return
	}

	notification := model2.NewNotificationFromEvent(&event)

	ctx := context.Background()
	if err := c.publisher.Publish(ctx, notification); err != nil {
		log.Printf("failed to publish notification: %v", err)
		d.Nack(false, true)
		return
	}
	log.Println("notification published successfully")

	d.Ack(false)
}
