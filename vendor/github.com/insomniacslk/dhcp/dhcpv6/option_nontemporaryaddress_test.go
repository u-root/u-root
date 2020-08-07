package dhcpv6

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestOptIANAParseOptIANA(t *testing.T) {
	data := []byte{
		1, 0, 0, 0, // IAID
		0, 0, 0, 1, // T1
		0, 0, 0, 2, // T2
		0, 5, 0, 0x18, 0x24, 1, 0xdb, 0, 0x30, 0x10, 0xc0, 0x8f, 0xfa, 0xce, 0, 0, 0, 0x44, 0, 0, 0, 0, 0xb2, 0x7a, 0, 0, 0xc0, 0x8a, // options
	}
	opt, err := ParseOptIANA(data)
	require.NoError(t, err)
	require.Equal(t, OptionIANA, opt.Code())
}

func TestOptIANAParseOptIANAInvalidLength(t *testing.T) {
	data := []byte{
		1, 0, 0, 0, // IAID
		0, 0, 0, 1, // T1
		// truncated from here
	}
	_, err := ParseOptIANA(data)
	require.Error(t, err)
}

func TestOptIANAParseOptIANAInvalidOptions(t *testing.T) {
	data := []byte{
		1, 0, 0, 0, // IAID
		0, 0, 0, 1, // T1
		0, 0, 0, 2, // T2
		0, 5, 0, 0x18, 0x24, 1, 0xdb, 0, 0x30, 0x10, 0xc0, 0x8f, 0xfa, 0xce, 0, 0, 0, 0x44, 0, 0, 0, 0, 0xb2, 0x7a, // truncated options
	}
	_, err := ParseOptIANA(data)
	require.Error(t, err)
}

func TestOptIANAGetOneOption(t *testing.T) {
	oaddr := &OptIAAddress{
		IPv6Addr: net.ParseIP("::1"),
	}
	opt := OptIANA{
		Options: IdentityOptions{[]Option{&OptStatusCode{}, oaddr}},
	}
	require.Equal(t, oaddr, opt.Options.OneAddress())
}

func TestOptIANAAddOption(t *testing.T) {
	opt := OptIANA{}
	opt.Options.Add(OptElapsedTime(0))
	require.Equal(t, 1, len(opt.Options.Options))
	require.Equal(t, OptionElapsedTime, opt.Options.Options[0].Code())
}

func TestOptIANAGetOneOptionMissingOpt(t *testing.T) {
	oaddr := &OptIAAddress{
		IPv6Addr: net.ParseIP("::1"),
	}
	opt := OptIANA{
		Options: IdentityOptions{[]Option{&OptStatusCode{}, oaddr}},
	}
	require.Equal(t, nil, opt.Options.GetOne(OptionDNSRecursiveNameServer))
}

func TestOptIANADelOption(t *testing.T) {
	optiaaddr := OptIAAddress{}
	optsc := OptStatusCode{}

	iana1 := OptIANA{
		Options: IdentityOptions{[]Option{
			&optsc,
			&optiaaddr,
			&optiaaddr,
		}},
	}
	iana1.Options.Del(OptionIAAddr)
	require.Equal(t, iana1.Options.Options, Options{&optsc})

	iana2 := OptIANA{
		Options: IdentityOptions{[]Option{
			&optiaaddr,
			&optsc,
			&optiaaddr,
		}},
	}
	iana2.Options.Del(OptionIAAddr)
	require.Equal(t, iana2.Options.Options, Options{&optsc})
}

func TestOptIANAToBytes(t *testing.T) {
	opt := OptIANA{
		IaId: [4]byte{1, 2, 3, 4},
		T1:   12345 * time.Second,
		T2:   54321 * time.Second,
		Options: IdentityOptions{[]Option{
			OptElapsedTime(10 * time.Millisecond),
		}},
	}
	expected := []byte{
		1, 2, 3, 4, // IA ID
		0, 0, 0x30, 0x39, // T1 = 12345
		0, 0, 0xd4, 0x31, // T2 = 54321
		0, 8, 0, 2, 0x00, 0x01,
	}
	require.Equal(t, expected, opt.ToBytes())
}

func TestOptIANAString(t *testing.T) {
	data := []byte{
		1, 0, 0, 0, // IAID
		0, 0, 0, 1, // T1
		0, 0, 0, 2, // T2
		0, 5, 0, 0x18, 0x24, 1, 0xdb, 0, 0x30, 0x10, 0xc0, 0x8f, 0xfa, 0xce, 0, 0, 0, 0x44, 0, 0, 0, 0, 0xb2, 0x7a, 0, 0, 0xc0, 0x8a, // options
	}
	opt, err := ParseOptIANA(data)
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
		"options={",
		"String() should return a list of options",
	)
}
