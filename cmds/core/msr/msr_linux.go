// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains support functions for msr access for Linux.
package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func msrList(n string) []string {
	m, err := filepath.Glob(filepath.Join("/dev/cpu", n, "msr"))
	// This err will be if the glob was bad.
	if err != nil {
		log.Fatalf("No MSRs matched %v: %v", n, err)
	}
	// len will be zero for any of a number of reasons.
	if len(m) == 0 {
		log.Fatalf("No msrs found. Make sure your kernel is compiled with msrs, and you may need to 'sudo modprobe msr'. To see available msrs, ls /dev/cpu.")
	}
	return m
}

func openAll(m []string, o int) ([]*os.File, []error) {
	var (
		f    = make([]*os.File, len(m))
		errs = make([]error, len(m))
	)
	for i := range m {
		f[i], errs[i] = os.OpenFile(m[i], o, 0)
	}
	return f, errs
}

func doio(msr *os.File, addr uint32, f func(*os.File) error) error {
	if _, err := msr.Seek(int64(addr), 0); err != nil {
		return fmt.Errorf("Bad address %v: %v", addr, err)
	}
	return f(msr)
}

func rdmsr(m []string, addr uint32) ([]uint64, []error) {
	var regs = make([]uint64, len(m))

	f, errs := openAll(m, os.O_RDONLY)
	for i := range m {
		if errs[i] != nil {
			continue
		}
		errs[i] = doio(f[i], addr, func(port *os.File) error {
			return binary.Read(port, binary.LittleEndian, &regs[i])
		})
	}
	return regs, errs
}

func wrmsr(m []string, addr uint32, data uint64) []error {
	f, errs := openAll(m, os.O_RDWR)

	for i := range m {
		if errs[i] != nil {
			continue
		}
		errs[i] = doio(f[i], addr, func(port *os.File) error {
			return binary.Write(port, binary.LittleEndian, data)
		})
	}
	return errs
}
