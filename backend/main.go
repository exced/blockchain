package main

import (
	"github.com/gorilla/mux"

	"gopkg.in/mgo.v2"

	"github.com/exced/blockchain/backend/api"
)

func main() {

	session, err := mgo.Dial("mongodb://localhost/simple-blockchain")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// user storage and user API
	userAPI := api.userAPI(session)

	// routes
	r := mux.NewRouter()
	r.HandleFunc("/login", userAPI.LoginUser).Methods("POST")
	r.HandleFunc("/signin", userAPI.SigninUser).Methods("POST")
	r.HandleFunc("/user/{id}", userAPI.GetUser).Methods("GET")
	r.HandleFunc("/user/{id}", userAPI.PostUser).Methods("POST")
	r.HandleFunc("/user/{id}", userAPI.PostUser).Methods("PUT")
	r.HandleFunc("/user/{id}", userAPI.DeleteUser).Methods("DELETE")
}
