// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9th || a || pointer || value
// +build !plan9th a pointer value

package brctl

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"syscall"
	"time"
	"unsafe"

	"github.com/tklauser/go-sysconf"
	"golang.org/x/sys/unix"
)

var errno0 = syscall.Errno(0)

// Helper for issuing raw ioctl wi//go:build !plan9
type ifreqptr struct {
	Ifrn [16]byte
	ptr  unsafe.Pointer
}

// BridgeInfo contains information about a bridge
// This information is not exhaustive, only the most important fields are included
// Feel free to add more fields if needed.
type BridgeInfo struct {
	Name       string
	BridgeID   string
	StpState   bool
	Interfaces []string
}

func sysconfhz() (int, error) {
	clktck, err := sysconf.Sysconf(sysconf.SC_CLK_TCK)
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}
	return int(clktck), nil
}

func getIfreqOption(ifreq *unix.Ifreq, ptr unsafe.Pointer) ifreqptr {
	i := ifreqptr{ptr: ptr}
	copy(i.Ifrn[:], ifreq.Name())
	return i
}

// ioctl helpers
// TODO: maybe use ifreq.withData for this?
func executeIoctlStr(fd int, req uint, raw string) (int, error) {
	local_bytes := append([]byte(raw), 0)
	err_int, _, errno := syscall.Syscall(unix.SYS_IOCTL, uintptr(fd), uintptr(req), uintptr(unsafe.Pointer(&local_bytes[0])))
	if !errors.Is(errno, errno0) {
		return int(err_int), fmt.Errorf("syscall.Syscall: %s", errno)
	}
	return int(err_int), nil
}

func ioctl(fd int, req uint, addr uintptr) (int, error) {
	err_int, _, errno := syscall.Syscall(unix.SYS_IOCTL, uintptr(fd), uintptr(req), addr)
	if !errors.Is(errno, errno0) {
		return int(err_int), fmt.Errorf("syscall.Syscall: %s", errno)
	}
	return int(err_int), nil
}

// https://github.com/WireGuard/wireguard-go/blob/master/tun/tun_linux.go#L217
func getIndexFromInterfaceName(ifname string) (int, error) {
	ifreq, err := unix.NewIfreq(ifname)
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}

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
func setBridgeValue(bridge string, name string, value uint64, ioctlcode uint64) error {
	err := os.WriteFile(BRCTL_SYS_NET+bridge+"/bridge/"+name, []byte(strconv.FormatUint(value, 10)), 0)
	if err != nil {
		log.Printf("br_set_val: %v", err)
		// TODO: 2. Use ioctl as fallback
		return nil
	}
	return nil
}

// Get values for the bridge
func getBridgeValue(bridge string, name string) (string, error) {
	out, err := os.ReadFile(BRCTL_SYS_NET + bridge + "/bridge/" + name)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// Set the value of a port in a bridge
func setBridgePort(bridge string, port string, name string, value uint64, ioctlcode uint64) error {
	err := os.WriteFile(BRCTL_SYS_NET+port+"/brport/"+bridge+"/"+name, []byte(strconv.FormatUint(value, 10)), 0)
	if err != nil {
		log.Printf("br_set_port: %v", err)
		return nil
	}
	return nil
}

// Get the value of a port in a bridge
func getBridgePort(bridge string, port string, name string) (string, error) {
	out, err := os.ReadFile(BRCTL_SYS_NET + port + "/brport/" + bridge + "/" + name)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// Convert a string representation of a time.Duration to jiffies
func stringToJiffies(in string) (int, error) {
	hz, err := sysconfhz()
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	tv, err := time.ParseDuration(in)
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	return int(tv.Seconds() * float64(hz)), nil
}
