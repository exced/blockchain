package core

// Transactionchain represents a valid chain of transactions
type Transactionchain []*Transaction

// IsValid checks if given transaction is valid with receiver transaction chain
func (tc *Transactionchain) IsValid(t *Transaction) bool {
	return false
}
