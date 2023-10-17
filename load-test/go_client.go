package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math/rand"
	"sync"
	pb "therealbroker/api/proto"
	"time"
)

const ALL_MESSAGES_COUNT = 3000000
const EXPECTED_RATE = 30000
const GOROUTINE_NUMS = 10

func randomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func workerFunc(messages chan pb.PublishRequest, client pb.BrokerClient) {

	tick := time.Tick((time.Second) / (5000 * 3))
	count := 0
	var now time.Time
	for msg := range messages {
		if count == 0 {
			now = time.Now()
		}
		<-tick
		response, err := client.Publish(
			context.Background(),
			&msg,
		)
		count++
		if count == 3000 {
			fmt.Println(time.Since(now).Seconds())
		}
		if err != nil {
			fmt.Printf("Error in publish %v \n", err)
		} else {
			fmt.Printf("Success %v\n", response)
		}
	}

}

func main() {

	tick := time.Tick((time.Second) / (EXPECTED_RATE * 4))

	//var now time.Time
	//log.Println("Hello")
	//for i := 0; i < 30000; i++ {
	//	if i == 0 {
	//		now = time.Now()
	//	}
	//	<-tick
	//	if i == 30000-1 {
	//		//log.Println(time.Now())
	//		log.Println(time.Since(now))
	//	}
	//}

	messages := make(chan pb.PublishRequest, ALL_MESSAGES_COUNT)
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Error in connection to grpc server")
	}

	client := pb.NewBrokerClient(conn)
	wg := sync.WaitGroup{}
	for w := 0; w < GOROUTINE_NUMS; w++ {
		wg.Add(1)
		go workerFunc(messages, client)
	}

	msg := pb.PublishRequest{
		Subject:           randomString(3),
		Body:              []byte(randomString(20)),
		ExpirationSeconds: 1000,
	}

	for j := 1; j <= ALL_MESSAGES_COUNT; j++ {
		<-tick
		messages <- msg
	}
	wg.Wait()

}
