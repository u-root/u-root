// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build tinygo

package cpuid

/*
#include <stdlib.h>
char *CpuidVendor() {
    unsigned int eax, ebx, ecx, edx;
    char *vendor = (char *)(calloc(13, sizeof(char)));

    // handle this in golang
    if (vendor == NULL) {
        return NULL;
    }

    eax = 0;
    __asm__ volatile(
        "cpuid"
        : "=b"(ebx), "=d"(edx), "=c"(ecx)
        : "a"(eax)
    );

    *((unsigned int*) &vendor[0]) = ebx;
    *((unsigned int*) &vendor[4]) = edx;
    *((unsigned int*) &vendor[8]) = ecx;
    vendor[12] = 0x00;

    return vendor;
}
*/
import "C"

import (
	"fmt"
	"unsafe"
)

const (
	ManufacturerIDAMD   = "AuthenticAMD"
	ManufacturerIDIntel = "GenuineIntel"
)

// Get the CPU Identification String and return it
func CPUManufacturerID() (string, error) {
	vendor := C.CpuidVendor()

	if vendor == nil {
		return "", fmt.Errorf("error allocating memory in CGO")
	}

	goVendor := C.GoString(vendor)
	C.free(unsafe.Pointer(vendor))
	return goVendor, nil
}
