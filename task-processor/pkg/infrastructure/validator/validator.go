package validator

import (
	"taskprocessor/pkg/app/model"
)

func ValidateEvent(event *model.TaskCreatedEvent) bool {
	if event.Data.TaskID == "" {
		return false
	}
	if event.Data.Email == "" && event.Data.Phone == "" {
		return false
	}
	return true
}
