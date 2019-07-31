package rtnetlink

import (
	"errors"
	"net"

	"github.com/mdlayher/netlink"
	"github.com/mdlayher/netlink/nlenc"
	"golang.org/x/sys/unix"
)

var (
	// errInvalidRouteMessage is returned when a RouteMessage is malformed.
	errInvalidRouteMessage = errors.New("rtnetlink RouteMessage is invalid or too short")

	// errInvalidRouteMessageAttr is returned when link attributes are malformed.
	errInvalidRouteMessageAttr = errors.New("rtnetlink RouteMessage has a wrong attribute data length")
)

var _ Message = &RouteMessage{}

type RouteMessage struct {
	Family    uint8 // Address family (current unix.AF_INET or unix.AF_INET6)
	DstLength uint8 // Length of destination prefix
	SrcLength uint8 // Length of source prefix
	Tos       uint8 // TOS filter
	Table     uint8 // Routing table ID
	Protocol  uint8 // Routing protocol
	Scope     uint8 // Distance to the destination
	Type      uint8 // Route type
	Flags     uint32

	Attributes RouteAttributes
}

func (m *RouteMessage) MarshalBinary() ([]byte, error) {
	b := make([]byte, unix.SizeofRtMsg)

	b[0] = m.Family
	b[1] = m.DstLength
	b[2] = m.SrcLength
	b[3] = m.Tos
	b[4] = m.Table
	b[5] = m.Protocol
	b[6] = m.Scope
	b[7] = m.Type
	nlenc.PutUint32(b[8:12], m.Flags)

	a, err := m.Attributes.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return append(b, a...), nil
}

func (m *RouteMessage) UnmarshalBinary(b []byte) error {
	l := len(b)
	if l < unix.SizeofRtMsg {
		return errInvalidRouteMessage
	}

	m.Family = uint8(b[0])
	m.DstLength = uint8(b[1])
	m.SrcLength = uint8(b[2])
	m.Tos = uint8(b[3])
	m.Table = uint8(b[4])
	m.Protocol = uint8(b[5])
	m.Scope = uint8(b[6])
	m.Type = uint8(b[7])
	m.Flags = nlenc.Uint32(b[8:12])

	if l > unix.SizeofRtMsg {
		m.Attributes = RouteAttributes{}
		err := m.Attributes.UnmarshalBinary(b[unix.SizeofRtMsg:])
		if err != nil {
			return err
		}
	}

	return nil
}

// rtMessage is an empty method to sattisfy the Message interface.
func (*RouteMessage) rtMessage() {}

type RouteService struct {
	c *Conn
}

// Add new route
func (r *RouteService) Add(req *RouteMessage) error {
	flags := netlink.Request | netlink.Create | netlink.Acknowledge | netlink.Excl
	_, err := r.c.Execute(req, unix.RTM_NEWROUTE, flags)
	if err != nil {
		return err
	}

	return nil
}

// Delete existing route
func (r *RouteService) Delete(req *RouteMessage) error {
	flags := netlink.Request | netlink.Acknowledge
	_, err := r.c.Execute(req, unix.RTM_DELROUTE, flags)
	if err != nil {
		return err
	}

	return nil
}

// Get Route(s)
func (r *RouteService) Get(req *RouteMessage) ([]RouteMessage, error) {
	flags := netlink.Request | netlink.DumpFiltered
	msgs, err := r.c.Execute(req, unix.RTM_GETROUTE, flags)
	if err != nil {
		return nil, err
	}

	routes := make([]RouteMessage, 0, len(msgs))
	for _, m := range msgs {
		route := (m).(*RouteMessage)
		routes = append(routes, *route)
	}

	return routes, nil
}

// List all routes
func (r *RouteService) List() ([]RouteMessage, error) {
	req := &RouteMessage{}

	flags := netlink.Request | netlink.Dump
	msgs, err := r.c.Execute(req, unix.RTM_GETROUTE, flags)
	if err != nil {
		return nil, err
	}

	routes := make([]RouteMessage, 0, len(msgs))
	for _, m := range msgs {
		route := (m).(*RouteMessage)
		routes = append(routes, *route)
	}

	return routes, nil
}

type RouteAttributes struct {
	Dst      net.IP
	Src      net.IP
	Gateway  net.IP
	OutIface uint32
	Priority uint32
	Table    uint32
}

func (a *RouteAttributes) UnmarshalBinary(b []byte) error {
	attrs, err := netlink.UnmarshalAttributes(b)
	if err != nil {
		return err
	}
	for _, attr := range attrs {
		switch attr.Type {
		case unix.RTA_UNSPEC:
		case unix.RTA_DST:
			if len(attr.Data) != 4 && len(attr.Data) != 16 {
				return errInvalidRouteMessageAttr
			}
			a.Dst = attr.Data
		case unix.RTA_PREFSRC:
			if len(attr.Data) != 4 && len(attr.Data) != 16 {
				return errInvalidRouteMessageAttr
			}
			a.Src = attr.Data
		case unix.RTA_GATEWAY:
			if len(attr.Data) != 4 && len(attr.Data) != 16 {
				return errInvalidRouteMessageAttr
			}
			a.Gateway = attr.Data
		case unix.RTA_OIF:
			if len(attr.Data) != 4 {
				return errInvalidRouteMessageAttr
			}
			a.OutIface = nlenc.Uint32(attr.Data)
		case unix.RTA_PRIORITY:
			if len(attr.Data) != 4 {
				return errInvalidRouteMessageAttr
			}
			a.Priority = nlenc.Uint32(attr.Data)
		case unix.RTA_TABLE:
			if len(attr.Data) != 4 {
				return errInvalidRouteMessageAttr
			}
			a.Table = nlenc.Uint32(attr.Data)
		}
	}

	return nil
}

func (a *RouteAttributes) MarshalBinary() ([]byte, error) {
	attrs := make([]netlink.Attribute, 0)

	if a.Dst != nil {
		if ipv4 := a.Dst.To4(); ipv4 == nil {
			// Dst Addr is IPv6
			attrs = append(attrs, netlink.Attribute{
				Type: unix.RTA_DST,
				Data: a.Dst,
			})
		} else {
			// Dst Addr is IPv4
			attrs = append(attrs, netlink.Attribute{
				Type: unix.RTA_DST,
				Data: ipv4,
			})
		}
	}

	if a.Src != nil {
		if ipv4 := a.Src.To4(); ipv4 == nil {
			// Src Addr is IPv6
			attrs = append(attrs, netlink.Attribute{
				Type: unix.RTA_PREFSRC,
				Data: a.Src,
			})
		} else {
			// Src Addr is IPv4
			attrs = append(attrs, netlink.Attribute{
				Type: unix.RTA_PREFSRC,
				Data: ipv4,
			})
		}
	}

	if a.Gateway != nil {
		if ipv4 := a.Gateway.To4(); ipv4 == nil {
			// Gateway Addr is IPv6
			attrs = append(attrs, netlink.Attribute{
				Type: unix.RTA_GATEWAY,
				Data: a.Gateway,
			})
		} else {
			// Gateway Addr is IPv4
			attrs = append(attrs, netlink.Attribute{
				Type: unix.RTA_GATEWAY,
				Data: ipv4,
			})
		}
	}

	if a.OutIface != 0 {
		attrs = append(attrs, netlink.Attribute{
			Type: unix.RTA_OIF,
			Data: nlenc.Uint32Bytes(a.OutIface),
		})
	}

	if a.Priority != 0 {
		attrs = append(attrs, netlink.Attribute{
			Type: unix.RTA_PRIORITY,
			Data: nlenc.Uint32Bytes(a.Priority),
		})
	}

	if a.Table != 0 {
		attrs = append(attrs, netlink.Attribute{
			Type: unix.RTA_TABLE,
			Data: nlenc.Uint32Bytes(a.Table),
		})
	}

	return netlink.MarshalAttributes(attrs)
}
