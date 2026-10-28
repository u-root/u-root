# CrabEFI from LinuxBoot

This loader starts CrabEFI's **64-bit-only amd64 payload**, not a PE/COFF EFI
application. It needs Linux and the standard Go toolchain; TinyGo and other
architectures return an unsupported-platform error.

## Entry ABI and payload build

The loader requires version 1 of CrabEFI's 64-bit entry ABI, identified by a
`CrabEFI` ELF note in a `PT_NOTE` program header. It passes the coreboot table
pointer in RDI and jumps to the ELF entry without leaving long mode. CrabEFI
owns its page tables, GDT and stack. There is **no 32-bit fallback**, Multiboot
adapter, external launcher or kexec-tools dependency. Old ELF64 payloads with
a 32-bit entry are intentionally rejected; rebuild the payload and u-root
together. See [doc.go](doc.go) for the exact contract.

The coordinated CrabEFI revision is
`e9a6da4ee8feb5497729a238a4e8a40537c4b9cb` (`x86-entry64`), or a later revision
retaining that entry ABI. From that CrabEFI checkout, build with:

```sh
cargo build -p crabefi-coreboot --release --target x86_64-unknown-none
```

Alternatively, `./crabefi build --arch x86-64` also rebuilds and audits its
Runtime Services image. Bundle the resulting ELF as `/boot/crabefi.elf`:

```sh
u-root -uinitcmd='boot -v' -skip-ldd \
    -files /path/to/CrabEFI/target/x86_64-unknown-none/release/crabefi:boot/crabefi.elf \
    ./cmds/core/init ./cmds/core/gosh ./cmds/boot/boot
```

An absent default payload leaves the existing boot menu unchanged. To choose
another payload, use `boot -crabefi=/path/to/payload`; an empty value disables
the firmware entry. `boot -no-load` only displays the selection;
`boot -no-exec` stages it without leaving Linux.

## First-kernel requirements

The first Linux kernel must provide:

- `kexec_load` (`CONFIG_KEXEC`) and permission to use it. A kernel that only
  permits signature-verified `kexec_file_load` cannot use this ELF loader.
- The original firmware map at `/sys/firmware/memmap`
  (`CONFIG_FIRMWARE_MEMMAP`), ACPI support, and readable firmware tables.
- DMI entry-point and table files at `/sys/firmware/dmi/tables`
  (`CONFIG_DMI`, `CONFIG_DMI_SYSFS`).
- `/dev/mem` (`CONFIG_DEVMEM`), with policy permitting ACPI reads.
- `/dev/fb0` backed by a physical, packed 32-bit true-color framebuffer with
  8-bit RGB channels, no panning, and an address below 64 GiB. Coreboot
  simplefb is suitable (`CONFIG_FB`, `CONFIG_FB_SIMPLE`, and the corresponding
  platform framebuffer setup); a shadow framebuffer with no physical address
  is not. The framebuffer must not overlap reusable firmware RAM.

RSDP discovery uses u-root's existing EFI/legacy ACPI discovery first. On
coreboot systems outside those discovery paths, the fallback needs the kernel
log's `ACPI: RSDP 0x...` message to remain available and readable. It does not
consume or clear the log. Missing firmware interfaces cause loading to fail;
the loader does not invent replacement ACPI or SMBIOS tables.

## Validation and limits

Host tests cover ELF ABI-note validation (including rejection of old payloads),
full-width entry operands, payload and handoff placement above 4 GiB, firmware
memory validation, coreboot table layout, SMBIOS relocation, signed read bounds,
framebuffer formats, and menu integration:

```sh
go test ./pkg/boot/crabefi
go test ./pkg/boot -run CrabEFI
go test ./cmds/boot/boot
```

End-to-end validation used QEMU Q35/KVM with 4 GiB RAM and NVMe storage:
coreboot -> Linux 7.2.0-rc6 -> u-root's normal boot menu -> CrabEFI's 64-bit
entry -> installed Windows 11 Pro 25H2 (build 26200.6584). The disposable test
disk's success and diagnostic files were deleted before each cold boot. The
full OS regenerated its success marker, reported `SYSTEMDRIVE=C:`, and recorded
`explorer.exe` running. This was not a WinPE boot or a new Windows installation.

These results do not establish hardware-wide compatibility. CrabEFI must
support the relevant hardware with its own drivers after Linux exits. When
reporting another boot, include payload/kernel revisions and build options,
board or QEMU configuration, storage type, framebuffer information, and logs.

This handoff does not authenticate or measure the first LinuxBoot stages,
carry a TPM event log, or supply persistent firmware-variable storage. Booting
Windows does not by itself establish Secure Boot or BitLocker compatibility.
