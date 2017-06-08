package api

import (
	"encoding/json"
	"net/http"

	"github.com/exced/blockchain/backend/model"
	mgo "gopkg.in/mgo.v2"
)

// UserAPI gives method to authenticate and access user model.
type UserAPI struct {
	storage *model.MgoUserStorage
}

// NewUserAPI returns a new UserAPI
func NewUserAPI(s *mgo.Session) *UserAPI {
	return &UserAPI{model.NewMgoUserStorage(s)}
}

func (api *UserAPI) LoginUser(w http.ResponseWriter, r *http.Request) {
	var u *model.User
	if r.Body == nil {
		respondWithError(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		respondWithError(w, err.Error(), 400)
		return
	}
	u, err = api.storage.Auth(u)
	if err != nil {
		respondWithError(w, err.Error(), 400)
		return
	}
	respondWithJSON(w, u, 200)
}

func (api *UserAPI) SigninUser(w http.ResponseWriter, r *http.Request) {
	var u *model.User
	if r.Body == nil {
		respondWithError(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		respondWithError(w, err.Error(), 400)
		return
	}
	err = api.storage.Create(u)
	if err != nil {
		respondWithError(w, err.Error(), 400)
		return
	}
	respondWithJSON(w, u, 200)
}

func (api *UserAPI) GetUser(w http.ResponseWriter, r *http.Request) {
	var u *model.User
	if r.Body == nil {
		respondWithError(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		respondWithError(w, err.Error(), 400)
		return
	}
	u, err = api.storage.Read(u)
	if err != nil {
		respondWithError(w, err.Error(), 400)
		return
	}
	respondWithJSON(w, u, 200)
}

func (api *UserAPI) PostUser(w http.ResponseWriter, r *http.Request) {

}

func (api *UserAPI) PutUser(w http.ResponseWriter, r *http.Request) {

}

func (api *UserAPI) DeleteUser(w http.ResponseWriter, r *http.Request) {

}
