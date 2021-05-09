package remote

import "net"

// Client connects to a remote parallel session and works on
// the given process.
type Client struct {
	conn net.Conn
}
