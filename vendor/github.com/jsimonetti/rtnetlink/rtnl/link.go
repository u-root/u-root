package rtnl

import (
	"net"

	"github.com/jsimonetti/rtnetlink/internal/unix"

	"github.com/jsimonetti/rtnetlink"
)

// Links return the list of interfaces available on the system.
func (c *Conn) Links() (r []*net.Interface, err error) {
	rx, err := c.Conn.Link.List()
	if err != nil {
		return nil, err
	}
	for _, m := range rx {
		ifc := linkmsgToInterface(&m)
		r = append(r, ifc)
	}
	return r, nil
}

// LinkByIndex returns an interface by its index. Similar to net.InterfaceByIndex.
func (c *Conn) LinkByIndex(ifindex int) (*net.Interface, error) {
	rx, err := c.Conn.Link.Get(uint32(ifindex))
	if err != nil {
		return nil, err
	}
	return linkmsgToInterface(&rx), nil
}

func linkmsgToInterface(m *rtnetlink.LinkMessage) *net.Interface {
	ifc := &net.Interface{
		Index:        int(m.Index),
		MTU:          int(m.Attributes.MTU),
		Name:         m.Attributes.Name,
		HardwareAddr: m.Attributes.Address,
		Flags:        linkFlags(m.Flags),
	}
	return ifc
}

func linkFlags(rawFlags uint32) net.Flags {
	var f net.Flags
	if rawFlags&unix.IFF_UP != 0 {
		f |= net.FlagUp
	}
	if rawFlags&unix.IFF_BROADCAST != 0 {
		f |= net.FlagBroadcast
	}
	if rawFlags&unix.IFF_LOOPBACK != 0 {
		f |= net.FlagLoopback
	}
	if rawFlags&unix.IFF_POINTOPOINT != 0 {
		f |= net.FlagPointToPoint
	}
	if rawFlags&unix.IFF_MULTICAST != 0 {
		f |= net.FlagMulticast
	}
	return f
}

// LinkSetHardwareAddr overrides the L2 address (MAC-address) for the interface.
func (c *Conn) LinkSetHardwareAddr(ifc *net.Interface, hw net.HardwareAddr) error {
	rx, err := c.Conn.Link.Get(uint32(ifc.Index))
	if err != nil {
		return err
	}
	tx := &rtnetlink.LinkMessage{
		Family: unix.AF_UNSPEC,
		Type:   rx.Type,
		Index:  uint32(ifc.Index),
		Flags:  rx.Flags,
		Change: 0, // rtnetlink(7) says it "should be always set to 0xFFFFFFFF" - BUG?
		Attributes: &rtnetlink.LinkAttributes{
			Address: hw,
			// some attributes are always included in ../link.go:/LinkAttributes/+/MarshalBinary/
			Name:      rx.Attributes.Name,
			MTU:       rx.Attributes.MTU,
			Type:      rx.Attributes.Type,
			QueueDisc: rx.Attributes.QueueDisc,
		},
	}
	return c.Conn.Link.Set(tx)
}

// LinkUp drives an inteface up, enabling the link.
func (c *Conn) LinkUp(ifc *net.Interface) error {
	rx, err := c.Conn.Link.Get(uint32(ifc.Index))
	if err != nil {
		return err
	}
	tx := &rtnetlink.LinkMessage{
		Family: unix.AF_UNSPEC,
		Type:   rx.Type,
		Index:  uint32(ifc.Index),
		Flags:  unix.IFF_UP,
		Change: unix.IFF_UP,
	}
	return c.Conn.Link.Set(tx)
}

// LinkDown takes an inteface down, disabling the link.
func (c *Conn) LinkDown(ifc *net.Interface) error {
	rx, err := c.Conn.Link.Get(uint32(ifc.Index))
	if err != nil {
		return err
	}
	tx := &rtnetlink.LinkMessage{
		Family: unix.AF_UNSPEC,
		Type:   rx.Type,
		Index:  uint32(ifc.Index),
		Flags:  0,
		Change: unix.IFF_UP,
	}
	return c.Conn.Link.Set(tx)
}
