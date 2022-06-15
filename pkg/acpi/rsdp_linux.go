// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package acpi

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// GetRSDPEFI finds the RSDP in the EFI System Table.
func GetRSDPEFI() (*RSDP, error) {
	file, err := os.Open("/sys/firmware/efi/systab")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	const (
		acpi20 = "ACPI20="
		acpi   = "ACPI="
	)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		start := ""
		if strings.HasPrefix(line, acpi20) {
			start = strings.TrimPrefix(line, acpi20)
		}
		if strings.HasPrefix(line, acpi) {
			start = strings.TrimPrefix(line, acpi)
		}
		if start == "" {
			continue
		}
		base, err := strconv.ParseInt(start, 0, 63)
		if err != nil {
			continue
		}
		rsdp, err := readRSDP(base)
		if err != nil {
			continue
		}
		return rsdp, nil
	}
	if err := scanner.Err(); err != nil {
		log.Printf("error while reading EFI systab: %v", err)
	}
	return nil, fmt.Errorf("invalid /sys/firmware/efi/systab file")
}

// You can change the getters if you wish for testing.
var rsdpgetters = []func() (*RSDP, error){GetRSDPEBDA, GetRSDPMem, GetRSDPEFI}
