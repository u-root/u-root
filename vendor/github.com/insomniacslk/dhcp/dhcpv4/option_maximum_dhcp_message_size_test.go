package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptMaximumDHCPMessageSize(t *testing.T) {
	o := OptMaxMessageSize(1500)
	require.Equal(t, OptionMaximumDHCPMessageSize, o.Code, "Code")
	require.Equal(t, []byte{5, 220}, o.Value.ToBytes(), "ToBytes")
	require.Equal(t, "Maximum DHCP Message Size: 1500", o.String())
}

func TestGetMaximumDHCPMessageSize(t *testing.T) {
	m, _ := New(WithGeneric(OptionMaximumDHCPMessageSize, []byte{5, 220}))
	o, err := m.MaxMessageSize()
	require.NoError(t, err)
	require.Equal(t, uint16(1500), o)

	// Short byte stream
	m, _ = New(WithGeneric(OptionMaximumDHCPMessageSize, []byte{2}))
	_, err = m.MaxMessageSize()
	require.Error(t, err, "should get error from short byte stream")

	// Bad length
	m, _ = New(WithGeneric(OptionMaximumDHCPMessageSize, []byte{2, 2, 2}))
	_, err = m.MaxMessageSize()
	require.Error(t, err, "should get error from bad length")
}
