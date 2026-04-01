package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"

	"notifier/pkg/app/consumer"
)

func main() {
	rabbitURL := getEnv("RABBITMQ_URL", "amqp://guest:guest@notifier-rabbitmq:5672/")

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatal("failed to connect to rabbitmq:", err)
	}
	defer conn.Close()

	notifConsumer, err := consumer.NewNotificationConsumer(conn, "notifier-queue", "notification.send")
	if err != nil {
		log.Fatal("failed to create consumer:", err)
	}
	if err := notifConsumer.Start(); err != nil {
		log.Fatal("failed to start consumer:", err)
	}

	log.Println("notifier started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
