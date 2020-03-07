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
		ofqdn := OptFQDN{Flags: flags, DomainName: domainname}
		d.AddOption(&ofqdn)
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
				iana.AddOption(&addr)
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
					Options: Options{},
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
