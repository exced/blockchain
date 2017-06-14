package core

import (
	"reflect"
	"testing"
)

func TestAppendBlock(t *testing.T) {
	transactions := []*Transaction{
		&Transaction{From: "0", To: "1", Amount: 5},
	}
	// case
	bc1 := NewBlockchain()
	bc2 := NewBlockchain()
	b := genesisBlock.GenNext(transactions)
	*bc2 = append(*bc2, b)
	// case
	bc3 := NewBlockchain()
	bc4 := NewBlockchain()
	b2 := bc3.GenNext(transactions)
	*bc4 = append(*bc4, b2)
	cases := []struct {
		blockchain *Blockchain
		block      *Block
		want       *Blockchain
	}{
		{bc1, b, bc2},
		{bc3, b2, bc4},
	}
	for _, c := range cases {
		c.blockchain.AppendBlock(c.block)
		if !reflect.DeepEqual(c.blockchain, c.want) {
			t.Errorf("(%#v).AppendBlock(%#v), want %#v", c.blockchain, c.block, c.want)
		}
	}
}

func TestIsValid(t *testing.T) {
	transactions := []*Transaction{
		&Transaction{From: "0", To: "1", Amount: 5},
	}
	bc1 := NewBlockchain()
	bc2 := NewBlockchain()
	bc2.AppendBlock(bc2.GenNext(transactions))
	bc3 := NewBlockchain()
	bc3.AppendBlock(bc3.GenNext(transactions))
	bc3.AppendBlock(bc3.GenNext(transactions))
	bc3.AppendBlock(bc3.GenNext(transactions))
	cases := []struct {
		blockchain *Blockchain
		want       bool
	}{
		{bc1, true},
		{bc2, true},
		{bc3, true},
	}
	for _, c := range cases {
		got := c.blockchain.IsValid()
		if got != c.want {
			t.Errorf("(%#v).IsValid() == %t, want %t", c.blockchain, got, c.want)
		}
	}
}

func TestFetch(t *testing.T) {
	bc1 := NewBlockchain()
	bc2 := NewBlockchain()
	transactions := []*Transaction{
		&Transaction{From: "0", To: "1", Amount: 5},
	}
	bc2.AppendBlock(bc2.GenNext(transactions))
	cases := []struct {
		blockchain *Blockchain
		other      *Blockchain
		want       *Blockchain
	}{
		{bc1, bc2, bc2},
	}
	for _, c := range cases {
		got := c.blockchain.Fetch(c.other)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("(%#v).Fetch(%#v) == %#v, want %#v", c.blockchain, c.other, got, c.want)
		}
	}
}
