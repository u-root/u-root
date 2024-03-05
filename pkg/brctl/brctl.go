// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package brctl

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

// Helper for issuing raw ioctl with a pointer value
type ifreqptr struct {
	Ifrn [16]byte
	ptr  unsafe.Pointer
}

// BridgeInfo contains information about a bridge
// This information is not exhaustive, only the most important fields are included
// Feel free to add more fields if needed.
type BridgeInfo struct {
	Name       string
	Bridge_id  string
	Stp_state  bool
	Interfaces []string
}

var (
	BRCTL_ADD_BRIDGE          = 2
	BRCTL_DEL_BRIDGE          = 3
	BRCTL_ADD_I               = 4
	BRCTL_DEL_I               = 5
	BRCTL_SET_AEGING_TIME     = 11
	BRCTL_SET_BRIDGE_PRIORITY = 15
	BRCTL_SET_PORT_PRIORITY   = 16
	BRCTL_SET_PATH_COST       = 17
)

func ifreq_option(ifreq *unix.Ifreq, ptr unsafe.Pointer) ifreqptr {
	i := ifreqptr{ptr: ptr}
	copy(i.Ifrn[:], ifreq.Name())
	return i
}

// ioctl helpers
// TODO: maybe use ifreq.withData for this?
func ioctl_str(fd int, req uint, raw string) (int, error) {
	local_bytes := append([]byte(raw), 0)
	err_int, _, err_str := syscall.Syscall(unix.SYS_IOCTL, uintptr(fd), uintptr(req), uintptr(unsafe.Pointer(&local_bytes[0])))
	return int(err_int), fmt.Errorf("%s", err_str)
}

func ioctl(fd int, req uint, addr uintptr) (int, error) {
	err_int, _, err_str := syscall.Syscall(unix.SYS_IOCTL, uintptr(fd), uintptr(req), addr)
	return int(err_int), fmt.Errorf("%s", err_str)
}

// https://github.com/WireGuard/wireguard-go/blob/master/tun/tun_linux.go#L217
func if_nametoindex(ifname string) (int, error) {
	ifreq, err := unix.NewIfreq(ifname)
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}
	fmt.Printf("ifr = %v, name = %v\n", ifreq, ifreq.Name())

	brctl_socket, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	err = unix.IoctlIfreq(brctl_socket, unix.SIOCGIFINDEX, ifreq)
	if err != nil {
		return 0, fmt.Errorf("%w %s", err, ifname)
	}

	ifr_ifindex := ifreq.Uint32()
	if ifr_ifindex == 0 {
		return 0, fmt.Errorf("interface %s not found", ifname)
	}

	return int(ifr_ifindex), nil
}

// set values for the bridge
// 1. Try to set the config options using the sysfs` bridge directory
// 2. Else use the ioctl interface
// TODO: hide this behind an interface and check in beforehand which to use
// TODO: how to parse the struct timeval and the jiffies?
/*
	@param bridge: name of the bridge
	@param name: name of the value to set
	@param value: value to set
	@param ioctlcode: old value to set, aka IOCTL control value, like BRCTL_SET_BRIDGE_MAX_AGE
*/
func br_set_val(bridge string, name string, value uint64, ioctlcode uint64) error {
	err := os.WriteFile("/sys/class/net/"+bridge+"/bridge/"+name, []byte(strconv.FormatUint(value, 10)), 0)
	if err != nil {
		log.Printf("br_set_val: %v", err)
		// TODO: 2. Use ioctl as fallback
		return nil
	}
	return nil
}

func br_set_port(bridge string, port string, name string, value uint64, ioctlcode uint64) error {
	err := os.WriteFile("/sys/class/net/"+port+"/brport/"+bridge+"/"+name, []byte(strconv.FormatUint(value, 10)), 0)
	if err != nil {
		log.Printf("br_set_port: %v", err)
		// TODO: 2. Use ioctl as fallback
		return nil
	}
	return nil
}

/*
Helper functions that convert seconds to jiffies and vice versa.
https://litux.nl/mirror/kerneldevelopment/0672327201/ch10lev1sec3.html
*/

func str_to_tv(in string) (unix.Timeval, error) {
	// cast string to long float
	time, err := strconv.ParseFloat(in, 64)
	if err != nil {
		return unix.Timeval{}, fmt.Errorf("%w", err)
	}

	// TODO: placeholder
	// 10^{-5}
	return unix.Timeval{
		Sec:  int64(time),
		Usec: int64(time / 100000),
	}, nil
}

// TODO: placeholder
func to_jiffies(tv unix.Timeval) int {
	return int(tv.Sec * 100)
}

func from_jiffies(jiffies int) int {
	return jiffies / 100
}

// subcommands
// https://elixir.bootlin.com/busybox/latest/source/networking/brctl.c#L583
// https://github.com/slackhq/nebula/blob/8822f1366c1111feb2f64fef229eed2024512104/overlay/tun_linux.go#L222
// can this also be done using IoctlIfreq?
func Addbr(name string) error {
	brctl_socket, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if _, err := ioctl_str(brctl_socket, unix.SIOCBRADDBR, name); err != nil {
		return fmt.Errorf("%w", err)
	}

	name_bytes, err := unix.ByteSliceFromString(name)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	args := []int64{int64(BRCTL_ADD_I), int64(uintptr(unsafe.Pointer(&name_bytes))), 0, 0}
	if _, err := ioctl(brctl_socket, unix.SIOCSIFBR, uintptr(unsafe.Pointer(&args))); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func Delbr(name string) error {
	brctl_socket, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if _, err := ioctl_str(brctl_socket, unix.SIOCBRDELBR, name); err != nil {
		return fmt.Errorf("%w", err)
	}

	name_bytes, err := unix.ByteSliceFromString(name)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	args := []int64{int64(BRCTL_DEL_BRIDGE), int64(uintptr(unsafe.Pointer(&name_bytes))), 0, 0}
	if _, err := ioctl(brctl_socket, unix.SIOCSIFBR, uintptr(unsafe.Pointer(&args))); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// create dummy device for testing: `sudo ip link add eth10 type dummy`
func Addif(name string, iface string) error {
	brctl_socket, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	ifr, err := unix.NewIfreq(name)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if_index, err := if_nametoindex(iface)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	ifr.SetUint32(uint32(if_index))

	//SIOCBRADDIF
	//SIOCGIFINDEX
	if err := unix.IoctlIfreq(brctl_socket, unix.SIOCBRADDIF, ifr); err != nil {
		return fmt.Errorf("%w", err)
	}

	// prepare args for ifr
	// apparently the go unix api does not support setting the second union value to a raw pointer
	// so we have to manually craft sth that looks like a struct ifreq
	var args = []int64{int64(BRCTL_ADD_I), int64(if_index), 0, 0}
	ifrd := ifreq_option(ifr, unsafe.Pointer(&args[0]))

	if _, err = ioctl(brctl_socket, unix.SIOCDEVPRIVATE, uintptr(unsafe.Pointer(&ifrd))); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func Delif(name string, iface string) error {
	brctl_socket, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	ifr, err := unix.NewIfreq(name)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if_index, err := if_nametoindex(iface)
	if err != nil || if_index == 0 {
		return fmt.Errorf("%w", err)
	}
	ifr.SetUint32(uint32(if_index))

	if err := unix.IoctlIfreq(brctl_socket, unix.SIOCBRDELIF, ifr); err != nil {
		return fmt.Errorf("%w", err)
	}

	// set args
	var args = []int64{int64(BRCTL_DEL_I), int64(if_index), 0, 0}
	ifrd := ifreq_option(ifr, unsafe.Pointer(&args[0]))

	if _, err = ioctl(brctl_socket, unix.SIOCDEVPRIVATE, uintptr(unsafe.Pointer(&ifrd))); err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

// All bridges are in the virtfs under /sys/class/net/<name>/bridge/<item>, read info from there
// Update this function if BridgeInfo struct changes
func getBridgeInfo(name string) (BridgeInfo, error) {
	base_path := "/sys/class/net/" + name + "/bridge/"
	bridge_id, err := os.ReadFile(base_path + "bridge_id")
	if err != nil {
		return BridgeInfo{}, fmt.Errorf("%w", err)
	}

	stp_enabled, err := os.ReadFile(base_path + "stp_state")
	if err != nil {
		return BridgeInfo{}, fmt.Errorf("%w", err)
	}

	stp_enabled_bool, err := strconv.ParseBool(strings.TrimSuffix(string(stp_enabled), "\n"))
	if err != nil {
		return BridgeInfo{}, fmt.Errorf("%w", err)
	}

	// TODO: get interfaces
	var interfaces = []string{"eth0", "eth1", "eth2"}

	return BridgeInfo{
		Name:       name,
		Bridge_id:  strings.TrimSuffix(string(bridge_id), "\n"),
		Stp_state:  stp_enabled_bool,
		Interfaces: interfaces,
	}, nil

}

// for now, only show essentials: bridge name, bridge id interfaces
func showBridge(name string) {
	info, err := getBridgeInfo(name)
	if err != nil {
		log.Fatalf("show_bridge: %v", err)
	}
	fmt.Printf("%s\t\t%s\t\t%v\t\t%v\n", info.Name, info.Bridge_id, info.Stp_state, info.Interfaces)
}

func Showmacs(name string) error {
	// The mac addresses are stored in the first 6 bytes of /sys/class/net/<name>/brforward,
	// The following format applies:
	// 00-05: MAC address
	// 06-08: port number
	// 09-10: is_local
	// 11-15: timeval (ignored for now)

	// parse sysf into 0x10 byte chunks
	brforward, err := os.ReadFile("/sys/class/net/" + name + "/brforward")
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	fmt.Printf("port no\tmac addr\t\tis_local?\n")

	for i := 0; i < len(brforward); i += 0x10 {
		chunk := brforward[i : i+0x10]
		mac := chunk[0:6]
		port_no := uint16(binary.BigEndian.Uint16(chunk[6:8]))
		is_local := uint8(chunk[9]) != 0

		fmt.Printf("%3d\t%2x:%2x:%2x:%2x:%2x:%2x\t%v\n", port_no, mac[0], mac[1], mac[2], mac[3], mac[4], mac[5], is_local)
	}

	return nil
}

func Show(names ...string) error {
	fmt.Println("bridge name\tbridge id\tSTP enabled\t\tinterfaces")
	if len(names) == 0 {
		devices, err := os.ReadDir("/sys/class/net")
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		for _, bridge := range devices {
			// check if device is bridge, aka if it has a bridge directory
			_, err := os.Stat("/sys/class/net/" + bridge.Name() + "/bridge/")
			if err == nil {
				showBridge(bridge.Name())
			}
		}

	} else {
		for _, name := range names {
			showBridge(name)
		}
	}
	return nil
}

// Spanning Tree Options
// TODO: this has a lot of boilerplate code. Maybe use a struct to store the values and use a switch statement to set the values
func Setageingtime(name string, time string) error {
	tv, err := str_to_tv(time)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	ageing_time := to_jiffies(tv)
	if err = br_set_val(name, "ageing_time", uint64(ageing_time), uint64(BRCTL_SET_AEGING_TIME)); err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

// weird, cannot find it in source code?
func Setgcint(name string, time string) error {
	// tv, err := str_to_tv(time)
	// if err != nil {
	// 	return fmt.Errorf("%w", err)
	// }

	// if err = br_set_val(name, "ageing_time", uint64(to_jiffies(tv)), uint64(BRCTL_SET_AEGING_TIME)); err != nil {
	// 	return fmt.Errorf("%w", err)
	// }
	return nil
}

// The manpage states:
// > If <state> is "on" or "yes"  the STP  will  be turned on, otherwise it will be turned off
// So this is actually the described behavior, not checking for "off" and "no"
// For coreutils see: brctl_cmd.c:284
func Stp(bridge string, state string) error {
	var stp_state uint64
	if state == "on" || state == "yes" {
		stp_state = 1
	} else {
		stp_state = 0
	}

	if err := br_set_val(bridge, "stp_state", stp_state, uint64(BRCTL_SET_BRIDGE_PRIORITY)); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// The manpage states only uint16 should be supplied, but brctl_cmd.c uses regular 'int'
// providing 2^16+1 results in 0 -> integer overflow
func Setbridgeprio(bridge string, bridge_priority string) error {
	prio, err := strconv.ParseUint(bridge_priority, 10, 16)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if err := br_set_val(bridge, "priority", prio, uint64(BRCTL_SET_AEGING_TIME)); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func Setfd(bridge string, time string) error {
	tv, err := str_to_tv(time)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	forward_delay := to_jiffies(tv)
	if err := br_set_val(bridge, "forward_delay", uint64(forward_delay), uint64(BRCTL_SET_AEGING_TIME)); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func Sethello(bridge string, time string) error {
	tv, err := str_to_tv(time)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	hello_time := to_jiffies(tv)
	if err := br_set_val(bridge, "hello_time", uint64(hello_time), uint64(BRCTL_SET_AEGING_TIME)); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func Setmaxage(bridge string, time string) error {
	tv, err := str_to_tv(time)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	max_age := to_jiffies(tv)
	if err := br_set_val(bridge, "max_age", uint64(max_age), uint64(BRCTL_SET_AEGING_TIME)); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// port ~= interface
func Setpathcost(bridge string, port string, cost string) error {
	path_cost, err := strconv.ParseUint(cost, 10, 64)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if err := br_set_port(bridge, port, "path_cost", path_cost, uint64(BRCTL_SET_PATH_COST)); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func Setportprio(bridge string, port string, prio string) error {
	port_priority, err := strconv.ParseUint(prio, 10, 64)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if err := br_set_port(bridge, port, "priority", port_priority, uint64(BRCTL_SET_PATH_COST)); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func Hairpin(bridge string, port string, hairpinmode string) error {
	var hairpin_mode uint64
	if hairpinmode == "on" {
		hairpin_mode = 1
	} else {
		hairpin_mode = 0
	}

	if err := br_set_port(bridge, port, "hairpin", hairpin_mode, 0); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
