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

	"golang.org/x/net/ipv4"
)

func (t *Trace) SendTracesICMP4() {
	conn, err := net.ListenPacket("ip4:icmp", t.SrcIP.String())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	rSocket, err := ipv4.NewRawConn(conn)
	if err != nil {
		log.Fatal("can not create raw socket:", err)
	}

	id := uint16(1)
	seq := id
	if t.DestPort != 0 {
		seq = t.DestPort
	}
	mod := uint16(1 << 15)

	for ttl := 1; ttl <= int(t.MaxHops); ttl++ {
		for j := 0; j < t.TracesPerHop; j++ {
			hdr, payload := t.BuildICMP4Pkt(uint8(ttl), id, seq, 0)
			rSocket.WriteTo(hdr, payload, nil)
			pb := &Probe{
				ID:       uint32(hdr.ID),
				Dest:     t.DestIP.To4(),
				TTL:      ttl,
				Sendtime: time.Now(),
			}
			t.SendChan <- pb
			id = (id + 1) % mod
			seq = (seq + 1) % mod
			go t.ReceiveTracesICMP4()
			time.Sleep(time.Microsecond * time.Duration(100000/t.PacketRate))
		}
	}
	for _, port := range []uint16{
		80,
		8080,
		443,
		8443,
	} {
		go t.IPv4TCPProbe(port)
	}
}

func (t *Trace) ReceiveTracesICMP4() {
	laddr := &net.IPAddr{IP: t.SrcIP.To4()}
	recvICMPConn, err := net.ListenIP("ip4:icmp", laddr)
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
		id := binary.BigEndian.Uint16(buf[32:34])
		dstip := net.IP(buf[24:28])
		// srcip := net.IP(buf[20:24])

		if dstip.Equal(t.DestIP) {
			pb := &Probe{
				ID:       uint32(id),
				Saddr:    net.ParseIP(raddr.String()),
				RecvTime: time.Now(),
			}
			t.ReceiveChan <- pb
		}
	}
}

func (t *Trace) BuildICMP4Pkt(ttl uint8, id, seq uint16, tos int) (*ipv4.Header, []byte) {
	iph := &ipv4.Header{
		Version:  ipv4.Version,
		TOS:      tos,
		Len:      ipv4.HeaderLen,
		TotalLen: 40,
		ID:       int(id),
		Flags:    0,
		FragOff:  0,
		TTL:      int(ttl),
		Protocol: 1,
		Checksum: 0,
		Src:      t.SrcIP,
		Dst:      t.DestIP,
	}

	h, err := iph.Marshal()
	if err != nil {
		log.Fatal(err)
	}
	iph.Checksum = int(checkSum(h))

	icmp := ICMPHeader{
		IType:    8, // Echo
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
	binary.Write(&b, binary.BigEndian, icmp)
	binary.Write(&b, binary.BigEndian, &payload)
	icmp.Checksum = checkSum(b.Bytes())

	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, &icmp)
	binary.Write(&buf, binary.BigEndian, &payload)
	return iph, buf.Bytes()
}
