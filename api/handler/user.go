package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/samuelssantos/user-service/domain/entity/user"

	"github.com/samuelssantos/user-service/domain"

	"github.com/samuelssantos/user-service/api/presenter"

	"github.com/samuelssantos/user-service/domain/entity"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

const (
	dateFormat = "2006-01-02"
)

func dateToString(date time.Time) string {

	if date.IsZero() {
		return ""
	}

	return date.Format(dateFormat)
}

func stringToDate(date string) time.Time {

	dateOf, _ := time.Parse(dateFormat, date)

	return dateOf
}

func listUsers(manager user.Manager) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error reading users"
		var data []*user.User
		var err error
		name := r.URL.Query().Get("name")
		switch {
		case name == "":
			data, err = manager.List()
		default:
			data, err = manager.Search(name)
		}
		w.Header().Set("Content-Type", "application/json")
		if err != nil && err != domain.ErrNotFound {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		if data == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(errorMessage))
			return
		}
		var toJ []*presenter.User
		for _, d := range data {

			toJ = append(toJ, &presenter.User{
				ID:          d.ID,
				Email:       d.Email,
				FirstName:   d.FirstName,
				LastName:    d.LastName,
				DateOfBirth: dateToString(d.DateOfBirth),
			})
		}
		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
		}
	})
}

func createUser(manager user.Manager) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error adding user"
		var input struct {
			Email       string `json:"email"`
			Password    string `json:"password"`
			FirstName   string `json:"first_name"`
			LastName    string `json:"last_name"`
			DateOfBirth string `json:"date_of_birth"`
		}
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		u := &user.User{
			ID:          entity.NewID(),
			Email:       input.Email,
			Password:    input.Password,
			FirstName:   input.FirstName,
			LastName:    input.LastName,
			DateOfBirth: stringToDate(input.DateOfBirth),
		}
		u.ID, err = manager.Create(u)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
		toJ := &presenter.User{
			ID:          u.ID,
			Email:       u.Email,
			FirstName:   u.FirstName,
			LastName:    u.LastName,
			DateOfBirth: dateToString(u.DateOfBirth),
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
	})
}

func getUser(manager user.Manager) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error reading user"
		vars := mux.Vars(r)
		id, err := entity.StringToID(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
		data, err := manager.Get(id)
		w.Header().Set("Content-Type", "application/json")
		if err != nil && err != domain.ErrNotFound {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		if data == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(errorMessage))
			return
		}

		date := data.DateOfBirth

		toJ := &presenter.User{
			ID:          data.ID,
			Email:       data.Email,
			FirstName:   data.FirstName,
			LastName:    data.LastName,
			DateOfBirth: dateToString(date),
		}
		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
		}
	})
}

func deleteUser(manager user.Manager) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error removing user"
		vars := mux.Vars(r)
		id, err := entity.StringToID(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
		err = manager.Delete(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
	})
}

//MakeUserHandlers make url handlers
func MakeUserHandlers(r *mux.Router, n negroni.Negroni, manager user.Manager) {
	r.Handle("/v1/user", n.With(
		negroni.Wrap(listUsers(manager)),
	)).Methods("GET", "OPTIONS").Name("listUsers")

	r.Handle("/v1/user", n.With(
		negroni.Wrap(createUser(manager)),
	)).Methods("POST", "OPTIONS").Name("createUser")

	r.Handle("/v1/user/{id}", n.With(
		negroni.Wrap(getUser(manager)),
	)).Methods("GET", "OPTIONS").Name("getUser")

	r.Handle("/v1/user/{id}", n.With(
		negroni.Wrap(deleteUser(manager)),
	)).Methods("DELETE", "OPTIONS").Name("deleteUser")
}
