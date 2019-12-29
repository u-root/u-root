package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptDomainName(t *testing.T) {
	o := OptDomainName("foo")
	require.Equal(t, OptionDomainName, o.Code, "Code")
	require.Equal(t, []byte{'f', 'o', 'o'}, o.Value.ToBytes(), "ToBytes")
	require.Equal(t, "Domain Name: foo", o.String())
}

func TestParseOptDomainName(t *testing.T) {
	m, _ := New(WithGeneric(OptionDomainName, []byte{'t', 'e', 's', 't'}))
	require.Equal(t, "test", m.DomainName())

	m, _ = New()
	require.Equal(t, "", m.DomainName())
}

func TestOptHostName(t *testing.T) {
	o := OptHostName("foo")
	require.Equal(t, OptionHostName, o.Code, "Code")
	require.Equal(t, []byte{'f', 'o', 'o'}, o.Value.ToBytes(), "ToBytes")
	require.Equal(t, "Host Name: foo", o.String())
}

func TestParseOptHostName(t *testing.T) {
	m, _ := New(WithGeneric(OptionHostName, []byte{'t', 'e', 's', 't'}))
	require.Equal(t, "test", m.HostName())

	m, _ = New()
	require.Equal(t, "", m.HostName())
}

func TestOptRootPath(t *testing.T) {
	o := OptRootPath("foo")
	require.Equal(t, OptionRootPath, o.Code, "Code")
	require.Equal(t, []byte{'f', 'o', 'o'}, o.Value.ToBytes(), "ToBytes")
	require.Equal(t, "Root Path: foo", o.String())
}

func TestParseOptRootPath(t *testing.T) {
	m, _ := New(WithGeneric(OptionRootPath, []byte{'t', 'e', 's', 't'}))
	require.Equal(t, "test", m.RootPath())

	m, _ = New()
	require.Equal(t, "", m.RootPath())
}

func TestOptBootFileName(t *testing.T) {
	o := OptBootFileName("foo")
	require.Equal(t, OptionBootfileName, o.Code, "Code")
	require.Equal(t, []byte{'f', 'o', 'o'}, o.Value.ToBytes(), "ToBytes")
	require.Equal(t, "Bootfile Name: foo", o.String())
}

func TestParseOptBootFileName(t *testing.T) {
	m, _ := New(WithGeneric(OptionBootfileName, []byte{'t', 'e', 's', 't'}))
	require.Equal(t, "test", m.BootFileNameOption())

	m, _ = New()
	require.Equal(t, "", m.BootFileNameOption())

	m, _ = New(WithGeneric(OptionBootfileName, []byte{'t', 'e', 's', 't', 0}))
	require.Equal(t, "test", m.BootFileNameOption())
}

func TestOptTFTPServerName(t *testing.T) {
	o := OptTFTPServerName("foo")
	require.Equal(t, OptionTFTPServerName, o.Code, "Code")
	require.Equal(t, []byte{'f', 'o', 'o'}, o.Value.ToBytes(), "ToBytes")
	require.Equal(t, "TFTP Server Name: foo", o.String())
}

func TestParseOptTFTPServerName(t *testing.T) {
	m, _ := New(WithGeneric(OptionTFTPServerName, []byte{'t', 'e', 's', 't'}))
	require.Equal(t, "test", m.TFTPServerName())

	m, _ = New()
	require.Equal(t, "", m.TFTPServerName())

	m, _ = New(WithGeneric(OptionTFTPServerName, []byte{'t', 'e', 's', 't', 0}))
	require.Equal(t, "test", m.TFTPServerName())
}

func TestOptClassIdentifier(t *testing.T) {
	o := OptClassIdentifier("foo")
	require.Equal(t, OptionClassIdentifier, o.Code, "Code")
	require.Equal(t, []byte{'f', 'o', 'o'}, o.Value.ToBytes(), "ToBytes")
	require.Equal(t, "Class Identifier: foo", o.String())
}

func TestParseOptClassIdentifier(t *testing.T) {
	m, _ := New(WithGeneric(OptionClassIdentifier, []byte{'t', 'e', 's', 't'}))
	require.Equal(t, "test", m.ClassIdentifier())

	m, _ = New()
	require.Equal(t, "", m.ClassIdentifier())
}

func TestOptUserClass(t *testing.T) {
	o := OptUserClass("linuxboot")
	require.Equal(t, OptionUserClassInformation, o.Code, "Code")
	expected := []byte{
		'l', 'i', 'n', 'u', 'x', 'b', 'o', 'o', 't',
	}
	require.Equal(t, expected, o.Value.ToBytes(), "ToBytes")
	require.Equal(t, "User Class Information: linuxboot", o.String())
}

func TestParseOptUserClass(t *testing.T) {
	m, _ := New(WithUserClass("linuxboot", false))
	require.Equal(t, []string{"linuxboot"}, m.UserClass())

	m, _ = New()
	require.Equal(t, 0, len(m.UserClass()))
}
