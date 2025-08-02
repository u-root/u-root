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
	ID       uint32
	Sendtime time.Time
	RecvTime time.Time
	Dest     net.IP
	Port     uint16
	TTL      int
	Saddr    net.IP
	Done     bool
}

func RunTraceroute(f *Flags) error {
	dAddr, err := DestAddr(f.Host, f.Proto)
	if err != nil {
		return err
	}

	sAddr, err := SrcAddr(f.Proto)
	if err != nil {
		return err
	}

	cc := Coms{
		SendChan: make(chan *Probe),
		RecvChan: make(chan *Probe),
	}

	mod := NewTrace(f.Proto, dAddr, *sAddr, cc, f)

	switch f.Proto {
	case "udp4":
		go mod.SendTracesUDP4()
	case "tcp4":
		go mod.SendTracesTCP4()
	case "icmp4":
		go mod.SendTracesICMP4()
	case "udp6":
		go mod.SendTracesUDP6()
	case "tcp6":
		go mod.SendTracesTCP6()
	case "icmp6":
		go mod.SendTracesICMP6()
	default:
		return fmt.Errorf("unsupported protocol: %s", f.Proto)
	}

	printMap := runTransmission(cc)

	destTTL := DestTTL(printMap)
	fmt.Printf("traceroute to %s (%s), %d hops max, %d byte packets\n",
		f.Host,
		dAddr.String(),
		mod.MaxHops,
		60)

	for i := 1; i <= destTTL; i++ {
		pbs := GetProbesByTLL(printMap, i)
		if len(pbs) == 0 {
			continue
		}
		fmt.Printf("TTL: %-5d", i)
		for _, pb := range pbs {
			fmt.Printf("%-20s (%-7.3fms) ", pb.Saddr, float64(pb.RecvTime.Sub(pb.Sendtime)/time.Microsecond)/1000)
		}
		fmt.Printf("\n")
	}

	return nil
}

func runTransmission(cc Coms) map[int]*Probe {
	sendProbes := make([]*Probe, 0)
	printMap := map[int]*Probe{}
	for {
		var p *Probe
		select {
		case p = <-cc.SendChan:
			sendProbes = append(sendProbes, p)
		case p = <-cc.RecvChan:
			for i, sp := range sendProbes {
				if sp.ID == p.ID {
					sendProbes[i].RecvTime = p.RecvTime
					sendProbes[i].Saddr = p.Saddr
					sendProbes[i].Done = true
					// Add to map
					printMap[int(sp.ID)] = sendProbes[i]
					if p.Saddr.Equal(sp.Dest) {
						return printMap
					}
				}
			}
		default:
			continue
		}
	}
}
