package dhcpv6

import (
	"log"
	"net"
	"testing"

	"github.com/insomniacslk/dhcp/iana"
	"github.com/stretchr/testify/require"
)

func TestWithClientID(t *testing.T) {
	duid := Duid{
		Type:          DUID_LL,
		HwType:        iana.HWTypeEthernet,
		LinkLayerAddr: net.HardwareAddr([]byte{0xfa, 0xce, 0xb0, 0x00, 0x00, 0x0c}),
	}
	m, err := NewMessage(WithClientID(duid))
	require.NoError(t, err)
	cid := m.Options.ClientID()
	require.Equal(t, cid, &duid)
}

func TestWithServerID(t *testing.T) {
	duid := Duid{
		Type:          DUID_LL,
		HwType:        iana.HWTypeEthernet,
		LinkLayerAddr: net.HardwareAddr([]byte{0xfa, 0xce, 0xb0, 0x00, 0x00, 0x0c}),
	}
	m, err := NewMessage(WithServerID(duid))
	require.NoError(t, err)
	sid := m.Options.ServerID()
	require.Equal(t, sid, &duid)
}

func TestWithRequestedOptions(t *testing.T) {
	// Check if ORO is created when no ORO present
	m, err := NewMessage(WithRequestedOptions(OptionClientID))
	require.NoError(t, err)
	oro := m.Options.RequestedOptions()
	require.ElementsMatch(t, oro, OptionCodes{OptionClientID})
	// Check if already set options are preserved
	WithRequestedOptions(OptionServerID)(m)
	oro = m.Options.RequestedOptions()
	require.ElementsMatch(t, oro, OptionCodes{OptionClientID, OptionServerID})
}

func TestWithIANA(t *testing.T) {
	var d Message
	WithIANA(OptIAAddress{
		IPv6Addr:          net.ParseIP("::1"),
		PreferredLifetime: 3600,
		ValidLifetime:     5200,
	})(&d)
	require.Equal(t, 1, len(d.Options.Options))
	require.Equal(t, OptionIANA, d.Options.Options[0].Code())
}

func TestWithDNS(t *testing.T) {
	var d Message
	WithDNS(
		net.ParseIP("fe80::1"),
		net.ParseIP("fe80::2"),
	)(&d)
	require.Equal(t, 1, len(d.Options.Options))
	dns := d.Options.DNS()
	log.Printf("DNS %+v", dns)
	require.Equal(t, 2, len(dns))
	require.Equal(t, net.ParseIP("fe80::1"), dns[0])
	require.Equal(t, net.ParseIP("fe80::2"), dns[1])
	require.NotEqual(t, net.ParseIP("fe80::1"), dns[1])
}

func TestWithDomainSearchList(t *testing.T) {
	var d Message
	WithDomainSearchList("slackware.it", "dhcp.slackware.it")(&d)
	require.Equal(t, 1, len(d.Options.Options))
	osl := d.Options.DomainSearchList()
	require.NotNil(t, osl)
	labels := osl.Labels
	require.Equal(t, 2, len(labels))
	require.Equal(t, "slackware.it", labels[0])
	require.Equal(t, "dhcp.slackware.it", labels[1])
}

func TestWithFQDN(t *testing.T) {
	var d Message
	WithFQDN(4, "cnos.localhost")(&d)
	require.Equal(t, 1, len(d.Options.Options))
	ofqdn := d.Options.FQDN()
	require.Equal(t, OptionFQDN, ofqdn.Code())
	require.Equal(t, uint8(4), ofqdn.Flags)
	require.Equal(t, "cnos.localhost", ofqdn.DomainName.Labels[0])
}

func TestWithDHCP4oDHCP6Server(t *testing.T) {
	var d Message
	WithDHCP4oDHCP6Server([]net.IP{
		net.ParseIP("fe80::1"),
		net.ParseIP("fe80::2"),
	}...)(&d)
	require.Equal(t, 1, len(d.Options.Options))
	opt := d.Options.DHCP4oDHCP6Server()
	require.Equal(t, OptionDHCP4oDHCP6Server, opt.Code())
	require.Equal(t, 2, len(opt.DHCP4oDHCP6Servers))
	require.Equal(t, net.ParseIP("fe80::1"), opt.DHCP4oDHCP6Servers[0])
	require.Equal(t, net.ParseIP("fe80::2"), opt.DHCP4oDHCP6Servers[1])
	require.NotEqual(t, net.ParseIP("fe80::1"), opt.DHCP4oDHCP6Servers[1])
}

func TestWithIAPD(t *testing.T) {
	var d Message
	_, pre, _ := net.ParseCIDR("2001:DB8:7689::/48")
	prefix := &OptIAPrefix{
		PreferredLifetime: 3600,
		ValidLifetime:     5200,
		Prefix:            pre,
	}
	WithIAPD([4]byte{1, 2, 3, 4}, prefix)(&d)
	opt := d.Options.IAPD()
	require.Equal(t, 1, len(opt))
	require.Equal(t, OptionIAPD, opt[0].Code())
}

func TestWithClientLinkLayerAddress(t *testing.T) {
	var d RelayMessage
	mac, _ := net.ParseMAC("a4:83:e7:e3:df:88")
	WithClientLinkLayerAddress(iana.HWTypeEthernet, mac)(&d)

	opt := d.Options.GetOne(OptionClientLinkLayerAddr)
	require.Equal(t, OptionClientLinkLayerAddr, opt.Code())

	llt, lla := d.Options.ClientLinkLayerAddress()
	require.Equal(t, iana.HWTypeEthernet, llt)
	require.Equal(t, mac, lla)
}
