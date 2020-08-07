package bsdp

import (
	"net"
	"testing"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/stretchr/testify/require"
)

func TestOptVendorSpecificInformationInterfaceMethods(t *testing.T) {
	o := OptVendorOptions(
		OptVersion(Version1_1),
		OptMessageType(MessageTypeList),
	)
	require.Equal(t, dhcpv4.OptionVendorSpecificInformation, o.Code, "Code")

	expectedBytes := []byte{
		1, 1, 1, // List option
		2, 2, 1, 1, // Version option
	}
	require.Equal(t, expectedBytes, o.Value.ToBytes(), "ToBytes")
}

func TestOptVendorSpecificInformationString(t *testing.T) {
	o := OptVendorOptions(
		OptMessageType(MessageTypeList),
		OptVersion(Version1_1),
	)
	expectedString := "Vendor Specific Information:\n    BSDP Message Type: LIST\n    BSDP Version: 1.1\n"
	require.Equal(t, expectedString, o.String())

	// Test more complicated string - sub options of sub options.
	o = OptVendorOptions(
		OptMessageType(MessageTypeList),
		OptBootImageList(
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
		),
		OptMachineName("foo"),
		OptServerIdentifier(net.IP{1, 1, 1, 1}),
		OptServerPriority(1234),
		OptReplyPort(1235),
		OptDefaultBootImageID(BootImageID{
			IsInstall: true,
			ImageType: BootImageTypeMacOS9,
			Index:     9009,
		}),
		OptSelectedBootImageID(BootImageID{
			IsInstall: true,
			ImageType: BootImageTypeMacOS9,
			Index:     9009,
		}),
	)
	expectedString = "Vendor Specific Information:\n" +
		"    BSDP Message Type: LIST\n" +
		"    BSDP Server Identifier: 1.1.1.1\n" +
		"    BSDP Server Priority: 1234\n" +
		"    BSDP Reply Port: 1235\n" +
		"    BSDP Default Boot Image ID: [9009] installable macOS 9 image\n" +
		"    BSDP Selected Boot Image ID: [9009] installable macOS 9 image\n" +
		"    BSDP Boot Image List: bsdp-1 [1001] uninstallable macOS image, bsdp-2 [9009] installable macOS 9 image\n" +
		"    BSDP Machine Name: foo\n"
	require.Equal(t, expectedString, o.String())
}
