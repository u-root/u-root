package dhcpv6

import (
	"bytes"
	"net"
	"testing"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/iana"
	"github.com/stretchr/testify/require"
)

var magicCookie = [4]byte{99, 130, 83, 99}

func TestParseOptDHCPv4Msg(t *testing.T) {
	data := []byte{
		1,                      // dhcp request
		1,                      // ethernet hw type
		6,                      // hw addr length
		3,                      // hop count
		0xaa, 0xbb, 0xcc, 0xdd, // transaction ID, big endian (network)
		0, 3, // number of seconds
		0, 1, // broadcast
		0, 0, 0, 0, // client IP address
		0, 0, 0, 0, // your IP address
		0, 0, 0, 0, // server IP address
		0, 0, 0, 0, // gateway IP address
		0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // client MAC address + padding
	}

	// server host name
	expectedHostname := []byte{}
	for i := 0; i < 64; i++ {
		expectedHostname = append(expectedHostname, 0)
	}
	data = append(data, expectedHostname...)
	// boot file name
	expectedBootfilename := []byte{}
	for i := 0; i < 128; i++ {
		expectedBootfilename = append(expectedBootfilename, 0)
	}
	data = append(data, expectedBootfilename...)
	// magic cookie, then no options
	data = append(data, magicCookie[:]...)

	opt, err := ParseOptDHCPv4Msg(data)
	d := opt.Msg
	require.NoError(t, err)
	require.Equal(t, d.OpCode, dhcpv4.OpcodeBootRequest)
	require.Equal(t, d.HWType, iana.HWTypeEthernet)
	require.Equal(t, d.HopCount, byte(3))
	require.Equal(t, d.TransactionID, dhcpv4.TransactionID{0xaa, 0xbb, 0xcc, 0xdd})
	require.Equal(t, d.NumSeconds, uint16(3))
	require.Equal(t, d.Flags, uint16(1))
	require.True(t, d.ClientIPAddr.Equal(net.IPv4zero))
	require.True(t, d.YourIPAddr.Equal(net.IPv4zero))
	require.True(t, d.GatewayIPAddr.Equal(net.IPv4zero))
	require.Equal(t, d.ClientHWAddr, net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff})
	require.Equal(t, d.ServerHostName, "")
	require.Equal(t, d.BootFileName, "")
	// no need to check Magic Cookie as it is already validated in FromBytes
	// above
}

func TestOptDHCPv4MsgToBytes(t *testing.T) {
	// the following bytes match what dhcpv4.New would create. Keep them in
	// sync!
	expected := []byte{
		1,                      // Opcode BootRequest
		1,                      // HwType Ethernet
		6,                      // HwAddrLen
		0,                      // HopCount
		0x11, 0x22, 0x33, 0x44, // TransactionID
		0, 0, // NumSeconds
		0, 0, // Flags
		0, 0, 0, 0, // ClientIPAddr
		0, 0, 0, 0, // YourIPAddr
		0, 0, 0, 0, // ServerIPAddr
		0, 0, 0, 0, // GatewayIPAddr
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // ClientHwAddr
	}
	// ServerHostName
	expected = append(expected, bytes.Repeat([]byte{0}, 64)...)
	// BootFileName
	expected = append(expected, bytes.Repeat([]byte{0}, 128)...)

	// Magic Cookie
	expected = append(expected, magicCookie[:]...)

	// Minimum message length padding.
	//
	// 236 + 4 byte cookie + 59 bytes padding + 1 byte end.
	expected = append(expected, bytes.Repeat([]byte{0}, 59)...)

	// End
	expected = append(expected, 0xff)

	d, err := dhcpv4.New()
	require.NoError(t, err)
	// fix TransactionID to match the expected one, since it's randomly
	// generated in New()
	d.TransactionID = dhcpv4.TransactionID{0x11, 0x22, 0x33, 0x44}
	opt := OptDHCPv4Msg{Msg: d}
	require.Equal(t, expected, opt.ToBytes())
}
