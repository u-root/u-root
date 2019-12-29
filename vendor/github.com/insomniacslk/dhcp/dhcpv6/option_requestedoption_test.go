package dhcpv6

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptRequestedOption(t *testing.T) {
	expected := []byte{0, 1, 0, 2}
	_, err := ParseOptRequestedOption(expected)
	require.NoError(t, err, "ParseOptRequestedOption() correct options should not error")
}

func TestOptRequestedOptionParseOptRequestedOptionTooShort(t *testing.T) {
	buf := []byte{0, 1, 0}
	_, err := ParseOptRequestedOption(buf)
	require.Error(t, err, "A short option should return an error (must be divisible by 2)")
}

func TestOptRequestedOptionString(t *testing.T) {
	buf := []byte{0, 1, 0, 2}
	opt, err := ParseOptRequestedOption(buf)
	require.NoError(t, err)
	require.Contains(
		t,
		opt.String(),
		"OPTION_CLIENTID, OPTION_SERVERID",
		"String() should contain the options specified",
	)
	opt.AddRequestedOption(12345)
	require.Contains(
		t,
		opt.String(),
		"unknown",
		"String() should contain 'Unknown' for an illegal option",
	)
}
