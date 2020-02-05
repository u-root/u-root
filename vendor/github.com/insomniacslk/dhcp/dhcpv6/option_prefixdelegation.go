package dhcpv6

import (
	"fmt"
	"time"

	"github.com/u-root/u-root/pkg/uio"
)

// OptIAForPrefixDelegation implements the identity association for prefix
// delegation option defined by RFC 3633, Section 9.
type OptIAForPrefixDelegation struct {
	IaId    [4]byte
	T1      time.Duration
	T2      time.Duration
	Options Options
}

// Code returns the option code
func (op *OptIAForPrefixDelegation) Code() OptionCode {
	return OptionIAPD
}

// ToBytes serializes the option and returns it as a sequence of bytes
func (op *OptIAForPrefixDelegation) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	buf.WriteBytes(op.IaId[:])

	t1 := Duration{op.T1}
	t1.Marshal(buf)
	t2 := Duration{op.T2}
	t2.Marshal(buf)

	buf.WriteBytes(op.Options.ToBytes())
	return buf.Data()
}

// String returns a string representation of the OptIAForPrefixDelegation data
func (op *OptIAForPrefixDelegation) String() string {
	return fmt.Sprintf("OptIAForPrefixDelegation{IAID=%v, t1=%v, t2=%v, options=%v}",
		op.IaId, op.T1, op.T2, op.Options)
}

// GetOneOption will get an option of the give type from the Options field, if
// it is present. It will return `nil` otherwise
func (op *OptIAForPrefixDelegation) GetOneOption(code OptionCode) Option {
	return op.Options.GetOne(code)
}

// DelOption will remove all the options that match a Option code.
func (op *OptIAForPrefixDelegation) DelOption(code OptionCode) {
	op.Options.Del(code)
}

// build an OptIAForPrefixDelegation structure from a sequence of bytes.
// The input data does not include option code and length bytes.
func ParseOptIAForPrefixDelegation(data []byte) (*OptIAForPrefixDelegation, error) {
	var opt OptIAForPrefixDelegation
	buf := uio.NewBigEndianBuffer(data)
	buf.ReadBytes(opt.IaId[:])

	var t1, t2 Duration
	t1.Unmarshal(buf)
	t2.Unmarshal(buf)
	opt.T1 = t1.Duration
	opt.T2 = t2.Duration

	if err := opt.Options.FromBytes(buf.ReadAll()); err != nil {
		return nil, err
	}
	return &opt, buf.FinError()
}
