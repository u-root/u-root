package dhcpv6

import (
	"fmt"
	"net"
	"time"

	"github.com/u-root/u-root/pkg/uio"
)

// OptIAPrefix implements the IAPrefix option.
//
// This module defines the OptIAPrefix structure.
// https://www.ietf.org/rfc/rfc3633.txt
type OptIAPrefix struct {
	PreferredLifetime time.Duration
	ValidLifetime     time.Duration
	prefixLength      byte
	ipv6Prefix        net.IP
	Options           Options
}

func (op *OptIAPrefix) Code() OptionCode {
	return OptionIAPrefix
}

// ToBytes marshals this option according to RFC 3633, Section 10.
func (op *OptIAPrefix) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)

	t1 := Duration{op.PreferredLifetime}
	t1.Marshal(buf)
	t2 := Duration{op.ValidLifetime}
	t2.Marshal(buf)

	buf.Write8(op.prefixLength)
	buf.WriteBytes(op.ipv6Prefix.To16())
	buf.WriteBytes(op.Options.ToBytes())
	return buf.Data()
}

func (op *OptIAPrefix) PrefixLength() byte {
	return op.prefixLength
}

func (op *OptIAPrefix) SetPrefixLength(pl byte) {
	op.prefixLength = pl
}

// IPv6Prefix returns the ipv6Prefix
func (op *OptIAPrefix) IPv6Prefix() net.IP {
	return op.ipv6Prefix
}

// SetIPv6Prefix sets the ipv6Prefix
func (op *OptIAPrefix) SetIPv6Prefix(p net.IP) {
	op.ipv6Prefix = p
}

func (op *OptIAPrefix) String() string {
	return fmt.Sprintf("OptIAPrefix{preferredlifetime=%v, validlifetime=%v, prefixlength=%v, ipv6prefix=%v, options=%v}",
		op.PreferredLifetime, op.ValidLifetime, op.PrefixLength(), op.IPv6Prefix(), op.Options)
}

// GetOneOption will get an option of the give type from the Options field, if
// it is present. It will return `nil` otherwise
func (op *OptIAPrefix) GetOneOption(code OptionCode) Option {
	return op.Options.GetOne(code)
}

// DelOption will remove all the options that match a Option code.
func (op *OptIAPrefix) DelOption(code OptionCode) {
	op.Options.Del(code)
}

// ParseOptIAPrefix an OptIAPrefix structure from a sequence of bytes. The
// input data does not include option code and length bytes.
func ParseOptIAPrefix(data []byte) (*OptIAPrefix, error) {
	buf := uio.NewBigEndianBuffer(data)
	var opt OptIAPrefix

	var t1, t2 Duration
	t1.Unmarshal(buf)
	t2.Unmarshal(buf)
	opt.PreferredLifetime = t1.Duration
	opt.ValidLifetime = t2.Duration

	opt.prefixLength = buf.Read8()
	opt.ipv6Prefix = net.IP(buf.CopyN(net.IPv6len))
	if err := opt.Options.FromBytes(buf.ReadAll()); err != nil {
		return nil, err
	}
	return &opt, buf.FinError()
}
