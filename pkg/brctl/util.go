// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9

package brctl

import (
	"fmt"
	"os"
	"path"
	"strconv"
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
	Name                     string
	DesignatedRoot           string
	BridgeID                 string
	RootPathCost             int
	MaxAge                   unix.Timeval
	HelloTime                unix.Timeval
	ForwardDelay             unix.Timeval
	BridgeMaxAge             unix.Timeval
	BridgeHelloTime          unix.Timeval
	BridgeForwardDelay       unix.Timeval
	RootPort                 uint16
	StpEnabled               bool
	TrillEnabled             bool
	TopologyChange           bool
	TopologyChangeDetected   bool
	AgingTime                unix.Timeval
	HelloTimerValue          unix.Timeval
	TCNTimerValue            unix.Timeval
	TopologyChangeTimerValue unix.Timeval
	GCTimerValue             unix.Timeval
	Interfaces               []string
}

func sysconfhz() (int, error) {
	clktck, err := sysconf.Sysconf(sysconf.SC_CLK_TCK)
	if err != nil {
		return 0, err
	}
	return int(clktck), nil
}

func executeIoctlStr(fd int, req uint, raw string) (int, error) {
	localBytes := append([]byte(raw), 0)
	_, _, errno := syscall.Syscall(unix.SYS_IOCTL, uintptr(fd), uintptr(req), uintptr(unsafe.Pointer(&localBytes[0])))
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

	brctlSocket, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		return 0, err
	}

	err = unix.IoctlIfreq(brctlSocket, unix.SIOCGIFINDEX, ifreq)
	if err != nil {
		return 0, err
	}

	ifrIfindex := ifreq.Uint32()
	if ifrIfindex == 0 {
		return 0, fmt.Errorf("interface %s not found", ifname)
	}

	return int(ifrIfindex), nil
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

func readID(p string, obj string) (string, error) {
	ret, err := os.ReadFile(path.Join(p, obj))
	if err != nil {
		return "", fmt.Errorf("os.ReadFile: %w", err)
	}

	return strings.TrimSuffix(string(ret), "\n"), nil
}

func readBool(p string, obj string) (bool, error) {
	valRaw, err := os.ReadFile(path.Join(p, obj))
	if err != nil {
		return false, fmt.Errorf("os.ReadFile: %w", err)
	}

	return strconv.ParseBool(strings.TrimSuffix(string(valRaw), "\n"))
}

func readInt(p string, obj string) (int, error) {
	valRaw, err := os.ReadFile(path.Join(p, obj))
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(strings.TrimSuffix(string(valRaw), "\n"))
}

func timerToString(t unix.Timeval) string {
	var s strings.Builder

	fmt.Fprintf(&s, "%4.2d.%-2.2d", t.Sec, t.Usec/10000)

	return s.String()
}
