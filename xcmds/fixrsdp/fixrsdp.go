// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This command copies the existing RSDP into the EBDA. This is because
// LinuxBoot tends to be EFI booted, which places the RSDP outside of the
// low 1M or the EBDA. If Linuxboot legacy boots the following operating systems,
// they may not have a way to find the RSDP afterwards.

package main

import (
	"log"
	"os"

	"github.com/u-root/u-root/pkg/acpi"
	"github.com/u-root/u-root/pkg/ebda"
)

func main() {
	// Find the RSDP.
	base, r, err := acpi.GetRSDP()
	if err != nil {
		log.Fatalf("Unable to find system RSDP, got: %v", err)
	}

	// Check if ACPI rsdp is already in low memory
	if base >= 0xe0000 && base < 0xffff0 {
		log.Printf("rsdp is already in low memory at 0x%X, no need to fix.", base)
		os.Exit(0)
	}

	rData := r.AllData()
	rLen := len(rData)

	// Find the EBDA
	f, err := os.Open("/dev/mem")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	e, err := ebda.ReadEBDA(f)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("EBDA starts at 0x%X, length 0x%X bytes", e.BaseOffset, e.Length)

	// Scan low 1K of EBDA for an empty spot that is 16 byte aligned
	emptyStart := 0
	for i := 0; i < (1024-rLen) && i < (int(e.Length)-rLen); i += 16 {
		found := true
		for j := i; j < acpi.HeaderLength; j++ {
			if e.Buf[j] != 0 {
				found = false
				break
			}
		}
		if found {
			emptyStart = i
			log.Printf("Found empty space at 0x%X offset into EBDA, will copy RSDP there", emptyStart)
			break
		}
	}

	if emptyStart == 0 {
		log.Fatal("Unable to find empty space to put RSDP")
	}

	copy(e.Buf[emptyStart:emptyStart+rLen], rData)

	if err = ebda.WriteEBDA(e, f); err != nil {
		log.Fatal(err)
	}
}
