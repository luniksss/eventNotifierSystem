package main

import (
	"log"
	"os"

	"taskprocessor/pkg/app"
)

func main() {
	cfg := app.Config{
		RabbitURL:             getEnv("RABBITMQ_URL", "amqp://guest:guest@notifier-rabbitmq:5672/"),
		NotificationsExchange: getEnv("NOTIFICATIONS_EXCHANGE", "notifications"),
		QueueName:             getEnv("TASK_QUEUE", "task-processor-queue"),
		RoutingKey:            getEnv("TASK_ROUTING_KEY", "task.created"),
	}

	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}
	defer func() {
		if err := application.Close(); err != nil {
			log.Printf("error closing app: %v", err)
		}
	}()

	if err := application.Start(); err != nil {
		log.Fatalf("app error: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
