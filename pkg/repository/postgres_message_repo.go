package repository

import (
	"sync"
	"therealbroker/internal/mapper"
	"therealbroker/internal/types"
	pkgBroker "therealbroker/pkg/broker"
	"therealbroker/pkg/database"
	"therealbroker/pkg/models"
)

type PostgresRepo struct {
	db              *database.PostgresDB
	insertLock      sync.Mutex
	messageId       int
	incrementIdLock sync.Mutex
}

func NewPostgresRepo(db *database.PostgresDB) PostgresRepo {
	return PostgresRepo{
		db:        db,
		messageId: 0,
	}
}

func (p *PostgresRepo) Add(message pkgBroker.CreateMessageDTO) types.CreatedMessage {
	dbMsg := mapper.CreateMessageDTOToDBMessage(message)

	//p.incrementIdLock.Lock()
	//dbMsg.ID = p.messageId
	//p.messageId++
	//p.incrementIdLock.Unlock()

	p.insertLock.Lock()
	p.db.DB.Create(&dbMsg)
	//logrus.Tracef("Created record %v", dbMsg)
	p.insertLock.Unlock()
	return mapper.DBMessageToCreatedMessage(dbMsg)
}

func (p *PostgresRepo) FetchUnexpiredBySubjectAndId(subject string, id int) (types.CreatedMessage, error) {
	var message models.PostgresMessage
	p.db.DB.Where("subject = ? AND id = ?", subject, id).First(&message)
	if message.IsExpired() {
		return types.CreatedMessage{}, pkgBroker.ErrExpiredID
	}
	return mapper.DBMessageToCreatedMessage(message), nil
}
