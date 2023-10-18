package repository

import (
	"github.com/gocql/gocql"
	log "github.com/sirupsen/logrus"
	"strconv"
	"sync"
	"therealbroker/internal/mapper"
	"therealbroker/internal/types"
	pkgBroker "therealbroker/pkg/broker"
	"therealbroker/pkg/database"
	"therealbroker/pkg/models"
	"time"
)

const CREATE_MESSAGE_TABLE_QUERY = "CREATE TABLE messages (id INT PRIMARY KEY, subject TEXT, body TEXT, expiration_time TIMESTAMP, created_at TIMESTAMP) ;\n"

type CassandraRepo struct {
	db              *database.CassandraDB
	insertLock      *sync.Mutex
	incrementIdLock *sync.Mutex
	id              int
	ticker          *time.Ticker
	batchMessages   []models.CassandraMessage
	creationAcks    []chan bool
}

func (cas *CassandraRepo) createMessageTable() {
	keySpace, _ := cas.db.Session.KeyspaceMetadata("broker")
	_, exists := keySpace.Tables["messages"]
	if exists == false {
		err := cas.db.Session.Query(CREATE_MESSAGE_TABLE_QUERY).Exec()
		if err != nil {
			log.Fatalf("Message table creation failed %v", err)
		}

	}
}

func NewCassandraRepo(db *database.CassandraDB) *CassandraRepo {
	casRepo := &CassandraRepo{
		db:              db,
		ticker:          time.NewTicker(time.Microsecond * 10),
		id:              -1,
		insertLock:      &sync.Mutex{},
		incrementIdLock: &sync.Mutex{},
		creationAcks:    make([]chan bool, 0),
	}
	casRepo.createMessageTable()
	log.Infof("After creating cas repo obj")
	casRepo.createMessagesInBatch()
	return casRepo
}

func (cas *CassandraRepo) NextId() int {
	if cas.id != -1 {
		cas.id++
	} else {
		var lastCreatedMessageId *int
		err := cas.db.Session.Query("SELECT * FROM messages ORDER BY created_at DESC LIMIT 1;").Scan(&lastCreatedMessageId)
		if lastCreatedMessageId == nil || err != nil {
			cas.id = 0
		} else {
			cas.id = *lastCreatedMessageId + 1
		}
	}
	return cas.id
}

func (cas *CassandraRepo) Add(message pkgBroker.CreateMessageDTO) types.CreatedMessage {
	dbMsg := mapper.CreateMessageDTOToCassandraMessage(message)
	cas.insertLock.Lock()
	dbMsg.ID = cas.NextId()
	ch := make(chan bool)
	cas.creationAcks = append(cas.creationAcks, ch)
	cas.batchMessages = append(cas.batchMessages, dbMsg)
	cas.insertLock.Unlock()
	return mapper.CassandraMessageToCreatedMessage(dbMsg)
}

func (cas *CassandraRepo) FetchUnexpiredBySubjectAndId(subject string, id int) (types.CreatedMessage, error) {
	var message models.CassandraMessage
	var body string
	var expirationTime time.Time
	var createdAt time.Time
	cas.db.Session.
		Query(
			"SELECT body, expirationTime, created_at FROM messages WHERE"+"id="+strconv.Itoa(id)+" AND "+"subject="+subject,
		).Scan(
		&body,
		&expirationTime,
		&createdAt,
	)
	message = models.CassandraMessage{
		ID:             id,
		Subject:        subject,
		Body:           body,
		ExpirationTime: expirationTime,
		CreatedAt:      createdAt,
	}
	if message.IsExpired() {
		return types.CreatedMessage{}, pkgBroker.ErrExpiredID
	}
	return mapper.CassandraMessageToCreatedMessage(message), nil
}

func (cas *CassandraRepo) createMessagesInBatch() {
	go func() {
		for {
			select {
			case <-cas.ticker.C:
				if len(cas.batchMessages) == 0 {
					log.Infof("Batch size is zero")
					continue
				}
				cas.insertLock.Lock()
				messagesForInsertion := cas.batchMessages
				ackChannels := cas.creationAcks
				cas.creationAcks = make([]chan bool, 0)
				cas.batchMessages = make([]models.CassandraMessage, 0)
				cas.insertLock.Unlock()
				batch := cas.db.Session.NewBatch(gocql.UnloggedBatch)
				log.Infof("After create new batch obj")
				stmt := "INSERT INTO messages (id, subject, body, expiration_time, created_at) VALUES (?, ?, ?, ?, toTimestamp(now()))"
				for i, msg := range messagesForInsertion {
					batch.Query(stmt, msg.ID, msg.Subject, msg.Body, msg.ExpirationTime)
					if i != 0 && i%500 == 0 {
						err := cas.db.Session.ExecuteBatch(batch)
						log.Infoln("After exec batch query")
						if err != nil {
							log.Errorln("ExecuteBatch error", err)
						}

					}
				}
				if batch.Size() > 0 {
					err := cas.db.Session.ExecuteBatch(batch)
					if err != nil {
						log.Errorln("ExecuteBatch error", err)
					}
				}
				log.Infof("ack channel size %v", len(ackChannels))
				for _, ch := range ackChannels {
					ch <- true
				}
				log.Infof("Ending one batch")
			}
		}
	}()
}
