package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math/rand"
	"sync"
	pb "therealbroker/api/proto"
)

func randomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {

	conn, err := grpc.Dial("192.168.49.2:30007", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Error in connection to grpc server")
	}

	client := pb.NewBrokerClient(conn)

	const CONCURRENT_PUBLISH_COUNT = 100000

	wg := sync.WaitGroup{}
	for j := 0; j < 20; j++ {
		for i := 0; i < CONCURRENT_PUBLISH_COUNT; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				response, err := client.Publish(
					context.Background(),
					&pb.PublishRequest{
						Subject:           randomString(2),
						Body:              []byte(randomString(20)),
						ExpirationSeconds: 1000,
					})
				if err != nil {
					fmt.Printf("Error in publish %v \n", err)
				} else {
					fmt.Printf("Success %v\n", response)
				}
			}()
		}
		wg.Wait()
	}

}
