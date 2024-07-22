// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package traceroute

import (
	"fmt"
	"net"
	"time"
)

func (t *Trace) IPv6TCPProbe(dport uint16) {
	seq := uint32(1000)
	mod := uint32(1 << 30)
	for i := 0; i < t.MaxHops; i++ {
		go t.IPv6TCPPing(seq, dport)
		seq = (seq + 4) % mod
		time.Sleep(time.Microsecond * time.Duration(200000/t.PacketRate))
	}
}

func (t *Trace) IPv6TCPPing(seq uint32, dport uint16) {
	pbs := &Probe{
		id:       seq,
		dest:     t.destIP,
		ttl:      0,
		sendtime: time.Now(),
	}
	t.SendChan <- pbs

	conn, err := net.DialTimeout("ip6:tcp", fmt.Sprintf("%s:%d", t.destIP.String(), dport), time.Second*2)
	if err != nil {
		return
	}
	conn.Close()

	fmt.Println("tcp probe")
	pbr := &Probe{
		id:       seq,
		saddr:    t.destIP,
		recvTime: time.Now(),
	}
	t.ReceiveChan <- pbr
}
