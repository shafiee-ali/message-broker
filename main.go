package main

import (
	"fmt"
	"therealbroker/api"
)

// Main requirements:
// 1. All tests should be passed
// 2. Your logs should be accessible in Graylog
// 3. Basic prometheus metrics ( latency, throughput, etc. ) should be implemented
// 	  for every base functionality ( publish, subscribe etc. )

func main() {
	fmt.Println("Hello!")
	api.StartGrpcServer()
}
