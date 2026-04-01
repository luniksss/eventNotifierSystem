package sender

import (
	"log"
)

type SMSSender struct {
	mockURL string
}

func NewSMSSender() *SMSSender {
	return &SMSSender{}
}

func (s *SMSSender) Send(to, text string) error {
	log.Printf("sms sent to %s: %s", to, text)
	return nil
}
