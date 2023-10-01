package broker

import (
	"context"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"therealbroker/pkg/broker"
	"therealbroker/pkg/repository"
	"time"
)

var (
	service broker.Broker
	mainCtx = context.Background()
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func TestPublishShouldSucceedAndReturnMessageId(t *testing.T) {
	msg := createMessageWithExpire("ALI", time.Millisecond*500)
	inMemoryRepo := repository.NewInMemoryDB()
	brokerService := NewModule(inMemoryRepo)
	id, err := brokerService.Publish(mainCtx, msg)
	assert.Nil(t, err)
	assert.Equal(t, 1, id)

}

func TestPublishShouldSucceedSecondTimeAndReturnTwoAsId(t *testing.T) {
	msg := createMessageWithExpire("ALI", time.Millisecond*500)
	inMemoryRepo := repository.NewInMemoryDB()
	brokerService := NewModule(inMemoryRepo)
	id, err := brokerService.Publish(mainCtx, msg)
	id, err = brokerService.Publish(mainCtx, msg)
	assert.Nil(t, err)
	assert.Equal(t, 2, id)
}

func randomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func createMessageWithExpire(subject string, duration time.Duration) broker.CreateMessageDTO {
	body := randomString(16)

	return broker.CreateMessageDTO{Subject: subject, Body: body, ExpirationTime: broker.CalcExpirationTime(duration)}
}
