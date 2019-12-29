package dhcpv6

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptIAForPrefixDelegationParseOptIAForPrefixDelegation(t *testing.T) {
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
	opt, err := ParseOptIAForPrefixDelegation(data)
	require.NoError(t, err)
	require.Equal(t, OptionIAPD, opt.Code())
	require.Equal(t, [4]byte{1, 0, 0, 0}, opt.IaId)
	require.Equal(t, uint32(1), opt.T1)
	require.Equal(t, uint32(2), opt.T2)
}

func TestOptIAForPrefixDelegationParseOptIAForPrefixDelegationInvalidLength(t *testing.T) {
	data := []byte{
		1, 0, 0, 0, // IAID
		0, 0, 0, 1, // T1
		// truncated from here
	}
	_, err := ParseOptIAForPrefixDelegation(data)
	require.Error(t, err)
}

func TestOptIAForPrefixDelegationParseOptIAForPrefixDelegationInvalidOptions(t *testing.T) {
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
	_, err := ParseOptIAForPrefixDelegation(data)
	require.Error(t, err)
}

func TestOptIAForPrefixDelegationGetOneOption(t *testing.T) {
	buf := []byte{
		0xaa, 0xbb, 0xcc, 0xdd, // preferredLifetime
		0xee, 0xff, 0x00, 0x11, // validLifetime
		36,                                             // prefixLength
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, // ipv6Prefix
	}
	oaddr, err := ParseOptIAPrefix(buf)
	if err != nil {
		t.Fatal(err)
	}
	opt := OptIAForPrefixDelegation{}
	opt.Options = append(opt.Options, oaddr)
	require.Equal(t, oaddr, opt.GetOneOption(OptionIAPrefix))
}

func TestOptIAForPrefixDelegationGetOneOptionMissingOpt(t *testing.T) {
	buf := []byte{
		0xaa, 0xbb, 0xcc, 0xdd, // preferredLifetime
		0xee, 0xff, 0x00, 0x11, // validLifetime
		36,                                             // prefixLength
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, // ipv6Prefix
	}
	oaddr, err := ParseOptIAPrefix(buf)
	if err != nil {
		t.Fatal(err)
	}
	opt := OptIAForPrefixDelegation{}
	opt.Options = append(opt.Options, oaddr)
	require.Equal(t, nil, opt.GetOneOption(OptionDNSRecursiveNameServer))
}

func TestOptIAForPrefixDelegationDelOption(t *testing.T) {
	optiana1 := OptIAForPrefixDelegation{}
	optiana2 := OptIAForPrefixDelegation{}
	optiaaddr := OptIAPrefix{}
	optsc := OptStatusCode{}

	optiana1.Options = append(optiana1.Options, &optsc)
	optiana1.Options = append(optiana1.Options, &optiaaddr)
	optiana1.Options = append(optiana1.Options, &optiaaddr)
	optiana1.DelOption(OptionIAPrefix)
	require.Equal(t, len(optiana1.Options), 1)
	require.Equal(t, optiana1.Options[0], &optsc)

	optiana2.Options = append(optiana2.Options, &optiaaddr)
	optiana2.Options = append(optiana2.Options, &optsc)
	optiana2.Options = append(optiana2.Options, &optiaaddr)
	optiana2.DelOption(OptionIAPrefix)
	require.Equal(t, len(optiana2.Options), 1)
	require.Equal(t, optiana2.Options[0], &optsc)
}

func TestOptIAForPrefixDelegationToBytes(t *testing.T) {
	oaddr := OptIAPrefix{}
	oaddr.PreferredLifetime = 0xaabbccdd
	oaddr.ValidLifetime = 0xeeff0011
	oaddr.SetPrefixLength(36)
	oaddr.SetIPv6Prefix(net.IPv6loopback)

	opt := OptIAForPrefixDelegation{}
	opt.IaId = [4]byte{1, 2, 3, 4}
	opt.T1 = 12345
	opt.T2 = 54321
	opt.Options = append(opt.Options, &oaddr)

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

func TestOptIAForPrefixDelegationString(t *testing.T) {
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
	opt, err := ParseOptIAForPrefixDelegation(data)
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
