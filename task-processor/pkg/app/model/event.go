package model

import "encoding/json"

type TaskCreatedEvent struct {
	Data struct {
		TaskID string `json:"taskId"`
		Email  string `json:"email"`
		Phone  string `json:"phone"`
		Title  string `json:"title"`
	} `json:"data"`
	Time string `json:"time"`
	Type string `json:"type"`
}

func (e *TaskCreatedEvent) FromJSON(data []byte) error {
	return json.Unmarshal(data, e)
}
