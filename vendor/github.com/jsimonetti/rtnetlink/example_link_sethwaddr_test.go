package rtnetlink_test

import (
	"log"
	"net"

	"github.com/jsimonetti/rtnetlink"
)

// Set the hw address of an interface
func Example_setLinkHWAddr() {
	// Gather the interface Index
	iface, _ := net.InterfaceByName("dummy0")
	// Get a hw addr to set the interface to
	hwAddr, _ := net.ParseMAC("ce:9c:5b:98:55:9c")

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

	// Set the hw address of the interfaces
	err = conn.Link.Set(&rtnetlink.LinkMessage{
		Family: msg.Family,
		Type:   msg.Type,
		Index:  uint32(iface.Index),
		Flags:  msg.Flags,
		Change: msg.Change,
		Attributes: &rtnetlink.LinkAttributes{
			Address: hwAddr,
		},
	})

	log.Fatal(err)
}
