package dhcpv6

import (
	"fmt"

	"github.com/u-root/u-root/pkg/uio"
)

// OptIANA implements the identity association for non-temporary addresses
// option.
//
// This module defines the OptIANA structure.
// https://www.ietf.org/rfc/rfc3633.txt
type OptIANA struct {
	IaId    [4]byte
	T1      uint32
	T2      uint32
	Options Options
}

func (op *OptIANA) Code() OptionCode {
	return OptionIANA
}

// ToBytes serializes IANA to DHCPv6 bytes.
func (op *OptIANA) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	buf.WriteBytes(op.IaId[:])
	buf.Write32(op.T1)
	buf.Write32(op.T2)
	buf.WriteBytes(op.Options.ToBytes())
	return buf.Data()
}

func (op *OptIANA) String() string {
	return fmt.Sprintf("OptIANA{IAID=%v, t1=%v, t2=%v, options=%v}",
		op.IaId, op.T1, op.T2, op.Options)
}

// AddOption adds an option at the end of the IA_NA options
func (op *OptIANA) AddOption(opt Option) {
	op.Options.Add(opt)
}

// GetOneOption will get an option of the give type from the Options field, if
// it is present. It will return `nil` otherwise
func (op *OptIANA) GetOneOption(code OptionCode) Option {
	return op.Options.GetOne(code)
}

// DelOption will remove all the options that match a Option code.
func (op *OptIANA) DelOption(code OptionCode) {
	op.Options.Del(code)
}

// ParseOptIANA builds an OptIANA structure from a sequence of bytes.  The
// input data does not include option code and length bytes.
func ParseOptIANA(data []byte) (*OptIANA, error) {
	var opt OptIANA
	buf := uio.NewBigEndianBuffer(data)
	buf.ReadBytes(opt.IaId[:])
	opt.T1 = buf.Read32()
	opt.T2 = buf.Read32()
	if err := opt.Options.FromBytes(buf.ReadAll()); err != nil {
		return nil, err
	}
	return &opt, buf.FinError()
}
