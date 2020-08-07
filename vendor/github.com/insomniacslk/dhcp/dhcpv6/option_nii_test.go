package dhcpv6

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptNetworkInterfaceIdParse(t *testing.T) {
	expected := []byte{
		1,  // type (UNDI)
		3,  // major revision
		20, // minor revision
	}
	var opt OptNetworkInterfaceID
	err := opt.FromBytes(expected)
	require.NoError(t, err, "ParseOptNetworkInterfaceId() should not return an error with correct bytes")
	require.Equal(t, OptionNII, opt.Code(), OptionNII, "Code() should return 62 for OptNetworkInterfaceId")
	require.Equal(t, NetworkInterfaceType(1), opt.Typ, "Typ should return 1 for UNDI")
	require.Equal(t, uint8(3), opt.Major, "Major should return 1 for UNDI")
	require.Equal(t, uint8(20), opt.Minor, "Minor should return 1 for UNDI")
}

func TestOptNetworkInterfaceIdToBytes(t *testing.T) {
	expected := []byte{
		1,  // type (UNDI)
		3,  // major revision
		20, // minor revision
	}
	var opt OptNetworkInterfaceID
	opt.Typ = NetworkInterfaceType(1)
	opt.Major = 3
	opt.Minor = 20
	require.Equal(t, expected, opt.ToBytes())
}

func TestOptNetworkInterfaceIdTooShort(t *testing.T) {
	buf := []byte{
		1, // type (UNDI)
		// missing major/minor revision bytes
	}
	var opt OptNetworkInterfaceID
	err := opt.FromBytes(buf)
	require.Error(t, err, "ParseOptNetworkInterfaceId() should return an error on truncated options")
}

func TestOptNetworkInterfaceIdString(t *testing.T) {
	buf := []byte{
		1,  // type (UNDI)
		3,  // major revision
		20, // minor revision
	}
	var opt OptNetworkInterfaceID
	err := opt.FromBytes(buf)
	require.NoError(t, err)
	require.Contains(
		t,
		opt.String(),
		"First gen. PXE boot ROMs (Revision 3.20)",
		"String() should contain the type and revision",
	)
	opt.Typ = NetworkInterfaceType(200)
	require.Contains(
		t, opt.String(),
		"NetworkInterfaceType(200, unknown)",
		"String() should contain unknown for an unknown type",
	)
}
