package model

type Notification struct {
	TaskID string `json:"taskId"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	Text   string `json:"text"`
}
