// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (!tinygo || tinygo.enable) && amd64 && linux

// The System Management Network (SMN, try to say it fast)
// is a parallel universe address space on newer AMD64 cpus.
// The address is 32 bits, as is data.
// Unfortunately it is only accessible via the classic
// index/register pair. Fortunately that pair is accessible
// in mmconfig space.
package main

import (
	"github.com/u-root/u-root/pkg/memio"
)

// This is a const at present but there are no guarantees.
const (
	pciBaseAddress = 0xe0000000
	usageSMM       = `io rs index # read from system management network on newer AMD CPUs.
io ws index value # write value to system management network on newer AMD CPUs.
`
)

func init() {
	usageMsg += usageSMM
	addCmd(readCmds, "rs", &cmd{smnRead, 32, 32})
	addCmd(writeCmds, "ws", &cmd{smnWrite, 32, 32})
}

func do(addr int64, data memio.UintN, op func(int64, memio.UintN) error) error {
	a := newInt(uint64(addr), 32)
	if err := memio.Write(pciBaseAddress+0xb8, a); err != nil {
		return err
	}
	return op(pciBaseAddress+0xbc, data)
}

func smnWrite(addr int64, data memio.UintN) error {
	return do(addr, data, memio.Write)
}

func smnRead(addr int64, data memio.UintN) error {
	return do(addr, data, memio.Read)
}
