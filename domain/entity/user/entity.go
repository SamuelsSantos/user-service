package user

import (
	"time"

	"github.com/samuelssantos/user-service/domain/entity"
)

//User data
type User struct {
	ID          entity.ID
	Email       string
	Password    string
	FirstName   string
	LastName    string
	DateOfBirth time.Time
}
