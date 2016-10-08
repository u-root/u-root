// Copyright (C) 2016 Kale Blankenship. All rights reserved.
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details

package trivialt

import "net"

import "sync"

// Server contains the configuration to run a TFTP server.
//
// A ReadHandler, WriteHandler, or both can be registered to the server. If one
// of the handlers isn't registered, the server will return errors to clients
// attempting to use them.
type Server struct {
	log     *logger
	net     string
	addrStr string
	addr    *net.UDPAddr
	connMu  sync.RWMutex
	conn    *net.UDPConn
	close   bool

	singlePort bool
	mgr        *connManager

	retransmit int // Per-packet retransmission limit

	rh ReadHandler
	wh WriteHandler
}

type connManager struct {
	sync.RWMutex
	reqMap map[string]chan []byte
}

func (m *connManager) New(addr net.Addr) chan []byte {
	reqChan := make(chan []byte, 64) // TODO (better value)
	m.Lock()
	m.reqMap[addr.String()] = reqChan
	m.Unlock()
	return reqChan
}

func (m *connManager) Get(addr net.Addr) (chan []byte, bool) {
	m.RLock()
	reqChan, ok := m.reqMap[addr.String()]
	m.RUnlock()
	return reqChan, ok
}

func (m *connManager) Remove(addr net.Addr) {
	m.Lock()
	delete(m.reqMap, addr.String())
	m.Unlock()
}

// NewServer returns a configured Server.
//
// Addr is the network address to listen on and is in the form "host:port".
// If a no host is given the server will listen on all interfaces.
//
// Any number of ServerOpts can be provided to configure optional values.
func NewServer(addr string, opts ...ServerOpt) (*Server, error) {
	s := &Server{
		log:        newLogger("server"),
		net:        defaultUDPNet,
		addrStr:    addr,
		retransmit: defaultRetransmit,
	}

	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

// Addr is the network address of the server. It is available
// after the server has been started.
func (s *Server) Addr() (*net.UDPAddr, error) {
	s.connMu.RLock()
	defer s.connMu.RUnlock()
	if s.conn == nil {
		return nil, ErrAddressNotAvailable
	}
	return s.conn.LocalAddr().(*net.UDPAddr), nil
}

// ReadHandler registers a ReadHandler for the server.
func (s *Server) ReadHandler(rh ReadHandler) {
	s.rh = rh
}

// WriteHandler registers a WriteHandler for the server.
func (s *Server) WriteHandler(wh WriteHandler) {
	s.wh = wh
}

// Serve starts the server using an existing UDPConn.
func (s *Server) Serve(conn *net.UDPConn) error {
	if s.rh == nil && s.wh == nil {
		return ErrNoRegisteredHandlers
	}

	s.connMu.Lock()
	s.conn = conn
	s.connMu.Unlock()

	s.connMu.RLock()
	defer s.connMu.RUnlock()
	buf := make([]byte, 65536) // Largest possible TFTP datagram
	for {
		numBytes, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			if s.close {
				return nil
			}
			return wrapError(err, "reading from conn")
		}

		// Make a copy of the received data
		b := make([]byte, numBytes)
		copy(b, buf)

		switch buf[1] {
		case 1: //RRQ
			go s.dispatchReadRequest(addr, b)
		case 2: //WRQ
			go s.dispatchWriteRequest(addr, b)
		default:
			go s.demuxToConn(addr, b)
		}
	}
}

// Connected is true if the server has started serving.
func (s *Server) Connected() bool {
	s.connMu.RLock()
	defer s.connMu.RUnlock()
	return s.conn != nil
}

// Close stops the server and closes the network connection.
func (s *Server) Close() error {
	s.connMu.RLock()
	defer s.connMu.RUnlock()
	s.close = true
	return s.conn.Close()
}

// dispatchReadRequest dispatches the read handler, if it is registered.
// If a handler is not registered the server sends an error to the client.
func (s *Server) dispatchReadRequest(addr *net.UDPAddr, buf []byte) {
	// Check for handler
	if s.rh == nil {
		s.log.debug("No read handler registered.")
		var err datagram
		err.writeError(ErrCodeIllegalOperation, "Server does not support read requests.")
		_, _ = s.conn.WriteTo(err.bytes(), addr) // Ignore error
		return
	}

	c, closer, err := s.newConn(addr, buf)
	if err != nil {
		return
	}
	defer errorDefer(closer, s.log, "error closing network connection in dispath")

	s.log.debug("New request from %v: %s", addr, c.rx)

	// Create request
	w := &readRequest{conn: c, name: c.rx.filename()}

	// execute handler
	s.rh.ServeTFTP(w)
}

// dispatchWriteRequest dispatches the read handler, if it is registered.
// If a handler is not registered the server sends an error to the client.
func (s *Server) dispatchWriteRequest(addr *net.UDPAddr, buf []byte) {
	// Check for handler
	if s.wh == nil {
		s.log.debug("No write handler registered.")
		var err datagram
		err.writeError(ErrCodeIllegalOperation, "Server does not support write requests.")
		_, _ = s.conn.WriteTo(err.bytes(), addr) // Ignore error
		return
	}

	c, closer, err := s.newConn(addr, buf)
	if err != nil {
		return
	}
	defer errorDefer(closer, s.log, "error closing network connection in dispath")

	s.log.debug("New request from %v: %s", addr, c.rx)

	// Create request
	w := &writeRequest{conn: c, name: c.rx.filename()}

	// parse options to get size
	c.log.trace("performing write setup")
	if err := c.readSetup(); err != nil {
		c.err = err
	}

	s.wh.ReceiveTFTP(w)
}

func (s *Server) demuxToConn(addr *net.UDPAddr, buf []byte) {
	if s.singlePort {
		if reqChan, ok := s.mgr.Get(addr); ok {
			reqChan <- buf
			return
		}
	}

	// RFC1350:
	// "If a source TID does not match, the packet should be
	// discarded as erroneously sent from somewhere else.  An error packet
	// should be sent to the source of the incorrect packet, while not
	// disturbing the transfer."
	dg := datagram{}
	dg.writeError(ErrCodeUnknownTransferID, "Unexpected TID")
	// Don't care about an error here, just a courtesy
	_, _ = s.conn.WriteTo(dg.bytes(), addr)
	s.log.debug("Unexpected datagram: %s", dg)
}

func (s *Server) newConn(addr *net.UDPAddr, buf []byte) (*conn, func() error, error) {
	var c *conn
	var err error
	var dg datagram

	dg.setBytes(buf)

	// Validate request datagram
	if err := dg.validate(); err != nil {
		s.log.debug("Error decoding new request: %v", err)
		return nil, nil, err
	}

	if s.singlePort {
		c = newSinglePortConn(addr, dg.mode(), s.conn, s.mgr.New(addr))
	} else {
		c, err = newConn(s.net, dg.mode(), addr) // Use empty mode until request has been parsed.
		if err != nil {
			s.log.err("Received error opening connection for new request: %v", err)
			return nil, nil, err
		}
	}

	c.rx = dg
	// Set retransmit
	c.retransmit = s.retransmit

	closer := func() error {
		err := c.Close()
		if s.singlePort {
			s.mgr.Remove(addr)
		}
		return err
	}

	return c, closer, nil
}

// ListenAndServe starts a configured server.
func (s *Server) ListenAndServe() error {
	addr, err := net.ResolveUDPAddr(s.net, s.addrStr)
	if err != nil {
		return wrapError(err, "resolving server address")
	}
	s.addr = addr

	conn, err := net.ListenUDP(s.net, s.addr)
	if err != nil {
		return wrapError(err, "opening network connection")
	}

	return wrapError(s.Serve(conn), "serving tftp")
}

// ServerOpt is a function that configures a Server.
type ServerOpt func(*Server) error

// ServerNet configures the network a server listens on.
// Must be one of: udp, udp4, udp6.
//
// Default: udp.
func ServerNet(net string) ServerOpt {
	return func(s *Server) error {
		if net != "udp" && net != "udp4" && net != "udp6" {
			return ErrInvalidNetwork
		}
		s.net = net
		return nil
	}
}

// ServerRetransmit configures the per-packet retransmission limit for all requests.
//
// Default: 10.
func ServerRetransmit(i int) ServerOpt {
	return func(s *Server) error {
		if i < 0 {
			return ErrInvalidRetransmit
		}
		s.retransmit = i
		return nil
	}
}

// ServerSinglePort enables the server to service all requests via a single port rather
// than the standard TFTP behavior of each client communicating on a separate port.
//
// This is an experimental feature.
//
// Default is disabled.
func ServerSinglePort(enable bool) ServerOpt {
	return func(s *Server) error {
		if enable {
			s.singlePort = true
			s.mgr = &connManager{reqMap: make(map[string]chan []byte)}
		}
		return nil
	}
}
