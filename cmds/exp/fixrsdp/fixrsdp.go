// Copyright 2019-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// fixrsdp copies the existing RSDP into the EBDA region in low mem.
//
// This is because LinuxBoot tends to be EFI booted, which places the RSDP
// outside of the low 1M or the EBDA. If Linuxboot legacy boots the following
// operating systems, such as with kexec, they may not have a way to find the
// RSDP afterwards.  All u-root commands that open /dev/mem should also flock
// it to ensure safe, sequential access.
package main

import (
	"bytes"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/acpi"
	"github.com/u-root/u-root/pkg/boot/ebda"
)

func main() {
	// Find the RSDP.
	r, err := acpi.GetRSDP()
	if err != nil {
		log.Fatalf("Unable to find system RSDP, got: %v", err)
	}

	rData := r.AllData()
	rLen := len(rData)

	base := r.RSDPAddr()
	// Check if ACPI rsdp is already in low memory
	if base >= 0xe0000 && base+int64(rLen) < 0xffff0 {
		log.Printf("RSDP is already in low memory at %#X, no need to fix.", base)
		return
	}

	// Find the EBDA
	f, err := os.OpenFile("/dev/mem", os.O_RDWR, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	e, err := ebda.ReadEBDA(f)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("EBDA starts at %#X, length %#X bytes", e.BaseOffset, e.Length)

	// Scan low 1K of EBDA for an empty spot that is 16 byte aligned
	emptyStart := 0
	for i := 16; i < (1024-rLen) && i < (int(e.Length)-rLen); i += 16 {
		// Check if there's an empty spot to put the RSDP table.
		if bytes.Equal(e.Data[i:i+rLen], make([]byte, rLen)) {
			emptyStart = i
			log.Printf("Found empty space at %#X offset into EBDA, will copy RSDP there", emptyStart)
			break
		}
	}

	if emptyStart == 0 {
		log.Fatal("Unable to find empty space to put RSDP")
	}

	copy(e.Data[emptyStart:emptyStart+rLen], rData)

	if err = ebda.WriteEBDA(e, f); err != nil {
		log.Fatal(err)
	}
	// Verify write, depending on the kernel settings like CONFIG_STRICT_DEVMEM, writes can silently fail.
	v, err := ebda.ReadEBDA(f)
	if err != nil {
		log.Fatalf("Error reading EBDA: %v", err)
	}
	res := bytes.Compare(e.Data, v.Data)
	if res != 0 {
		log.Fatal("Write verification failed !")
	}
}
