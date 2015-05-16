package pop3

import (
	"bufio"
	"net"
)

// Client holds the net conn and read/write buffer objects.
type Client struct {

	// Net Conn
	conn net.Conn

	// Read Buffer
	r *bufio.Reader

	// Write buffer
	w *bufio.Writer
}

// MessageList represents the metadata returned by the server for a message stored in the maildrop.
type MessageList struct {

	// Non unique id reported by the server
	ID int

	// Size of the message
	Size int
}

// MessageUidl represents the metadata returned by the server for a message stored in the maildrop.
type MessageUidl struct {

	// Non unique id reported by the server
	ID int

	// Unique id reported by the server
	UID string
}
