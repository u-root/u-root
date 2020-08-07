package dhcpv6

import (
	"testing"

	"github.com/insomniacslk/dhcp/iana"
	"github.com/stretchr/testify/require"
)

func TestParseOptClientArchType(t *testing.T) {
	data := []byte{
		0, 6, // EFI_IA32
	}
	opt, err := parseOptClientArchType(data)
	require.NoError(t, err)
	require.Equal(t, iana.EFI_IA32, opt.Archs[0])
}

func TestParseOptClientArchTypeInvalid(t *testing.T) {
	data := []byte{42}
	_, err := parseOptClientArchType(data)
	require.Error(t, err)
}

func TestOptClientArchTypeParseAndToBytes(t *testing.T) {
	data := []byte{
		0, 8, // EFI_XSCALE
	}
	opt, err := parseOptClientArchType(data)
	require.NoError(t, err)
	require.Equal(t, data, opt.ToBytes())
}

func TestOptClientArchType(t *testing.T) {
	opt := OptClientArchType(iana.EFI_ITANIUM)
	require.Equal(t, OptionClientArchType, opt.Code())
	require.Contains(t, opt.String(), "EFI Itanium", "String() should contain the correct ArchType output")
}
