# CrabEFI from LinuxBoot

This loader starts CrabEFI's **amd64 coreboot payload**, not a PE/COFF EFI
application. It needs Linux and the standard Go toolchain; TinyGo and other
architectures return an unsupported-platform error.

## Building and bundling the payload

The entry contract is described in [doc.go](doc.go). It was checked against
[CrabEFI bb83e6eb](https://github.com/ArthurHeymans/CrabEFI/tree/bb83e6eb9a79a99b0be76b13806c4dc2c40de4e4).
For that revision, the payload can be built with:

```sh
git clone https://github.com/ArthurHeymans/CrabEFI
cd CrabEFI
git checkout bb83e6eb9a79a99b0be76b13806c4dc2c40de4e4
# Use CrabEFI's documented Rust toolchain/development environment.
cargo build -p crabefi-coreboot --release --target x86_64-unknown-none
```

Bundle `target/x86_64-unknown-none/release/crabefi` in the u-root initramfs as
`/boot/crabefi.elf`:

```sh
u-root -uinitcmd=boot -skip-ldd \
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

Host tests exercise ELF validation, allocation (including the last arena below
4 GiB), coreboot table layout, SMBIOS relocation, framebuffer formats, signed
physical-read bounds, and menu integration:

```sh
go test ./pkg/boot/crabefi
go test ./pkg/boot -run CrabEFI
go test ./cmds/boot/boot
```

These tests do **not** establish that a particular CrabEFI revision boots
Windows on a particular machine. When reporting a successful end-to-end boot,
include the CrabEFI revision/build options, kernel revision/configuration,
board or QEMU configuration, storage/controller type, framebuffer information,
and serial log through the EFI boot manager. CrabEFI must support the relevant
hardware with its own drivers after Linux exits.

This handoff does not authenticate or measure the first LinuxBoot stages,
carry a TPM event log, or supply persistent firmware-variable storage. Booting
Windows does not by itself establish Secure Boot or BitLocker compatibility.
