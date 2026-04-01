package model

type Notification struct {
	TaskID string `json:"taskId"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	Text   string `json:"text"`
}

func NewNotificationFromEvent(event *TaskCreatedEvent) *Notification {
	return &Notification{
		TaskID: event.Data.TaskID,
		Email:  event.Data.Email,
		Phone:  event.Data.Phone,
		Text:   "Task created: " + event.Data.Title,
	}
}
