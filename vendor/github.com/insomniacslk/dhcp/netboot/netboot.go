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
				// don't wait at the end of the last attempt
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
				// don't wait at the end of the last attempt
				break
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
func ConversationToNetconf(conversation []dhcpv6.DHCPv6) (*NetConf, string, error) {
	var reply dhcpv6.DHCPv6
	for _, m := range conversation {
		// look for a REPLY
		if m.Type() == dhcpv6.MessageTypeReply {
			reply = m
			break
		}
	}
	if reply == nil {
		return nil, "", errors.New("no REPLY received")
	}
	netconf, err := GetNetConfFromPacketv6(reply.(*dhcpv6.Message))
	if err != nil {
		return nil, "", fmt.Errorf("cannot get netconf from packet: %v", err)
	}
	// look for boot file
	var (
		opt      dhcpv6.Option
		bootfile string
	)
	opt = reply.GetOneOption(dhcpv6.OptionBootfileURL)
	if opt == nil {
		log.Printf("no bootfile URL option found in REPLY, looking for it in ADVERTISE")
		// as a fallback, look for bootfile URL in the advertise
		var advertise dhcpv6.DHCPv6
		for _, m := range conversation {
			// look for an ADVERTISE
			if m.Type() == dhcpv6.MessageTypeAdvertise {
				advertise = m
				break
			}
		}
		if advertise == nil {
			return nil, "", errors.New("no ADVERTISE found")
		}
		opt = advertise.GetOneOption(dhcpv6.OptionBootfileURL)
		if opt == nil {
			return nil, "", errors.New("no bootfile URL option found in ADVERTISE")
		}
	}
	if opt != nil {
		obf := opt.(*dhcpv6.OptBootFileURL)
		bootfile = string(obf.BootFileURL)
	}
	return netconf, bootfile, nil
}

// ConversationToNetconfv4 extracts network configuration and boot file URL from a
// DHCPv4 4-way conversation and returns them, or an error if any.
func ConversationToNetconfv4(conversation []*dhcpv4.DHCPv4) (*NetConf, string, error) {
	var reply *dhcpv4.DHCPv4
	var bootFileURL string
	for _, m := range conversation {
		// look for a BootReply packet of type Offer containing the bootfile URL.
		// Normally both packets with Message Type OFFER or ACK do contain
		// the bootfile URL.
		if m.OpCode == dhcpv4.OpcodeBootReply && m.MessageType() == dhcpv4.MessageTypeOffer {
			bootFileURL = m.BootFileName
			reply = m
			break
		}
	}
	if reply == nil {
		return nil, "", errors.New("no OFFER with valid bootfile URL received")
	}
	netconf, err := GetNetConfFromPacketv4(reply)
	if err != nil {
		return nil, "", fmt.Errorf("could not get netconf: %v", err)
	}
	return netconf, bootFileURL, nil
}
