package dhcp6client

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/mdlayher/dhcp6"
	"golang.org/x/sys/unix"
)

type packetSock struct {
	fd      int
	ifindex int
}

var bcastMAC = []byte{255, 255, 255, 255, 255, 255}

func NewPacketSock(ifindex int) (*packetSock, error) {
	fd, err := unix.Socket(unix.AF_PACKET, unix.SOCK_DGRAM, int(swap16(unix.ETH_P_IPV6)))
	if err != nil {
		return nil, err
	}
	addr := unix.SockaddrLinklayer{
		Ifindex:  ifindex,
		Protocol: swap16(unix.ETH_P_IPV6),
	}
	if err = unix.Bind(fd, &addr); err != nil {
		return nil, err
	}
	return &packetSock{
		fd:      fd,
		ifindex: ifindex,
	}, nil
}

// Write dhcpv6 requests
func (pc *packetSock) Write(pb []byte) error {
	// Define linke layer
	lladdr := unix.SockaddrLinklayer{
		Ifindex:  pc.ifindex,
		Protocol: swap16(unix.ETH_P_IPV6),
		Halen:    uint8(len(bcastMAC)),
	}
	copy(lladdr.Addr[:], bcastMAC)

	// Wrap up request
	//req, err := dhcp6.ParseRequest(pb, addr)
	//if err != nil {
	//	return err
	//}

	//rb, err := MarshalBinary(req)
	//if err != nil {
	//	return err
	//}
	// Send out request from link layer
	return unix.Sendto(pc.fd, pb, 0, &lladdr)
}

func (pc *packetSock) ReadFrom() {
	fmt.Printf("starts reading\n")
	pb := make([]byte, 200) // pkt of size 100 bytes, for now
	n, _, err := unix.Recvfrom(pc.fd, pb, 0)
	packet := dhcp6.Packet{}
	UnmarshalBinary(&packet, pb)
	fmt.Printf("response: %v\n", packet)
	fmt.Printf("read from server: %v, %v, %v\n", n, pb, err)
}

func (pc *packetSock) Close() error {
	return unix.Close(pc.fd)
}

func UnmarshalBinary(p *dhcp6.Packet, b []byte) error {
	// Packet must contain at least a message type and transaction ID
	if len(b) < 4 {
		return dhcp6.ErrInvalidPacket
	}
	p.MessageType = dhcp6.MessageType(b[0])
	txID := [3]byte{}
	copy(txID[:], b[1:4])
	p.TransactionID = txID

	options, err := parseOptions(b[4:])
	if err != nil {
		// Invalid options means an invalid packet
		return dhcp6.ErrInvalidPacket
	}
	p.Options = options
	return nil
}

func parseOptions(b []byte) (dhcp6.Options, error) {
	var length int
	options := make(dhcp6.Options)
	buf := bytes.NewBuffer(b)
	for buf.Len() > 3 {
		// 2 bytes: option code
		o := option{}
		code := dhcp6.OptionCode(binary.BigEndian.Uint16(buf.Next(2)))
		// If code is 0, bytes are empty after this point
		if code == 0 {
			return options, nil
		}

		o.Code = code
		// 2 bytes: option length
		length = int(binary.BigEndian.Uint16(buf.Next(2)))

		// If length indicated is zero, skip to next iteration
		if length == 0 {
			continue
		}

		// N bytes: option data
		o.Data = buf.Next(length)
		// Set slice's max for option's data
		o.Data = o.Data[:len(o.Data):len(o.Data)]

		// If option data has less bytes than indicated by length,
		// return an error
		if len(o.Data) < length {
			return nil, errors.New("invalid options data")
		}

		addRaw(options, o.Code, o.Data)
	}
	// Report error for any trailing bytes
	if buf.Len() != 0 {
		return nil, errors.New("invalid options data")
	}
	fmt.Printf("options: %v\n", options)
	return options, nil
}

func swap16(x uint16) uint16 {
	var b [2]byte
	binary.BigEndian.PutUint16(b[:], x)
	return binary.LittleEndian.Uint16(b[:])
}

//func MarshalBinary(req *dhcp6.Request) ([]byte, error) {
//	r := *req
//	opts := enumerate(r.Options)
//	addrbyte := []byte(r.RemoteAddr)
//	b := make([]byte, 6+opts.count()+len(addrbyte))
//	b[0] = byte(r.MessageType)
//	copy(b[1:4], r.TransactionID[:])
//	opts.write(b[4 : 4+opts.count()])
//	copy(b[4+opts.count():], addrbyte[:])
//
//	return b, nil
//}
