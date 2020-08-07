package dhcpv4

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptBroadcastAddress(t *testing.T) {
	ip := net.IP{192, 168, 0, 1}
	o := OptBroadcastAddress(ip)

	require.Equal(t, OptionBroadcastAddress, o.Code, "Code")
	require.Equal(t, []byte(ip), o.Value.ToBytes(), "ToBytes")
	require.Equal(t, "Broadcast Address: 192.168.0.1", o.String(), "String")
}

func TestGetIPs(t *testing.T) {
	o := Options{102: []byte{}}
	i := GetIPs(optionCode(102), o)
	require.Nil(t, i)

	o = Options{102: []byte{192, 168, 0}}
	i = GetIPs(optionCode(102), o)
	require.Nil(t, i)

	o = Options{102: []byte{192, 168, 0, 1}}
	i = GetIPs(optionCode(102), o)
	require.Equal(t, i, []net.IP{{192, 168, 0, 1}})

	o = Options{102: []byte{192, 168, 0, 1, 192, 168, 0, 2}}
	i = GetIPs(optionCode(102), o)
	require.Equal(t, i, []net.IP{{192, 168, 0, 1}, {192, 168, 0, 2}})
}

func TestParseIP(t *testing.T) {
	var ip IP
	err := ip.FromBytes([]byte{})
	require.Error(t, err, "empty byte stream")

	err = ip.FromBytes([]byte{192, 168, 0})
	require.Error(t, err, "wrong IP length")

	err = ip.FromBytes([]byte{192, 168, 0, 1})
	require.NoError(t, err)
	require.Equal(t, net.IP{192, 168, 0, 1}, net.IP(ip))
}

func TestOptRequestedIPAddress(t *testing.T) {
	ip := net.IP{192, 168, 0, 1}
	o := OptRequestedIPAddress(ip)

	require.Equal(t, OptionRequestedIPAddress, o.Code, "Code")
	require.Equal(t, []byte(ip), o.Value.ToBytes(), "ToBytes")
	require.Equal(t, "Requested IP Address: 192.168.0.1", o.String(), "String")
}

func TestOptServerIdentifier(t *testing.T) {
	ip := net.IP{192, 168, 0, 1}
	o := OptServerIdentifier(ip)

	require.Equal(t, OptionServerIdentifier, o.Code, "Code")
	require.Equal(t, []byte(ip), o.Value.ToBytes(), "ToBytes")
	require.Equal(t, "Server Identifier: 192.168.0.1", o.String(), "String")
}
