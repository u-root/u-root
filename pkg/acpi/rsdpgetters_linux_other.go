// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux && !amd64 && !386

package acpi

// On non-x86 architectures, there are no legacy PC BIOS memory locations (e.g. EBDA, 0xe0000).
// Attempting to read these unbacked physical addresses via /dev/mem can cause fatal hardware faults.
// Thus, we only query the EFI system table.
var rsdpgetters = []func() (*RSDP, error){GetRSDPEFI}
