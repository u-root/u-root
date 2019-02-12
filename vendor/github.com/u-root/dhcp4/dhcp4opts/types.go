// Copyright 2018 the u-root Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dhcp4opts

import (
	"io"
	"net"

	"github.com/u-root/dhcp4"
	"github.com/u-root/u-root/pkg/uio"
)

// DHCPMessageType implements encoding.BinaryMarshaler and encapsulates binary
// encoding and decoding methods for DHCP message types as specified by RFC
// 2132, Section 9.6.
type DHCPMessageType uint8

// Legal values of DHCP message types as per RFC 2132, Section 9.6.
const (
	DHCPDiscover DHCPMessageType = 1
	DHCPOffer    DHCPMessageType = 2
	DHCPRequest  DHCPMessageType = 3
	DHCPDecline  DHCPMessageType = 4
	DHCPACK      DHCPMessageType = 5
	DHCPNAK      DHCPMessageType = 6
	DHCPRelease  DHCPMessageType = 7
	DHCPInform   DHCPMessageType = 8
)

// MarshalBinary marshals the DHCP message type option to binary.
func (d DHCPMessageType) MarshalBinary() ([]byte, error) {
	return []byte{byte(d)}, nil
}

// UnmarshalBinary unmarshals the DHCP message type option from binary.
func (d *DHCPMessageType) UnmarshalBinary(p []byte) error {
	buf := uio.NewBigEndianBuffer(p)
	*d = DHCPMessageType(buf.Read8())
	return buf.FinError()
}

// SubnetMask implements encoding.BinaryMarshaler and encapsulates binary
// encoding and decoding methods for a subnet mask as specified by RFC 2132,
// Section 3.3.
type SubnetMask net.IPMask

// MarshalBinary writes the subnet mask option to binary.
func (s SubnetMask) MarshalBinary() ([]byte, error) {
	return []byte(s[:net.IPv4len]), nil
}

// UnmarshalBinary reads the subnet mask option from binary.
func (s *SubnetMask) UnmarshalBinary(p []byte) error {
	buf := uio.NewBigEndianBuffer(p)
	*s = buf.CopyN(net.IPv4len)
	return buf.FinError()
}

// IP implements encoding.BinaryMarshaler and encapsulates binary encoding and
// decoding for an IPv4 IP as defined by RFC 2132 for the options in Sections
// 3.18, 5.3, 5.7, 9.1, and 9.5.
type IP net.IP

// MarshalBinary writes the IP address to binary.
func (i IP) MarshalBinary() ([]byte, error) {
	return []byte(i[:net.IPv4len]), nil
}

// UnmarshalBinary reads the IP address from binary.
func (i *IP) UnmarshalBinary(p []byte) error {
	if len(p) < net.IPv4len {
		return io.ErrUnexpectedEOF
	}

	*i = make([]byte, net.IPv4len)
	copy(*i, p[:net.IPv4len])
	return nil
}

// GetIP returns the IP encoded in `code` option of `o`, if there is one.
func GetIP(code dhcp4.OptionCode, o dhcp4.Options) IP {
	v := o.Get(code)
	if v == nil {
		return nil
	}
	var ip IP
	if err := (&ip).UnmarshalBinary(v); err != nil {
		return nil
	}
	return ip
}

// IPs implements encoding.BinaryMarshaler and encapsulates binary encoding and
// decoding methods for a list of IPs as used by RFC 2132 for options in
// Sections 3.5 through 3.13, 8.2, 8.3, 8.5, 8.6, 8.9, and 8.10.
type IPs []net.IP

// MarshalBinary writes the list of IPs to binary.
func (i IPs) MarshalBinary() ([]byte, error) {
	b := uio.NewBigEndianBuffer(make([]byte, 0, net.IPv4len*len(i)))
	for _, ip := range i {
		b.WriteBytes(ip.To4())
	}
	return b.Data(), nil
}

// UnmarshalBinary reads a list of IPs from binary.
func (i *IPs) UnmarshalBinary(p []byte) error {
	b := uio.NewBigEndianBuffer(p)

	*i = make([]net.IP, 0, b.Len()/net.IPv4len)
	for b.Has(net.IPv4len) {
		*i = append(*i, net.IP(b.CopyN(net.IPv4len)))
	}
	return b.FinError()
}

// GetIPs returns the list of IPs encoded in `code` option of `o`.
func GetIPs(code dhcp4.OptionCode, o dhcp4.Options) IPs {
	v := o.Get(code)
	if v == nil {
		return nil
	}

	var i IPs
	if err := i.UnmarshalBinary(v); err != nil {
		return nil
	}
	return i
}

// String implements encoding.BinaryMarshaler and encapsulates binary encoding
// and decoding of strings as specified by RFC 2132 in Sections 3.14, 3.16,
// 3.17, 3.19, and 3.20.
type String string

// MarshalBinary writes the string to binary.
func (s String) MarshalBinary() ([]byte, error) {
	return []byte(s), nil
}

// GetString returns the string encoded in the `code` option of `o`.
func GetString(code dhcp4.OptionCode, o dhcp4.Options) string {
	v := o.Get(code)
	if v == nil {
		return ""
	}
	return string(v)
}

// OptionCodes implements encoding.BinaryMarshaler and encapsulates binary
// encoding and decoding methods of DHCP option codes as specified in RFC 2132
// Section 9.8.
type OptionCodes []dhcp4.OptionCode

// MarshalBinary writes the option code list to binary.
func (o OptionCodes) MarshalBinary() ([]byte, error) {
	b := uio.NewBigEndianBuffer(nil)
	for _, code := range o {
		b.Write8(uint8(code))
	}
	return b.Data(), nil
}

// UnmarshalBinary reads the option code list from binary.
func (o *OptionCodes) UnmarshalBinary(p []byte) error {
	b := uio.NewBigEndianBuffer(p)
	*o = make(OptionCodes, 0, b.Len())
	for b.Has(1) {
		*o = append(*o, dhcp4.OptionCode(b.Read8()))
	}
	return b.FinError()
}

// Uint16 implements encoding.BinaryMarshaler and encapsulates binary encoding
// and decoding methods of uint16s as defined by RFC 2132 Section 9.10.
type Uint16 uint16

// MarshalBinary writes the uint16 to binary.
func (u Uint16) MarshalBinary() ([]byte, error) {
	b := uio.NewBigEndianBuffer(nil)
	b.Write16(uint16(u))
	return b.Data(), nil
}

// UnmarshalBinary reads the uint16 from binary.
func (u *Uint16) UnmarshalBinary(p []byte) error {
	b := uio.NewBigEndianBuffer(p)
	*u = Uint16(b.Read16())
	return b.FinError()
}

// Uint32 implements encoding.BinaryMarshaler and encapsulates binary encoding
// and decoding methods of uint32s as defined by RFC 2132 Section 9.2.
type Uint32 uint32

// MarshalBinary writes the uint32 to binary.
func (u Uint32) MarshalBinary() ([]byte, error) {
	b := uio.NewBigEndianBuffer(nil)
	b.Write32(uint32(u))
	return b.Data(), nil
}

// UnmarshalBinary reads the uint32 from binary.
func (u *Uint32) UnmarshalBinary(p []byte) error {
	b := uio.NewBigEndianBuffer(p)
	*u = Uint32(b.Read32())
	return b.FinError()
}
