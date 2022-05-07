package models

import "time"

type User struct {
	ID             int64
	ChatID         int64
	UserName       string
	FirstName      string
	LastName       string
	DateRegistered time.Time
}
