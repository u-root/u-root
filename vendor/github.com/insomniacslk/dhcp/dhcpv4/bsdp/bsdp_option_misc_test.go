package bsdp

import (
	"net"
	"testing"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/stretchr/testify/require"
)

func TestOptReplyPort(t *testing.T) {
	o := OptReplyPort(1234)
	require.Equal(t, OptionReplyPort, o.Code, "Code")
	require.Equal(t, []byte{4, 210}, o.Value.ToBytes(), "ToBytes")
	require.Equal(t, "BSDP Reply Port: 1234", o.String())
}

func TestGetReplyPort(t *testing.T) {
	o := VendorOptions{dhcpv4.OptionsFromList(OptReplyPort(1234))}
	port, err := o.ReplyPort()
	require.NoError(t, err)
	require.Equal(t, uint16(1234), port)

	o = VendorOptions{dhcpv4.Options{}}
	_, err = o.ReplyPort()
	require.Error(t, err, "no reply port present")
}

func TestOptServerPriority(t *testing.T) {
	o := OptServerPriority(1234)
	require.Equal(t, OptionServerPriority, o.Code, "Code")
	require.Equal(t, []byte{4, 210}, o.Value.ToBytes(), "ToBytes")
	require.Equal(t, "BSDP Server Priority: 1234", o.String())
}

func TestGetServerPriority(t *testing.T) {
	o := VendorOptions{dhcpv4.OptionsFromList(OptServerPriority(1234))}
	prio, err := o.ServerPriority()
	require.NoError(t, err)
	require.Equal(t, uint16(1234), prio)

	o = VendorOptions{dhcpv4.Options{}}
	_, err = o.ServerPriority()
	require.Error(t, err, "no server prio present")
}

func TestOptMachineName(t *testing.T) {
	o := OptMachineName("foo")
	require.Equal(t, OptionMachineName, o.Code, "Code")
	require.Equal(t, []byte("foo"), o.Value.ToBytes(), "ToBytes")
	require.Equal(t, "BSDP Machine Name: foo", o.String())
}

func TestGetMachineName(t *testing.T) {
	o := VendorOptions{dhcpv4.OptionsFromList(OptMachineName("foo"))}
	require.Equal(t, "foo", o.MachineName())

	o = VendorOptions{dhcpv4.Options{}}
	require.Equal(t, "", o.MachineName())
}

func TestOptVersion(t *testing.T) {
	o := OptVersion(Version1_1)
	require.Equal(t, OptionVersion, o.Code, "Code")
	require.Equal(t, []byte{1, 1}, o.Value.ToBytes(), "ToBytes")
	require.Equal(t, "BSDP Version: 1.1", o.String())
}

func TestGetVersion(t *testing.T) {
	o := VendorOptions{dhcpv4.OptionsFromList(OptVersion(Version1_1))}
	ver, err := o.Version()
	require.NoError(t, err)
	require.Equal(t, ver, Version1_1)

	o = VendorOptions{dhcpv4.Options{}}
	_, err = o.Version()
	require.Error(t, err, "no version present")

	o = VendorOptions{dhcpv4.Options{OptionVersion.Code(): []byte{}}}
	_, err = o.Version()
	require.Error(t, err, "empty version field")

	o = VendorOptions{dhcpv4.Options{OptionVersion.Code(): []byte{1}}}
	_, err = o.Version()
	require.Error(t, err, "version option too short")

	o = VendorOptions{dhcpv4.Options{OptionVersion.Code(): []byte{1, 2, 3}}}
	_, err = o.Version()
	require.Error(t, err, "version option too long")
}

func TestOptServerIdentifier(t *testing.T) {
	o := OptServerIdentifier(net.IP{1, 1, 1, 1})
	require.Equal(t, OptionServerIdentifier, o.Code, "Code")
	require.Equal(t, []byte{1, 1, 1, 1}, o.Value.ToBytes(), "ToBytes")
	require.Equal(t, "BSDP Server Identifier: 1.1.1.1", o.String())
}

func TestGetServerIdentifier(t *testing.T) {
	o := VendorOptions{dhcpv4.OptionsFromList(OptServerIdentifier(net.IP{1, 1, 1, 1}))}
	require.Equal(t, net.IP{1, 1, 1, 1}, o.ServerIdentifier())

	o = VendorOptions{dhcpv4.Options{}}
	require.Nil(t, o.ServerIdentifier())
}
