package core

import (
	"testing"
)

func TestBlockIsValid(t *testing.T) {
	transactions := NewTransactions()
	ts := make([]*Transaction, 4)
	ts[0] = &Transaction{From: "0", To: "1", Amount: 5}
	ts[1] = &Transaction{From: "0", To: "1", Amount: 5}
	ts[2] = &Transaction{From: "2", To: "3", Amount: 5}
	ts[3] = &Transaction{From: "3", To: "1", Amount: 5}
	transactions.Transactions = ts
	g := genesisBlock
	a := g.GenNext(transactions)
	b := a.GenNext(transactions)
	cases := []struct {
		block    *Block
		previous *Block
		want     bool
	}{
		{a, g, true},
		{b, a, true},
		{b, g, false},
	}
	for _, c := range cases {
		got := c.block.IsValid(c.previous)
		if got != c.want {
			t.Errorf("(%#v).IsValid(%#v) == %t, want %t", c.block, c.previous, got, c.want)
		}
	}
}
