# Integration Tests

These are VM based tests for core u-root functionality such as:

-   retrieving and kexec'ing a Linux kernel,
-   DHCP client tests,
-   uinit (user init), and
-   running unit tests requiring root privileges.

## Overview

All tests are in the integration/ directory. Within that, there are a
few subdirectories:

* generic-tests/ : most tests can be put under this.
* golang-tests/ : this is for Go unit tests that can be run inside the VM.
* testcmd/ : this contains custom uinits for tests.
* testdata/ : this contains any extra files for tests.

### Generic Tests

This is where most tests live. main\_test.go is the test orchestrator. It is 
responsible for running all the tests in the directory.

All tests must have an entry in the tests table in main\_test.go. Each test 
entry has 2 things: a name, and a runner function. This runner function is where
the actual test is taking place.

main\_test.go will then go through the table and call the runner function
for each test.

When main\_test.go is invoked using `go test`, it takes a few flags:
* kernel: path to the Linux kernel binary to use. eg. `-kernel="$HOME/linux/arch/x86/boot/bzImage"`
* qemu: path to the QEMU binary and args to use. eg. `-qemu="$HOME/bin/qemu-system-x86_64 -enable-kvm"`
* [optional] initramfs: path to an initramfs to use for all the tests. If one is not
  provided, there is a default generic initramfs that includes all cmds and simply runs the test
  commands given (using this is preferred).
* [optional] testarch: test architecture to use (amd64 or arm).

### Golang Tests 

This is for running Go unit tests from all u-root packages in the VM.
gotest\_test.go finds the tests from all packages (except those that are
blacklisted) and executes the test.

## VM Testing Infrastructure

This is a brief look at what's happening under the hood.

Our VM testing infrastructure starts a QEMU virtual machine that boots with
our given kernel and initramfs, and runs the uinit or commands that we want to
test.

The test architecture, kernel and QEMU binary can be set using environment
variables, or flags.

Testing mainly relies on 2 packages: [pkg/vmtest](/pkg/vmtest) and
[pkg/qemu](/pkg/qemu).

pkg/vmtest takes in integration test options, and given those and the
environment variables, uses pkg/qemu to start a QEMU VM with the correct command
line and configuration.

Files that need to be shared with the VM are written to a temp dir which is
exposed as a Plan 9 (9p) filesystem in the VM.

To check for the correct behavior, we use the go expect package to find
expected output in QEMU's serial output within a given timeout.

## Running Tests

### Generic Tests

To run all tests:

```sh
cd generic-tests/
go test [-v] -kernel </path/to/kernel> -qemu </path/to/qemu>
```

You have an option to specify the test architecure (as of now, amd64 and arm are
supported. defaults to host architecture), and custom initramfs as
additional flags.

The verbose flag is useful to see the QEMU command line being used and the full
serial output. It is also useful to see which tests are being skipped and why
(particularly for ARM, where many tests are currently skipped).

Unless you want to wait a long time for all tests to complete, run just the
specific test you want.

To run a specific test:

```sh
go test [-v] -kernel </path/to/kernel> -qemu </path/to/qemu> -test.run=TestGeneric/NAME
```

Here, NAME is the name from the test entry in main\_test.go.

### Golang Tests 

```sh
cd golang-tests/
go test [-v]
```

## Writing Tests

###  Generic Tests

1. Write a test runner function.

This is of the type `func(t *testing.T, initramfs string)`, where initramfs is
the optional flag to test orchestrator. In the runner, you have the option to
override this initramfs, or have another default besides the generic. (see
dhclient\_test.go for examples).

  a. Call `vmtest.QEMUTest` to start the VM.

  `vmtest.QEMUTest` starts the QEMU VM and returns the VM struct. There, provide
  the test options for your use case. You **must** pass in the Initramfs option to
  use a non-generic initramfs (either caller-provided, or your own). The
  TestCmds option is for specifying the actual commands you want to test. 

  b. Use Expect to check for the correct results.

  The VM struct returned by `vmtest.QEMUTest` represents a running QEMU virtual
  machine. Use its family of Expect methods to check for the correct result.

2. Add the test entry to main\_test.go.


### Golang Tests

Write Go unit tests in the package they are meant to test.

