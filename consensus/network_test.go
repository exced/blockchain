package consensus

import "testing"

var peer1 = NewPeer()

func TestBroadcast(t *testing.T) {
	network1 := NewNetwork()
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
