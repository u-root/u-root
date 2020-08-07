package dhcpv6

import (
	"bytes"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestOptIAPrefix(t *testing.T) {
	buf := []byte{
		0xaa, 0xbb, 0xcc, 0xdd, // preferredLifetime
		0xee, 0xff, 0x00, 0x11, // validLifetime
		36,                                             // prefixLength
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, // ipv6Prefix
	}
	opt, err := ParseOptIAPrefix(buf)
	if err != nil {
		t.Fatal(err)
	}
	want := &OptIAPrefix{
		PreferredLifetime: 0xaabbccdd * time.Second,
		ValidLifetime:     0xeeff0011 * time.Second,
		Prefix: &net.IPNet{
			Mask: net.CIDRMask(36, 128),
			IP:   net.IPv6loopback,
		},
		Options: PrefixOptions{[]Option{}},
	}
	if !reflect.DeepEqual(want, opt) {
		t.Errorf("parseIAPrefix = %v, want %v", opt, want)
	}
}

func TestOptIAPrefixToBytes(t *testing.T) {
	buf := []byte{
		0xaa, 0xbb, 0xcc, 0xdd, // preferredLifetime
		0xee, 0xff, 0x00, 0x11, // validLifetime
		36,                                             // prefixLength
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // ipv6Prefix
		0, 8, 0, 2, 0x00, 0x01, // options
	}
	opt := OptIAPrefix{
		PreferredLifetime: 0xaabbccdd * time.Second,
		ValidLifetime:     0xeeff0011 * time.Second,
		Prefix: &net.IPNet{
			Mask: net.CIDRMask(36, 128),
			IP:   net.IPv6zero,
		},
		Options: PrefixOptions{[]Option{OptElapsedTime(10 * time.Millisecond)}},
	}
	toBytes := opt.ToBytes()
	if !bytes.Equal(toBytes, buf) {
		t.Fatalf("Invalid ToBytes result. Expected %v, got %v", buf, toBytes)
	}
}

func TestOptIAPrefixToBytesDefault(t *testing.T) {
	buf := []byte{
		0, 0, 0, 0, // preferredLifetime
		0, 0, 0, 0, // validLifetime
		0,                                              // prefixLength
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // ipv6Prefix
	}
	opt := OptIAPrefix{}
	toBytes := opt.ToBytes()
	if !bytes.Equal(toBytes, buf) {
		t.Fatalf("Invalid ToBytes result. Expected %v, got %v", buf, toBytes)
	}
}

func TestOptIAPrefixParseInvalidTooShort(t *testing.T) {
	buf := []byte{
		0xaa, 0xbb, 0xcc, 0xdd, // preferredLifetime
		0xee, 0xff, 0x00, 0x11, // validLifetime
		36,                  // prefixLength
		0, 0, 0, 0, 0, 0, 0, // truncated ipv6Prefix
	}
	if opt, err := ParseOptIAPrefix(buf); err == nil {
		t.Fatalf("ParseOptIAPrefix: Expected error on truncated option, got %v", opt)
	}
}

func TestOptIAPrefixString(t *testing.T) {
	buf := []byte{
		0x00, 0x00, 0x00, 60, // preferredLifetime
		0x00, 0x00, 0x00, 50, // validLifetime
		36,                                                         // prefixLength
		0x20, 0x01, 0x0d, 0xb8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // ipv6Prefix
	}
	opt, err := ParseOptIAPrefix(buf)
	require.NoError(t, err)

	str := opt.String()
	require.Contains(
		t, str,
		"Prefix=2001:db8::/36",
		"String() should return the ipv6addr",
	)
	require.Contains(
		t, str,
		"PreferredLifetime=1m",
		"String() should return the preferredlifetime",
	)
	require.Contains(
		t, str,
		"ValidLifetime=50s",
		"String() should return the validlifetime",
	)
}
