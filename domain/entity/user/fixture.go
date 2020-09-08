package user

import (
	"time"

	"github.com/samuelssantos/user-service/domain/entity"
)

const (
	layoutISO = "2006-01-02"
)

// NewFixtureUser ...
func NewFixtureUser() *User {
	dateOfBirth, _ := time.Parse(layoutISO, "2006-01-02")
	return &User{
		ID:          entity.NewID(),
		Password:    "123456",
		FirstName:   "zelda",
		LastName:    "zica",
		DateOfBirth: dateOfBirth,
	}
}
