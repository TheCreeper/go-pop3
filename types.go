package pop3

import (
	"bufio"
	"net"
)

type Client struct {

	// Net Conn
	conn net.Conn

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
