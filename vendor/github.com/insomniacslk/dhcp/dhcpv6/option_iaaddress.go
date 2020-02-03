package dhcpv6

import (
	"fmt"
	"net"
	"time"

	"github.com/u-root/u-root/pkg/uio"
)

// OptIAAddress represents an OptionIAAddr.
//
// This module defines the OptIAAddress structure.
// https://www.ietf.org/rfc/rfc3633.txt
type OptIAAddress struct {
	IPv6Addr          net.IP
	PreferredLifetime time.Duration
	ValidLifetime     time.Duration
	Options           Options
}

// Code returns the option's code
func (op *OptIAAddress) Code() OptionCode {
	return OptionIAAddr
}

// ToBytes serializes the option and returns it as a sequence of bytes
func (op *OptIAAddress) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	buf.WriteBytes(op.IPv6Addr.To16())

	t1 := Duration{op.PreferredLifetime}
	t1.Marshal(buf)
	t2 := Duration{op.ValidLifetime}
	t2.Marshal(buf)

	buf.WriteBytes(op.Options.ToBytes())
	return buf.Data()
}

func (op *OptIAAddress) String() string {
	return fmt.Sprintf("OptIAAddress{ipv6addr=%v, preferredlifetime=%v, validlifetime=%v, options=%v}",
		op.IPv6Addr, op.PreferredLifetime, op.ValidLifetime, op.Options)
}

// ParseOptIAAddress builds an OptIAAddress structure from a sequence
// of bytes. The input data does not include option code and length
// bytes.
func ParseOptIAAddress(data []byte) (*OptIAAddress, error) {
	var opt OptIAAddress
	buf := uio.NewBigEndianBuffer(data)
	opt.IPv6Addr = net.IP(buf.CopyN(net.IPv6len))

	var t1, t2 Duration
	t1.Unmarshal(buf)
	t2.Unmarshal(buf)
	opt.PreferredLifetime = t1.Duration
	opt.ValidLifetime = t2.Duration

	if err := opt.Options.FromBytes(buf.ReadAll()); err != nil {
		return nil, err
	}
	return &opt, buf.FinError()
}
