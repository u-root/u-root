package dhcpv6

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseOptDHCP4oDHCP6Server(t *testing.T) {
	data := []byte{
		0x2a, 0x03, 0x28, 0x80, 0xff, 0xfe, 0x00, 0x0c, 0xfa, 0xce, 0xb0, 0x0c, 0x00, 0x00, 0x00, 0x35,
	}
	expected := []net.IP{
		net.IP(data),
	}
	opt, err := ParseOptDHCP4oDHCP6Server(data)
	require.NoError(t, err)
	require.Equal(t, expected, opt.DHCP4oDHCP6Servers)
	require.Equal(t, OptionDHCP4oDHCP6Server, opt.Code())
	require.Contains(t, opt.String(), "4o6-servers=[2a03:2880:fffe:c:face:b00c:0:35]", "String() should contain the correct DHCP4-over-DHCP6 server output")
}

func TestOptDHCP4oDHCP6ServerToBytes(t *testing.T) {
	ip1 := net.ParseIP("2a03:2880:fffe:c:face:b00c:0:35")
	ip2 := net.ParseIP("2001:4860:4860::8888")
	servers := []net.IP{ip1, ip2}
	expected := append([]byte{}, []byte(ip1)...)
	expected = append(expected, []byte(ip2)...)
	opt := OptDHCP4oDHCP6Server{DHCP4oDHCP6Servers: servers}
	require.Equal(t, expected, opt.ToBytes())
}

func TestParseOptDHCP4oDHCP6ServerParseNoAddr(t *testing.T) {
	data := []byte{
	}
	var expected []net.IP
	opt, err := ParseOptDHCP4oDHCP6Server(data)
	require.NoError(t, err)
	require.Equal(t, expected, opt.DHCP4oDHCP6Servers)
}

func TestOptDHCP4oDHCP6ServerToBytesNoAddr(t *testing.T) {
	expected := []byte(nil)
	opt := OptDHCP4oDHCP6Server{}
	require.Equal(t, expected, opt.ToBytes())
}

func TestParseOptDHCP4oDHCP6ServerParseBogus(t *testing.T) {
	data := []byte{
		0x2a, 0x03, 0x28, 0x80, 0xff, 0xfe, 0x00, 0x0c, // invalid IPv6 address
	}
	_, err := ParseOptDHCP4oDHCP6Server(data)
	require.Error(t, err, "An invalid IPv6 address should return an error")
}
