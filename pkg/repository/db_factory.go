package repository

import "therealbroker/pkg/database"

func RepoFactory(dbName string) IMessageRepository {
	if dbName == "CASSANDRA" {
		db := database.NewCassandra()
		return NewCassandraRepo(db)
	}
	if dbName == "POSTGRES" {
		dbConfig := database.DBConfig{
			User:   "postgres",
			Pass:   "password",
			DbName: "broker",
			Port:   5432,
			Host:   "postgres",
		}
		db := database.NewPostgres(dbConfig)
		return NewPostgresRepo(db)
	}
	if dbName == "MEMORY" {
		return NewInMemoryMessageDB()
	}
	return nil
}
