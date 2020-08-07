package dhcpv6

import (
	"bytes"
	"net"
	"testing"

	"github.com/insomniacslk/dhcp/iana"
	"github.com/stretchr/testify/require"
)

func TestParseOptClientLinkLayerAddress(t *testing.T) {
	data := []byte{
		0, 1, // LinkLayerType
		164, 131, 231, 227, 223, 136,
	}
	opt, err := parseOptClientLinkLayerAddress(data)

	require.NoError(t, err)
	require.Equal(t, OptionClientLinkLayerAddr, opt.Code())
	require.Equal(t, iana.HWTypeEthernet, opt.LinkLayerType)
	require.Equal(t, net.HardwareAddr(data[2:]), opt.LinkLayerAddress)
	require.Equal(t, "ClientLinkLayerAddress: Type=Ethernet LinkLayerAddress=a4:83:e7:e3:df:88", opt.String())
}

func TestOptClientLinkLayerAddressToBytes(t *testing.T) {
	mac, _ := net.ParseMAC("a4:83:e7:e3:df:88")
	opt := optClientLinkLayerAddress{
		LinkLayerType:    iana.HWTypeEthernet,
		LinkLayerAddress: mac,
	}
	want := []byte{
		0, 1, // LinkLayerType
		164, 131, 231, 227, 223, 136,
	}
	b := opt.ToBytes()
	if !bytes.Equal(b, want) {
		t.Fatalf("opt.ToBytes()=%v, want %v", b, want)
	}
}
