package repository

import (
	"sync"
	"therealbroker/internal/types"
	pkgBroker "therealbroker/pkg/broker"
	"therealbroker/pkg/models"
	"time"
)

type InMemoryDB struct {
	dbLock sync.Mutex
	items  []models.Message
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		sync.Mutex{},
		make([]models.Message, 0, 0),
	}
}

type IMessageRepository interface {
	Add(message pkgBroker.CreateMessageDTO) types.CreatedMessage
	FetchUnexpiredBySubjectAndId(subject string, id int) (types.CreatedMessage, error)
}

func (db *InMemoryDB) GenerateNewId() int {
	return len(db.items) + 1
}

func (db *InMemoryDB) FetchUnexpiredBySubjectAndId(subject string, id int) (types.CreatedMessage, error) {
	for _, msg := range db.items {
		if msg.Id == id && msg.Subject == subject {
			if msg.IsExpired() {
				return types.EmptyCreatedMessage(), pkgBroker.ErrExpiredID
			} else {
				return types.NewCreatedMessage(msg), nil
			}
		}
	}
	return types.EmptyCreatedMessage(), pkgBroker.ErrInvalidID
}

func (db *InMemoryDB) Add(message pkgBroker.CreateMessageDTO) types.CreatedMessage {
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
	return types.NewCreatedMessage(newMessage)
}
