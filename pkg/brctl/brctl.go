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
	brctl_socket, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		return fmt.Errorf("unix.Socket: %w", err)
	}

	if _, err := executeIoctlStr(brctl_socket, unix.SIOCBRDELBR, name); err != nil {
		return fmt.Errorf("executeIoctlStr: %w", err)
	}

	return nil
}

// Addif adds an interface to the bridge provided
func Addif(bridge string, iface string) error {
	brctl_socket, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		return fmt.Errorf("unix.Socket: %w", err)
	}

	ifr, err := unix.NewIfreq(bridge)
	if err != nil {
		return fmt.Errorf("unix.NewIfreq: %w", err)
	}

	if_index, err := getIndexFromInterfaceName(iface)
	if err != nil {
		return fmt.Errorf("getIndexFromInterfaceName: %w", err)
	}
	ifr.SetUint32(uint32(if_index))

	if err := unix.IoctlIfreq(brctl_socket, unix.SIOCBRADDIF, ifr); err != nil {
		return fmt.Errorf("unix.IoctlIfreq: %w", err)
	}

	return nil
}

// Delif deleted a given interface from the bridge
func Delif(bridge string, iface string) error {
	brctl_socket, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		return fmt.Errorf("unix.Socket: %w", err)
	}

	ifr, err := unix.NewIfreq(bridge)
	if err != nil {
		return fmt.Errorf("unix.NewIfreq: %w", err)
	}

	if_index, err := getIndexFromInterfaceName(iface)
	if err != nil || if_index == 0 {
		return fmt.Errorf("getIndexFromInterfaceName: %w", err)
	}
	ifr.SetUint32(uint32(if_index))

	if err := unix.IoctlIfreq(brctl_socket, unix.SIOCBRDELIF, ifr); err != nil {
		return fmt.Errorf("unix.IoctlIfreq: %w", err)
	}

	return nil
}

// All bridges are in the virtfs under /sys/class/net/<name>/bridge/<item>, read info from there
// Update this function if BridgeInfo struct changes
func getBridgeInfo(name string) (BridgeInfo, error) {
	base_path := path.Join(BRCTL_SYS_NET, name, "bridge")
	bridge_id, err := os.ReadFile(path.Join(base_path, "bridge_id"))
	if err != nil {
		return BridgeInfo{}, fmt.Errorf("os.ReadFile: %w", err)
	}

	stp_enabled, err := os.ReadFile(path.Join(base_path, "stp_state"))
	if err != nil {
		return BridgeInfo{}, fmt.Errorf("os.ReadFile: %w", err)
	}

	stp_enabled_bool, err := strconv.ParseBool(strings.TrimSuffix(string(stp_enabled), "\n"))
	if err != nil {
		return BridgeInfo{}, fmt.Errorf("strconv.ParseBool: %w", err)
	}

	// get interfaceDir from sysfs
	interfaceDir, err := os.ReadDir(path.Join(BRCTL_SYS_NET, name, "brif"))
	if err != nil {
		return BridgeInfo{}, fmt.Errorf("os.ReadDir: %w", err)
	}

	interfaces := []string{}
	for i := range interfaceDir {
		interfaces = append(interfaces, interfaceDir[i].Name())
	}

	return BridgeInfo{
		Name:       name,
		BridgeID:   strings.TrimSuffix(string(bridge_id), "\n"),
		StpState:   stp_enabled_bool,
		Interfaces: interfaces,
	}, nil

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

	fmt.Fprintf(out, "%s\t\t%s\t\t%v\t\t%v\n", info.Name, info.BridgeID, info.StpState, ifaceString)
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
	brforward, err := os.ReadFile(path.Join(BRCTL_SYS_NET, bridge, "brforward"))
	if err != nil {
		return fmt.Errorf("Readfile(%q): %w", path.Join(BRCTL_SYS_NET, bridge, "brforward"), err)
	}

	fmt.Fprintf(out, "port no\tmac addr\t\tis_local?\n")

	for i := 0; i < len(brforward); i += 0x10 {
		chunk := brforward[i : i+0x10]
		mac := chunk[0:6]
		port_no := uint16(binary.BigEndian.Uint16(chunk[6:8]))
		is_local := uint8(chunk[9]) != 0

		fmt.Fprintf(out, "%3d\t%2x:%2x:%2x:%2x:%2x:%2x\t%v\n", port_no, mac[0], mac[1], mac[2], mac[3], mac[4], mac[5], is_local)
	}

	return nil
}

// Show will show some information on the bridge and its attached ports.
func Show(out io.Writer, names ...string) error {
	fmt.Fprint(out, "bridge name\tbridge id\tSTP enabled\t\tinterfaces")
	if len(names) == 0 {
		devices, err := os.ReadDir(BRCTL_SYS_NET)
		if err != nil {
			return fmt.Errorf("ReadDir(%q)= %w", BRCTL_SYS_NET, err)
		}

		for _, bridge := range devices {
			// check if device is bridge, aka if it has a bridge directory
			_, err := os.Stat(BRCTL_SYS_NET + bridge.Name() + "/bridge/")
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
	ageing_time, err := stringToJiffies(time)
	if err != nil {
		return fmt.Errorf("stringToJiffies(%q) = %w", time, err)
	}

	if err = setBridgeValue(name, BRCTL_AGEING_TIME, []byte(strconv.Itoa(ageing_time)), uint64(BRCTL_SET_AEGING_TIME)); err != nil {
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
	var stp_state int
	if state == "on" || state == "yes" {
		stp_state = 1
	} else {
		stp_state = 0
	}

	if err := setBridgeValue(bridge, BRCTL_STP_STATE, []byte(strconv.Itoa(stp_state)), uint64(BRCTL_SET_BRIDGE_PRIORITY)); err != nil {
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
	forward_delay, err := stringToJiffies(time)
	if err != nil {
		return fmt.Errorf("stringToJiffies(%q) = %w", time, err)
	}

	if err := setBridgeValue(bridge, BRCTL_FORWARD_DELAY, []byte(strconv.Itoa(forward_delay)), 0); err != nil {
		return fmt.Errorf("setBridgeValue: %w", err)
	}

	return nil
}

// Sethello sets the bridge's 'bridge hello time' to <time> seconds.
func Sethello(bridge string, time string) error {
	hello_time, err := stringToJiffies(time)
	if err != nil {
		return fmt.Errorf("stringToJiffies(%q) = %w", time, err)
	}

	if err := setBridgeValue(bridge, BRCTL_HELLO_TIME, []byte(strconv.Itoa(hello_time)), 0); err != nil {
		return fmt.Errorf("setBridgeValue: %w", err)
	}

	return nil
}

// Setmaxage sets the bridge's 'maximum message age' to <time> seconds.
func Setmaxage(bridge string, time string) error {
	max_age, err := stringToJiffies(time)
	if err != nil {
		return fmt.Errorf("stringToJiffies(%q) = %w", time, err)
	}

	if err := setBridgeValue(bridge, BRCTL_MAX_AGE, []byte(strconv.Itoa(max_age)), 0); err != nil {
		return fmt.Errorf("setBridgeValue: %w", err)
	}

	return nil
}

// Setpathcost sets the port cost of the port <port> to <cost>. This is a dimensionless metric.
func Setpathcost(bridge string, port string, cost string) error {
	path_cost, err := strconv.ParseUint(cost, 10, 64)
	if err != nil {
		return err
	}

	err = setPortBrportValue(port, BRCTL_PATH_COST, append([]byte(strconv.FormatUint(path_cost, 10)), BRCTL_SYS_SUFFIX))
	if err != nil {
		return fmt.Errorf("setPortBrportValue: %w", err)
	}

	return nil
}

// Setportprio sets the port <port>'s priority to <priority>.
// The priority value is an unsigned 8-bit quantity (a number between 0 and 255),
// and has no dimension. This metric is used in the designated port and root port selection algorithms.
func Setportprio(bridge string, port string, prio string) error {
	port_priority, err := strconv.Atoi(prio)
	if err != nil {
		return err
	}

	return setPortBrportValue(port, BRCTL_PRIORITY, []byte(strconv.Itoa(port_priority)))
}

// Hairpin sets the hairpin mode of the <port> attached to <bridge>
func Hairpin(bridge string, port string, hairpinmode string) error {
	var hairpin_mode string
	if hairpinmode == "on" {
		hairpin_mode = "1"
	} else {
		hairpin_mode = "0"
	}

	if err := setPortBrportValue(port, BRCTL_HAIRPIN, []byte(hairpin_mode)); err != nil {
		return fmt.Errorf("setPortBrportValue: %w", err)
	}

	return nil
}
