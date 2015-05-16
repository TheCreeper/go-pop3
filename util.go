package pop3

import "strings"

func IsOK(s string) bool {

	if strings.Fields(s)[0] != OK {

		return false
	}

	return true
}

func IsErr(s string) bool {

	if strings.Fields(s)[0] != ERR {

		return false
	}

	return true
}
