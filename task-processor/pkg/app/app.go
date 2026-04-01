package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"taskprocessor/pkg/infrastructure/consumer"
	"taskprocessor/pkg/infrastructure/publisher"
)

type Config struct {
	RabbitURL             string
	NotificationsExchange string
	QueueName             string
	RoutingKey            string
}

type App struct {
	conn     *amqp.Connection
	pub      *publisher.NotificationPublisher
	consumer *consumer.TaskCreatedConsumer
}

func NewApp(cfg Config) (*App, error) {
	conn, err := amqp.Dial(cfg.RabbitURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	pub, err := publisher.NewNotificationPublisher(conn, cfg.NotificationsExchange)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create notification publisher: %w", err)
	}

	taskConsumer, err := consumer.NewTaskCreatedConsumer(conn, cfg.QueueName, cfg.RoutingKey, pub)
	if err != nil {
		pub.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}
	log.Println("consumer created")

	return &App{
		conn:     conn,
		pub:      pub,
		consumer: taskConsumer,
	}, nil
}

func (a *App) Start() error {
	if err := a.consumer.Start(); err != nil {
		return fmt.Errorf("failed to start consumer: %w", err)
	}
	log.Println("waiting for messages")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := a.consumer.Stop(); err != nil {
		log.Printf("Consumer stop error: %v", err)
	}
	a.consumer.Wait()
	return nil
}

func (a *App) Close() error {
	var errs []error
	if err := a.consumer.Close(); err != nil {
		errs = append(errs, fmt.Errorf("consumer close: %w", err))
	}
	if err := a.pub.Close(); err != nil {
		errs = append(errs, fmt.Errorf("publisher close: %w", err))
	}
	if err := a.conn.Close(); err != nil {
		errs = append(errs, fmt.Errorf("connection close: %w", err))
	}
	if len(errs) > 0 {
		return fmt.Errorf("close errors: %v", errs)
	}
	return nil
}
