package tftp

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

// NewServer creates TFTP server. It requires two functions to handle
// read and write requests.
// In case nil is provided for read or write handler the respective
// operation is disabled.
func NewServer(readHandler func(filename string, rf io.ReaderFrom) error,
	writeHandler func(filename string, wt io.WriterTo) error) *Server {
	s := &Server{
		mu:                &sync.Mutex{},
		timeout:           defaultTimeout,
		retries:           defaultRetries,
		smartBlock:        true,
		packetReadTimeout: 100 * time.Millisecond,
		readHandler:       readHandler,
		writeHandler:      writeHandler,
		wg:                &sync.WaitGroup{},
	}
	s.cancel, s.cancelFn = context.WithCancel(context.Background())
	return s
}

// RequestPacketInfo provides a method of getting the local IP address
// that is handling a UDP request.  It relies for its accuracy on the
// OS providing methods to inspect the underlying UDP and IP packets
// directly.
type RequestPacketInfo interface {
	// LocalIP returns the IP address we are servicing the request on.
	// If it is unable to determine what address that is, the returned
	// net.IP will be nil.
	LocalIP() net.IP
}

// Server is an instance of a TFTP server
type Server struct {
	mu           *sync.Mutex
	readHandler  func(filename string, rf io.ReaderFrom) error
	writeHandler func(filename string, wt io.WriterTo) error
	hook         Hook
	backoff      backoffFunc
	conn         net.PacketConn
	conn6        *ipv6.PacketConn
	conn4        *ipv4.PacketConn
	wg           *sync.WaitGroup
	timeout      time.Duration
	retries      int
	maxBlockLen  int
	smartBlock   bool
	sendAEnable  bool /* senderAnticipate enable by server */
	sendAWinSz   uint
	// Single port fields
	singlePort        bool
	handlers          map[string]chan []byte
	packetReadTimeout time.Duration
	cancel            context.Context
	cancelFn          context.CancelFunc
}

// TransferStats contains details about a single TFTP transfer
type TransferStats struct {
	RemoteAddr              net.IP
	Filename                string
	Tid                     int
	SenderAnticipateEnabled bool
	Mode                    string
	Opts                    options
	Duration                time.Duration
	DatagramsSent           int
	DatagramsAcked          int
}

// Hook is an interface used to provide the server with success and failure hooks
type Hook interface {
	OnSuccess(stats TransferStats)
	OnFailure(stats TransferStats, err error)
}

// SetAnticipate provides an experimental feature in which when a packets
// is requested the server will keep sending a number of packets before
// checking whether an ack has been received. It improves tftp downloading
// speed by a few times.
// The argument winsz specifies how many packets will be sent before
// waiting for an ack packet.
// When winsz is bigger than 1, the feature is enabled, and the server
// runs through a different experimental code path. When winsz is 0 or 1,
// the feature is disabled.
func (s *Server) SetAnticipate(winsz uint) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if winsz > 1 {
		s.sendAEnable = true
		s.sendAWinSz = winsz
	} else {
		s.sendAEnable = false
		s.sendAWinSz = 1
	}
}

// SetHook sets the Hook for success and failure of transfers
func (s *Server) SetHook(hook Hook) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.hook = hook
}

// EnableSinglePort enables an experimental mode where the server will
// serve all connections on port 69 only. There will be no random TIDs
// on the server side.
//
// Enabling this will negatively impact performance
func (s *Server) EnableSinglePort() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.singlePort = true
	s.handlers = make(map[string]chan []byte)
	if s.maxBlockLen == 0 {
		s.maxBlockLen = blockLength
	}
}

// SetTimeout sets maximum time server waits for single network
// round-trip to succeed.
// Default is 5 seconds.
func (s *Server) SetTimeout(t time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if t <= 0 {
		s.timeout = defaultTimeout
	} else {
		s.timeout = t
	}
}

// SetBlockSize sets the maximum size of an individual data block.
// This must be a value between 512 (the default block size for TFTP)
// and 65465 (the max size a UDP packet payload can be).
//
// This is an advisory value -- it will be clamped to the smaller of
// the block size the client wants and the MTU of the interface being
// communicated over munis overhead.
func (s *Server) SetBlockSize(i int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if i > 512 && i < 65465 {
		s.maxBlockLen = i
	}
}

// SetBlockSizeNegotiation controls MTU-based blksize negotiation.
//
// When enabled (default), the server clamps negotiated block size to the
// interface MTU (minus protocol overhead).
//
// When disabled, the server honors client-requested blksize (still bounded
// by RFC limits and optional SetBlockSize cap).
func (s *Server) SetBlockSizeNegotiation(enabled bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.smartBlock = enabled
}

// SetRetries sets maximum number of attempts server made to transmit a
// packet.
// Default is 5 attempts.
func (s *Server) SetRetries(count int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if count < 1 {
		s.retries = defaultRetries
	} else {
		s.retries = count
	}
}

// SetBackoff sets a user provided function that is called to provide a
// backoff duration prior to retransmitting an unacknowledged packet.
func (s *Server) SetBackoff(h backoffFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.backoff = h
}

// ListenAndServe binds to address provided and start the server.
// ListenAndServe returns when Shutdown is called.
func (s *Server) ListenAndServe(addr string) error {
	a, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", a)
	if err != nil {
		return err
	}
	return s.Serve(conn)
}

// Serve starts server provided already opened UDP connection. It is
// useful for the case when you want to run server in separate goroutine
// but still want to be able to handle any errors opening connection.
// Serve returns when Shutdown is called.
func (s *Server) Serve(conn net.PacketConn) error {
	laddr := conn.LocalAddr()
	host, _, err := net.SplitHostPort(laddr.String())
	if err != nil {
		return err
	}
	s.mu.Lock()
	s.conn = conn
	s.mu.Unlock()
	// Having seperate control paths for IP4 and IP6 is annoying,
	// but necessary at this point.
	addr := net.ParseIP(host)
	if addr == nil {
		return fmt.Errorf("failed to determine IP class of listening address")
	}

	if conn, ok := s.conn.(*net.UDPConn); ok {
		if addr.To4() != nil {
			s.conn4 = ipv4.NewPacketConn(conn)
			if err := s.conn4.SetControlMessage(ipv4.FlagDst|ipv4.FlagInterface, true); err != nil {
				s.conn4 = nil
			}
		} else {
			s.conn6 = ipv6.NewPacketConn(conn)
			if err := s.conn6.SetControlMessage(ipv6.FlagDst|ipv6.FlagInterface, true); err != nil {
				s.conn6 = nil
			}
		}
	}

	if s.singlePort {
		return s.singlePortProcessRequests()
	}
	for {
		select {
		case <-s.cancel.Done():
			// Stop server because Shutdown was called
			return nil
		default:
			var err error
			if s.conn4 != nil {
				err = s.processRequest4()
			} else if s.conn6 != nil {
				err = s.processRequest6()
			} else {
				err = s.processRequest()
			}
			if err != nil && s.hook != nil {
				s.hook.OnFailure(TransferStats{
					SenderAnticipateEnabled: s.sendAEnable,
				}, err)
			}
		}
	}
}

// Yes, I don't really like having separate IPv4 and IPv6 variants,
// bit we are relying on the low-level packet control channel info to
// get a reliable source address, and those have different types and
// the struct itself is not easily interface-ized or embedded.
//
// If control is nil for whatever reason (either things not being
// implemented on a target OS or whatever other reason), localIP
// (and hence LocalIP()) will return a nil IP address.
func (s *Server) processRequest4() error {
	buf := make([]byte, datagramLength)
	cnt, control, srcAddr, err := s.conn4.ReadFrom(buf)
	if err != nil {
		return nil
	}
	maxSz := blockLength
	var localAddr net.IP
	if control != nil {
		localAddr = control.Dst
		if intf, err := net.InterfaceByIndex(control.IfIndex); err == nil {
			// mtu - ipv4 overhead - udp overhead
			maxSz = intf.MTU - 28
		}
	}
	return s.handlePacket(localAddr, srcAddr.(*net.UDPAddr), buf, cnt, maxSz, nil)
}

func (s *Server) processRequest6() error {
	buf := make([]byte, datagramLength)
	cnt, control, srcAddr, err := s.conn6.ReadFrom(buf)
	if err != nil {
		return nil
	}
	maxSz := blockLength
	var localAddr net.IP
	if control != nil {
		localAddr = control.Dst
		if intf, err := net.InterfaceByIndex(control.IfIndex); err == nil {
			// mtu - ipv6 overhead - udp overhead
			maxSz = intf.MTU - 48
		}
	}
	return s.handlePacket(localAddr, srcAddr.(*net.UDPAddr), buf, cnt, maxSz, nil)
}

// Fallback if we had problems opening a ipv4/6 control channel
func (s *Server) processRequest() error {
	buf := make([]byte, datagramLength)
	cnt, srcAddr, err := s.conn.ReadFrom(buf)
	if err != nil {
		return fmt.Errorf("reading UDP: %v", err)
	}
	return s.handlePacket(nil, srcAddr.(*net.UDPAddr), buf, cnt, blockLength, nil)
}

// Shutdown make server stop listening for new requests, allows
// server to finish outstanding transfers and stops server.
// Shutdown blocks until all outstanding requests are processed or timed out.
// Calling Shutdown from the handler or hook might cause deadlock.
func (s *Server) Shutdown() {
	if s.singlePort {
		s.cancelFn()
		s.mu.Lock()
		conn := s.conn
		s.mu.Unlock()
		if conn != nil {
			// Wake the single-port read loop so it can observe cancellation.
			_ = conn.SetReadDeadline(time.Now())
		}
		return
	}
	s.mu.Lock()
	// Connection could not exist if Serve or
	// ListenAndServe was never called.
	if s.conn != nil {
		s.conn.Close()
	}
	s.mu.Unlock()
	s.cancelFn()
	s.wg.Wait()
}

func (s *Server) handlePacket(localAddr net.IP, remoteAddr *net.UDPAddr, buffer []byte, n, maxBlockLen int, listener chan []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Cope with packets received on the broadcast address
	// We can't use this address as the source address in responses
	// so fallback to the OS default.
	if localAddr.Equal(net.IPv4bcast) {
		localAddr = net.IPv4zero
	}
	if !s.smartBlock {
		// Disable MTU-based clamping.
		maxBlockLen = 65464
	}

	// handlePacket is always called with maxBlockLen = blockLength (above, in processRequest).
	// As a result, the block size would always be capped at 512 bytes, even when the tftp
	// client indicated to use a larger value.  So override that value.  And make sure to
	// use that value below, when allocating buffers.  (Happening on Windows Server 2016.)
	//	if s.maxBlockLen > 0 {
	if s.maxBlockLen > 0 && s.maxBlockLen < maxBlockLen {
		maxBlockLen = s.maxBlockLen
	}
	if maxBlockLen < blockLength {
		maxBlockLen = blockLength
	}
	p, err := parsePacket(buffer[:n])
	if err != nil {
		return err
	}
	listenAddr := &net.UDPAddr{IP: localAddr}
	switch p := p.(type) {
	case pWRQ:
		filename, mode, opts, err := unpackRQ(p)
		if err != nil {
			return fmt.Errorf("unpack WRQ while starting transfer: %v", err)
		}
		wt := &receiver{
			send:        make([]byte, datagramLength),
			receive:     make([]byte, datagramLength),
			retry:       &backoff{handler: s.backoff},
			timeout:     s.timeout,
			retries:     s.retries,
			addr:        remoteAddr,
			localIP:     localAddr,
			mode:        mode,
			opts:        opts,
			maxBlockLen: maxBlockLen,
			hook:        s.hook,
			filename:    filename,
			startTime:   time.Now(),
		}
		if s.singlePort {
			wt.conn = &chanConnection{
				server:  s,
				srcAddr: listenAddr,
				addr:    remoteAddr,
				channel: listener,
				timeout: s.timeout,
			}
			wt.singlePort = true
		} else {
			conn, err := net.ListenUDP("udp", listenAddr)
			if err != nil {
				return err
			}
			wt.conn = &connConnection{conn: conn}
		}
		s.wg.Add(1)
		go func() {
			if s.writeHandler != nil {
				err := s.writeHandler(filename, wt)
				if err != nil {
					wt.abort(err)
				} else {
					wt.terminate()
				}
			} else {
				wt.abort(fmt.Errorf("server does not support write requests"))
			}
			s.wg.Done()
		}()
	case pRRQ:
		filename, mode, opts, err := unpackRQ(p)
		if err != nil {
			return fmt.Errorf("unpack RRQ: %v", err)
		}
		//fmt.Printf("got RRQ (filename=%s, mode=%s, opts=%v)\n", filename, mode, opts)
		rf := &sender{
			send:        make([]byte, datagramLength),
			sendA:       senderAnticipate{enabled: false},
			receive:     make([]byte, datagramLength),
			tid:         remoteAddr.Port,
			retry:       &backoff{handler: s.backoff},
			timeout:     s.timeout,
			retries:     s.retries,
			addr:        remoteAddr,
			localIP:     localAddr,
			mode:        mode,
			opts:        opts,
			maxBlockLen: maxBlockLen,
			hook:        s.hook,
			filename:    filename,
			startTime:   time.Now(),
		}
		if s.singlePort {
			rf.conn = &chanConnection{
				server:  s,
				srcAddr: listenAddr,
				addr:    remoteAddr,
				channel: listener,
				timeout: s.timeout,
			}
		} else {
			conn, err := net.ListenUDP("udp", listenAddr)
			if err != nil {
				return err
			}
			rf.conn = &connConnection{conn: conn}
		}
		if s.sendAEnable { /* senderAnticipate if enabled in server */
			rf.sendA.enabled = true /* pass enable from server to sender */
			sendAInit(&rf.sendA, datagramLength, s.sendAWinSz)
		}
		s.wg.Add(1)
		go func(rh func(string, io.ReaderFrom) error, rf *sender, wg *sync.WaitGroup) {
			if s.readHandler != nil {
				err := s.readHandler(filename, rf)
				if err != nil {
					rf.abort(err)
				}
			} else {
				rf.abort(fmt.Errorf("server does not support read requests"))
			}
			s.wg.Done()
		}(s.readHandler, rf, s.wg)
	default:
		return fmt.Errorf("unexpected %T", p)
	}
	return nil
}
