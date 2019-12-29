// +build go1.12

package server4

import (
	"context"
	"log"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/nclient4"
	"github.com/insomniacslk/dhcp/interfaces"
	"github.com/stretchr/testify/require"
)

func init() {
	// initialize seed. This is generally bad, but "good enough"
	// to generate random ports for these tests
	rand.Seed(time.Now().UTC().UnixNano())
}

// DORAHandler is a server handler suitable for DORA transactions
func DORAHandler(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
	if m == nil {
		log.Printf("Packet is nil!")
		return
	}
	if m.OpCode != dhcpv4.OpcodeBootRequest {
		log.Printf("Not a BootRequest!")
		return
	}
	reply, err := dhcpv4.NewReplyFromRequest(m)
	if err != nil {
		log.Printf("NewReplyFromRequest failed: %v", err)
		return
	}
	reply.UpdateOption(dhcpv4.OptServerIdentifier(net.IP{1, 2, 3, 4}))
	switch mt := m.MessageType(); mt {
	case dhcpv4.MessageTypeDiscover:
		reply.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeOffer))
	case dhcpv4.MessageTypeRequest:
		reply.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeAck))
	default:
		log.Printf("Unhandled message type: %v", mt)
		return
	}

	if _, err := conn.WriteTo(reply.ToBytes(), peer); err != nil {
		log.Printf("Cannot reply to client: %v", err)
	}
}

// utility function to set up a client and a server instance and run it in
// background. The caller needs to call Server.Close() once finished.
func setUpClientAndServer(t *testing.T, iface net.Interface, handler Handler) (*nclient4.Client, *Server) {
	// strong assumption, I know
	loAddr := net.ParseIP("127.0.0.1")
	saddr := &net.UDPAddr{
		IP:   loAddr,
		Port: 0,
	}
	caddr := net.UDPAddr{
		IP:   loAddr,
		Port: 0,
	}
	s, err := NewServer("", saddr, handler)
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		_ = s.Serve()
	}()
	saddr = s.conn.LocalAddr().(*net.UDPAddr)

	clientConn, err := NewIPv4UDPConn("", &caddr)
	if err != nil {
		t.Fatal(err)
	}
	c, err := nclient4.NewWithConn(clientConn, iface.HardwareAddr, nclient4.WithServerAddr(saddr))
	if err != nil {
		t.Fatal(err)
	}
	return c, s
}

func TestServer(t *testing.T) {
	ifaces, err := interfaces.GetLoopbackInterfaces()
	require.NoError(t, err)
	require.NotEqual(t, 0, len(ifaces))

	// lo has a HardwareAddr of "nil". The client will drop all packets
	// that don't match the HWAddr of the client interface.
	hwaddr := net.HardwareAddr{1, 2, 3, 4, 5, 6}
	ifaces[0].HardwareAddr = hwaddr

	c, s := setUpClientAndServer(t, ifaces[0], DORAHandler)
	defer func() {
		require.Nil(t, s.Close())
	}()

	xid := dhcpv4.TransactionID{0xaa, 0xbb, 0xcc, 0xdd}

	modifiers := []dhcpv4.Modifier{
		dhcpv4.WithTransactionID(xid),
		dhcpv4.WithHwAddr(ifaces[0].HardwareAddr),
	}

	offer, ack, err := c.Request(context.Background(), modifiers...)
	require.NoError(t, err)
	require.NotNil(t, offer, ack)
	for _, p := range []*dhcpv4.DHCPv4{offer, ack} {
		require.Equal(t, xid, p.TransactionID)
		require.Equal(t, ifaces[0].HardwareAddr, p.ClientHWAddr)
	}
}

func TestBadAddrFamily(t *testing.T) {
	saddr := &net.UDPAddr{
		IP:   net.IPv6loopback,
		Port: 0,
	}
	_, err := NewServer("", saddr, DORAHandler)
	if err == nil {
		t.Fatal("Expected server4.NewServer to fail with an IPv6 address")
	}
}
