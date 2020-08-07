package main

import (
	"log"
	"net"

	"github.com/insomniacslk/dhcp/dhcpv6"
	"github.com/insomniacslk/dhcp/iana"
)

func main() {
	// In this example we create and manipulate a DHCPv6 solicit packet
	// and encapsulate it in a relay packet. To to this, we use
	// `dhcpv6.Message` and `dhcpv6.DHCPv6Relay`, two structures
	// that implement the `dhcpv6.DHCPv6` interface.
	// Then print the wire-format representation of the packet.

	iface, err := net.InterfaceByName("eth0")
	if err != nil {
		log.Fatal(err)
	}

	// Create the DHCPv6 Solicit first, using the interface "eth0"
	// to get the MAC address
	msg, err := dhcpv6.NewSolicit(iface.HardwareAddr)
	if err != nil {
		log.Fatal(err)
	}

	// In this example I want to redact the MAC address of my
	// network interface, so instead of replacing it manually,
	// I will show how to use modifiers for the purpose.
	// A Modifier is simply a function that can be applied on
	// a DHCPv6 object to manipulate it. Here we use it to
	// replace the MAC address with a dummy one.
	// Modifiers can be passed to many functions, for example
	// to constructors, `Exchange()`, `Solicit()`, etc. Check
	// the source code to know where to use them.
	// Existing modifiers are implemented in dhcpv6/modifiers.go .
	mac, err := net.ParseMAC("00:fa:ce:b0:0c:00")
	if err != nil {
		log.Fatal(err)
	}
	duid := dhcpv6.Duid{
		Type:          dhcpv6.DUID_LLT,
		HwType:        iana.HWTypeEthernet,
		Time:          dhcpv6.GetTime(),
		LinkLayerAddr: mac,
	}
	// As suggested above, an alternative is to call
	// dhcpv6.NewSolicitForInterface("eth0", dhcpv6.WithClientID(duid))
	dhcpv6.WithClientID(duid)(msg)

	// Now encapsulate the message in a DHCPv6 relay.
	// As per RFC3315, the link-address and peer-address have
	// to be set by the relay agent. We use dummy values here.
	linkAddr := net.ParseIP("2001:0db8::1")
	peerAddr := net.ParseIP("2001:0db8::2")
	relay, err := dhcpv6.EncapsulateRelay(msg, dhcpv6.MessageTypeRelayForward, linkAddr, peerAddr)
	if err != nil {
		log.Fatal(err)
	}

	// Print a verbose representation of the relay packet, that will also
	// show a short representation of the inner Solicit message.
	// To print a detailed summary of the inner packet, extract it
	// first from the relay using `relay.GetInnerMessage()`.
	log.Print(relay.Summary())

	// And finally, print the bytes that would be sent on the wire
	log.Print(relay.ToBytes())

	// Note: there are many more functions in the library, check them
	// out in the source code. For example, if you want to decode a
	// byte stream into a DHCPv6 message or relay, you can use
	// `dhcpv6.FromBytes`.
}
