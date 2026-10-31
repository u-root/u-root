// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package crabefi loads the x86-64 CrabEFI payload from Linux.
// It is supported on amd64 Linux with the standard Go toolchain.
//
// # Entry ABI
//
// Only version 1 of CrabEFI's 64-bit entry ABI is supported. The payload must
// be an identity-linked, little-endian ELF64 executable with a PT_NOTE owned
// by "CrabEFI\x00", type 1, and a four-byte little-endian version number of 1.
// CrabEFI's default coreboot build has a 32-bit entry and is not suitable for
// this loader. Build the payload with --features entry64 or ./crabefi build
// --entry64. There is no 32-bit fallback or use of the multiboot trampoline.
//
// Linux's x86-64 kexec hands over in long mode with identity-mapped segments.
// A small entry stub passes the physical coreboot table pointer in RDI and
// jumps to the ELF entry, using full-width addresses throughout. CrabEFI
// installs its own GDT, stack and page tables without leaving long mode.
// Four- and five-level paging are supported; payload and handoff storage must
// fit within CrabEFI's initial 64-GiB identity map. SMBIOS 2 still requires its
// relocated table below 4 GiB, as required by that table format, not the entry ABI.
//
// Coreboot structures follow commonlib/coreboot_tables.h and
// commonlib/bsd/cbmem_id.h. The 64-memory-region limit is a CrabEFI parser
// constraint, not a coreboot ABI requirement. The arena contains the entry
// code, coreboot tables and relocated SMBIOS data. It is reserved in the
// memory map passed to CrabEFI; no original CBMEM area is advertised.
//
// # Firmware and security
//
// Firmware RAM is page-aligned inward; other ranges are aligned outward.
// Overlapping resulting ranges are rejected rather than guessing ownership.
// ACPI remains in preserved firmware memory, while SMBIOS is copied out of
// sysfs and its entry point is relocated using package smbios.
//
// The payload and first LinuxBoot stages must be trusted by the caller.
// Authentication, measurement, TPM event-log handoff and persistent firmware
// variables are not provided by this package. PrepareLoad only stages data;
// the caller remains responsible for kexec_load, cleanup and execution.
package crabefi
