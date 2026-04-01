package main

import (
	"log"
	"os"

	"api-gateway/pkg/app"
)

func main() {
	cfg := app.Config{
		DBDSN:      getEnv("DB_DSN", "test:1234@tcp(notifier-mysql:3306)/tasksdb?parseTime=true"),
		RabbitURL:  getEnv("RABBITMQ_URL", "amqp://guest:guest@notifier-rabbitmq:5672/"),
		Exchange:   getEnv("RABBITMQ_EXCHANGE", "tasks"),
		RoutingKey: getEnv("RABBITMQ_ROUTING_KEY", "task.created"),
		HTTPAddr:   getEnv("HTTP_ADDR", ":8080"),
	}

	app, err := app.NewApp(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer app.Close()

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
