package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptParameterRequestListInterfaceMethods(t *testing.T) {
	opts := []OptionCode{OptionBootfileName, OptionNameServer}
	o := OptParameterRequestList(opts...)

	require.Equal(t, OptionParameterRequestList, o.Code, "Code")

	expectedBytes := []byte{67, 5}
	require.Equal(t, expectedBytes, o.Value.ToBytes(), "ToBytes")

	expectedString := "Parameter Request List: Name Server, Bootfile Name"
	require.Equal(t, expectedString, o.String(), "String")
}

func TestParseOptParameterRequestList(t *testing.T) {
	var o OptionCodeList
	err := o.FromBytes([]byte{67, 5})
	require.NoError(t, err)
	expectedOpts := OptionCodeList{OptionBootfileName, OptionNameServer}
	require.Equal(t, expectedOpts, o)
}
