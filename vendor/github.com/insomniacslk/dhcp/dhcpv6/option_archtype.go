package dhcpv6

import (
	"fmt"
	"strings"

	"github.com/insomniacslk/dhcp/iana"
	"github.com/u-root/u-root/pkg/uio"
)

// OptClientArchType represents an option CLIENT_ARCH_TYPE
//
// This module defines the OptClientArchType structure.
// https://www.ietf.org/rfc/rfc5970.txt
type OptClientArchType struct {
	ArchTypes []iana.Arch
}

func (op *OptClientArchType) Code() OptionCode {
	return OptionClientArchType
}

// ToBytes marshals the client arch type as defined by RFC 5970.
func (op *OptClientArchType) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	for _, at := range op.ArchTypes {
		buf.Write16(uint16(at))
	}
	return buf.Data()
}

func (op *OptClientArchType) String() string {
	atStrings := make([]string, 0)
	for _, at := range op.ArchTypes {
		atStrings = append(atStrings, at.String())
	}
	return fmt.Sprintf("OptClientArchType{archtype=%v}", strings.Join(atStrings, ", "))
}

// ParseOptClientArchType builds an OptClientArchType structure from
// a sequence of bytes The input data does not include option code and
// length bytes.
func ParseOptClientArchType(data []byte) (*OptClientArchType, error) {
	var opt OptClientArchType
	buf := uio.NewBigEndianBuffer(data)
	for buf.Has(2) {
		opt.ArchTypes = append(opt.ArchTypes, iana.Arch(buf.Read16()))
	}
	return &opt, buf.FinError()
}
