package rtnetlink

import (
	"errors"
	"fmt"
	"net"

	"github.com/mdlayher/netlink"
	"github.com/mdlayher/netlink/nlenc"
	"golang.org/x/sys/unix"
)

var (
	// errInvalidLinkMessage is returned when a LinkMessage is malformed.
	errInvalidLinkMessage = errors.New("rtnetlink LinkMessage is invalid or too short")

	// errInvalidLinkMessageAttr is returned when link attributes are malformed.
	errInvalidLinkMessageAttr = errors.New("rtnetlink LinkMessage has a wrong attribute data length")
)

var _ Message = &LinkMessage{}

// A LinkMessage is a route netlink link message.
type LinkMessage struct {
	// Always set to AF_UNSPEC (0)
	Family uint16

	// Device Type
	Type uint16

	// Unique interface index, using a nonzero value with
	// NewLink will instruct the kernel to create a
	// device with the given index (kernel 3.7+ required)
	Index uint32

	// Contains device flags, see netdevice(7)
	Flags uint32

	// Change Flags, specifies which flags will be affected by the Flags field
	Change uint32

	// Attributes List
	Attributes *LinkAttributes
}

// MarshalBinary marshals a LinkMessage into a byte slice.
func (m *LinkMessage) MarshalBinary() ([]byte, error) {
	b := make([]byte, unix.SizeofIfInfomsg)

	b[0] = 0 //Family
	b[1] = 0 //reserved
	nlenc.PutUint16(b[2:4], m.Type)
	nlenc.PutUint32(b[4:8], m.Index)
	nlenc.PutUint32(b[8:12], m.Flags)
	nlenc.PutUint32(b[12:16], m.Change)

	if m.Attributes != nil {
		a, err := m.Attributes.MarshalBinary()
		if err != nil {
			return nil, err
		}

		return append(b, a...), nil
	}

	return b, nil
}

// UnmarshalBinary unmarshals the contents of a byte slice into a LinkMessage.
func (m *LinkMessage) UnmarshalBinary(b []byte) error {
	l := len(b)
	if l < unix.SizeofIfInfomsg {
		return errInvalidLinkMessage
	}

	m.Family = nlenc.Uint16(b[0:2])
	m.Type = nlenc.Uint16(b[2:4])
	m.Index = nlenc.Uint32(b[4:8])
	m.Flags = nlenc.Uint32(b[8:12])
	m.Change = nlenc.Uint32(b[12:16])

	if l > unix.SizeofIfInfomsg {
		m.Attributes = &LinkAttributes{}
		err := m.Attributes.UnmarshalBinary(b[16:])
		if err != nil {
			return err
		}
	}

	return nil
}

// rtMessage is an empty method to sattisfy the Message interface.
func (*LinkMessage) rtMessage() {}

// LinkService is used to retrieve rtnetlink family information.
type LinkService struct {
	c *Conn
}

// New creates a new interface using the LinkMessage information.
func (l *LinkService) New(req *LinkMessage) error {
	flags := netlink.Request | netlink.Create | netlink.Acknowledge | netlink.Excl
	_, err := l.c.Execute(req, unix.RTM_NEWLINK, flags)
	if err != nil {
		return err
	}

	return nil
}

// Delete removes an interface by index.
func (l *LinkService) Delete(index uint32) error {
	req := &LinkMessage{
		Index: index,
	}

	flags := netlink.Request | netlink.Acknowledge
	_, err := l.c.Execute(req, unix.RTM_DELLINK, flags)
	if err != nil {
		return err
	}

	return nil
}

// Get retrieves interface information by index.
func (l *LinkService) Get(index uint32) (LinkMessage, error) {
	req := &LinkMessage{
		Index: index,
	}

	flags := netlink.Request | netlink.DumpFiltered
	msg, err := l.c.Execute(req, unix.RTM_GETLINK, flags)
	if err != nil {
		return LinkMessage{}, err
	}

	if len(msg) != 1 {
		return LinkMessage{}, fmt.Errorf("too many/little matches, expected 1")
	}

	link := (msg[0]).(*LinkMessage)
	return *link, nil
}

// Set sets interface attributes according to the LinkMessage information.
func (l *LinkService) Set(req *LinkMessage) error {
	flags := netlink.Request | netlink.Acknowledge
	_, err := l.c.Execute(req, unix.RTM_SETLINK, flags)
	if err != nil {
		return err
	}

	return nil
}

// List retrieves all interfaces.
func (l *LinkService) List() ([]LinkMessage, error) {
	req := &LinkMessage{}

	flags := netlink.Request | netlink.Dump
	msgs, err := l.c.Execute(req, unix.RTM_GETLINK, flags)
	if err != nil {
		return nil, err
	}

	links := make([]LinkMessage, 0, len(msgs))
	for _, m := range msgs {
		link := (m).(*LinkMessage)
		links = append(links, *link)
	}

	return links, nil
}

// LinkAttributes contains all attributes for an interface.
type LinkAttributes struct {
	Address          net.HardwareAddr // Interface L2 address
	Broadcast        net.HardwareAddr // L2 broadcast address
	Name             string           // Device name
	MTU              uint32           // MTU of the device
	Type             uint32           // Link type
	QueueDisc        string           // Queueing discipline
	OperationalState OperationalState // Interface operation state
	Stats            *LinkStats       // Interface Statistics
	Info             *LinkInfo        // Detailed Interface Information
}

// OperationalState represents an interface's operational state.
type OperationalState uint8

// Constants that represent operational state of an interface
//
// Adapted from https://elixir.bootlin.com/linux/v4.19.2/source/include/uapi/linux/if.h#L166
const (
	OperStateUnknown        OperationalState = iota // status could not be determined
	OperStateNotPresent                             // down, due to some missing component (typically hardware)
	OperStateDown                                   // down, either administratively or due to a fault
	OperStateLowerLayerDown                         // down, due to lower-layer interfaces
	OperStateTesting                                // operationally down, in some test mode
	OperStateDormant                                // down, waiting for some external event
	OperStateUp                                     // interface is in a state to send and receive packets
)

// UnmarshalBinary unmarshals the contents of a byte slice into a LinkMessage.
func (a *LinkAttributes) UnmarshalBinary(b []byte) error {
	attrs, err := netlink.UnmarshalAttributes(b)
	if err != nil {
		return err
	}

	for _, attr := range attrs {
		switch attr.Type {
		case unix.IFLA_UNSPEC:
			//unused attribute
		case unix.IFLA_ADDRESS:
			l := len(attr.Data)
			if l < 4 || l > 32 {
				return errInvalidLinkMessageAttr
			}
			a.Address = attr.Data
		case unix.IFLA_BROADCAST:
			l := len(attr.Data)
			if l < 4 || l > 32 {
				return errInvalidLinkMessageAttr
			}
			a.Broadcast = attr.Data
		case unix.IFLA_IFNAME:
			a.Name = nlenc.String(attr.Data)
		case unix.IFLA_MTU:
			if len(attr.Data) != 4 {
				return errInvalidLinkMessageAttr
			}
			a.MTU = nlenc.Uint32(attr.Data)
		case unix.IFLA_LINK:
			if len(attr.Data) != 4 {
				return errInvalidLinkMessageAttr
			}
			a.Type = nlenc.Uint32(attr.Data)
		case unix.IFLA_QDISC:
			a.QueueDisc = nlenc.String(attr.Data)
		case unix.IFLA_OPERSTATE:
			if len(attr.Data) != 1 {
				return errInvalidLinkMessageAttr
			}
			a.OperationalState = OperationalState(nlenc.Uint8(attr.Data))
		case unix.IFLA_STATS:
			a.Stats = &LinkStats{}
			err := a.Stats.UnmarshalBinary(attr.Data)
			if err != nil {
				return err
			}
		case unix.IFLA_LINKINFO:
			a.Info = &LinkInfo{}
			err := a.Info.UnmarshalBinary(attr.Data)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// MarshalBinary marshals a LinkAttributes into a byte slice.
func (a *LinkAttributes) MarshalBinary() ([]byte, error) {
	attrs := []netlink.Attribute{
		{
			Type: unix.IFLA_UNSPEC,
			Data: nlenc.Uint16Bytes(0),
		},
		{
			Type: unix.IFLA_IFNAME,
			Data: nlenc.Bytes(a.Name),
		},
		{
			Type: unix.IFLA_MTU,
			Data: nlenc.Uint32Bytes(a.MTU),
		},
		{
			Type: unix.IFLA_LINK,
			Data: nlenc.Uint32Bytes(a.Type),
		},
		{
			Type: unix.IFLA_QDISC,
			Data: nlenc.Bytes(a.QueueDisc),
		},
	}

	if len(a.Address) != 0 {
		attrs = append(attrs, netlink.Attribute{
			Type: unix.IFLA_ADDRESS,
			Data: a.Address,
		})
	}

	if len(a.Broadcast) != 0 {
		attrs = append(attrs, netlink.Attribute{
			Type: unix.IFLA_BROADCAST,
			Data: a.Broadcast,
		})
	}

	if a.OperationalState != OperStateUnknown {
		attrs = append(attrs, netlink.Attribute{
			Type: unix.IFLA_OPERSTATE,
			Data: nlenc.Uint8Bytes(uint8(a.OperationalState)),
		})
	}

	if a.Info != nil {
		info, err := a.Info.MarshalBinary()
		if err != nil {
			return nil, err
		}
		attrs = append(attrs, netlink.Attribute{
			Type: unix.IFLA_LINKINFO,
			Data: info,
		})
	}

	return netlink.MarshalAttributes(attrs)
}

//LinkStats contains packet statistics
type LinkStats struct {
	RXPackets  uint32 // total packets received
	TXPackets  uint32 // total packets transmitted
	RXBytes    uint32 // total bytes received
	TXBytes    uint32 // total bytes transmitted
	RXErrors   uint32 // bad packets received
	TXErrors   uint32 // packet transmit problems
	RXDropped  uint32 // no space in linux buffers
	TXDropped  uint32 // no space available in linux
	Multicast  uint32 // multicast packets received
	Collisions uint32

	// detailed rx_errors:
	RXLengthErrors uint32
	RXOverErrors   uint32 // receiver ring buff overflow
	RXCRCErrors    uint32 // recved pkt with crc error
	RXFrameErrors  uint32 // recv'd frame alignment error
	RXFIFOErrors   uint32 // recv'r fifo overrun
	RXMissedErrors uint32 // receiver missed packet

	// detailed tx_errors
	TXAbortedErrors   uint32
	TXCarrierErrors   uint32
	TXFIFOErrors      uint32
	TXHeartbeatErrors uint32
	TXWindowErrors    uint32

	// for cslip etc
	RXCompressed uint32
	TXCompressed uint32

	RXNoHandler uint32 // dropped, no handler found
}

// UnmarshalBinary unmarshals the contents of a byte slice into a LinkMessage.
func (a *LinkStats) UnmarshalBinary(b []byte) error {
	l := len(b)
	if l != 92 && l != 96 {
		return fmt.Errorf("incorrect size, want: 92 or 96")
	}

	a.RXPackets = nlenc.Uint32(b[0:4])
	a.TXPackets = nlenc.Uint32(b[4:8])
	a.RXBytes = nlenc.Uint32(b[8:12])
	a.TXBytes = nlenc.Uint32(b[12:16])
	a.RXErrors = nlenc.Uint32(b[16:20])
	a.TXErrors = nlenc.Uint32(b[20:24])
	a.RXDropped = nlenc.Uint32(b[24:28])
	a.TXDropped = nlenc.Uint32(b[28:32])
	a.Multicast = nlenc.Uint32(b[32:36])
	a.Collisions = nlenc.Uint32(b[36:40])

	a.RXLengthErrors = nlenc.Uint32(b[40:44])
	a.RXOverErrors = nlenc.Uint32(b[44:48])
	a.RXCRCErrors = nlenc.Uint32(b[48:52])
	a.RXFrameErrors = nlenc.Uint32(b[52:56])
	a.RXFIFOErrors = nlenc.Uint32(b[56:60])
	a.RXMissedErrors = nlenc.Uint32(b[60:64])

	a.TXAbortedErrors = nlenc.Uint32(b[64:68])
	a.TXCarrierErrors = nlenc.Uint32(b[68:72])
	a.TXFIFOErrors = nlenc.Uint32(b[72:76])
	a.TXHeartbeatErrors = nlenc.Uint32(b[76:80])
	a.TXWindowErrors = nlenc.Uint32(b[80:84])

	a.RXCompressed = nlenc.Uint32(b[84:88])
	a.TXCompressed = nlenc.Uint32(b[88:92])

	if l == 96 {
		a.RXNoHandler = nlenc.Uint32(b[92:96])
	}

	return nil
}

const (
	iflaInfoKind      = 0x1
	iflaInfoData      = 0x2
	iflaInfoSlaveKind = 0x3
	iflaInfoSlaveData = 0x4
)

// LinkInfo contains data for specific network types
type LinkInfo struct {
	Kind      string // Driver name
	Data      []byte // Driver specific configuration stored as nested Netlink messages
	SlaveKind string // Slave driver name
	SlaveData []byte // Slave driver specific configuration
}

// UnmarshalBinary unmarshals the contents of a byte slice into a LinkInfo.
func (i *LinkInfo) UnmarshalBinary(b []byte) error {
	attrs, err := netlink.UnmarshalAttributes(b)
	if err != nil {
		return err
	}

	for _, attr := range attrs {
		switch attr.Type {
		case iflaInfoKind:
			i.Kind = nlenc.String(attr.Data)
		case iflaInfoSlaveKind:
			i.SlaveKind = nlenc.String(attr.Data)
		case iflaInfoData:
			i.Data = attr.Data
		case iflaInfoSlaveData:
			i.SlaveData = attr.Data
		}
	}

	return nil
}

// MarshalBinary marshals a LinkInfo into a byte slice.
func (i *LinkInfo) MarshalBinary() ([]byte, error) {
	attrs := []netlink.Attribute{
		{
			Type: iflaInfoKind,
			Data: nlenc.Bytes(i.Kind),
		},
		{
			Type: iflaInfoData,
			Data: i.Data,
		},
	}

	if len(i.SlaveData) > 0 {
		attrs = append(attrs,
			netlink.Attribute{
				Type: iflaInfoSlaveKind,
				Data: nlenc.Bytes(i.SlaveKind),
			},
			netlink.Attribute{
				Type: iflaInfoSlaveData,
				Data: i.SlaveData,
			},
		)
	}

	return netlink.MarshalAttributes(attrs)
}
