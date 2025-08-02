// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package brctl provides a Go interface to the Linux bridge control
// (brctl) command. It allows you to manage Ethernet bridges and their
// interfaces, including adding and deleting bridges, adding and
// deleting interfaces, and configuring various bridge parameters.
//
// The original C implementation offers the ability to issue the bridge
// configuration in 2 ways: 1. `ioctl` and 2. `sysfs`. Since all modern
// systems deploy the `sysfs` we use it to configure the bridges whenever
// possible. The create and deletion of bridges and their interfaces is
// achieved via `ioctl`.
package brctl

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
)

// Addbr adds a bridge with the provided name.
func Addbr(name string) error {
	if len(name) >= unix.IFNAMSIZ {
		return fmt.Errorf("bridge name too long, %d bytes allowed", unix.IFNAMSIZ-1)
	}

	brctlSocket, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		return fmt.Errorf("open unix socket: %w", err)
	}

	if _, err := executeIoctlStr(brctlSocket, unix.SIOCBRADDBR, name); err != nil {
		return fmt.Errorf("can't add bridge %q: %w", name, err)
	}

	return nil
}

// Delbr deletes a bridge with the provided name.
func Delbr(name string) error {
	brctlSocket, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		return fmt.Errorf("open unix socket: %w", err)
	}

	if _, err := executeIoctlStr(brctlSocket, unix.SIOCBRDELBR, name); err != nil {
		return fmt.Errorf("can't delete bridge %q: %w", name, err)
	}

	return nil
}

// Addif adds an interface to bridge.
func Addif(bridge string, iface string) error {
	brctlSocket, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		return fmt.Errorf("open unix socket: %w", err)
	}

	ifr, err := unix.NewIfreq(bridge)
	if err != nil {
		return fmt.Errorf("bridge name exceeds max length of %d", unix.IFNAMSIZ-1)
	}

	ifIndex, err := getIndexFromInterfaceName(iface)
	if err != nil {
		return fmt.Errorf("interface %q: %w", iface, err)
	}
	ifr.SetUint32(uint32(ifIndex))

	if err := unix.IoctlIfreq(brctlSocket, unix.SIOCBRADDIF, ifr); err != nil {
		if errors.Is(err, syscall.ENODEV) { // no such device
			return fmt.Errorf("bridge %q: %w", bridge, err)
		}
		if errors.Is(err, syscall.EBUSY) { // resource busy
			return fmt.Errorf("device %q is already a member of a bridge; can't add it to bridge %q", iface, bridge)
		}
		return fmt.Errorf("can't add %q to bridge %q: %w", iface, bridge, err)
	}

	return nil
}

// Delif deletes an interface from bridge.
func Delif(bridge string, iface string) error {
	brctlSocket, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		return fmt.Errorf("open unix socket: %w", err)
	}

	ifr, err := unix.NewIfreq(bridge)
	if err != nil {
		return fmt.Errorf("bridge name exceeds max length of %d", unix.IFNAMSIZ-1)
	}

	ifIndex, err := getIndexFromInterfaceName(iface)
	if err != nil {
		return fmt.Errorf("interface %q: %w", iface, err)
	}
	ifr.SetUint32(uint32(ifIndex))

	if err := unix.IoctlIfreq(brctlSocket, unix.SIOCBRDELIF, ifr); err != nil {
		if errors.Is(err, syscall.ENODEV) { // no such device
			return fmt.Errorf("bridge %q: %w", bridge, err)
		}
		if errors.Is(err, syscall.EINVAL) { // invalid argument
			return fmt.Errorf("device %q is not a member of bridge %q", iface, bridge)
		}
		return fmt.Errorf("can't delete %q to bridge %q: %w", iface, bridge, err)
	}

	return nil
}

// print a line to out with essential bridge information
func showBridge(name string, out io.Writer) error {
	info, err := NewInfo(name)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("bridge %q does not exist", name)
		}
		return fmt.Errorf("read info for bridge %q: %w", name, err)
	}

	ifStr := ""
	for _, s := range info.Interfaces {
		ifStr += s.Name + " "
	}

	fmt.Fprintf(out, ShowBridgeFmt, info.Name, info.BridgeID, info.StpEnabled, ifStr)

	return nil
}

// ShowMACs shows a list of learned MAC addresses for this bridge.
// The following byte format applies according to the kernel source [1]
// 0-5:   MAC address
// 6:     port number
// 7:     is_local
// 8-11:  ageing timer
// 12-15: unused in this context
//
// [1] https://github.com/torvalds/linux/blob/master/include/uapi/linux/if_bridge.h#L93
func ShowMACs(bridge string, out io.Writer) error {
	b2str := func(b bool) string {
		if b {
			return "yes"
		}
		return "no"
	}

	brforward, err := os.ReadFile(filepath.Join(BRCTL_SYS_NET, bridge, BRCTL_BRFORWARD))
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("read forward table: %w", ErrBridgeNotExist)
	} else if err != nil {
		return fmt.Errorf("read forward table: %w", err)
	}

	fmt.Fprintf(out, "port no\tmac addr\t\tis_local?\tageing timer\n")

	// parse sysf in 16 byte chunks
	for i := 0; i < len(brforward); i += 0x10 {
		chunk := brforward[i : i+0x10]
		mac := chunk[0:6]
		portNo := chunk[6]
		isLocal := chunk[7] != 0
		ageingTimer := uint16(binary.BigEndian.Uint16(chunk[8:12]))

		fmt.Fprintf(out, "%3d\t%02x:%02x:%02x:%02x:%02x:%02x\t%s\t\t%.2f\n", portNo, mac[0], mac[1], mac[2], mac[3], mac[4], mac[5], b2str(isLocal), float64(ageingTimer)/100)
	}

	return nil
}

const ShowBridgeFmt = "%-15s %23s %15v %20v\n"

// Show performs the brctl show command.
// If no names are provided, it will show all bridges.
// If names are provided, it will show the specified bridges.
func Show(out io.Writer, names ...string) error {
	fmt.Fprintf(out, ShowBridgeFmt, "bridge name", "bridge id", "STP enabled", "interfaces")

	if len(names) == 0 {
		devices, err := os.ReadDir(BRCTL_SYS_NET)
		if err != nil {
			return fmt.Errorf("listing devices (%q): %w", BRCTL_SYS_NET, err)
		}

		for _, dev := range devices {
			// check if device is bridge, thus if it has a bridge directory
			if _, err := os.Stat(filepath.Join(BRCTL_SYS_NET, dev.Name(), BRCTL_BRIDGE_DIR)); err == nil {
				if err = showBridge(dev.Name(), out); err != nil {
					return fmt.Errorf("show bridge %q: %w", dev.Name(), err)
				}
			}
		}
	} else {
		for _, name := range names {
			if err := showBridge(name, out); err != nil {
				return fmt.Errorf("show bridge %q: %w", name, err)
			}
		}
	}
	return nil
}

// SetAgeingTime sets the ethernet (MAC) address ageing time, in seconds.
// After <time> seconds of not having seen a frame coming from a certain address,
// the bridge will time out (delete) that address from the Forwarding DataBase (fdb).
func SetAgeingTime(name string, time string) error {
	ageingTime, err := stringToJiffies(time)
	if err != nil {
		return fmt.Errorf("convert time (%q): %w", time, err)
	}

	if err = setBridgeValue(name, BRCTL_AGEING_TIME, []byte(strconv.Itoa(ageingTime)), uint64(BRCTL_SET_AEGING_TIME)); err != nil {
		return fmt.Errorf("set ageing time: %w", err)
	}
	return nil
}

// SetSTP set the STP state of the bridge to on or off
// Enable using "on" or "yes", disable by providing anything else
// The manpage states:
// > If <state> is "on" or "yes"  the STP  will  be turned on, otherwise it will be turned off
// So this is actually the described behavior, not checking for "off" and "no"
func SetSTP(bridge string, state string) error {
	var stpState int
	if state == "on" || state == "yes" {
		stpState = 1
	} else {
		stpState = 0
	}

	if err := setBridgeValue(bridge, BRCTL_STP_STATE, []byte(strconv.Itoa(stpState)), uint64(BRCTL_SET_BRIDGE_PRIORITY)); err != nil {
		return fmt.Errorf("set STP: %w", err)
	}

	return nil
}

// SetBridgePrio sets the port <port>'s priority to <priority>.
// The priority value is an unsigned 8-bit quantity (a number between 0 and 255),
// and has no dimension. This metric is used in the designated port and root port selection algorithms.
func SetBridgePrio(bridge string, bridgePriority string) error {
	prio, err := strconv.Atoi(bridgePriority)
	if err != nil {
		return err
	}

	if err := setBridgeValue(bridge, BRCTL_BRIDGE_PRIO, []byte(strconv.Itoa(prio)), 0); err != nil {
		return fmt.Errorf("set bridge prio: %w", err)
	}

	return nil
}

// SetForwardDelay sets the bridge's 'bridge forward delay' to <time> seconds.
func SetForwardDelay(bridge string, time string) error {
	forwardDelay, err := stringToJiffies(time)
	if err != nil {
		return fmt.Errorf("convert time (%q): %w", time, err)
	}

	if err := setBridgeValue(bridge, BRCTL_FORWARD_DELAY, []byte(strconv.Itoa(forwardDelay)), 0); err != nil {
		return fmt.Errorf("set forward delay: %w", err)
	}

	return nil
}

// SetHello sets the bridge's 'bridge hello time' to <time> seconds.
func SetHello(bridge string, time string) error {
	helloTime, err := stringToJiffies(time)
	if err != nil {
		return fmt.Errorf("convert time (%q): %w", time, err)
	}

	if err := setBridgeValue(bridge, BRCTL_HELLO_TIME, []byte(strconv.Itoa(helloTime)), 0); err != nil {
		return fmt.Errorf("set hello time: %w", err)
	}

	return nil
}

// SetMaxAge sets the bridge's 'maximum message age' to <time> seconds.
func SetMaxAge(bridge string, time string) error {
	maxAge, err := stringToJiffies(time)
	if err != nil {
		return fmt.Errorf("convert time (%q): %w", time, err)
	}

	if err := setBridgeValue(bridge, BRCTL_MAX_AGE, []byte(strconv.Itoa(maxAge)), 0); err != nil {
		return fmt.Errorf("set max age: %w", err)
	}

	return nil
}

var errBadValue = errors.New("bad value")

// Setpathcost sets the port cost of the port <port> to <cost>. This is a dimensionless metric.
func SetPathCost(bridge string, port string, cost string) error {
	pathCost, err := strconv.ParseUint(cost, 10, 64)
	if err != nil {
		return fmt.Errorf("set path cost: %w", errBadValue)
	}

	err = setPortValue(port, BRCTL_PATH_COST, append([]byte(strconv.FormatUint(pathCost, 10)), BRCTL_SYS_SUFFIX))
	if err != nil {
		return fmt.Errorf("set path cost: %w", err)
	}

	return nil
}

// SetPortPrio sets the port <port>'s priority to <priority>.
// The priority value is an unsigned 8-bit quantity (a number between 0 and 255),
// and has no dimension. This metric is used in the designated port and root port selection algorithms.
func SetPortPrio(bridge string, port string, prio string) error {
	portPriority, err := strconv.Atoi(prio)
	if err != nil {
		return fmt.Errorf("set port prio: %w", errBadValue)
	}

	err = setPortValue(port, BRCTL_PRIORITY, []byte(strconv.Itoa(portPriority)))
	if err != nil {
		return fmt.Errorf("set port prio: %w", err)
	}

	return nil
}

// Hairpin sets the hairpin mode of the <port> attached to <bridge>
func Hairpin(bridge string, port string, hairpinmode string) error {
	var hairpinMode string
	if hairpinmode == "on" {
		hairpinMode = "1"
	} else {
		hairpinMode = "0"
	}

	err := setPortValue(port, BRCTL_HAIRPIN, []byte(hairpinMode))
	if err != nil {
		return fmt.Errorf("set hairpin mode: %w", err)
	}

	return nil
}

func ShowStp(out io.Writer, bridge string) error {
	bridgeInfo, err := NewInfo(bridge)
	if err != nil {
		return err
	}

	var s strings.Builder

	fmt.Fprintf(&s, "%s\n", bridge)
	fmt.Fprintf(&s, " bridge id\t\t%s\n", bridgeInfo.BridgeID)
	fmt.Fprintf(&s, " designated root\t%s\n", bridgeInfo.RootID)
	fmt.Fprintf(&s, " root port\t\t   %d\t\t\t", bridgeInfo.RootPort)
	fmt.Fprintf(&s, "path cost\t\t   %d\n", bridgeInfo.RootPathCost)
	fmt.Fprintf(&s, " max age\t\t%s", timerToString(bridgeInfo.MaxAge))
	fmt.Fprintf(&s, "\t\t\tbridge max age\t\t%s\n", timerToString(bridgeInfo.BridgeMaxAge))
	fmt.Fprintf(&s, " hello time\t\t%s", timerToString(bridgeInfo.HelloTime))
	fmt.Fprintf(&s, "\t\t\tbridge hello time\t%s\n", timerToString(bridgeInfo.BridgeHelloTime))
	fmt.Fprintf(&s, " forward delay\t\t%s", timerToString(bridgeInfo.ForwardDelay))
	fmt.Fprintf(&s, "\t\t\tbridge forward delay\t%s\n", timerToString(bridgeInfo.BridgeForwardDelay))
	fmt.Fprintf(&s, " aging time\t\t%s\n", timerToString(bridgeInfo.AgingTime))
	fmt.Fprintf(&s, " hello timer\t\t%s", timerToString(bridgeInfo.HelloTimer))
	fmt.Fprintf(&s, "\t\t\ttcn timer\t\t%s\n", timerToString(bridgeInfo.TCNTimer))
	fmt.Fprintf(&s, " topology change timer\t%s", timerToString(bridgeInfo.TopologyChangeTimer))
	fmt.Fprintf(&s, "\t\t\tgc timer\t\t%s\n", timerToString(bridgeInfo.GCTimer))
	fmt.Fprintf(&s, " flags\t\t\t")
	if bridgeInfo.TopologyChange {
		fmt.Fprintf(&s, "TOPOLOGY_CHANGE ")
	}
	if bridgeInfo.TopologyChangeDetected {
		fmt.Fprintf(&s, "TOPOLOGY_CHANGE_DETECTED ")
	}
	fmt.Fprint(&s, "\n\n\n")

	for _, portInfo := range bridgeInfo.Interfaces {
		fmt.Fprintf(&s, "%s (%d)\n", portInfo.Name, portInfo.PortNumber)
		fmt.Fprintf(&s, " port id\t\t%s", portInfo.PortID)
		fmt.Fprintf(&s, "\t\t\tport state\t\t  %d\n", portInfo.State) // TODO: How is the mapping to string (disabled, blocking, listening, learning, forwarding)?
		fmt.Fprintf(&s, " designated root\t%s", portInfo.DesignatedRoot)
		fmt.Fprintf(&s, "\tpath cost\t\t  %d\n", portInfo.PathCost)
		fmt.Fprintf(&s, " designated bridge\t%s", portInfo.DesignatedBridge)
		fmt.Fprintf(&s, "\tmessage age timer\t%s\n", timerToString(portInfo.MessageAgeTimer))
		fmt.Fprintf(&s, " designated port\t%s", portInfo.DesignatedPort)
		fmt.Fprintf(&s, "\t\t\tforward delay timer\t%s\n", timerToString(portInfo.ForwardDelayTimer))
		fmt.Fprintf(&s, " designated cost\t%d", portInfo.DesignatedCost)
		fmt.Fprintf(&s, "\t\t\thold timer\t\t%s\n", timerToString(portInfo.HoldTimer))
		fmt.Fprintf(&s, " flags\t\t\n")
		if portInfo.HairpinMode {
			fmt.Fprintf(&s, " hairpin mode\t\t1\n")
		}
		fmt.Fprint(&s, "\n\n")
	}

	fmt.Fprintf(out, "%s", s.String())

	return nil
}
