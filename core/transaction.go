package core

// Transaction represents a money transaction between 2 users
type Transaction struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount int64  `json:"amount"`
}