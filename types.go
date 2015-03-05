package pop3

import (
	"bufio"
	"io"
	"net"
)

type Client struct {

	// Pluggable Dialer
	Dial func(network, addr string) (net.Conn, error)

	// Net Conn
	conn io.ReadWriteCloser

	// Read Buffer
	r *bufio.Reader

	// Write buffer
	w *bufio.Writer
}

type MessageList struct {

	// Non unique id reported by the server
	ID int

	// Size of the message
	Size int
}

type MessageUidl struct {

	// Non unique id reported by the server
	ID int

	// Unique id reported by the server
	UID string
}
