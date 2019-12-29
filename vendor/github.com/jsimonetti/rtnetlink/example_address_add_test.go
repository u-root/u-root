// +build linux

package rtnetlink_test

import (
	"encoding/binary"
	"log"
	"net"

	"github.com/jsimonetti/rtnetlink"
	"golang.org/x/sys/unix"
)

// Add IP address '127.0.0.2/8' to an interface 'lo'
func Example_addAddress() {
	// Gather the interface Index
	iface, _ := net.InterfaceByName("lo")
	// Get an ip address to add to the interface
	addr, cidr, _ := net.ParseCIDR("127.0.0.2/8")

	// Dial a connection to the rtnetlink socket
	conn, err := rtnetlink.Dial(nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Test for the right address family for addr
	family := unix.AF_INET6
	to4 := cidr.IP.To4()
	if to4 != nil {
		family = unix.AF_INET
	}
	// Calculate the prefix length
	ones, _ := cidr.Mask.Size()

	// Calculate the broadcast IP
	// Only used when family is AF_INET
	var brd net.IP
	if to4 != nil {
		brd = make(net.IP, len(to4))
		binary.BigEndian.PutUint32(brd, binary.BigEndian.Uint32(to4)|^binary.BigEndian.Uint32(net.IP(cidr.Mask).To4()))
	}

	// Send the message using the rtnetlink.Conn
	err = conn.Address.New(&rtnetlink.AddressMessage{
		Family:       uint8(family),
		PrefixLength: uint8(ones),
		Scope:        unix.RT_SCOPE_UNIVERSE,
		Index:        uint32(iface.Index),
		Attributes: rtnetlink.AddressAttributes{
			Address:   addr,
			Local:     addr,
			Broadcast: brd,
		},
	})

	log.Fatal(err)
}
