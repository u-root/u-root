package rtnl

import (
	"net"

	"golang.org/x/sys/unix"

	"github.com/jsimonetti/rtnetlink"
)

// RouteAdd adds infomation about a network route.
func (c *Conn) RouteAdd(ifc *net.Interface, dst net.IPNet, gw net.IP) (err error) {
	return c.RouteAddSrc(ifc, dst, nil, gw)
}

// RouteAddSrc adds infomation about a network route with the given destination
// and source. If source is `nil` it's ignored.
func (c *Conn) RouteAddSrc(ifc *net.Interface, dst net.IPNet, src *net.IPNet, gw net.IP) (err error) {
	af, err := addrFamily(dst.IP)
	if err != nil {
		return err
	}
	dstlen, _ := dst.Mask.Size()
	scope := addrScope(dst.IP)
	if len(dst.IP) == net.IPv6len && dst.IP.To4() == nil {
		scope = unix.RT_SCOPE_UNIVERSE
	}
	attr := rtnetlink.RouteAttributes{
		Dst:      dst.IP,
		OutIface: uint32(ifc.Index),
	}
	var srclen int
	if src != nil {
		srclen, _ = src.Mask.Size()
		attr.Src = src.IP
	}
	if gw != nil {
		attr.Gateway = gw
	}
	tx := &rtnetlink.RouteMessage{
		Family:     uint8(af),
		Table:      unix.RT_TABLE_MAIN,
		Protocol:   unix.RTPROT_BOOT,
		Type:       unix.RTN_UNICAST,
		Scope:      uint8(scope),
		DstLength:  uint8(dstlen),
		SrcLength:  uint8(srclen),
		Attributes: attr,
	}
	return c.Conn.Route.Add(tx)
}

// RouteDel deletes the route to the given destination.
func (c *Conn) RouteDel(ifc *net.Interface, dst net.IPNet) error {
	af, err := addrFamily(dst.IP)
	if err != nil {
		return err
	}
	prefixlen, _ := dst.Mask.Size()
	attr := rtnetlink.RouteAttributes{
		Dst:      dst.IP,
		OutIface: uint32(ifc.Index),
	}
	tx := &rtnetlink.RouteMessage{
		Family:     uint8(af),
		Table:      unix.RT_TABLE_MAIN,
		DstLength:  uint8(prefixlen),
		Attributes: attr,
	}
	return c.Conn.Route.Delete(tx)
}
