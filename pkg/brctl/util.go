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

// Info contains details about a bridge device and its ports.
// This information is not exhaustive, only the most important fields are included
type Info struct {
	Name                   string
	RootID                 string
	BridgeID               string
	RootPathCost           int
	MaxAge                 unix.Timeval
	HelloTime              unix.Timeval
	ForwardDelay           unix.Timeval
	BridgeMaxAge           unix.Timeval
	BridgeHelloTime        unix.Timeval
	BridgeForwardDelay     unix.Timeval
	RootPort               uint16
	StpEnabled             bool
	TopologyChange         bool
	TopologyChangeDetected bool
	AgingTime              unix.Timeval
	HelloTimer             unix.Timeval
	TCNTimer               unix.Timeval
	TopologyChangeTimer    unix.Timeval
	GCTimer                unix.Timeval
	Interfaces             []PortInfo
}

// NewInfo returns a new Info struct populated with the details of the bridge
// device with the given name and its ports. Information is read from sysfs.
func NewInfo(name string) (Info, error) {
	info := Info{
		Name: name,
	}
	var err error

	basePath := filepath.Join(BRCTL_SYS_NET, name, BRCTL_BRIDGE_DIR)

	info.RootID, err = readString(basePath, BRCTL_ROOT_ID)
	if err != nil {
		return Info{}, fmt.Errorf("root id: %w", err)
	}

	info.BridgeID, err = readString(basePath, BRCTL_BRIDGE_ID)
	if err != nil {
		return Info{}, fmt.Errorf("bridge id: %w", err)
	}

	info.RootPathCost, err = readInt(basePath, BRCTL_ROOT_PATH_COST)
	if err != nil {
		return Info{}, fmt.Errorf("root path cost: %w", err)
	}

	info.MaxAge, err = readTimeVal(basePath, BRCTL_MAX_AGE)
	if err != nil {
		return Info{}, fmt.Errorf("max age: %w", err)
	}

	info.BridgeMaxAge, err = readTimeVal(basePath, BRCTL_MAX_AGE)
	if err != nil {
		return Info{}, fmt.Errorf("bridge max age: %w", err)
	}

	info.HelloTime, err = readTimeVal(basePath, BRCTL_HELLO_TIME)
	if err != nil {
		return Info{}, fmt.Errorf("hello time: %w", err)
	}

	info.BridgeHelloTime, err = readTimeVal(basePath, BRCTL_HELLO_TIME)
	if err != nil {
		return Info{}, fmt.Errorf("bridge hello time: %w", err)
	}

	info.ForwardDelay, err = readTimeVal(basePath, BRCTL_FORWARD_DELAY)
	if err != nil {
		return Info{}, fmt.Errorf("forward delay: %w", err)
	}

	info.BridgeForwardDelay, err = readTimeVal(basePath, BRCTL_FORWARD_DELAY)
	if err != nil {
		return Info{}, fmt.Errorf("bridge forward delay: %w", err)
	}

	info.AgingTime, err = readTimeVal(basePath, BRCTL_AGEING_TIME)
	if err != nil {
		return Info{}, fmt.Errorf("aging time: %w", err)
	}

	info.HelloTimer, err = readTimeVal(basePath, BRCTL_HELLO_TIMER)
	if err != nil {
		return Info{}, fmt.Errorf("hello timer: %w", err)
	}

	info.TCNTimer, err = readTimeVal(basePath, BRCTL_TCN_TIMER)
	if err != nil {
		return Info{}, fmt.Errorf("tcn timer: %w", err)
	}

	info.TopologyChangeTimer, err = readTimeVal(basePath, BRCTL_TOPOLOGY_CHANGE_TIMER)
	if err != nil {
		return Info{}, fmt.Errorf("topology change timer: %w", err)
	}

	info.GCTimer, err = readTimeVal(basePath, BRCTL_GC_TIMER)
	if err != nil {
		return Info{}, fmt.Errorf("gc timer: %w", err)
	}

	rootPort, err := readInt(basePath, BRCTL_ROOT_PORT)
	if err != nil {
		return Info{}, fmt.Errorf("root port: %w", err)
	}

	if rootPort > 0xFFFF {
		return Info{}, fmt.Errorf("root port %w", strconv.ErrRange)
	}
	info.RootPort = uint16(rootPort)

	info.StpEnabled, err = readBool(basePath, BRCTL_STP_STATE)
	if err != nil {
		return Info{}, fmt.Errorf("stp state: %w", err)
	}

	info.TopologyChange, err = readBool(basePath, BRCTL_TOPOLOGY_CHANGE)
	if err != nil {
		return Info{}, fmt.Errorf("topology change: %w", err)
	}

	info.TopologyChangeDetected, err = readBool(basePath, BRCTL_TOPOLOGY_CHANGE_DETECTED)
	if err != nil {
		return Info{}, fmt.Errorf("topology change detected: %w", err)
	}

	// get respective ports from sysfs
	interfaceDir, err := os.ReadDir(filepath.Join(BRCTL_SYS_NET, name, BRCTL_BRIDGE_INTERFACE_DIR))
	if err != nil {
		return Info{}, fmt.Errorf("listing bridge interfaces: %w", err)
	}

	info.Interfaces = make([]PortInfo, 0, len(interfaceDir))
	for i := range interfaceDir {
		bridgeName := name
		portName := interfaceDir[i].Name()
		portInfo, err := NewPortInfo(bridgeName, portName)
		if err != nil {
			return Info{}, fmt.Errorf("getting port info for %q: %w", portName, err)
		}
		info.Interfaces = append(info.Interfaces, portInfo)
	}

	return info, nil
}

type PortInfo struct {
	Name              string
	PortID            string
	PortNumber        int
	State             int
	PathCost          int
	DesignatedRoot    string
	DesignatedCost    int
	DesignatedBridge  string
	DesignatedPort    string
	MessageAgeTimer   unix.Timeval
	ForwardDelayTimer unix.Timeval
	HoldTimer         unix.Timeval
	HairpinMode       bool
}

// NewPortInfo returns a new PortInfo struct populated with the details of the
// port with the given name on the bridge device with the given name. Information
// is read from sysfs.
func NewPortInfo(bridge, port string) (PortInfo, error) {
	info := PortInfo{
		Name: port,
	}

	basePath := filepath.Join(BRCTL_SYS_NET, bridge, BRCTL_BRIDGE_INTERFACE_DIR, port)

	var err error

	info.PortID, err = readString(basePath, BRCTL_PORT_ID)
	if err != nil {
		return PortInfo{}, fmt.Errorf("port id: %w", err)
	}

	portNumberStr, err := readString(basePath, BRCTL_PORT_NO)
	if err != nil {
		return PortInfo{}, fmt.Errorf("port number: %w", err)
	}

	portNumberStr = strings.TrimPrefix(portNumberStr, "0x")
	portNumber, err := strconv.ParseUint(portNumberStr, 16, 16)
	if err != nil {
		return PortInfo{}, fmt.Errorf("port number: %w", err)
	}
	info.PortNumber = int(portNumber)

	info.State, err = readInt(basePath, BRCTL_PORT_STATE)
	if err != nil {
		return PortInfo{}, fmt.Errorf("port state: %w", err)
	}

	info.PathCost, err = readInt(basePath, BRCTL_PATH_COST)
	if err != nil {
		return PortInfo{}, fmt.Errorf("path cost: %w", err)
	}

	info.DesignatedRoot, err = readString(basePath, BRCTL_DESIGNATED_ROOT)
	if err != nil {
		return PortInfo{}, fmt.Errorf("designated root: %w", err)
	}

	info.DesignatedCost, err = readInt(basePath, BRCTL_DESIGNATED_COST)
	if err != nil {
		return PortInfo{}, fmt.Errorf("designated cost: %w", err)
	}

	info.DesignatedBridge, err = readString(basePath, BRCTL_DESIGNATED_BRIDGE)
	if err != nil {
		return PortInfo{}, fmt.Errorf("designated bridge: %w", err)
	}

	info.DesignatedPort, err = readString(basePath, BRCTL_DESIGNATED_PORT)
	if err != nil {
		return PortInfo{}, fmt.Errorf("designated port: %w", err)
	}

	info.MessageAgeTimer, err = readTimeVal(basePath, BRCTL_MSG_AGE_TIMER)
	if err != nil {
		return PortInfo{}, fmt.Errorf("message age timer: %w", err)
	}

	info.ForwardDelayTimer, err = readTimeVal(basePath, BRCTL_FORWARD_DELAY_TIMER)
	if err != nil {
		return PortInfo{}, fmt.Errorf("forward delay timer: %w", err)
	}

	info.HoldTimer, err = readTimeVal(basePath, BRCTL_HOLD_TIMER)
	if err != nil {
		return PortInfo{}, fmt.Errorf("hold timer: %w", err)
	}

	info.HairpinMode, err = readBool(basePath, BRCTL_HAIRPIN)
	if err != nil {
		return PortInfo{}, fmt.Errorf("hairpin mode: %w", err)
	}

	return info, nil
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

	tv, err := time.ParseDuration(in + "s") // add 's' for a valid duration, brctl expects seconds on all time values
	if err != nil {
		return 0, fmt.Errorf("parse duration (%q): %w", in, err)
	}

	return int(tv.Seconds() * float64(hz)), nil
}

func readString(p string, obj string) (string, error) {
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
