// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package traceroute

import (
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"golang.org/x/net/ipv4"
)

type UDP4Trace struct {
	Dest     net.IP
	destPort uint16
	src      net.IP
	//srcPort     uint16
	PortOffset   int32
	MaxHops      int
	SendChan     chan<- *Probe
	ReceiveChan  chan<- *Probe
	exitChan     chan<- bool
	debug        bool
	TracesPerHop int
}

// SendTrace in a routine
func (u *UDP4Trace) SendTraces() {
	id := uint16(1)
	dport := uint16(int32(u.destPort) + rand.Int31n(64))
	sport := uint16(1000 + u.PortOffset + rand.Int31n(500))
	mod := uint16(1 << 15)

	for ttl := 1; ttl <= int(u.MaxHops); ttl++ {
		for j := 0; j < u.TracesPerHop; j++ {
			conn, err := net.ListenPacket("ip4:udp", "")
			if err != nil {
				log.Printf("net.ListenPacket() = %v", err)
				return
			}
			defer conn.Close()

			rSock, err := ipv4.NewRawConn(conn)
			if err != nil {
				log.Printf("ipv4.NewRawConn() = %v", err)
				return
			}

			pb := &Probe{
				id:   id,
				dest: [4]byte(u.Dest.To4()),
				port: dport,
				ttl:  ttl,
			}
			hdr, pl := u.BuildIPv4UDPkt(sport, dport, uint8(ttl), id, 0)

			pb.sendtime = time.Now()
			if err := rSock.WriteTo(hdr, pl, nil); err != nil {
				log.Fatal(err)
			}

			u.SendChan <- pb
			dport = uint16(int32(u.destPort) + rand.Int31n(64))
			id = (id + 1) % mod
			go u.ReceiveTraces()
			time.Sleep(time.Microsecond * time.Duration(100000))
		}

	}
}

func (u *UDP4Trace) ReceiveTraces() {
	dest := u.Dest.To4()
	var err error
	recvICMPConn, err := net.ListenIP("ip4:icmp", nil)
	if err != nil {
		log.Fatal("bind failure:", err)
	}

	buf := make([]byte, 1500)
	n, raddr, err := recvICMPConn.ReadFrom(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	icmpType := buf[0]
	if (icmpType == 11 || (icmpType == 3 && buf[1] == 3)) && (n >= 36) { //TTL Exceeded or Port Unreachable

		id := binary.BigEndian.Uint16(buf[12:14])
		dstip := net.IP(buf[24:28])
		//srcip := net.IP(buf[20:24])
		_ = binary.BigEndian.Uint16(buf[28:30])
		_ = binary.BigEndian.Uint16(buf[30:32])
		if dstip.Equal(dest) { // && dstPort == t.dstPort {
			recvProbe := &Probe{
				id:       id,
				saddr:    raddr.String(),
				recvTime: time.Now(),
			}
			u.ReceiveChan <- recvProbe
		}
	}
}
