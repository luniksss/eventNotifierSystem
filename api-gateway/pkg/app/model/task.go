package model

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          string    `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	Email       string    `json:"email" db:"email"`
	Phone       string    `json:"phone" db:"phone"`
}

func NewTask(title, desc, email, phone string) *Task {
	return &Task{
		ID:          uuid.New().String(),
		Title:       title,
		Description: desc,
		Email:       email,
		Phone:       phone,
		CreatedAt:   time.Now(),
	}
}
