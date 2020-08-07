// +build linux

package rtnetlink

import (
	"encoding"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/mdlayher/netlink"
	"golang.org/x/sys/unix"
)

func TestConnExecute(t *testing.T) {
	skipBigEndian(t)

	req := &LinkMessage{}

	wantnl := netlink.Message{
		Header: netlink.Header{
			Type:  unix.RTM_GETLINK,
			Flags: netlink.Request,
			// Sequence and PID not set because we are mocking the underlying
			// netlink connection.
		},
		Data: mustMarshal(req),
	}
	wantrt := []LinkMessage{
		{
			Family: 0x101,
			Type:   0x0,
			Index:  0x4030201,
			Flags:  0x101,
			Change: 0x4030201,
		},
	}

	c, tc := testConn(t)
	tc.receive = []netlink.Message{{
		Header: netlink.Header{
			Length: 16,
			Type:   unix.RTM_GETLINK,
			// Sequence and PID not set because we are mocking the underlying
			// netlink connection.
		},
		Data: []byte{
			0x01, 0x01, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04,
			0x01, 0x01, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04,
		},
	}}

	msgs, err := c.Execute(req, unix.RTM_GETLINK, netlink.Request)
	if err != nil {
		t.Fatalf("failed to execute: %v", err)
	}

	links := make([]LinkMessage, 0, len(msgs))
	for _, m := range msgs {
		link := (m).(*LinkMessage)
		links = append(links, *link)
	}

	if want, got := wantnl, tc.send; !reflect.DeepEqual(want, got) {
		t.Fatalf("unexpected request:\n- want: %#v\n-  got: %#v",
			want, got)
	}
	if want, got := wantrt, links; !reflect.DeepEqual(want, got) {
		t.Fatalf("unexpected replies:\n- want: %#v\n-  got: %#v",
			want, got)
	}
}

func TestConnSend(t *testing.T) {
	skipBigEndian(t)

	req := &LinkMessage{}

	c, tc := testConn(t)

	nlreq, err := c.Send(req, unix.RTM_GETLINK, netlink.Request)
	if err != nil {
		t.Fatalf("failed to send: %v", err)
	}

	reqb, err := req.MarshalBinary()
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	want := netlink.Message{
		Header: netlink.Header{
			Type:  unix.RTM_GETLINK,
			Flags: netlink.Request,
		},
		Data: reqb,
	}

	if got := tc.send; !reflect.DeepEqual(want, got) {
		t.Fatalf("unexpected output message from Conn.Send:\n- want: %#v\n-  got: %#v",
			want, got)
	}
	if got := nlreq; !reflect.DeepEqual(want, got) {
		t.Fatalf("unexpected modified message:\n- want: %#v\n-  got: %#v",
			want, got)
	}
}

func TestConnReceive(t *testing.T) {
	skipBigEndian(t)

	c, tc := testConn(t)
	tc.receive = []netlink.Message{
		{
			Header: netlink.Header{
				Length:   16,
				Sequence: 1,
				Type:     unix.RTM_GETLINK,
				PID:      uint32(os.Getpid()),
			},
			Data: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			Header: netlink.Header{
				Length:   16,
				Sequence: 1,
				Type:     unix.RTM_GETLINK,
				PID:      uint32(os.Getpid()),
			},
			Data: []byte{
				0x02, 0x01, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04,
				0x02, 0x01, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04,
			},
		},
	}

	wantnl := tc.receive
	wantrt := []LinkMessage{
		{
			Family: 0x0,
			Type:   0x0,
			Index:  0x0,
			Flags:  0x0,
			Change: 0x0,
		},
		{
			Family: 0x102,
			Type:   0x0,
			Index:  0x4030201,
			Flags:  0x102,
			Change: 0x4030201,
		},
	}

	rtmsgs, nlmsgs, err := c.Receive()
	if err != nil {
		t.Fatalf("failed to receive messages: %v", err)
	}

	links := make([]LinkMessage, 0, len(rtmsgs))
	for _, m := range rtmsgs {
		link := (m).(*LinkMessage)
		links = append(links, *link)
	}

	if want, got := wantnl, nlmsgs; !reflect.DeepEqual(want, got) {
		t.Fatalf("unexpected netlink.Messages from Conn.Receive:\n- want: %#v\n-  got: %#v",
			want, got)
	}

	if want, got := wantrt, links; !reflect.DeepEqual(want, got) {
		t.Fatalf("unexpected Messages from Conn.Receive:\n- want: %#v\n-  got: %#v",
			want, got)
	}
}

func testConn(t *testing.T) (*Conn, *testNetlinkConn) {
	c := &testNetlinkConn{}
	return newConn(c), c
}

type testNetlinkConn struct {
	send    netlink.Message
	receive []netlink.Message

	noopConn
}

func (c *testNetlinkConn) Send(m netlink.Message) (netlink.Message, error) {
	c.send = m
	return m, nil
}

func (c *testNetlinkConn) Receive() ([]netlink.Message, error) {
	return c.receive, nil
}

func (c *testNetlinkConn) Execute(m netlink.Message) ([]netlink.Message, error) {
	c.send = m
	return c.receive, nil
}

type noopConn struct{}

func (c *noopConn) Close() error                                         { return nil }
func (c *noopConn) Send(_ netlink.Message) (netlink.Message, error)      { return netlink.Message{}, nil }
func (c *noopConn) Receive() ([]netlink.Message, error)                  { return nil, nil }
func (c *noopConn) Execute(m netlink.Message) ([]netlink.Message, error) { return nil, nil }
func (c *noopConn) SetReadDeadline(t time.Time) error                    { return nil }

func mustMarshal(m encoding.BinaryMarshaler) []byte {
	b, err := m.MarshalBinary()
	if err != nil {
		panic(fmt.Sprintf("failed to marshal binary: %v", err))
	}

	return b
}
