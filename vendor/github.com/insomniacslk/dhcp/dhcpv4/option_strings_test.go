package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseStringsMultiple(t *testing.T) {
	var opt Strings
	expected := []byte{
		9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
		4, 't', 'e', 's', 't',
	}
	err := opt.FromBytes(expected)
	require.NoError(t, err)
	require.Equal(t, len(opt), 2)
	require.Equal(t, "linuxboot", opt[0])
	require.Equal(t, "test", opt[1])
}

func TestParseStringsNone(t *testing.T) {
	var opt Strings
	expected := []byte{}
	err := opt.FromBytes(expected)
	require.Error(t, err)
}

func TestParseStrings(t *testing.T) {
	var opt Strings
	expected := []byte{
		9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
	}
	err := opt.FromBytes(expected)
	require.NoError(t, err)
	require.Equal(t, 1, len(opt))
	require.Equal(t, "linuxboot", opt[0])
}

func TestParseStringsZeroLength(t *testing.T) {
	var opt Strings
	err := opt.FromBytes([]byte{0, 0})
	require.Error(t, err)
}

func TestOptRFC3004UserClass(t *testing.T) {
	opt := OptRFC3004UserClass(Strings([]string{"linuxboot"}))
	data := opt.Value.ToBytes()
	expected := []byte{
		9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
	}
	require.Equal(t, expected, data)
}

func TestOptRFC3004UserClassMultiple(t *testing.T) {
	opt := OptRFC3004UserClass(
		[]string{
			"linuxboot",
			"test",
		},
	)
	data := opt.Value.ToBytes()
	expected := []byte{
		9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
		4, 't', 'e', 's', 't',
	}
	require.Equal(t, expected, data)
}
