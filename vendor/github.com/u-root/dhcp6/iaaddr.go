package dhcp6

import (
	"encoding/binary"
	"io"
	"net"
	"time"
)

// IAAddr represents an Identity Association Address, as defined in RFC 3315,
// Section 22.6.
//
// DHCP clients use identity assocation addresses (IAAddrs) to request IPv6
// addresses from a DHCP server, using the lifetimes specified in the preferred
// lifetime and valid lifetime fields.  Multiple IAAddrs may be present in a
// single DHCP request, but only enscapsulated within an IANA or IATA options
// field.
type IAAddr struct {
	// IP specifies the IPv6 address to offer to a client.  The validity of the
	// address is controlled by the PreferredLifetime and ValidLifetime fields.
	IP net.IP

	// PreferredLifetime specifies the preferred lifetime of an IPv6 address.
	// When the preferred lifetime of an address expires, the address becomes
	// deprecated, and should not be used in new communications.
	//
	// The preferred lifetime of an address must not be greater than its
	// valid lifetime.
	PreferredLifetime time.Duration

	// ValidLifetime specifies the valid lifetime of an IPv6 address.  When the
	// valid lifetime of an address expires, the address should not be used for
	// any further communication.
	//
	// The valid lifetime of an address must be greater than its preferred
	// lifetime.
	ValidLifetime time.Duration

	// Options specifies a map of DHCP options specific to this IAAddr.
	// Its methods can be used to retrieve data from an incoming IAAddr, or
	// send data with an outgoing IAAddr.
	Options Options
}

// NewIAAddr creates a new IAAddr from an IPv6 address, preferred and valid lifetime
// durations, and an optional Options map.
//
// The IP must be exactly 16 bytes, the correct length for an IPv6 address.
// The preferred lifetime duration must be less than the valid lifetime
// duration.  Failure to meet either of these conditions will result in an error.
// If an Options map is not specified, a new one will be allocated.
func NewIAAddr(ip net.IP, preferred time.Duration, valid time.Duration, options Options) (*IAAddr, error) {
	// From documentation: If ip is not an IPv4 address, To4 returns nil.
	if ip.To4() != nil {
		return nil, ErrInvalidIP
	}

	// Preferred lifetime must always be less than valid lifetime.
	if preferred > valid {
		return nil, ErrInvalidLifetimes
	}

	// If no options set, make empty map
	if options == nil {
		options = make(Options)
	}

	return &IAAddr{
		IP:                ip,
		PreferredLifetime: preferred,
		ValidLifetime:     valid,
		Options:           options,
	}, nil
}

// MarshalBinary allocates a byte slice containing the data from a IAAddr.
func (i *IAAddr) MarshalBinary() ([]byte, error) {
	// 16 bytes: IPv6 address
	//  4 bytes: preferred lifetime
	//  4 bytes: valid lifetime
	//  N bytes: options
	opts := i.Options.enumerate()
	b := make([]byte, 24+opts.count())

	copy(b[0:16], i.IP)
	binary.BigEndian.PutUint32(b[16:20], uint32(i.PreferredLifetime/time.Second))
	binary.BigEndian.PutUint32(b[20:24], uint32(i.ValidLifetime/time.Second))
	opts.write(b[24:])

	return b, nil
}

// UnmarshalBinary unmarshals a raw byte slice into a IAAddr.
//
// If the byte slice does not contain enough data to form a valid IAAddr,
// io.ErrUnexpectedEOF is returned.  If the preferred lifetime value in the
// byte slice is less than the valid lifetime, ErrInvalidLifetimes is returned.
func (i *IAAddr) UnmarshalBinary(b []byte) error {
	if len(b) < 24 {
		return io.ErrUnexpectedEOF
	}

	ip := make(net.IP, 16)
	copy(ip, b[0:16])
	i.IP = ip

	i.PreferredLifetime = time.Duration(binary.BigEndian.Uint32(b[16:20])) * time.Second
	i.ValidLifetime = time.Duration(binary.BigEndian.Uint32(b[20:24])) * time.Second

	// Preferred lifetime must always be less than valid lifetime.
	if i.PreferredLifetime > i.ValidLifetime {
		return ErrInvalidLifetimes
	}

	options, err := parseOptions(b[24:])
	if err != nil {
		return err
	}
	i.Options = options

	return nil
}
