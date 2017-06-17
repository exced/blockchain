package core

import "testing"

func TestTransactionIsValid(t *testing.T) {
	transactions := NewTransactions()
	transactions.Append(&Transaction{From: "0", To: "1", Amount: 5})
	transactions.Append(&Transaction{From: "0", To: "1", Amount: 5})
	transactions.Append(&Transaction{From: "2", To: "3", Amount: 5})
	transactions.Append(&Transaction{From: "3", To: "1", Amount: 5})
	cases := []struct {
		transactions *Transactions
		transaction  *Transaction
		want         bool
	}{
		{transactions, &Transaction{From: "1", To: "0", Amount: 5}, true},
		{transactions, &Transaction{From: "0", To: "1", Amount: 5}, false},
		{transactions, &Transaction{From: "1", To: "3", Amount: 15}, true},
		{transactions, &Transaction{From: "1", To: "3", Amount: 20}, false},
		{transactions, &Transaction{From: "HSBC", To: "4", Amount: 100}, false},
	}
	for _, c := range cases {
		got := c.transactions.IsValid(c.transaction)
		if got != c.want {
			t.Errorf("(%q).IsValid(%#v) == %t, want %t", c.transactions, c.transaction, got, c.want)
		}
	}
}

func TestTransactionAreValid(t *testing.T) {
	transactions := NewTransactions()
	transactions.Append(&Transaction{From: "0", To: "1", Amount: 5})
	transactions.Append(&Transaction{From: "0", To: "1", Amount: 5})
	transactions.Append(&Transaction{From: "2", To: "3", Amount: 5})
	transactions.Append(&Transaction{From: "3", To: "1", Amount: 5})
	transactions1 := NewTransactions()
	transactions1.Append(&Transaction{From: "1", To: "0", Amount: 5})
	transactions1.Append(&Transaction{From: "1", To: "0", Amount: 5})
	transactions1.Append(&Transaction{From: "1", To: "0", Amount: 5})
	transactions1.Append(&Transaction{From: "1", To: "0", Amount: 10})
	transactions2 := NewTransactions()
	transactions2.Append(&Transaction{From: "1", To: "0", Amount: 2})
	transactions2.Append(&Transaction{From: "1", To: "0", Amount: 2})
	cases := []struct {
		received *Transactions
		given    *Transactions
		want     bool
	}{
		{transactions, transactions1, false},
		{transactions, transactions2, true},
	}
	for _, c := range cases {
		got := c.received.AreValid(c.given)
		if got != c.want {
			t.Errorf("(%q).AreValid(%#v) == %t, want %t", c.received, c.given, got, c.want)
		}
	}
}
