package pcap

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
	"time"
	"unsafe"

	"github.com/gopacket/gopacket"
	"golang.org/x/net/bpf"

	"github.com/packetcap/go-pcap/filter"
)

const (
	// DefaultSyscalls whether the default is to use syscalls or not
	DefaultSyscalls = defaultSyscalls
)

// Packet a single packet returned by a listen call
type Packet struct {
	B     []byte
	Info  gopacket.CaptureInfo
	Error error
}

type BpfProgram struct {
	Len    uint16
	Filter *bpf.RawInstruction
}

// OpenLive open a live capture. Returns a Handle that implements https://godoc.org/github.com/gopacket/gopacket#PacketDataSource
// so you can pass it there.
func OpenLive(device string, snaplen int32, promiscuous bool, timeout time.Duration, syscalls bool) (handle *Handle, _ error) {
	return openLive(device, snaplen, promiscuous, timeout, syscalls)
}

// Listen simple one-step command to listen and send packets over a returned channel
func (h Handle) Listen() chan Packet {
	c := make(chan Packet, 50)
	go func() {
		for {
			b, ci, err := h.ReadPacketData()
			c <- Packet{
				B:     b,
				Info:  ci,
				Error: err,
			}
		}
	}()
	return c
}

// set a classic BPF filter on the listener. filter must be compliant with
// tcpdump syntax.
func (h *Handle) SetBPFFilter(expr string) error {
	expr2 := strings.TrimSpace(expr)
	// empty strings are not of interest
	if expr2 == "" {
		return nil
	}
	e := filter.NewExpression(expr2)
	if e == nil {
		return fmt.Errorf("no expression received for filter '%s'", expr)
	}
	f := e.Compile()
	instructions, err := f.Compile()
	if err != nil {
		return fmt.Errorf("failed to compile filter into instructions: %v", err)
	}
	raw, err := bpf.Assemble(instructions)
	if err != nil {
		return fmt.Errorf("bpf assembly failed: %v", err)
	}
	return h.SetRawBPFFilter(raw)
}

func (h *Handle) SetRawBPFFilter(raw []bpf.RawInstruction) error {
	h.filter = raw
	return h.setFilter()
}

// LinkType return the link type, compliant with pcap-linktype(7) and http://www.tcpdump.org/linktypes.html.
// For now, we just support Ethernet; some day we may support more
func (h Handle) LinkType() uint8 {
	return LinkTypeEthernet
}

// getEndianness discover the endianness of our current system
func getEndianness() (binary.ByteOrder, error) {
	buf := [2]byte{}
	*(*uint16)(unsafe.Pointer(&buf[0])) = uint16(0xABCD)

	switch buf {
	case [2]byte{0xCD, 0xAB}:
		return binary.LittleEndian, nil
	case [2]byte{0xAB, 0xCD}:
		return binary.BigEndian, nil
	default:
		return nil, errors.New("could not determine native endianness")
	}
}

// nolint: unused
func htons(in uint16) uint16 {
	return (in<<8)&0xff00 | in>>8
}
