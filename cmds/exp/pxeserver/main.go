// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// pxeserver is a test & lab PXE server that supports TFTP, HTTP, and DHCPv4.
//
// pxeserver can either respond to *all* DHCP requests, or a DHCP request from
// a specific MAC. In either case, it will supply the same IP in all answers.
package main

import (
	"bytes"
	"flag"
	"log"
	"math"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/server4"
	"pack.ag/tftp"
)

var (
	mac = flag.String("mac", "", "MAC address to respond to. Responds to all requests if unspecified.")

	// DHCPv4-specific
	selfIP       = flag.String("ip", "192.168.0.1", "DHCPv4 IP of self")
	yourIP       = flag.String("your-ip", "192.168.0.2/24", "The one and only CIDR to give to all DHCPv4 clients")
	rootpath     = flag.String("rootpath", "", "RootPath option to serve via DHCPv4")
	bootfilename = flag.String("bootfilename", "pxelinux.0", "Boot file to serve via DHCPv4")
	inf          = flag.String("interface", "eth0", "Interface to serve DHCPv4 on")

	// File serving
	tftpDir = flag.String("tftp-dir", "", "Directory to serve over TFTP")
	httpDir = flag.String("http-dir", "", "Directory to serve over HTTP")
)

type server struct {
	mac          net.HardwareAddr
	yourIP       net.IP
	submask      net.IPMask
	self         net.IP
	bootfilename string
	rootpath     string
}

func (s *server) dhcpHandler(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
	log.Printf("Handling request %v for peer %v", m, peer)

	var replyType dhcpv4.MessageType
	switch mt := m.MessageType(); mt {
	case dhcpv4.MessageTypeDiscover:
		replyType = dhcpv4.MessageTypeOffer
	case dhcpv4.MessageTypeRequest:
		replyType = dhcpv4.MessageTypeAck
	default:
		log.Printf("Can't handle type %v", mt)
		return
	}
	if s.mac != nil && !bytes.Equal(m.ClientHWAddr, s.mac) {
		log.Printf("Not responding to DHCP request for mac %s, which does not match %s", m.ClientHWAddr, s.mac)
		return
	}

	reply, err := dhcpv4.NewReplyFromRequest(m,
		dhcpv4.WithMessageType(replyType),
		dhcpv4.WithServerIP(s.self),
		dhcpv4.WithRouter(s.self),
		dhcpv4.WithNetmask(s.submask),
		dhcpv4.WithYourIP(s.yourIP),
		// RFC 2131, Section 4.3.1. Server Identifier: MUST
		dhcpv4.WithOption(dhcpv4.OptServerIdentifier(s.self)),
		// RFC 2131, Section 4.3.1. IP lease time: MUST
		dhcpv4.WithOption(dhcpv4.OptIPAddressLeaseTime(time.Duration(math.MaxUint32)*time.Second)),
	)
	// RFC 6842, MUST include Client Identifier if client specified one.
	if val := m.Options.Get(dhcpv4.OptionClientIdentifier); len(val) > 0 {
		reply.UpdateOption(dhcpv4.OptGeneric(dhcpv4.OptionClientIdentifier, val))
	}
	if len(s.bootfilename) > 0 {
		reply.BootFileName = s.bootfilename
	}
	if len(s.rootpath) > 0 {
		reply.UpdateOption(dhcpv4.OptRootPath(s.rootpath))
	}
	if err != nil {
		log.Printf("Could not create reply for %v: %v", m, err)
		return
	}
	log.Printf("Sending %v to %v", reply.Summary(), peer)
	if _, err := conn.WriteTo(reply.ToBytes(), peer); err != nil {
		log.Printf("Could not write %v: %v", reply, err)
	}
}

func main() {
	flag.Parse()

	var maca net.HardwareAddr
	if len(*mac) > 0 {
		var err error
		maca, err = net.ParseMAC(*mac)
		if err != nil {
			log.Fatal(err)
		}
	}
	yourIP, yourNet, err := net.ParseCIDR(*yourIP)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	if len(*tftpDir) != 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			server, err := tftp.NewServer(":69")
			if err != nil {
				log.Fatalf("Could not start TFTP server: %v", err)
			}

			log.Println("starting file server")
			server.ReadHandler(tftp.FileServer(*tftpDir))
			log.Fatal(server.ListenAndServe())
		}()
	}
	if len(*httpDir) != 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			http.Handle("/", http.FileServer(http.Dir(*httpDir)))
			log.Fatal(http.ListenAndServe(":80", nil))
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		s := &server{
			mac:          maca,
			self:         net.ParseIP(*selfIP),
			yourIP:       yourIP,
			submask:      yourNet.Mask,
			bootfilename: *bootfilename,
			rootpath:     *rootpath,
		}

		laddr := &net.UDPAddr{Port: dhcpv4.ServerPort}
		server, err := server4.NewServer(*inf, laddr, s.dhcpHandler)
		if err != nil {
			log.Fatal(err)
		}
		log.Fatal(server.Serve())
	}()

	wg.Wait()
}
