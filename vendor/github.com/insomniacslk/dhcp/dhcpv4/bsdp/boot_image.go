package bsdp

import (
	"fmt"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/u-root/u-root/pkg/uio"
)

// BootImageType represents the different BSDP boot image types.
type BootImageType byte

// Different types of BootImages - e.g. for different flavors of macOS.
const (
	BootImageTypeMacOS9              BootImageType = 0
	BootImageTypeMacOSX              BootImageType = 1
	BootImageTypeMacOSXServer        BootImageType = 2
	BootImageTypeHardwareDiagnostics BootImageType = 3
	// 4 - 127 are reserved for future use.
)

// bootImageTypeToString maps the different BootImageTypes to human-readable
// representations.
var bootImageTypeToString = map[BootImageType]string{
	BootImageTypeMacOS9:              "macOS 9",
	BootImageTypeMacOSX:              "macOS",
	BootImageTypeMacOSXServer:        "macOS Server",
	BootImageTypeHardwareDiagnostics: "Hardware Diagnostic",
}

// BootImageID describes a boot image ID - whether it's an install image and
// what kind of boot image (e.g. OS 9, macOS, hardware diagnostics)
type BootImageID struct {
	IsInstall bool
	ImageType BootImageType
	Index     uint16
}

// ToBytes implements dhcpv4.OptionValue.
func (b BootImageID) ToBytes() []byte {
	return uio.ToBigEndian(b)
}

// FromBytes reads data into b.
func (b *BootImageID) FromBytes(data []byte) error {
	return uio.FromBigEndian(b, data)
}

// Marshal writes the binary representation to buf.
func (b BootImageID) Marshal(buf *uio.Lexer) {
	var byte0 byte
	if b.IsInstall {
		byte0 |= 0x80
	}
	byte0 |= byte(b.ImageType)
	buf.Write8(byte0)
	buf.Write8(byte(0))
	buf.Write16(b.Index)
}

// String converts a BootImageID to a human-readable representation.
func (b BootImageID) String() string {
	s := fmt.Sprintf("[%d]", b.Index)
	if b.IsInstall {
		s += " installable"
	} else {
		s += " uninstallable"
	}
	t, ok := bootImageTypeToString[b.ImageType]
	if !ok {
		t = "unknown"
	}
	return s + " " + t + " image"
}

// Unmarshal reads b's binary representation from buf.
func (b *BootImageID) Unmarshal(buf *uio.Lexer) error {
	byte0 := buf.Read8()
	_ = buf.Read8()
	b.IsInstall = byte0&0x80 != 0
	b.ImageType = BootImageType(byte0 & 0x7f)
	b.Index = buf.Read16()
	return buf.Error()
}

// BootImage describes a boot image - contains the boot image ID and the name.
type BootImage struct {
	ID   BootImageID
	Name string
}

// Marshal write a BootImage to buf.
func (b BootImage) Marshal(buf *uio.Lexer) {
	b.ID.Marshal(buf)
	buf.Write8(uint8(len(b.Name)))
	buf.WriteBytes([]byte(b.Name))
}

// String converts a BootImage to a human-readable representation.
func (b BootImage) String() string {
	return fmt.Sprintf("%v %v", b.Name, b.ID.String())
}

// Unmarshal reads data from buf into b.
func (b *BootImage) Unmarshal(buf *uio.Lexer) error {
	if err := (&b.ID).Unmarshal(buf); err != nil {
		return err
	}
	nameLength := buf.Read8()
	b.Name = string(buf.Consume(int(nameLength)))
	return buf.Error()
}

func getBootImageID(code dhcpv4.OptionCode, o dhcpv4.Options) *BootImageID {
	v := o.Get(code)
	if v == nil {
		return nil
	}
	var b BootImageID
	if err := uio.FromBigEndian(&b, v); err != nil {
		return nil
	}
	return &b
}

// OptDefaultBootImageID returns a new default boot image ID option as per
// BSDP.
func OptDefaultBootImageID(b BootImageID) dhcpv4.Option {
	return dhcpv4.Option{Code: OptionDefaultBootImageID, Value: b}
}

// OptSelectedBootImageID returns a new selected boot image ID option as per
// BSDP.
func OptSelectedBootImageID(b BootImageID) dhcpv4.Option {
	return dhcpv4.Option{Code: OptionSelectedBootImageID, Value: b}
}
