package broker

import (
	"context"
	"sync"
	"therealbroker/pkg/broker"
	"therealbroker/pkg/repository"
)

type Module struct {
	repository repository.IMessageRepository
	lock       sync.Mutex
	isClosed   bool
}

func NewModule(repo repository.IMessageRepository) broker.Broker {
	return &Module{repository: repo}
}

func (m *Module) Close() error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.isClosed {
		return broker.ErrUnavailable
	}
	m.isClosed = true
	return nil
}

func (m *Module) Publish(ctx context.Context, msg broker.CreateMessageDTO) (int, error) {
	err := m.checkServerDown()
	if err != nil {
		return -1, err
	}
	createdMessage := m.repository.Add(msg)
	return createdMessage.Id, nil
}

func (m *Module) Subscribe(ctx context.Context, subject string) (<-chan broker.CreateMessageDTO, error) {
	panic("implement me")
}

func (m *Module) Fetch(ctx context.Context, subject string, id int) (broker.CreateMessageDTO, error) {
	panic("implement me")
}

func (m *Module) checkServerDown() error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.isClosed {
		return broker.ErrUnavailable
	}
	return nil
}
