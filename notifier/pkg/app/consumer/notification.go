package consumer

import (
	"encoding/json"
	"log"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"

	"notifier/pkg/app/model"
	"notifier/pkg/app/sender"
)

type NotificationConsumer struct {
	channel     *amqp.Channel
	queueName   string
	emailSender *sender.EmailSender
	smsSender   *sender.SMSSender
}

func NewNotificationConsumer(conn *amqp.Connection, queueName, bindingKey string) (*NotificationConsumer, error) {
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
		"notifications",
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &NotificationConsumer{
		channel:     ch,
		queueName:   q.Name,
		emailSender: sender.NewEmailSender(),
		smsSender:   sender.NewSMSSender(),
	}, nil
}

func (c *NotificationConsumer) Start() error {
	msgs, err := c.channel.Consume(
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

	go func() {
		for d := range msgs {
			c.handleMessage(d)
		}
	}()
	return nil
}

func (c *NotificationConsumer) handleMessage(d amqp.Delivery) {
	var notif model.Notification
	if err := json.Unmarshal(d.Body, &notif); err != nil {
		log.Printf("failed to parse notification: %v", err)
		d.Nack(false, false)
		return
	}

	errCh := make(chan error, 2)
	var wg sync.WaitGroup

	if notif.Email != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := c.emailSender.Send(notif.Email, notif.Text); err != nil {
				errCh <- err
			}
		}()
	}

	if notif.Phone != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := c.smsSender.Send(notif.Phone, notif.Text); err != nil {
				errCh <- err
			}
		}()
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	var errors []error
	for err := range errCh {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		log.Printf("failed to send notifications for task %s: %v", notif.TaskID, errors)
		d.Nack(false, true)
		return
	}

	d.Ack(false)
	log.Printf("notifications sent for task %s", notif.TaskID)
}
