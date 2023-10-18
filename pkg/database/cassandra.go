package database

import (
	"fmt"
	"github.com/gocql/gocql"
	log "github.com/sirupsen/logrus"
)

type CassandraDB struct {
	Session *gocql.Session
}

func NewCassandra() *CassandraDB {
	cluster := gocql.NewCluster("cassandra")
	fmt.Println("Cas connected", cluster)
	cluster.Keyspace = "broker"
	cluster.Consistency = gocql.Quorum
	//cluster.Authenticator = gocql.PasswordAuthenticator{
	//	Username: "cassandra",
	//	Password: "password",
	//}

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Unbale connect to cassandra")
	}
	fmt.Println("Session is created")
	return &CassandraDB{
		Session: session,
	}
}
