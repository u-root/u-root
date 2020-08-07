package bsdp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptMessageTypeInterfaceMethods(t *testing.T) {
	o := OptMessageType(MessageTypeList)
	require.Equal(t, OptionMessageType, o.Code, "Code")
	require.Equal(t, []byte{1}, o.Value.ToBytes(), "ToBytes")
}

func TestParseOptMessageType(t *testing.T) {
	var o MessageType
	data := []byte{1} // DISCOVER
	err := o.FromBytes(data)
	require.NoError(t, err)
	require.Equal(t, MessageTypeList, o)
}

func TestOptMessageTypeString(t *testing.T) {
	// known
	o := OptMessageType(MessageTypeList)
	require.Equal(t, "BSDP Message Type: LIST", o.String())

	// unknown
	o = OptMessageType(99)
	require.Equal(t, "BSDP Message Type: unknown (99)", o.String())
}
