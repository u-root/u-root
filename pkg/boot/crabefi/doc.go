// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package crabefi loads the x86-64 CrabEFI coreboot payload from Linux.
// It is supported on amd64 Linux with the standard Go toolchain.
//
// # Entry ABI
//
// The payload is a little-endian, identity-linked ELF64 executable, but its
// ELF entry point executes in 32-bit protected mode. CrabEFI's _start reads
// the coreboot table pointer from [esp+4], builds identity page tables, and
// enters long mode itself. ELFCLASS64 describes the file, not the CPU mode
// required at its entry. The internal long_mode_start label is not an
// independent boot ABI: it depends on the page tables, GDT and registers set
// up by _start. See CrabEFI's crabefi-coreboot/src/arch/x86_64/entry.rs and
// crabefi-coreboot/x86_64-coreboot.ld.
//
// PrepareLoad reuses u-root's multiboot trampoline to leave long mode, then
// adapts its EBX argument to coreboot's cdecl stack argument. Payload segments
// and the handoff arena must therefore fit below 4 GiB. The arena contains
// an entry adapter, one trampoline page, two temporary stack pages, and
// coreboot tables with relocated SMBIOS data. It is reserved in the memory
// map passed to CrabEFI; no original CBMEM area is reused or advertised.
//
// Coreboot structures follow commonlib/coreboot_tables.h and
// commonlib/bsd/cbmem_id.h. The 64-memory-region limit and 64-GiB framebuffer
// mapping limit are CrabEFI implementation constraints, not coreboot ABI
// requirements. This loader does not implement a generic coreboot payload ABI.
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
