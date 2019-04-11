package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"

	"golang.org/x/sys/unix"

	"github.com/jsimonetti/rtnetlink"
)

func addrAdd(iface *net.Interface, addr *net.IPNet) error {
	conn, err := rtnetlink.Dial(nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	family := unix.AF_INET6
	if addr.IP.To4() != nil {
		family = unix.AF_INET
	}
	ones, _ := addr.Mask.Size()

	brd, err := broadcastAddr(addr)
	if err != nil {
		return err
	}

	err = conn.Address.New(&rtnetlink.AddressMessage{
		Family:       uint8(family),
		PrefixLength: uint8(ones),
		Scope:        unix.RT_SCOPE_UNIVERSE,
		Index:        uint32(iface.Index),
		Attributes: rtnetlink.AddressAttributes{
			Address:   addr.IP,
			Local:     addr.IP,
			Broadcast: brd,
		},
	})

	return err
}

func addrDel(iface *net.Interface, addr *net.IPNet) error {
	conn, err := rtnetlink.Dial(nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	family := unix.AF_INET6
	if addr.IP.To4() != nil {
		family = unix.AF_INET
	}
	ones, _ := addr.Mask.Size()

	req := &rtnetlink.AddressMessage{}

	msg, err := addrList(iface, uint8(family))
	for _, v := range msg {
		if v.PrefixLength == uint8(ones) && v.Attributes.Address.Equal(addr.IP) {
			req = &rtnetlink.AddressMessage{
				Family:       uint8(family),
				PrefixLength: uint8(ones),
				Index:        uint32(iface.Index),
				Attributes: rtnetlink.AddressAttributes{
					Address:   addr.IP,
					Broadcast: v.Attributes.Broadcast,
				},
			}
		}
	}

	err = conn.Address.Delete(req)

	return err
}

func addrList(iface *net.Interface, family uint8) ([]rtnetlink.AddressMessage, error) {
	conn, err := rtnetlink.Dial(nil)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	msg, err := conn.Address.List()
	if err != nil {
		return nil, err
	}

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

	return addr, nil
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
	conn, err := rtnetlink.Dial(nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	attr := rtnetlink.RouteAttributes{
		Dst:      dst.IP,
		OutIface: uint32(iface.Index),
	}
	if gw == nil {
		attr.Gateway = gw
	}

	// TODO: fix this code
	err = conn.Route.Add(&rtnetlink.RouteMessage{
		Family:     unix.AF_INET,
		Table:      unix.RT_TABLE_MAIN,
		Protocol:   unix.RTPROT_BOOT,
		Scope:      unix.RT_SCOPE_LINK,
		Type:       unix.RTN_UNICAST,
		Attributes: attr,
	})

	return err
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

var addrFlagNames = map[int]string{
	unix.IFA_F_SECONDARY:      "secondary",
	unix.IFA_F_NODAD:          "nodad",
	unix.IFA_F_OPTIMISTIC:     "optimistic",
	unix.IFA_F_DADFAILED:      "dadfailed",
	unix.IFA_F_HOMEADDRESS:    "homeaddress",
	unix.IFA_F_DEPRECATED:     "deprecated",
	unix.IFA_F_TENTATIVE:      "tentative",
	unix.IFA_F_PERMANENT:      "permanent",
	unix.IFA_F_MANAGETEMPADDR: "managetempaddr",
	unix.IFA_F_NOPREFIXROUTE:  "noprefixroute",
	unix.IFA_F_MCAUTOJOIN:     "mcautojoin",
	unix.IFA_F_STABLE_PRIVACY: "stable-privacy",
}

func addrFlags(f uint32) string {
	s := ""
	for i, name := range addrFlagNames {
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

// TODO: fix this to work for ipv6 too
func broadcastAddr(n *net.IPNet) (net.IP, error) {
	if n.IP.To4() == nil {
		return net.IP{}, errors.New("does not support IPv6 addresses.")
	}
	ip := make(net.IP, len(n.IP.To4()))
	binary.BigEndian.PutUint32(ip, binary.BigEndian.Uint32(n.IP.To4())|^binary.BigEndian.Uint32(net.IP(n.Mask).To4()))
	return ip, nil
}
