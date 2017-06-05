package core

import (
	"testing"
)

func TestIsValid(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"", true},
		{"abcd", true},
		{"aa", false},
		{"aba", false},
	}
	for _, c := range cases {
		got := isValid(c.in)
		if got != c.want {
			t.Errorf("isUnique(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}