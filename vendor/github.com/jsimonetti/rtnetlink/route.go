package rtnetlink

import (
	"errors"
	"net"
	"unsafe"

	"github.com/jsimonetti/rtnetlink/internal/unix"

	"github.com/mdlayher/netlink"
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
	nativeEndian.PutUint32(b[8:12], m.Flags)

	ae := netlink.NewAttributeEncoder()
	err := m.Attributes.encode(ae)
	if err != nil {
		return nil, err
	}

	a, err := ae.Encode()
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
	m.Flags = nativeEndian.Uint32(b[8:12])

	if l > unix.SizeofRtMsg {
		m.Attributes = RouteAttributes{}
		ad, err := netlink.NewAttributeDecoder(b[unix.SizeofRtMsg:])
		if err != nil {
			return err
		}
		err = m.Attributes.decode(ad)
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

// Replace or add new route
func (r *RouteService) Replace(req *RouteMessage) error {
	flags := netlink.Request | netlink.Create | netlink.Replace | netlink.Acknowledge
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
	Dst       net.IP
	Src       net.IP
	Gateway   net.IP
	OutIface  uint32
	Priority  uint32
	Table     uint32
	Mark      uint32
	Expires   *uint32
	Metrics   *RouteMetrics
	Multipath []NextHop
}

func (a *RouteAttributes) decode(ad *netlink.AttributeDecoder) error {

	for ad.Next() {
		switch ad.Type() {
		case unix.RTA_UNSPEC:
			//unused attribute
		case unix.RTA_DST:
			l := len(ad.Bytes())
			if l != 4 && l != 16 {
				return errInvalidRouteMessageAttr
			}
			a.Dst = ad.Bytes()
		case unix.RTA_PREFSRC:
			l := len(ad.Bytes())
			if l != 4 && l != 16 {
				return errInvalidRouteMessageAttr
			}
			a.Src = ad.Bytes()
		case unix.RTA_GATEWAY:
			l := len(ad.Bytes())
			if l != 4 && l != 16 {
				return errInvalidRouteMessageAttr
			}
			a.Gateway = ad.Bytes()
		case unix.RTA_OIF:
			a.OutIface = ad.Uint32()
		case unix.RTA_PRIORITY:
			a.Priority = ad.Uint32()
		case unix.RTA_TABLE:
			a.Table = ad.Uint32()
		case unix.RTA_MARK:
			a.Mark = ad.Uint32()
		case unix.RTA_EXPIRES:
			timeout := ad.Uint32()
			a.Expires = &timeout
		case unix.RTA_METRICS:
			a.Metrics = &RouteMetrics{}
			ad.Nested(a.Metrics.decode)
		case unix.RTA_MULTIPATH:
			ad.Do(a.parseMultipath)
		}
	}

	return ad.Err()
}

func (a *RouteAttributes) encode(ae *netlink.AttributeEncoder) error {
	if a.Dst != nil {
		if ipv4 := a.Dst.To4(); ipv4 == nil {
			// Dst Addr is IPv6
			ae.Bytes(unix.RTA_DST, a.Dst)
		} else {
			// Dst Addr is IPv4
			ae.Bytes(unix.RTA_DST, ipv4)
		}
	}

	if a.Src != nil {
		if ipv4 := a.Src.To4(); ipv4 == nil {
			// Src Addr is IPv6
			ae.Bytes(unix.RTA_PREFSRC, a.Src)
		} else {
			// Src Addr is IPv4
			ae.Bytes(unix.RTA_PREFSRC, ipv4)
		}
	}

	if a.Gateway != nil {
		if ipv4 := a.Gateway.To4(); ipv4 == nil {
			// Gateway Addr is IPv6
			ae.Bytes(unix.RTA_GATEWAY, a.Gateway)
		} else {
			// Gateway Addr is IPv4
			ae.Bytes(unix.RTA_GATEWAY, ipv4)
		}
	}

	if a.OutIface != 0 {
		ae.Uint32(unix.RTA_OIF, a.OutIface)
	}

	if a.Priority != 0 {
		ae.Uint32(unix.RTA_PRIORITY, a.Priority)
	}

	if a.Table != 0 {
		ae.Uint32(unix.RTA_TABLE, a.Table)
	}

	if a.Mark != 0 {
		ae.Uint32(unix.RTA_MARK, a.Mark)
	}

	if a.Expires != nil {
		ae.Uint32(unix.RTA_EXPIRES, *a.Expires)
	}

	if a.Metrics != nil {
		ae.Nested(unix.RTA_METRICS, a.Metrics.encode)
	}

	if len(a.Multipath) > 0 {
		ae.Do(unix.RTA_MULTIPATH, a.encodeMultipath)
	}

	return nil
}

// RouteMetrics holds some advanced metrics for a route
type RouteMetrics struct {
	AdvMSS   uint32
	Features uint32
	InitCwnd uint32
	MTU      uint32
}

func (rm *RouteMetrics) decode(ad *netlink.AttributeDecoder) error {
	for ad.Next() {
		switch ad.Type() {
		case unix.RTAX_ADVMSS:
			rm.AdvMSS = ad.Uint32()
		case unix.RTAX_FEATURES:
			rm.Features = ad.Uint32()
		case unix.RTAX_INITCWND:
			rm.InitCwnd = ad.Uint32()
		case unix.RTAX_MTU:
			rm.MTU = ad.Uint32()
		}
	}

	// ad.Err call handled by Nested method in calling attribute decoder.
	return nil
}

func (rm *RouteMetrics) encode(ae *netlink.AttributeEncoder) error {
	if rm.AdvMSS != 0 {
		ae.Uint32(unix.RTAX_ADVMSS, rm.AdvMSS)
	}

	if rm.Features != 0 {
		ae.Uint32(unix.RTAX_FEATURES, rm.Features)
	}

	if rm.InitCwnd != 0 {
		ae.Uint32(unix.RTAX_INITCWND, rm.InitCwnd)
	}

	if rm.MTU != 0 {
		ae.Uint32(unix.RTAX_MTU, rm.MTU)
	}

	return nil
}

// TODO(mdlayher): probably eliminate Length field from the API to avoid the
// caller possibly tampering with it since we can compute it.

// RTNextHop represents the netlink rtnexthop struct (not an attribute)
type RTNextHop struct {
	Length  uint16 // length of this hop including nested values
	Flags   uint8  // flags defined in rtnetlink.h line 311
	Hops    uint8
	IfIndex uint32 // the interface index number
}

// NextHop wraps struct rtnexthop to provide access to nested attributes
type NextHop struct {
	Hop     RTNextHop // a rtnexthop struct
	Gateway net.IP    // that struct's nested Gateway attribute
}

func (a *RouteAttributes) encodeMultipath() ([]byte, error) {
	var b []byte
	for _, nh := range a.Multipath {
		// Encode the attributes first so their total length can be used to
		// compute the length of each (rtnexthop, attributes) pair.
		ae := netlink.NewAttributeEncoder()

		if a.Gateway != nil {
			// TODO(mdlayher): more validation.
			ae.Bytes(unix.RTA_GATEWAY, nh.Gateway)
		}

		ab, err := ae.Encode()
		if err != nil {
			return nil, err
		}

		// Assume the caller wants the length updated so they don't have to
		// keep track of it themselves when encoding attributes.
		nh.Hop.Length = unix.SizeofRtNexthop + uint16(len(ab))
		var nhb [unix.SizeofRtNexthop]byte

		copy(
			nhb[:],
			(*(*[unix.SizeofRtNexthop]byte)(unsafe.Pointer(&nh.Hop)))[:],
		)

		// rtnexthop first, then attributes.
		b = append(b, nhb[:]...)
		b = append(b, ab...)
	}

	return b, nil
}

// parseMultipath consumes RTA_MULTIPATH data into RouteAttributes.
func (a *RouteAttributes) parseMultipath(b []byte) error {
	// We cannot retain b after the function returns, so make a copy of the
	// bytes up front for the multipathParser.
	buf := make([]byte, len(b))
	copy(buf, b)

	// Iterate until no more bytes remain in the buffer or an error occurs.
	mpp := &multipathParser{b: buf}
	for mpp.Next() {
		// Each iteration reads a fixed length RTNextHop structure immediately
		// followed by its associated netlink attributes with optional data.
		nh := NextHop{Hop: mpp.RTNextHop()}
		if err := nh.decode(mpp.AttributeDecoder()); err != nil {
			return err
		}

		// Stop iteration early if the data was malformed, or otherwise append
		// this NextHop to the Multipath field.
		if err := mpp.Err(); err != nil {
			return err
		}

		a.Multipath = append(a.Multipath, nh)
	}

	// Check the error when Next returns false.
	return mpp.Err()
}

// rtnexthop payload is at least one nested attribute RTA_GATEWAY
// possibly others?
func (nh *NextHop) decode(ad *netlink.AttributeDecoder) error {
	if ad == nil {
		// Invalid decoder, do nothing.
		return nil
	}

	for ad.Next() {
		switch ad.Type() {
		case unix.RTA_GATEWAY:
			l := len(ad.Bytes())
			if l != 4 && l != 16 {
				return errInvalidRouteMessageAttr
			}

			nh.Gateway = ad.Bytes()
		}
	}

	return ad.Err()
}

// A multipathParser parses packed RTNextHop and netlink attributes into
// multipath attributes for an rtnetlink route.
type multipathParser struct {
	// Any errors which occurred during parsing.
	err error

	// The underlying buffer and a pointer to the reading position.
	b []byte
	i int

	// The length of the next set of netlink attributes.
	alen int
}

// Next continues iteration until an error occurs or no bytes remain.
func (mpp *multipathParser) Next() bool {
	if mpp.err != nil {
		return false
	}

	// Are there enough bytes left for another RTNextHop, or 0 for EOF?
	n := len(mpp.b[mpp.i:])
	switch {
	case n == 0:
		// EOF.
		return false
	case n >= unix.SizeofRtNexthop:
		return true
	default:
		mpp.err = errInvalidRouteMessageAttr
		return false
	}
}

// Err returns any errors encountered while parsing.
func (mpp *multipathParser) Err() error { return mpp.err }

// RTNextHop parses the next RTNextHop structure from the buffer.
func (mpp *multipathParser) RTNextHop() RTNextHop {
	if mpp.err != nil {
		return RTNextHop{}
	}

	if len(mpp.b)-mpp.i < unix.SizeofRtNexthop {
		// Out of bounds access, not enough data for a valid RTNextHop.
		mpp.err = errInvalidRouteMessageAttr
		return RTNextHop{}
	}

	// Consume an RTNextHop from the buffer by copying its bytes into an output
	// structure while also verifying that the size of each structure is equal
	// to avoid any out-of-bounds unsafe memory access.
	var rtnh RTNextHop
	next := mpp.b[mpp.i : mpp.i+unix.SizeofRtNexthop]

	if unix.SizeofRtNexthop != len(next) {
		panic("rtnetlink: invalid RTNextHop structure size, panicking to avoid out-of-bounds unsafe access")
	}

	copy(
		(*(*[unix.SizeofRtNexthop]byte)(unsafe.Pointer(&rtnh)))[:],
		(*(*[unix.SizeofRtNexthop]byte)(unsafe.Pointer(&next[0])))[:],
	)

	if rtnh.Length < unix.SizeofRtNexthop {
		// Length value is invalid.
		mpp.err = errInvalidRouteMessageAttr
		return RTNextHop{}
	}

	// Compute the length of the next set of attributes using the Length value
	// in the RTNextHop, minus the size of that fixed length structure itself.
	// Then, advance the pointer to be ready to read those attributes.
	mpp.alen = int(rtnh.Length) - unix.SizeofRtNexthop
	mpp.i += unix.SizeofRtNexthop

	return rtnh
}

// AttributeDecoder returns a netlink.AttributeDecoder pointed at the next set
// of netlink attributes from the buffer.
func (mpp *multipathParser) AttributeDecoder() *netlink.AttributeDecoder {
	if mpp.err != nil {
		return nil
	}

	// Ensure the attributes length value computed while parsing the rtnexthop
	// fits within the actual slice.
	if len(mpp.b[mpp.i:]) < mpp.alen {
		mpp.err = errInvalidRouteMessageAttr
		return nil
	}

	// Consume the next set of netlink attributes from the buffer and advance
	// the pointer to the next RTNextHop or EOF once that is complete.
	ad, err := netlink.NewAttributeDecoder(mpp.b[mpp.i : mpp.i+mpp.alen])
	if err != nil {
		mpp.err = err
		return nil
	}

	mpp.i += mpp.alen

	return ad
}
