// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// brctl - ethernet bridge administration
//
// Synopsis:
//
// brctl addbr <name> creates a new instance of the ethernet bridge
// brctl delbr <name> deletes the instance <name> of the ethernet bridge
// brctl show shows all current instances of the ethernet bridge
//
// brctl addif <brname> <ifname> will make the interface <ifname> a port of the bridge <brname>
// brctl delif <brname> <ifname> will detach the interface <ifname> from the bridge <brname>
// brctl show <brname> will show some information on the bridge and its attached ports
//
// brctl showmacs <brname> shows a list of learned MAC addresses for this bridge
// brctl setageingtime <brname> <time> sets the ethernet (MAC) address ageing time, in seconds [OPT]
// brctl setgcint <brname> <time> sets the garbage collection interval for the bridge <brname> to <time> seconds [OPT]
//
// TODO: Spanning Tree Protocol
// See: https://elixir.bootlin.com/busybox/latest/source/networking/brctl.c
// Author:
//
//	Leon Gross (leon.gross@9elements.com)
package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

var (
	BRCTL_ADD_BRIDGE = 2
	BRCTL_DEL_BRIDGE = 3
	BRCTL_ADD_I      = 4
	BRCTL_DEL_I      = 5
)

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

// cli
const usage = "brctl [commands]"

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

// subcommands
// https://elixir.bootlin.com/busybox/latest/source/networking/brctl.c#L583
// https://github.com/slackhq/nebula/blob/8822f1366c1111feb2f64fef229eed2024512104/overlay/tun_linux.go#L222
// can this also be done using IoctlIfreq?
func addbr(name string) error {
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

func delbr(name string) error {
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

// create dummy device for testing: sudo ip link add eth10 type dummy
func addif(name string, iface string) error {
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

	// SIOCBRADDIF
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

func delif(name string, iface string) error {
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

func showmacs(name string) error {
	// The mac addresses are stored in the first 6 bytes of /sys/class/net/<name/brforward,
	// The following format applies:
	// 00-05: MAC address
	// 06-08: port number
	// 08-09: is_local
	// 25-128: timeval (ignored for now)

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
		is_local := uint8(chunk[8]) != 0

		fmt.Printf("%3d\t%2x:%2x:%2x:%2x:%2x:%2x\t%v\n", port_no, mac[0], mac[1], mac[2], mac[3], mac[4], mac[5], is_local)
	}

	return nil
}

func show(names ...string) error {
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

// runner
func run(out io.Writer, argv []string) error {
	var err error

	command := argv[0]
	args := argv[1:]

	switch command {
	case "addbr":
		if len(args) != 1 {
			return fmt.Errorf("too few args")
		}
		err = addbr(args[0])

	case "delbr":
		if len(args) != 1 {
			return fmt.Errorf("too few args")
		}
		err = delbr(args[0])

	case "addif":
		if len(args) != 2 {
			return fmt.Errorf("too few args")
		}
		err = addif(args[0], args[1])

	case "delif":
		if len(args) != 2 {
			return fmt.Errorf("too few args")
		}
		err = delif(args[0], args[1])

	case "show":
		err = show(args...)

	case "showmacs":
		if len(args) != 1 {
			return fmt.Errorf("too few args")
		}
		err = showmacs(args[0])

	default:
		return fmt.Errorf("unknown command: %s", command)
	}

	return err
}

func main() {
	argv := os.Args

	if len(argv) < 2 {
		log.Fatal(usage)
		os.Exit(1)
	}

	if err := run(os.Stdout, argv[1:]); err != nil {
		log.Fatalf("brctl: %v", err)
	}
}
