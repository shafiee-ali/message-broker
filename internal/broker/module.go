package broker

import (
	"context"
	log "github.com/sirupsen/logrus"
	"sync"
	"therealbroker/internal/types"
	"therealbroker/pkg/broker"
	"therealbroker/pkg/repository"
)

type Module struct {
	subscribers       map[string][]*types.Subscriber
	repository        repository.IMessageRepository
	statusLock        sync.Mutex
	addSubscriberLock sync.Mutex
	publishLock       sync.RWMutex
	isClosed          bool
}

func NewModule(repo repository.IMessageRepository) broker.Broker {
	return &Module{
		repository:  repo,
		subscribers: make(map[string][]*types.Subscriber),
		isClosed:    false,
	}
}

func (m *Module) Close() error {
	m.statusLock.Lock()
	defer m.statusLock.Unlock()
	m.isClosed = true
	return nil
}

func (m *Module) Publish(ctx context.Context, msg broker.CreateMessageDTO) (int, error) {
	err := m.checkServerDown()
	if err != nil {
		log.Infof("Server is closed...")
		return -1, err
	}
	//go func() {
	//	m.repository.Add(msg)
	//}()
	createdMessage := m.repository.Add(msg)
	//_,
	m.publishLock.RLock()
	m.SendMessageToSubscribers(createdMessage)
	//log.Tracef("Message send to subscribers")
	m.publishLock.RUnlock()
	return createdMessage.Id, nil
}

func (m *Module) SendMessageToSubscribers(msg types.CreatedMessage) {
	msgWithoutId := types.NewCreatedMessageWithoutId(msg)
	wg := sync.WaitGroup{}
	for _, subscriber := range m.subscribers[msg.Subject] {
		wg.Add(1)
		go func(sub *types.Subscriber) {
			defer wg.Done()
			sub.Stream <- *msgWithoutId
		}(subscriber)
	}
	wg.Wait()
}

func (m *Module) Subscribe(ctx context.Context, subject string) (<-chan types.CreatedMessageWithoutId, error) {
	err := m.checkServerDown()
	if err != nil {
		return nil, broker.ErrUnavailable
	}

	m.addSubscriberLock.Lock()
	newSub := types.NewSubscriber()
	m.subscribers[subject] = append(m.subscribers[subject], newSub)
	m.addSubscriberLock.Unlock()
	return newSub.Stream, nil
}

func (m *Module) Fetch(ctx context.Context, subject string, id int) (types.CreatedMessage, error) {
	err := m.checkServerDown()
	if err != nil {
		return types.EmptyCreatedMessage(), broker.ErrUnavailable
	}
	return m.repository.FetchUnexpiredBySubjectAndId(subject, id)
}

func (m *Module) checkServerDown() error {
	m.statusLock.Lock()
	defer m.statusLock.Unlock()
	if m.isClosed {
		return broker.ErrUnavailable
	}
	return nil
}
