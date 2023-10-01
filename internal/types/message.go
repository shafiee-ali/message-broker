package types

import (
	"therealbroker/pkg/models"
	"time"
)

type CreatedMessage struct {
	Id             int
	Subject        string
	Body           string
	ExpirationTime time.Time
}

func NewCreatedMessage(msgModel models.Message) CreatedMessage {
	return CreatedMessage{
		Id:             msgModel.Id,
		Subject:        msgModel.Subject,
		Body:           msgModel.Body,
		ExpirationTime: msgModel.ExpirationTime,
	}
}

func EmptyCreatedMessage() CreatedMessage {
	return CreatedMessage{}
}
