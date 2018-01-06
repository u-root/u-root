package dhcp4

import (
	"net"
	"strings"

	"github.com/u-root/dhcp4/util"
)

const (
	minPacketLen = 236

	// Maximum length of the CHAddr (client hardware address) according to
	// RFC 2131, Section 2. This is the link-layer destination a server
	// must send responses to.
	chaddrLen = 16
)

var (
	magicCookie = []byte{99, 130, 83, 99}
)

const (
	flagBroadcast = 1 << 15
)

// Packet is a DHCPv4 packet as described in RFC 2131 Section 2.
type Packet struct {
	Op            OpCode
	HType         uint8
	Hops          uint8
	TransactionID [4]byte
	Secs          uint16
	Broadcast     bool

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

	ServerName string
	BootFile   string
	Options    Options
}

func NewPacket(op OpCode) *Packet {
	return &Packet{
		Op:      op,
		HType:   1, /* ethernet */
		Options: make(Options),
	}
}

func (p *Packet) writeIP(b *util.Buffer, ip net.IP) {
	var zeros [net.IPv4len]byte
	if ip == nil {
		b.WriteBytes(zeros[:])
	} else {
		b.WriteBytes(ip[:net.IPv4len])
	}
}

func (p *Packet) MarshalBinary() ([]byte, error) {
	b := util.NewBuffer(make([]byte, 0, minPacketLen))
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

	p.writeIP(b, p.CIAddr)
	p.writeIP(b, p.YIAddr)
	p.writeIP(b, p.SIAddr)
	p.writeIP(b, p.GIAddr)
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
	b.WriteBytes(magicCookie)

	p.Options.Marshal(b)
	// TODO pad to 272 bytes for really old crap.
	return b.Data(), nil
}

func (p *Packet) UnmarshalBinary(q []byte) error {
	b := util.NewBuffer(q)
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

	// Read the cookie and then fucking ignore it.
	var cookie [4]byte
	b.ReadBytes(cookie[:])

	return (&p.Options).Unmarshal(b)
}
