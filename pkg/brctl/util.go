// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9

package brctl

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
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
		return 0, errno
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
		return 0, errors.New("interface has invalid index 0")
	}

	return int(ifrIfindex), nil
}

var ErrBridgeNotExist = errors.New("bridge does not exist")

// setBridgeValue sets the value of a bridge setting in the sysfs.
func setBridgeValue(bridge string, setting string, value []byte, _ uint64) error {
	bridgePath := filepath.Join(BRCTL_SYS_NET, bridge)
	settingPath := filepath.Join(bridgePath, "bridge", setting)
	val := append(value, BRCTL_SYS_SUFFIX) // values in sysfs are <bytes> + '\n'

	if _, err := os.Stat(bridgePath); err != nil {
		return ErrBridgeNotExist
	}

	if _, err := os.Stat(settingPath); err != nil {
		return fmt.Errorf("invalid setting %q", setting)
	}

	err := os.WriteFile(settingPath, val, 0)
	if errors.Is(err, os.ErrPermission) {
		return errors.New("permission denied")
	}

	return err
}

var ErrPortNotExist = errors.New("port does not exist")

// setPortValue sets the value of a port setting in the sysfs.
func setPortValue(port string, setting string, value []byte) error {
	ifacePath := filepath.Join(BRCTL_SYS_NET, port)
	settingPath := filepath.Join(ifacePath, "brport", setting)
	val := append(value, BRCTL_SYS_SUFFIX) // values in sysfs are <bytes> + '\n'

	if _, err := os.Stat(ifacePath); err != nil {
		return ErrPortNotExist
	}

	if _, err := os.Stat(settingPath); err != nil {
		return fmt.Errorf("invalid setting %q", setting)
	}

	err := os.WriteFile(settingPath, val, 0)
	if errors.Is(err, os.ErrPermission) {
		return errors.New("permission denied")
	}

	return err
}

// Convert a string representation of a time.Duration to jiffies
func stringToJiffies(in string) (int, error) {
	hz, err := sysconfhz()
	if err != nil {
		return 0, fmt.Errorf("eval system clock frequency :%w", err)
	}

	tv, err := time.ParseDuration(in)
	if err != nil {
		return 0, fmt.Errorf("parse duration (%q): %w", in, err)
	}

	return int(tv.Seconds() * float64(hz)), nil
}

func readID(p string, obj string) (string, error) {
	ret, err := os.ReadFile(filepath.Join(p, obj))
	if err != nil {
		return "", fmt.Errorf("os.ReadFile: %w", err)
	}

	return strings.TrimSuffix(string(ret), "\n"), nil
}

func readBool(p string, obj string) (bool, error) {
	valRaw, err := os.ReadFile(filepath.Join(p, obj))
	if err != nil {
		return false, fmt.Errorf("os.ReadFile: %w", err)
	}

	return strconv.ParseBool(strings.TrimSuffix(string(valRaw), "\n"))
}

func readInt(p string, obj string) (int, error) {
	valRaw, err := os.ReadFile(filepath.Join(p, obj))
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
