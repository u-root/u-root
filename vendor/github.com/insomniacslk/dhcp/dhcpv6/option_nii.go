package dhcpv6

import (
	"fmt"

	"github.com/u-root/u-root/pkg/uio"
)

// see rfc4578
const (
	NII_LANDESK_NOPXE   = 0
	NII_PXE_GEN_I       = 1
	NII_PXE_GEN_II      = 2
	NII_UNDI_NOEFI      = 3
	NII_UNDI_EFI_GEN_I  = 4
	NII_UNDI_EFI_GEN_II = 5
)

var niiToStringMap = map[uint8]string{
	NII_LANDESK_NOPXE:   "LANDesk service agent boot ROMs. No PXE",
	NII_PXE_GEN_I:       "First gen. PXE boot ROMs",
	NII_PXE_GEN_II:      "Second gen. PXE boot ROMs",
	NII_UNDI_NOEFI:      "UNDI 32/64 bit. UEFI drivers, no UEFI runtime",
	NII_UNDI_EFI_GEN_I:  "UNDI 32/64 bit. UEFI runtime 1st gen",
	NII_UNDI_EFI_GEN_II: "UNDI 32/64 bit. UEFI runtime 2nd gen",
}

// OptNetworkInterfaceId implements the NIC ID option for network booting as
// defined by RFC 4578 Section 2.2 and RFC 5970 Section 3.4.
type OptNetworkInterfaceId struct {
	type_        uint8
	major, minor uint8 // revision number
}

func (op *OptNetworkInterfaceId) Code() OptionCode {
	return OptionNII
}

func (op *OptNetworkInterfaceId) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	buf.Write8(op.type_)
	buf.Write8(op.major)
	buf.Write8(op.minor)
	return buf.Data()
}

func (op *OptNetworkInterfaceId) Type() uint8 {
	return op.type_
}

func (op *OptNetworkInterfaceId) SetType(type_ uint8) {
	op.type_ = type_
}

func (op *OptNetworkInterfaceId) Major() uint8 {
	return op.major
}

func (op *OptNetworkInterfaceId) SetMajor(major uint8) {
	op.major = major
}

func (op *OptNetworkInterfaceId) Minor() uint8 {
	return op.minor
}

func (op *OptNetworkInterfaceId) SetMinor(minor uint8) {
	op.minor = minor
}

func (op *OptNetworkInterfaceId) String() string {
	typeName, ok := niiToStringMap[op.type_]
	if !ok {
		typeName = "Unknown"
	}
	return fmt.Sprintf("OptNetworkInterfaceId{type=%v, revision=%v.%v}",
		typeName, op.major, op.minor,
	)
}

// build an OptNetworkInterfaceId structure from a sequence of bytes.
// The input data does not include option code and length bytes.
func ParseOptNetworkInterfaceId(data []byte) (*OptNetworkInterfaceId, error) {
	buf := uio.NewBigEndianBuffer(data)
	var opt OptNetworkInterfaceId
	opt.type_ = buf.Read8()
	opt.major = buf.Read8()
	opt.minor = buf.Read8()
	return &opt, buf.FinError()
}
