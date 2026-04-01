package repo

import (
	"database/sql"
	"fmt"

	"api-gateway/pkg/app/model"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (tr *TaskRepository) Create(task *model.Task) error {
	query := `INSERT INTO tasks (id, title, description, email, phone, created_at) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := tr.db.Exec(query, task.ID, task.Title, task.Description, task.Email, task.Phone, task.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}
	return nil
}
