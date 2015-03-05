package pop3

import (
	"net"
	"strings"
)

func LookupMailServer(addr string) (string, error) {

	host, port, err := net.SplitHostPort(addr)
	if err != nil {

		return "", err
	}

	mxs, err := net.LookupMX(host)
	if err != nil {

		return "", err
	}
	if len(mxs) == 0 {

		return "", err
	}
	return net.JoinHostPort(mxs[0].Host, port), nil
}

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
