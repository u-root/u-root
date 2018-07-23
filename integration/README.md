# Integration Tests

This tests core use cases for u-root such as:

- retrieving and kexec'ing a Linux kernel,
- uinit (user init), and
- running unit tests requiring root privileges.

## Usage

Run the tests with:

    go test

## Requirements

- QEMU
  - Path and arguments must be set with `UROOT_QEMU`.
  - Example: `export UROOT_QEMU="$HOME/bin/qemu-system-x86_64 -L ."`
- Linux kernel
  - Path and arguments must be set with `UROOT_KERNEL`.
  - Example: `export UROOT_KERNEL="$HOME/linux/arch/x86_64/boot/bzImage"`

## To Dos

1. Support testing on architectures besides amd64. This is currently limited by
   the Linux bzImage uploaded.
