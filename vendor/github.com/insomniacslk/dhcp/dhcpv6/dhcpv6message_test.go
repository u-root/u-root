package dhcpv6

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsNetboot(t *testing.T) {
	msg1 := Message{}
	require.False(t, msg1.IsNetboot())

	msg2 := Message{}
	msg2.AddOption(OptRequestedOption(OptionBootfileURL))
	require.True(t, msg2.IsNetboot())

	msg3 := Message{}
	optbf := OptBootFileURL("")
	msg3.AddOption(optbf)
	require.True(t, msg3.IsNetboot())
}

func TestIsOptionRequested(t *testing.T) {
	msg1 := Message{}
	require.False(t, msg1.IsOptionRequested(OptionDNSRecursiveNameServer))

	msg2 := Message{}
	msg2.AddOption(OptRequestedOption(OptionDNSRecursiveNameServer))
	require.True(t, msg2.IsOptionRequested(OptionDNSRecursiveNameServer))
}
