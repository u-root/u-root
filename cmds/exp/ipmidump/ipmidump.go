// Copyright 2019-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Synopsis:
//
//	ipmidump [-option]
//
// Description:
//
// Options:
//
//	-chassis : Print chassis power status.
//	-sel     : Print SEL information.
//	-lan     : Print IP information.
//	-device  : Print device information.
//	-raw     : Send raw command and print response.
//	-help    : Print help message.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/u-root/u-root/pkg/ipmi"
)

const cmd = "ipmidump [options] "

var (
	flagChassis = flag.Bool("chassis", false, "print chassis power status")
	flagSEL     = flag.Bool("sel", false, "print SEL information")
	flagLan     = flag.Bool("lan", false, "Print IP address")
	flagRaw     = flag.Bool("raw", false, "Send IPMI raw command")
	flagHelp    = flag.Bool("help", false, "print help message")
	flagDev     = flag.Bool("device", false, "print device information")
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

	if *flagLan {
		lanConfig()
	}

	if *flagDev {
		deviceID()
	}

	if *flagRaw {
		sendRawCmd(flag.Args())
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

func lanConfig() {
	const (
		setInProgress byte = iota
		_
		_
		IPAddress
		IPAddressSrc
		MACAddress
	)

	setInProgressStr := []string{
		"Set Complete", "Set In Progress", "Commit Write", "Reserved",
	}

	IPAddressSrcStr := []string{
		"Unspecified", "Static Address", "DHCP Address", "BIOS Assigned Address",
	}

	ipmi, err := ipmi.Open(0)
	if err != nil {
		log.Fatal(err)
	}
	defer ipmi.Close()

	// data 1	completion code
	// data 2	parameter revision, 0x11
	// data 3:N	data

	// set in progress
	if buf, err := ipmi.GetLanConfig(1, setInProgress); err != nil {
		fmt.Printf("Failed to get LAN config: %v\n", err)
	} else {
		fmt.Printf("Set In Progress   : ")
		if int(buf[2]) < len(setInProgressStr) {
			fmt.Println(setInProgressStr[buf[2]])
		} else {
			fmt.Println("Unknown")
			fmt.Printf("%v\n", buf)
		}
	}

	// ip address source
	if buf, err := ipmi.GetLanConfig(1, IPAddressSrc); err != nil {
		fmt.Printf("Failed to get LAN config: %v\n", err)
	} else {
		fmt.Printf("IP Address Source : ")
		if int(buf[2]) < len(IPAddressSrcStr) {
			fmt.Println(IPAddressSrcStr[buf[2]])
		} else {
			fmt.Println("Other")
			fmt.Printf("%v\n", buf)
		}
	}

	// ip address
	if buf, err := ipmi.GetLanConfig(1, IPAddress); err != nil {
		fmt.Printf("Failed to get LAN config: %v\n", err)
	} else {
		fmt.Printf("IP Address        : ")
		if len(buf) == 6 {
			fmt.Printf("%d.%d.%d.%d\n", buf[2], buf[3], buf[4], buf[5])
		} else {
			fmt.Printf("Unknown\n")
		}
	}

	// MAC address
	if buf, err := ipmi.GetLanConfig(1, MACAddress); err != nil {
		fmt.Printf("Failed to get LAN config: %v\n", err)
	} else {
		fmt.Printf("MAC Address       : ")
		if len(buf) == 8 {
			fmt.Printf("%02x:%02x:%02x:%02x:%02x:%02x\n", buf[2], buf[3], buf[4], buf[5], buf[6], buf[7])
		} else {
			fmt.Printf("Unknown\n")
		}
	}
}

func deviceID() {
	status := map[byte]string{
		0x80: "yes",
		0x00: "no",
	}

	adtlDevSupport := []string{
		"Sensor Device",         /* bit 0 */
		"SDR Repository Device", /* bit 1 */
		"SEL Device",            /* bit 2 */
		"FRU Inventory Device",  /* bit 3 */
		"IPMB Event Receiver",   /* bit 4 */
		"IPMB Event Generator",  /* bit 5 */
		"Bridge",                /* bit	6 */
		"Chassis Device",        /* bit 7 */
	}

	ipmi, err := ipmi.Open(0)
	if err != nil {
		fmt.Printf("Failed to open ipmi device: %v\n", err)
	}
	defer ipmi.Close()

	if info, err := ipmi.GetDeviceID(); err != nil {
		fmt.Printf("Failed to get device ID information: %v\n", err)
	} else {
		fmt.Println("Device ID information")
		fmt.Printf("%-26s: %d\n", "Device ID", info.DeviceID)
		fmt.Printf("%-26s: %d\n", "Device Revision", (info.DeviceRevision & 0x0F))
		fmt.Printf("%-26s: %d.%02x\n", "Firmware Revision",
			(info.FwRev1 & 0x3F), info.FwRev2)

		spec := uint8(info.IpmiVersion)
		fmt.Printf("%-26s: %x.%x\n", "IPMI Version", spec&0x0F, (spec&0xF0)>>4)

		var mid uint32
		mid = uint32(info.ManufacturerID[2]) << 16
		mid |= uint32(info.ManufacturerID[1]) << 8
		mid |= uint32(info.ManufacturerID[0])

		fmt.Printf("%-26s: %d (0x%04X)\n", "Manufacturer ID", mid, mid)

		var pid uint16
		pid = uint16(info.ProductID[1]) << 8
		pid |= uint16(info.ProductID[0])

		fmt.Printf("%-26s: %d (0x%04X)\n", "Product ID", pid, pid)

		// bit 7 == 0 indicates normal operation
		fmt.Printf("%-26s: %s\n", "Device Available", status[(^info.FwRev1&0x80)])
		fmt.Printf("%-26s: %s\n", "Provides Device SDRs",
			status[(info.DeviceRevision&0x80)])

		fmt.Printf("%-26s:\n", "Additional Device Support")
		for i := 0; i < 8; i++ {
			if (info.AdtlDeviceSupport & (1 << i)) != 0 {
				fmt.Printf("    %s\n", adtlDevSupport[i])
			}
		}

		// This field is optional, 4 bytes.
		fmt.Printf("%-26s:\n", "Aux Firmware Rev Info")
		for _, val := range info.AuxFwRev {
			fmt.Printf("    0x%02x\n", val)
		}
	}
}

func sendRawCmd(cmds []string) {
	ipmi, err := ipmi.Open(0)
	if err != nil {
		log.Fatal(err)
	}
	defer ipmi.Close()

	data := make([]byte, 0)

	for _, cmd := range cmds {
		val, err := strconv.ParseInt(cmd, 0, 8)
		if err != nil {
			fmt.Printf("Invalid syntax: \"%s\"\n", cmd)
			return
		}
		data = append(data, byte(val))
	}

	if buf, err := ipmi.RawCmd(data); err != nil {
		fmt.Printf("Unable to send RAW command: %v\n", err)
	} else {
		for i, x := range buf {
			fmt.Printf("| 0x%-2x ", x)
			if i%8 == 7 || i == len(buf)-1 {
				fmt.Printf("|\n")
			}
		}
	}
}
