package repository

import (
	log "github.com/sirupsen/logrus"
	"sync"
	"therealbroker/internal/mapper"
	"therealbroker/internal/types"
	pkgBroker "therealbroker/pkg/broker"
	"therealbroker/pkg/database"
	"therealbroker/pkg/models"
	"time"
)

type PostgresRepo struct {
	db              *database.PostgresDB
	insertLock      *sync.Mutex
	incrementIdLock *sync.Mutex
	id              int
	ticker          *time.Ticker
	batchMessages   []models.PostgresMessage
}

func NewPostgresRepo(db *database.PostgresDB) *PostgresRepo {

	pgRepo := &PostgresRepo{
		db:              db,
		ticker:          time.NewTicker(time.Millisecond * 500),
		id:              -1,
		insertLock:      &sync.Mutex{},
		incrementIdLock: &sync.Mutex{},
	}
	log.Infof("After creating pg repo obj")
	pgRepo.createMessagesInBatch()
	return pgRepo
}

func (p *PostgresRepo) NextId() int {
	if p.id != -1 {
		p.id++
	} else {
		var lastCreatedMessage *models.PostgresMessage
		p.db.DB.Raw("SELECT * FROM postgres_messages ORDER BY created_at DESC LIMIT 1;").Scan(&lastCreatedMessage)
		log.Infof("Founded row %v", lastCreatedMessage)
		if lastCreatedMessage == nil {
			p.id = 0
		} else {
			p.id = lastCreatedMessage.ID + 1
		}
	}
	return p.id
}

func (p *PostgresRepo) Add(message pkgBroker.CreateMessageDTO) types.CreatedMessage {
	dbMsg := mapper.CreateMessageDTOToDBMessage(message)
	p.insertLock.Lock()
	dbMsg.ID = p.NextId()
	p.batchMessages = append(p.batchMessages, dbMsg)
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

func (p *PostgresRepo) createMessagesInBatch() {
	go func() {
		for {
			select {
			case <-p.ticker.C:
				if len(p.batchMessages) == 0 {
					continue
				}
				log.Infof("Create batch with length %v", len(p.batchMessages))
				p.insertLock.Lock()
				messagesForInsertion := p.batchMessages
				p.batchMessages = make([]models.PostgresMessage, 0)
				p.insertLock.Unlock()
				p.db.DB.CreateInBatches(messagesForInsertion, len(messagesForInsertion))
			}
		}
	}()
}
