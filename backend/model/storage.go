package model

// Storage defines a CRUD API to access our resources.
type Storage interface {
	Create(i *interface{}) error
	Read(i *interface{}) (*interface{}, error)
	Update(o *interface{}, n *interface{}) error
	Delete(i *interface{}) error
}
