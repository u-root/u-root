# Integration Tests

These are VM based tests for core u-root functionality such as:

* retrieving and kexec'ing a Linux kernel,
* DHCP client tests,
* uinit (user init), and
* running unit tests requiring root privileges.

All tests are in the integration/ directory. Within that, there are a
few subdirectories:

* generic-tests/ : most tests can be put under this.
* gotests/ : this is for Go unit tests that can be run inside the VM.
* testcmd/ : this contains custom uinits for tests.
* testdata/ : this contains any extra files for tests.

To learn more about how these tests work under the hood, see the next section,
otherwise jump ahead to the sections on how to write and run these tests.

## VM Testing Infrastructure

Our VM testing infrastructure starts a QEMU virtual machine that boots with
our given kernel and initramfs, and runs the uinit or commands that we want to
test.

Testing mainly relies on [hugelgupf/vmtest](https://github.com/hugelgupf/vmtest)
API.

vmtest takes in integration test options, and given those and the environment
variables starts a QEMU VM with the correct command line and configuration.

E.g. files can be shared with the VM using the `vmtest.WithSharedDir` option,
which exposes them as a Plan 9 (9p) filesystem in the VM. This includes the
initramfs being used for the VM.

The test architecture, kernel and QEMU binary are set using environment
variables, `VMTEST_ARCH`, `VMTEST_KERNEL`, and `VMTEST_QEMU`.

The initramfs can come from the following sources:
* Default: u-root initramfs built with `gosh` and `init`. Go test
  configuration can request more commands and files to be added.
* User overridden: when the `VMTEST_INITRAMFS_OVERRIDE` environment variable is
  used to override the initramfs. The user is responsible for ensuring the
  initramfs contains the binaries required by the test.

To check for the correct behavior, we use expect scripting package to find
expected output in QEMU's serial output within a given timeout.

## Running Tests

Most tests run on amd64, arm64, or arm. All should be available on amd64.

To reproduce the same config that GitHub Actions use, the `runvmtest` tool will
set up `VMTEST_KERNEL` and `VMTEST_QEMU` for you.

```
go install github.com/hugelgupf/vmtest/tools/runvmtest@latest
runvmtest -- go test [-v] [pattern]
```

If necessary to reproduce, pass `--keep-artifacts` to runvmtest:

```
runvmtest --keep-artifacts -- go test [-v] [pattern]
```

To try a different guest architecture, use `VMTEST_ARCH` with GOARCH values:

```
VMTEST_ARCH=arm64 runvmtest -- go test [-v] [pattern]
```

### Supplying custom kernel or QEMU

`VMTEST_QEMU` points to a QEMU binary and args. In the following example,
runvmtest only sets up `VMTEST_KERNEL`:

```sh
VMTEST_QEMU="$HOME/bin/qemu-system-x86_64 -enable-kvm" runvmtest -- go test -v
```

`VMTEST_KERNEL` points to a Linux kernel binary, e.g.

```sh
VMTEST_KERNEL="$HOME/linux/arch/x86/boot/bzImage"
```

Other variables:

-   (optional) `VMTEST_INITRAMFS_OVERRIDE` is a custom initramfs to use for all
    tests. This will override all other initramfs options defined by the tests.

## Kernel Code Coverage

For kernel code coverage, build your kernel with the CONFIGs here:
https://www.kernel.org/doc/html/v4.14/dev-tools/gcov.html

With these configs enabled, a folder containing gcda files (raw coverage data)
will appear in the VM as /sys/kernel/debug/gcov/.

To collect kernel code coverage, set `VMTEST_KERNEL_COVERAGE_DIR` to a
directory. Per-test coverage will be saved as

```
${VMTEST_KERNEL_COVERAGE_DIR}/{{testname}}/{{instance}}/kernel_coverage.tar
```

## Writing a New Test

Take a look at the [vmtest
API](https://pkg.go.dev/github.com/hugelgupf/vmtest?utm_source=godoc).

Examples to look at in u-root:

* Running Go unit tests in a VM: [/pkg/gpio/gpio_integration_test.go](pkg/gpio)
* Running commands and using expect to look for output:
[/integration/generic-tests/pxeboot_test.go](pxeboot test)
