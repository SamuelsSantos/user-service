package user

import (
	"strings"

	"github.com/samuelssantos/user-service/pkg/password"

	"github.com/samuelssantos/user-service/domain/entity"
)

//manager  interface
type manager struct {
	repo repository
	pwd  password.Service
}

//NewManager create new repository
func NewManager(r repository, pwd password.Service) *manager {
	return &manager{
		repo: r,
		pwd:  pwd,
	}
}

//Create an user
func (s *manager) Create(e *User) (entity.ID, error) {
	e.ID = entity.NewID()
	pwd, err := s.pwd.Generate(e.Password)
	if err != nil {
		return e.ID, err
	}
	e.Password = pwd
	return s.repo.Create(e)
}

//Get an user
func (s *manager) Get(id entity.ID) (*User, error) {
	return s.repo.Get(id)
}

//Search users
func (s *manager) Search(query string) ([]*User, error) {
	return s.repo.Search(strings.ToLower(query))
}

//List users
func (s *manager) List() ([]*User, error) {
	return s.repo.List()
}

//Delete an user
func (s *manager) Delete(id entity.ID) error {
	_, err := s.Get(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(id)
}

//Update an user
func (s *manager) Update(e *User) error {
	return s.repo.Update(e)
}
