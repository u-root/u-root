package rtnetlink_test

import (
	"log"

	"github.com/jsimonetti/rtnetlink"
)

// List all interfaces
func Example_listLink() {
	// Dial a connection to the rtnetlink socket
	conn, err := rtnetlink.Dial(nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Request a list of interfaces
	msg, err := conn.Link.List()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%#v", msg)
}
