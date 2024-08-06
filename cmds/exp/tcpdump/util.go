// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"time"

	"github.com/vishvananda/netlink"
)

var wellKnownPortsMap = map[string]string{
	"20":  "ftp-data",
	"21":  "ftp",
	"22":  "ssh-scp",
	"23":  "telnet",
	"25":  "smtp",
	"53":  "domain",
	"80":  "http",
	"88":  "kerberos",
	"110": "pop3",
	"119": "nntp",
	"123": "ntp",
	"143": "imap",
	"443": "https",
	"465": "smtps",
	"563": "nntps",
	"989": "ftps-data",
	"990": "ftps",
	"993": "imaps",
	"995": "pop3s",
}

// wellKnownPorts returns the well-known name of the port or the port number itself.
func (cmd cmd) wellKnownPorts(port string) string {
	if name, ok := wellKnownPortsMap[port]; ok && !cmd.Opts.numerical {
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
func (cmd *cmd) parseTimeStamp(currentTimestamp, lastTimeStamp time.Time) (timeStamp string) {
	switch {
	case cmd.Opts.t:
		return ""
	case cmd.Opts.tt:
		return fmt.Sprintf("%d", currentTimestamp.Unix())
	case cmd.Opts.ttt, cmd.Opts.ttttt:
		switch cmd.Opts.timeStampInNanoSeconds {
		case true:
			if !cmd.Opts.firstPacketProcessed {
				cmd.Opts.firstPacketProcessed = true
				return "00:00:00.000000000"
			}
			return time.Unix(0, 0).Add(currentTimestamp.Sub(lastTimeStamp)).Format("15:04:05.000000000")
		default:
			if !cmd.Opts.firstPacketProcessed {
				cmd.Opts.firstPacketProcessed = true
				return "00:00:00.000000"
			}
			return time.Unix(0, 0).Add(currentTimestamp.Sub(lastTimeStamp)).Format("15:04:05.000000")
		}
	case cmd.Opts.tttt:
		midnight := time.Now().Truncate(24 * time.Hour)
		return currentTimestamp.Sub(midnight).String()
	}

	return currentTimestamp.Format("15:04:05.000000")
}
