package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"gopkg.in/mgo.v2"

	"github.com/exced/blockchain/backend/api"
)

func main() {
	httpAddr := flag.String("http", ":3000", "HTTP listen address")
	dbAddr := flag.String("db", "mongodb://localhost/blockchain", "DB listen address")
	flag.Parse()

	// MongoDB Dial
	session, err := mgo.Dial(*dbAddr)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// user storage and user API
	userAPI := api.NewUserAPI(session)

	// routes
	r := mux.NewRouter()
	r.HandleFunc("/login", userAPI.LoginUser).Methods("POST")
	r.HandleFunc("/signin", userAPI.SigninUser).Methods("POST")
	r.HandleFunc("/user/{id}", userAPI.GetUser).Methods("GET")
	r.HandleFunc("/user/{id}", userAPI.PostUser).Methods("POST")
	r.HandleFunc("/user/{id}", userAPI.PostUser).Methods("PUT")
	r.HandleFunc("/user/{id}", userAPI.DeleteUser).Methods("DELETE")

	// http serve
	log.Println("http server started on", *httpAddr)
	err = http.ListenAndServe(*httpAddr, r)
	if err != nil {
		log.Fatal("Could not serve http: ", err)
	}
}
