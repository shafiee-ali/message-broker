package broker

import (
	"context"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sync"
	"testing"
	"therealbroker/internal/types"
	"therealbroker/pkg/broker"
	"therealbroker/pkg/repository"
	"time"
)

var (
	service broker.Broker
	mainCtx = context.Background()
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func TestMain(m *testing.M) {
	rand.Seed(time.Now().Unix())
	inMemoryRepo := repository.NewInMemoryDB()
	service = NewModule(inMemoryRepo)
	m.Run()
}

func Setup() {
	inMemoryRepo := repository.NewInMemoryDB()
	service = NewModule(inMemoryRepo)
}

func TestPublishShouldFailOnClosed(t *testing.T) {
	Setup()
	msg := createMessage("ali")

	err := service.Close()
	assert.Nil(t, err)

	_, err = service.Publish(mainCtx, msg)
	assert.Equal(t, broker.ErrUnavailable, err)
}

func TestSubscribeShouldFailOnClosed(t *testing.T) {
	Setup()
	err := service.Close()
	assert.Nil(t, err)

	_, err = service.Subscribe(mainCtx, "ali")
	assert.Equal(t, broker.ErrUnavailable, err)
}

func TestFetchShouldFailOnClosed(t *testing.T) {
	Setup()
	err := service.Close()
	assert.Nil(t, err)

	_, err = service.Fetch(mainCtx, "ali", rand.Intn(100))
	assert.Equal(t, broker.ErrUnavailable, err)
}

func TestPublishShouldNotFail(t *testing.T) {
	Setup()
	msg := createMessage("ali")

	_, err := service.Publish(mainCtx, msg)

	assert.Equal(t, nil, err)
}

func TestSubscribeShouldNotFail(t *testing.T) {
	Setup()
	sub, err := service.Subscribe(mainCtx, "ali")

	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, sub)
}

func TestPublishShouldSendMessageToSubscribedChan(t *testing.T) {
	Setup()
	msg := createMessage("ali")

	sub, _ := service.Subscribe(mainCtx, "ali")
	_, _ = service.Publish(mainCtx, msg)
	in := <-sub

	assert.Equal(t, msg.Body, in.Body)
	assert.Equal(t, msg.Subject, in.Subject)
	assert.Equal(t, msg.ExpirationTime, in.ExpirationTime)
}

func TestPublishShouldSendMessageToSubscribedChans(t *testing.T) {
	Setup()
	msg := createMessage("ali")

	sub1, _ := service.Subscribe(mainCtx, "ali")
	sub2, _ := service.Subscribe(mainCtx, "ali")
	sub3, _ := service.Subscribe(mainCtx, "ali")
	_, _ = service.Publish(mainCtx, msg)
	in1 := <-sub1
	in2 := <-sub2
	in3 := <-sub3

	assert.Equal(t, msg.Subject, in1.Subject)
	assert.Equal(t, msg.Body, in1.Body)
	assert.Equal(t, msg.ExpirationTime, in1.ExpirationTime)

	assert.Equal(t, msg.Subject, in2.Subject)
	assert.Equal(t, msg.Body, in2.Body)
	assert.Equal(t, msg.ExpirationTime, in2.ExpirationTime)

	assert.Equal(t, msg.Subject, in3.Subject)
	assert.Equal(t, msg.Body, in3.Body)
	assert.Equal(t, msg.ExpirationTime, in3.ExpirationTime)

}

func TestPublishShouldPreserveOrder(t *testing.T) {
	Setup()
	n := 50
	messages := make([]broker.CreateMessageDTO, n)
	sub, _ := service.Subscribe(mainCtx, "ali")
	for i := 0; i < n; i++ {
		messages[i] = createMessage("ali")
		_, _ = service.Publish(mainCtx, messages[i])
	}

	for i := 0; i < n; i++ {
		msg := <-sub
		assert.Equal(t, messages[i].Subject, msg.Subject)
		assert.Equal(t, messages[i].Body, msg.Body)
		assert.Equal(t, messages[i].ExpirationTime, msg.ExpirationTime)
	}
}

func TestPublishShouldNotSendToOtherSubscriptions(t *testing.T) {
	Setup()
	msg := createMessage("ali")
	ali, _ := service.Subscribe(mainCtx, "ali")
	maryam, _ := service.Subscribe(mainCtx, "maryam")

	_, _ = service.Publish(mainCtx, msg)
	select {
	case m := <-ali:
		assert.Equal(t, msg.Subject, m.Subject)
		assert.Equal(t, msg.Body, m.Body)
		assert.Equal(t, msg.ExpirationTime, m.ExpirationTime)
	case <-maryam:
		assert.Fail(t, "Wrong message received")
	}
}

func TestNonExpiredMessageShouldBeFetchable(t *testing.T) {
	Setup()
	msg := createMessageWithExpire("ali", time.Second*10)
	id, _ := service.Publish(mainCtx, msg)
	fMsg, _ := service.Fetch(mainCtx, "ali", id)

	assert.Equal(t, msg.Subject, fMsg.Subject)
	assert.Equal(t, msg.Body, fMsg.Body)
	assert.Equal(t, msg.ExpirationTime, fMsg.ExpirationTime)
}

func TestExpiredMessageShouldNotBeFetchable(t *testing.T) {
	Setup()
	msg := createMessageWithExpire("ali", time.Millisecond*500)
	id, _ := service.Publish(mainCtx, msg)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	<-ticker.C
	fMsg, err := service.Fetch(mainCtx, "ali", id)
	assert.Equal(t, broker.ErrExpiredID, err)
	assert.Equal(t, types.CreatedMessage{}, fMsg)
}

func TestNewSubscriptionShouldNotGetPreviousMessages(t *testing.T) {
	Setup()
	msg := createMessage("ali")
	_, _ = service.Publish(mainCtx, msg)
	sub, _ := service.Subscribe(mainCtx, "ali")

	select {
	case <-sub:
		assert.Fail(t, "Got previous message")
	default:
	}
}

func TestConcurrentSubscribesOnOneSubjectShouldNotFail(t *testing.T) {
	Setup()
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	var wg sync.WaitGroup

	for {
		select {
		case <-ticker.C:
			wg.Wait()
			return

		default:
			wg.Add(1)
			go func() {
				defer wg.Done()

				_, err := service.Subscribe(mainCtx, "ali")
				assert.Nil(t, err)
			}()
		}
	}
}

func TestConcurrentSubscribesShouldNotFail(t *testing.T) {
	Setup()
	ticker := time.NewTicker(2000 * time.Millisecond)
	defer ticker.Stop()
	var wg sync.WaitGroup

	for {
		select {
		case <-ticker.C:
			wg.Wait()
			return

		default:
			wg.Add(1)
			go func() {
				defer wg.Done()

				_, err := service.Subscribe(mainCtx, randomString(4))
				assert.Nil(t, err)
			}()
		}
	}
}

func TestConcurrentPublishOnOneSubjectShouldNotFail(t *testing.T) {
	Setup()
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	var wg sync.WaitGroup

	msg := createMessage("ali")

	for {
		select {
		case <-ticker.C:
			wg.Wait()
			return

		default:
			wg.Add(1)
			go func() {
				defer wg.Done()

				_, err := service.Publish(mainCtx, msg)
				assert.Nil(t, err)
			}()
		}
	}
}

func TestConcurrentPublishShouldNotFail(t *testing.T) {
	Setup()
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	var wg sync.WaitGroup

	msg := createMessage(randomString(4))

	for {
		select {
		case <-ticker.C:
			wg.Wait()
			return

		default:
			wg.Add(1)
			go func() {
				defer wg.Done()

				_, err := service.Publish(mainCtx, msg)
				assert.Nil(t, err)
			}()
		}
	}
}

func TestDataRace(t *testing.T) {
	Setup()
	duration := 500 * time.Millisecond
	ticker := time.NewTicker(duration)
	defer ticker.Stop()
	var wg sync.WaitGroup

	ids := make(chan int, 100000)

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-ticker.C:
				return

			default:
				id, err := service.Publish(mainCtx, createMessageWithExpire("ali", duration))
				ids <- id
				assert.Nil(t, err)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-ticker.C:
				return

			default:
				_, err := service.Subscribe(mainCtx, "ali")
				assert.Nil(t, err)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-ticker.C:
				return

			case id := <-ids:
				_, err := service.Fetch(mainCtx, "ali", id)
				assert.Nil(t, err)
			}
		}
	}()

	wg.Wait()
}

func BenchmarkPublish(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := service.Publish(mainCtx, createMessage(randomString(2)))
		assert.Nil(b, err)
	}
}

func BenchmarkSubscribe(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := service.Subscribe(mainCtx, randomString(2))
		assert.Nil(b, err)
	}
}
func createMessage(subject string) broker.CreateMessageDTO {
	body := randomString(16)

	return broker.NewCreateMessageDTO(subject, body, int32(0))
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
