package pop3

// Non optional POP3 commands as extracted from rfc1939 section 5 and 6.
const (
	USER = "USER"
	PASS = "PASS"
	QUIT = "QUIT"
	STAT = "STAT"
	LIST = "LIST"
	RETR = "RETR"
	DELE = "DELE"
	NOOP = "NOOP"
	RSET = "RSET"
)

// Optional POP3 commands as extracted from rfc1939 section 7.
const (
	ATOP = "ATOP"
	TOP  = "TOP"
	UIDL = "UIDL"
)

// POP3 replies as extracted from rfc1939 section 9.
const (
	OK  = "+OK"
	ERR = "-ERR"
)
