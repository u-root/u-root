package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetRelayAgentInformation(t *testing.T) {
	m, _ := New(WithGeneric(OptionRelayAgentInformation, []byte{
		1, 5, 'l', 'i', 'n', 'u', 'x',
		2, 4, 'b', 'o', 'o', 't',
	}))

	opt := m.RelayAgentInfo()
	require.NotNil(t, opt)
	require.Equal(t, len(opt.Options), 2)

	circuit := opt.Get(GenericOptionCode(1))
	remote := opt.Get(GenericOptionCode(2))
	require.Equal(t, circuit, []byte("linux"))
	require.Equal(t, remote, []byte("boot"))

	// Empty.
	m, _ = New()
	require.Nil(t, m.RelayAgentInfo())

	// Invalid contents.
	m, _ = New(WithGeneric(OptionRelayAgentInformation, []byte{
		1, 7, 'l', 'i', 'n', 'u', 'x',
	}))
	require.Nil(t, m.RelayAgentInfo())
}

func TestOptRelayAgentInfo(t *testing.T) {
	opt := OptRelayAgentInfo(
		OptGeneric(GenericOptionCode(1), []byte("linux")),
		OptGeneric(GenericOptionCode(2), []byte("boot")),
	)
	wantBytes := []byte{
		1, 5, 'l', 'i', 'n', 'u', 'x',
		2, 4, 'b', 'o', 'o', 't',
	}
	wantString := "Relay Agent Information:\n    unknown (1): [108 105 110 117 120]\n    unknown (2): [98 111 111 116]\n"
	require.Equal(t, wantBytes, opt.Value.ToBytes())
	require.Equal(t, OptionRelayAgentInformation, opt.Code)
	require.Equal(t, wantString, opt.String())
}
