// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package traceroute

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"golang.org/x/net/ipv4"
)

func (t *Trace) SendTracesTCP4() {
	sport := uint16(1000 + t.PortOffset + rand.Int31n(500))
	conn, err := net.ListenPacket("ip4:tcp", t.srcIP.String())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	rSocket, err := ipv4.NewRawConn(conn)
	if err != nil {
		log.Fatal("can not create raw socket:", err)
	}
	seq := uint32(1000)
	mod := uint32(1 << 30)
	for ttl := 1; ttl <= int(t.MaxHops); ttl++ {
		for j := 0; j < t.TracesPerHop; j++ {
			hdr, payload := t.BuildTCP4SYNPkt(sport, t.destPort, uint8(ttl), seq, 0)
			rSocket.WriteTo(hdr, payload, nil)
			pb := &Probe{
				id:       seq,
				dest:     t.destIP,
				ttl:      ttl,
				sendtime: time.Now(),
			}
			t.SendChan <- pb
			seq = (seq + 4) % mod
			go t.ReceiveTracesTCP4ICMP()
			go t.ReceiveTracesTCP4()
			time.Sleep(time.Microsecond * time.Duration(200000/t.PacketRate))
		}

	}
}

func (t *Trace) ReceiveTracesTCP4() {
	laddr := &net.IPAddr{IP: t.srcIP}
	recvTCPConn, err := net.ListenIP("ip4:tcp", laddr)
	if err != nil {
		log.Fatal("bind TCP failure:", err)
	}
	buf := make([]byte, 1500)
	n, raddr, err := recvTCPConn.ReadFrom(buf)
	if err != nil {
		return
	}

	if (n >= 20) && (n <= 100) {
		if (buf[13] == TCP_ACK+TCP_SYN) && (raddr.String() == t.destIP.String()) {
			//no need to generate RST message, Linux will automatically send rst
			//sport := binary.BigEndian.Uint16(buf[0:2])
			//dport := binary.BigEndian.Uint16(buf[2:4])
			ack := binary.BigEndian.Uint32(buf[8:12]) - 1
			pb := &Probe{
				id:       ack,
				saddr:    net.ParseIP(raddr.String()),
				recvTime: time.Now(),
			}
			t.ReceiveChan <- pb
		}
	}
}

func (t *Trace) ReceiveTracesTCP4ICMP() {
	//laddr := &net.IPAddr{IP: t.srcIP}
	recvICMPConn, err := net.ListenIP("ip4:icmp", &net.IPAddr{IP: t.srcIP})
	if err != nil {
		log.Fatal("bind failure:", err)
	}
	for {
		buf := make([]byte, 1500)
		n, raddr, err := recvICMPConn.ReadFrom(buf)
		if err != nil {
			break
		}

		icmpType := buf[0]
		if (icmpType == 11 || (icmpType == 3 && buf[1] == 3)) && (n >= 36) { //TTL Exceeded or Port Unreachable
			seq := binary.BigEndian.Uint32(buf[32:36])
			dstip := net.IP(buf[24:28])
			//srcip := net.IP(buf[20:24])
			//srcPort := binary.BigEndian.Uint16(buf[28:30])
			//dstPort := binary.BigEndian.Uint16(buf[30:32])
			if dstip.Equal(t.destIP) { // && dstPort == t.dstPort {
				pb := &Probe{
					id:       seq,
					saddr:    net.ParseIP(raddr.String()),
					recvTime: time.Now(),
				}
				t.ReceiveChan <- pb
			}
		}
	}
}

func (t *Trace) IPv4TCPProbe(dport uint16) {
	seq := uint32(1000)
	mod := uint32(1 << 30)
	for i := 0; i < t.MaxHops; i++ {
		go t.IPv4TCPPing(seq, dport)
		seq = (seq + 4) % mod
		time.Sleep(time.Microsecond * time.Duration(200000/t.PacketRate))
	}

}

func (t *Trace) IPv4TCPPing(seq uint32, dport uint16) {
	pbs := &Probe{
		id:       seq,
		dest:     t.destIP,
		ttl:      0,
		sendtime: time.Now(),
	}
	t.SendChan <- pbs

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", t.destIP.String(), dport), time.Second*2)
	if err != nil {
		return
	}
	conn.Close()
	pbr := &Probe{
		id:       seq,
		saddr:    t.destIP,
		recvTime: time.Now(),
	}
	t.ReceiveChan <- pbr
}

func (t *Trace) BuildTCP4SYNPkt(srcPort uint16, dstPort uint16, ttl uint8, seq uint32, tos int) (*ipv4.Header, []byte) {
	iph := &ipv4.Header{
		Version:  ipv4.Version,
		TOS:      tos,
		Len:      ipv4.HeaderLen,
		TotalLen: 60,
		ID:       0,
		Flags:    0,
		FragOff:  0,
		TTL:      int(ttl),
		Protocol: 6,
		Checksum: 0,
		Src:      t.srcIP,
		Dst:      t.destIP,
	}

	h, err := iph.Marshal()
	if err != nil {
		log.Fatal(err)
	}
	iph.Checksum = int(checkSum(h))

	tcp := TCPHeader{
		Src:        srcPort,
		Dst:        dstPort,
		SeqNum:     seq,
		AckNum:     0,
		DataOffset: 160,
		Flags:      TCP_SYN,
		Window:     64240,
		Urgent:     0,
	}

	//payload is TCP Optionheader
	payload := []byte{0x02, 0x04, 0x05, 0xb4, 0x04, 0x02, 0x08, 0x0a, 0x7f, 0x73, 0xf9, 0x3a, 0x00, 0x00, 0x00, 0x00, 0x01, 0x03, 0x03, 0x07}
	tcp.checksum(iph, payload)

	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, &tcp)
	binary.Write(&buf, binary.BigEndian, &payload)
	return iph, buf.Bytes()
}
