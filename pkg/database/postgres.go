package database

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"therealbroker/pkg/models"
)

type DBConfig struct {
	Host   string
	Port   int
	User   string
	Pass   string
	DbName string
}

type PostgresDB struct {
	DB *gorm.DB
}

func NewPostgres(config DBConfig) *PostgresDB {
	dbConfig := fmt.Sprintf(
		"host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Pass,
		config.DbName,
	)

	dbConnection, err := gorm.Open(postgres.Open(dbConfig), &gorm.Config{})
	if err != nil {
		log.Fatalf("Unable connect to database %v", err)
	}

	dbConnection.AutoMigrate(models.PostgresMessage{})
	return &PostgresDB{dbConnection}
}