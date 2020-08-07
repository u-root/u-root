package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptMessageType(t *testing.T) {
	o := OptMessageType(MessageTypeDiscover)
	require.Equal(t, OptionDHCPMessageType, o.Code, "Code")
	require.Equal(t, []byte{1}, o.Value.ToBytes(), "ToBytes")
	require.Equal(t, "DHCP Message Type: DISCOVER", o.String())

	// unknown
	o = OptMessageType(99)
	require.Equal(t, "DHCP Message Type: unknown (99)", o.String())
}

func TestParseOptMessageType(t *testing.T) {
	var m MessageType
	data := []byte{1} // DISCOVER
	err := m.FromBytes(data)
	require.NoError(t, err)
	require.Equal(t, MessageTypeDiscover, m)

	// Bad length
	data = []byte{1, 2}
	err = m.FromBytes(data)
	require.Error(t, err, "should get error from bad length")
}

func TestGetMessageType(t *testing.T) {
	m, _ := New(WithMessageType(MessageTypeDiscover))
	require.Equal(t, MessageTypeDiscover, m.MessageType())

	m, _ = New()
	require.Equal(t, MessageTypeNone, m.MessageType())
}
