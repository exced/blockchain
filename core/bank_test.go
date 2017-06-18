package core

import "testing"

func TestExistsBank(t *testing.T) {
	cases := []struct {
		bank string
		want bool
	}{
		{"CA", true},
		{"HSBC", true},
		{"LCL", false},
	}
	for _, c := range cases {
		got := ExistsBank(c.bank)
		if got != c.want {
			t.Errorf("ExistsBank(%s) == %t, want %t", c.bank, got, c.want)
		}
	}
}
