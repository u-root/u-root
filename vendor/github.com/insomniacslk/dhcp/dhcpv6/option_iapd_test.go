package dhcpv6

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestOptIAPDParseOptIAPD(t *testing.T) {
	data := []byte{
		1, 0, 0, 0, // IAID
		0, 0, 0, 1, // T1
		0, 0, 0, 2, // T2
		0, 26, 0, 25, // 26 = IAPrefix Option, 25 = length
		0xaa, 0xbb, 0xcc, 0xdd, // IAPrefix preferredLifetime
		0xee, 0xff, 0x00, 0x11, // IAPrefix validLifetime
		36,                                             // IAPrefix prefixLength
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, // IAPrefix ipv6Prefix
	}
	opt, err := ParseOptIAPD(data)
	require.NoError(t, err)
	require.Equal(t, OptionIAPD, opt.Code())
	require.Equal(t, [4]byte{1, 0, 0, 0}, opt.IaId)
	require.Equal(t, time.Second, opt.T1)
	require.Equal(t, 2*time.Second, opt.T2)
}

func TestOptIAPDParseOptIAPDInvalidLength(t *testing.T) {
	data := []byte{
		1, 0, 0, 0, // IAID
		0, 0, 0, 1, // T1
		// truncated from here
	}
	_, err := ParseOptIAPD(data)
	require.Error(t, err)
}

func TestOptIAPDParseOptIAPDInvalidOptions(t *testing.T) {
	data := []byte{
		1, 0, 0, 0, // IAID
		0, 0, 0, 1, // T1
		0, 0, 0, 2, // T2
		0, 26, 0, 25, // 26 = IAPrefix Option, 25 = length
		0xaa, 0xbb, 0xcc, 0xdd, // IAPrefix preferredLifetime
		0xee, 0xff, 0x00, 0x11, // IAPrefix validLifetime
		36,                                          // IAPrefix prefixLength
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // IAPrefix ipv6Prefix missing last byte
	}
	_, err := ParseOptIAPD(data)
	require.Error(t, err)
}

func TestOptIAPDToBytes(t *testing.T) {
	oaddr := OptIAPrefix{
		PreferredLifetime: 0xaabbccdd * time.Second,
		ValidLifetime:     0xeeff0011 * time.Second,
		Prefix: &net.IPNet{
			Mask: net.CIDRMask(36, 128),
			IP:   net.IPv6loopback,
		},
	}
	opt := OptIAPD{
		IaId:    [4]byte{1, 2, 3, 4},
		T1:      12345 * time.Second,
		T2:      54321 * time.Second,
		Options: PDOptions{[]Option{&oaddr}},
	}

	expected := []byte{
		1, 2, 3, 4, // IA ID
		0, 0, 0x30, 0x39, // T1 = 12345
		0, 0, 0xd4, 0x31, // T2 = 54321
		0, 26, 0, 25, // 26 = IAPrefix Option, 25 = length
		0xaa, 0xbb, 0xcc, 0xdd, // IAPrefix preferredLifetime
		0xee, 0xff, 0x00, 0x11, // IAPrefix validLifetime
		36,                                             // IAPrefix prefixLength
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, // IAPrefix ipv6Prefix
	}
	require.Equal(t, expected, opt.ToBytes())
}

func TestOptIAPDString(t *testing.T) {
	data := []byte{
		1, 0, 0, 0, // IAID
		0, 0, 0, 1, // T1
		0, 0, 0, 2, // T2
		0, 26, 0, 25, // 26 = IAPrefix Option, 25 = length
		0xaa, 0xbb, 0xcc, 0xdd, // IAPrefix preferredLifetime
		0xee, 0xff, 0x00, 0x11, // IAPrefix validLifetime
		36,                                             // IAPrefix prefixLength
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, // IAPrefix ipv6Prefix
	}
	opt, err := ParseOptIAPD(data)
	require.NoError(t, err)

	str := opt.String()
	require.Contains(
		t, str,
		"IAID=[1 0 0 0]",
		"String() should return the IAID",
	)
	require.Contains(
		t, str,
		"t1=1s, t2=2s",
		"String() should return the T1/T2 options",
	)
	require.Contains(
		t, str,
		"Options=[",
		"String() should return a list of options",
	)
}
