// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netstat

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"golang.org/x/sys/unix"
)

var ProcNetDevPath = "/proc/net/dev"

func PrintInterfaceTable(ifstr string, cont bool, out io.Writer) error {
	for {
		var s strings.Builder

		fmt.Fprintf(&s, "%s\n", "Kernel Interface table")
		fmt.Fprintf(&s, "%-16s %-8s %-8s %-8s %-8s %-8s %-8s %-8s %-8s %-8s %s\n",
			"Iface",
			"MTU",
			"Rx-OK",
			"Rx-ERR",
			"Rx-DRP",
			"Rx-OVR",
			"TX-OK",
			"TX-ERR",
			"TX-DRP",
			"TX-OVR",
			"Flg",
		)

		ifdata, err := readProcNetDevData()
		if err != nil {
			return fmt.Errorf("read network device: %w", err)
		}

		for _, iface := range ifdata {
			if ifstr == iface.IfName || ifstr == "" {
				fmt.Fprintf(&s, "%-16s %-8d %-8d %-8d %-8d %-8d %-8d %-8d %-8d %-8d %s\n",
					iface.IfName,
					iface.MTU,
					iface.RxPackets,
					iface.RxErrs,
					iface.RxDrops,
					iface.RxFifo,
					iface.TxPackets,
					iface.TxErrs,
					iface.TxDrops,
					iface.TxFifo,
					iface.Flags,
				)
			}
		}

		fmt.Fprintf(out, "%s\n", s.String())
		if !cont {
			break
		}
		time.Sleep(2 * time.Second)
	}

	return nil
}

type ifData struct {
	IfName    string
	RxBytes   uint64
	RxPackets uint64
	RxErrs    uint64
	RxDrops   uint64
	RxFifo    uint64
	RxFrame   uint64
	RxCompr   uint64
	RxMulti   uint64
	TxBytes   uint64
	TxPackets uint64
	TxErrs    uint64
	TxDrops   uint64
	TxFifo    uint64
	TxColls   uint64
	TxCarrier uint64
	TxCompr   uint64
	MTU       uint32
	Flags     string
}

func readProcNetDevData() ([]ifData, error) {
	ret := make([]ifData, 0)
	// Get interface names and some statistics
	procNetDev, err := os.Open(ProcNetDevPath)
	if err != nil {
		return nil, err
	}

	s := bufio.NewScanner(procNetDev)

	// First two lines are garbage
	s.Scan()
	s.Scan()

	skfd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		return nil, err
	}

	for s.Scan() {
		d := ifData{}
		line := s.Text()
		_, err := fmt.Sscanf(line, "%s %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d",
			&d.IfName,
			&d.RxBytes,
			&d.RxPackets,
			&d.RxErrs,
			&d.RxDrops,
			&d.RxFifo,
			&d.RxFrame,
			&d.RxCompr,
			&d.RxMulti,
			&d.TxBytes,
			&d.TxPackets,
			&d.TxErrs,
			&d.TxDrops,
			&d.TxFifo,
			&d.TxColls,
			&d.TxCarrier,
			&d.TxCompr,
		)
		if err != nil {
			return nil, fmt.Errorf("parse line: %w", err)
		}

		// Remove : from Ifname
		d.IfName, _ = strings.CutSuffix(d.IfName, ":")

		ifreq, err := unix.NewIfreq(d.IfName)
		if err != nil {
			return nil, err
		}

		// Request MTU
		if err := unix.IoctlIfreq(skfd, unix.SIOCGIFMTU, ifreq); err != nil {
			return nil, err
		}

		// Read MTU
		d.MTU = ifreq.Uint32()

		// Request Flags
		if err := unix.IoctlIfreq(skfd, unix.SIOCGIFFLAGS, ifreq); err != nil {
			return nil, err
		}

		// Read Flags
		flags := ifreq.Uint16()

		// Parse flags accordingly
		d.Flags = parseIfFlags(flags)

		ret = append(ret, d)
	}

	return ret, nil
}

// parseIfFlags
func parseIfFlags(flags uint16) string {
	var s strings.Builder

	if flags == 0x0 {
		s.WriteString("[NO FLAGS]")
	}
	if (flags & unix.IFF_BROADCAST) > 0 {
		s.WriteString("B")
	}
	if (flags & unix.IFF_DEBUG) > 0 {
		s.WriteString("D")
	}
	if (flags & unix.IFF_DYNAMIC) > 0 {
		s.WriteString("d")
	}
	if (flags & unix.IFF_LOOPBACK) > 0 {
		s.WriteString("L")
	}
	if (flags & unix.IFF_UP) > 0 {
		s.WriteString("U")
	}
	if (flags & unix.IFF_MULTICAST) > 0 {
		s.WriteString("M")
	}
	if (flags & unix.IFF_MASTER) > 0 {
		s.WriteString("m")
	}
	if (flags & unix.IFF_NOTRAILERS) > 0 {
		s.WriteString("N")
	}
	if (flags & unix.IFF_NOARP) > 0 {
		s.WriteString("O")
	}
	if (flags & unix.IFF_POINTOPOINT) > 0 {
		s.WriteString("p")
	}
	if (flags & unix.IFF_PROMISC) > 0 {
		s.WriteString("P")
	}

	if (flags & unix.IFF_RUNNING) > 0 {
		s.WriteString("R")
	}
	if (flags & unix.IFF_SLAVE) > 0 {
		s.WriteString("s")
	}

	return s.String()
}
