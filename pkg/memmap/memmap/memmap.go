// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package memmap contains funtions that let us figure out system memory.
package memmap

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"unicode"
)

type MemMapType int

// These types and consts are compatible with kexec. Mistake?
const (
	Ram MemMapType = iota
	Reserved
	ACPITables
	ACPINonVolatileStorage
	Uncached
)

type MemoryRange struct {
	Start uint64
	End   uint64
	Type  MemMapType
}

var (
	Types = map[string]MemMapType{
		"System RAM":                Ram,
		"reserved":                  Reserved,
		"ACPI Tables":               ACPITables,
		"ACPI Non-volatile Storage": ACPINonVolatileStorage,
		"uncached":                  Uncached,
	}

	Names = []string{"System RAM", "reserved", "ACPI Tables", "ACPI Non-volatile Storage", "uncached"}
)

func mmVal(dir, file string) (uint64, error) {
	s, err := ioutil.ReadFile(filepath.Join(dir, file))
	if err != nil {
		return 0, err
	}
	// strip newline at end of string. Bleah.
	return strconv.ParseUint(string(s)[0:len(s)-1], 0, 64)
}

func mmr(name string) (*MemoryRange, error) {
	var m = &MemoryRange{Type: Reserved}
	var err error
	m.Start, err = mmVal(name, "start")
	if err != nil {
		return nil, err
	}
	m.End, err = mmVal(name, "end")
	if err != nil {
		return nil, err
	}
	s, err := ioutil.ReadFile(filepath.Join(name, "type"))
	if err != nil {
		return nil, err
	}
	if t, ok := Types[string(s)[0:len(s)-1]]; ok {
		m.Type = t
	}
	return m, nil
}

// This stringer prints out what you'd get from
// cat /sys/firmware/memmap/*/{start, end, type}
// and hence is useful for testing.
func (m *MemoryRange) String() string {
	var s string
	s = s + fmt.Sprintf("0x%x\n0x%x\n%v\n", m.Start, m.End, Names[m.Type])
	return s
}

// MemMap reads /sys/firmware/*/files and returns an array of MemoryRange structs or an error.
// The question of error handling is a bit messy, but for now we've decided the if even one
// map entry is readable, that will be enough, and that we'll skip ones we can't read. It's
// unlikely that any such errors will happen, however.
func Ranges() ([]MemoryRange, error) {
	var mr []MemoryRange
	// Walk the /sys/firmware/memmap tree. For each directory, read the start, end, and type entries.
	err := filepath.Walk("/sys/firmware/memmap", func(name string, fi os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("%v: %v\n", name, err)
			return err
		}
		// this should never happen, unless they add weird non-directory things in the future.
		if !fi.IsDir() || !unicode.IsDigit(rune(fi.Name()[0])) {
			return nil
		}
		m, err := mmr(name)
		if err == nil {
			mr = append(mr, *m)
			return filepath.SkipDir
		}
		return err
	})

	return mr, err
}
