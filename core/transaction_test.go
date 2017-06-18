package core

import "testing"

func TestTransactionsIsValid(t *testing.T) {
	transactions := NewTransactions()
	transactions.Append(&Transaction{From: "HSBC", To: "0", Amount: 10})
	transactions.Append(&Transaction{From: "0", To: "1", Amount: 5})
	transactions.Append(&Transaction{From: "0", To: "1", Amount: 5})
	cases := []struct {
		transactions *Transactions
		transaction  *Transaction
		want         bool
	}{
		{transactions, &Transaction{From: "1", To: "0", Amount: 5}, true},
		{transactions, &Transaction{From: "0", To: "1", Amount: 5}, false},
		{transactions, &Transaction{From: "1", To: "3", Amount: 10}, true},
		{transactions, &Transaction{From: "1", To: "3", Amount: 20}, false},
		{transactions, &Transaction{From: "HSBC", To: "4", Amount: 100}, true},
	}
	for _, c := range cases {
		got := c.transactions.IsValid(c.transaction)
		if got != c.want {
			t.Errorf("(%q).IsValid(%#v) == %t, want %t", c.transactions, c.transaction, got, c.want)
		}
	}
}

func TestTransactionsAreValid(t *testing.T) {
	transactions := NewTransactions()
	ts := make([]*Transaction, 5)
	ts[0] = &Transaction{From: "HSBC", To: "0", Amount: 10}
	ts[1] = &Transaction{From: "0", To: "1", Amount: 5}
	ts[2] = &Transaction{From: "0", To: "1", Amount: 5}
	ts[3] = &Transaction{From: "2", To: "3", Amount: 5}
	ts[4] = &Transaction{From: "3", To: "1", Amount: 5}
	transactions.Transactions = ts
	transactions1 := NewTransactions()
	ts1 := make([]*Transaction, 4)
	ts1[0] = &Transaction{From: "1", To: "0", Amount: 5}
	ts1[1] = &Transaction{From: "1", To: "0", Amount: 5}
	ts1[2] = &Transaction{From: "1", To: "0", Amount: 5}
	ts1[3] = &Transaction{From: "1", To: "0", Amount: 10}
	transactions1.Transactions = ts1
	transactions2 := NewTransactions()
	ts2 := make([]*Transaction, 2)
	ts2[0] = &Transaction{From: "1", To: "0", Amount: 2}
	ts2[1] = &Transaction{From: "1", To: "0", Amount: 2}
	transactions2.Transactions = ts2
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
