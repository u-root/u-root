// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9
// +build !plan9

package brctl

import (
	"fmt"
	"os"
	"path"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/tklauser/go-sysconf"
	"golang.org/x/sys/unix"
)

var errno0 = syscall.Errno(0)

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
		return 0, err
	}
	return int(clktck), nil
}

// ioctl helpers
// TODO: maybe use ifreq.withData for this?
func executeIoctlStr(fd int, req uint, raw string) (int, error) {
	local_bytes := append([]byte(raw), 0)
	_, _, errno := syscall.Syscall(unix.SYS_IOCTL, uintptr(fd), uintptr(req), uintptr(unsafe.Pointer(&local_bytes[0])))
	if errno != 0 {
		return 0, fmt.Errorf("syscall.Syscall: %w", errno)
	}
	return 0, nil
}

func getIndexFromInterfaceName(ifname string) (int, error) {
	ifreq, err := unix.NewIfreq(ifname)
	if err != nil {
		return 0, err
	}

	brctl_socket, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		return 0, err
	}

	err = unix.IoctlIfreq(brctl_socket, unix.SIOCGIFINDEX, ifreq)
	if err != nil {
		return 0, err
	}

	ifr_ifindex := ifreq.Uint32()
	if ifr_ifindex == 0 {
		return 0, fmt.Errorf("interface %s not found", ifname)
	}

	return int(ifr_ifindex), nil
}

// set values for the bridge
// all values in the sysfs are of type <bytes> + '\n'
func setBridgeValue(bridge string, name string, value []byte, _ uint64) error {
	err := os.WriteFile(path.Join(BRCTL_SYS_NET, bridge, "bridge", name), append(value, BRCTL_SYS_SUFFIX), 0)
	if err != nil {
		return err
	}
	return nil
}

// Get values for the bridge
// For some reason these values have a '\n' (0x0a) as a suffix, so we need to trim it
func getBridgeValue(bridge string, name string) (string, error) {
	out, err := os.ReadFile(path.Join(BRCTL_SYS_NET, bridge, "bridge", name))
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(out), "\n"), nil
}

func setPortBrportValue(port string, name string, value []byte) error {
	err := os.WriteFile(path.Join(BRCTL_SYS_NET, port, "brport", name), append(value, BRCTL_SYS_SUFFIX), 0)
	if err != nil {
		return err
	}
	return nil
}

func getPortBrportValue(port string, name string) (string, error) {
	out, err := os.ReadFile(path.Join(BRCTL_SYS_NET, port, "brport", name))
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// Convert a string representation of a time.Duration to jiffies
func stringToJiffies(in string) (int, error) {
	hz, err := sysconfhz()
	if err != nil {
		return 0, fmt.Errorf("sysconfhz():%w", err)
	}

	tv, err := time.ParseDuration(in)
	if err != nil {
		return 0, fmt.Errorf("ParseDuration(%q) = %w", in, err)
	}

	return int(tv.Seconds() * float64(hz)), nil
}
