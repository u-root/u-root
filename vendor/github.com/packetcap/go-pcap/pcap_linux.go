package pcap

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"sync/atomic"
	"time"
	"unsafe"

	"golang.org/x/net/bpf"
	syscall "golang.org/x/sys/unix"

	"github.com/gopacket/gopacket"
	log "github.com/sirupsen/logrus"
)

const (
	//defaultFrameSize = 4096
	defaultFrameSize = 65632 //nolint:unused
	//defaultBlockNumbers = 128
	defaultBlockNumbers = 32
	//defaultBlockSize = defaultFrameSize * defaultBlockNumbers
	defaultBlockSize = 131072 //nolint:unused
	//defaultFramesPerBlock = defaultBlockSize / defaultFrameSize
	defaultFramesPerBlock = 32
	EthHlen               = 0x10
	// defaultSyscalls default setting for using syscalls
	defaultSyscalls     = false
	offsetToBlockStatus = 4 + 4

	tpacketAuxdataSize = 20
)

var (
	packetRALLSize           int32
	alignedTpacketHdrSize    int32
	alignedTpacketRALLSize   int32
	alignedTpacketAllHdrSize int32
)

type blockHeader struct {
	Version      uint32
	OffsetToPriv uint32
	H1           syscall.TpacketHdrV1
}

type captured struct {
	data []byte
	ci   gopacket.CaptureInfo
}

// Handle states
// Writer (OpenLive + Close) transitions:
//   - closed <-> open
//   - reading / polling -> canceling
//   - canceled / gone -> closed
//
// Reader (ReadPacketData) transitions:
//   - open <-> reading
//   - reading <-> polling
//   - reading -> gone
//   - canceling -> canceled
const (
	closed    uint32 = 0
	open      uint32 = 1
	reading   uint32 = 2
	polling   uint32 = 3
	canceling uint32 = 4
	canceled  uint32 = 5
	gone      uint32 = 6
)

const pollIntervalMs = 60 * 1000 // 1 minute

type Handle struct {
	// this must be first for atomic to behave nicely
	state           uint32
	syscalls        bool
	promiscuous     bool
	index           int
	iface           string
	snaplen         int32
	fd              int
	ring            []byte
	framePtr        int
	framesPerBuffer uint32
	frameIndex      uint32 //nolint:unused
	frameSize       uint32
	frameNumbers    uint32
	blockNumbers    int
	blockSize       int
	pollfd          []syscall.PollFd
	endian          binary.ByteOrder
	filter          []bpf.RawInstruction
	cache           []captured
}

func (h *Handle) ReadPacketData() (data []byte, ci gopacket.CaptureInfo, err error) {
	if !atomic.CompareAndSwapUint32(&h.state, open, reading) {
		return data, ci, io.EOF
	}
	defer func() {
		if !atomic.CompareAndSwapUint32(&h.state, reading, open) {
			if atomic.CompareAndSwapUint32(&h.state, canceling, canceled) {
				logger := log.WithFields(log.Fields{
					"iface": h.iface,
				})
				logger.Debugf("packet read was canceled")
			}
		}
	}()
	if h.syscalls {
		return h.readPacketDataSyscall()
	}
	// mmap can return multiple packets, so we can cache extras, and return if there are

	// if there already was one in the cache, return it
	if len(h.cache) > 0 {
		cap := h.cache[0]
		h.cache = h.cache[1:]
		return cap.data, cap.ci, nil
	}
	// there was not, so read a new one
	caps, err := h.readPacketDataMmap()
	if err != nil {
		return nil, ci, err
	}
	switch len(caps) {
	case 0:
		return nil, ci, nil
	case 1:
		return caps[0].data, caps[0].ci, nil
	}
	h.cache = caps
	cap := h.cache[0]
	h.cache = h.cache[1:]
	return cap.data, cap.ci, nil
}

func writeVLANTag(data []byte, tci, tpid uint16) ([]byte, []byte) {
	buf := make([]byte, 4)
	if tpid == 0 || binary.BigEndian.Uint16(data[12:14]) != 0x8100 {
		tpid = binary.BigEndian.Uint16(data[12:14])
		binary.BigEndian.PutUint16(data[12:14], 0x8100) // set ethernet frame type to VLAN
	}
	binary.BigEndian.PutUint16(buf[:2], tci)
	binary.BigEndian.PutUint16(buf[2:], tpid)
	return data, buf
}

func (h *Handle) readPacketDataSyscall() (data []byte, ci gopacket.CaptureInfo, err error) {
	b := make([]byte, h.snaplen)
	oob := make([]byte, syscall.CmsgSpace(tpacketAuxdataSize))
	n, _, _, _, err := syscall.Recvmsg(h.fd, b, oob, 0)
	if err != nil {
		return nil, ci, fmt.Errorf("error reading packets: %w", err)
	}

	var auxData syscall.TpacketAuxdata
	cmsgs, err := syscall.ParseSocketControlMessage(oob)
	if err != nil {
		return nil, ci, fmt.Errorf("error reading socket control messages: %w", err)
	}
	for _, cmsg := range cmsgs {
		if cmsg.Header.Level == syscall.SOL_PACKET && cmsg.Header.Type == syscall.PACKET_AUXDATA && cmsg.Header.Len >= tpacketAuxdataSize {
			auxData.Vlan_tci = binary.BigEndian.Uint16(cmsg.Data[len(cmsg.Data)-5 : len(cmsg.Data)-3])
			auxData.Vlan_tpid = binary.BigEndian.Uint16(cmsg.Data[len(cmsg.Data)-3:])
			break
		}
	}
	if auxData.Vlan_tci != 0 {
		var aux []byte
		b, aux = writeVLANTag(b, auxData.Vlan_tci, auxData.Vlan_tpid)
		b = append(append(b[:14], aux...), b[14:]...)
		n = n + 4
	}
	// TODO: add CaptureInfo, specifically:
	//    capture timestamp
	//    original packet length
	ci = gopacket.CaptureInfo{
		CaptureLength:  n,
		InterfaceIndex: h.index,
	}
	return b, ci, nil
}

func (h *Handle) readPacketDataMmap() ([]captured, error) {
	logger := log.WithFields(log.Fields{
		"method": "mmap",
		"iface":  h.iface,
	})
	logger.Debugf(
		"started: framesPerBuffer %d, blockSize %d, frameSize %d, frameNumbers %d, blockNumbers %d",
		h.framesPerBuffer,
		h.blockSize,
		h.frameSize,
		h.frameNumbers,
		h.blockNumbers,
	)
	// we check the bit setting on the pointer
	blockBase := h.framePtr * h.blockSize
	// add a loop, so that we do not just rely on the polling, but instead the actual flag bit
	flagIndex := blockBase + offsetToBlockStatus
	for atomic.LoadUint32(&h.state) == reading {
		logger.Debugf("checking for packet at block %d, buffer starting position %d, flagIndex %d ring pointer %p", h.framePtr, blockBase, flagIndex, h.ring)
		if h.ring[flagIndex]&syscall.TP_STATUS_USER == syscall.TP_STATUS_USER {
			return h.processMmapPackets(blockBase, flagIndex)
		}
		logger.Debugf("packet not ready at block %d position %d, polling via %#v", h.framePtr, blockBase, h.pollfd)
		var err error
		var val int
		// Just repeat Poll when we get timeout, do not even log anything.
		for err == nil && val == 0 {
			if !atomic.CompareAndSwapUint32(&h.state, reading, polling) {
				// the state is cancelling
				logger.Debugf("polling was canceled for ring %p", h.ring)
				return nil, io.EOF
			}
			// We need to have some timeout to eventually detect closed socket.
			// Listening for syscall.POLLERR and syscall.POLLNVAL events
			// does not seem to always do the job.
			val, err = syscall.Poll(h.pollfd, pollIntervalMs)
			if !atomic.CompareAndSwapUint32(&h.state, polling, reading) {
				// the state is cancelling
				logger.Debugf("polling was canceled for ring %p", h.ring)
				return nil, io.EOF
			}
		}
		logger.Debugf("poll returned val %v with pollfd %#v", val, h.pollfd)

		switch {
		case err != nil && err == syscall.EINTR:
			logger.Debug("poll returned with EINTR, continuing")
			continue
		case err != nil:
			logger.Errorf("error polling socket: %v", err)
			return nil, fmt.Errorf("error polling socket: %v", err)
		case val < 0:
			logger.Error("negative return value from polling socket")
			return nil, errors.New("negative return value from polling socket")
		case h.pollfd[0].Revents&syscall.POLLIN == syscall.POLLIN:
			continue
		case h.pollfd[0].Revents&syscall.POLLERR == syscall.POLLERR:
			// unknown error; we need to check the socket itself to find out
			logger.Debug("received POLLERR, checking actual error from socket")
			sockOptVal, err := syscall.GetsockoptInt(h.fd, syscall.SOL_SOCKET, syscall.SO_ERROR)
			if err != nil {
				logger.Errorf("could not get sockopt to check poll error; sockopt error: %v", err)
				return nil, fmt.Errorf("could not get sockopt to check poll error; sockopt error: %v", err)
			}
			if sockOptVal == int(syscall.ENETDOWN) {
				logger.Errorf("interface %s is down, marking handle as gone and returning", h.iface)
				atomic.StoreUint32(&h.state, gone)
				return nil, nil
			}
			// we have no idea what it was, so just return
			logger.Error("unknown error returned from socket")
			return nil, errors.New("unknown error returned from socket")
		case h.pollfd[0].Revents&syscall.POLLNVAL == syscall.POLLNVAL:
			logger.Error("socket closed")
			atomic.StoreUint32(&h.state, gone)
			return nil, io.EOF
		}
	}
	// if we got here, it was not 0, so EOF
	return nil, io.EOF
}

func (h *Handle) processMmapPackets(blockBase, flagIndex int) ([]captured, error) {
	logger := log.WithFields(log.Fields{
		"method": "mmap-process",
		"iface":  h.iface,
	})
	// read the header
	logger.Debugf("reading block header into b slice from position %d to position %d", blockBase, blockBase+h.blockSize)
	b := h.ring[blockBase : blockBase+h.blockSize]
	buf := bytes.NewBuffer(b[:])
	bHdr := blockHeader{}
	logger.Debugf("binary parsing block header of size %d", buf.Len())
	if err := binary.Read(buf, h.endian, &bHdr); err != nil {
		logger.Errorf("error reading block header: %v", err)
		return nil, fmt.Errorf("error reading block header: %v", err)
	}
	logger.Debugf("block header %#v", bHdr)
	// now we need to get the packets themselves
	numPkts := int(bHdr.H1.Num_pkts)
	packets := make([]captured, numPkts)

	nextOffset := bHdr.H1.Offset_to_first_pkt
	for i := 0; i < numPkts; i++ {
		hdr := syscall.Tpacket3Hdr{}
		logger.Debugf("packet number %d/%d at position %d in block", i, numPkts, nextOffset)
		b = b[nextOffset:]
		buf := bytes.NewBuffer(b[:alignedTpacketHdrSize])
		logger.Debugf("binary parsing packet header of size %d", buf.Len())
		if err := binary.Read(buf, h.endian, &hdr); err != nil {
			msg := fmt.Sprintf("error reading tpacket3 header on byte %d: %v", i, err)
			logger.Errorf(msg)
			return nil, fmt.Errorf(msg)
		}
		logger.Debugf("tpacket3 header %#v", hdr)
		nextOffset = hdr.Next_offset
		logger.Debugf("setting next offset to %d", nextOffset)

		// read the sockaddr_ll
		// unfortunately, we cannot do binary.Read() because syscall.SockaddrLinklayer has an embedded slice
		// so we have to read it manually
		// use b[hdr.Mac:hdr.Mac+alignedTpacketRALLSize] instead?
		sall, err := parseSocketAddrLinkLayer(b[alignedTpacketHdrSize:alignedTpacketAllHdrSize], h.endian)
		if err != nil {
			logger.Errorf("error parsing sockaddr_ll: %v", err)
			return nil, fmt.Errorf("error parsing sockaddr_ll for packet %d: %v", i, err)
		}

		ci := gopacket.CaptureInfo{
			Length:         int(hdr.Len),
			CaptureLength:  int(hdr.Snaplen),
			Timestamp:      time.Unix(int64(hdr.Sec), int64(hdr.Nsec)),
			InterfaceIndex: int(sall.Ifindex),
		}
		// We need to copy packet data because as soon as ReadPacketData returns,
		// the ring buffer could be un-mapped by Close. If the caller of ReadPacketData
		// does not process packet data quickly enough and a call to Handle.Close()
		// interleaves, it could find itself reading from invalid memory segments.
		// If you are using go-pcap with gopacket/gopacket, it is better to have packet
		// data copied here and then enable NoCopy for PacketSource, i.e.:
		//   packetSource := gopacket.NewPacketSource(...)
		//   packetSource.NoCopy = true
		data := make([]byte, hdr.Snaplen)
		copy(data, b[hdr.Mac:uint32(hdr.Mac)+hdr.Snaplen])
		if hdr.Hv1.Vlan_tci != 0 {
			var vlanTag []byte
			data, vlanTag = writeVLANTag(data, uint16(hdr.Hv1.Vlan_tci), uint16(hdr.Hv1.Vlan_tpid))
			data = append(data[:14], append(vlanTag, data[14:]...)...)
		}
		packets[i] = captured{
			ci:   ci,
			data: data,
		}

		logger.Debugf("raw packet for packet %d: %d\n ", i, data)
	}

	// indicate we are done with this frame, send back to the kernel
	logger.Debugf("returning block at pos %d to kernel", h.framePtr)
	h.ring[flagIndex] = syscall.TP_STATUS_KERNEL

	h.framePtr = (h.framePtr + 1) % h.blockNumbers
	logger.Debugf("final block: %d", h.framePtr)

	return packets, nil
}

// Close close sockets and release resources
func (h *Handle) Close() {
	logger := log.WithFields(log.Fields{
		"iface": h.iface,
	})
	// Wait for reader to finish before unmapping memory with the ring buffer.
	closeAttempts := 0
	for !atomic.CompareAndSwapUint32(&h.state, open, closed) {
		closeAttempts += 1
		if closeAttempts >= 1000 {
			// Stopping before we become an infinite loop
			logger.Tracef("Swapping on Stop tried for %d times, giving up now", closeAttempts)
			break
		}
		state := atomic.LoadUint32(&h.state)
		if state == canceled || state == gone {
			atomic.StoreUint32(&h.state, closed)
			break
		}
		if atomic.CompareAndSwapUint32(&h.state, reading, canceling) {
			logger.Debugf("cancelling ongoig packet read")
		}
		if atomic.CompareAndSwapUint32(&h.state, polling, canceling) {
			// When polling is interrupted it is safe to go ahead and unmap the ring buffer.
			// Reader will eventually detect canceled polling and will exit without accessing
			// the buffer anymore.
			logger.Debugf("cancelling ongoing socket polling; not waiting for poll to return")
			break
		}
	}
	if h.ring != nil {
		if err := syscall.Munmap(h.ring); err != nil {
			logger.Errorf("error unmapping mmap at %p ; nothing to do", h.ring)
		}
	}
	// close the socket
	if err := syscall.Close(h.fd); err != nil {
		logger.Errorf("error closing file descriptor %d ; nothing to do", h.fd)
	}
}

// set a classic BPF filter on the listener. filter must be compliant with
// tcpdump syntax.
func (h *Handle) setFilter() error {

	/*
	 * Try to install the kernel filter.
	 */
	prog := syscall.SockFprog{
		Len:    uint16(len(h.filter)),
		Filter: (*syscall.SockFilter)(unsafe.Pointer(&h.filter[0])),
	}

	if err := syscall.SetsockoptSockFprog(h.fd, syscall.SOL_SOCKET, syscall.SO_ATTACH_FILTER, &prog); err != nil {
		return fmt.Errorf("unable to set filter: %v", err)
	}
	return nil
}

func tpacketAlign(base int32) int32 {
	return (base + syscall.TPACKET_ALIGNMENT - 1) &^ (syscall.TPACKET_ALIGNMENT - 1)
}

func openLive(iface string, snaplen int32, promiscuous bool, timeout time.Duration, syscalls bool) (handle *Handle, _ error) {
	logger := log.WithFields(log.Fields{
		"iface":       iface,
		"snaplen":     snaplen,
		"promiscuous": promiscuous,
		"timeout":     timeout,
		"syscalls":    syscalls,
	})
	logger.Debug("started")
	h := Handle{
		// we start with it not open
		state:    closed,
		snaplen:  snaplen,
		syscalls: syscalls,
		iface:    iface,
	}
	// we need to know our endianness
	endianness, err := getEndianness()
	if err != nil {
		return nil, err
	}
	h.endian = endianness

	// because syscall package does not provide this
	rall := syscall.RawSockaddrLinklayer{}
	packetRALLSize = int32(unsafe.Sizeof(rall))
	alignedTpacketHdrSize = tpacketAlign(syscall.SizeofTpacket3Hdr)
	alignedTpacketRALLSize = tpacketAlign(packetRALLSize)
	alignedTpacketAllHdrSize = alignedTpacketHdrSize + alignedTpacketRALLSize

	// set up the socket - remember to switch to network socket order for the protocol int
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(htons(syscall.ETH_P_ALL)))
	if err != nil {
		logger.Errorf("failed opening raw socket: %v", err)
		return nil, fmt.Errorf("failed opening raw socket: %v", err)
	}
	h.fd = fd
	h.pollfd = []syscall.PollFd{{
		Fd:     int32(h.fd),
		Events: syscall.POLLIN | syscall.POLLERR | syscall.POLLNVAL}}
	if err := syscall.SetNonblock(fd, false); err != nil {
		return nil, fmt.Errorf("failed to set socket as blocking: %v", err)
	}
	if err = syscall.SetsockoptInt(fd, syscall.SOL_PACKET, syscall.PACKET_AUXDATA, 1); err != nil {
		return nil, fmt.Errorf("failed to set packet auxilary data: %w", err)
	}
	if iface != "" {
		// get our interface
		in, err := net.InterfaceByName(iface)
		if err != nil {
			logger.Errorf("unknown interface %s: %v", iface, err)
			return nil, fmt.Errorf("unknown interface %s: %v", iface, err)
		}
		// check the interface is up
		if in.Flags&net.FlagUp != net.FlagUp {
			logger.Errorf("interface %s is not up", iface)
			return nil, fmt.Errorf("interface %s is not up", iface)
		}
		h.index = in.Index

		// create the sockaddr_ll
		sa := syscall.SockaddrLinklayer{
			Protocol: htons(syscall.ETH_P_ALL),
			Ifindex:  in.Index,
		}
		// bind to it
		if err = syscall.Bind(fd, &sa); err != nil {
			return nil, fmt.Errorf("failed to bind")
		}
		if promiscuous {
			h.promiscuous = true
			mreq := syscall.PacketMreq{
				Ifindex: int32(in.Index),
				Type:    syscall.PACKET_MR_PROMISC,
			}
			if err = syscall.SetsockoptPacketMreq(fd, syscall.SOL_PACKET, syscall.PACKET_ADD_MEMBERSHIP, &mreq); err != nil {
				logger.Errorf("failed to set promiscuous for %s: %v", iface, err)
				return nil, fmt.Errorf("failed to set promiscuous for %s: %v", iface, err)
			}
		}
	}
	if !syscalls {
		if err = syscall.SetsockoptInt(fd, syscall.SOL_PACKET, syscall.PACKET_VERSION, syscall.TPACKET_V3); err != nil {
			logger.Errorf("failed to set TPACKET_V3: %v", err)
			return nil, fmt.Errorf("failed to set TPACKET_V3: %v", err)
		}
		// set up the ring
		var (
			frameSize           = uint32(tpacketAlign(syscall.SizeofTpacket3Hdr+EthHlen) + tpacketAlign(snaplen))
			pageSize            = syscall.Getpagesize()
			blockSize           = uint32(pageSize)
			blockNumbers uint32 = defaultBlockNumbers
		)
		for {
			if blockSize > frameSize {
				break
			}
			blockSize = blockSize << 1
		}
		// we use the default - for now

		framesPerBuffer := blockSize / frameSize
		frameNumbers := blockNumbers * framesPerBuffer

		tpreq := syscall.TpacketReq3{
			Block_size: blockSize,
			Block_nr:   blockNumbers,
			Frame_size: frameSize,
			Frame_nr:   frameNumbers,
		}
		logger.Debugf("creating mmap buffer with tpreq %#v", tpreq)
		if err = syscall.SetsockoptTpacketReq3(fd, syscall.SOL_PACKET, syscall.PACKET_RX_RING, &tpreq); err != nil {
			logger.Errorf("failed to set tpacket req: %v", err)
			return nil, fmt.Errorf("failed to set tpacket req: %v", err)
		}
		totalSize := int(tpreq.Block_size * tpreq.Block_nr)
		data, err := syscall.Mmap(fd, 0, totalSize, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
		if err != nil {
			logger.Errorf("error mmapping: %v", err)
			return nil, fmt.Errorf("error mmapping: %v", err)
		}
		logger.Infof("mmap buffer created at %p with size %d", data, len(data))
		h.framesPerBuffer = framesPerBuffer
		h.blockSize = int(blockSize)
		h.frameSize = frameSize
		h.frameNumbers = frameNumbers
		h.blockNumbers = int(blockNumbers)
		h.ring = data
		h.cache = make([]captured, 0, blockSize/frameSize)
	}
	atomic.StoreUint32(&h.state, open)
	return &h, nil
}

// parseSocketAddrLinkLayer parse byte data to get a RawSockAddrLinkLayer
func parseSocketAddrLinkLayer(b []byte, endian binary.ByteOrder) (*syscall.RawSockaddrLinklayer, error) {
	if len(b) < int(packetRALLSize) {
		return nil, fmt.Errorf("bytes of length %d shorter than mandated %d", len(b), packetRALLSize)
	}
	var addr [8]byte
	copy(addr[:], b[11:19])
	sall := syscall.RawSockaddrLinklayer{
		Family:   endian.Uint16(b[0:2]),
		Protocol: endian.Uint16(b[2:4]),
		Ifindex:  int32(endian.Uint32(b[4:8])),
		Hatype:   endian.Uint16(b[8:10]),
		Pkttype:  b[10],
		Halen:    b[11],
		Addr:     addr,
	}
	return &sall, nil
}
