// Copyright 2019-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Synopsis:
//     ipmidump [-option]
//
// Description:
//
// Options:
//     -chassis : Print chassis power status.
//     -sel     : Print SEL information.
//     -help    : Print help message.
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/u-root/u-root/pkg/ipmi"
)

const cmd = "ipmidump [options] "

var (
	flagChassis = flag.Bool("chassis", false, "print chassis power status")
	flagSEL     = flag.Bool("sel", false, "print SEL information")
	flagHelp    = flag.Bool("help", false, "print help message")
)

func itob(i int) bool { return i != 0 }

func init() {
	defUsage := flag.Usage
	flag.Usage = func() {
		os.Args[0] = cmd
		defUsage()
	}
}

func main() {
	flag.Parse()

	if *flagHelp {
		flag.Usage()
		os.Exit(1)
	}

	if *flagChassis {
		chassisInfo()
	}

	if *flagSEL {
		selInfo()
	}
}

func chassisInfo() {
	allow := map[bool]string{true: "allowed", false: "not allowed"}
	act := map[bool]string{true: "active", false: "inactive"}
	state := map[bool]string{true: "true", false: "false"}

	policy := map[int]string{
		0x0: "always-off",
		0x1: "previous",
		0x2: "always-on",
		0x3: "unknown",
	}

	event := map[int]string{
		0x10: "IPMI command",
		0x08: "power fault",
		0x04: "power interlock",
		0x02: "power overload",
		0x01: "AC failed",
		0x00: "none",
	}

	ipmi, err := ipmi.Open(0)
	if err != nil {
		fmt.Printf("Failed to open ipmi device: %v\n", err)
	}
	defer ipmi.Close()

	if status, err := ipmi.GetChassisStatus(); err != nil {
		fmt.Printf("Failed to get chassis power status: %v\n", err)
	} else {
		// Current power status
		data := int(status.CurrentPowerState)
		fmt.Println("Chassis power status")
		fmt.Println("Power Restore Policy:", policy[(data>>5)&0x03])
		fmt.Println("Power Control Fault :", state[itob(data&0x10)])
		fmt.Println("Power Fault         :", state[itob(data&0x08)])
		fmt.Println("Power Interlock     :", act[itob(data&0x04)])
		fmt.Println("Power Overload      :", state[itob(data&0x02)])
		fmt.Printf("Power Status        : ")
		if (data & 0x01) != 0 {
			fmt.Println("on")
		} else {
			fmt.Println("off")
		}

		// Last power event
		data = int(status.LastPowerEvent)
		fmt.Println("Last Power Event    :", event[data&0x1F])

		// Misc. chassis state
		data = int(status.MiscChassisState)
		fmt.Println("Misc. chassis state")
		fmt.Println("Cooling/Fan Fault   :", state[itob(data&0x08)])
		fmt.Println("Drive Fault         :", state[itob(data&0x04)])
		fmt.Println("Front Panel Lockout :", act[itob(data&0x02)])
		fmt.Println("Chass Intrusion     :", act[itob(data&0x01)])

		// Front panel button (optional)
		data = int(status.FrontPanelButton)
		if status.FrontPanelButton != 0 {
			fmt.Println("Front Panel Button")
			fmt.Println("Standby Button Disable    :", allow[itob(data&0x80)])
			fmt.Println("Diagnostic Buttton Disable:", allow[itob(data&0x40)])
			fmt.Println("Reset Button Disable      :", allow[itob(data&0x20)])
			fmt.Println("Power-off Button Disable  :", allow[itob(data&0x10)])

			fmt.Println("Standby Button            :", state[itob(data&0x08)])
			fmt.Println("Diagnostic Buttton        :", state[itob(data&0x04)])
			fmt.Println("Reset Button              :", state[itob(data&0x02)])
			fmt.Println("Power-off Button          :", state[itob(data&0x01)])
		} else {
			fmt.Println("Front Panel Button  : none")
		}
	}
}

func selInfo() {
	support := map[bool]string{true: "supported", false: "unsupported"}

	ipmi, err := ipmi.Open(0)
	if err != nil {
		fmt.Printf("Failed to open ipmi device: %v\n", err)
	}
	defer ipmi.Close()

	if info, err := ipmi.GetSELInfo(); err != nil {
		fmt.Printf("Failed to get SEL information: %v\n", err)
	} else {
		fmt.Println("SEL information")

		switch info.Version {
		case 0x51:
			fallthrough
		case 0x02:
			fmt.Printf("Version        : %d.%d (1.5, 2.0 compliant)\n", info.Version&0x0F, info.Version>>4)
		default:
			fmt.Println("Version        : unknown")
		}

		fmt.Println("Entries        :", info.Entries)
		fmt.Printf("Free Space     : %d bytes\n", info.FreeSpace)

		// Most recent addition/erase timestamp
		fmt.Printf("Last Add Time  : ")
		if info.LastAddTime != 0xFFFFFFFF {
			fmt.Println(time.Unix(int64(info.LastAddTime), 0))
		} else {
			fmt.Println("not available")
		}

		fmt.Printf("Last Del Time  : ")
		if info.LastDelTime != 0xFFFFFFFF {
			fmt.Println(time.Unix(int64(info.LastDelTime), 0))
		} else {
			fmt.Println("not available")
		}

		// Operation Support
		fmt.Printf("Overflow       : ")
		if (info.OpSupport & 0x80) != 0 {
			fmt.Println("true")
		} else {
			fmt.Println("false")
		}

		data := int(info.OpSupport)
		if (data & 0x0F) != 0 {
			fmt.Println("Supported cmds")
			fmt.Println("Delete         :", support[itob(data&0x08)])
			fmt.Println("Partial Add    :", support[itob(data&0x04)])
			fmt.Println("Reserve        :", support[itob(data&0x02)])
			fmt.Println("Get Alloc Info :", support[itob(data&0x01)])
		} else {
			fmt.Println("Supported cmds : none")
		}
	}
}
