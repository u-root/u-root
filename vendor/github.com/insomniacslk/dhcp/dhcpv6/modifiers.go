package dhcpv6

import (
	"log"
	"net"

	"github.com/insomniacslk/dhcp/iana"
	"github.com/insomniacslk/dhcp/rfc1035label"
)

// WithClientID adds a client ID option to a DHCPv6 packet
func WithClientID(duid Duid) Modifier {
	return func(d DHCPv6) {
		cid := OptClientId{Cid: duid}
		d.UpdateOption(&cid)
	}
}

// WithServerID adds a client ID option to a DHCPv6 packet
func WithServerID(duid Duid) Modifier {
	return func(d DHCPv6) {
		sid := OptServerId{Sid: duid}
		d.UpdateOption(&sid)
	}
}

// WithNetboot adds bootfile URL and bootfile param options to a DHCPv6 packet.
func WithNetboot(d DHCPv6) {
	msg, ok := d.(*Message)
	if !ok {
		log.Printf("WithNetboot: not a Message")
		return
	}
	// add OptionBootfileURL and OptionBootfileParam
	opt := msg.GetOneOption(OptionORO)
	if opt == nil {
		opt = &OptRequestedOption{}
	}
	// TODO only add options if they are not there already
	oro := opt.(*OptRequestedOption)
	oro.AddRequestedOption(OptionBootfileURL)
	oro.AddRequestedOption(OptionBootfileParam)
	msg.UpdateOption(oro)
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
		ao := OptClientArchType{ArchTypes: []iana.Arch{at}}
		d.AddOption(&ao)
	}
}

// WithIANA adds or updates an OptIANA option with the provided IAAddress
// options
func WithIANA(addrs ...OptIAAddress) Modifier {
	return func(d DHCPv6) {
		opt := d.GetOneOption(OptionIANA)
		if opt == nil {
			opt = &OptIANA{}
		}
		iaNa := opt.(*OptIANA)
		for _, addr := range addrs {
			iaNa.AddOption(&addr)
		}
		d.UpdateOption(iaNa)
	}
}

// WithIAID updates an OptIANA option with the provided IAID
func WithIAID(iaid [4]byte) Modifier {
	return func(d DHCPv6) {
		opt := d.GetOneOption(OptionIANA)
		if opt == nil {
			opt = &OptIANA{
				Options: Options{},
			}
		}
		iaNa := opt.(*OptIANA)
		copy(iaNa.IaId[:], iaid[:])
		d.UpdateOption(iaNa)
	}
}

// WithDNS adds or updates an OptDNSRecursiveNameServer
func WithDNS(dnses ...net.IP) Modifier {
	return func(d DHCPv6) {
		odns := OptDNSRecursiveNameServer{
			NameServers: append([]net.IP{}, dnses[:]...),
		}
		d.UpdateOption(&odns)
	}
}

// WithDomainSearchList adds or updates an OptDomainSearchList
func WithDomainSearchList(searchlist ...string) Modifier {
	return func(d DHCPv6) {
		osl := OptDomainSearchList{
			DomainSearchList: &rfc1035label.Labels{
				Labels: searchlist,
			},
		}
		d.UpdateOption(&osl)
	}
}

// WithRapidCommit adds the rapid commit option to a message.
func WithRapidCommit(d DHCPv6) {
	d.UpdateOption(&OptionGeneric{OptionCode: OptionRapidCommit})
}

// WithRequestedOptions adds requested options to the packet
func WithRequestedOptions(optionCodes ...OptionCode) Modifier {
	return func(d DHCPv6) {
		opt := d.GetOneOption(OptionORO)
		if opt == nil {
			opt = &OptRequestedOption{}
		}
		oro := opt.(*OptRequestedOption)
		for _, optionCode := range optionCodes {
			oro.AddRequestedOption(optionCode)
		}
		d.UpdateOption(oro)
	}
}
