package core

import (
	"testing"
)

func TestFetch(t *testing.T) {
	bc1 := NewBlockchain()
	bc2 := NewBlockchain()
	tr1 := make([]*Transaction, 0)
	bc2.AppendBlock(bc2.GenNext(tr1))
	cases := []struct {
		blockchain *Blockchain
		other      *Blockchain
		want       *Blockchain
	}{
		{bc1, bc2, bc2},
	}
	for _, c := range cases {
		c.blockchain.Fetch(c.other)
		if c.blockchain != c.want {
			t.Errorf("(%#v).Fetch(%#v), want %#v", c.blockchain, c.other, c.want)
		}
	}
}
