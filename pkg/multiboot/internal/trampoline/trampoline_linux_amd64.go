// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Trampoline sets machine to a specific state defined
// by multiboot v1 spec and boots the final kernel.
// https://www.gnu.org/software/grub/manual/multiboot/multiboot.html#Machine-state.
package trampoline

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/u-root/u-root/pkg/ubinary"
)

const (
	trampolineBeginLow  = "u-root-trampoline"
	trampolineBeginHigh = "-begin"
	trampolineEnd       = "u-root-trampoline-end"

	trampolineEntry = "u-root-entry-long"
	trampolineInfo  = "u-root-info-long"
)

var trampolineBegin []byte

var alwaysFalse bool

func start()

func init() {
	// Cannot use "u-root-trampoline-begin" directly, because then there
	// would be two locations of "u-root-trampoline-begin" sequence.
	trampolineBegin = append([]byte(trampolineBeginLow), []byte(trampolineBeginHigh)...)

	if alwaysFalse {
		// Can never happen, but still need it for linker to include assembly to a binary.
		start()
	}

	// Hope Go compiler will never get smart enough to optimize this function.
}

// alignUp aligns x to a 0x10 bytes boundary.
// go compiler aligns TEXT parts at 0x10 bytes boundary.
func alignUp(x int) int {
	const mask = 0x10 - 1
	return (x + mask) & ^mask
}

// Setup scans file for trampoline code and sets
// values for multiboot info address and kernel entry point.
func Setup(path string, infoAddr, entryPoint uintptr) ([]byte, error) {
	d, err := extract(path)
	if err != nil {
		return nil, err
	}
	return patch(d, infoAddr, entryPoint)
}

// extract extracts trampoline segment from file.
// trampoline segment begins after "u-root-trampoline-begin" byte sequence + padding,
// and ends at "u-root-trampoline-end" byte sequence.
func extract(path string) ([]byte, error) {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read trampoline file: %v", err)
	}

	begin := bytes.Index(d, trampolineBegin)
	if begin == -1 {
		return nil, io.ErrUnexpectedEOF
	}
	if begin = alignUp(begin + len(trampolineBegin)); begin > len(d) {
		return nil, io.ErrUnexpectedEOF
	}
	if ind := bytes.Index(d[begin:], trampolineBegin); ind != -1 {
		return nil, fmt.Errorf("multiple definition of %q label", trampolineBegin)
	}

	end := bytes.Index(d[begin:], []byte(trampolineEnd))
	if end == -1 {
		return nil, io.ErrUnexpectedEOF
	}
	return d[begin : begin+end], nil
}

// patch patches the trampoline code to store value for multiboot info address
// after "u-root-header-long" byte sequence + padding and value
// for kernel entry point, after "u-root-entry-long" byte sequence + padding.
func patch(trampoline []byte, infoAddr, entryPoint uintptr) ([]byte, error) {
	replace := func(d, label []byte, val uint32) error {
		buf := make([]byte, 4)
		ubinary.NativeEndian.PutUint32(buf, val)

		ind := bytes.Index(d, label)
		if ind == -1 {
			return fmt.Errorf("%q label not found in file", label)
		}
		ind = alignUp(ind + len(label))
		if len(d) < ind+len(buf) {
			return io.ErrUnexpectedEOF
		}
		copy(d[ind:], buf)
		return nil
	}

	if err := replace(trampoline, []byte(trampolineInfo), uint32(infoAddr)); err != nil {
		return nil, err
	}
	if err := replace(trampoline, []byte(trampolineEntry), uint32(entryPoint)); err != nil {
		return nil, err
	}
	return trampoline, nil
}
