package netboot

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/client4"
	"github.com/insomniacslk/dhcp/dhcpv6"
	"github.com/insomniacslk/dhcp/dhcpv6/client6"
)

var sleeper = func(d time.Duration) {
	time.Sleep(d)
}

// BootConf is a structure describes everything a host needs to know to boot over network
type BootConf struct {
	// NetConf is the network configuration of the client
	NetConf

	// BootfileURL is "where is the image (kernel)".
	// See RFC5970 section 3.1 for IPv6 and RFC2132 section 9.5 ("Bootfile name") for IPv4
	BootfileURL string

	// BootfileParam is "what arguments should we pass (cmdline)".
	// See RFC5970 section 3.2 for IPv6.
	BootfileParam []string
}

// RequestNetbootv6 sends a netboot request via DHCPv6 and returns the exchanged packets. Additional modifiers
// can be passed to manipulate both solicit and advertise packets.
func RequestNetbootv6(ifname string, timeout time.Duration, retries int, modifiers ...dhcpv6.Modifier) ([]dhcpv6.DHCPv6, error) {
	var (
		conversation []dhcpv6.DHCPv6
		err          error
	)
	modifiers = append(modifiers, dhcpv6.WithNetboot)
	delay := 2 * time.Second
	for i := 0; i <= retries; i++ {
		log.Printf("sending request, attempt #%d", i+1)

		client := client6.NewClient()
		client.ReadTimeout = timeout
		conversation, err = client.Exchange(ifname, modifiers...)
		if err != nil {
			log.Printf("Client.Exchange failed: %v", err)
			if i >= retries {
				return nil, fmt.Errorf("netboot failed after %d attempts: %v", retries+1, err)
			}
			log.Printf("sleeping %v before retrying", delay)
			sleeper(delay)
			// TODO add random splay
			delay = delay * 2
			continue
		}
		break
	}
	return conversation, nil
}

// RequestNetbootv4 sends a netboot request via DHCPv4 and returns the exchanged packets. Additional modifiers
// can be passed to manipulate both the discover and offer packets.
func RequestNetbootv4(ifname string, timeout time.Duration, retries int, modifiers ...dhcpv4.Modifier) ([]*dhcpv4.DHCPv4, error) {
	var (
		conversation []*dhcpv4.DHCPv4
		err          error
	)
	delay := 2 * time.Second
	modifiers = append(modifiers, dhcpv4.WithNetboot)
	for i := 0; i <= retries; i++ {
		log.Printf("sending request, attempt #%d", i+1)
		client := client4.NewClient()
		client.ReadTimeout = timeout
		conversation, err = client.Exchange(ifname, modifiers...)
		if err != nil {
			log.Printf("Client.Exchange failed: %v", err)
			log.Printf("sleeping %v before retrying", delay)
			if i >= retries {
				return nil, fmt.Errorf("netboot failed after %d attempts: %v", retries+1, err)
			}
			sleeper(delay)
			// TODO add random splay
			delay = delay * 2
			continue
		}
		break
	}
	return conversation, nil
}

// ConversationToNetconf extracts network configuration and boot file URL from a
// DHCPv6 4-way conversation and returns them, or an error if any.
func ConversationToNetconf(conversation []dhcpv6.DHCPv6) (*BootConf, error) {
	var advertise, reply *dhcpv6.Message
	for _, m := range conversation {
		switch m.Type() {
		case dhcpv6.MessageTypeAdvertise:
			advertise = m.(*dhcpv6.Message)
		case dhcpv6.MessageTypeReply:
			reply = m.(*dhcpv6.Message)
		}
	}
	if reply == nil {
		return nil, errors.New("no REPLY received")
	}

	bootconf := &BootConf{}
	netconf, err := GetNetConfFromPacketv6(reply)
	if err != nil {
		return nil, fmt.Errorf("cannot get netconf from packet: %v", err)
	}
	bootconf.NetConf = *netconf

	if u := reply.Options.BootFileURL(); len(u) > 0 {
		bootconf.BootfileURL = u
		bootconf.BootfileParam = reply.Options.BootFileParam()
	} else {
		log.Printf("no bootfile URL option found in REPLY, fallback to ADVERTISE's value")
		if u := advertise.Options.BootFileURL(); len(u) > 0 {
			bootconf.BootfileURL = u
			bootconf.BootfileParam = advertise.Options.BootFileParam()
		}
	}
	if len(bootconf.BootfileURL) == 0 {
		return nil, errors.New("no bootfile URL option found")
	}
	return bootconf, nil
}

// ConversationToNetconfv4 extracts network configuration and boot file URL from a
// DHCPv4 4-way conversation and returns them, or an error if any.
func ConversationToNetconfv4(conversation []*dhcpv4.DHCPv4) (*BootConf, error) {
	var reply *dhcpv4.DHCPv4
	for _, m := range conversation {
		// look for a BootReply packet of type Offer containing the bootfile URL.
		// Normally both packets with Message Type OFFER or ACK do contain
		// the bootfile URL.
		if m.OpCode == dhcpv4.OpcodeBootReply && m.MessageType() == dhcpv4.MessageTypeOffer {
			reply = m
			break
		}
	}
	if reply == nil {
		return nil, errors.New("no OFFER with valid bootfile URL received")
	}

	bootconf := &BootConf{}
	netconf, err := GetNetConfFromPacketv4(reply)
	if err != nil {
		return nil, fmt.Errorf("could not get netconf: %v", err)
	}
	bootconf.NetConf = *netconf

	bootconf.BootfileURL = reply.BootFileName
	// TODO: should we support bootfile parameters here somehow? (see netconf.BootfileParam)
	return bootconf, nil
}
