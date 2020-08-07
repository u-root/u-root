package server6

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/insomniacslk/dhcp/dhcpv6"
	"github.com/insomniacslk/dhcp/dhcpv6/nclient6"
	"github.com/insomniacslk/dhcp/interfaces"
	"github.com/stretchr/testify/require"
)

// Turns a connected UDP conn into an "unconnected" UDP conn.
type unconnectedConn struct {
	*net.UDPConn
}

func (f unconnectedConn) WriteTo(b []byte, _ net.Addr) (int, error) {
	return f.UDPConn.Write(b)
}

func (f unconnectedConn) ReadFrom(b []byte) (int, net.Addr, error) {
	n, err := f.Read(b)
	return n, nil, err
}

// utility function to set up a client and a server instance and run it in
// background. The caller needs to call Server.Close() once finished.
func setUpClientAndServer(handler Handler) (*nclient6.Client, *Server) {
	laddr := &net.UDPAddr{
		IP:   net.ParseIP("::1"),
		Port: 0,
	}
	s, err := NewServer("", laddr, handler)
	if err != nil {
		panic(err)
	}
	go func() {
		_ = s.Serve()
	}()

	clientConn, err := net.DialUDP("udp6", &net.UDPAddr{IP: net.ParseIP("::1")}, s.conn.LocalAddr().(*net.UDPAddr))
	if err != nil {
		panic(err)
	}

	c, err := nclient6.NewWithConn(unconnectedConn{clientConn}, net.HardwareAddr{1, 2, 3, 4, 5, 6})
	if err != nil {
		panic(err)
	}
	return c, s
}

func TestServer(t *testing.T) {
	handler := func(conn net.PacketConn, peer net.Addr, m dhcpv6.DHCPv6) {
		msg := m.(*dhcpv6.Message)
		adv, err := dhcpv6.NewAdvertiseFromSolicit(msg)
		if err != nil {
			log.Printf("NewAdvertiseFromSolicit failed: %v", err)
			return
		}
		if _, err := conn.WriteTo(adv.ToBytes(), peer); err != nil {
			log.Printf("Cannot reply to client: %v", err)
		}
	}

	c, s := setUpClientAndServer(handler)
	defer s.Close()

	ifaces, err := interfaces.GetLoopbackInterfaces()
	require.NoError(t, err)
	require.NotEqual(t, 0, len(ifaces))

	_, err = c.Solicit(context.Background(), dhcpv6.WithRapidCommit)
	require.NoError(t, err)
}
