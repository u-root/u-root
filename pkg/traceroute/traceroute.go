// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package traceroute

import (
	"fmt"
	"net"
	"time"
)

type Probe struct {
	id       uint32
	sendtime time.Time
	recvTime time.Time
	dest     net.IP
	port     uint16
	ttl      int
	saddr    net.IP
	done     bool
}

func RunTraceroute(host string, prot string, f *Flags) error {
	dAddr, err := destAddr(host, prot)
	if err != nil {
		return err
	}

	fmt.Println(dAddr)

	sAddr, err := srcAddr(prot)
	if err != nil {
		return err
	}

	cc := coms{
		sendChan: make(chan *Probe),
		recvChan: make(chan *Probe),
		exitChan: make(chan bool),
	}

	mod := NewTrace(prot, dAddr, sAddr, cc, f)

	switch prot {
	case "udp4":
		go mod.SendTracesUDP4()
	case "tcp4":
		go mod.SendTracesTCP4()
	case "icmp4":
		go mod.SendTracesICMP4()
	case "udp6":
	case "tcp6":
		go mod.SendTracesTCP6()
	case "icmp6":
		go mod.SendTracesICMP6()
	}

	printMap := runTransmission(cc)

	destTTL := findDestinationTTL(printMap)
	fmt.Printf("traceroute to %s (%s), %d hops max, %d byte packets\n",
		host,
		dAddr.String(),
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
			//fmt.Println(p.id)
		case p = <-cc.recvChan:
			//fmt.Println(p.id)
			for i, sp := range sendProbes {
				if sp.id == p.id {
					sendProbes[i].recvTime = p.recvTime
					sendProbes[i].saddr = p.saddr
					sendProbes[i].done = true
					// Add to map
					printMap[int(sp.id)] = sendProbes[i]
					if p.saddr.Equal(sp.dest) {
						fmt.Println("final")
						return printMap
					}
				}
			}
		default:
			continue
		}
	}
}
