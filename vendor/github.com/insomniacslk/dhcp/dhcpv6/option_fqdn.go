package dhcpv6

import (
	"fmt"

	"github.com/u-root/u-root/pkg/uio"
)

// OptFQDN implements OptionFQDN option.
//
// https://tools.ietf.org/html/rfc4704
type OptFQDN struct {
	Flags      uint8
	DomainName string
}

// Code returns the option code.
func (op *OptFQDN) Code() OptionCode {
	return OptionFQDN
}

// ToBytes serializes the option and returns it as a sequence of bytes
func (op *OptFQDN) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	buf.Write8(op.Flags)
	buf.WriteBytes([]byte(op.DomainName))
	return buf.Data()
}

func (op *OptFQDN) String() string {
	return fmt.Sprintf("OptFQDN{flags=%d, domainname=%s}", op.Flags, op.DomainName)
}

// ParseOptFQDN deserializes from bytes to build a OptFQDN structure.
func ParseOptFQDN(data []byte) (*OptFQDN, error) {
	var opt OptFQDN
	buf := uio.NewBigEndianBuffer(data)
	opt.Flags = buf.Read8()
	opt.DomainName = string(buf.ReadAll())
	return &opt, buf.FinError()
}
