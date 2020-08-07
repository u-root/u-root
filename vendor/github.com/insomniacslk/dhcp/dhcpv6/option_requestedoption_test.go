package dhcpv6

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptRequestedOption(t *testing.T) {
	expected := []byte{0, 1, 0, 2}
	var o optRequestedOption
	err := o.FromBytes(expected)
	require.NoError(t, err, "ParseOptRequestedOption() correct options should not error")
}

func TestOptRequestedOptionParseOptRequestedOptionTooShort(t *testing.T) {
	buf := []byte{0, 1, 0}
	var o optRequestedOption
	err := o.FromBytes(buf)
	require.Error(t, err, "A short option should return an error (must be divisible by 2)")
}

func TestOptRequestedOptionString(t *testing.T) {
	buf := []byte{0, 1, 0, 2}
	var o optRequestedOption
	err := o.FromBytes(buf)
	require.NoError(t, err)
	require.Contains(
		t,
		o.String(),
		"Client Identifier, Server Identifier",
		"String() should contain the options specified",
	)
	o.OptionCodes = append(o.OptionCodes, 12345)
	require.Contains(
		t,
		o.String(),
		"unknown",
		"String() should contain 'Unknown' for an illegal option",
	)
}
