package client

// Client represents a client by its key and password
type Client struct {
	publicKey string
	password  string
}

// NewClient creates a new client
func NewClient(password string) *Client {
	return &Client{
		publicKey: "",
		password:  password,
	}
}
