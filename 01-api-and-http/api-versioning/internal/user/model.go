package user

import "time"

type Status string

const (
	StatusActive   Status = "ACTIVE"
	StatusInactive Status = "INACTIVE"
)

type User struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	Status    Status
	CreatedAt time.Time
}
