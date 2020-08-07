package dhcpv6

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestOptIAAddressParse(t *testing.T) {
	ipaddr := []byte{0x24, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	data := append(ipaddr, []byte{
		0xa, 0xb, 0xc, 0xd, // preferred lifetime
		0xe, 0xf, 0x1, 0x2, // valid lifetime
		0, 8, 0, 2, 0xaa, 0xbb, // options
	}...)
	opt, err := ParseOptIAAddress(data)
	require.NoError(t, err)
	require.Equal(t, net.IP(ipaddr), opt.IPv6Addr)
	require.Equal(t, 0x0a0b0c0d*time.Second, opt.PreferredLifetime)
	require.Equal(t, 0x0e0f0102*time.Second, opt.ValidLifetime)
}

func TestOptIAAddressParseInvalidTooShort(t *testing.T) {
	data := []byte{
		0x24, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		0xa, 0xb, 0xc, 0xd, // preferred lifetime
		// truncated here
	}
	_, err := ParseOptIAAddress(data)
	require.Error(t, err)
}

func TestOptIAAddressParseInvalidBrokenOptions(t *testing.T) {
	data := []byte{
		0x24, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		0xa, 0xb, 0xc, 0xd, // preferred lifetime
		0xe, 0xf, 0x1, 0x2, // valid lifetime
		0, 8, 0, 2, 0xaa, // broken options
	}
	_, err := ParseOptIAAddress(data)
	require.Error(t, err)
}

func TestOptIAAddressToBytesDefault(t *testing.T) {
	want := []byte{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // IP
		0, 0, 0, 0, // preferred lifetime
		0, 0, 0, 0, // valid lifetime
	}
	opt := OptIAAddress{}
	require.Equal(t, opt.ToBytes(), want)
}

func TestOptIAAddressToBytes(t *testing.T) {
	ipBytes := []byte{0x24, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	expected := append(ipBytes, []byte{
		0xa, 0xb, 0xc, 0xd, // preferred lifetime
		0xe, 0xf, 0x1, 0x2, // valid lifetime
		0, 8, 0, 2, 0x00, 0x01, // options
	}...)
	opt := OptIAAddress{
		IPv6Addr:          net.IP(ipBytes),
		PreferredLifetime: 0x0a0b0c0d * time.Second,
		ValidLifetime:     0x0e0f0102 * time.Second,
		Options: AddressOptions{[]Option{
			OptElapsedTime(10 * time.Millisecond),
		}},
	}
	require.Equal(t, expected, opt.ToBytes())
}

func TestOptIAAddressString(t *testing.T) {
	ipaddr := []byte{0x24, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	data := append(ipaddr, []byte{
		0x00, 0x00, 0x00, 70, // preferred lifetime
		0x00, 0x00, 0x00, 50, // valid lifetime
		0, 8, 0, 2, 0xaa, 0xbb, // options
	}...)
	opt, err := ParseOptIAAddress(data)
	require.NoError(t, err)

	str := opt.String()
	require.Contains(
		t, str,
		"IP=2401:203:405:607:809:a0b:c0d:e0f",
		"String() should return the ipv6addr",
	)
	require.Contains(
		t, str,
		"PreferredLifetime=1m10s",
		"String() should return the preferredlifetime",
	)
	require.Contains(
		t, str,
		"ValidLifetime=50s",
		"String() should return the validlifetime",
	)
}
