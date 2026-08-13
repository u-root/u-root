//go:build darwin || freebsd

package pcap

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"golang.org/x/net/bpf"
	"golang.org/x/sys/unix"

	"github.com/gopacket/gopacket"
	log "github.com/sirupsen/logrus"
)

const (
	enable = 1
	// defaultSyscalls default setting for using syscalls
	defaultSyscalls = true
)

type Handle struct {
	context     context.Context
	syscalls    bool
	close       sync.Once
	closed      atomic.Bool
	promiscuous bool //nolint: unused
	timeout     time.Duration
	index       int
	fd          int
	buf         []byte
	endian      binary.ByteOrder
	filter      []bpf.RawInstruction
	linkType    uint32
}

type BpfProgram struct {
	Len    uint32
	Filter *bpf.RawInstruction
}

func (h *Handle) ReadPacketData() (data []byte, ci gopacket.CaptureInfo, err error) {
	if h.syscalls {
		return h.readPacketDataSyscall()
	}
	return h.readPacketDataMmap()
}

func (h *Handle) readPacketDataSyscall() (data []byte, ci gopacket.CaptureInfo, err error) {
	// must memset the buffer
	h.buf = make([]byte, len(h.buf))

	var pipefd [2]int
	if err := unix.Pipe(pipefd[:]); err != nil {
		return nil, ci, fmt.Errorf("pipe: %w", err)
	}
	rfd := pipefd[0]
	wfd := pipefd[1]

	defer func() {
		if closeErr := unix.Close(rfd); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close read pipe: %w", closeErr))
		}
	}()
	defer func() {
		if closeErr := unix.Close(wfd); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close write pipe: %w", closeErr))
		}
	}()

	// Make pipe non-blocking
	if err := unix.SetNonblock(rfd, true); err != nil {
		return nil, ci, fmt.Errorf("set read pipe non-blocking: %w", err)
	}
	if err := unix.SetNonblock(wfd, true); err != nil {
		return nil, ci, fmt.Errorf("set write pipe non-blocking: %w", err)
	}

	// Goroutine to signal cancellation
	done := make(chan struct{})
	go func() {
		select {
		case <-h.context.Done():
			// Write any value to wake poll
			_, _ = unix.Write(wfd, []byte{1})
		case <-done:
		}
	}()
	defer close(done)

	// pollfd to handle events, like idle timeout or context message
	pfd := []unix.PollFd{
		{
			Fd:     int32(h.fd),
			Events: unix.POLLIN,
		},
		{
			Fd:     int32(rfd),
			Events: unix.POLLIN,
		},
	}

	ms := -1
	if h.timeout > 0 {
		ms = int(h.timeout.Milliseconds())
	}
	n, err := unix.Poll(pfd, ms)
	if err != nil {
		if errors.Is(err, unix.EINTR) {
			return nil, ci, h.context.Err()
		}
		return nil, ci, err
	}
	if n == 0 {
		return nil, ci, context.DeadlineExceeded
	}
	// Context canceled → eventfd readable
	if pfd[1].Revents&unix.POLLIN != 0 {
		return nil, ci, h.context.Err()
	}

	read, err := unix.Read(h.fd, h.buf)
	if err != nil {
		return nil, ci, fmt.Errorf("read packet: %w", err)
	}
	if read <= 0 {
		return nil, ci, fmt.Errorf("read no packets")
	}
	// separate the header and packet body
	hdr := unix.BpfHdr{}
	buf := bytes.NewBuffer(h.buf[:unix.SizeofBpfHdr])
	err = binary.Read(buf, h.endian, &hdr)
	if err != nil {
		return nil, ci, fmt.Errorf("read BPF header: %w", err)
	}
	ci = gopacket.CaptureInfo{
		Timestamp:      time.Now(),
		CaptureLength:  int(hdr.Caplen),
		Length:         int(hdr.Datalen),
		InterfaceIndex: h.index,
	}
	return h.buf[hdr.Hdrlen : uint32(hdr.Hdrlen)+hdr.Caplen], ci, nil
}

func (h *Handle) readPacketDataMmap() (data []byte, ci gopacket.CaptureInfo, err error) {
	return nil, ci, errors.New("mmap unsupported on Darwin")
}

// Close close sockets and release resources
// Close is idempotent, and uses sync.Once to ensure it only runs once.
func (h *Handle) Close() {
	// close the socket
	h.close.Do(func() {
		_ = unix.Close(h.fd)
		h.closed.Store(true)
	})
}

// set a classic BPF filter on the listener. filter must be compliant with
// tcpdump syntax.
func (h *Handle) setFilter() error {
	if len(h.filter) == 0 {
		return errors.New("cannot install empty BPF filter")
	}

	/*
	 * Try to install the kernel filter.
	 */
	prog := BpfProgram{
		Len: uint32(len(h.filter)),
		// #nosec G103 -- BIOCSETF synchronously copies this non-empty Go-owned slice.
		Filter: (*bpf.RawInstruction)(unsafe.Pointer(&h.filter[0])),
	}
	// #nosec G103 -- ioctl receives the address of prog only for this synchronous call.
	if err := ioctlPtr(h.fd, unix.BIOCSETF, unsafe.Pointer(&prog)); err != nil {
		return fmt.Errorf("set BPF filter: %w", err)
	}

	return nil
}

func openLive(ctx context.Context, iface string, snaplen int32, promiscuous bool, timeout time.Duration, syscalls bool) (handle *Handle, _ error) {
	var (
		fd  = -1
		err error
	)
	logger := log.WithFields(log.Fields{
		"iface":       iface,
		"snaplen":     snaplen,
		"promiscuous": promiscuous,
		"timeout":     timeout,
		"syscalls":    syscalls,
	})
	logger.Debug("started")
	h := Handle{
		context:  ctx,
		syscalls: syscalls,
	}
	// we need to know our endianness
	endianness, err := getEndianness()
	if err != nil {
		return nil, err
	}
	h.endian = endianness

	// open the bpf device
	for i := 0; i < 255; i++ {
		dev := fmt.Sprintf("/dev/bpf%d", i)
		fd, err = unix.Open(dev, unix.O_RDWR, 0o000)
		if fd > -1 {
			break
		}
		if errors.Is(err, unix.EBUSY) {
			continue
		}
		return nil, fmt.Errorf("open device %s: %w", dev, err)
	}
	if fd <= -1 {
		return nil, errors.New("failed to get valid bpf device")
	}
	h.fd = fd

	// set the options
	if err = SetBpfInterface(fd, iface); err != nil {
		return nil, fmt.Errorf("set BPF interface: %w", err)
	}
	if err = SetBpfHeadercmpl(fd, enable); err != nil {
		return nil, fmt.Errorf("set BPF header complete option: %w", err)
	}
	if err = SetBpfMonitor(fd, enable); err != nil {
		return nil, fmt.Errorf("set BPF monitor option: %w", err)
	}
	if err = SetBpfImmediate(fd, enable); err != nil {
		return nil, fmt.Errorf("set BPF immediate return option: %w", err)
	}
	size, err := BpfBuflen(fd)
	if err != nil {
		return nil, fmt.Errorf("read buffer length: %w", err)
	}
	h.buf = make([]byte, size)

	linkType, err := getLinkType(fd)
	if err != nil {
		return nil, fmt.Errorf("get link type: %w", err)
	}
	h.linkType = linkType

	h.timeout = timeout

	return &h, nil
}

// because they deprecated all of the below from "syscall" and redirected to "golang.org/x/net/bpf" but did not
// create a replacement. Sigh.

type ivalue struct {
	name  [unix.IFNAMSIZ]byte
	value int16
}

func SetBpfInterface(fd int, name string) error {
	var iv ivalue
	copy(iv.name[:], name)
	// #nosec G103 -- ioctl reads iv only for the duration of this call.
	return ioctlPtr(fd, unix.BIOCSETIF, unsafe.Pointer(&iv))
}

func SetBpfHeadercmpl(fd, m int) error {
	return unix.IoctlSetPointerInt(fd, unix.BIOCSHDRCMPLT, m)
}

func SetBpfImmediate(fd, m int) error {
	return unix.IoctlSetPointerInt(fd, unix.BIOCIMMEDIATE, m)
}

func SetBpfMonitor(fd, m int) error {
	return unix.IoctlSetPointerInt(fd, unix.BIOCSSEESENT, m)
}

func BpfBuflen(fd int) (int, error) {
	return unix.IoctlGetInt(fd, unix.BIOCGBLEN)
}

func ioctlPtr(fd, arg int, valPtr unsafe.Pointer) error {
	//nolint:staticcheck // unix.SYS_IOCTL is deprecated, but golang does not provide a better alternative
	// as of this writing for passing pointers
	_, _, errno := unix.RawSyscall(unix.SYS_IOCTL, uintptr(fd), uintptr(arg), uintptr(valPtr))
	if errno != 0 {
		return fmt.Errorf("ioctl %#x: %w", arg, errno)
	}
	return nil
}

func getLinkType(fd int) (uint32, error) {
	linkType, err := unix.IoctlGetInt(fd, unix.BIOCGDLT)
	if err != nil {
		return 0xffffffff, fmt.Errorf("get link type: %w", err)
	}
	return uint32(linkType), nil
}

// LinkType return the link type, compliant with pcap-linktype(7) and http://www.tcpdump.org/linktypes.html.
// For now, we just support Null and Ethernet; some day we may support more
func (h *Handle) LinkType() uint32 {
	return h.linkType
}
