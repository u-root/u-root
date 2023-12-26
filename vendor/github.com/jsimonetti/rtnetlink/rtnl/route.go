package rtnl

import (
	"errors"
	"fmt"
	"net"

	"github.com/jsimonetti/rtnetlink/internal/unix"

	"github.com/jsimonetti/rtnetlink"
)

// Route represents a route table entry
type Route struct {
	Destination *net.IPNet
	Gateway     net.IP
	Interface   *net.Interface
	Metric      uint32
}

// generating route message
func genRouteMessage(ifc *net.Interface, dst net.IPNet, gw net.IP, options ...RouteOption) (rm *rtnetlink.RouteMessage, err error) {
	opts := DefaultRouteOptions(ifc, dst, gw)

	for _, option := range options {
		option(opts)
	}

	af, err := addrFamily(dst.IP)
	if err != nil {
		return nil, err
	}

	// Determine scope
	var scope uint8
	switch {
	case gw != nil:
		scope = unix.RT_SCOPE_UNIVERSE
	case len(dst.IP) == net.IPv6len && dst.IP.To4() == nil:
		scope = unix.RT_SCOPE_UNIVERSE
	default:
		// Set default scope to LINK
		scope = unix.RT_SCOPE_LINK
	}

	var srclen int
	if opts.Src != nil {
		srclen, _ = opts.Src.Mask.Size()
		opts.Attrs.Src = opts.Src.IP
	}

	dstlen, _ := dst.Mask.Size()

	tx := &rtnetlink.RouteMessage{
		Family:     uint8(af),
		Table:      unix.RT_TABLE_MAIN,
		Protocol:   unix.RTPROT_BOOT,
		Type:       unix.RTN_UNICAST,
		Scope:      scope,
		DstLength:  uint8(dstlen),
		SrcLength:  uint8(srclen),
		Attributes: opts.Attrs,
	}
	return tx, nil
}

// RouteAdd adds infomation about a network route.
func (c *Conn) RouteAdd(ifc *net.Interface, dst net.IPNet, gw net.IP, options ...RouteOption) (err error) {
	rm, err := genRouteMessage(ifc, dst, gw, options...)
	if err != nil {
		return err
	}

	return c.Conn.Route.Add(rm)
}

// RouteReplace adds or replace information about a network route.
func (c *Conn) RouteReplace(ifc *net.Interface, dst net.IPNet, gw net.IP, options ...RouteOption) (err error) {
	rm, err := genRouteMessage(ifc, dst, gw, options...)
	if err != nil {
		return err
	}
	return c.Conn.Route.Replace(rm)
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

// RouteGet gets a single route to the given destination address.
func (c *Conn) RouteGet(dst net.IP) (*Route, error) {
	list, err := c.RouteGetAll(dst)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, errors.New("route wrong length")
	}

	return list[0], nil
}

// RouteGetAll returns all routes to the given destination IP in the main routing table.
func (c *Conn) RouteGetAll(dst net.IP) (ret []*Route, err error) {
	af, err := addrFamily(dst)
	if err != nil {
		return nil, err
	}

	attr := rtnetlink.RouteAttributes{
		Dst: dst,
	}

	tx := &rtnetlink.RouteMessage{
		Family:     uint8(af),
		Table:      unix.RT_TABLE_MAIN,
		Attributes: attr,
	}

	rx, err := c.Conn.Route.Get(tx)
	if err != nil {
		return nil, err
	}

	for _, rt := range rx {
		ifindex := int(rt.Attributes.OutIface)

		iface, err := c.LinkByIndex(ifindex)
		if err != nil {
			return nil, fmt.Errorf("failed to get link by interface index: %w", err)
		}

		_, dstNet, err := net.ParseCIDR(fmt.Sprintf("%s/%d", rt.Attributes.Dst.String(), rt.DstLength))
		if err != nil {
			return nil, fmt.Errorf("failed to construct CIDR from route destination address and length: %w", err)
		}

		ret = append(ret, &Route{
			Destination: dstNet,
			Gateway:     rt.Attributes.Gateway,
			Interface:   iface,
			Metric:      rt.Attributes.Priority,
		})
	}

	return ret, nil
}
