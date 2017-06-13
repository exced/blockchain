package core

import (
	"testing"
)

func TestIsTransactionValid(t *testing.T) {
	b := genesisBlock
	b.Transactions = []*Transaction{
		&Transaction{From: "0", To: "1", Amount: 5},
		&Transaction{From: "0", To: "1", Amount: 5},
		&Transaction{From: "2", To: "3", Amount: 5},
		&Transaction{From: "3", To: "1", Amount: 5},
	}
	cases := []struct {
		block       *Block
		transaction *Transaction
		want        bool
	}{
		{b, &Transaction{From: "1", To: "0", Amount: 5}, true},
		{b, &Transaction{From: "0", To: "1", Amount: 5}, false},
		{b, &Transaction{From: "1", To: "3", Amount: 15}, true},
		{b, &Transaction{From: "1", To: "3", Amount: 20}, false},
	}
	for _, c := range cases {
		got := c.block.IsTransactionValid(c.transaction)
		if got != c.want {
			t.Errorf("(%q).IsTransactionValid(%#v) == %t, want %t", c.block, c.transaction, got, c.want)
		}
	}
}

func TestAreTransactionsValid(t *testing.T) {
	b := genesisBlock
	b.Transactions = []*Transaction{
		&Transaction{From: "0", To: "1", Amount: 5},
		&Transaction{From: "0", To: "1", Amount: 5},
		&Transaction{From: "2", To: "3", Amount: 5},
		&Transaction{From: "3", To: "1", Amount: 5},
	}
	transactions1 := []*Transaction{
		&Transaction{From: "1", To: "0", Amount: 5},
		&Transaction{From: "1", To: "0", Amount: 5},
		&Transaction{From: "1", To: "0", Amount: 5},
		&Transaction{From: "1", To: "0", Amount: 10},
	}
	transactions2 := []*Transaction{
		&Transaction{From: "1", To: "0", Amount: 2},
		&Transaction{From: "1", To: "0", Amount: 2},
	}
	cases := []struct {
		block        *Block
		transactions []*Transaction
		want         bool
	}{
		{b, transactions1, false},
		{b, transactions2, true},
	}
	for _, c := range cases {
		got := c.block.areTransactionsValid(c.transactions)
		if got != c.want {
			t.Errorf("(%q).areTransactionsValid(%#v) == %t, want %t", c.block, c.transactions, got, c.want)
		}
	}
}
