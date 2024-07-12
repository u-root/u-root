// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package traceroute

import (
	"fmt"
	"log"
	"net"
	"time"
)

type Probe struct {
	id       uint32
	sendtime time.Time
	recvTime time.Time
	dest     [4]byte
	port     uint16
	ttl      int
	saddr    string
	done     bool
}

func RunTraceroute(host string, prot string, debug bool) error {
	dAddr, err := destAddr(host)
	if err != nil {
		return err
	}

	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		log.Fatal(err)
	}
	sAddr := conn.LocalAddr().(*net.UDPAddr)
	conn.Close()

	cc := coms{
		sendChan: make(chan *Probe),
		recvChan: make(chan *Probe),
		exitChan: make(chan bool),
	}

	mod := NewTrace(prot, dAddr, sAddr, cc, debug)

	switch prot {
	case "udp4":
		go mod.SendTracesUDP4()
	case "tcp4":
		go mod.SendTracesTCP4()
	case "icmp4":
		go mod.SendTracesICMP4()

	}

	printMap := runTransmission(cc)

	destTTL := findDestinationTTL(printMap, dAddr)
	fmt.Printf("traceroute to %s (%s), %d hops max, %d byte packets\n",
		host,
		net.IPv4(dAddr[0], dAddr[1], dAddr[2], dAddr[3]).String(),
		mod.MaxHops,
		60)

	for i := 1; i <= destTTL; i++ {
		pbs := getProbesByTLL(printMap, i)
		if len(pbs) == 0 {
			continue
		}
		fmt.Printf("TTL: %-5d", i)
		for _, pb := range pbs {
			fmt.Printf("%-20s (%-7.3fms) ", pb.saddr, float64(pb.recvTime.Sub(pb.sendtime)/time.Microsecond)/1000)
		}
		fmt.Printf("\n")
	}

	return nil
}

func runTransmission(cc coms) map[int]*Probe {
	sendProbes := make([]*Probe, 0)
	printMap := map[int]*Probe{}
	for {
		var p *Probe
		select {
		case p = <-cc.sendChan:
			sendProbes = append(sendProbes, p)
		case p = <-cc.recvChan:
			for i, sp := range sendProbes {
				if sp.id == p.id {
					sendProbes[i].recvTime = p.recvTime
					sendProbes[i].saddr = p.saddr
					sendProbes[i].done = true
					// Add to map
					printMap[int(sp.id)] = sendProbes[i]

					if p.saddr == net.IPv4(sp.dest[0], sp.dest[1], sp.dest[2], sp.dest[3]).String() {
						//fmt.Println(p.saddr)
						return printMap
					}
				}
			}
		default:
			continue
		}
	}
}
