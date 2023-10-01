package models

import "time"

type Message struct {
	Id             int
	Subject        string
	Body           string
	ExpirationTime time.Time
	CreatedAt      time.Time
}
