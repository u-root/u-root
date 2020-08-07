package dhcpv6

import (
	"net"
	"testing"

	"github.com/insomniacslk/dhcp/iana"
	"github.com/stretchr/testify/require"
)

func TestParseOptServerID(t *testing.T) {
	data := []byte{
		0, 3, // DUID_LL
		0, 1, // hwtype ethernet
		0, 1, 2, 3, 4, 5, // hw addr
	}
	opt, err := parseOptServerID(data)
	require.NoError(t, err)
	require.Equal(t, DUID_LL, opt.Type)
	require.Equal(t, iana.HWTypeEthernet, opt.HwType)
	require.Equal(t, net.HardwareAddr{0, 1, 2, 3, 4, 5}, opt.LinkLayerAddr)
}

func TestOptServerIdToBytes(t *testing.T) {
	opt := OptServerID(
		Duid{
			Type:          DUID_LL,
			HwType:        iana.HWTypeEthernet,
			LinkLayerAddr: net.HardwareAddr{5, 4, 3, 2, 1, 0},
		},
	)
	expected := []byte{
		0, 3, // DUID_LL
		0, 1, // hwtype ethernet
		5, 4, 3, 2, 1, 0, // hw addr
	}
	require.Equal(t, expected, opt.ToBytes())
}

func TestOptServerIdDecodeEncode(t *testing.T) {
	data := []byte{
		0, 3, // DUID_LL
		0, 1, // hwtype ethernet
		5, 4, 3, 2, 1, 0, // hw addr
	}
	opt, err := parseOptServerID(data)
	require.NoError(t, err)
	require.Equal(t, data, opt.ToBytes())
}

func TestOptionServerId(t *testing.T) {
	opt := OptServerID(
		Duid{
			Type:          DUID_LL,
			HwType:        iana.HWTypeEthernet,
			LinkLayerAddr: net.HardwareAddr{0xde, 0xad, 0, 0, 0xbe, 0xef},
		},
	)
	require.Equal(t, OptionServerID, opt.Code())
	require.Contains(
		t,
		opt.String(),
		"ServerID: DUID{type=DUID-LL hwtype=Ethernet hwaddr=de:ad:00:00:be:ef}",
		"String() should contain the correct sid output",
	)
}

func TestOptServerIdparseOptServerIDBogusDUID(t *testing.T) {
	data := []byte{
		0, 4, // DUID_UUID
		1, 2, 3, 4, 5, 6, 7, 8, 9, // a UUID should be 18 bytes not 17
		10, 11, 12, 13, 14, 15, 16, 17,
	}
	_, err := parseOptServerID(data)
	require.Error(t, err, "A truncated OptServerId DUID should return an error")
}

func TestOptServerIdparseOptServerIDInvalidTooShort(t *testing.T) {
	data := []byte{
		0, // truncated: DUIDs are at least 2 bytes
	}
	_, err := parseOptServerID(data)
	require.Error(t, err, "A truncated OptServerId should return an error")
}
