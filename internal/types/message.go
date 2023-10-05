package types

import (
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

func EmptyCreatedMessage() CreatedMessage {
	return CreatedMessage{}
}
