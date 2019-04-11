package main

import (
	"fmt"
	"net"

	"golang.org/x/sys/unix"

	"github.com/jsimonetti/rtnetlink"
)

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
	conn, err := rtnetlink.Dial(nil)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	msg, err := conn.Link.List()
	if err != nil {
		return nil, err
	}

	// Cycle through to make sure the msg.Attributes is not nil
	// as we might access this later on
	for i := 0; i < len(msg); i++ {
		if msg[i].Attributes == nil {
			msg[i].Attributes = &rtnetlink.LinkAttributes{}
		}
	}
	return msg, nil
}

func linkSetHardwareAddr(iface *net.Interface, hwAddr net.HardwareAddr) error {
	conn, err := rtnetlink.Dial(nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	msg, err := conn.Link.Get(uint32(iface.Index))
	if err != nil {
		return err
	}

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

	return err
}

func linkSetUp(iface *net.Interface) error {
	conn, err := rtnetlink.Dial(nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	msg, err := conn.Link.Get(uint32(iface.Index))
	if err != nil {
		return err
	}

	state := msg.Attributes.OperationalState
	// If the link is already up, return immediately
	if state == rtnetlink.OperStateUp || state == rtnetlink.OperStateUnknown {
		return nil
	}

	err = conn.Link.Set(&rtnetlink.LinkMessage{
		Family: msg.Family,
		Type:   msg.Type,
		Index:  uint32(iface.Index),
		Flags:  unix.IFF_UP,
		Change: unix.IFF_UP,
	})

	return err
}

func linkSetDown(iface *net.Interface) error {
	conn, err := rtnetlink.Dial(nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	msg, err := conn.Link.Get(uint32(iface.Index))
	if err != nil {
		return err
	}

	state := msg.Attributes.OperationalState
	// If the link is already down, return immediately
	if state == rtnetlink.OperStateDown {
		return nil
	}

	err = conn.Link.Set(&rtnetlink.LinkMessage{
		Family: msg.Family,
		Type:   msg.Type,
		Index:  uint32(iface.Index),
		Flags:  0x0,
		Change: 0x1,
	})

	return err
}

func routeAdd(iface *net.Interface, dst net.IPNet, gw net.IP) error {
	return nil
}

func neighList(iface *net.Interface, family uint16) ([]rtnetlink.NeighMessage, error) {
	conn, err := rtnetlink.Dial(nil)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	msg, err := conn.Neigh.List()
	if err != nil {
		return nil, err
	}

	var neigh []rtnetlink.NeighMessage
	for _, v := range msg {
		add := true
		if iface != nil && v.Index != uint32(iface.Index) {
			add = false
		}
		if family != 0 && v.Family != family {
			add = false
		}
		if add {
			neigh = append(neigh, v)
		}
	}

	return neigh, nil
}

var flagNames = map[int]string{
	unix.IFF_UP:          "UP",
	unix.IFF_BROADCAST:   "BROADCAST",
	unix.IFF_POINTOPOINT: "POINTTOPOINT",
	unix.IFF_LOOPBACK:    "LOOPBACK",
	unix.IFF_LOWER_UP:    "LOWER_UP",
	unix.IFF_NOARP:       "NOARP",
	unix.IFF_MULTICAST:   "MULTICAST",
}

func linkFlags(f uint32) string {
	s := ""
	for i, name := range flagNames {
		if f&uint32(i) != 0 {
			if s != "" {
				s += ","
			}
			s += name
		}
	}
	if s == "" {
		s = "0"
	}
	return s
}

func encapType(t uint16) string {
	switch t {
	case 0:
		return "generic"
	case unix.ARPHRD_ETHER:
		return "ether"
	case unix.ARPHRD_ATM:
		return "atm"
	case unix.ARPHRD_PPP:
		return "ppp"
	case unix.ARPHRD_TUNNEL:
		return "ipip"
	case unix.ARPHRD_TUNNEL6:
		return "tunnel6"
	case unix.ARPHRD_LOOPBACK:
		return "loopback"
	case unix.ARPHRD_SIT:
		return "sit"
	case unix.ARPHRD_IPGRE:
		return "gre"
	case unix.ARPHRD_NETLINK:
		return "netlink"
	}
	return fmt.Sprintf("unknown%d", t)
}

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

var linkStates = map[rtnetlink.OperationalState]string{
	rtnetlink.OperStateUnknown:        "unknown",
	rtnetlink.OperStateNotPresent:     "not-present",
	rtnetlink.OperStateDown:           "down",
	rtnetlink.OperStateLowerLayerDown: "lower-layer-down",
	rtnetlink.OperStateTesting:        "testing",
	rtnetlink.OperStateDormant:        "dormant",
	rtnetlink.OperStateUp:             "up",
}
