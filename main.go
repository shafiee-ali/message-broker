package main

import (
	"fmt"
	"therealbroker/api"
	"therealbroker/internal/broker"
	"therealbroker/pkg/database"
	"therealbroker/pkg/metrics"
	"therealbroker/pkg/repository"
)

// Main requirements:
// 1. All tests should be passed
// 2. Your logs should be accessible in Graylog
// 3. Basic prometheus metrics ( latency, throughput, etc. ) should be implemented
// 	  for every base functionality ( publish, subscribe etc. )

func main() {
	fmt.Println("Server started...")
	dbConfig := database.DBConfig{
		User:   "root",
		Pass:   "password",
		DbName: "broker",
		Port:   5434,
		Host:   "localhost",
	}
	db := database.NewPostgres(dbConfig)
	postgresRepo := repository.NewPostgresRepo(db)
	broker := broker.NewModule(&postgresRepo)
	metrics := metrics.StartPrometheus()
	api.StartGrpcServer(broker, metrics)
}
