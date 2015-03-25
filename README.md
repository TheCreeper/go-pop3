go-pop3
=====================

[![go-pop3](https://godoc.org/github.com/TheCreeper/go-pop3?status.png)](http://godoc.org/github.com/TheCreeper/go-pop3)

Package pop3 provides an implementation of the [Post Office Protocol - Version 3](https://www.ietf.org/rfc/rfc1939.txt).

## Example

```
package main

import (
	"log"

	"github.com/TheCreeper/go-pop3"
)

func main() {

	// Create a connection to the server
	c, err := pop3.DialTLS("pop3.riseup.net:993")
	if err != nil {

		log.Fatal(err)
	}

	// Authenticate with the server
	if err = c.Auth("username", "password"); err != nil {

		log.Fatal(err)
	}

	// Print the UID of all messages in the maildrop
	messages, err := c.UidlAll()
	if err != nil {

		log.Fatal(err)
	}
	for _, v := range messages {

		log.Print(v.UID)
	}

	// Send the QUIT command and close the connection
	if err = c.Quit(); err != nil {

		log.Fatal(err)
	}
}
```