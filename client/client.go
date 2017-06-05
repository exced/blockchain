package client

// Client represents a client by its key and password
type Client struct {
	privateKey string
	password   string
}

// NewClient creates a new client
func NewClient(password string) *Client {
	return &Client{
		privateKey: "",
		password:   password,
	}
}
