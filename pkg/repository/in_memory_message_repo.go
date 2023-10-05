package repository

import (
	"sync"
	"therealbroker/internal/mapper"
	"therealbroker/internal/types"
	pkgBroker "therealbroker/pkg/broker"
	"therealbroker/pkg/models"
	"time"
)

type InMemoryMessageDB struct {
	dbLock sync.Mutex
	items  []models.Message
}

func NewInMemoryMessageDB() *InMemoryMessageDB {
	return &InMemoryMessageDB{
		sync.Mutex{},
		make([]models.Message, 0, 0),
	}
}

func (db *InMemoryMessageDB) GenerateNewId() int {
	return len(db.items) + 1
}

func (db *InMemoryMessageDB) FetchUnexpiredBySubjectAndId(subject string, id int) (types.CreatedMessage, error) {
	for _, msg := range db.items {
		if msg.Id == id && msg.Subject == subject {
			if msg.IsExpired() {
				return types.EmptyCreatedMessage(), pkgBroker.ErrExpiredID
			} else {
				return mapper.InMemoryDBMessageToCreatedMessage(msg), nil
			}
		}
	}
	return types.EmptyCreatedMessage(), pkgBroker.ErrInvalidID
}

func (db *InMemoryMessageDB) Add(message pkgBroker.CreateMessageDTO) types.CreatedMessage {
	db.dbLock.Lock()
	defer db.dbLock.Unlock()
	id := db.GenerateNewId()
	newMessage := models.Message{
		Id:             id,
		Body:           message.Body,
		Subject:        message.Subject,
		ExpirationTime: message.ExpirationTime,
		CreatedAt:      time.Now(),
	}
	db.items = append(db.items, newMessage)
	return mapper.InMemoryDBMessageToCreatedMessage(newMessage)
}
