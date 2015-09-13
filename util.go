package pop3

import "strings"

// IsOk checks to see if the reply from the server contains +OK.
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
