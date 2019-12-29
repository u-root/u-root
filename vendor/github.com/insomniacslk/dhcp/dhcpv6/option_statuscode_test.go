package dhcpv6

import (
	"testing"

	"github.com/insomniacslk/dhcp/iana"
	"github.com/stretchr/testify/require"
)

func TestParseOptStatusCode(t *testing.T) {
	data := []byte{
		0, 5, // StatusUseMulticast
		'u', 's', 'e', ' ', 'm', 'u', 'l', 't', 'i', 'c', 'a', 's', 't',
	}
	opt, err := ParseOptStatusCode(data)
	require.NoError(t, err)
	require.Equal(t, iana.StatusUseMulticast, opt.StatusCode)
	require.Equal(t, []byte("use multicast"), opt.StatusMessage)
}

func TestOptStatusCodeToBytes(t *testing.T) {
	expected := []byte{
		0, 0, // StatusSuccess
		's', 'u', 'c', 'c', 'e', 's', 's',
	}
	opt := OptStatusCode{
		iana.StatusSuccess,
		[]byte("success"),
	}
	actual := opt.ToBytes()
	require.Equal(t, expected, actual)
}

func TestOptStatusCodeParseOptStatusCodeTooShort(t *testing.T) {
	_, err := ParseOptStatusCode([]byte{0})
	require.Error(t, err, "ParseOptStatusCode: Expected error on truncated option")
}

func TestOptStatusCodeString(t *testing.T) {
	data := []byte{
		0, 5, // StatusUseMulticast
		'u', 's', 'e', ' ', 'm', 'u', 'l', 't', 'i', 'c', 'a', 's', 't',
	}
	opt, err := ParseOptStatusCode(data)
	require.NoError(t, err)

	require.Contains(
		t,
		opt.String(),
		"code=UseMulticast (5), message=use multicast",
		"String() should contain the code and message",
	)
}
