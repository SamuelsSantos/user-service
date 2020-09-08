package presenter

import (
	"github.com/samuelssantos/user-service/domain/entity"
)

//User data
type User struct {
	ID          entity.ID `json:"id"`
	Email       string    `json:"email"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	DateOfBirth string    `json:"date_of_birth,omitempty"`
}
