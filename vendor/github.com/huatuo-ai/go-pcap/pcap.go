package pcap

import (
	"context"
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"github.com/gopacket/gopacket"
	"golang.org/x/net/bpf"
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

// OpenLive opens a live capture using the platform's default capture path.
// The returned Handle implements gopacket.PacketDataSource.
func OpenLive(device string, snaplen int32, promisc bool, timeout time.Duration) (*Handle, error) {
	return OpenLiveWithContext(
		context.Background(),
		device,
		snaplen,
		promisc,
		timeout,
		DefaultSyscalls,
	)
}

// OpenLiveWithContext opens a live capture with cancellation and capture-path control.
// A nil context is treated as context.Background for compatibility.
func OpenLiveWithContext(
	ctx context.Context,
	device string,
	snaplen int32,
	promisc bool,
	timeout time.Duration,
	syscalls bool,
) (*Handle, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	return openLive(ctx, device, snaplen, promisc, timeout, syscalls)
}

// Listen simple one-step command to listen and send packets over a returned channel
func (h *Handle) Listen() chan Packet {
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
	instructions, err := compileLiveFilter(expr2)
	if err != nil {
		return fmt.Errorf("compile filter into instructions: %w", err)
	}
	raw, err := bpf.Assemble(instructions)
	if err != nil {
		return fmt.Errorf("assemble BPF instructions: %w", err)
	}
	return h.SetRawBPFFilter(raw)
}

func (h *Handle) SetRawBPFFilter(raw []bpf.RawInstruction) error {
	h.filter = raw
	return h.setFilter()
}

// getEndianness discover the endianness of our current system
func getEndianness() (binary.ByteOrder, error) {
	return binary.NativeEndian, nil
}

//nolint:unused
func htons(in uint16) uint16 {
	return (in<<8)&0xff00 | in>>8
}
