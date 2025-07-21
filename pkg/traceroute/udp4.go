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

	"golang.org/x/net/ipv4"
)

// SendTrace in a routine
func (t *Trace) SendTracesUDP4() {
	id := uint16(1)
	dport := uint16(int32(t.DestPort) + rand.Int31n(64))
	sport := uint16(1000 + t.PortOffset + rand.Int31n(500))
	mod := uint16(1 << 15)

	for ttl := 1; ttl <= int(t.MaxHops); ttl++ {
		for j := 0; j < t.TracesPerHop; j++ {
			conn, err := net.ListenPacket("ip4:udp", "")
			if err != nil {
				log.Fatalf("net.ListenPacket() = %v", err)
			}
			defer conn.Close()

			rSock, err := ipv4.NewRawConn(conn)
			if err != nil {
				log.Fatalf("ipv4.NewRawConn() = %v", err)
			}

			pb := &Probe{
				ID:   uint32(id),
				Dest: t.DestIP,
				Port: dport,
				TTL:  ttl,
			}
			hdr, pl := t.BuildUDP4Pkt(sport, dport, uint8(ttl), id, 0)

			pb.Sendtime = time.Now()
			if err := rSock.WriteTo(hdr, pl, nil); err != nil {
				log.Fatal(err)
			}

			t.SendChan <- pb

			dport = uint16(int32(t.DestPort) + rand.Int31n(64))
			id = (id + 1) % mod
			go t.ReceiveTracesUDP4()
			time.Sleep(time.Microsecond * time.Duration(100000))
		}
	}
}

func (t *Trace) ReceiveTracesUDP4() {
	dest := t.DestIP.To4()
	var err error
	recvICMPConn, err := net.ListenIP("ip4:icmp", nil)
	if err != nil {
		log.Fatal("bind failure:", err)
	}

	buf := make([]byte, 1500)
	n, raddr, err := recvICMPConn.ReadFrom(buf)
	if err != nil {
		log.Fatal(err)
	}

	icmpType := buf[0]
	if (icmpType == 11 || (icmpType == 3 && buf[1] == 3)) && (n >= 36) { // TTL Exceeded or Port Unreachable
		id := binary.BigEndian.Uint16(buf[12:14])
		dstip := net.IP(buf[24:28])
		// srcip := net.IP(buf[20:24])
		_ = binary.BigEndian.Uint16(buf[28:30])
		_ = binary.BigEndian.Uint16(buf[30:32])
		if dstip.Equal(dest) { // && dstPort == t.dstPort {
			recvProbe := &Probe{
				ID:       uint32(id),
				Saddr:    net.ParseIP(raddr.String()),
				RecvTime: time.Now(),
			}
			t.ReceiveChan <- recvProbe
		}
	}
}

func (t *Trace) BuildUDP4Pkt(srcPort uint16, dstPort uint16, ttl uint8, id uint16, tos int) (*ipv4.Header, []byte) {
	iph := &ipv4.Header{
		Version:  ipv4.Version,
		TOS:      tos,
		Len:      ipv4.HeaderLen,
		TotalLen: 60,
		ID:       int(id),
		Flags:    0,
		FragOff:  0,
		TTL:      int(ttl),
		Protocol: 17,
		Checksum: 0,
		Src:      t.SrcIP.To4(),
		Dst:      t.DestIP.To4(),
	}

	h, err := iph.Marshal()
	if err != nil {
		log.Fatal(err)
	}
	iph.Checksum = int(checkSum(h))

	udp := UDPHeader{
		Src: srcPort,
		Dst: dstPort,
	}

	payload := make([]byte, 32)
	for i := range 32 {
		payload[i] = uint8(i + 64)
	}
	udp.Length = uint16(len(payload) + 8)
	udp.checksum(iph, payload)

	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, &udp)
	binary.Write(&buf, binary.BigEndian, &payload)
	return iph, buf.Bytes()
}
