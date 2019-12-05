// Package client4 is deprecated. Use "nclient4" instead.
package client4

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"reflect"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"golang.org/x/net/ipv4"
	"golang.org/x/sys/unix"
)

// MaxUDPReceivedPacketSize is the (arbitrary) maximum UDP packet size supported
// by this library. Theoretically could be up to 65kb.
const (
	MaxUDPReceivedPacketSize = 8192
)

var (
	// DefaultReadTimeout is the time to wait after listening in which the
	// exchange is considered failed.
	DefaultReadTimeout = 3 * time.Second

	// DefaultWriteTimeout is the time to wait after sending in which the
	// exchange is considered failed.
	DefaultWriteTimeout = 3 * time.Second
)

// Client is the object that actually performs the DHCP exchange. It currently
// only has read and write timeout values, plus (optional) local and remote
// addresses.
type Client struct {
	ReadTimeout, WriteTimeout time.Duration
	RemoteAddr                net.Addr
	LocalAddr                 net.Addr
}

// NewClient generates a new client to perform a DHCP exchange with, setting the
// read and write timeout fields to defaults.
func NewClient() *Client {
	return &Client{
		ReadTimeout:  DefaultReadTimeout,
		WriteTimeout: DefaultWriteTimeout,
	}
}

// MakeRawUDPPacket converts a payload (a serialized DHCPv4 packet) into a
// raw UDP packet for the specified serverAddr from the specified clientAddr.
func MakeRawUDPPacket(payload []byte, serverAddr, clientAddr net.UDPAddr) ([]byte, error) {
	udp := make([]byte, 8)
	binary.BigEndian.PutUint16(udp[:2], uint16(clientAddr.Port))
	binary.BigEndian.PutUint16(udp[2:4], uint16(serverAddr.Port))
	binary.BigEndian.PutUint16(udp[4:6], uint16(8+len(payload)))
	binary.BigEndian.PutUint16(udp[6:8], 0) // try to offload the checksum

	h := ipv4.Header{
		Version:  4,
		Len:      20,
		TotalLen: 20 + len(udp) + len(payload),
		TTL:      64,
		Protocol: 17, // UDP
		Dst:      serverAddr.IP,
		Src:      clientAddr.IP,
	}
	ret, err := h.Marshal()
	if err != nil {
		return nil, err
	}
	ret = append(ret, udp...)
	ret = append(ret, payload...)
	return ret, nil
}

// makeRawSocket creates a socket that can be passed to unix.Sendto.
func makeRawSocket(ifname string) (int, error) {
	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_RAW, unix.IPPROTO_RAW)
	if err != nil {
		return fd, err
	}
	err = unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_REUSEADDR, 1)
	if err != nil {
		return fd, err
	}
	err = unix.SetsockoptInt(fd, unix.IPPROTO_IP, unix.IP_HDRINCL, 1)
	if err != nil {
		return fd, err
	}
	err = dhcpv4.BindToInterface(fd, ifname)
	if err != nil {
		return fd, err
	}
	return fd, nil
}

// MakeBroadcastSocket creates a socket that can be passed to unix.Sendto
// that will send packets out to the broadcast address.
func MakeBroadcastSocket(ifname string) (int, error) {
	fd, err := makeRawSocket(ifname)
	if err != nil {
		return fd, err
	}
	err = unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_BROADCAST, 1)
	if err != nil {
		return fd, err
	}
	return fd, nil
}

// MakeListeningSocket creates a listening socket on 0.0.0.0 for the DHCP client
// port and returns it.
func MakeListeningSocket(ifname string) (int, error) {
	return makeListeningSocketWithCustomPort(ifname, dhcpv4.ClientPort)
}

func htons(v uint16) uint16 {
	var tmp [2]byte
	binary.BigEndian.PutUint16(tmp[:], v)
	return binary.LittleEndian.Uint16(tmp[:])
}

func makeListeningSocketWithCustomPort(ifname string, port int) (int, error) {
	fd, err := unix.Socket(unix.AF_PACKET, unix.SOCK_DGRAM, int(htons(unix.ETH_P_IP)))
	if err != nil {
		return fd, err
	}
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		return fd, err
	}
	llAddr := unix.SockaddrLinklayer{
		Ifindex:  iface.Index,
		Protocol: htons(unix.ETH_P_IP),
	}
	err = unix.Bind(fd, &llAddr)
	return fd, err
}

func toUDPAddr(addr net.Addr, defaultAddr *net.UDPAddr) (*net.UDPAddr, error) {
	var uaddr *net.UDPAddr
	if addr == nil {
		uaddr = defaultAddr
	} else {
		if addr, ok := addr.(*net.UDPAddr); ok {
			uaddr = addr
		} else {
			return nil, fmt.Errorf("could not convert to net.UDPAddr, got %v instead", reflect.TypeOf(addr))
		}
	}
	if uaddr.IP.To4() == nil {
		return nil, fmt.Errorf("'%s' is not a valid IPv4 address", uaddr.IP)
	}
	return uaddr, nil
}

func (c *Client) getLocalUDPAddr() (*net.UDPAddr, error) {
	defaultLocalAddr := &net.UDPAddr{IP: net.IPv4zero, Port: dhcpv4.ClientPort}
	laddr, err := toUDPAddr(c.LocalAddr, defaultLocalAddr)
	if err != nil {
		return nil, fmt.Errorf("Invalid local address: %s", err)
	}
	return laddr, nil
}

func (c *Client) getRemoteUDPAddr() (*net.UDPAddr, error) {
	defaultRemoteAddr := &net.UDPAddr{IP: net.IPv4bcast, Port: dhcpv4.ServerPort}
	raddr, err := toUDPAddr(c.RemoteAddr, defaultRemoteAddr)
	if err != nil {
		return nil, fmt.Errorf("Invalid remote address: %s", err)
	}
	return raddr, nil
}

// Exchange runs a full DORA transaction: Discover, Offer, Request, Acknowledge,
// over UDP. Does not retry in case of failures. Returns a list of DHCPv4
// structures representing the exchange. It can contain up to four elements,
// ordered as Discovery, Offer, Request and Acknowledge. In case of errors, an
// error is returned, and the list of DHCPv4 objects will be shorted than 4,
// containing all the sent and received DHCPv4 messages.
func (c *Client) Exchange(ifname string, modifiers ...dhcpv4.Modifier) ([]*dhcpv4.DHCPv4, error) {
	conversation := make([]*dhcpv4.DHCPv4, 0)
	raddr, err := c.getRemoteUDPAddr()
	if err != nil {
		return nil, err
	}
	laddr, err := c.getLocalUDPAddr()
	if err != nil {
		return nil, err
	}
	// Get our file descriptor for the raw socket we need.
	var sfd int
	// If the address is not net.IPV4bcast, use a unicast socket. This should
	// cover the majority of use cases, but we're essentially ignoring the fact
	// that the IP could be the broadcast address of a specific subnet.
	if raddr.IP.Equal(net.IPv4bcast) {
		sfd, err = MakeBroadcastSocket(ifname)
	} else {
		sfd, err = makeRawSocket(ifname)
	}
	if err != nil {
		return conversation, err
	}
	rfd, err := makeListeningSocketWithCustomPort(ifname, laddr.Port)
	if err != nil {
		return conversation, err
	}

	defer func() {
		// close the sockets
		if err := unix.Close(sfd); err != nil {
			log.Printf("unix.Close(sendFd) failed: %v", err)
		}
		if sfd != rfd {
			if err := unix.Close(rfd); err != nil {
				log.Printf("unix.Close(recvFd) failed: %v", err)
			}
		}
	}()

	// Discover
	discover, err := dhcpv4.NewDiscoveryForInterface(ifname, modifiers...)
	if err != nil {
		return conversation, err
	}
	conversation = append(conversation, discover)

	// Offer
	offer, err := c.SendReceive(sfd, rfd, discover, dhcpv4.MessageTypeOffer)
	if err != nil {
		return conversation, err
	}
	conversation = append(conversation, offer)

	// Request
	request, err := dhcpv4.NewRequestFromOffer(offer, modifiers...)
	if err != nil {
		return conversation, err
	}
	conversation = append(conversation, request)

	// Ack
	ack, err := c.SendReceive(sfd, rfd, request, dhcpv4.MessageTypeAck)
	if err != nil {
		return conversation, err
	}
	conversation = append(conversation, ack)

	return conversation, nil
}

// SendReceive sends a packet (with some write timeout) and waits for a
// response up to some read timeout value. If the message type is not
// MessageTypeNone, it will wait for a specific message type
func (c *Client) SendReceive(sendFd, recvFd int, packet *dhcpv4.DHCPv4, messageType dhcpv4.MessageType) (*dhcpv4.DHCPv4, error) {
	raddr, err := c.getRemoteUDPAddr()
	if err != nil {
		return nil, err
	}
	laddr, err := c.getLocalUDPAddr()
	if err != nil {
		return nil, err
	}
	packetBytes, err := MakeRawUDPPacket(packet.ToBytes(), *raddr, *laddr)
	if err != nil {
		return nil, err
	}

	// Create a goroutine to perform the blocking send, and time it out after
	// a certain amount of time.
	var (
		destination [net.IPv4len]byte
		response    *dhcpv4.DHCPv4
	)
	copy(destination[:], raddr.IP.To4())
	remoteAddr := unix.SockaddrInet4{Port: laddr.Port, Addr: destination}
	recvErrors := make(chan error, 1)
	go func(errs chan<- error) {
		// set read timeout
		timeout := unix.NsecToTimeval(c.ReadTimeout.Nanoseconds())
		if innerErr := unix.SetsockoptTimeval(recvFd, unix.SOL_SOCKET, unix.SO_RCVTIMEO, &timeout); innerErr != nil {
			errs <- innerErr
			return
		}
		for {
			buf := make([]byte, MaxUDPReceivedPacketSize)
			n, _, innerErr := unix.Recvfrom(recvFd, buf, 0)
			if innerErr != nil {
				errs <- innerErr
				return
			}

			var iph ipv4.Header
			if err := iph.Parse(buf[:n]); err != nil {
				// skip non-IP data
				continue
			}
			if iph.Protocol != 17 {
				// skip non-UDP packets
				continue
			}
			udph := buf[iph.Len:n]
			// check source and destination ports
			srcPort := int(binary.BigEndian.Uint16(udph[0:2]))
			expectedSrcPort := dhcpv4.ServerPort
			if c.RemoteAddr != nil {
				expectedSrcPort = c.RemoteAddr.(*net.UDPAddr).Port
			}
			if srcPort != expectedSrcPort {
				continue
			}
			dstPort := int(binary.BigEndian.Uint16(udph[2:4]))
			expectedDstPort := dhcpv4.ClientPort
			if c.RemoteAddr != nil {
				expectedDstPort = c.LocalAddr.(*net.UDPAddr).Port
			}
			if dstPort != expectedDstPort {
				continue
			}
			// UDP checksum is not checked
			pLen := int(binary.BigEndian.Uint16(udph[4:6]))
			payload := buf[iph.Len+8 : iph.Len+8+pLen]

			response, innerErr = dhcpv4.FromBytes(payload)
			if innerErr != nil {
				errs <- innerErr
				return
			}
			// check that this is a response to our message
			if response.TransactionID != packet.TransactionID {
				continue
			}
			// wait for a response message
			if response.OpCode != dhcpv4.OpcodeBootReply {
				continue
			}
			// if we are not requested to wait for a specific message type,
			// return what we have
			if messageType == dhcpv4.MessageTypeNone {
				break
			}
			// break if it's a reply of the desired type, continue otherwise
			if response.MessageType() == messageType {
				break
			}
		}
		recvErrors <- nil
	}(recvErrors)

	// send the request while the goroutine waits for replies
	if err = unix.Sendto(sendFd, packetBytes, 0, &remoteAddr); err != nil {
		return nil, err
	}

	select {
	case err = <-recvErrors:
		if err == unix.EAGAIN {
			return nil, errors.New("timed out while listening for replies")
		}
		if err != nil {
			return nil, err
		}
	case <-time.After(c.ReadTimeout):
		return nil, errors.New("timed out while listening for replies")
	}

	return response, nil
}
