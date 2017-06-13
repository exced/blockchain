package consensus

import "github.com/gorilla/websocket"

// Connect connects to peer address and await for its consensus and network response.
func Connect(url string) (*Consensus, *Network, error) {
	var dialer *websocket.Dialer
	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		return nil, nil, err
	}
	c := &Consensus{}
	err = conn.ReadJSON(c)
	if err != nil {
		return nil, nil, err
	}
	n := &Network{}
	err = conn.ReadJSON(n)
	if err != nil {
		return nil, nil, err
	}
	return c, n, nil
}
