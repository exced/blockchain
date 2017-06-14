package crypto

import (
	"testing"
)

func TestMatchHash(t *testing.T) {
	cases := []struct {
		hash       string
		difficulty int
		want       bool
	}{
		{"00023", 3, true},
		{"0023", 3, false},
		{"023", 3, false},
		{"000023", 3, true},
		{"11123", 3, false},
		{"000023", 4, true},
	}
	for _, c := range cases {
		got := MatchHash(c.hash, c.difficulty)
		if got != c.want {
			t.Errorf("MatchHash(%q, %d) == %t, want %t", c.hash, c.difficulty, got, c.want)
		}
	}
}
