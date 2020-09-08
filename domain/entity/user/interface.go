package user

import "github.com/samuelssantos/user-service/domain/entity"

//Reader interface
type Reader interface {
	Get(id entity.ID) (*User, error)
	Search(query string) ([]*User, error)
	List() ([]*User, error)
}

//Writer user writer
type Writer interface {
	Create(e *User) (entity.ID, error)
	Update(e *User) error
	Delete(id entity.ID) error
}

//repository interface
type repository interface {
	Reader
	Writer
}

//Manager interface
type Manager interface {
	repository
}
