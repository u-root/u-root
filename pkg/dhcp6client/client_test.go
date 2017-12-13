package dhcp6client

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"syscall"
	"testing"
	"time"

	"github.com/mdlayher/dhcp6"
)

type timeoutErr struct{}

func (timeoutErr) Error() string {
	return "i/o timeout"
}

func (timeoutErr) Timeout() bool {
	return true
}

type udpPacket struct {
	source  *net.UDPAddr
	dest    *net.UDPAddr
	payload []byte
}

// mockUDPConn implements net.PacketConn.
type mockUDPConn struct {
	// This'll just be nil for all the methods we don't implement.

	// in is the queue of packets ReadFromUDP reads from.
	//
	// ReadFromUDP returns io.EOF when in is closed.
	in chan udpPacket

	inTimer *time.Timer

	// out is the queue of packets WriteTo writes to.
	out chan<- udpPacket

	closed bool
}

func newMockUDPConn(in chan udpPacket, out chan<- udpPacket) *mockUDPConn {
	return &mockUDPConn{
		in:  in,
		out: out,
	}
}

// SetReadDeadline implements PacketConn.SetReadDeadline.
func (m *mockUDPConn) SetReadDeadline(t time.Time) error {
	duration := t.Sub(time.Now())
	if duration < 0 {
		return fmt.Errorf("deadline must be in the future")
	}
	m.inTimer = time.NewTimer(duration)
	return nil
}

func (m *mockUDPConn) LocalAddr() net.Addr {
	panic("unused")
}

func (m *mockUDPConn) SetWriteDeadline(t time.Time) error {
	panic("unused")
}

func (m *mockUDPConn) SetDeadline(t time.Time) error {
	panic("unused")
}

// Close implements PacketConn.Close.
func (m *mockUDPConn) Close() error {
	m.closed = true
	close(m.out)
	return nil
}

// ReadFrom is a mock for PacketConn.ReadFromUDP.
func (m *mockUDPConn) ReadFrom(b []byte) (int, net.Addr, error) {
	// Make sure we don't have data waiting.
	select {
	case p, ok := <-m.in:
		if !ok {
			// Connection was closed.
			return 0, nil, nil
		}
		return copy(b, p.payload), p.source, nil
	default:
	}

	select {
	case p, ok := <-m.in:
		if !ok {
			return 0, nil, nil
		}
		return copy(b, p.payload), p.source, nil
	case <-m.inTimer.C:
		// This net.OpError will return true for Timeout().
		return 0, nil, &net.OpError{Err: timeoutErr{}}
	}
}

// WriteTo is a mock for PacketConn.WriteTo.
func (m *mockUDPConn) WriteTo(b []byte, dest net.Addr) (int, error) {
	if m.closed {
		return 0, syscall.EBADF
	}

	m.out <- udpPacket{
		dest:    dest.(*net.UDPAddr),
		payload: b,
	}
	return len(b), nil
}

type server struct {
	in  chan udpPacket
	out chan udpPacket

	received []*dhcp6.Packet

	// Each received packet can have more than one response (in theory,
	// from different servers sending different Advertise, for example).
	responses [][]*dhcp6.Packet
}

func (s *server) serve(ctx context.Context) {
	go func() {
		select {
		case udpPkt, ok := <-s.in:
			if !ok {
				break
			}

			// What did we get?
			var pkt dhcp6.Packet
			if err := (&pkt).UnmarshalBinary(udpPkt.payload); err != nil {
				panic(fmt.Sprintf("invalid dhcp6 packet %q: %v", udpPkt.payload, err))
			}
			s.received = append(s.received, &pkt)

			if len(s.responses) > 0 {
				resps := s.responses[0]
				// What should we send in response?
				for _, resp := range resps {
					bin, err := resp.MarshalBinary()
					if err != nil {
						panic(fmt.Sprintf("failed to serialize dhcp6 packet %v: %v", resp, err))
					}
					s.out <- udpPacket{
						source:  udpPkt.dest,
						payload: bin,
					}
				}
				s.responses = s.responses[1:]
			}

		case <-ctx.Done():
			break
		}

		// We're done sending stuff.
		close(s.out)
	}()

}

func ComparePacket(got *dhcp6.Packet, want *dhcp6.Packet) error {
	aa, err := got.MarshalBinary()
	if err != nil {
		panic(err)
	}
	bb, err := want.MarshalBinary()
	if err != nil {
		panic(err)
	}
	if bytes.Compare(aa, bb) != 0 {
		return fmt.Errorf("packet got %v, want %v", got, want)
	}
	return nil
}

func pktsExpected(got []*dhcp6.Packet, want []*dhcp6.Packet) error {
	if len(got) != len(want) {
		return fmt.Errorf("got %d packets, want %d packets", len(got), len(want))
	}

	for i := range got {
		if err := ComparePacket(got[i], want[i]); err != nil {
			return err
		}
	}
	return nil
}

func serveAndClient(ctx context.Context, responses [][]*dhcp6.Packet) (*Client, *mockUDPConn) {
	// These are the client's channels.
	in := make(chan udpPacket, 100)
	out := make(chan udpPacket, 100)

	mockConn := &mockUDPConn{
		in:  in,
		out: out,
	}

	mc := &Client{
		srcMAC:  []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab},
		conn:    mockConn,
		retry:   1,
		timeout: time.Second,
	}

	// Of course, for the server they are reversed.
	s := &server{
		in:        out,
		out:       in,
		responses: responses,
	}
	go s.serve(ctx)

	return mc, mockConn
}

func TestSendAndRead(t *testing.T) {
	for _, tt := range []struct {
		desc   string
		send   *dhcp6.Packet
		server []*dhcp6.Packet

		// If want is nil, we assume server contains what is wanted.
		want    []*dhcp6.Packet
		wantErr error
	}{
		{
			desc: "two response packets",
			send: &dhcp6.Packet{
				MessageType:   dhcp6.MessageTypeSolicit,
				TransactionID: [3]byte{0x33, 0x33, 0x33},
			},
			server: []*dhcp6.Packet{
				{
					MessageType:   dhcp6.MessageTypeAdvertise,
					TransactionID: [3]byte{0x33, 0x33, 0x33},
				},
				{
					MessageType:   dhcp6.MessageTypeAdvertise,
					TransactionID: [3]byte{0x33, 0x33, 0x33},
				},
			},
		},
		{
			desc: "one response packet",
			send: &dhcp6.Packet{
				MessageType:   dhcp6.MessageTypeSolicit,
				TransactionID: [3]byte{0x33, 0x33, 0x33},
			},
			server: []*dhcp6.Packet{
				{
					MessageType:   dhcp6.MessageTypeAdvertise,
					TransactionID: [3]byte{0x33, 0x33, 0x33},
				},
			},
		},
		{
			desc: "one response packet, one invalid XID",
			send: &dhcp6.Packet{
				MessageType:   dhcp6.MessageTypeSolicit,
				TransactionID: [3]byte{0x33, 0x33, 0x33},
			},
			server: []*dhcp6.Packet{
				{
					MessageType:   dhcp6.MessageTypeAdvertise,
					TransactionID: [3]byte{0x33, 0x33, 0x33},
				},
				{
					MessageType:   dhcp6.MessageTypeAdvertise,
					TransactionID: [3]byte{0x77, 0x77, 0x77},
				},
			},
			want: []*dhcp6.Packet{
				{
					MessageType:   dhcp6.MessageTypeAdvertise,
					TransactionID: [3]byte{0x33, 0x33, 0x33},
				},
			},
		},
		{
			desc: "discard wrong XID",
			send: &dhcp6.Packet{
				MessageType:   dhcp6.MessageTypeSolicit,
				TransactionID: [3]byte{0x33, 0x33, 0x33},
			},
			server: []*dhcp6.Packet{
				{
					MessageType:   dhcp6.MessageTypeAdvertise,
					TransactionID: [3]byte{0x00, 0x00, 0x00},
				},
			},
			want:    []*dhcp6.Packet{}, // Explicitly empty.
			wantErr: context.DeadlineExceeded,
		},
		{
			desc: "no response, timeout",
			send: &dhcp6.Packet{
				MessageType:   dhcp6.MessageTypeSolicit,
				TransactionID: [3]byte{0x33, 0x33, 0x33},
			},
			wantErr: context.DeadlineExceeded,
		},
	} {
		// Both server and client only get 2 seconds.
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		mc, _ := serveAndClient(ctx, [][]*dhcp6.Packet{tt.server})
		defer mc.conn.Close()

		out, errCh := mc.SendAndRead(ctx, DefaultServers, tt.send)

		var rcvd []*dhcp6.Packet
		for packet := range out {
			rcvd = append(rcvd, packet.Packet)
		}

		if err, ok := <-errCh; ok && err.Err != tt.wantErr {
			t.Errorf("SendAndRead(%v): got %v, want %v", tt.send, err.Err, tt.wantErr)
		} else if !ok && tt.wantErr != nil {
			t.Errorf("got no error, want %v", tt.wantErr)
		}

		want := tt.want
		if want == nil {
			want = tt.server
		}
		if err := pktsExpected(rcvd, want); err != nil {
			t.Errorf("got unexpected packets: %v", err)
		}
	}
}

func TestSendAndReadHandleCancel(t *testing.T) {
	pkt := &dhcp6.Packet{
		MessageType:   dhcp6.MessageTypeSolicit,
		TransactionID: [3]byte{0x33, 0x33, 0x33},
	}

	responses := []*dhcp6.Packet{
		{
			MessageType:   dhcp6.MessageTypeAdvertise,
			TransactionID: [3]byte{0x33, 0x33, 0x33},
		},
		{
			MessageType:   dhcp6.MessageTypeRelayRepl,
			TransactionID: [3]byte{0x33, 0x33, 0x33},
		},
		{
			MessageType:   dhcp6.MessageTypeInformationRequest,
			TransactionID: [3]byte{0x33, 0x33, 0x33},
		},
		{
			MessageType:   dhcp6.MessageTypeReply,
			TransactionID: [3]byte{0x33, 0x33, 0x33},
		},
	}

	// Both the server and client only get 2 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	mc, udpConn := serveAndClient(ctx, [][]*dhcp6.Packet{responses})
	defer mc.conn.Close()

	out, errCh := mc.SendAndRead(ctx, DefaultServers, pkt)

	var counter int
	for range out {
		counter++
		if counter == 2 {
			cancel()
		}
	}

	if err, ok := <-errCh; ok {
		t.Errorf("got %v, want nil error", err)
	}

	// Make sure that two packets are still in the queue to be read.
	for packet := range udpConn.in {
		bin, err := responses[counter].MarshalBinary()
		if err != nil {
			panic(err)
		}
		if bytes.Compare(packet.payload, bin) != 0 {
			t.Errorf("SendAndRead read more packets than expected!")
		}
		counter++
	}
}

func TestSendAndReadDiscardGarbage(t *testing.T) {
	pkt := &dhcp6.Packet{
		MessageType:   dhcp6.MessageTypeSolicit,
		TransactionID: [3]byte{0x33, 0x33, 0x33},
	}

	responses := []*dhcp6.Packet{
		{
			MessageType:   dhcp6.MessageTypeAdvertise,
			TransactionID: [3]byte{0x33, 0x33, 0x33},
		},
	}

	// Both the server and client only get 2 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	mc, udpConn := serveAndClient(ctx, [][]*dhcp6.Packet{responses})
	defer mc.conn.Close()

	udpConn.in <- udpPacket{
		payload: []byte{0x01}, // Too short for valid DHCPv6 packet.
	}

	out, errCh := mc.SendAndRead(ctx, DefaultServers, pkt)

	var i int
	for recvd := range out {
		if err := ComparePacket(recvd.Packet, responses[i]); err != nil {
			t.Error(err)
		}
		i++
	}

	if err, ok := <-errCh; ok {
		t.Errorf("SendAndRead(%v): got %v %v, want %v", pkt, ok, err, nil)
	}
	if i != len(responses) {
		t.Errorf("should have received %d valid packet, counter is %d", len(responses), i)
	}
}

func TestSendAndReadDiscardGarbageTimeout(t *testing.T) {
	pkt := &dhcp6.Packet{
		MessageType:   dhcp6.MessageTypeSolicit,
		TransactionID: [3]byte{0x33, 0x33, 0x33},
	}

	// Both the server and client only get 2 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	mc, udpConn := serveAndClient(ctx, nil)
	defer mc.conn.Close()

	udpConn.in <- udpPacket{
		payload: []byte{0x01}, // Too short for valid DHCPv6 packet.
	}

	out, errCh := mc.SendAndRead(ctx, DefaultServers, pkt)

	var counter int
	for range out {
		counter++
	}

	if err, ok := <-errCh; !ok || err == nil || err.Err != context.DeadlineExceeded {
		t.Errorf("SendAndRead(%v): got %v %v, want %v", pkt, ok, err, context.DeadlineExceeded)
	}
	if counter != 0 {
		t.Errorf("should not have received a valid packet, counter is %d", counter)
	}
}
