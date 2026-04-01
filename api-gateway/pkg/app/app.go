package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	amqp "github.com/rabbitmq/amqp091-go"

	"api-gateway/pkg/infrastructure/handler"
	"api-gateway/pkg/infrastructure/producer"
	"api-gateway/pkg/infrastructure/repo"
)

type Config struct {
	DBDSN      string
	RabbitURL  string
	Exchange   string
	RoutingKey string
	HTTPAddr   string
}

type App struct {
	server *http.Server
	db     *sql.DB
	pub    *producer.TaskEventProducer
}

func NewApp(cfg Config) (*App, error) {
	db, err := sql.Open("mysql", cfg.DBDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("db ping failed: %w", err)
	}

	conn, err := amqp.Dial(cfg.RabbitURL)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	pub, err := producer.NewTaskEventProducer(conn, cfg.Exchange, cfg.RoutingKey)
	if err != nil {
		db.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to create publisher: %w", err)
	}

	repo := repo.NewTaskRepository(db)
	taskHandler := handler.NewTaskHandler(repo, pub)

	http.HandleFunc("/tasks", taskHandler.CreateTask)

	server := &http.Server{
		Addr:         cfg.HTTPAddr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &App{
		server: server,
		db:     db,
		pub:    pub,
	}, nil
}

func (a *App) Start() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("api-gateway started on %s", a.server.Addr)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-quit
	log.Println("shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		log.Printf("server shutdown failed: %v", err)
	}

	if err := a.pub.Close(); err != nil {
		log.Printf("publisher close failed: %v", err)
	}
	if err := a.db.Close(); err != nil {
		log.Printf("db close failed: %v", err)
	}

	log.Println("server stopped")
	return nil
}

func (a *App) Close() error {
	if err := a.pub.Close(); err != nil {
		return err
	}
	return a.db.Close()
}
