package rtnetlink_test

import (
	"log"
	"net"

	"github.com/jsimonetti/rtnetlink"
	"golang.org/x/sys/unix"
)

// Set the operational state an interface to Up
func Example_setLinkUp() {
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
	// If the link is already up, return immediately
	if state == rtnetlink.OperStateUp || state == rtnetlink.OperStateUnknown {
		return
	}

	// Set the interface operationally UP
	err = conn.Link.Set(&rtnetlink.LinkMessage{
		Family: msg.Family,
		Type:   msg.Type,
		Index:  uint32(iface.Index),
		Flags:  unix.IFF_UP,
		Change: unix.IFF_UP,
	})

	log.Fatal(err)
}
