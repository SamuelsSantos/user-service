package user

import (
	"database/sql"
	"log"

	"github.com/samuelssantos/user-service/domain"

	"github.com/samuelssantos/user-service/domain/entity"
)

//SQLRepo repository
type SQLRepo struct {
	db *sql.DB
}

//NewSQLRepoRepository create new repository
func NewSQLRepoRepository(db *sql.DB) *SQLRepo {
	return &SQLRepo{
		db: db,
	}
}

//Create an user
func (r *SQLRepo) Create(e *User) (entity.ID, error) {

	var stmt *sql.Stmt
	var err error

	if e.DateOfBirth.IsZero() {
		stmt, err = r.db.Prepare(`insert into public.user (id, "password", email, first_name, last_name) values ($1, $2, $3, $4, $5)`)
	} else {
		stmt, err = r.db.Prepare(`insert into public.user (id, "password", email, first_name, last_name, date_of_birth) values ($1, $2, $3, $4, $5, $6)`)
	}

	if err != nil {
		return e.ID, err
	}

	if e.DateOfBirth.IsZero() {
		_, err = stmt.Exec(
			e.ID,
			e.Password,
			e.Email,
			e.FirstName,
			e.LastName,
		)

	} else {
		_, err = stmt.Exec(
			e.ID,
			e.Password,
			e.Email,
			e.FirstName,
			e.LastName,
			e.DateOfBirth,
		)
	}

	if err != nil {
		log.Fatal(err)
		return e.ID, err
	}
	err = stmt.Close()
	if err != nil {
		return e.ID, err
	}
	return e.ID, nil

}

//Get an user
func (r *SQLRepo) Get(id entity.ID) (*User, error) {
	return getUser(id, r.db)
}

func getUser(id entity.ID, db *sql.DB) (*User, error) {
	stmt, err := db.Prepare(`select id, email, first_name, last_name, date_of_birth from public.user where id = $1`)
	if err != nil {
		return nil, err
	}
	var u User
	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.DateOfBirth)
	}
	return &u, nil
}

//Update an user
func (r *SQLRepo) Update(e *User) error {
	if e.DateOfBirth.IsZero() {
		_, err := r.db.Exec("update public.user set email = $1, password = $2, first_name = $3, last_name = $4 where id = $5",
			e.Email,
			e.Password,
			e.FirstName,
			e.LastName,
			e.ID,
		)
		if err != nil {
			return err
		}
	} else {
		_, err := r.db.Exec("update public.user set email = $1, password = $2, first_name = $3, last_name = $4, date_of_birth = $5 where id = $6",
			e.Email,
			e.Password,
			e.FirstName,
			e.LastName,
			e.DateOfBirth,
			e.ID,
		)
		if err != nil {
			return err
		}

	}

	return nil
}

//Search users
func (r *SQLRepo) Search(query string) ([]*User, error) {
	stmt, err := r.db.Prepare(`select id from public.user where name like $1`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var ids []entity.ID
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var i entity.ID
		err = rows.Scan(&i)
		if err != nil {
			return nil, err
		}
		ids = append(ids, i)
	}
	if len(ids) == 0 {
		return nil, domain.ErrNotFound
	}
	var users []*User
	for _, id := range ids {
		u, err := getUser(id, r.db)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

//List users
func (r *SQLRepo) List() ([]*User, error) {
	stmt, err := r.db.Prepare(`select id from public.user`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var ids []entity.ID
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var i entity.ID
		err = rows.Scan(&i)
		if err != nil {
			return nil, err
		}
		ids = append(ids, i)
	}
	if len(ids) == 0 {
		return nil, domain.ErrNotFound
	}
	var users []*User
	for _, id := range ids {
		u, err := getUser(id, r.db)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

//Delete an user
func (r *SQLRepo) Delete(id entity.ID) error {
	_, err := r.db.Exec("delete from public.user where id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
