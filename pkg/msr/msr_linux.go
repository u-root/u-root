// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains support functions for msr access for Linux.
package msr

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

// CPUs is a slice of the various cpus to read or write the MSR to.
type CPUs []uint64

// AllCPUs is a helper variable that lists all the cpus available on the machine.
var AllCPUs CPUs

func init() {
	var errs []error

	if AllCPUs, errs = getAllCPUs(); errs != nil {
		// Should I fatal here? panic?
		log.Print("Failed to find all CPUs from /dev/cpu")
		log.Print(errs)
	}
}

// GlobCPUs allow the user to specify CPUs using a glob as one would in /dev/cpu
func GlobCPUs(g string) (CPUs, []error) {
	var hadErr bool

	f, err := filepath.Glob(filepath.Join("/dev/cpu", g))
	if err != nil {
		return nil, []error{err}
	}

	c := make([]uint64, len(f))
	errs := make([]error, len(f))
	for i, v := range f {
		c[i], errs[i] = strconv.ParseUint(filepath.Base(v), 0, 64)
		if errs[i] != nil {
			hadErr = true
		}
	}
	if hadErr {
		return nil, errs
	}
	return c, nil
}

func getAllCPUs() (CPUs, []error) {
	return GlobCPUs("*")
}

// MSR is the address of the MSR we want to target.
type MSR uint32

func (m MSR) String() string {
	return fmt.Sprintf("%#x", uint32(m))
}

func (c CPUs) paths() []string {
	var p = make([]string, len(c))

	for i, v := range c {
		p[i] = filepath.Join("/dev/cpu", strconv.Itoa(int(v)), "msr")
	}
	return p
}

func (m MSR) Read(c CPUs) ([]uint64, []error) {
	var hadErr bool
	var regs = make([]uint64, len(c))

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

// Write writes the corresponding data to their specified msrs
func (m MSR) Write(c CPUs, data ...uint64) []error {
	var hadErr bool

	if len(data) != len(c) && len(data) != 1 {
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

// MaskBits takes a mask of bits to clear and to set, and applies them to the specified MSR in
// each of the CPUs.
func (m MSR) MaskBits(c CPUs, clearMask uint64, setMask uint64) []error {
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
			v &= ^clearMask
			v |= setMask
			return binary.Write(port, binary.LittleEndian, v)
		})
	}
	return errs
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
		return fmt.Errorf("bad address %v: %v", addr, err)
	}
	return f(msr)
}
