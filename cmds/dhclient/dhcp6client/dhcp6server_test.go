// Copyright 2015 dhcp6 Author. All Rights Reserved.
// https://github.com/mdlayher/dhcp6
// Author: Matt Layher <mdlayher@gmail.com>

package dhcp6client

import (
	"bytes"
	"encoding/hex"
	"errors"
	"log"
	"net"
	"time"

	"github.com/mdlayher/dhcp6"
	"golang.org/x/net/ipv6"
)

var errClosing = errors.New("use of closed network connection")

type PacketConn dhcp6.PacketConn

/**
* testMessage that implements messages between server and client.
 */
type testMessage struct {
	b    bytes.Buffer
	cm   *ipv6.ControlMessage
	addr net.Addr
}

/**
* testPacketConn that implement interface dhcp6.PacketConn.
 */

type testPacketConn struct {
	r *testMessage
	w *testMessage

	closed bool
	joined []net.Addr
	left   []net.Addr
	flags  map[ipv6.ControlFlags]bool
}

func (c *testPacketConn) Close() error {
	c.closed = true
	return nil
}

func (c *testPacketConn) JoinGroup(iface *net.Interface, group net.Addr) error {
	c.joined = append(c.joined, group)
	return nil
}

func (c *testPacketConn) LeaveGroup(iface *net.Interface, group net.Addr) error {
	c.left = append(c.left, group)
	return nil
}

func (c *testPacketConn) SetControlMessage(cf ipv6.ControlFlags, on bool) error {
	c.flags[cf] = on
	return nil
}

func (c *testPacketConn) ReadFrom(b []byte) (int, *ipv6.ControlMessage, net.Addr, error) {
	n, err := c.r.b.Read(b)
	cm := c.r.cm
	src := c.r.addr

	return n, cm, src, err
}

func (c *testPacketConn) WriteTo(b []byte, cm *ipv6.ControlMessage, dst net.Addr) (int, error) {
	n, err := c.w.b.Write(b)
	c.w.cm = cm
	c.w.addr = dst

	return n, err
}

/**
* oneReadPacketConn that allows the server to read only once.
 */
type oneReadPacketConn struct {
	PacketConn

	err    error
	txDone bool

	// Once the read/write is done, channel will close and stop blocking.
	readDoneChan  chan struct{}
	writeDoneChan chan struct{}
}

func (c *oneReadPacketConn) ReadFrom(b []byte) (int, *ipv6.ControlMessage, net.Addr, error) {
	if c.txDone {
		return 0, nil, nil, errClosing
	}
	c.txDone = true
	n, cm, addr, err := c.PacketConn.ReadFrom(b)
	close(c.readDoneChan)
	log.Printf("read from: %v, %v, %v, %v\n", n, cm, addr, err)
	return n, cm, addr, err
}

func (c *oneReadPacketConn) WriteTo(b []byte, cm *ipv6.ControlMessage, addr net.Addr) (int, error) {
	n, err := c.PacketConn.WriteTo(b, cm, addr)
	close(c.writeDoneChan)
	return n, err
}

/**
* Serve and Handle
 */
// ServeDHCP is a dhcp6.Handler which invokes an internal handler that
// allows errors to be returned and handled in one place.
func (h *Handler) ServeDHCP(w dhcp6.ResponseSender, r *dhcp6.Request) {
	if err := h.handler(h.ip, w, r); err != nil {
		log.Println(err)
	}
}

type Handler struct {
	ip      net.IP
	handler handler
}

type handler func(ip net.IP, w dhcp6.ResponseSender, r *dhcp6.Request) error

// handle is a handler which assigns IPv6 addresses using DHCPv6.
func handle(ip net.IP, w dhcp6.ResponseSender, r *dhcp6.Request) error {
	// Accept only Solicit, Request, or Confirm, since this server
	// does not handle Information Request or other message types
	valid := map[dhcp6.MessageType]struct{}{
		dhcp6.MessageTypeSolicit: struct{}{},
		dhcp6.MessageTypeRequest: struct{}{},
		dhcp6.MessageTypeConfirm: struct{}{},
	}
	if _, ok := valid[r.MessageType]; !ok {
		return nil
	}

	// Make sure client sent a client ID
	duid, ok := r.Options.Get(dhcp6.OptionClientID)
	if !ok {
		return nil
	}

	// Log information about the incoming request.
	log.Printf("[%s] id: %s, type: %d, len: %d, tx: %s",
		hex.EncodeToString(duid),
		r.RemoteAddr,
		r.MessageType,
		r.Length,
		hex.EncodeToString(r.TransactionID[:]),
	)

	// Print out options the client has requested
	if opts, ok, err := r.Options.OptionRequest(); err == nil && ok {
		log.Println("\t- requested:")
		for _, o := range opts {
			log.Printf("\t\t - %s", o)
		}
	}

	// Client must send a IANA to retrieve an IPv6 address
	ianas, ok, err := r.Options.IANA()
	if err != nil {
		return err
	}
	if !ok {
		log.Println("no IANAs provided")
		return nil
	}

	// Only accept one IANA
	if len(ianas) > 1 {
		log.Println("can only handle one IANA")
		return nil
	}
	ia := ianas[0]

	log.Printf("\tIANA: %s (%s, %s), opts: %v",
		hex.EncodeToString(ia.IAID[:]),
		ia.T1,
		ia.T2,
		ia.Options,
	)

	// Instruct client to prefer this server unconditionally
	_ = w.Options().Add(dhcp6.OptionPreference, dhcp6.Preference(255))

	// IANA may already have an IAAddr if an address was already assigned.
	// If not, assign a new one.
	iaaddrs, ok, err := ia.Options.IAAddr()
	if err != nil {
		return err
	}

	// Client did not indicate a previous address, and is soliciting.
	// Advertise a new IPv6 address.
	if !ok && r.MessageType == dhcp6.MessageTypeSolicit {
		return newIAAddr(ia, ip, w, r)
	} else if !ok {
		// Client did not indicate an address and is not soliciting.  Ignore.
		return nil
	}

	// Confirm or renew an existing IPv6 address

	// Must have an IAAddr, but we ignore if more than one is present
	if len(iaaddrs) == 0 {
		return nil
	}
	iaa := iaaddrs[0]

	log.Printf("\t\tIAAddr: %s (%s, %s), opts: %v",
		iaa.IP,
		iaa.PreferredLifetime,
		iaa.ValidLifetime,
		iaa.Options,
	)

	// Add IAAddr inside IANA, add IANA to options
	_ = ia.Options.Add(dhcp6.OptionIAAddr, iaa)
	_ = w.Options().Add(dhcp6.OptionIANA, ia)

	// Send reply to client
	_, err = w.Send(dhcp6.MessageTypeReply)
	return err
}

// newIAAddr creates a IAAddr for a IANA using the specified IPv6 address,
// and advertises it to a client.
func newIAAddr(ia *dhcp6.IANA, ip net.IP, w dhcp6.ResponseSender, r *dhcp6.Request) error {
	// Send IPv6 address with 60 second preferred lifetime,
	// 90 second valid lifetime, no extra options
	iaaddr, err := dhcp6.NewIAAddr(ip, 60*time.Second, 90*time.Second, nil)
	if err != nil {
		return err
	}

	// Add IAAddr inside IANA, add IANA to options
	_ = ia.Options.Add(dhcp6.OptionIAAddr, iaaddr)
	_ = w.Options().Add(dhcp6.OptionIANA, ia)

	// Advertise address to soliciting clients
	log.Printf("advertising IP: %s", ip)
	_, err = w.Send(dhcp6.MessageTypeAdvertise)
	return err
}

// serve calls server.Serve() that servers a request from client
func serve(r *testMessage) (*dhcp6.Packet, error) {
	s := &dhcp6.Server{}
	s.Iface = &net.Interface{
		Name:  "foo",
		Index: 0,
	}
	s.Handler = &Handler{
		ip:      net.ParseIP("::0"),
		handler: handle,
	}

	tc := &testPacketConn{
		r: r,
		w: &testMessage{},

		closed: false,
		joined: make([]net.Addr, 0),
		left:   make([]net.Addr, 0),
		flags:  make(map[ipv6.ControlFlags]bool),
	}

	c := &oneReadPacketConn{
		PacketConn: tc,
		txDone:     false,

		readDoneChan:  make(chan struct{}, 0),
		writeDoneChan: make(chan struct{}, 0),
	}

	err := s.Serve(c)
	<-c.readDoneChan
	<-c.writeDoneChan

	if err == errClosing {
		p := &dhcp6.Packet{}
		if err0 := p.UnmarshalBinary(tc.w.b.Bytes()); err0 != nil {
			return nil, err0
		}
		return p, nil
	}
	return nil, err
}
