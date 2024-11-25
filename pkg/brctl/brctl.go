// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package brctl

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"
)

// Addbr adds a bridge with the provided name.
func Addbr(name string) error {
	brctlSocket, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		return fmt.Errorf("unix.Socket: %w", err)
	}

	if _, err := executeIoctlStr(brctlSocket, unix.SIOCBRADDBR, name); err != nil {
		return fmt.Errorf("executeIoctlStr: %w", err)
	}

	return nil
}

// Delbr deletes a bridge with the name provided.
func Delbr(name string) error {
	brctlSocket, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		return fmt.Errorf("unix.Socket: %w", err)
	}

	if _, err := executeIoctlStr(brctlSocket, unix.SIOCBRDELBR, name); err != nil {
		return fmt.Errorf("executeIoctlStr: %w", err)
	}

	return nil
}

// Addif adds an interface to the bridge provided
func Addif(bridge string, iface string) error {
	brctlSocket, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		return fmt.Errorf("unix.Socket: %w", err)
	}

	ifr, err := unix.NewIfreq(bridge)
	if err != nil {
		return fmt.Errorf("unix.NewIfreq: %w", err)
	}

	ifIndex, err := getIndexFromInterfaceName(iface)
	if err != nil {
		return fmt.Errorf("getIndexFromInterfaceName: %w", err)
	}
	ifr.SetUint32(uint32(ifIndex))

	if err := unix.IoctlIfreq(brctlSocket, unix.SIOCBRADDIF, ifr); err != nil {
		return fmt.Errorf("unix.IoctlIfreq: %w", err)
	}

	return nil
}

// Delif deleted a given interface from the bridge
func Delif(bridge string, iface string) error {
	brctlSocket, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		return fmt.Errorf("unix.Socket: %w", err)
	}

	ifr, err := unix.NewIfreq(bridge)
	if err != nil {
		return fmt.Errorf("unix.NewIfreq: %w", err)
	}

	ifIndex, err := getIndexFromInterfaceName(iface)
	if err != nil || ifIndex == 0 {
		return fmt.Errorf("getIndexFromInterfaceName: %w", err)
	}
	ifr.SetUint32(uint32(ifIndex))

	if err := unix.IoctlIfreq(brctlSocket, unix.SIOCBRDELIF, ifr); err != nil {
		return fmt.Errorf("unix.IoctlIfreq: %w", err)
	}

	return nil
}

// All bridges are in the virtfs under /sys/class/net/<name>/bridge/<item>, read info from there
// Update this function if BridgeInfo struct changes
func getBridgeInfo(name string) (BridgeInfo, error) {
	ret := BridgeInfo{
		Name: name,
	}
	var err error

	basePath := path.Join(BRCTL_SYS_NET, name, "bridge")

	// Read designated Root
	ret.DesignatedRoot, err = readID(basePath, BRCTL_DESIGNATED_ROOT)
	if err != nil {
		return BridgeInfo{}, err
	}

	// Read bridge id
	ret.BridgeID, err = readID(basePath, BRCTL_BRIDGEID)
	if err != nil {
		return BridgeInfo{}, err
	}

	// Read root path cost
	ret.RootPathCost, err = readInt(basePath, BRCTL_ROOT_PATH_COST)
	if err != nil {
		return BridgeInfo{}, err
	}

	// Read max age
	ret.MaxAge, err = readTimeVal(basePath, BRCTL_MAX_AGE)
	if err != nil {
		return BridgeInfo{}, err
	}

	// Read bridge max age
	ret.BridgeMaxAge, err = readTimeVal(basePath, BRCTL_MAX_AGE)
	if err != nil {
		return BridgeInfo{}, err
	}

	// Read hello time
	ret.HelloTime, err = readTimeVal(basePath, BRCTL_HELLO_TIME)
	if err != nil {
		return BridgeInfo{}, err
	}

	// Read bridge hello time
	ret.BridgeHelloTime, err = readTimeVal(basePath, BRCTL_HELLO_TIME)
	if err != nil {
		return BridgeInfo{}, err
	}

	// Read forward delay
	ret.ForwardDelay, err = readTimeVal(basePath, BRCTL_FORWARD_DELAY)
	if err != nil {
		return BridgeInfo{}, err
	}

	// Read bridge forward delay
	ret.BridgeForwardDelay, err = readTimeVal(basePath, BRCTL_FORWARD_DELAY)
	if err != nil {
		return BridgeInfo{}, err
	}

	// Read aging time
	ret.AgingTime, err = readTimeVal(basePath, BRCTL_AGEING_TIME)
	if err != nil {
		return BridgeInfo{}, err
	}

	// Read hello timer value
	ret.HelloTimerValue, err = readTimeVal(basePath, BRCTL_HELLO_TIMER_VALUE)
	if err != nil {
		return BridgeInfo{}, err
	}

	// Read tcn timer value
	ret.TCNTimerValue, err = readTimeVal(basePath, BRCTL_TCN_TIMER)
	if err != nil {
		return BridgeInfo{}, err
	}

	// Read Topology change timer
	ret.TopologyChangeTimerValue, err = readTimeVal(basePath, BRCTL_TOPOLOGY_CHANGE_TIMER)
	if err != nil {
		return BridgeInfo{}, nil
	}

	// Read GC Timer
	ret.GCTimerValue, err = readTimeVal(basePath, BRCTL_GC_TIMER_VALUE)
	if err != nil {
		return BridgeInfo{}, err
	}

	// Read root port
	rport, err := readInt(basePath, BRCTL_ROOT_PORT)
	if err != nil {
		return BridgeInfo{}, err
	}

	if rport > 0xFFFF {
		return BridgeInfo{}, strconv.ErrRange
	}
	ret.RootPort = uint16(rport)

	// Read STP state
	ret.StpEnabled, err = readBool(basePath, BRCTL_STP_STATE)
	if err != nil {
		return BridgeInfo{}, err
	}

	// Read trill enabled
	// ret.TrillEnabled, err = readBool(basePath, BRCTL_TRILL_ENABLED)
	// if err != nil {
	//  	return BridgeInfo{}, err
	// }

	// Read topology change
	ret.TopologyChange, err = readBool(basePath, BRCTL_TOPOLOGY_CHANGE)
	if err != nil {
		return BridgeInfo{}, err
	}

	// Read topology change detected
	ret.TopologyChangeDetected, err = readBool(basePath, BRCTL_TOPOLOGY_CHANGE_DETECTED)
	if err != nil {
		return BridgeInfo{}, err
	}

	// get interfaceDir from sysfs
	interfaceDir, err := os.ReadDir(path.Join(BRCTL_SYS_NET, name, BRCTL_BRIDGE_INTERFACE))
	if err != nil {
		return BridgeInfo{}, fmt.Errorf("os.ReadDir: %w", err)
	}

	ret.Interfaces = make([]string, 0)
	for i := range interfaceDir {
		ret.Interfaces = append(ret.Interfaces, interfaceDir[i].Name())
	}

	return ret, nil

}

// for now, only show essentials: bridge name, bridge id interfaces.
func showBridge(name string, out io.Writer) {
	info, err := getBridgeInfo(name)
	if err != nil {
		log.Fatalf("show_bridge: %v", err)
	}

	ifaceString := ""
	for _, iface := range info.Interfaces {
		ifaceString += iface + " "
	}

	fmt.Fprintf(out, ShowBridgeFmt, info.Name, info.BridgeID, info.StpEnabled, ifaceString)
}

// Showmacs shows a list of learned MAC addresses for this bridge.
// The mac addresses are stored in the first 6 bytes of /sys/class/net/<name>/brforward,
// The following format applies:
// 00-05: MAC address
// 06-08: port number
// 09-10: is_local
// 11-15: timeval (ignored for now)
func Showmacs(bridge string, out io.Writer) error {
	// parse sysf into 0x10 byte chunks
	brforward, err := os.ReadFile(path.Join(BRCTL_SYS_NET, bridge, BRCTL_BRFORWARD))
	if err != nil {
		return fmt.Errorf("Readfile(%q): %w", path.Join(BRCTL_SYS_NET, bridge, BRCTL_BRFORWARD), err)
	}

	fmt.Fprintf(out, "port no\tmac addr\t\tis_local?\n")

	for i := 0; i < len(brforward); i += 0x10 {
		chunk := brforward[i : i+0x10]
		mac := chunk[0:6]
		portNo := uint16(binary.BigEndian.Uint16(chunk[6:8]))
		isLocal := uint8(chunk[9]) != 0

		fmt.Fprintf(out, "%3d\t%2x:%2x:%2x:%2x:%2x:%2x\t%v\n", portNo, mac[0], mac[1], mac[2], mac[3], mac[4], mac[5], isLocal)
	}

	return nil
}

const ShowBridgeFmt = "%-15s %23s %15v %20v\n"

// Show will show some information on the bridge and its attached ports.
func Show(out io.Writer, names ...string) error {
	fmt.Fprintf(out, ShowBridgeFmt, "bridge name", "bridge id", "STP enabled", "interfaces")

	if len(names) == 0 {
		devices, err := os.ReadDir(BRCTL_SYS_NET)
		if err != nil {
			return fmt.Errorf("ReadDir(%q)= %w", BRCTL_SYS_NET, err)
		}

		for _, bridge := range devices {
			// check if device is bridge, aka if it has a bridge directory
			_, err := os.Stat(filepath.Join(BRCTL_SYS_NET, bridge.Name(), "bridge"))
			if err == nil {
				showBridge(bridge.Name(), out)
			}
		}
	} else {
		for _, name := range names {
			showBridge(name, out)
		}
	}
	return nil
}

// Setageingtime sets the ethernet (MAC) address ageing time, in seconds.
// After <time> seconds of not having seen a frame coming from a certain address,
// the bridge will time out (delete) that address from the Forwarding DataBase (fdb).
func Setageingtime(name string, time string) error {
	ageingTime, err := stringToJiffies(time)
	if err != nil {
		return fmt.Errorf("stringToJiffies(%q) = %w", time, err)
	}

	if err = setBridgeValue(name, BRCTL_AGEING_TIME, []byte(strconv.Itoa(ageingTime)), uint64(BRCTL_SET_AEGING_TIME)); err != nil {
		return fmt.Errorf("setBridgeValue: %w", err)
	}
	return nil
}

// Stp set the STP state of the bridge to on or off
// Enable using "on" or "yes", disable by providing anything else
// The manpage states:
// > If <state> is "on" or "yes"  the STP  will  be turned on, otherwise it will be turned off
// So this is actually the described behavior, not checking for "off" and "no"
func Stp(bridge string, state string) error {
	var stpState int
	if state == "on" || state == "yes" {
		stpState = 1
	} else {
		stpState = 0
	}

	if err := setBridgeValue(bridge, BRCTL_STP_STATE, []byte(strconv.Itoa(stpState)), uint64(BRCTL_SET_BRIDGE_PRIORITY)); err != nil {
		return fmt.Errorf("setBridgeValue: %w", err)
	}

	return nil
}

// Setbridgeprio sets the port <port>'s priority to <priority>.
// The priority value is an unsigned 8-bit quantity (a number between 0 and 255),
// and has no dimension. This metric is used in the designated port and root port selection algorithms.
func Setbridgeprio(bridge string, bridgePriority string) error {
	// parse bridgePriority to int
	prio, err := strconv.Atoi(bridgePriority)
	if err != nil {
		return err
	}

	if err := setBridgeValue(bridge, BRCTL_BRIDGE_PRIO, []byte(strconv.Itoa(prio)), 0); err != nil {
		return fmt.Errorf("setBridgeValue %w", err)
	}

	return nil
}

// Setfd sets the bridge's 'bridge forward delay' to <time> seconds.
func Setfd(bridge string, time string) error {
	forwardDelay, err := stringToJiffies(time)
	if err != nil {
		return fmt.Errorf("stringToJiffies(%q) = %w", time, err)
	}

	if err := setBridgeValue(bridge, BRCTL_FORWARD_DELAY, []byte(strconv.Itoa(forwardDelay)), 0); err != nil {
		return fmt.Errorf("setBridgeValue: %w", err)
	}

	return nil
}

// Sethello sets the bridge's 'bridge hello time' to <time> seconds.
func Sethello(bridge string, time string) error {
	helloTime, err := stringToJiffies(time)
	if err != nil {
		return fmt.Errorf("stringToJiffies(%q) = %w", time, err)
	}

	if err := setBridgeValue(bridge, BRCTL_HELLO_TIME, []byte(strconv.Itoa(helloTime)), 0); err != nil {
		return fmt.Errorf("setBridgeValue: %w", err)
	}

	return nil
}

// Setmaxage sets the bridge's 'maximum message age' to <time> seconds.
func Setmaxage(bridge string, time string) error {
	maxAge, err := stringToJiffies(time)
	if err != nil {
		return fmt.Errorf("stringToJiffies(%q) = %w", time, err)
	}

	if err := setBridgeValue(bridge, BRCTL_MAX_AGE, []byte(strconv.Itoa(maxAge)), 0); err != nil {
		return fmt.Errorf("setBridgeValue: %w", err)
	}

	return nil
}

// Setpathcost sets the port cost of the port <port> to <cost>. This is a dimensionless metric.
func Setpathcost(bridge string, port string, cost string) error {
	pathCost, err := strconv.ParseUint(cost, 10, 64)
	if err != nil {
		return err
	}

	err = setPortBrportValue(port, BRCTL_PATH_COST, append([]byte(strconv.FormatUint(pathCost, 10)), BRCTL_SYS_SUFFIX))
	if err != nil {
		return fmt.Errorf("setPortBrportValue: %w", err)
	}

	return nil
}

// Setportprio sets the port <port>'s priority to <priority>.
// The priority value is an unsigned 8-bit quantity (a number between 0 and 255),
// and has no dimension. This metric is used in the designated port and root port selection algorithms.
func Setportprio(bridge string, port string, prio string) error {
	portPriority, err := strconv.Atoi(prio)
	if err != nil {
		return err
	}

	return setPortBrportValue(port, BRCTL_PRIORITY, []byte(strconv.Itoa(portPriority)))
}

// Hairpin sets the hairpin mode of the <port> attached to <bridge>
func Hairpin(bridge string, port string, hairpinmode string) error {
	var hairpinMode string
	if hairpinmode == "on" {
		hairpinMode = "1"
	} else {
		hairpinMode = "0"
	}

	if err := setPortBrportValue(port, BRCTL_HAIRPIN, []byte(hairpinMode)); err != nil {
		return fmt.Errorf("setPortBrportValue: %w", err)
	}

	return nil
}

func ShowStp(out io.Writer, bridge string) error {
	bridgeInfo, err := getBridgeInfo(bridge)
	if err != nil {
		return err
	}

	var s strings.Builder

	fmt.Fprintf(&s, "%s\n", bridge)
	fmt.Fprintf(&s, " bridge id\t\t%s\n", bridgeInfo.BridgeID)
	fmt.Fprintf(&s, " designated root\t%s\n", bridgeInfo.DesignatedRoot)
	fmt.Fprintf(&s, " root port\t\t   %d\t\t\t", bridgeInfo.RootPort)
	fmt.Fprintf(&s, "path cost\t\t   %d\n", bridgeInfo.RootPathCost)
	fmt.Fprintf(&s, " max age\t\t%s", timerToString(bridgeInfo.MaxAge))
	fmt.Fprintf(&s, "\t\t\tbridge max age\t\t%s\n", timerToString(bridgeInfo.BridgeMaxAge))
	fmt.Fprintf(&s, " hello time\t\t%s", timerToString(bridgeInfo.HelloTime))
	fmt.Fprintf(&s, "\t\t\tbridge hello time\t%s\n", timerToString(bridgeInfo.BridgeHelloTime))
	fmt.Fprintf(&s, " forward delay\t\t%s", timerToString(bridgeInfo.ForwardDelay))
	fmt.Fprintf(&s, "\t\t\tbridge forward delay\t%s\n", timerToString(bridgeInfo.BridgeForwardDelay))
	fmt.Fprintf(&s, " aging time\t\t%s\n", timerToString(bridgeInfo.AgingTime))
	fmt.Fprintf(&s, " hello timer\t\t%s", timerToString(bridgeInfo.HelloTimerValue))
	fmt.Fprintf(&s, "\t\t\ttcn timer\t\t%s\n", timerToString(bridgeInfo.TCNTimerValue))
	fmt.Fprintf(&s, " topology change timer\t%s", timerToString(bridgeInfo.TopologyChangeTimerValue))
	fmt.Fprintf(&s, "\t\t\tgc timer\t\t%s\n", timerToString(bridgeInfo.GCTimerValue))
	fmt.Fprintf(&s, " flags\t\t\t")
	if bridgeInfo.TopologyChange {
		fmt.Fprintf(&s, "TOPOLOGY_CHANGE ")
	}
	if bridgeInfo.TopologyChangeDetected {
		fmt.Fprintf(&s, "TOPOLOGY_CHANGE_DETECTED ")
	}
	fmt.Fprint(&s, "\n\n\n")

	fmt.Fprintf(out, "%s", s.String())

	return nil
}
