// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package traceroute

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"time"

	"golang.org/x/net/ipv6"
)

func (t *Trace) SendTracesICMP6() {
	conn, err := net.ListenPacket("ip6:ipv6-icmp", t.SrcIP.String())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	pktconn := ipv6.NewPacketConn(conn)
	id := uint16(1)
	seq := id
	if t.DestPort != 0 {
		seq = t.DestPort
	}
	mod := uint16(1 << 15)

	for ttl := 1; ttl < int(t.MaxHops); ttl++ {
		for j := 0; j < t.TracesPerHop; j++ {
			cm, payload := t.BuildICMP6Pkt(ttl, id, seq, 0)
			pktconn.WriteTo(payload, cm, &net.UDPAddr{IP: t.DestIP})
			pb := &Probe{
				ID:       uint32(id),
				Dest:     t.DestIP,
				TTL:      ttl,
				Sendtime: time.Now(),
			}
			t.SendChan <- pb
			id = (id + 1) % mod
			seq = (seq + 1) % mod
			go t.ReceiveTraceICMP6()
			time.Sleep(time.Microsecond * time.Duration(100000/t.PacketRate))
		}
	}
	for _, port := range []uint16{
		80,
		8080,
		443,
		8443,
	} {
		go t.IPv6TCPProbe(port)
	}
}

func (t *Trace) ReceiveTraceICMP6() {
	recvICMPConn, err := net.ListenIP("ip6:ipv6-icmp", nil)
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 1500)
	n, raddr, err := recvICMPConn.ReadFrom(buf)
	if err != nil {
		log.Fatal(err)
	}

	icmpType := buf[0]
	if (icmpType == 1 || (icmpType == 3 && buf[1] == 0)) && (n >= 36) { // destination unreachable or ttl exceed in transit

		id := binary.BigEndian.Uint16(buf[14+ipv6.HeaderLen : 16+ipv6.HeaderLen])
		ipv6hdr, _ := ipv6.ParseHeader(buf[8:])
		if ipv6hdr.Dst.Equal(t.DestIP) {
			pb := &Probe{
				ID:       uint32(id),
				Saddr:    net.ParseIP(raddr.String()),
				RecvTime: time.Now(),
			}
			t.ReceiveChan <- pb
		}

	}
}

func (t *Trace) BuildICMP6Pkt(ttl int, id uint16, seq uint16, tc int) (*ipv6.ControlMessage, []byte) {
	ctlmsg := &ipv6.ControlMessage{
		TrafficClass: 0,
		HopLimit:     int(ttl),
	}

	icmppkt := ICMPHeader{
		IType:    128,
		ICode:    0,
		Checksum: 0,
		ID:       id,
		Seq:      seq,
	}

	payload := make([]byte, 32)
	for i := range 32 {
		payload[i] = uint8(i + 64)
	}

	var b bytes.Buffer
	binary.Write(&b, binary.BigEndian, icmppkt)
	binary.Write(&b, binary.BigEndian, &payload)
	return ctlmsg, b.Bytes()
}
