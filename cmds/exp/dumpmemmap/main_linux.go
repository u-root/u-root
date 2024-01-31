// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command dumpmemmap prints different kernel interpretations of physical
// memory address space.
//
// Support for:
//
//   - /proc/iomem (exists on all systems)
//   - /sys/firmware/memmap (exists on x86 systems)
//   - /sys/kernel/debug/memblock (exists on systems with CONFIG_ARCH_KEEP_MEMBLOCK, in particular arm64)
//   - /sys/firmware/fdt (exists on systems with device trees)
package main

import (
	"fmt"
	"log"

	"github.com/dustin/go-humanize"
	"github.com/u-root/u-root/pkg/boot/kexec"
	"github.com/u-root/u-root/pkg/dt"
)

func printMM(mm kexec.MemoryMap) {
	for _, r := range mm {
		fmt.Println(" ", r, " ", humanize.Bytes(uint64(r.Range.Size)))
	}
}

func main() {
	memmap, err := kexec.MemoryMapFromSysfsMemmap()
	if err != nil {
		log.Printf("/sys/firmware/memmap: %v", err)
	} else {
		fmt.Println("/sys/firmware/memmap:")
		printMM(memmap)
	}

	memblock, err := kexec.MemoryMapFromMemblock()
	if err != nil {
		log.Printf("/sys/kernel/debug/memblock: %v", err)
	} else {
		fmt.Println("/sys/kernel/debug/memblock:")
		printMM(memblock)
	}

	iomem, err := kexec.MemoryMapFromIOMem()
	if err != nil {
		log.Printf("/proc/iomem: %v", err)
	} else {
		fmt.Println("/proc/iomem:")
		printMM(iomem)
	}

	fdt, err := dt.LoadFDT(nil)
	if err != nil {
		log.Printf("loadFDT: %v", err)
		return
	}

	// Prepare segments.
	fdtMap, err := kexec.MemoryMapFromFDT(fdt)
	if err != nil {
		log.Printf("MemoryMapFromFDT(%v): %v", fdt, err)
		return
	}

	fmt.Println("/sys/firmware/fdt:")
	printMM(fdtMap)
}
