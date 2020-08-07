package dhcpv6

import (
	"net"
	"testing"

	"github.com/insomniacslk/dhcp/iana"
	"github.com/stretchr/testify/require"
)

func TestParseOptClientID(t *testing.T) {
	data := []byte{
		0, 3, // DUID_LL
		0, 1, // hwtype ethernet
		0, 1, 2, 3, 4, 5, // hw addr
	}
	opt, err := parseOptClientID(data)
	require.NoError(t, err)
	require.Equal(t, DUID_LL, opt.Type)
	require.Equal(t, iana.HWTypeEthernet, opt.HwType)
	require.Equal(t, net.HardwareAddr([]byte{0, 1, 2, 3, 4, 5}), opt.LinkLayerAddr)
}

func TestOptClientIdToBytes(t *testing.T) {
	opt := OptClientID(
		Duid{
			Type:          DUID_LL,
			HwType:        iana.HWTypeEthernet,
			LinkLayerAddr: net.HardwareAddr([]byte{5, 4, 3, 2, 1, 0}),
		},
	)
	expected := []byte{
		0, 3, // DUID_LL
		0, 1, // hwtype ethernet
		5, 4, 3, 2, 1, 0, // hw addr
	}
	require.Equal(t, expected, opt.ToBytes())
}

func TestOptClientIdDecodeEncode(t *testing.T) {
	data := []byte{
		0, 3, // DUID_LL
		0, 1, // hwtype ethernet
		5, 4, 3, 2, 1, 0, // hw addr
	}
	opt, err := parseOptClientID(data)
	require.NoError(t, err)
	require.Equal(t, data, opt.ToBytes())
}

func TestOptionClientId(t *testing.T) {
	opt := OptClientID(
		Duid{
			Type:          DUID_LL,
			HwType:        iana.HWTypeEthernet,
			LinkLayerAddr: net.HardwareAddr([]byte{0xde, 0xad, 0, 0, 0xbe, 0xef}),
		},
	)
	require.Equal(t, OptionClientID, opt.Code())
	require.Contains(
		t,
		opt.String(),
		"ClientID: DUID{type=DUID-LL hwtype=Ethernet hwaddr=de:ad:00:00:be:ef}",
		"String() should contain the correct cid output",
	)
}

func TestOptClientIdparseOptClientIDBogusDUID(t *testing.T) {
	data := []byte{
		0, 4, // DUID_UUID
		1, 2, 3, 4, 5, 6, 7, 8, 9, // a UUID should be 18 bytes not 17
		10, 11, 12, 13, 14, 15, 16, 17,
	}
	_, err := parseOptClientID(data)
	require.Error(t, err, "A truncated OptClientId DUID should return an error")
}

func TestOptClientIdparseOptClientIDInvalidTooShort(t *testing.T) {
	data := []byte{
		0, // truncated: DUIDs are at least 2 bytes
	}
	_, err := parseOptClientID(data)
	require.Error(t, err, "A truncated OptClientId should return an error")
}
