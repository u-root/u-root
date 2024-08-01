// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package traceroute_test

import (
	"bytes"
	"encoding/binary"
	"net"
	"testing"

	"github.com/u-root/u-root/pkg/traceroute"
)

func TestUDP4Packet(t *testing.T) {
	tr := traceroute.Trace{
		DestIP: net.IPv4(127, 0, 0, 1),
		SrcIP:  net.IPv4(127, 0, 0, 1),
	}

	_, _ = tr.BuildUDP4Pkt(0, 0, 1, 0, 0)
}

func TestUDP6Packet(t *testing.T) {
	tr := traceroute.Trace{
		DestIP: net.IPv4(127, 0, 0, 1),
		SrcIP:  net.IPv4(127, 0, 0, 1),
	}

	_, _ = tr.BuildUDP6Pkt(0, 0, 1, 0, 0)
}

func TestTCP4Packet(t *testing.T) {
	tr := traceroute.Trace{
		DestIP: net.IPv4(127, 0, 0, 1),
		SrcIP:  net.IPv4(127, 0, 0, 1),
	}

	_, _ = tr.BuildTCP4SYNPkt(0, 0, 1, 0, 0)
}

func TestTCP6Packet(t *testing.T) {
	tr := traceroute.Trace{
		DestIP: net.IPv4(127, 0, 0, 1),
		SrcIP:  net.IPv4(127, 0, 0, 1),
	}

	_, _ = tr.BuildTCP6SYNPkt(0, 0, 1, 0, 0)
}

func TestICMP4Packet(t *testing.T) {
	tr := traceroute.Trace{
		DestIP: net.IPv4(127, 0, 0, 1),
		SrcIP:  net.IPv4(127, 0, 0, 1),
	}

	_, _ = tr.BuildICMP4Pkt(1, 0, 0, 0)
}

func TestICMP6Packet(t *testing.T) {
	tr := traceroute.Trace{
		DestIP: net.IPv4(127, 0, 0, 1),
		SrcIP:  net.IPv4(127, 0, 0, 1),
	}

	_, _ = tr.BuildICMP6Pkt(1, 0, 0, 0)
}

func TestNewTrace(t *testing.T) {
	destIP := net.IPv4(127, 0, 0, 1)
	srcIP := net.IPv4(127, 0, 0, 1)
	flgs := traceroute.Flags{}
	cc := traceroute.Coms{
		SendChan: make(chan *traceroute.Probe),
		RecvChan: make(chan *traceroute.Probe),
	}

	for _, tt := range []struct {
		prot string
		dest net.IP
		src  net.IP
	}{
		{
			prot: "udp4",
			dest: destIP,
			src:  srcIP,
		},
		{
			prot: "tcp4",
			dest: destIP,
			src:  srcIP,
		},
		{
			prot: "icmp4",
			dest: destIP,
			src:  srcIP,
		},
		{
			prot: "udp6",
			dest: destIP,
			src:  srcIP,
		},
		{
			prot: "tcp6",
			dest: destIP,
			src:  srcIP,
		},
		{
			prot: "icmp6",
			dest: destIP,
			src:  srcIP,
		},
	} {
		t.Run(tt.prot, func(t *testing.T) {
			_ = traceroute.NewTrace(tt.prot, tt.dest, tt.src, cc, &flgs)
		})
	}
}

func TestParseTCP(t *testing.T) {
	hdr := traceroute.TCPHeader{
		Src:        0,
		Dst:        0,
		SeqNum:     1000,
		AckNum:     0,
		DataOffset: 160,
		Flags:      traceroute.TCP_SYN,
		Window:     64240,
		Urgent:     0,
	}

	var data bytes.Buffer

	binary.Write(&data, binary.BigEndian, hdr)

	newhdr, err := traceroute.ParseTCP(data.Bytes())
	if err != nil {
		t.Errorf("ParseTCP() = %v, not nil", err)
	}

	if newhdr.Src != hdr.Src {
		t.Errorf("source address not equal")
	}

	if newhdr.Dst != hdr.Dst {
		t.Errorf("destination address not equal")
	}

	if newhdr.SeqNum != hdr.SeqNum {
		t.Errorf("sequence numbers not equal")
	}

	if newhdr.AckNum != hdr.AckNum {
		t.Errorf("acknowledge numbers not equal")
	}

	if newhdr.DataOffset != hdr.DataOffset {
		t.Errorf("data offsets not equal")
	}

	if newhdr.Flags != hdr.Flags {
		t.Errorf("flags not equal")
	}

	if newhdr.Window != hdr.Window {
		t.Errorf("window not equal")
	}

	if newhdr.Urgent != hdr.Urgent {
		t.Errorf("urgent not equal")
	}
}

func TestFindDestinationTTL(t *testing.T) {
	pbMap := map[int]*traceroute.Probe{
		1: {
			TTL: 0,
		},
		3: {
			TTL: 2,
		},
		7: {
			TTL: 3,
		},
		6: {
			TTL: 14,
		},
	}
	ttl := traceroute.DestTTL(pbMap)

	if ttl != 15 {
		t.Errorf("FindDestinationTTL() = %d, not %d", ttl, 15)
	}
}

func TestGetProbeByTTL(t *testing.T) {
	pbMap := map[int]*traceroute.Probe{
		1: {
			TTL: 1,
		},
		3: {
			TTL: 2,
		},
		7: {
			TTL: 3,
		},
		6: {
			TTL: 14,
		},
		8: {
			TTL: 14,
		},
		11: {
			TTL: 14,
		},
		14: {
			TTL: 14,
		},
		19: {
			TTL: 14,
		},
	}

	pbs := traceroute.GetProbesByTLL(pbMap, 14)

	if len(pbs) != 5 {
		t.Errorf("len(GetProbesByTTL()) = %d, not %d", len(pbs), 5)
	}
}
