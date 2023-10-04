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

type CreatedMessageWithoutId struct {
	Subject        string
	Body           string
	ExpirationTime time.Time
}

type Subscriber struct {
	Stream chan CreatedMessageWithoutId
}

func NewCreatedMessageWithoutId(msg CreatedMessage) *CreatedMessageWithoutId {
	return &CreatedMessageWithoutId{
		Subject:        msg.Subject,
		Body:           msg.Body,
		ExpirationTime: msg.ExpirationTime,
	}
}

func NewSubscriber() *Subscriber {
	return &Subscriber{
		Stream: make(chan CreatedMessageWithoutId, 100),
	}
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
