package model

import (
	"golang.org/x/crypto/bcrypt"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// User represents a user of our blockchain.
type User struct {
	Key      string `json:"key"`
	Password []byte
	Wallet   Wallet `json:"wallet"`
}

// MgoUserStorage uses mongoDB to store data.
type MgoUserStorage struct {
	c *mgo.Collection
}

// NewMgoUserStorage creates and retrieves a new DBUserStorage object.
func NewMgoUserStorage(s *mgo.Session) *MgoUserStorage {
	return &MgoUserStorage{s.DB("store").C("user")}
}

// Create CRUD create resource
func (db *MgoUserStorage) Create(user *User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = hash
	return db.c.Insert(user)
}

// Read CRUD read resource
func (db *MgoUserStorage) Read(user *User) (*User, error) {
	res := User{}
	err := db.c.Find(bson.M{"Key": user.Key}).One(&res)
	if err = bcrypt.CompareHashAndPassword(res.Password, []byte(user.Password)); err != nil {
		return nil, err
	}
	return &res, nil
}

// Update CRUD update resource
func (db *MgoUserStorage) Update(o *User, n *User) error {
	return db.c.Update(o, n)
}

// Delete CRUD delete resource
func (db *MgoUserStorage) Delete(user *User) error {
	return db.c.Remove(user)
}
