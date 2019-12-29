package dhcpv6

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpt4RDNonMapRuleParse(t *testing.T) {
	data := []byte{0x81, 0xaa, 0x05, 0xd4}
	opt, err := ParseOpt4RDNonMapRule(data)
	require.NoError(t, err)
	require.True(t, opt.HubAndSpoke)
	require.NotNil(t, opt.TrafficClass)
	require.EqualValues(t, 0xaa, *opt.TrafficClass)
	require.EqualValues(t, 1492, opt.DomainPMTU)

	// Remove the TrafficClass flag and check value is ignored
	data[0] = 0x80
	opt, err = ParseOpt4RDNonMapRule(data)
	require.NoError(t, err)
	require.True(t, opt.HubAndSpoke)
	require.Nil(t, opt.TrafficClass)
	require.EqualValues(t, 1492, opt.DomainPMTU)
}

func TestOpt4RDNonMapRuleToBytes(t *testing.T) {
	var tClass uint8 = 0xaa
	opt := Opt4RDNonMapRule{
		HubAndSpoke:  true,
		TrafficClass: &tClass,
		DomainPMTU:   1492,
	}
	expected := []byte{0x81, 0xaa, 0x05, 0xd4}

	require.Equal(t, expected, opt.ToBytes())

	// Unsetting TrafficClass should zero the corresponding bytes in the output
	opt.TrafficClass = nil
	expected[0], expected[1] = 0x80, 0x00

	require.Equal(t, expected, opt.ToBytes())
}

func TestOpt4RDNonMapRuleString(t *testing.T) {
	var tClass uint8 = 120
	opt := Opt4RDNonMapRule{
		HubAndSpoke:  true,
		TrafficClass: &tClass,
		DomainPMTU:   9000,
	}

	str := opt.String()

	require.Contains(t, str, "HubAndSpoke=true",
		"String() should contain the HubAndSpoke flag value")
	require.Contains(t, str, "TrafficClass=120",
		"String() should contain the TrafficClass flag value")
	require.Contains(t, str, "DomainPMTU=9000",
		"String() should contain the domain PMTU")
}

func TestOpt4RDMapRuleParse(t *testing.T) {
	ip6addr, ip6net, err := net.ParseCIDR("2001:db8::1234:5678:0:aabb/64")
	ip6net.IP = ip6addr // We want to keep the entire address however, not apply the mask
	require.NoError(t, err)
	ip4addr, ip4net, err := net.ParseCIDR("100.64.0.234/10")
	ip4net.IP = ip4addr.To4()
	require.NoError(t, err)
	data := append([]byte{
		10,   // IPv4 prefix length
		64,   // IPv6 prefix length
		32,   // EA-bits
		0x80, // WKPs authorized
	},
		append(ip4addr.To4(), ip6addr...)...,
	)

	opt, err := ParseOpt4RDMapRule(data)
	require.NoError(t, err)
	require.EqualValues(t, *ip6net, opt.Prefix6)
	require.EqualValues(t, *ip4net, opt.Prefix4)
	require.EqualValues(t, 32, opt.EABitsLength)
	require.True(t, opt.WKPAuthorized)
}

func TestOpt4RDMapRuleToBytes(t *testing.T) {
	opt := Opt4RDMapRule{
		Prefix4: net.IPNet{
			IP:   net.IPv4(100, 64, 0, 238),
			Mask: net.CIDRMask(24, 32),
		},
		Prefix6: net.IPNet{
			IP:   net.ParseIP("2001:db8::1234:5678:0:aabb"),
			Mask: net.CIDRMask(80, 128),
		},
		EABitsLength:  32,
		WKPAuthorized: true,
	}

	expected := append([]byte{
		24,   // v4 prefix length
		80,   // v6 prefix length
		32,   // EA-bits
		0x80, // WKPs authorized
	},
		append(opt.Prefix4.IP.To4(), opt.Prefix6.IP.To16()...)...,
	)

	require.Equal(t, expected, opt.ToBytes())
}

// FIXME: Invalid packets are serialized without error

func TestOpt4RDMapRuleString(t *testing.T) {
	opt := Opt4RDMapRule{
		Prefix4: net.IPNet{
			IP:   net.IPv4(100, 64, 0, 238),
			Mask: net.CIDRMask(24, 32),
		},
		Prefix6: net.IPNet{
			IP:   net.ParseIP("2001:db8::1234:5678:0:aabb"),
			Mask: net.CIDRMask(80, 128),
		},
		EABitsLength:  32,
		WKPAuthorized: true,
	}

	str := opt.String()
	require.Contains(t, str, "WKPAuthorized=true", "String() should write the flag values")
	require.Contains(t, str, "Prefix6=2001:db8::1234:5678:0:aabb/80",
		"String() should include the IPv6 prefix")
	require.Contains(t, str, "Prefix4=100.64.0.238/24",
		"String() should include the IPv4 prefix")
	require.Contains(t, str, "EA-Bits=32", "String() should include the value for EA-Bits")
}

// This test round-trip serialization/deserialization of both kinds of 4RD
// options, and the container option
func TestOpt4RDRoundTrip(t *testing.T) {
	var tClass uint8 = 0xaa
	opt := Opt4RD{
		&Opt4RDMapRule{
			Prefix4: net.IPNet{
				IP:   net.IPv4(100, 64, 0, 238).To4(),
				Mask: net.CIDRMask(24, 32),
			},
			Prefix6: net.IPNet{
				IP:   net.ParseIP("2001:db8::1234:5678:0:aabb"),
				Mask: net.CIDRMask(80, 128),
			},
			EABitsLength:  32,
			WKPAuthorized: true,
		},
		&Opt4RDNonMapRule{
			HubAndSpoke:  true,
			TrafficClass: &tClass,
			DomainPMTU:   9000,
		},
	}

	rtOpt, err := ParseOpt4RD(opt.ToBytes())

	require.NoError(t, err)
	require.NotNil(t, rtOpt)
	require.Equal(t, opt, *rtOpt)
}
