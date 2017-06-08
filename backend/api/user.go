package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/exced/simple-blockchain/backend/model"
	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
)

type UserAPI struct {
	Storage *model.MgoUserStorage
}

func NewUserAPI(s *mgo.Session) *UserAPI {
	return &UserAPI{model.NewMgoUserStorage(s)}
}

func (api *UserAPI) LoginUser(w http.ResponseWriter, r *http.Request) {

}

func (api *UserAPI) SigninUser(w http.ResponseWriter, r *http.Request) {
	var u User
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

func (api *UserAPI) GetUser(w http.ResponseWriter, r *http.Request) {

}

func (api *UserAPI) PostUser(w http.ResponseWriter, r *http.Request) {

}

func (api *UserAPI) PutUser(w http.ResponseWriter, r *http.Request) {

}

func (api *UserAPI) DeleteUser(w http.ResponseWriter, r *http.Request) {

}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(members)
	w.Write(j)
}

func postMembersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var m Member
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &m)

	members = append(members, m)

	j, _ := json.Marshal(m)
	w.Write(j)
}

func GetUser(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(people)
}

func CreateUser(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	var user User
	_ = json.NewDecoder(req.Body).Decode(&person)
	person.ID = params["id"]
	people = append(people, person)
	json.NewEncoder(w).Encode(people)
}

func DeletePersonEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for index, item := range people {
		if item.ID == params["id"] {
			people = append(people[:index], people[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(people)
}
