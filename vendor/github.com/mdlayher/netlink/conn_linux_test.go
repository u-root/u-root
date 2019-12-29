//+build linux

package netlink

import (
	"errors"
	"math"
	"os"
	"reflect"
	"testing"
	"time"

	"golang.org/x/sys/unix"
)

func TestLinuxConn_bindOK(t *testing.T) {
	s := &testSocket{}
	if _, _, err := bind(s, &Config{}); err != nil {
		t.Fatalf("failed to bind: %v", err)
	}

	addr := &unix.SockaddrNetlink{
		Family: unix.AF_NETLINK,
	}

	if want, got := addr, s.bind; !reflect.DeepEqual(want, got) {
		t.Fatalf("unexpected bind address:\n- want: %#v\n-  got: %#v",
			want, got)
	}
}

func TestLinuxConn_bindBindErrorCloseSocket(t *testing.T) {
	// Trigger an error during bind with Bind, meaning that the socket should be
	// closed to avoid leaking file descriptors.
	s := &testSocket{
		bindErr: errors.New("cannot bind"),
	}

	if _, _, err := bind(s, &Config{}); err == nil {
		t.Fatal("no error occurred, but expected one")
	}

	if want, got := true, s.closed; want != got {
		t.Fatalf("unexpected socket closed:\n- want: %v\n-  got: %v",
			want, got)
	}
}

func TestLinuxConn_bindGetsocknameErrorCloseSocket(t *testing.T) {
	// Trigger an error during bind with Getsockname, meaning that the socket
	// should be closed to avoid leaking file descriptors.
	s := &testSocket{
		getsocknameErr: errors.New("cannot get socket name"),
	}

	if _, _, err := bind(s, &Config{}); err == nil {
		t.Fatal("no error occurred, but expected one")
	}

	if want, got := true, s.closed; want != got {
		t.Fatalf("unexpected socket closed:\n- want: %v\n-  got: %v",
			want, got)
	}
}

func TestLinuxConnSend(t *testing.T) {
	c, s := testLinuxConn(t, nil)

	req := Message{
		Header: Header{
			Length:   uint32(nlmsgAlign(nlmsgLength(2))),
			Flags:    Request | Acknowledge,
			Sequence: 1,
			PID:      uint32(os.Getpid()),
		},
		Data: []byte{0x01, 0x02},
	}

	if err := c.Send(req); err != nil {
		t.Fatalf("error while sending: %v", err)
	}

	// Pad data to 4 bytes as is done when marshaling for later comparison
	req.Data = append(req.Data, 0x00, 0x00)

	to := &unix.SockaddrNetlink{
		Family: unix.AF_NETLINK,
	}

	if want, got := 0, s.sendmsg.flags; want != got {
		t.Fatalf("unexpected sendmsg flags:\n- want: %v\n-  got: %v",
			want, got)
	}
	if want, got := to, s.sendmsg.to; !reflect.DeepEqual(want, got) {
		t.Fatalf("unexpected sendmsg address:\n- want: %#v\n-  got: %#v",
			want, got)
	}

	var out Message
	if err := (&out).UnmarshalBinary(s.sendmsg.p); err != nil {
		t.Fatalf("failed to unmarshal sendmsg buffer into message: %v", err)
	}

	if want, got := req, out; !reflect.DeepEqual(want, got) {
		t.Fatalf("unexpected output message:\n- want: %#v\n-  got: %#v",
			want, got)
	}
}

func TestLinuxConnReceive(t *testing.T) {
	// The request we sent netlink in the previous test; it will be echoed
	// back to us as part of this test
	req := Message{
		Header: Header{
			Length:   uint32(nlmsgAlign(nlmsgLength(4))),
			Flags:    Request | Acknowledge,
			Sequence: 1,
			PID:      uint32(os.Getpid()),
		},
		Data: []byte{0x01, 0x02, 0x00, 0x00},
	}
	reqb, err := req.MarshalBinary()
	if err != nil {
		t.Fatalf("failed to marshal request to binary: %v", err)
	}

	res := Message{
		Header: Header{
			// 16 bytes: header
			//  4 bytes: error code
			// 20 bytes: request message
			Length:   uint32(nlmsgAlign(nlmsgLength(24))),
			Type:     Error,
			Sequence: 1,
			PID:      uint32(os.Getpid()),
		},
		// Error code "success", and copy of request
		Data: append([]byte{0x00, 0x00, 0x00, 0x00}, reqb...),
	}
	resb, err := res.MarshalBinary()
	if err != nil {
		t.Fatalf("failed to marshal response to binary: %v", err)
	}

	c, s := testLinuxConn(t, nil)

	from := &unix.SockaddrNetlink{
		Family: unix.AF_NETLINK,
	}

	s.recvmsg.p = resb
	s.recvmsg.from = from

	msgs, err := c.Receive()
	if err != nil {
		t.Fatalf("failed to receive messages: %v", err)
	}

	if want, got := from, s.recvmsg.from; !reflect.DeepEqual(want, got) {
		t.Fatalf("unexpected recvmsg address:\n- want: %#v\n-  got: %#v",
			want, got)
	}

	// Expect a MSG_PEEK and then no flags on second call
	if want, got := 2, len(s.recvmsg.flags); want != got {
		t.Fatalf("unexpected number of calls to recvmsg:\n- want: %v\n-  got: %v",
			want, got)
	}
	if want, got := unix.MSG_PEEK, s.recvmsg.flags[0]; want != got {
		t.Fatalf("unexpected first recvmsg flags:\n- want: %v\n-  got: %v",
			want, got)
	}
	if want, got := 0, s.recvmsg.flags[1]; want != got {
		t.Fatalf("unexpected second recvmsg flags:\n- want: %v\n-  got: %v",
			want, got)
	}

	if want, got := 1, len(msgs); want != got {
		t.Fatalf("unexpected number of messages:\n- want: %v\n-  got: %v",
			want, got)
	}

	if want, got := res, msgs[0]; !reflect.DeepEqual(want, got) {
		t.Fatalf("unexpected output message:\n- want: %#v\n-  got: %#v",
			want, got)
	}
}

func TestLinuxConnReceiveLargeMessage(t *testing.T) {
	n := os.Getpagesize() * 4

	res := Message{
		Header: Header{
			Length:   uint32(nlmsgAlign(nlmsgLength(n))),
			Type:     Error,
			Sequence: 1,
			PID:      uint32(os.Getpid()),
		},
		Data: make([]byte, n),
	}
	resb, err := res.MarshalBinary()
	if err != nil {
		t.Fatalf("failed to marshal response to binary: %v", err)
	}

	c, s := testLinuxConn(t, nil)

	from := &unix.SockaddrNetlink{
		Family: unix.AF_NETLINK,
	}

	s.recvmsg.p = resb
	s.recvmsg.from = from

	if _, err := c.Receive(); err != nil {
		t.Fatalf("failed to receive messages: %v", err)
	}

	// Expect several MSG_PEEK and then no flags
	want := []int{
		unix.MSG_PEEK,
		unix.MSG_PEEK,
		unix.MSG_PEEK,
		unix.MSG_PEEK,
		0,
	}

	if got := s.recvmsg.flags; !reflect.DeepEqual(want, got) {
		t.Fatalf("unexpected number recvmsg flags:\n- want: %v\n-  got: %v",
			want, got)
	}
}

func TestLinuxConnReceiveMultipleMessagesLastUnaligned(t *testing.T) {
	// This test checks if syscall.ParseNetlinkMessage allows the final
	// message in a sequence to be unaligned.  Apparently, nfnetlink can
	// do this at times.
	//
	// Reference: https://golang.org/cl/35531/.

	c, s := testLinuxConn(t, nil)

	s.recvmsg.p = []byte{
		// First message, aligned
		0x10, 0x00, 0x00, 0x00,
		0x00, 0x00,
		0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		// Final message, unaligned
		0x11, 0x00, 0x00, 0x00,
		0x00, 0x00,
		0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0xff,
	}

	s.recvmsg.from = &unix.SockaddrNetlink{
		Family: unix.AF_NETLINK,
	}

	msgs, err := c.Receive()
	if err != nil {
		t.Fatalf("failed to receive messages: %v", err)
	}

	// TODO(mdlayher): check received messages.
	_ = msgs
}

func TestLinuxConnJoinLeaveGroup(t *testing.T) {
	c, s := testLinuxConn(t, nil)

	group := uint32(1)

	if err := c.JoinGroup(group); err != nil {
		t.Fatalf("failed to join group: %v", err)
	}

	if err := c.LeaveGroup(group); err != nil {
		t.Fatalf("failed to leave group: %v", err)
	}

	want := []setSockopt{
		{
			level: unix.SOL_NETLINK,
			opt:   unix.NETLINK_ADD_MEMBERSHIP,
			value: int(group),
		},
		{
			level: unix.SOL_NETLINK,
			opt:   unix.NETLINK_DROP_MEMBERSHIP,
			value: int(group),
		},
	}

	if got := s.setSockopt; !reflect.DeepEqual(want, got) {
		t.Fatalf("unexpected socket options:\n- want: %v\n-  got: %v",
			want, got)
	}
}

func TestLinuxConnSetOption(t *testing.T) {
	tests := []struct {
		name   string
		option ConnOption
		enable bool

		want setSockopt
		err  error
	}{
		{
			name:   "invalid",
			option: 999,
			enable: true,
			err:    os.NewSyscallError("setsockopt", unix.ENOPROTOOPT),
		},
		{
			name:   "packet info on",
			option: PacketInfo,
			enable: true,
			want: setSockopt{
				opt:   unix.NETLINK_PKTINFO,
				value: 1,
			},
		},
		{
			name:   "packet info off",
			option: PacketInfo,
			enable: false,
			want: setSockopt{
				opt:   unix.NETLINK_PKTINFO,
				value: 0,
			},
		},
		{
			name:   "broadcast error",
			option: BroadcastError,
			enable: true,
			want: setSockopt{
				opt:   unix.NETLINK_BROADCAST_ERROR,
				value: 1,
			},
		},
		{
			name:   "no ENOBUFS",
			option: NoENOBUFS,
			enable: true,
			want: setSockopt{
				opt:   unix.NETLINK_NO_ENOBUFS,
				value: 1,
			},
		},
		{
			name:   "listen all NSID",
			option: ListenAllNSID,
			enable: true,
			want: setSockopt{
				opt:   unix.NETLINK_LISTEN_ALL_NSID,
				value: 1,
			},
		},
		{
			name:   "cap acknowledge",
			option: CapAcknowledge,
			enable: true,
			want: setSockopt{
				opt:   unix.NETLINK_CAP_ACK,
				value: 1,
			},
		},
		{
			name:   "extended ACK reporting",
			option: ExtendedAcknowledge,
			enable: true,
			want: setSockopt{
				opt:   unix.NETLINK_EXT_ACK,
				value: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, s := testLinuxConn(t, nil)

			// Pre-populate fixed values.
			tt.want.level = unix.SOL_NETLINK

			if err := c.SetOption(tt.option, tt.enable); err != nil {
				if want, got := tt.err, err; !reflect.DeepEqual(want, got) {
					t.Fatalf("unexpected error:\n- want: %v\n-  got: %v",
						want, got)
				}

				return
			}

			if want, got := []setSockopt{tt.want}, s.setSockopt; !reflect.DeepEqual(want, got) {
				t.Fatalf("unexpected socket options:\n- want: %v\n-  got: %v",
					want, got)
			}
		})
	}
}

func TestLinuxConnSetDeadlines(t *testing.T) {
	c, s := testLinuxConn(t, nil)

	rwd := time.Now().Add(1 * time.Second)
	if err := c.SetDeadline(rwd); err != nil {
		t.Fatalf("failed to set deadline: %v", err)
	}
	if !s.deadline.Equal(rwd) {
		t.Fatalf("set deadline %v, want %v", s.deadline, rwd)
	}

	rd := time.Now().Add(2 * time.Second)
	if err := c.SetReadDeadline(rd); err != nil {
		t.Fatalf("failed to set read deadline: %v", err)
	}
	if !s.readDeadline.Equal(rd) {
		t.Fatalf("set read deadline to %v, want %v", s.readDeadline, rd)
	}

	wd := time.Now().Add(1 * time.Second)
	if err := c.SetWriteDeadline(wd); err != nil {
		t.Fatalf("failed to set write deadline: %v", err)
	}
	if !s.writeDeadline.Equal(wd) {
		t.Fatalf("set write deadline to %v, want %v", s.writeDeadline, wd)
	}
}

func TestLinuxConnSetBuffers(t *testing.T) {
	c, s := testLinuxConn(t, nil)

	n := 64

	if err := c.SetReadBuffer(n); err != nil {
		t.Fatalf("failed to set read buffer size: %v", err)
	}

	if err := c.SetWriteBuffer(n); err != nil {
		t.Fatalf("failed to set write buffer size: %v", err)
	}

	want := []setSockopt{
		{
			level: unix.SOL_SOCKET,
			opt:   unix.SO_RCVBUF,
			value: n,
		},
		{
			level: unix.SOL_SOCKET,
			opt:   unix.SO_SNDBUF,
			value: n,
		},
	}

	if got := s.setSockopt; !reflect.DeepEqual(want, got) {
		t.Fatalf("unexpected socket options:\n- want: %v\n-  got: %v",
			want, got)
	}
}

func TestLinuxConnConfig(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		groups uint32
	}{
		{
			name:   "Default Config",
			config: &Config{},
			groups: 0x0,
		},
		{
			name:   "Config with Groups RTMGRP_IPV4_IFADDR",
			config: &Config{Groups: 0x10},
			groups: 0x10,
		},
		{
			name:   "Config with Groups RTMGRP_IPV4_IFADDR | RTMGRP_IPV4_ROUTE",
			config: &Config{Groups: 0x10 | 0x40},
			groups: 0x50,
		},
		{
			name: "Config with DisableNSLockThread and Groups RTMGRP_IPV4_IFADDR | RTMGRP_IPV4_ROUTE",
			config: &Config{
				Groups:              0x10 | 0x40,
				DisableNSLockThread: true,
			},
			groups: 0x50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testLinuxConn(t, tt.config)

			if want, got := tt.groups, c.sa.Groups; want != got {
				t.Fatalf("unexpected error:\n- want: %v\n-  got: %v",
					want, got)
			}
		})
	}
}

func Test_newLockedNetNSGoroutineNetNSDisabled(t *testing.T) {
	tests := []struct {
		name       string
		ns         int
		ok         bool
		lockThread bool
	}{
		{
			// Network namespaces are disabled but none is set: this should
			// succeed.
			name:       "not set",
			ok:         true,
			lockThread: true,
		},
		{
			// Network namespaces are disabled but one is set explicitly:
			// this should fail.
			name:       "set",
			ns:         1,
			lockThread: true,
		},
		{
			// thread locking is disabled but an ns is provided.
			// this should fail.
			name:       "disable lock thread with ns defined",
			ns:         1,
			lockThread: false,
		},
		{
			// thread locking is disabled but an ns is not provided.
			// this should succeed.
			name:       "disable lock thread without ns defined",
			lockThread: false,
			ok:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := newLockedNetNSGoroutine(tt.ns, func() (*netNS, error) {
				// Network namespaces should be disabled due to a non-existent
				// file.
				return fileNetNS("/netlinktestdoesnotexist")
			}, tt.lockThread)
			if err != nil {
				if tt.ok {
					t.Fatalf("failed to create goroutine: %v", err)
				}

				return
			}
			defer g.stop()

			if !tt.ok {
				t.Fatal("expected an error, but none occurred")
			}
		})
	}
}

func testLinuxConn(t *testing.T, config *Config) (*conn, *testSocket) {
	s := &testSocket{}
	c, _, err := bind(s, config)
	if err != nil {
		t.Fatalf("failed to bind: %v", err)
	}

	return c, s
}

type testSocket struct {
	bind           unix.Sockaddr
	bindErr        error
	closed         bool
	getsockname    unix.Sockaddr
	getsocknameErr error
	sendmsg        struct {
		p     []byte
		oob   []byte
		to    unix.Sockaddr
		flags int
	}
	recvmsg struct {
		// Received from caller
		flags []int
		// Sent to caller
		p         []byte
		oob       []byte
		recvflags int
		from      unix.Sockaddr
	}
	deadline      time.Time
	readDeadline  time.Time
	writeDeadline time.Time
	setSockopt    []setSockopt
}

type setSockopt struct {
	level int
	opt   int
	value int
}

func (s *testSocket) Bind(sa unix.Sockaddr) error {
	s.bind = sa
	return s.bindErr
}

func (s *testSocket) Close() error {
	s.closed = true
	return nil
}

func (s *testSocket) FD() int { return 0 }

func (s *testSocket) File() *os.File { return nil }

func (s *testSocket) Getsockname() (unix.Sockaddr, error) {
	if s.getsockname == nil {
		return &unix.SockaddrNetlink{}, s.getsocknameErr
	}

	return s.getsockname, s.getsocknameErr
}

func (s *testSocket) Recvmsg(p, oob []byte, flags int) (int, int, int, unix.Sockaddr, error) {
	s.recvmsg.flags = append(s.recvmsg.flags, flags)
	n := copy(p, s.recvmsg.p)
	oobn := copy(oob, s.recvmsg.oob)

	return n, oobn, s.recvmsg.recvflags, s.recvmsg.from, nil
}

func (s *testSocket) Sendmsg(p, oob []byte, to unix.Sockaddr, flags int) error {
	s.sendmsg.p = p
	s.sendmsg.oob = oob
	s.sendmsg.to = to
	s.sendmsg.flags = flags
	return nil
}

func (s *testSocket) SetDeadline(t time.Time) error {
	s.deadline = t
	return nil
}

func (s *testSocket) SetReadDeadline(t time.Time) error {
	s.readDeadline = t
	return nil
}

func (s *testSocket) SetWriteDeadline(t time.Time) error {
	s.writeDeadline = t
	return nil
}

func (s *testSocket) SetSockoptInt(level, opt, value int) error {
	// Value must be in range of a C integer.
	if value < math.MinInt32 || value > math.MaxInt32 {
		return unix.EINVAL
	}

	s.setSockopt = append(s.setSockopt, setSockopt{
		level: level,
		opt:   opt,
		value: value,
	})

	return nil
}

func (s *testSocket) SetSockoptSockFprog(_, _ int, _ *unix.SockFprog) error {
	panic("netlink: testSocket.SetSockoptSockFprog not currently implemented")
}
