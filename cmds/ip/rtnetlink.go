package main

import (
	"net"

	"github.com/jsimonetti/rtnetlink"
)

// TODO: add these constants to golang.org/x/sys/unix
const (
	NUD_NONE       uint16 = 0x00
	NUD_INCOMPLETE uint16 = 0x01
	NUD_REACHABLE  uint16 = 0x02
	NUD_STALE      uint16 = 0x04
	NUD_DELAY      uint16 = 0x08
	NUD_PROBE      uint16 = 0x10
	NUD_FAILED     uint16 = 0x20
	NUD_NOARP      uint16 = 0x40
	NUD_PERMANENT  uint16 = 0x80
)

var neighStates = map[uint16]string{
	NUD_NONE:       "NONE",
	NUD_INCOMPLETE: "INCOMPLETE",
	NUD_REACHABLE:  "REACHABLE",
	NUD_STALE:      "STALE",
	NUD_DELAY:      "DELAY",
	NUD_PROBE:      "PROBE",
	NUD_FAILED:     "FAILED",
	NUD_NOARP:      "NOARP",
	NUD_PERMANENT:  "PERMANENT",
}

func addrAdd(iface *net.Interface, addr *net.IPNet) error {
	return nil
}

func addrDel(iface *net.Interface, addr *net.IPNet) error {
	return nil
}

func addrList(link *net.Interface, family uint8) ([]rtnetlink.AddressMessage, error) {
	return nil, nil
}

func linkList() ([]rtnetlink.LinkMessage, error) {
	var list []rtnetlink.LinkMessage

	for i := 0; i < len(list); i++ {
		if list[i].Attributes == nil {
			list[i].Attributes = &rtnetlink.LinkAttributes{}
		}
	}
	return nil, nil
}

func linkSetHardwareAddr(iface *net.Interface, hwAddr net.HardwareAddr) error {
	return nil
}

func linkSetUp(iface *net.Interface) error {
	return nil
}

func linkSetDown(iface *net.Interface) error {
	return nil
}

func routeAdd(iface *net.Interface, dst net.IPNet, gw net.IP) error {
	return nil
}

func neighList(iface *net.Interface, family uint8) ([]rtnetlink.NeighMessage, error) {
	return nil, nil
}
