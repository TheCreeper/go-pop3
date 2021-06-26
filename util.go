package pop3

import (
	"errors"
	"strings"
)

// IsOK checks to see if the reply from the server contains +OK.
func IsOK(s string) bool {
	if strings.Fields(s)[0] != OK {
		return false
	}
	return true
}

// IsErr checks to see if the reply from the server contains +Err.
func IsErr(s string) bool {
	if strings.Fields(s)[0] != ERR {
		return false
	}
	return true
}

// GetErr checks the reply from the server and returns an error
// if it's not an OK response.
func GetErr(s string) error {
	f := strings.Fields(s)
	if f[0] != ERR {
		return nil
	}

	return errors.New(strings.Join(f[1:], " "))
}
