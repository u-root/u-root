package bsdp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptBootImageListInterfaceMethods(t *testing.T) {
	bs := []BootImage{
		BootImage{
			ID: BootImageID{
				IsInstall: false,
				ImageType: BootImageTypeMacOSX,
				Index:     1001,
			},
			Name: "bsdp-1",
		},
		BootImage{
			ID: BootImageID{
				IsInstall: true,
				ImageType: BootImageTypeMacOS9,
				Index:     9009,
			},
			Name: "bsdp-2",
		},
	}
	o := OptBootImageList(bs...)
	require.Equal(t, OptionBootImageList, o.Code, "Code")
	expectedBytes := []byte{
		// boot image 1
		0x1, 0x0, 0x03, 0xe9, // ID
		6, // name length
		'b', 's', 'd', 'p', '-', '1',
		// boot image 1
		0x80, 0x0, 0x23, 0x31, // ID
		6, // name length
		'b', 's', 'd', 'p', '-', '2',
	}
	require.Equal(t, expectedBytes, o.Value.ToBytes(), "ToBytes")
}

func TestParseOptBootImageList(t *testing.T) {
	data := []byte{
		// boot image 1
		0x1, 0x0, 0x03, 0xe9, // ID
		6, // name length
		'b', 's', 'd', 'p', '-', '1',
		// boot image 1
		0x80, 0x0, 0x23, 0x31, // ID
		6, // name length
		'b', 's', 'd', 'p', '-', '2',
	}
	var o BootImageList
	err := o.FromBytes(data)
	require.NoError(t, err)
	expectedBootImages := BootImageList{
		BootImage{
			ID: BootImageID{
				IsInstall: false,
				ImageType: BootImageTypeMacOSX,
				Index:     1001,
			},
			Name: "bsdp-1",
		},
		BootImage{
			ID: BootImageID{
				IsInstall: true,
				ImageType: BootImageTypeMacOS9,
				Index:     9009,
			},
			Name: "bsdp-2",
		},
	}
	require.Equal(t, expectedBootImages, o)

	// Error parsing boot image (malformed)
	data = []byte{
		// boot image 1
		0x1, 0x0, 0x03, 0xe9, // ID
		4, // name length
		'b', 's', 'd', 'p', '-', '1',
		// boot image 2
		0x80, 0x0, 0x23, 0x31, // ID
		6, // name length
		'b', 's', 'd', 'p', '-', '2',
	}
	err = o.FromBytes(data)
	require.Error(t, err, "should get error from bad boot image")
}

func TestOptBootImageListString(t *testing.T) {
	bs := []BootImage{
		BootImage{
			ID: BootImageID{
				IsInstall: false,
				ImageType: BootImageTypeMacOSX,
				Index:     1001,
			},
			Name: "bsdp-1",
		},
		BootImage{
			ID: BootImageID{
				IsInstall: true,
				ImageType: BootImageTypeMacOS9,
				Index:     9009,
			},
			Name: "bsdp-2",
		},
	}
	o := OptBootImageList(bs...)
	expectedString := "BSDP Boot Image List: bsdp-1 [1001] uninstallable macOS image, bsdp-2 [9009] installable macOS 9 image"
	require.Equal(t, expectedString, o.String())
}
