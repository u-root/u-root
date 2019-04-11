package rtnetlink

import (
	"errors"
	"fmt"
	"net"

	"github.com/mdlayher/netlink"
	"github.com/mdlayher/netlink/nlenc"
	"golang.org/x/sys/unix"
)

var (
	// errInvalidaddressMessage is returned when a AddressMessage is malformed.
	errInvalidAddressMessage = errors.New("rtnetlink AddressMessage is invalid or too short")

	// errInvalidAddressMessageAttr is returned when link attributes are malformed.
	errInvalidAddressMessageAttr = errors.New("rtnetlink AddressMessage has a wrong attribute data length")
)

var _ Message = &AddressMessage{}

// A AddressMessage is a route netlink address message.
type AddressMessage struct {
	// Address family (current unix.AF_INET or unix.AF_INET6)
	Family uint8

	// Prefix length
	PrefixLength uint8

	// Contains address flags
	Flags uint8

	// Address Scope
	Scope uint8

	// Interface index
	Index uint32

	// Attributes List
	Attributes AddressAttributes
}

// MarshalBinary marshals a AddressMessage into a byte slice.
func (m *AddressMessage) MarshalBinary() ([]byte, error) {
	b := make([]byte, unix.SizeofIfAddrmsg)

	b[0] = m.Family
	b[1] = m.PrefixLength
	b[2] = m.Flags
	b[3] = m.Scope
	nlenc.PutUint32(b[4:8], m.Index)

	a, err := m.Attributes.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return append(b, a...), nil
}

// UnmarshalBinary unmarshals the contents of a byte slice into a AddressMessage.
func (m *AddressMessage) UnmarshalBinary(b []byte) error {
	l := len(b)
	if l < unix.SizeofIfAddrmsg {
		return errInvalidAddressMessage
	}

	m.Family = uint8(b[0])
	m.PrefixLength = uint8(b[1])
	m.Flags = uint8(b[2])
	m.Scope = uint8(b[3])
	m.Index = nlenc.Uint32(b[4:8])

	if l > unix.SizeofIfAddrmsg {
		m.Attributes = AddressAttributes{}
		err := m.Attributes.UnmarshalBinary(b[unix.SizeofIfAddrmsg:])
		if err != nil {
			return err
		}
	}

	return nil
}

// rtMessage is an empty method to sattisfy the Message interface.
func (*AddressMessage) rtMessage() {}

// AddressService is used to retrieve rtnetlink family information.
type AddressService struct {
	c *Conn
}

// New creates a new address using the AddressMessage information.
func (a *AddressService) New(req *AddressMessage) error {
	flags := netlink.Request | netlink.Create | netlink.Acknowledge | netlink.Excl
	_, err := a.c.Execute(req, unix.RTM_NEWADDR, flags)
	if err != nil {
		return err
	}

	return nil
}

// Delete removes an address using the AddressMessage information.
func (a *AddressService) Delete(req *AddressMessage) error {
	flags := netlink.Request | netlink.Acknowledge
	_, err := a.c.Execute(req, unix.RTM_DELADDR, flags)
	if err != nil {
		return err
	}

	return nil
}

// List retrieves all addresses.
func (a *AddressService) List() ([]AddressMessage, error) {
	req := &AddressMessage{}

	flags := netlink.Request | netlink.Dump
	msgs, err := a.c.Execute(req, unix.RTM_GETADDR, flags)
	if err != nil {
		return nil, err
	}

	addresses := make([]AddressMessage, 0, len(msgs))
	for _, m := range msgs {
		address := (m).(*AddressMessage)
		addresses = append(addresses, *address)
	}
	return addresses, nil
}

// AddressAttributes contains all attributes for an interface.
type AddressAttributes struct {
	Address   net.IP // Interface Ip address
	Local     net.IP // Local Ip address
	Label     string
	Broadcast net.IP    // Broadcast Ip address
	Anycast   net.IP    // Anycast Ip address
	CacheInfo CacheInfo // Address information
	Multicast net.IP    // Multicast Ip address
	Flags     uint32    // Address flags
}

// UnmarshalBinary unmarshals the contents of a byte slice into a AddressMessage.
func (a *AddressAttributes) UnmarshalBinary(b []byte) error {
	attrs, err := netlink.UnmarshalAttributes(b)
	if err != nil {
		return err
	}
	for _, attr := range attrs {
		switch attr.Type {
		case unix.IFA_UNSPEC:
			//unused attribute
		case unix.IFA_ADDRESS:
			if len(attr.Data) != 4 && len(attr.Data) != 16 {
				return errInvalidAddressMessageAttr
			}
			a.Address = attr.Data
		case unix.IFA_LOCAL:
			if len(attr.Data) != 4 {
				return errInvalidAddressMessageAttr
			}
			a.Local = attr.Data
		case unix.IFA_LABEL:
			a.Label = nlenc.String(attr.Data)
		case unix.IFA_BROADCAST:
			if len(attr.Data) != 4 {
				return errInvalidAddressMessageAttr
			}
			a.Broadcast = attr.Data
		case unix.IFA_ANYCAST:
			if len(attr.Data) != 4 && len(attr.Data) != 16 {
				return errInvalidAddressMessageAttr
			}
			a.Anycast = attr.Data
		case unix.IFA_CACHEINFO:
			if len(attr.Data) != 16 {
				return errInvalidAddressMessageAttr
			}
			err := a.CacheInfo.UnmarshalBinary(attr.Data)
			if err != nil {
				return err
			}
		case unix.IFA_MULTICAST:
			if len(attr.Data) != 4 && len(attr.Data) != 16 {
				return errInvalidAddressMessageAttr
			}
			a.Multicast = attr.Data
		case unix.IFA_FLAGS:
			if len(attr.Data) != 4 {
				return errInvalidAddressMessageAttr
			}
			a.Flags = nlenc.Uint32(attr.Data)
		}
	}

	return nil
}

// MarshalBinary marshals a AddressAttributes into a byte slice.
func (a *AddressAttributes) MarshalBinary() ([]byte, error) {
	attrs := []netlink.Attribute{
		{
			Type: unix.IFA_UNSPEC,
			Data: nlenc.Uint16Bytes(0),
		},
		{
			Type: unix.IFA_ADDRESS,
			Data: a.Address,
		},
		{
			Type: unix.IFA_BROADCAST,
			Data: a.Broadcast,
		},
		{
			Type: unix.IFA_ANYCAST,
			Data: a.Anycast,
		},
		{
			Type: unix.IFA_MULTICAST,
			Data: a.Multicast,
		},
		{
			Type: unix.IFA_FLAGS,
			Data: nlenc.Uint32Bytes(a.Flags),
		},
	}

	if a.Local != nil {
		attrs = append(attrs, netlink.Attribute{
			Type: unix.IFA_LOCAL,
			Data: a.Local,
		})
	}

	return netlink.MarshalAttributes(attrs)
}

// CacheInfo contains address information
type CacheInfo struct {
	Prefered uint32
	Valid    uint32
	Created  uint32
	Updated  uint32
}

// UnmarshalBinary unmarshals the contents of a byte slice into a LinkMessage.
func (c *CacheInfo) UnmarshalBinary(b []byte) error {
	if len(b) != 16 {
		return fmt.Errorf("incorrect size, want: 16, got: %d", len(b))
	}

	c.Prefered = nlenc.Uint32(b[0:4])
	c.Valid = nlenc.Uint32(b[4:8])
	c.Created = nlenc.Uint32(b[8:12])
	c.Updated = nlenc.Uint32(b[12:16])

	return nil
}
