package broker

import (
	"context"
	"io"
	"therealbroker/internal/types"
	"time"
)

type CreateMessageDTO struct {
	Subject        string
	Body           string
	ExpirationTime time.Time
}

func NewCreateMessageDTO(subject string, body string, expiration int32) CreateMessageDTO {
	return CreateMessageDTO{subject, body, CalcExpirationTime(time.Duration(expiration))}
}

type Broker interface {
	io.Closer
	Publish(ctx context.Context, msg CreateMessageDTO) (int, error)
	Subscribe(ctx context.Context, subject string) (<-chan types.CreatedMessage, error)
	Fetch(ctx context.Context, subject string, id int) (types.CreatedMessage, error)
}

func CalcExpirationTime(d time.Duration) time.Time {
	return time.Now().Add(d)
}
