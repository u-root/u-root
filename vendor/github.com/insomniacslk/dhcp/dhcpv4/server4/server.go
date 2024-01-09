// Package server4 is a basic, extensible DHCPv4 server.
//
// To use the DHCPv4 server code you have to call NewServer with two arguments:
//  - an interface to listen on,
//  - an address to listen on, and
//  - a handler function, that will be called every time a valid DHCPv4 packet is
//    received.
//
// The address to listen on is used to know IP address, port and optionally the
// scope to create and UDP socket to listen on for DHCPv4 traffic.
//
// The handler is a function that takes as input a packet connection, that can
// be used to reply to the client; a peer address, that identifies the client
// sending the request, and the DHCPv4 packet itself. Just implement your
// custom logic in the handler.
//
// Optionally, NewServer can receive options that will modify the server
// object. Some options already exist, for example WithConn. If this option is
// passed with a valid connection, the listening address argument is ignored.
//
// Example program:
//
//	package main
//
//	import (
//		"log"
//		"net"
//
//		"github.com/insomniacslk/dhcp/dhcpv4"
//		"github.com/insomniacslk/dhcp/dhcpv4/server4"
//	)
//
//	func handler(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
//		// this function will just print the received DHCPv4 message, without replying
//		log.Print(m.Summary())
//	}
//
//	func main() {
//		laddr := net.UDPAddr{
//			IP:   net.ParseIP("127.0.0.1"),
//			Port: 67,
//		}
//		server, err := server4.NewServer("eth0", &laddr, handler)
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		// This never returns. If you want to do other stuff, dump it into a
//		// goroutine.
//		server.Serve()
//	}
//
package server4

import (
	"log"
	"net"
	"os"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

// Handler is a type that defines the handler function to be called every time a
// valid DHCPv4 message is received
type Handler func(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4)

// Server represents a DHCPv4 server object
type Server struct {
	conn    net.PacketConn
	Handler Handler
	logger  Logger
}

// Serve serves requests.
func (s *Server) Serve() error {
	s.logger.Printf("Server listening on %s", s.conn.LocalAddr())
	s.logger.Printf("Ready to handle requests")

	defer s.Close()
	for {
		rbuf := make([]byte, 4096) // FIXME this is bad
		n, peer, err := s.conn.ReadFrom(rbuf)
		if err != nil {
			s.logger.Printf("Error reading from packet conn: %v", err)
			return err
		}
		s.logger.Printf("Handling request from %v", peer)

		m, err := dhcpv4.FromBytes(rbuf[:n])
		if err != nil {
			s.logger.Printf("Error parsing DHCPv4 request: %v", err)
			continue
		}

		upeer, ok := peer.(*net.UDPAddr)
		if !ok {
			s.logger.Printf("Not a UDP connection? Peer is %s", peer)
			continue
		}
		// Set peer to broadcast if the client did not have an IP.
		if upeer.IP == nil || upeer.IP.To4().Equal(net.IPv4zero) {
			upeer = &net.UDPAddr{
				IP:   net.IPv4bcast,
				Port: upeer.Port,
			}
		}
		go s.Handler(s.conn, upeer, m)
	}
}

// Close sends a termination request to the server, and closes the UDP listener.
func (s *Server) Close() error {
	return s.conn.Close()
}

// ServerOpt adds optional configuration to a server.
type ServerOpt func(s *Server)

// WithConn configures the server with the given connection.
func WithConn(c net.PacketConn) ServerOpt {
	return func(s *Server) {
		s.conn = c
	}
}

// NewServer initializes and returns a new Server object
func NewServer(ifname string, addr *net.UDPAddr, handler Handler, opt ...ServerOpt) (*Server, error) {
	s := &Server{
		Handler: handler,
		logger:  EmptyLogger{},
	}

	for _, o := range opt {
		o(s)
	}
	if s.conn == nil {
		var err error
		conn, err := NewIPv4UDPConn(ifname, addr)
		if err != nil {
			return nil, err
		}
		s.conn = conn
	}
	return s, nil
}

// WithSummaryLogger logs one-line DHCPv4 message summaries when sent & received.
func WithSummaryLogger() ServerOpt {
	return func(s *Server) {
		s.logger = ShortSummaryLogger{
			Printfer: log.New(os.Stderr, "[dhcpv4] ", log.LstdFlags),
		}
	}
}

// WithDebugLogger logs multi-line full DHCPv4 messages when sent & received.
func WithDebugLogger() ServerOpt {
	return func(s *Server) {
		s.logger = DebugLogger{
			Printfer: log.New(os.Stderr, "[dhcpv4] ", log.LstdFlags),
		}
	}
}

// WithLogger set the logger (see interface Logger).
func WithLogger(newLogger Logger) ServerOpt {
	return func(s *Server) {
		s.logger = newLogger
	}
}
