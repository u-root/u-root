package dhcpv6

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestOptIATAParseOptIATA(t *testing.T) {
	data := []byte{
		1, 0, 0, 0, // IAID
		0, 5, 0, 0x18, 0x24, 1, 0xdb, 0, 0x30, 0x10, 0xc0, 0x8f, 0xfa, 0xce, 0, 0, 0, 0x44, 0, 0, 0, 0, 0xb2, 0x7a, 0, 0, 0xc0, 0x8a, // options
	}
	opt, err := ParseOptIATA(data)
	require.NoError(t, err)
	require.Equal(t, OptionIATA, opt.Code())
}

func TestOptIATAParseOptIATAInvalidLength(t *testing.T) {
	data := []byte{
		1, 0, 0, // truncated IAID
	}
	_, err := ParseOptIATA(data)
	require.Error(t, err)
}

func TestOptIATAParseOptIATAInvalidOptions(t *testing.T) {
	data := []byte{
		1, 0, 0, 0, // IAID
		0, 5, 0, 0x18, 0x24, 1, 0xdb, 0, 0x30, 0x10, 0xc0, 0x8f, 0xfa, 0xce, 0, 0, 0, 0x44, 0, 0, 0, 0, 0xb2, 0x7a, // truncated options
	}
	_, err := ParseOptIATA(data)
	require.Error(t, err)
}

func TestOptIATAGetOneOption(t *testing.T) {
	oaddr := &OptIAAddress{
		IPv6Addr: net.ParseIP("::1"),
	}
	opt := OptIATA{
		Options: IdentityOptions{[]Option{&OptStatusCode{}, oaddr}},
	}
	require.Equal(t, oaddr, opt.Options.OneAddress())
}

func TestOptIATAAddOption(t *testing.T) {
	opt := OptIATA{}
	opt.Options.Add(OptElapsedTime(0))
	require.Equal(t, 1, len(opt.Options.Options))
	require.Equal(t, OptionElapsedTime, opt.Options.Options[0].Code())
}

func TestOptIATAGetOneOptionMissingOpt(t *testing.T) {
	oaddr := &OptIAAddress{
		IPv6Addr: net.ParseIP("::1"),
	}
	opt := OptIATA{
		Options: IdentityOptions{[]Option{&OptStatusCode{}, oaddr}},
	}
	require.Equal(t, nil, opt.Options.GetOne(OptionDNSRecursiveNameServer))
}

func TestOptIATADelOption(t *testing.T) {
	optiaaddr := OptIAAddress{}
	optsc := OptStatusCode{}

	iana1 := OptIATA{
		Options: IdentityOptions{[]Option{
			&optsc,
			&optiaaddr,
			&optiaaddr,
		}},
	}
	iana1.Options.Del(OptionIAAddr)
	require.Equal(t, iana1.Options.Options, Options{&optsc})

	iana2 := OptIATA{
		Options: IdentityOptions{[]Option{
			&optiaaddr,
			&optsc,
			&optiaaddr,
		}},
	}
	iana2.Options.Del(OptionIAAddr)
	require.Equal(t, iana2.Options.Options, Options{&optsc})
}

func TestOptIATAToBytes(t *testing.T) {
	opt := OptIATA{
		IaId: [4]byte{1, 2, 3, 4},
		Options: IdentityOptions{[]Option{
			OptElapsedTime(10 * time.Millisecond),
		}},
	}
	expected := []byte{
		1, 2, 3, 4, // IA ID
		0, 8, 0, 2, 0x00, 0x01,
	}
	require.Equal(t, expected, opt.ToBytes())
}

func TestOptIATAString(t *testing.T) {
	data := []byte{
		1, 0, 0, 0, // IAID
		0, 5, 0, 0x18, 0x24, 1, 0xdb, 0, 0x30, 0x10, 0xc0, 0x8f, 0xfa, 0xce, 0, 0, 0, 0x44, 0, 0, 0, 0, 0xb2, 0x7a, 0, 0, 0xc0, 0x8a, // options
	}
	opt, err := ParseOptIATA(data)
	require.NoError(t, err)

	str := opt.String()
	require.Contains(
		t, str,
		"IAID=[1 0 0 0]",
		"String() should return the IAID",
	)
	require.Contains(
		t, str,
		"options={",
		"String() should return a list of options",
	)
}
