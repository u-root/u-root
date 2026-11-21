# CrabEFI from LinuxBoot

This loader starts CrabEFI's **opt-in `entry64` amd64 payload**, not a PE/COFF EFI
application. It needs Linux and the standard Go toolchain; TinyGo and other
architectures return an unsupported-platform error.

## Entry ABI and payload build

The loader requires version 1 of CrabEFI's 64-bit entry ABI, identified by a
`CrabEFI` ELF note in a `PT_NOTE` program header. It passes the coreboot table
pointer in RDI and jumps to the ELF entry without leaving long mode. CrabEFI
owns its page tables, GDT and stack. There is **no 32-bit fallback**, Multiboot
adapter, external launcher or kexec-tools dependency in this loader. CrabEFI's
**default build keeps its 32-bit coreboot entry**, which remains supported for
unmodified coreboot but is not suitable for u-root. See [doc.go](doc.go) for the
64-bit caller contract.

The 64-bit entry was added by
[CrabEFI#222](https://github.com/ArthurHeymans/CrabEFI/pull/222) (merge commit
`1cfd94fdb57ac6ce19dd70d86bed7680a02ad425`). Use that or a later revision, and
explicitly select its `entry64` feature:

```sh
./crabefi build --arch x86-64 --entry64
```

This rebuilds/audits Runtime Services and writes
`target/entry64/x86_64-unknown-none/release/crabefi`, separately from the default
coreboot payload. Alternatively, use the checked-in runtime bundle with Cargo:

```sh
cargo build -p crabefi-coreboot --release --target x86_64-unknown-none --features entry64
```

Direct Cargo builds write the usual `target/x86_64-unknown-none/release/crabefi`
path; do not confuse that explicitly flagged output with a default coreboot
build. Bundle the wrapper's entry64 ELF as `/boot/crabefi.elf`:

```sh
u-root -uinitcmd='boot -v' -skip-ldd \
    -files /path/to/CrabEFI/target/entry64/x86_64-unknown-none/release/crabefi:boot/crabefi.elf \
    ./cmds/core/init ./cmds/core/gosh ./cmds/boot/boot
```

An absent default payload leaves the existing boot menu unchanged. To choose
another payload, use `boot -crabefi=/path/to/payload`; an empty value disables
the firmware entry. `boot -no-load` only lists the entries; `boot -no-exec`
stages the selected one without leaving Linux.

## First-kernel requirements

The first Linux kernel must provide:

- `kexec_load` (`CONFIG_KEXEC`) and permission to use it. A kernel that only
  permits signature-verified `kexec_file_load` cannot use this ELF loader.
- The original firmware map at `/sys/firmware/memmap`
  (`CONFIG_FIRMWARE_MEMMAP`), ACPI support, and readable firmware tables.
- DMI entry-point and table files at `/sys/firmware/dmi/tables`
  (`CONFIG_DMI`, `CONFIG_DMI_SYSFS`).
- `/dev/mem` (`CONFIG_DEVMEM`), with policy permitting ACPI reads.

The framebuffer is optional. CrabEFI has no display drivers of its own; its
Graphics Output Protocol (GOP) and graphical boot menu use the framebuffer
Linux inherited from firmware. If `/dev/fb0` does not exist, the loader logs a
warning and omits the framebuffer record. The handoff carries no serial console
record either, so CrabEFI then has no console and boots its default entry after
the menu timeout.

When `/dev/fb0` exists, it must be a physical, packed 32-bit true-color linear
framebuffer with 8-bit RGB channels, no panning, and an address below 64 GiB
that does not overlap reusable RAM. Otherwise loading fails rather than
silently dropping the display. Any fbdev driver exposing the firmware
framebuffer with its physical address is suitable, such as `simplefb`, `efifb`
or `vesafb`; DRM fbdev emulation using a shadow buffer without a physical
address is not.
`/dev/fb0` also needs `CONFIG_FB_DEVICE`.

On coreboot, the tested configuration uses coreboot's framebuffer driver:
`CONFIG_GOOGLE_FIRMWARE`, `CONFIG_GOOGLE_COREBOOT_TABLE`,
`CONFIG_GOOGLE_FRAMEBUFFER_COREBOOT`, `CONFIG_FB_SIMPLE` and
`CONFIG_FB_DEVICE`, with `CONFIG_SYSFB_SIMPLEFB` disabled. With
`CONFIG_SYSFB_SIMPLEFB`, a Linux payload added by cbfstool can end up with an
`efi-framebuffer` device, which needs `efifb` and thus `CONFIG_EFI`.

RSDP discovery uses u-root's existing EFI/legacy ACPI discovery first. On
coreboot systems outside those discovery paths, the fallback needs the kernel
log's `ACPI: RSDP 0x...` message to remain available and readable. It does not
consume or clear the log. The RSDP and its XSDT must be below 64 GiB. Other
missing firmware interfaces cause loading to fail; the loader does not invent
replacement ACPI or SMBIOS tables.

## Validation and limits

Host tests cover the ELF and ABI-note checks, handoff placement, firmware
memory and ACPI validation, coreboot table layout, SMBIOS relocation,
framebuffer formats and menu integration:

```sh
go test ./pkg/boot/crabefi
go test ./pkg/boot -run CrabEFI
go test ./cmds/boot/boot
```

In QEMU Q35/KVM (coreboot, Linux 7.2-rc6, CrabEFI `4e3a3e24` built with
`--entry64`), this path deployed Windows 11 Pro 25H2 from WinPE, then booted
the installed OS to its desktop, with and without a framebuffer.

These results do not establish hardware-wide compatibility. CrabEFI must
support the relevant hardware with its own drivers after Linux exits. When
reporting another boot, include payload/kernel revisions and build options,
board or QEMU configuration, storage type, framebuffer information, and logs.

This handoff does not authenticate or measure the first LinuxBoot stages,
carry a TPM event log, or supply persistent firmware-variable storage. Booting
Windows does not by itself establish Secure Boot or BitLocker compatibility.
