package models

import "time"

type Message struct {
	Id             int
	Subject        string
	Body           string
	expirationTime time.Time
}
