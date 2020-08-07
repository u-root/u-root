package bsdp

import (
	"strings"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/u-root/u-root/pkg/uio"
)

// BootImageList contains a list of boot images presented by a netboot server.
//
// Implements the BSDP option listing the boot images.
type BootImageList []BootImage

// FromBytes deserializes data into bil.
func (bil *BootImageList) FromBytes(data []byte) error {
	buf := uio.NewBigEndianBuffer(data)

	for buf.Has(5) {
		var image BootImage
		if err := image.Unmarshal(buf); err != nil {
			return err
		}
		*bil = append(*bil, image)
	}
	return nil
}

// ToBytes returns a serialized stream of bytes for this option.
func (bil BootImageList) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	for _, image := range bil {
		image.Marshal(buf)
	}
	return buf.Data()
}

// String returns a human-readable string for this option.
func (bil BootImageList) String() string {
	s := make([]string, 0, len(bil))
	for _, image := range bil {
		s = append(s, image.String())
	}
	return strings.Join(s, ", ")
}

// OptBootImageList returns a new BSDP boot image list.
func OptBootImageList(b ...BootImage) dhcpv4.Option {
	return dhcpv4.Option{
		Code:  OptionBootImageList,
		Value: BootImageList(b),
	}
}
