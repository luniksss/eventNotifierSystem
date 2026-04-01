package sender

import (
	"log"
)

type EmailSender struct {
	mockURL string
}

func NewEmailSender() *EmailSender {
	return &EmailSender{}
}

func (s *EmailSender) Send(to, text string) error {
	log.Printf("email sent to %s: %s", to, text)
	return nil
}
