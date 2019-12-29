// +build

package rtnetlink_test

import (
	"log"
	"net"

	"github.com/jsimonetti/rtnetlink"
	"golang.org/x/sys/unix"
)

// Add a route
func Example_addRoute() {
	// Gather the interface Index
	iface, _ := net.InterfaceByName("lo")
	// Get a route to add
	_, dst, _ := net.ParseCIDR("192.168.0.0/16")
	// Get a gw to use
	gw := net.ParseIP("127.0.0.1")

	// Dial a connection to the rtnetlink socket
	conn, err := rtnetlink.Dial(nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	attr := rtnetlink.RouteAttributes{
		Dst:      dst.IP,
		OutIface: uint32(iface.Index),
	}
	if gw == nil {
		attr.Gateway = gw
	}
	ones, _ := dst.Mask.Size()

	err = conn.Route.Add(&rtnetlink.RouteMessage{
		Family:     unix.AF_INET,
		Table:      unix.RT_TABLE_MAIN,
		Protocol:   unix.RTPROT_BOOT,
		Scope:      unix.RT_SCOPE_LINK,
		Type:       unix.RTN_UNICAST,
		DstLength:  uint8(ones),
		Attributes: attr,
	})

	log.Fatal(err)
}
