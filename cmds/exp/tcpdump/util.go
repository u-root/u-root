// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"fmt"
	"time"

	"github.com/vishvananda/netlink"
)

//go:generate go run gen.go

// wellKnownPorts returns the well-known name of the port or the port number itself.
func (cmd cmd) wellKnownPorts(port string) string {
	if name, ok := wellKnownPortsMap[port]; ok && !cmd.Opts.Numerical {
		return name
	}

	return port
}

// listDevices lists all the network devices which can be listed to.
func listDevices() error {
	links, err := netlink.LinkList()
	if err != nil {
		return err
	}

	for idx, link := range links {
		fmt.Printf("%d.%s [%s]\n", idx, link.Attrs().Name, link.Attrs().OperState)
	}

	return nil
}

// parseTimeStamp returns the timestamp in the format specified by the user using the -t, -tt, -ttt, -tttt, -ttttt, -nano flags.
func (cmd *cmd) parseTimeStamp(currentTimestamp, lastTimeStamp time.Time) string {
	if cmd.Opts.T {
		return ""
	}
	if cmd.Opts.TT {
		return fmt.Sprintf("%d", currentTimestamp.Unix())
	}
	if cmd.Opts.TTT {
		diff := currentTimestamp.Sub(lastTimeStamp)
		if cmd.Opts.TimeStampInNanoSeconds {
			return fmt.Sprintf("%02d:%02d:%02d.%09d", int(diff.Hours()), int(diff.Minutes())%60, int(diff.Seconds())%60, diff.Nanoseconds()%1e9)
		}

		return fmt.Sprintf("%02d:%02d:%02d.%06d", int(diff.Hours()), int(diff.Minutes())%60, int(diff.Seconds())%60, diff.Microseconds()%1e6)
	}
	if cmd.Opts.TTTT {
		diff := currentTimestamp.Sub(lastTimeStamp)
		return fmt.Sprintf("%02d:%02d:%02d", int(diff.Hours()), int(diff.Minutes())%60, int(diff.Seconds())%60)
	}
	return currentTimestamp.Format("15:04:05.000000")
}

func formatPacketData(data []byte) string {
	var result string
	for i := 0; i < len(data); i += 16 {
		// Print the offset
		result += fmt.Sprintf("0x%04x:  ", i)

		// Print the hex values
		for j := range 16 {
			if i+j < len(data) {
				result += fmt.Sprintf("%02x", data[i+j])
			} else {
				result += "  "
			}
			if j%2 == 1 {
				result += " "
			}
		}
		result += "\n"
	}
	return result
}
