package dhcpv4

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptSubnetMask(t *testing.T) {
	o := OptSubnetMask(net.IPMask{255, 255, 255, 0})
	require.Equal(t, o.Code, OptionSubnetMask, "Code")
	require.Equal(t, "Subnet Mask: ffffff00", o.String(), "String")
	require.Equal(t, []byte{255, 255, 255, 0}, o.Value.ToBytes(), "ToBytes")
}

func TestGetSubnetMask(t *testing.T) {
	m, _ := New(WithOption(OptSubnetMask(net.IPMask{})))
	mask := m.SubnetMask()
	require.Nil(t, mask, "empty byte stream")

	m, _ = New(WithOption(OptSubnetMask(net.IPMask{255})))
	mask = m.SubnetMask()
	require.Nil(t, mask, "short byte stream")

	m, _ = New(WithOption(OptSubnetMask(net.IPMask{255, 255, 255, 0})))
	mask = m.SubnetMask()
	require.Equal(t, net.IPMask{255, 255, 255, 0}, mask)
}
