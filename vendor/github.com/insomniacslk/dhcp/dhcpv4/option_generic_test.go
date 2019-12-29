package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptionGenericCode(t *testing.T) {
	o := OptGeneric(OptionDHCPMessageType, []byte{byte(MessageTypeDiscover)})
	require.Equal(t, OptionDHCPMessageType, o.Code)
	require.Equal(t, []byte{1}, o.Value.ToBytes())
	require.Equal(t, "DHCP Message Type: [1]", o.String())
}

func TestOptionGenericStringUnknown(t *testing.T) {
	o := OptGeneric(optionCode(102), []byte{byte(MessageTypeDiscover)})
	require.Equal(t, "unknown (102): [1]", o.String())
}
