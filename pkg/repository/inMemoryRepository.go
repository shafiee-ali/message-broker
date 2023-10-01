package repository

import (
	"therealbroker/internal/types"
	pkgBroker "therealbroker/pkg/broker"
	"therealbroker/pkg/models"
	"time"
)

type InMemoryDB struct {
	items []models.Message
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{make([]models.Message, 0, 0)}
}

type IMessageRepository interface {
	Add(message pkgBroker.CreateMessageDTO) types.CreatedMessage
}

func (db *InMemoryDB) GenerateNewId() int {
	return len(db.items) + 1
}

func (db *InMemoryDB) Add(message pkgBroker.CreateMessageDTO) types.CreatedMessage {
	id := db.GenerateNewId()
	newMessage := models.Message{
		Id:             id,
		Body:           message.Body,
		Subject:        message.Subject,
		ExpirationTime: message.ExpirationTime,
		CreatedAt:      time.Now(),
	}
	db.items = append(db.items, newMessage)
	return types.NewCreatedMessage(newMessage)
}
