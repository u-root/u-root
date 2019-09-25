# Integration Tests

This tests core use cases for u-root such as:

-   retrieving and kexec'ing a Linux kernel,
-   DHCP client tests,
-   uinit (user init), and
-   running unit tests requiring root privileges.

## Requirements

These tests only run on Linux on amd64 and arm.

Environment variables:

-   `UROOT_QEMU` points to a QEMU binary and args, e.g.

```sh
export UROOT_QEMU="$HOME/bin/qemu-system-x86_64 -enable-kvm"
```

-   `UROOT_KERNEL` points to a Linux kernel binary, e.g.

```sh
export UROOT_KERNEL="$HOME/linux/arch/x86/boot/bzImage"
```

-   (optional) `UROOT_TESTARCH` (defaults to host architecture) is the
    architecture to test. Only `arm` and `amd64` are supported.

-   (optional) `UROOT_QEMU_TIMEOUT_X` (defaults to 1.0) can be used to multiply
    the timeouts for each test in case QEMU on your machine is slower. For
    example, if you cannot turn on `-enable-kvm`, use `UROOT_QEMU_TIMEOUT_X=2`
    as our test automation does.

The kernel used in our automated CI is a `v4.17` kernel built from
[this .config file](/.circleci/images/test-image-amd64/config_linux4.17_x86_64.txt).
You can look at the [Dockerfile](/.circleci/images/test-image-amd64/Dockerfile)
to see what exactly is needed to build it.

## Usage

Run the tests with:

```sh
go test
```

Unless you want to wait a long time for all tests to complete, run just the test
you need with

```sh
go test -test.run=TestDhclient
```
