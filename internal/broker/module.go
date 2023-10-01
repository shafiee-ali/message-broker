package broker

import (
	"context"
	"therealbroker/pkg/broker"
	"therealbroker/pkg/repository"
)

type Module struct {
	repository repository.IMessageRepository
}

func NewModule(repo repository.IMessageRepository) broker.Broker {
	return &Module{repository: repo}
}

func (m *Module) Close() error {
	panic("implement me")
}

func (m *Module) Publish(ctx context.Context, msg broker.CreateMessageDTO) (int, error) {
	createdMessage := m.repository.Add(msg)
	return createdMessage.Id, nil
}

func (m *Module) Subscribe(ctx context.Context, subject string) (<-chan broker.CreateMessageDTO, error) {
	panic("implement me")
}

func (m *Module) Fetch(ctx context.Context, subject string, id int) (broker.CreateMessageDTO, error) {
	panic("implement me")
}
