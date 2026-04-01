package handler

import (
	"encoding/json"
	"net/http"

	"api-gateway/pkg/app/model"
	"api-gateway/pkg/infrastructure/producer"
	"api-gateway/pkg/infrastructure/repo"
)

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
}

type TaskHandler struct {
	repo     *repo.TaskRepository
	producer *producer.TaskEventProducer
}

func NewTaskHandler(repo *repo.TaskRepository, prod *producer.TaskEventProducer) *TaskHandler {
	return &TaskHandler{
		repo:     repo,
		producer: prod,
	}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	req, err := decodeJSONBody[CreateTaskRequest](r)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	task := model.NewTask(req.Title, req.Description, req.Email, req.Phone)
	if err := h.repo.Create(task); err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	eventData := map[string]interface{}{
		"taskId": task.ID,
		"email":  req.Email,
		"phone":  req.Phone,
		"title":  req.Title,
	}
	ctx := r.Context()
	if err := h.producer.PublishTaskCreated(ctx, eventData); err != nil {
		http.Error(w, "failed to publish task", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": task.ID})
}

func decodeJSONBody[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, err
	}
	return v, nil
}
