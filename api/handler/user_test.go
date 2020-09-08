package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/codegangsta/negroni"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/samuelssantos/user-service/api/presenter"
	"github.com/samuelssantos/user-service/domain"
	"github.com/samuelssantos/user-service/domain/entity"
	"github.com/samuelssantos/user-service/domain/entity/user"
	"github.com/samuelssantos/user-service/domain/entity/user/mock"
	"github.com/stretchr/testify/assert"
)

func Test_listUsers(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	m := mock.NewMockManager(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakeUserHandlers(r, *n, m)
	path, err := r.GetRoute("listUsers").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/user", path)
	u := user.NewFixtureUser()
	m.EXPECT().
		List().
		Return([]*user.User{u}, nil)
	ts := httptest.NewServer(listUsers(m))
	defer ts.Close()
	res, err := http.Get(ts.URL)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func Test_listUsers_NotFound(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	m := mock.NewMockManager(controller)
	ts := httptest.NewServer(listUsers(m))
	defer ts.Close()
	m.EXPECT().
		Search("zeca").
		Return(nil, domain.ErrNotFound)
	res, err := http.Get(ts.URL + "?name=zeca")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func Test_listUsers_Search(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	m := mock.NewMockManager(controller)
	u := user.NewFixtureUser()
	m.EXPECT().
		Search("zema").
		Return([]*user.User{u}, nil)
	ts := httptest.NewServer(listUsers(m))
	defer ts.Close()
	res, err := http.Get(ts.URL + "?name=zema")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func Test_createUser(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	m := mock.NewMockManager(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakeUserHandlers(r, *n, m)
	path, err := r.GetRoute("createUser").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/user", path)

	m.EXPECT().
		Create(gomock.Any()).
		Return(entity.NewID(), nil)
	h := createUser(m)

	ts := httptest.NewServer(h)
	defer ts.Close()
	payload := fmt.Sprintf(`{
		"password" : "challenge",
		"email" : "zelda@hash.com.br",
		"first_name" : "zelda",
		"last_name" : "zica",
		"date_of_birth" : "2006-01-02"
	}`)
	resp, _ := http.Post(ts.URL+"/v1/user", "application/json", strings.NewReader(payload))
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var u *presenter.User
	json.NewDecoder(resp.Body).Decode(&u)
	assert.Equal(t, "zelda zica", fmt.Sprintf("%s %s", u.FirstName, u.LastName))
}

func Test_getUser(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	m := mock.NewMockManager(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakeUserHandlers(r, *n, m)
	path, err := r.GetRoute("getUser").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/user/{id}", path)
	u := user.NewFixtureUser()
	m.EXPECT().
		Get(u.ID).
		Return(u, nil)
	handler := getUser(m)
	r.Handle("/v1/user/{id}", handler)
	ts := httptest.NewServer(r)
	defer ts.Close()
	res, err := http.Get(ts.URL + "/v1/user/" + u.ID.String())
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	var d *presenter.User
	json.NewDecoder(res.Body).Decode(&d)
	assert.NotNil(t, d)
	assert.Equal(t, u.ID, d.ID)
}

func Test_deleteUser(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	m := mock.NewMockManager(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakeUserHandlers(r, *n, m)
	path, err := r.GetRoute("deleteUser").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/user/{id}", path)
	u := user.NewFixtureUser()
	m.EXPECT().Delete(u.ID).Return(nil)
	handler := deleteUser(m)
	req, _ := http.NewRequest("DELETE", "/v1/user/"+u.ID.String(), nil)
	r.Handle("/v1/user/{id}", handler).Methods("DELETE", "OPTIONS")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}
