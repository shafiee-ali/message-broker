package main

import (
	log "github.com/sirupsen/logrus"
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
	DB_TYPE := database.CASSANDRA

	repo := repository.RepoFactory(DB_TYPE)
	broker := broker.NewModule(repo)
	metrics := metrics.StartPrometheus()

	log.SetLevel(log.TraceLevel)
	api.StartGrpcServer(broker, metrics)
}
