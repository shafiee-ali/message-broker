package models

import (
	"fmt"
	"time"
)

type PostgresMessage struct {
	ID             int       `gorm:"primary_key"`
	Subject        string    `gorm:"not null"`
	Body           string    `gorm:"not null"`
	ExpirationTime time.Time `gorm:"not null"`
	CreatedAt      time.Time
}

func (msg *PostgresMessage) IsExpired() bool {
	x := time.Now()
	y := msg.ExpirationTime
	fmt.Println(x, y)
	return msg.ExpirationTime.Before(time.Now())
}
