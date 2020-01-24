package dhcpv6

import (
	"fmt"
	"time"

	"github.com/u-root/u-root/pkg/uio"
)

// Duration is a duration as embedded in IA messages (IAPD, IANA, IATA).
type Duration struct {
	time.Duration
}

// Marshal encodes the time in uint32 seconds as defined by RFC 3315 for IANA
// messages.
func (d Duration) Marshal(buf *uio.Lexer) {
	buf.Write32(uint32(d.Duration.Round(time.Second) / time.Second))
}

// Unmarshal decodes time from uint32 seconds as defined by RFC 3315 for IANA
// messages.
func (d *Duration) Unmarshal(buf *uio.Lexer) {
	t := buf.Read32()
	d.Duration = time.Duration(t) * time.Second
}

// OptIANA implements the identity association for non-temporary addresses
// option.
//
// This module defines the OptIANA structure.
// https://www.ietf.org/rfc/rfc3633.txt
type OptIANA struct {
	IaId    [4]byte
	T1      time.Duration
	T2      time.Duration
	Options Options
}

func (op *OptIANA) Code() OptionCode {
	return OptionIANA
}

// ToBytes serializes IANA to DHCPv6 bytes.
func (op *OptIANA) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	buf.WriteBytes(op.IaId[:])
	t1 := Duration{op.T1}
	t1.Marshal(buf)
	t2 := Duration{op.T2}
	t2.Marshal(buf)
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
