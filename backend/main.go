package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"gopkg.in/mgo.v2"

	"github.com/exced/simple-blockchain/backend/model"
)

var (
	Database *mgo.Database
)

func WithDB(h http.Handler, storage api.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		h.ServeHTTP(w, r)

	})
}

func main() {

	session, err = mgo.Dial("mongodb://localhost/simple-blockchain")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// user storage
	userStorage := model.NewMgoUserStorage(session)

	r := mux.NewRouter()
	r.HandleFunc("/user", GetPeopleEndpoint).Methods("GET")
	r.HandleFunc("/user/{id}", GetPersonEndpoint).Methods("GET")
	r.HandleFunc("/user/{id}", CreatePersonEndpoint).Methods("POST")
	r.HandleFunc("/user/{id}", DeletePersonEndpoint).Methods("DELETE")
}
