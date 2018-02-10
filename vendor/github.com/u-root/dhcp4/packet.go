package dhcp4

import (
	"fmt"
	"net"
	"strings"

	"github.com/u-root/dhcp4/internal/buffer"
)

const (
	minPacketLen = 236

	// Maximum length of the CHAddr (client hardware address) according to
	// RFC 2131, Section 2. This is the link-layer destination a server
	// must send responses to.
	chaddrLen = 16

	// flagBroadcast is the broadcast bit in the flag field as defined by
	// RFC 2131, Section 2, Figure 2.
	flagBroadcast = 1 << 15
)

var (
	// This is the magic cookie for BOOTP/DHCP packets as defined in RFC
	// 1497 and RFC 2131, Section 3.
	magicCookie = [4]byte{99, 130, 83, 99}
)

// Packet is a DHCPv4 packet as described in RFC 2131, Section 2.
type Packet struct {
	// Op is the BOOTP message op code / message type.
	//
	// This is not to be confused with the DHCP message type, which is
	// defined as an option value.
	Op OpCode

	// HType is the hardware type.
	//
	// The possible values are listed in the IANA ARP assigned numbers.
	HType uint8

	// Hops is the number of hops this packet has taken.
	Hops uint8

	// TransactionID is a random number used to associate server responses
	// with client requests.
	TransactionID [4]byte

	// Secs is the number of seconds elapsed since the client began address
	// acquisition or renewal process.
	Secs uint16

	// Broadcast is the broadcast flag of the flags field.
	Broadcast bool

	// Client IP address.
	CIAddr net.IP

	// Your IP address.
	YIAddr net.IP

	// Server IP address.
	SIAddr net.IP

	// Gateway IP address.
	GIAddr net.IP

	// Client hardware address.
	CHAddr net.HardwareAddr

	// ServerName is an optional server host name.
	ServerName string

	// BootFile is a fully qualified directory path to the boot file.
	BootFile string

	// Options is the list of vendor-specific extensions.
	Options Options
}

// NewPacket returns a new DHCP packet with the given op code.
func NewPacket(op OpCode) *Packet {
	return &Packet{
		Op:      op,
		HType:   1, /* ethernet */
		Options: make(Options),
	}
}

func writeIP(b *buffer.Buffer, ip net.IP) {
	var zeros [net.IPv4len]byte
	if ip == nil {
		b.WriteBytes(zeros[:])
	} else {
		b.WriteBytes(ip[:net.IPv4len])
	}
}

// MarshalBinary writes the packet to binary.
func (p *Packet) MarshalBinary() ([]byte, error) {
	b := buffer.New(make([]byte, 0, minPacketLen))
	b.Write8(uint8(p.Op))
	b.Write8(p.HType)

	// HLen
	b.Write8(uint8(len(p.CHAddr)))
	b.Write8(p.Hops)
	b.WriteBytes(p.TransactionID[:])
	b.Write16(p.Secs)

	var flags uint16
	if p.Broadcast {
		flags |= flagBroadcast
	}
	b.Write16(flags)

	writeIP(b, p.CIAddr)
	writeIP(b, p.YIAddr)
	writeIP(b, p.SIAddr)
	writeIP(b, p.GIAddr)
	copy(b.WriteN(chaddrLen), p.CHAddr)

	var sname [64]byte
	copy(sname[:], []byte(p.ServerName))
	sname[len(p.ServerName)] = 0
	b.WriteBytes(sname[:])

	var file [128]byte
	copy(file[:], []byte(p.BootFile))
	file[len(p.BootFile)] = 0
	b.WriteBytes(file[:])

	// The magic cookie.
	b.WriteBytes(magicCookie[:])

	p.Options.Marshal(b)
	// TODO pad to 272 bytes for really old crap.
	return b.Data(), nil
}

// UnmarshalBinary reads the packet from binary.
func (p *Packet) UnmarshalBinary(q []byte) error {
	b := buffer.New(q)
	if b.Len() < minPacketLen {
		return ErrInvalidPacket
	}

	p.Op = OpCode(b.Read8())
	p.HType = b.Read8()
	hlen := b.Read8()
	p.Hops = b.Read8()
	b.ReadBytes(p.TransactionID[:])
	p.Secs = b.Read16()

	flags := b.Read16()
	if flags&flagBroadcast != 0 {
		p.Broadcast = true
	}

	p.CIAddr = make(net.IP, net.IPv4len)
	b.ReadBytes(p.CIAddr)
	p.YIAddr = make(net.IP, net.IPv4len)
	b.ReadBytes(p.YIAddr)
	p.SIAddr = make(net.IP, net.IPv4len)
	b.ReadBytes(p.SIAddr)
	p.GIAddr = make(net.IP, net.IPv4len)
	b.ReadBytes(p.GIAddr)

	if hlen > chaddrLen {
		hlen = chaddrLen
	}
	// Always read 16 bytes, but only use hlen of them.
	p.CHAddr = make(net.HardwareAddr, chaddrLen)
	b.ReadBytes(p.CHAddr)
	p.CHAddr = p.CHAddr[:hlen]

	var sname [64]byte
	b.ReadBytes(sname[:])
	length := strings.Index(string(sname[:]), "\x00")
	if length == -1 {
		length = 64
	}
	p.ServerName = string(sname[:length])

	var file [128]byte
	b.ReadBytes(file[:])
	length = strings.Index(string(file[:]), "\x00")
	if length == -1 {
		length = 128
	}
	p.BootFile = string(file[:length])

	var cookie [4]byte
	b.ReadBytes(cookie[:])
	if cookie != magicCookie {
		return fmt.Errorf("malformed DHCP packet: got magic cookie %v, want %v", cookie[:], magicCookie[:])
	}

	return (&p.Options).Unmarshal(b)
}
