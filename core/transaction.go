package core

// Transaction represents a money transaction between 2 users
type Transaction struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount int64  `json:"amount"`
}

// NewTransaction returns a new Transaction from someone to someone else
func NewTransaction(from, to string, a int64) *Transaction {
	return &Transaction{from, to, a}
}

// Transactions represent a list of transaction.
// Transactions are received by peers before receiving next valid block.
// These are flushed and appended to the new next block each time you receive a new valid block.
type Transactions struct {
	Transactions []*Transaction
}

// NewTransactions returns an empty list of transaction
func NewTransactions() *Transactions {
	return &Transactions{}
}

// Append add given transaction at the end of received list, if it is valid
func (ts *Transactions) Append(t *Transaction) {
	if ts.IsValid(t) {
		ts.Transactions = append(ts.Transactions, t)
	}
}

func (ts *Transactions) IsValid(other *Transaction) bool {
	balances := make(map[string]int64) // key - amount
	for _, t := range ts.Transactions {

		balances[t.From] -= t.Amount
		balances[t.To] += t.Amount
	}
	// ignore banks
	if !ExistsBank(other.From) {
		if balances[other.From] < other.Amount {
			return false
		}
	}
	return true
}

func (ts *Transactions) AreValid(others *Transactions) bool {
	balances := make(map[string]int64) // key - amount
	for _, t := range ts.Transactions {
		// ignore banks
		balances[t.From] -= t.Amount
		balances[t.To] += t.Amount
	}
	for _, t := range others.Transactions {
		// ignore banks
		balances[t.From] -= t.Amount
		if !ExistsBank(t.From) {
			if balances[t.From] < 0 {
				return false
			}
		}
		balances[t.To] += t.Amount
	}
	return true
}
