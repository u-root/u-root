package dhcpv6

import (
	"net"
	"testing"

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
		Options: []Option{&OptElapsedTime{}, oaddr},
	}
	require.Equal(t, oaddr, opt.GetOneOption(OptionIAAddr))
}

func TestOptIANAAddOption(t *testing.T) {
	opt := OptIANA{}
	opt.AddOption(&OptElapsedTime{})
	require.Equal(t, 1, len(opt.Options))
	require.Equal(t, OptionElapsedTime, opt.Options[0].Code())
}

func TestOptIANAGetOneOptionMissingOpt(t *testing.T) {
	oaddr := &OptIAAddress{
		IPv6Addr: net.ParseIP("::1"),
	}
	opt := OptIANA{
		Options: []Option{&OptElapsedTime{}, oaddr},
	}
	require.Equal(t, nil, opt.GetOneOption(OptionDNSRecursiveNameServer))
}

func TestOptIANADelOption(t *testing.T) {
	optiana1 := OptIANA{}
	optiana2 := OptIANA{}
	optiaaddr := OptIAAddress{}
	optsc := OptStatusCode{}

	optiana1.Options = append(optiana1.Options, &optsc)
	optiana1.Options = append(optiana1.Options, &optiaaddr)
	optiana1.Options = append(optiana1.Options, &optiaaddr)
	optiana1.DelOption(OptionIAAddr)
	require.Equal(t, optiana1.Options, Options{&optsc})

	optiana2.Options = append(optiana2.Options, &optiaaddr)
	optiana2.Options = append(optiana2.Options, &optsc)
	optiana2.Options = append(optiana2.Options, &optiaaddr)
	optiana2.DelOption(OptionIAAddr)
	require.Equal(t, optiana2.Options, Options{&optsc})
}

func TestOptIANAToBytes(t *testing.T) {
	opt := OptIANA{
		IaId: [4]byte{1, 2, 3, 4},
		T1:   12345,
		T2:   54321,
		Options: []Option{
			&OptElapsedTime{
				ElapsedTime: 0xaabb,
			},
		},
	}
	expected := []byte{
		1, 2, 3, 4, // IA ID
		0, 0, 0x30, 0x39, // T1 = 12345
		0, 0, 0xd4, 0x31, // T2 = 54321
		0, 8, 0, 2, 0xaa, 0xbb,
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
		"t1=1, t2=2",
		"String() should return the T1/T2 options",
	)
	require.Contains(
		t, str,
		"options=[",
		"String() should return a list of options",
	)
}
