// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux && (amd64 || 386)

package acpi

// On x86 architectures, we probe legacy PC BIOS memory locations first, then EFI.
// You can change the getters if you wish for testing.
var rsdpgetters = []func() (*RSDP, error){GetRSDPEBDA, GetRSDPMem, GetRSDPEFI}
