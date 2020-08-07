package rtnetlink_test

import (
	"log"
	"net"

	"github.com/jsimonetti/rtnetlink"
	"golang.org/x/sys/unix"
)

// List all IPv4 addresses configured on interface 'lo'
func Example_listAddress() {
	// Gather the interface Index
	iface, _ := net.InterfaceByName("lo")
	// Get an ip address to add to the interface
	family := uint8(unix.AF_INET)

	// Dial a connection to the rtnetlink socket
	conn, err := rtnetlink.Dial(nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Request a list of addresses
	msg, err := conn.Address.List()
	if err != nil {
		log.Fatal(err)
	}

	// Filter out the wanted messages and put them in the 'addr' slice.
	var addr []rtnetlink.AddressMessage
	for _, v := range msg {
		add := true
		if iface != nil && v.Index != uint32(iface.Index) {
			add = false
		}
		if family != 0 && v.Family != family {
			add = false
		}
		if add {
			addr = append(addr, v)
		}
	}

	log.Printf("%#v", addr)
}
