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
	batchMessage    []models.PostgresMessage
	creationAcks    []chan bool
}

func NewPostgresRepo(db *database.PostgresDB) *PostgresRepo {

	pgRepo := &PostgresRepo{
		db:              db,
		ticker:          time.NewTicker(time.Millisecond * 500),
		id:              -1,
		insertLock:      &sync.Mutex{},
		incrementIdLock: &sync.Mutex{},
		creationAcks:    make([]chan bool, 0),
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
	dbMsg := mapper.CreateMessageDTOToPostgresMessage(message)
	p.insertLock.Lock()
	dbMsg.ID = p.NextId()
	ch := make(chan bool)
	p.creationAcks = append(p.creationAcks, ch)
	p.batchMessage = append(p.batchMessage, dbMsg)
	p.insertLock.Unlock()

	<-ch
	return mapper.PostgresMessageToCreatedMessage(dbMsg)
}

func (p *PostgresRepo) FetchUnexpiredBySubjectAndId(subject string, id int) (types.CreatedMessage, error) {
	var message models.PostgresMessage
	p.db.DB.Where("subject = ? AND id = ?", subject, id).First(&message)
	if message.IsExpired() {
		return types.CreatedMessage{}, pkgBroker.ErrExpiredID
	}
	return mapper.PostgresMessageToCreatedMessage(message), nil
}

func (p *PostgresRepo) createMessagesInBatch() {
	go func() {
		for {
			select {
			case <-p.ticker.C:
				if len(p.batchMessage) == 0 {
					continue
				}
				log.Infof("Create batch with length %v", len(p.batchMessage))
				p.insertLock.Lock()
				messagesForInsertion := p.batchMessage
				ackChannels := p.creationAcks
				p.batchMessage = make([]models.PostgresMessage, 0)
				p.creationAcks = make([]chan bool, 0)
				p.insertLock.Unlock()
				p.db.DB.CreateInBatches(messagesForInsertion, len(messagesForInsertion))
				for _, ch := range ackChannels {
					ch <- true
				}
			}
		}
	}()
}
