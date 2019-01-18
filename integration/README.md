# Integration Tests

This tests core use cases for u-root such as:

- retrieving and kexec'ing a Linux kernel,
- uinit (user init), and
- running unit tests requiring root privileges.

## Usage

Run the tests with:

    go test

When the QEMU arch is not amd64, set the `UROOT_TESTARCH` variable. For
example:

    UROOT_TESTARCH=arm go test

Currently, only amd64 and arm are supported.

## Requirements

- QEMU
  - Path and arguments must be set with `UROOT_QEMU`.
  - Example: `export UROOT_QEMU="$HOME/bin/qemu-system-x86_64 -L ."`
- Linux kernel
  - Path and arguments must be set with `UROOT_KERNEL`.
  - Example: `export UROOT_KERNEL="$HOME/linux/arch/x86_64/boot/bzImage"`
