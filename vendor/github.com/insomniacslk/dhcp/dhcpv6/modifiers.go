package dhcpv6

import (
	"net"

	"github.com/insomniacslk/dhcp/iana"
	"github.com/insomniacslk/dhcp/rfc1035label"
)

// WithOption adds the specific option to the DHCPv6 message.
func WithOption(o Option) Modifier {
	return func(d DHCPv6) {
		d.UpdateOption(o)
	}
}

// WithClientID adds a client ID option to a DHCPv6 packet
func WithClientID(duid Duid) Modifier {
	return WithOption(OptClientID(duid))
}

// WithServerID adds a client ID option to a DHCPv6 packet
func WithServerID(duid Duid) Modifier {
	return WithOption(OptServerID(duid))
}

// WithNetboot adds bootfile URL and bootfile param options to a DHCPv6 packet.
func WithNetboot(d DHCPv6) {
	WithRequestedOptions(OptionBootfileURL, OptionBootfileParam)(d)
}

// WithFQDN adds a fully qualified domain name option to the packet
func WithFQDN(flags uint8, domainname string) Modifier {
	return func(d DHCPv6) {
		d.UpdateOption(&OptFQDN{
			Flags: flags,
			DomainName: &rfc1035label.Labels{
				Labels: []string{domainname},
			},
		})
	}
}

// WithUserClass adds a user class option to the packet
func WithUserClass(uc []byte) Modifier {
	// TODO let the user specify multiple user classes
	return func(d DHCPv6) {
		ouc := OptUserClass{UserClasses: [][]byte{uc}}
		d.AddOption(&ouc)
	}
}

// WithArchType adds an arch type option to the packet
func WithArchType(at iana.Arch) Modifier {
	return func(d DHCPv6) {
		d.AddOption(OptClientArchType(at))
	}
}

// WithIANA adds or updates an OptIANA option with the provided IAAddress
// options
func WithIANA(addrs ...OptIAAddress) Modifier {
	return func(d DHCPv6) {
		if msg, ok := d.(*Message); ok {
			iana := msg.Options.OneIANA()
			if iana == nil {
				iana = &OptIANA{}
			}
			for _, addr := range addrs {
				iana.Options.Add(&addr)
			}
			msg.UpdateOption(iana)
		}
	}
}

// WithIAID updates an OptIANA option with the provided IAID
func WithIAID(iaid [4]byte) Modifier {
	return func(d DHCPv6) {
		if msg, ok := d.(*Message); ok {
			iana := msg.Options.OneIANA()
			if iana == nil {
				iana = &OptIANA{
					Options: IdentityOptions{Options: []Option{}},
				}
			}
			copy(iana.IaId[:], iaid[:])
			d.UpdateOption(iana)
		}
	}
}

// WithDNS adds or updates an OptDNSRecursiveNameServer
func WithDNS(dnses ...net.IP) Modifier {
	return WithOption(OptDNS(dnses...))
}

// WithDomainSearchList adds or updates an OptDomainSearchList
func WithDomainSearchList(searchlist ...string) Modifier {
	return func(d DHCPv6) {
		d.UpdateOption(OptDomainSearchList(
			&rfc1035label.Labels{
				Labels: searchlist,
			},
		))
	}
}

// WithRapidCommit adds the rapid commit option to a message.
func WithRapidCommit(d DHCPv6) {
	d.UpdateOption(&OptionGeneric{OptionCode: OptionRapidCommit})
}

// WithRequestedOptions adds requested options to the packet
func WithRequestedOptions(codes ...OptionCode) Modifier {
	return func(d DHCPv6) {
		if msg, ok := d.(*Message); ok {
			oro := msg.Options.RequestedOptions()
			for _, c := range codes {
				oro.Add(c)
			}
			d.UpdateOption(OptRequestedOption(oro...))
		}
	}
}

// WithDHCP4oDHCP6Server adds or updates an OptDHCP4oDHCP6Server
func WithDHCP4oDHCP6Server(addrs ...net.IP) Modifier {
	return func(d DHCPv6) {
		opt := OptDHCP4oDHCP6Server{
			DHCP4oDHCP6Servers: addrs,
		}
		d.UpdateOption(&opt)
	}
}

// WithIAPD adds or updates an IAPD option with the provided IAID and
// prefix options to a DHCPv6 packet.
func WithIAPD(iaid [4]byte, prefixes ...*OptIAPrefix) Modifier {
	return func(d DHCPv6) {
		if msg, ok := d.(*Message); ok {
			opt := msg.Options.OneIAPD()
			if opt == nil {
				opt = &OptIAPD{}
			}
			copy(opt.IaId[:], iaid[:])

			for _, prefix := range prefixes {
				opt.Options.Add(prefix)
			}
			d.UpdateOption(opt)
		}
	}
}
