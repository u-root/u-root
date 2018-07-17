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
  - Path and arguments must be set with `TEST_QEMU`.
  - Example: `export TEST_QEMU="$USER/bin/qemu-system-x86_64" -L .`

## To Dos

1. Support testing on architectures besides amd64. This is currently limited by
   the Linux bzImage uploaded.
