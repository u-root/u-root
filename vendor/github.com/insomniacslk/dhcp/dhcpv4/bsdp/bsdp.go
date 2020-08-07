package bsdp

import (
	"errors"
	"fmt"
	"net"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

// MaxDHCPMessageSize is the size set in DHCP option 57 (DHCP Maximum Message Size).
// BSDP includes its own sub-option (12) to indicate to NetBoot servers that the
// client can support larger message sizes, and modern NetBoot servers will
// prefer this BSDP-specific option over the DHCP standard option.
const MaxDHCPMessageSize = 1500

// AppleVendorID is the string constant set in the vendor class identifier (DHCP
// option 60) that is sent by the server.
const AppleVendorID = "AAPLBSDPC"

// ReplyConfig is a struct containing some common configuration values for a
// BSDP reply (ACK).
type ReplyConfig struct {
	ServerIP                     net.IP
	ServerHostname, BootFileName string
	ServerPriority               uint16
	Images                       []BootImage
	DefaultImage, SelectedImage  *BootImage
}

// ParseBootImageListFromAck parses the list of boot images presented in the
// ACK[LIST] packet and returns them as a list of BootImages.
func ParseBootImageListFromAck(ack *dhcpv4.DHCPv4) ([]BootImage, error) {
	vendorOpts := GetVendorOptions(ack.Options)
	if vendorOpts == nil {
		return nil, errors.New("ParseBootImageListFromAck: could not find vendor-specific option")
	}
	return vendorOpts.BootImageList(), nil
}

func needsReplyPort(replyPort uint16) bool {
	return replyPort != 0 && replyPort != dhcpv4.ClientPort
}

// MessageTypeFromPacket extracts the BSDP message type (LIST, SELECT) from the
// vendor-specific options and returns it. If the message type option cannot be
// found, returns false.
func MessageTypeFromPacket(packet *dhcpv4.DHCPv4) MessageType {
	vendorOpts := GetVendorOptions(packet.Options)
	if vendorOpts == nil {
		return MessageTypeNone
	}
	return vendorOpts.MessageType()
}

// Packet is a BSDP packet wrapper around a DHCPv4 packet in order to print the
// correct vendor-specific BSDP information in String().
type Packet struct {
	dhcpv4.DHCPv4
}

// PacketFor returns a wrapped BSDP Packet given a DHCPv4 packet.
func PacketFor(d *dhcpv4.DHCPv4) *Packet {
	return &Packet{*d}
}

func (p Packet) v4() *dhcpv4.DHCPv4 {
	return &p.DHCPv4
}

func (p Packet) String() string {
	return p.DHCPv4.String()
}

// Summary prints the BSDP packet with the correct vendor-specific options.
func (p Packet) Summary() string {
	return p.DHCPv4.SummaryWithVendor(&VendorOptions{})
}

// NewInformListForInterface creates a new INFORM packet for interface ifname
// with configuration options specified by config.
func NewInformListForInterface(ifname string, replyPort uint16) (*Packet, error) {
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		return nil, err
	}
	// Get currently configured IP.
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}
	localIPs, err := dhcpv4.GetExternalIPv4Addrs(addrs)
	if err != nil {
		return nil, fmt.Errorf("could not get local IPv4 addr for %s: %v", iface.Name, err)
	}
	if len(localIPs) == 0 {
		return nil, fmt.Errorf("could not get local IPv4 addr for %s", iface.Name)
	}
	return NewInformList(iface.HardwareAddr, localIPs[0], replyPort)
}

// NewInformList creates a new INFORM packet for interface with hardware address
// `hwaddr` and IP `localIP`. Packet will be sent out on port `replyPort`.
func NewInformList(hwaddr net.HardwareAddr, localIP net.IP, replyPort uint16, modifiers ...dhcpv4.Modifier) (*Packet, error) {
	// Validate replyPort first
	if needsReplyPort(replyPort) && replyPort >= 1024 {
		return nil, errors.New("replyPort must be a privileged port")
	}

	vendorClassID, err := MakeVendorClassIdentifier()
	if err != nil {
		return nil, err
	}

	// These are vendor-specific options used to pass along BSDP information.
	vendorOpts := []dhcpv4.Option{
		OptMessageType(MessageTypeList),
		OptVersion(Version1_1),
	}
	if needsReplyPort(replyPort) {
		vendorOpts = append(vendorOpts, OptReplyPort(replyPort))
	}

	d, err := dhcpv4.NewInform(hwaddr, localIP,
		dhcpv4.PrependModifiers(modifiers, dhcpv4.WithRequestedOptions(
			dhcpv4.OptionVendorSpecificInformation,
			dhcpv4.OptionClassIdentifier,
		),
			dhcpv4.WithOption(dhcpv4.OptMaxMessageSize(MaxDHCPMessageSize)),
			dhcpv4.WithOption(dhcpv4.OptClassIdentifier(vendorClassID)),
			dhcpv4.WithOption(OptVendorOptions(vendorOpts...)),
		)...)
	if err != nil {
		return nil, err
	}
	return PacketFor(d), nil
}

// InformSelectForAck constructs an INFORM[SELECT] packet given an ACK to the
// previously-sent INFORM[LIST].
func InformSelectForAck(ack *Packet, replyPort uint16, selectedImage BootImage) (*Packet, error) {
	if needsReplyPort(replyPort) && replyPort >= 1024 {
		return nil, errors.New("replyPort must be a privileged port")
	}

	// Data for OptionSelectedBootImageID
	vendorOpts := []dhcpv4.Option{
		OptMessageType(MessageTypeSelect),
		OptVersion(Version1_1),
		OptSelectedBootImageID(selectedImage.ID),
	}

	// Validate replyPort if requested.
	if needsReplyPort(replyPort) {
		vendorOpts = append(vendorOpts, OptReplyPort(replyPort))
	}

	// Find server IP address
	serverIP := ack.ServerIdentifier()
	if serverIP.To4() == nil {
		return nil, fmt.Errorf("could not parse server identifier from ACK")
	}
	vendorOpts = append(vendorOpts, OptServerIdentifier(serverIP))

	vendorClassID, err := MakeVendorClassIdentifier()
	if err != nil {
		return nil, err
	}

	d, err := dhcpv4.New(dhcpv4.WithReply(ack.v4()),
		dhcpv4.WithOption(dhcpv4.OptClassIdentifier(vendorClassID)),
		dhcpv4.WithRequestedOptions(
			dhcpv4.OptionSubnetMask,
			dhcpv4.OptionRouter,
			dhcpv4.OptionBootfileName,
			dhcpv4.OptionVendorSpecificInformation,
			dhcpv4.OptionClassIdentifier,
		),
		dhcpv4.WithMessageType(dhcpv4.MessageTypeInform),
		dhcpv4.WithOption(OptVendorOptions(vendorOpts...)),
	)
	if err != nil {
		return nil, err
	}
	return PacketFor(d), nil
}

// NewReplyForInformList constructs an ACK for the INFORM[LIST] packet `inform`
// with additional options in `config`.
func NewReplyForInformList(inform *Packet, config ReplyConfig) (*Packet, error) {
	if config.DefaultImage == nil {
		return nil, errors.New("NewReplyForInformList: no default boot image ID set")
	}
	if config.Images == nil || len(config.Images) == 0 {
		return nil, errors.New("NewReplyForInformList: no boot images provided")
	}
	reply, err := dhcpv4.NewReplyFromRequest(&inform.DHCPv4)
	if err != nil {
		return nil, err
	}
	reply.ClientIPAddr = inform.ClientIPAddr
	reply.GatewayIPAddr = inform.GatewayIPAddr
	reply.ServerIPAddr = config.ServerIP
	reply.ServerHostName = config.ServerHostname

	reply.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeAck))
	reply.UpdateOption(dhcpv4.OptServerIdentifier(config.ServerIP))
	reply.UpdateOption(dhcpv4.OptClassIdentifier(AppleVendorID))

	// BSDP opts.
	vendorOpts := []dhcpv4.Option{
		OptMessageType(MessageTypeList),
		OptServerPriority(config.ServerPriority),
		OptDefaultBootImageID(config.DefaultImage.ID),
		OptBootImageList(config.Images...),
	}
	if config.SelectedImage != nil {
		vendorOpts = append(vendorOpts, OptSelectedBootImageID(config.SelectedImage.ID))
	}
	reply.UpdateOption(OptVendorOptions(vendorOpts...))
	return PacketFor(reply), nil
}

// NewReplyForInformSelect constructs an ACK for the INFORM[Select] packet
// `inform` with additional options in `config`.
func NewReplyForInformSelect(inform *Packet, config ReplyConfig) (*Packet, error) {
	if config.SelectedImage == nil {
		return nil, errors.New("NewReplyForInformSelect: no selected boot image ID set")
	}
	if config.Images == nil || len(config.Images) == 0 {
		return nil, errors.New("NewReplyForInformSelect: no boot images provided")
	}
	reply, err := dhcpv4.NewReplyFromRequest(&inform.DHCPv4)
	if err != nil {
		return nil, err
	}

	reply.ClientIPAddr = inform.ClientIPAddr
	reply.GatewayIPAddr = inform.GatewayIPAddr
	reply.ServerIPAddr = config.ServerIP
	reply.ServerHostName = config.ServerHostname
	reply.BootFileName = config.BootFileName

	reply.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeAck))
	reply.UpdateOption(dhcpv4.OptServerIdentifier(config.ServerIP))
	reply.UpdateOption(dhcpv4.OptClassIdentifier(AppleVendorID))

	// BSDP opts.
	reply.UpdateOption(OptVendorOptions(
		OptMessageType(MessageTypeSelect),
		OptSelectedBootImageID(config.SelectedImage.ID),
	))
	return PacketFor(reply), nil
}
