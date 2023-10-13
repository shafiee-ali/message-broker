package main

import (
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
	dbConfig := database.DBConfig{
		User:   "postgres",
		Pass:   "password",
		DbName: "broker",
		Port:   5432,
		Host:   "postgres",
	}
	db := database.NewPostgres(dbConfig)
	postgresRepo := repository.NewPostgresRepo(db)
	broker := broker.NewModule(&postgresRepo)
	metrics := metrics.StartPrometheus()
	//log.SetLevel(log.TraceLevel)
	api.StartGrpcServer(broker, metrics)
}
