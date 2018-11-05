// Copyright 2018 the u-root Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dhcp4server

import (
	"fmt"
	"log"
	"net"

	"github.com/u-root/dhcp4"
	"github.com/u-root/dhcp4/dhcp4opts"
)

type macAddr [16]byte

const maxMessageSize = 1500

type Server struct {
	// whoami
	ip net.IP

	ips   *ipAllocator
	conns map[macAddr]net.IP

	sname, filename string
}

func New(ip net.IP, subnet *net.IPNet, sname, filename string) *Server {
	return &Server{
		ip:       ip.To4(),
		ips:      newIPAllocator(subnet),
		conns:    make(map[macAddr]net.IP),
		sname:    sname,
		filename: filename,
	}
}

func (s *Server) responsePacket(request *dhcp4.Packet, typ dhcp4opts.DHCPMessageType) *dhcp4.Packet {
	packet := dhcp4.NewPacket(dhcp4.BootReply)
	packet.HType = request.HType

	// Must include the following according to RFC 2131 Section 4.3.1 Table 3.
	packet.TransactionID = request.TransactionID
	packet.CHAddr = request.CHAddr
	packet.Broadcast = request.Broadcast
	packet.GIAddr = request.GIAddr
	packet.Options.Add(dhcp4.OptionDHCPMessageType, typ)
	packet.Options.Add(dhcp4.OptionServerIdentifier, dhcp4opts.IP(s.ip))

	// IP of next bootstrap server (us).
	packet.SIAddr = s.ip

	// Optional.
	packet.Options.Add(dhcp4.OptionMaximumDHCPMessageSize, dhcp4opts.Uint16(maxMessageSize))
	return packet
}

func (s *Server) getIP(haddr net.HardwareAddr) net.IP {
	mac := getMac(haddr)

	// Already allocated an IP to this client.
	if ip, ok := s.conns[mac]; ok {
		return ip
	}
	return nil
}

func (s *Server) newIP(haddr net.HardwareAddr) net.IP {
	mac := getMac(haddr)

	// Already allocated an IP to this client.
	if ip, ok := s.conns[mac]; ok {
		return ip
	}

	ip := s.ips.alloc()
	if ip == nil {
		return nil
	}
	s.conns[mac] = ip
	return ip
}

func (s *Server) grabIP(haddr net.HardwareAddr, ip net.IP) bool {
	mac := getMac(haddr)

	if _, ok := s.conns[mac]; ok {
		return false
	}

	if !s.ips.grab(ip) {
		return false
	}
	s.conns[mac] = ip
	return true
}

func getMac(haddr net.HardwareAddr) macAddr {
	var mac macAddr
	copy(mac[:], haddr)
	return mac
}

func (s *Server) release(haddr net.HardwareAddr) {
	mac := getMac(haddr)
	ip, ok := s.conns[mac]
	if !ok {
		return
	}
	delete(s.conns, mac)
	s.ips.free(ip)
}

func (s *Server) writePacket(conn net.PacketConn, addr net.Addr, p *dhcp4.Packet) error {
	pkt, err := p.MarshalBinary()
	if err != nil {
		return err
	}

	uaddr, ok := addr.(*net.UDPAddr)
	if !ok {
		return fmt.Errorf("DHCP send: addr %v is not a UDP address", addr)
	}
	if uaddr.IP.Equal(net.IPv4zero) {
		// Broadcast instead
		_, err = conn.WriteTo(pkt, &net.UDPAddr{IP: net.IPv4bcast, Port: uaddr.Port})
		return err
	}
	_, err = conn.WriteTo(pkt, addr)
	return err
}

func (s *Server) Serve(logger *log.Logger, conn net.PacketConn) error {
	var buf [maxMessageSize]byte
	for {
		n, addr, err := conn.ReadFrom(buf[:])
		if err != nil {
			return err
		}

		pkt, err := dhcp4.ParsePacket(buf[:n])
		if err != nil {
			logger.Printf("Invalid DHCP packet from %v: %v", addr, err)
			continue
		}

		switch typ := dhcp4opts.GetDHCPMessageType(pkt.Options); typ {
		case dhcp4opts.DHCPDiscover:
			offer := s.responsePacket(pkt, dhcp4opts.DHCPOffer)
			mac := getMac(pkt.CHAddr)

			if ip, ok := s.conns[mac]; ok {
				// Already has an IP allocated.
				offer.YIAddr = ip
			} else if rip := dhcp4opts.GetRequestedIPAddress(pkt.Options); s.grabIP(pkt.CHAddr, net.IP(rip)) {
				// Requested IP is available.
				offer.YIAddr = net.IP(rip)
			} else if ip := s.newIP(pkt.CHAddr); ip != nil {
				// Grab a random new IP.
				offer.YIAddr = ip
			}

			offer.ServerName = s.sname
			offer.BootFile = s.filename
			if offer.YIAddr != nil {
				if err := s.writePacket(conn, addr, offer); err != nil {
					// TODO Undo address assignment.
					return err
				}
			} else {
				// TODO: send rejection.
			}

		case dhcp4opts.DHCPRequest:
			offered := s.getIP(pkt.CHAddr)

			rip := dhcp4opts.GetRequestedIPAddress(pkt.Options)
			var re *dhcp4.Packet
			if !net.IP(rip).Equal(offered) {
				// Client is confused about IP offered?
				re = s.responsePacket(pkt, dhcp4opts.DHCPNAK)
			} else {
				re = s.responsePacket(pkt, dhcp4opts.DHCPACK)
				re.CIAddr = pkt.CIAddr
				re.YIAddr = offered
				re.ServerName = s.sname
				re.BootFile = s.filename
			}

			if err := s.writePacket(conn, addr, re); err != nil {
				// TODO: Undo address assignment.
				return err
			}

		case dhcp4opts.DHCPDecline, dhcp4opts.DHCPRelease:
			// TODO
			s.release(pkt.CHAddr)

		case dhcp4opts.DHCPInform:
			// TODO

		case dhcp4opts.DHCPOffer, dhcp4opts.DHCPACK, dhcp4opts.DHCPNAK:
			// DHCP servers ignore these according to RFC 2131,
			// Section 4.3.

		default:
			logger.Printf("DHCP message with unknown type %v", typ)
			continue
		}
	}
}
