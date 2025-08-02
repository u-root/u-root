// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package traceroute

import (
	"bytes"
	"encoding/binary"
	"log"
	"math/rand"
	"net"
	"time"

	"golang.org/x/net/ipv6"
)

func (t *Trace) SendTracesUDP6() {
	id := uint16(1)
	dport := uint16(int32(t.DestPort) + rand.Int31n(64))
	sport := uint16(1000 + t.PortOffset + rand.Int31n(500))
	mod := uint16(1 << 15)

	for ttl := 1; ttl <= int(t.MaxHops); ttl++ {
		for j := 0; j < t.TracesPerHop; j++ {
			conn, err := net.ListenPacket("ip6:udp", "")
			if err != nil {
				log.Fatalf("net.ListenPacket() = %v", err)
			}
			defer conn.Close()

			rSock := ipv6.NewPacketConn(conn)
			rSock.SetChecksum(true, 0x8)

			pb := &Probe{
				ID:   uint32(id),
				Dest: t.DestIP,
				Port: dport,
				TTL:  ttl,
			}
			cm, payload := t.BuildUDP6Pkt(sport, dport, uint8(ttl), id, 0)

			pb.Sendtime = time.Now()
			rSock.WriteTo(payload, cm, &net.IPAddr{IP: t.DestIP})
			t.SendChan <- pb

			dport = uint16(int32(t.DestPort) + rand.Int31n(64))
			id = (id + 1) % mod
			go t.ReceiveTracesUDP6()
			time.Sleep(time.Microsecond * time.Duration(100000))
		}
	}
}

func (t *Trace) ReceiveTracesUDP6() {
	var err error
	recvICMPConn, err := net.ListenIP("ip6:ipv6-icmp", nil)
	if err != nil {
		log.Fatal("bind failure:", err)
	}

	buf := make([]byte, 1500)
	n, raddr, err := recvICMPConn.ReadFrom(buf)
	if err != nil {
		log.Fatal(err)
	}

	ip6hdr, _ := ipv6.ParseHeader(buf[8:])

	icmpType := buf[0]
	if (icmpType == 1 || (icmpType == 3 && buf[1] == 0)) && (n >= 36) { // TTL Exceeded or Port Unreachable
		id := binary.BigEndian.Uint16(buf[46+ipv6.HeaderLen : 48+ipv6.HeaderLen])
		if ip6hdr.Dst.Equal(t.DestIP) { // && dstPort == t.dstPort {
			recvProbe := &Probe{
				ID:       uint32(id),
				Saddr:    net.ParseIP(raddr.String()),
				RecvTime: time.Now(),
			}
			t.ReceiveChan <- recvProbe
		}
	}
}

func (t *Trace) BuildUDP6Pkt(sport, dport uint16, ttl uint8, id uint16, tos int) (*ipv6.ControlMessage, []byte) {
	cm := &ipv6.ControlMessage{
		HopLimit: int(ttl),
	}

	udphdr := UDPHeader{
		Src: sport,
		Dst: dport,
	}

	payload := make([]byte, 30)
	for i := range 30 {
		payload[i] = uint8(i + 64)
	}

	// Place the ID at the end of the payload.
	idBin := make([]byte, 2)
	binary.BigEndian.PutUint16(idBin, id)
	payload = append(payload, idBin...)

	udphdr.Length = uint16(len(payload) + 8)

	// Set the checksum to 0xffff, which indicates that it should be calculated.
	udphdr.Chksum = 0xffff

	var b bytes.Buffer
	binary.Write(&b, binary.BigEndian, &udphdr)
	binary.Write(&b, binary.BigEndian, &payload)
	return cm, b.Bytes()
}
