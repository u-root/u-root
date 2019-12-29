package rtnetlink_test

import (
	"log"
	"net"

	"github.com/jsimonetti/rtnetlink"
	"golang.org/x/sys/unix"
)

// List all neighbors on interface 'lo'
func Example_listNeighbors() {
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

	// Request all neighbors
	msg, err := conn.Neigh.List()
	if err != nil {
		log.Fatal(err)
	}

	// Filter neighbors by family and interface index
	var neigh []rtnetlink.NeighMessage
	for _, v := range msg {
		add := true
		if iface != nil && v.Index != uint32(iface.Index) {
			add = false
		}
		if family != 0 && v.Family != uint16(family) {
			add = false
		}
		if add {
			neigh = append(neigh, v)
		}
	}

	log.Printf("%#v", neigh)
}
