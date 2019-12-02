# Integration Tests

These are VM based tests for core u-root functionality such as:

-   retrieving and kexec'ing a Linux kernel,
-   DHCP client tests,
-   uinit (user init), and
-   running unit tests requiring root privileges.

To learn more about how these tests work under the hood, see the next section,
otherwise jump ahead to the sections on how to write and run these tests.

## VM Testing Infrastructure

Our VM testing infrastructure starts a QEMU virtual machine that boots with
our given kernel and initramfs, and runs the uinit or commands that we want to
test.

The test architecture, kernel and QEMU binary are set using environment
variables.

Testing mainly relies on 2 packages: [pkg/vmtest](/pkg/vmtest) and
[pkg/qemu](/pkg/qemu).

pkg/vmtest takes in integration test options, and given those and the
environment variables, uses pkg/qemu to start a QEMU VM with the correct command
line and configuration.

There are a couple of ways to test:
* Custom initramfs: provide an initramfs that will be used in the VM.
* Custom uinit: provide a uinit. The testing setup will generate an initramfs 
  with that uinit.
* Test commands: provide the set of commands to be tested. The testing setup
  will generate an initramfs that runs those commands.

Files that need to be shared with the VM are written to a temp dir which is
exposed as a Plan 9 (9p) filesystem in the VM.

To check for the correct behavior, we use the go expect package to find
expected output in QEMU's serial output within a given timeout.

## Running Tests

These tests only run on Linux on amd64 and arm.

1. Set Environment variables:

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


Our automated CI uses Dockerfiles to build a kernel and QEMU and set these
environment variables. You can see the Dockerfile and the config file used to
build the kernel for each supported architecture [here](/.circleci/images).

2. Run the tests with:

```sh
go test [-v]
```

The verbose flag is useful to see the QEMU command line being used and the full
serial output. It is also useful to see which tests are being skipped and why
(particularly for ARM, where many tests are currently skipped).

Unless you want to wait a long time for all tests to complete, run just the
specific test you want, e.g.

```sh
go test [-v] -test.run=TestDhclient
```

## Writing a New Test

To write a new test, first decide which of the options from the previous
section best fit your case (custom initramfs, custom uinit, test commands).

`vmtest.QEMUTest` is the function that starts the QEMU VM and returns the VM
struct. There, provide the test options for your use case.

The VM struct returned by `vmtest.QEMUTest` represents a running QEMU virtual
machine. Use its family of Expect methods to check for the correct result.

