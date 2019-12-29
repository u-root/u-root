package dhcpv4

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseIPs(t *testing.T) {
	var i IPs
	data := []byte{
		192, 168, 0, 10, // DNS #1
		192, 168, 0, 20, // DNS #2
	}
	err := i.FromBytes(data)
	require.NoError(t, err)
	servers := []net.IP{
		net.IP{192, 168, 0, 10},
		net.IP{192, 168, 0, 20},
	}
	require.Equal(t, servers, []net.IP(i))

	// Bad length
	data = []byte{1, 1, 1}
	err = i.FromBytes(data)
	require.Error(t, err, "should get error from bad length")

	// RFC2132 requires that at least one IP is specified for each IP field.
	err = i.FromBytes([]byte{})
	require.Error(t, err)
}

func TestOptDomainNameServer(t *testing.T) {
	o := OptDNS(net.IPv4(192, 168, 0, 1), net.IPv4(192, 168, 0, 10))
	require.Equal(t, OptionDomainNameServer, o.Code)
	require.Equal(t, []byte{192, 168, 0, 1, 192, 168, 0, 10}, o.Value.ToBytes())
	require.Equal(t, "Domain Name Server: 192.168.0.1, 192.168.0.10", o.String())
}

func TestGetDomainNameServer(t *testing.T) {
	ips := []net.IP{
		net.IP{192, 168, 0, 1},
		net.IP{192, 168, 0, 10},
	}
	m, _ := New(WithOption(OptDNS(ips...)))
	require.Equal(t, ips, m.DNS())

	m, _ = New()
	require.Nil(t, m.DNS())
}

func TestOptNTPServers(t *testing.T) {
	o := OptNTPServers(net.IPv4(192, 168, 0, 1), net.IPv4(192, 168, 0, 10))
	require.Equal(t, OptionNTPServers, o.Code)
	require.Equal(t, []byte{192, 168, 0, 1, 192, 168, 0, 10}, o.Value.ToBytes())
	require.Equal(t, "NTP Servers: 192.168.0.1, 192.168.0.10", o.String())
}

func TestGetNTPServers(t *testing.T) {
	ips := []net.IP{
		net.IP{192, 168, 0, 1},
		net.IP{192, 168, 0, 10},
	}
	m, _ := New(WithOption(OptNTPServers(ips...)))
	require.Equal(t, ips, m.NTPServers())

	m, _ = New()
	require.Nil(t, m.NTPServers())
}

func TestOptRouter(t *testing.T) {
	o := OptRouter(net.IPv4(192, 168, 0, 1), net.IPv4(192, 168, 0, 10))
	require.Equal(t, OptionRouter, o.Code)
	require.Equal(t, []byte{192, 168, 0, 1, 192, 168, 0, 10}, o.Value.ToBytes())
	require.Equal(t, "Router: 192.168.0.1, 192.168.0.10", o.String())
}

func TestGetRouter(t *testing.T) {
	ips := []net.IP{
		net.IP{192, 168, 0, 1},
		net.IP{192, 168, 0, 10},
	}
	m, _ := New(WithOption(OptRouter(ips...)))
	require.Equal(t, ips, m.Router())

	m, _ = New()
	require.Nil(t, m.Router())
}
