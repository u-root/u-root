// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains support functions for msr access for Linux.
package msr

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/u-root/cpuid"
)

// CPUs is a slice of the various cpus to read or write the MSR to.
type CPUs []uint64

func parseCPUs(s string) (CPUs, error) {
	cpus := make(CPUs, 0)
	// We expect the format to be "0-5,7-8..." or we could also get just one cpu.
	// We're unlikely to get more than one range since we're looking at present cpus,
	// but handle it just in case.
	ranges := strings.Split(strings.TrimSpace(s), ",")
	for _, r := range ranges {
		if len(r) == 0 {
			continue
		}
		// Split on a - if it exists.
		cs := strings.Split(r, "-")
		switch len(cs) {
		case 1:
			u, err := strconv.ParseUint(cs[0], 0, 64)
			if err != nil {
				return nil, fmt.Errorf("unknown cpu range: %v, failed to parse %w", r, err)
			}
			cpus = append(cpus, uint64(u))
		case 2:
			ul, err := strconv.ParseUint(cs[0], 0, 64)
			if err != nil {
				return nil, fmt.Errorf("unknown cpu range: %v, failed to parse %w", r, err)
			}
			uh, err := strconv.ParseUint(cs[1], 0, 64)
			if err != nil {
				return nil, fmt.Errorf("unknown cpu range: %v, failed to parse %w", r, err)
			}
			if ul > uh {
				return nil, fmt.Errorf("invalid cpu range, upper bound greater than lower: %v", r)
			}
			for i := ul; i <= uh; i++ {
				cpus = append(cpus, uint64(i))
			}
		default:
			return nil, fmt.Errorf("unknown cpu range: %v", r)
		}
	}
	if len(cpus) == 0 {
		return nil, fmt.Errorf("no cpus found, input was %v", s)
	}
	sort.Slice(cpus, func(i, j int) bool { return cpus[i] < cpus[j] })
	// Remove duplicates
	for i := 0; i < len(cpus)-1; i++ {
		if cpus[i] == cpus[i+1] {
			cpus = append(cpus[:i], cpus[i+1:]...)
			i--
		}
	}
	return cpus, nil
}

// AllCPUs searches for actual present CPUs instead of relying on the glob.
// This is more accurate than what's presented in /dev/cpu/*/msr
func AllCPUs() (CPUs, error) {
	v, err := os.ReadFile("/sys/devices/system/cpu/present")
	if err != nil {
		return nil, err
	}
	return parseCPUs(string(v))
}

// GlobCPUs allow the user to specify CPUs using a glob as one would in /dev/cpu
func GlobCPUs(g string) (CPUs, []error) {
	var hadErr bool

	f, err := filepath.Glob(filepath.Join("/dev/cpu", g, "msr"))
	if err != nil {
		return nil, []error{err}
	}

	c := make([]uint64, len(f))
	errs := make([]error, len(f))
	for i, v := range f {
		c[i], errs[i] = strconv.ParseUint(filepath.Base(filepath.Dir(v)), 0, 64)
		if errs[i] != nil {
			hadErr = true
		}
	}
	if hadErr {
		return nil, errs
	}
	return c, nil
}

// String implements String() for MSR.
func (m MSR) String() string {
	return fmt.Sprintf("%#x", uint32(m))
}

// String pretty prints the list of CPUs. For example: 1-2,4
func (c CPUs) String() string {
	if len(c) == 0 {
		return "nil"
	}
	sort.Slice(c, func(i, j int) bool { return c[i] < c[j] })

	var s []string
	for i := 0; i < len(c); i++ {
		// Find the last CPU in this continuous range.
		j := i
		for j+1 < len(c) && c[j]+1 == c[j+1] {
			j++
		}

		if i == j {
			// Continuous set of size 1.
			s = append(s, fmt.Sprintf("%d", c[i]))
		} else {
			// Multiple CPUs in continous set.
			s = append(s, fmt.Sprintf("%d-%d", c[i], c[j]))
		}

		i = j // Skip over set.
	}
	return strings.Join(s, ",")
}

func (c CPUs) paths() []string {
	p := make([]string, len(c))

	for i, v := range c {
		p[i] = filepath.Join("/dev/cpu", strconv.FormatUint(v, 10), "msr")
	}
	return p
}

// Read reads an MSR from a set of CPUs.
func (m MSR) Read(c CPUs) ([]uint64, []error) {
	var hadErr bool
	regs := make([]uint64, len(c))

	paths := c.paths()
	f, errs := openAll(paths, os.O_RDONLY)
	if errs != nil {
		return nil, errs
	}
	errs = make([]error, len(f))
	for i := range f {
		defer f[i].Close()
		errs[i] = doIO(f[i], m, func(port *os.File) error {
			return binary.Read(port, binary.LittleEndian, &regs[i])
		})
		if errs[i] != nil {
			hadErr = true
		}
	}
	if hadErr {
		return nil, errs
	}

	return regs, nil
}

// Write writes values to an MSR on a set of CPUs.
// The data must be passed as a scalar (single value) or a slice.
// If the data slice has more than one element,
// the length of the data slice and the CPU slice must be the same.
// If a single value is given, it will be written to all the CPUs.
// If multiple data values are given, each will be written to its corresponding
// CPU.
func (m MSR) Write(c CPUs, data ...uint64) []error {
	var hadErr bool

	if len(data) == 1 {
		// Expand value to all cpus
		for i := 1; i < len(c); i++ {
			data = append(data, data[0])
		}
	}
	if len(data) != len(c) {
		return []error{fmt.Errorf("mismatched lengths: cpus %v, data %v", c, data)}
	}

	paths := c.paths()
	f, errs := openAll(paths, os.O_RDWR)

	if errs != nil {
		return errs
	}
	errs = make([]error, len(f))
	for i := range f {
		defer f[i].Close()
		errs[i] = doIO(f[i], m, func(port *os.File) error {
			return binary.Write(port, binary.LittleEndian, data[i])
		})
		if errs[i] != nil {
			hadErr = true
		}
	}
	if hadErr {
		return errs
	}
	return nil
}

// testAndSetMaybe takes a mask of bits to clear and to set, and applies them to the specified MSR in
// each of the CPUs. It will set the MSR only if the value is different and a set is requested.
// If the MSR is different for any reason that is an error.
func (m MSR) testAndSetMaybe(c CPUs, clearMask uint64, setMask uint64, set bool) []error {
	var hadErr bool
	paths := c.paths()
	f, errs := openAll(paths, os.O_RDWR)

	if errs != nil {
		return errs
	}
	errs = make([]error, len(f))
	for i := range f {
		defer f[i].Close()
		errs[i] = doIO(f[i], m, func(port *os.File) error {
			var v uint64
			err := binary.Read(port, binary.LittleEndian, &v)
			if err != nil {
				return err
			}
			n := v & ^clearMask
			n |= setMask
			// We write only if there is a change. This is to avoid
			// cases where we try to set a lock bit again, but the bit is
			// already set
			if n != v && set {
				return binary.Write(port, binary.LittleEndian, n)
			}
			if n != v {
				return fmt.Errorf("%#x", v)
			}
			return nil
		})
		if errs[i] != nil {
			hadErr = true
		}
	}
	if hadErr {
		return errs
	}
	return nil
}

// Test takes a mask of bits to clear and to set, and returns an error for those
// that do not match.
func (m MSR) Test(c CPUs, clearMask uint64, setMask uint64) []error {
	return m.testAndSetMaybe(c, clearMask, setMask, false)
}

// TestAndSet takes a mask of bits to clear and to set, and applies them to the specified MSR in
// each of the CPUs. Note that TestAndSet does not write if the mask does not change the MSR.
func (m MSR) TestAndSet(c CPUs, clearMask uint64, setMask uint64) []error {
	return m.testAndSetMaybe(c, clearMask, setMask, true)
}

// Locked verifies that for all MSRVal's for the CPU vendor, the MSRs are locked.
// TODO: this is another Intel-specific function at present.
func Locked() error {
	vendor := cpuid.VendorIdentificatorString
	// TODO: support more than Intel. Use the vendor id to look up msrs.
	if vendor != "GenuineIntel" {
		return fmt.Errorf("sorry, this package only supports Intel at present")
	}

	cpus, err := AllCPUs()
	if err != nil {
		return err
	}

	var allerrors string
	for _, m := range LockIntel {
		Debug("MSR %v on cpus %v, clearmask 0x%8x, setmask 0x%8x", m.Addr, cpus, m.Clear, m.Set)
		if m.WriteOnly {
			continue
		}
		errs := m.Addr.Test(cpus, m.Clear, m.Set)

		for i, e := range errs {
			if e != nil {
				allerrors += fmt.Sprintf("[cpu%d(%s)%v ", cpus[i], m.String(), e)
			}
		}
	}

	if allerrors != "" {
		return fmt.Errorf("%s: %v", vendor, allerrors)
	}
	return nil
}

func openAll(m []string, o int) ([]*os.File, []error) {
	var (
		f      = make([]*os.File, len(m))
		errs   = make([]error, len(m))
		hadErr bool
	)
	for i := range m {
		f[i], errs[i] = os.OpenFile(m[i], o, 0)
		if errs[i] != nil {
			hadErr = true
			f[i] = nil // Not sure if I need to do this, it doesn't seem guaranteed.
		}
	}
	if hadErr {
		for i := range f {
			if f[i] != nil {
				f[i].Close()
			}
		}
		return nil, errs
	}
	return f, nil
}

func doIO(msr *os.File, addr MSR, f func(*os.File) error) error {
	if _, err := msr.Seek(int64(addr), 0); err != nil {
		return fmt.Errorf("bad address %v: %w", addr, err)
	}
	return f(msr)
}
