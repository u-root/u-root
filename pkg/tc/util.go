// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trafficctl

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/florianl/go-tc"
)

const (
	TimeUnitsPerSecs = 1000000
)

var ErrNoDevice = errors.New("no such device")

func getDevice(dev string) (net.Interface, error) {
	var ret net.Interface
	devs, err := net.Interfaces()
	if err != nil {
		return ret, err
	}

	var found bool
	for _, iface := range devs {
		if iface.Name == dev {
			ret = iface
			found = true
		}
	}

	if !found {
		return ret, fmt.Errorf("available devices: %q, but '%s': %w", devs, dev, ErrNoDevice)
	}

	return ret, nil
}

func parseTime(t string) (uint32, error) {
	var cutstring string
	multiplier := TimeUnitsPerSecs
	if strings.HasSuffix(t, "sec") {
		cutstring, _ = strings.CutSuffix(t, "sec")
		multiplier = TimeUnitsPerSecs
	} else if strings.HasSuffix(t, "secs") {
		cutstring, _ = strings.CutSuffix(t, "secs")
		multiplier = TimeUnitsPerSecs
	} else if strings.HasSuffix(t, "s") {
		cutstring, _ = strings.CutSuffix(t, "s")
		multiplier = TimeUnitsPerSecs
	}

	if strings.HasSuffix(t, "ms") {
		cutstring, _ = strings.CutSuffix(t, "ms")
		multiplier = TimeUnitsPerSecs / 1000
	} else if strings.HasSuffix(t, "msec") {
		cutstring, _ = strings.CutSuffix(t, "msec")
		multiplier = TimeUnitsPerSecs / 1000
	} else if strings.HasSuffix(t, "msecs") {
		cutstring, _ = strings.CutSuffix(t, "msecs")
		multiplier = TimeUnitsPerSecs / 1000
	}

	if strings.HasSuffix(t, "us") {
		cutstring, _ = strings.CutSuffix(t, "us")
		multiplier = TimeUnitsPerSecs / 1000000
	} else if strings.HasSuffix(t, "usec") {
		cutstring, _ = strings.CutSuffix(t, "usec")
		multiplier = TimeUnitsPerSecs / 1000000
	} else if strings.HasSuffix(t, "usecs") {
		cutstring, _ = strings.CutSuffix(t, "usecs")
		multiplier = TimeUnitsPerSecs / 1000000
	}

	val, err := strconv.Atoi(cutstring)
	if err != nil {
		return 0, err
	}
	if val < 0x0 || val >= 0x7FFFFFFF {
		return 0, ErrOutOfBounds
	}

	ret := uint32(val) * uint32(multiplier)

	return ret, nil
}

// ParseHandle takes a string in the form of XXXX: and returns the XXXX value as
// uint32 type shifted left by 16 bits.
func ParseHandle(h string) (uint32, error) {
	// split the string at :
	maj, _, ok := strings.Cut(h, ":")
	if !ok {
		return 0, ErrInvalidArg
	}

	major, err := strconv.ParseUint(maj, 16, 16)
	if err != nil {
		return 0, err
	}

	return uint32(major) << 16, nil
}

// ParseClassID takes a string which can have three forms:
// Form 1: "root", which returns tc.HandleRoot
// Form 2: "none", which returns 0 as classid
// Form 3: "XXXX:XXXX" (Major:Minor), will return uint32(major<<16)|uint32(minor)
func ParseClassID(p string) (uint32, error) {
	if p == "root" {
		return tc.HandleRoot, nil
	}

	if p == "none" {
		return 0, nil
	}

	// split the string at :
	mj, mn, ok := strings.Cut(p, ":")
	if !ok {
		return 0, ErrInvalidArg
	}

	major, err := strconv.ParseUint(mj, 16, 16)
	if err != nil {
		major = 0
	}

	if mn == "" {
		return uint32(major << 16), nil
	}
	minor, err := strconv.ParseUint(mn, 16, 16)
	if err != nil {
		return 0, err
	}

	return uint32(major<<16) | uint32(minor), nil
}

var ErrUnknownLinkLayer = errors.New("unknown linklayer value provided")

// RenderClassID is the inverse of ParseClassID.
func RenderClassID(classID uint32, printParent bool) string {
	if classID == tc.HandleRoot {
		return "root"
	}

	var parent string
	if printParent {
		parent = "parent "
	} else {
		parent = ""
	}

	major := classID >> 16
	minor := classID & 0xFFFF
	if minor == 0 {
		return fmt.Sprintf("%s%x:", parent, major)
	}

	return fmt.Sprintf("%s%x:%x", parent, major, minor)
}

// ParseLinkLayer takes a string of LinkLayer name and returns the
// equivalent uint8 representation.
func ParseLinkLayer(l string) (uint8, error) {
	for _, ll := range []struct {
		name string
		val  uint8
	}{
		{name: "ethernet", val: 1},
		{name: "atm", val: 2},
		{name: "ads1", val: 2},
	} {
		if ll.name == l {
			return ll.val, nil
		}
	}
	return 0xFF, ErrUnknownLinkLayer
}

// ParseSize takes a string of the form `0123456789gkmbit` and returns
// the equivalent size as uint64.
func ParseSize(s string) (uint64, error) {
	sizeStr := strings.TrimRight(s, "gkmbit")

	sz, err := strconv.ParseUint(sizeStr, 10, 32)
	if err != nil {
		return 0, err
	}

	unitMuliplier := strings.TrimLeft(s, "0123456789")

	switch unitMuliplier {
	case "k", "kb":
		sz *= 1024
	case "m", "mb":
		sz *= 1024 * 1024
	case "g", "gb":
		sz *= 1024 * 1024 * 1024
	case "kbit":
		sz *= 1024 / 8
	case "mbit":
		sz *= 1024 * 1024 / 8
	case "gbit":
		sz *= 1024 * 1024 * 1024 / 8
	}

	return sz, nil
}

// ParseRate takes a string of the form `0123456789bBgGKkMmTitps` and returns
// the equivalent rate as uint64.
func ParseRate(arg string) (uint64, error) {
	unit := strings.TrimLeft(arg, "0123456789")

	sizeStr := strings.TrimRight(arg, "bBgGKkMmTitps")
	sz, err := strconv.ParseUint(sizeStr, 10, 32)
	if err != nil {
		return 0, err
	}

	for _, entry := range []struct {
		unit  string
		value uint64
	}{
		{unit: "bit", value: 1},
		{unit: "Kibit", value: 1024},
		{unit: "mibit", value: 1024 * 1024},
		{unit: "gibit", value: 1024 * 1024 * 1024},
		{unit: "tibit", value: 1024 * 1024 * 1024 * 1024},
		{unit: "kbit", value: 1000},
		{unit: "mbit", value: 1000 * 1000},
		{unit: "gbit", value: 1000 * 1000 * 1000},
		{unit: "tit", value: 1000 * 1000 * 1000 * 1000},
		{unit: "Bps", value: 8},
		{unit: "KiBps", value: 8 * 1024},
		{unit: "Mibit", value: 8 * 1024 * 1024},
		{unit: "Gibit", value: 8 * 1024 * 1024 * 1024},
		{unit: "TiBps", value: 8 * 1024 * 1024 * 1024 * 1024},
		{unit: "KBps", value: 8 * 1000},
		{unit: "MBps", value: 8 * 1000 * 1000},
		{unit: "GBps", value: 8 * 1000 * 1000 * 1000},
		{unit: "TBps", value: 8 * 1000 * 1000 * 1000 * 1000},
	} {
		if entry.unit == unit {
			return (sz * entry.value) / 8, nil
		}
	}
	return 0, ErrInvalidArg
}

// GetHz reads the psched rate from /proc/net/psched and returns it.
func GetHz() (int, error) {
	const HZdef = 100
	psched, err := os.Open("/proc/net/psched")
	if err != nil {
		return 0, err
	}
	defer psched.Close()

	var gb1, gb2, nom, denom int

	fmt.Fscanf(psched, "%8x %8x %8x %8x",
		&gb1,
		&gb2,
		&nom,
		&denom)

	if nom == 1000000 {
		return denom, nil
	}

	return HZdef, nil
}

// CalcXMitTime takes a rate and size of uint64 and calculates the XMitTime.
func CalcXMitTime(rate uint64, size uint32) (uint32, error) {
	ret := TimeUnitsPerSecs * (float64(size) / float64(rate))
	if ret >= 0xFFFF_FFFF {
		ret = maxUint32
	}

	tickInUsec, err := getTickInUsec()
	if err != nil {
		return 0, err
	}

	return uint32(math.Ceil(ret * tickInUsec)), nil
}

// CalcXMitSize is the inverse of CalcXMitTime
func CalcXMitSize(rate uint64, ticks uint32) (uint32, error) {
	tickInUsec, err := getTickInUsec()
	if err != nil {
		return 0, err
	}
	usecs := float64(ticks) / tickInUsec
	return uint32(float64(rate) * usecs / TimeUnitsPerSecs), nil
}

func getTickInUsec() (float64, error) {
	psched, err := os.Open("/proc/net/psched")
	if err != nil {
		return 0, err
	}
	defer psched.Close()

	var t2us, us2t, clockRes, gb int

	fmt.Fscanf(psched, "%8x %8x %8x %8x",
		&t2us,
		&us2t,
		&clockRes,
		&gb)

	if clockRes == 1000000000 {
		t2us = us2t
	}

	clockFactor := float64(clockRes) / float64(TimeUnitsPerSecs)

	return float64(t2us) / float64(us2t) * clockFactor, nil
}

func getClockfactor() (uint32, error) {
	psched, err := os.Open("/proc/net/psched")
	if err != nil {
		return 0, err
	}
	defer psched.Close()

	var t2us, us2t, clockRes, gb int

	fmt.Fscanf(psched, "%8x %8x %8x %8x",
		&t2us,
		&us2t,
		&clockRes,
		&gb)

	if clockRes == 1000000000 {
		t2us = us2t
	}

	return uint32(clockRes / TimeUnitsPerSecs), nil
}

var ErrNoValidProto = errors.New("invalid protocol name")

// ParseProto takes an EtherType protocol string and returns the equivalent
// uint16 representation in network byte order.
func ParseProto(prot string) (uint16, error) {
	for _, p := range []struct {
		name string
		prot uint16
	}{
		{"802_3", 0x0001},
		{"802_2", 0x0004},
		{"ip", 0x800},
		{"arp", 0x806},
		{"aarp", 0x80F3},
		{"ipx", 0x8137},

		{"ipv6", 0x86DD},
	} {
		if p.name == prot {
			return HToNS(p.prot), nil
		}
	}
	return 0, ErrNoValidProto
}

// GetProtoFromInfo extracts the uint16 EtherType protocol value (in network
// byte order) from the tc.Object's Info field.
func GetProtoFromInfo(info uint32) uint16 {
	return uint16(info & 0xFFFF)
}

// RenderProto takes the uint16 representation of an EtherType protocol in
// network byte order and returns the equivalent string.
func RenderProto(proto uint16) string {
	pNr := NToHS(proto)

	for _, p := range []struct {
		name string
		prot uint16
	}{
		{"802_3", 0x0001},
		{"802_2", 0x0004},
		{"ip", 0x800},
		{"arp", 0x806},
		{"aarp", 0x80F3},
		{"ipx", 0x8137},
		{"ipv6", 0x86DD},
	} {
		if p.prot == pNr {
			return p.name
		}
	}

	return ""
}

// GetPrefFromInfo takes the uint32 representation of the Info field of
// tc.Object and returns the preference/priority value as uint16.
func GetPrefFromInfo(info uint32) uint16 {
	return uint16(info >> 16)
}

// GetInfoFromPrefAndProto combines the uint16 preference/priority value (in
// host byte order) and the uint16 EtherType protocol value (in network byte
// order) such that the combined value can be stored in the Info field of
// tc.Object.
func GetInfoFromPrefAndProto(hostPref, netProto uint16) uint32 {
	return (uint32(hostPref) << 16) | uint32(netProto)
}

// HToNS converts a uint16 value from host (native) byte order to network (big
// endian) byte order.
func HToNS(hostShort uint16) uint16 {
	netBytes := make([]byte, 2)

	// serialize hostShort into netBytes (this is where bytes may be swapped)
	binary.BigEndian.PutUint16(netBytes, hostShort)

	// reinterpret netBytes as a native value
	return binary.NativeEndian.Uint16(netBytes)
}

// NToHS converts a uint16 value from network (big endian) byte order to host
// (native) byte order.
func NToHS(netShort uint16) uint16 {
	netBytes := make([]byte, 2)

	// serialize netShort transparently
	binary.NativeEndian.PutUint16(netBytes, netShort)

	// parse netBytes into native value (this is where bytes may be swapped)
	return binary.BigEndian.Uint16(netBytes)
}

// HToNL converts a uint32 value from host (native) byte order to network (big
// endian) byte order.
func HToNL(hostLong uint32) uint32 {
	netBytes := make([]byte, 4)

	// serialize hostLong into netBytes (this is where bytes may be swapped)
	binary.BigEndian.PutUint32(netBytes, hostLong)

	// reinterpret netBytes as a native value
	return binary.NativeEndian.Uint32(netBytes)
}

// NToHL converts a uint32 value from network (big endian) byte order to host
// (native) byte order.
func NToHL(netLong uint32) uint32 {
	netBytes := make([]byte, 4)

	// serialize netLong transparently
	binary.NativeEndian.PutUint32(netBytes, netLong)

	// parse netBytes into native value (this is where bytes may be swapped)
	return binary.BigEndian.Uint32(netBytes)
}
