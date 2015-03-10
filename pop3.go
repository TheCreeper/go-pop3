package pop3

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/mail"
	"strconv"
	"strings"
)

// Some errors
var (
	ErrServerNotReady = errors.New("Server did not respond with +OK on the initial connection")
	ErrBadCommand     = errors.New("Server did not respond with +OK after sending a command")
)

type Error struct {

	// The reply sent by the server
	Response string

	// Error message
	Err error
}

func (e *Error) Error() string {

	return fmt.Sprintf("pop3: %s: %s\n", e.Err.Error(), e.Response)
}

func Dial(address string) (c *Client, err error) {

	conn, err := net.Dial("tcp", address)
	if err != nil {

		return
	}

	return NewClient(conn)
}

func DialTLS(address string) (c *Client, err error) {

	conn, err := tls.Dial("tcp", address, nil)
	if err != nil {

		return
	}

	return NewClient(conn)
}

func NewClient(conn io.ReadWriteCloser) (c *Client, err error) {

	c = &Client{

		conn: conn,
		r:    bufio.NewReader(conn),
		w:    bufio.NewWriter(conn),
	}

	// Make sure we receive the server greeting
	line, err := c.ReadLine()
	if err != nil {

		return
	}
	if !IsOK(line) {

		return nil, &Error{string(line), ErrServerNotReady}
	}

	return
}

func (c *Client) ReadLine() (line string, err error) {

	b, _, err := c.r.ReadLine()
	if err == io.EOF {

		return
	}
	if err != nil {

		return
	}
	line = string(b)

	return
}

func (c *Client) ReadLines() (lines []string, err error) {

	for {

		line, err := c.ReadLine()
		if err != nil {

			return nil, err
		}

		// Look for a dot to indicate the end of a message
		// from the server.
		if line == "." {

			break
		}
		lines = append(lines, line)
	}

	return
}

func (c *Client) Send(format string, args ...interface{}) (err error) {

	_, err = c.w.WriteString(fmt.Sprintf(format, args...))
	if err != nil {

		return
	}

	err = c.w.Flush()
	if err != nil {

		return
	}

	return
}

func (c *Client) Cmd(format string, args ...interface{}) (line string, err error) {

	err = c.Send(format, args...)
	if err != nil {

		return
	}

	line, err = c.ReadLine()
	if err != nil {

		return
	}
	if !IsOK(line) {

		return "", &Error{string(line), ErrBadCommand}
	}

	return
}

func (c *Client) User(u string) (err error) {

	_, err = c.Cmd("%s %s\r\n", USER, u)
	if err != nil {

		return
	}

	return
}

func (c *Client) Pass(p string) (err error) {

	_, err = c.Cmd("%s %s\r\n", PASS, p)
	if err != nil {

		return
	}

	return
}

func (c *Client) Quit() (err error) {

	// Send the server a QUIT message
	err = c.Send("%s\r\n", QUIT)
	if err != nil {

		return
	}

	// Close the socket
	return c.conn.Close()
}

func (c *Client) Auth(u, p string) (err error) {

	err = c.User(u)
	if err != nil {

		return
	}

	err = c.Pass(p)
	if err != nil {

		return
	}

	// issue a dud command. Server might
	// not respond with invalid auth unless a cmd is issued
	return c.Noop()
}

func (c *Client) Stat() (count, size int, err error) {

	line, err := c.Cmd("%s\r\n", STAT)
	if err != nil {

		return
	}

	// Number of messages in maildrop
	count, err = strconv.Atoi(strings.Fields(line)[1])
	if err != nil {

		return
	}
	if count < 1 {

		return
	}

	// Total size of messages in bytes
	size, err = strconv.Atoi(strings.Fields(line)[2])
	if err != nil {

		return
	}
	if size < 1 {

		return
	}

	return
}

func (c *Client) List(msg int) (list MessageList, err error) {

	line, err := c.Cmd("%s %s\r\n", LIST, msg)
	if err != nil {

		return
	}

	id, err := strconv.Atoi(strings.Fields(line)[0])
	if err != nil {

		return
	}

	size, err := strconv.Atoi(strings.Fields(line)[1])
	if err != nil {

		return
	}

	return MessageList{id, size}, nil
}

func (c *Client) ListAll() (list []MessageList, err error) {

	_, err = c.Cmd("%s\r\n", LIST)
	if err != nil {

		return
	}

	lines, err := c.ReadLines()
	if err != nil {

		return
	}
	for _, v := range lines {

		id, err := strconv.Atoi(strings.Fields(v)[0])
		if err != nil {

			return nil, err
		}

		size, err := strconv.Atoi(strings.Fields(v)[1])
		if err != nil {

			return nil, err
		}

		list = append(list, MessageList{id, size})
	}

	return
}

func (c *Client) Retr(msg int) (m *mail.Message, err error) {

	_, err = c.Cmd("%s %s\r\n", RETR, msg)
	if err != nil {

		return
	}

	m, err = mail.ReadMessage(c.r)
	if err != nil {

		return
	}

	// mail.ReadMessage does not consume the message end dot in the buffer
	// so we must move the buffer along. Need to find a better way of doing this.
	line, err := c.ReadLine()
	if err != nil {

		return
	}
	if line != "." {

		err = c.r.UnreadByte()
		if err != nil {

			return
		}
	}

	return
}

func (c *Client) Dele(msg int) (err error) {

	_, err = c.Cmd("%s %s\r\n", DELE, msg)
	if err != nil {

		return
	}

	return
}

func (c *Client) Noop() (err error) {

	_, err = c.Cmd("%s\r\n", NOOP)
	if err != nil {

		return
	}

	return
}

func (c *Client) Rset() (err error) {

	_, err = c.Cmd("%s\r\n", RSET)
	if err != nil {

		return
	}

	return
}

func (c *Client) Top(msg int, n int) (m *mail.Message, err error) {

	_, err = c.Cmd("%s %d %d\r\n", TOP, msg, n)
	if err != nil {

		return
	}

	m, err = mail.ReadMessage(c.r)
	if err != nil {

		return
	}

	// mail.ReadMessage does not consume the message end dot in the buffer
	// so we must move the buffer along. Need to find a better way of doing this.
	line, err := c.ReadLine()
	if err != nil {

		return
	}
	if line != "." {

		err = c.r.UnreadByte()
		if err != nil {

			return
		}
	}

	return
}

func (c *Client) Uidl(msg int) (list MessageUidl, err error) {

	line, err := c.Cmd("%s %s\r\n", UIDL, msg)
	if err != nil {

		return
	}

	id, err := strconv.Atoi(strings.Fields(line)[1])
	if err != nil {

		return
	}

	return MessageUidl{id, strings.Fields(line)[2]}, nil
}

func (c *Client) UidlAll() (list []MessageUidl, err error) {

	_, err = c.Cmd("%s\r\n", UIDL)
	if err != nil {

		return
	}

	lines, err := c.ReadLines()
	if err != nil {

		return
	}
	for _, v := range lines {

		id, err := strconv.Atoi(strings.Fields(v)[0])
		if err != nil {

			return nil, err
		}

		list = append(list, MessageUidl{id, strings.Fields(v)[1]})
	}

	return
}
