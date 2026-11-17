package tftp

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/net/ipv6"

	"golang.org/x/net/ipv4"
)

type connectionError struct {
	error
	timeout   bool
	temporary bool
}

func (t *connectionError) Timeout() bool {
	return t.timeout
}

func (t *connectionError) Temporary() bool {
	return t.temporary
}

type connection interface {
	sendTo([]byte, *net.UDPAddr) error
	readFrom([]byte) (int, *net.UDPAddr, error)
	setDeadline(time.Duration) error
	close()
}

type connConnection struct {
	conn *net.UDPConn
}

type chanConnection struct {
	server        *Server
	channel       chan []byte
	srcAddr, addr *net.UDPAddr
	timeout       time.Duration
}

func (c *chanConnection) sendTo(data []byte, addr *net.UDPAddr) error {
	var err error
	c.server.mu.Lock()
	defer c.server.mu.Unlock()
	if conn, ok := c.server.conn.(*net.UDPConn); ok {
		srcAddr := c.srcAddr.IP.To4()
		var cmm []byte
		if srcAddr != nil {
			cm := &ipv4.ControlMessage{Src: srcAddr}
			cmm = cm.Marshal()
		} else {
			cm := &ipv6.ControlMessage{Src: c.srcAddr.IP}
			cmm = cm.Marshal()
		}
		_, _, err = conn.WriteMsgUDP(data, cmm, c.addr)
	} else {
		_, err = c.server.conn.WriteTo(data, addr)
	}
	return err
}

func (c *chanConnection) readFrom(buffer []byte) (int, *net.UDPAddr, error) {
	select {
	case data := <-c.channel:
		copy(buffer, data)
		return len(data), c.addr, nil
	case <-time.After(c.timeout):
		return 0, nil, makeError(c.addr.String())
	}
}

func (c *chanConnection) setDeadline(deadline time.Duration) error {
	c.timeout = deadline
	return nil
}

func (c *chanConnection) close() {
	c.server.mu.Lock()
	defer c.server.mu.Unlock()
	close(c.channel)
	delete(c.server.handlers, c.addr.String())
}

func (c *connConnection) sendTo(data []byte, addr *net.UDPAddr) error {
	_, err := c.conn.WriteToUDP(data, addr)
	return err
}

func makeError(addr string) net.Error {
	error := connectionError{
		timeout:   true,
		temporary: true,
	}
	error.error = fmt.Errorf("channel timeout: %v", addr)
	return &error
}

func (c *connConnection) readFrom(buffer []byte) (int, *net.UDPAddr, error) {
	return c.conn.ReadFromUDP(buffer)
}

func (c *connConnection) setDeadline(deadline time.Duration) error {
	return c.conn.SetReadDeadline(time.Now().Add(deadline))
}

func (c *connConnection) close() {
	c.conn.Close()
}
