package models

import (
	"fmt"
	"time"
)

type CassandraMessage struct {
	ID             int
	Subject        string
	Body           string
	ExpirationTime time.Time
	CreatedAt      time.Time
}

func (msg *CassandraMessage) IsExpired() bool {
	x := time.Now()
	y := msg.ExpirationTime
	fmt.Println(x, y)
	return msg.ExpirationTime.Before(time.Now())
}
