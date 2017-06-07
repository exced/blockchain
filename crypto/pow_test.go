package crypto

import (
	"testing"
)

func TestMatchHash(t *testing.T) {
	cases := []struct {
		pow  *PoW
		hash string
		want bool
	}{
		{NewPoW(3, 4), "00023", true},
		{NewPoW(3, 4), "0023", false},
		{NewPoW(3, 4), "023", false},
		{NewPoW(3, 4), "000023", true},
		{NewPoW(3, 4), "11123", false},
	}
	for _, c := range cases {
		got := c.pow.MatchHash(c.hash)
		if got != c.want {
			t.Errorf("(%v).MatchHash(%q) == %t, want %t", c.pow, c.hash, got, c.want)
		}
	}
}
