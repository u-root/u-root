package dhcpv6

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseOptUserClass(t *testing.T) {
	expected := []byte{
		0, 9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
	}
	opt, err := ParseOptUserClass(expected)
	require.NoError(t, err)
	require.Equal(t, 1, len(opt.UserClasses))
	require.Equal(t, []byte("linuxboot"), opt.UserClasses[0])
}

func TestParseOptUserClassMultiple(t *testing.T) {
	expected := []byte{
		0, 9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
		0, 4, 't', 'e', 's', 't',
	}
	opt, err := ParseOptUserClass(expected)
	require.NoError(t, err)
	require.Equal(t, len(opt.UserClasses), 2)
	require.Equal(t, []byte("linuxboot"), opt.UserClasses[0])
	require.Equal(t, []byte("test"), opt.UserClasses[1])
}

func TestParseOptUserClassNone(t *testing.T) {
	expected := []byte{}
	_, err := ParseOptUserClass(expected)
	require.Error(t, err)
}

func TestOptUserClassToBytes(t *testing.T) {
	opt := OptUserClass{
		UserClasses: [][]byte{[]byte("linuxboot")},
	}
	data := opt.ToBytes()
	expected := []byte{
		0, 9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
	}
	require.Equal(t, expected, data)
}

func TestOptUserClassToBytesMultiple(t *testing.T) {
	opt := OptUserClass{
		UserClasses: [][]byte{
			[]byte("linuxboot"),
			[]byte("test"),
		},
	}
	data := opt.ToBytes()
	expected := []byte{
		0, 9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
		0, 4, 't', 'e', 's', 't',
	}
	require.Equal(t, expected, data)
}

func TestOptUserClassParseOptUserClassTooShort(t *testing.T) {
	buf := []byte{
		0, 9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
		0, 4, 't', 'e',
	}
	_, err := ParseOptUserClass(buf)
	require.Error(t, err, "ParseOptUserClass() should error if given truncated user classes")

	buf = []byte{
		0, 9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
		0,
	}
	_, err = ParseOptUserClass(buf)
	require.Error(t, err, "ParseOptUserClass() should error if given a truncated length")
}

func TestOptUserClassString(t *testing.T) {
	data := []byte{
		0, 9, 'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
		0, 4, 't', 'e', 's', 't',
	}
	opt, err := ParseOptUserClass(data)
	require.NoError(t, err)

	require.Contains(
		t,
		opt.String(),
		"userclass=[linuxboot, test]",
		"String() should contain the list of user classes",
	)
}
