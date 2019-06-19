package rtnl

import (
	"net"

	"github.com/jsimonetti/rtnetlink"
)

// Neigh represents a neighbour table entry (e.g. an entry in the ARP table)
type Neigh struct {
	HwAddr    net.HardwareAddr // Link-layer address
	IP        net.IP           // Network-layer address
	Interface *net.Interface   // Network interface
}

// Neighbours lists entries from the neighbor table (e.g. the ARP table).
func (c *Conn) Neighbours(ifc *net.Interface, family int) (r []*Neigh, err error) {
	rx, err := c.Conn.Neigh.List()
	if err != nil {
		return nil, err
	}
	match := func(v *rtnetlink.NeighMessage, ifc *net.Interface, family int) bool {
		if ifc != nil && v.Index != uint32(ifc.Index) {
			return false
		}
		if family != 0 && v.Family != uint16(family) {
			return false
		}
		return true
	}
	ifcache := map[int]*net.Interface{}
	for _, m := range rx {
		if !match(&m, ifc, family) {
			continue
		}
		ifindex := int(m.Index)
		iface, ok := ifcache[ifindex]
		if !ok {
			iface, err = c.LinkByIndex(ifindex)
			if err != nil {
				return nil, err
			}
			ifcache[ifindex] = iface
		}
		p := &Neigh{
			HwAddr:    m.Attributes.LLAddress,
			IP:        m.Attributes.Address,
			Interface: iface,
		}
		r = append(r, p)
	}
	return r, nil
}
