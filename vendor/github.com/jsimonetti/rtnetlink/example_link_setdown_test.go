package rtnetlink_test

import (
	"log"
	"net"

	"github.com/jsimonetti/rtnetlink"
)

// Set the operational state an interface to Down
func Example_setLinkDown() {
	// Gather the interface Index
	iface, _ := net.InterfaceByName("dummy0")

	// Dial a connection to the rtnetlink socket
	conn, err := rtnetlink.Dial(nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Request the details of the interface
	msg, err := conn.Link.Get(uint32(iface.Index))
	if err != nil {
		log.Fatal(err)
	}

	state := msg.Attributes.OperationalState
	// If the link is already down, return immediately
	if state == rtnetlink.OperStateDown {
		return
	}

	// Set the interface operationally Down
	err = conn.Link.Set(&rtnetlink.LinkMessage{
		Family: msg.Family,
		Type:   msg.Type,
		Index:  uint32(iface.Index),
		Flags:  0x0,
		Change: 0x1,
	})

	log.Fatal(err)
}
