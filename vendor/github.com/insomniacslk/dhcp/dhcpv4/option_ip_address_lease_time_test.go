package dhcpv4

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestOptIPAddressLeaseTime(t *testing.T) {
	o := OptIPAddressLeaseTime(43200 * time.Second)
	require.Equal(t, OptionIPAddressLeaseTime, o.Code, "Code")
	require.Equal(t, []byte{0, 0, 168, 192}, o.Value.ToBytes(), "ToBytes")
	require.Equal(t, "IP Addresses Lease Time: 12h0m0s", o.String(), "String")
}

func TestGetIPAddressLeaseTime(t *testing.T) {
	m, _ := New(WithGeneric(OptionIPAddressLeaseTime, []byte{0, 0, 168, 192}))
	leaseTime := m.IPAddressLeaseTime(0)
	require.Equal(t, 43200*time.Second, leaseTime)

	// Too short.
	m, _ = New(WithGeneric(OptionIPAddressLeaseTime, []byte{168, 192}))
	leaseTime = m.IPAddressLeaseTime(0)
	require.Equal(t, time.Duration(0), leaseTime)

	// Too long.
	m, _ = New(WithGeneric(OptionIPAddressLeaseTime, []byte{1, 1, 1, 1, 1}))
	leaseTime = m.IPAddressLeaseTime(0)
	require.Equal(t, time.Duration(0), leaseTime)

	// Empty.
	m, _ = New()
	require.Equal(t, time.Duration(10), m.IPAddressLeaseTime(10))
}
