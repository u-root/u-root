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

Testing mainly relies on 2 packages: [pkg/vmtest](/pkg/vmtest) and
[pkg/qemu](/pkg/qemu).

pkg/vmtest takes in integration test options, and given those and the
environment variables, uses pkg/qemu to start a QEMU VM with the correct command
line and configuration.

Files that need to be shared with the VM are written to a temp dir which is
exposed as a Plan 9 (9p) filesystem in the VM. This includes the kernel and
initramfs being used for the VM.

The test architecture, kernel and QEMU binary are set using environment
variables.

The initramfs can come from the following sources:
* User overridden: when the `UROOT_INITRAMFS` environment variable is used to
  override the initramfs. The user is responsible for ensuring the initramfs
  contains the correct binaries.
* Custom u-root opts: define u-root opts in the test itself (eg. custom uinit).
  The testing setup will generate an initramfs with those options.
* Default: provide the set of commands to be tested. The commands are written to
  an elvish script in the shared dir. The testing setup will generate a generic
  initramfs that mounts the shared 9p filesystem as '/testdata', and then finds
  and runs the elvish script.

To check for the correct behavior, we use the go expect package to find
expected output in QEMU's serial output within a given timeout.

## Running Tests

These tests only run on Linux on amd64 and arm.

1. **Set Environment Variables**

-   `UROOT_QEMU` points to a QEMU binary and args, e.g.

```sh
export UROOT_QEMU="$HOME/bin/qemu-system-x86_64 -enable-kvm"
```

-   `UROOT_KERNEL` points to a Linux kernel binary, e.g.

```sh
export UROOT_KERNEL="$HOME/linux/arch/x86/boot/bzImage"
```
-   (optional) `UROOT_INITRAMFS` is a custom initramfs to use for all tests.
    This will override all other initramfs options defined by the tests.

-   (optional) `UROOT_TESTARCH` (defaults to host architecture) is the
    architecture to test. Only `arm` and `amd64` are supported.

-   (optional) `UROOT_QEMU_TIMEOUT_X` (defaults to 1.0) can be used to multiply
    the timeouts for each test in case QEMU on your machine is slower. For
    example, if you cannot turn on `-enable-kvm`, use `UROOT_QEMU_TIMEOUT_X=2`
    as our test automation does.


Our automated CI uses Dockerfiles to build a kernel and QEMU and set these
environment variables. You can see the Dockerfile and the config file used to
build the kernel for each supported architecture [here](/.circleci/images).

If you don't want to deal with version differences in QEMU and the kernel, you
can use the docker image get both. Inside /.circleci/images/test-image-amd64 (or
whatever arch you have), run

```
cd test-image-$GOARCH
docker build . -t uroottest/test-image-$GOARCH:$VERSION
docker run uroottest/test-image-$GOARCH:$VERSION
docker container list -a
```

Then look for the container id for your newly built container, and

```
docker cp $CONTAINER_ID:bzImage <target>
docker cp $CONTAINER_ID:qemu-system-x86_64 <target>
docker cp $CONTAINER_ID:pc-bios <target>
```

The pc bios needs to be passed into qemu with the -L flag for this built version
of qemu.


2. **Run Tests**

Recall that there are 2 subdirectories with tests, generic-tests/ and gotests/.
To run tests in both directories, run:

```sh
go test [-v] ./...
```

The verbose flag is useful to see the QEMU command line being used and the full
serial output. It is also useful to see which tests are being skipped and why
(particularly for ARM, where many tests are currently skipped).

Unless you want to wait a long time for all tests to complete, run just the
specific test you want from inside the correct directory e.g.

```sh
cd generic-tests/
go test [-v] -test.run=TestDhclient
```

*To avoid having to do this every time, check the instructions for the RUNLOCAL
script.*

## Writing a New Test

To write a new test, first decide which of the options from the previous
section best fit your case (custom initramfs, custom uinit, test commands).

`vmtest.QEMUTest` is the function that starts the QEMU VM and returns the VM
struct. There, provide the test options for your use case.

The VM struct returned by `vmtest.QEMUTest` represents a running QEMU virtual
machine. Use its family of Expect methods to check for the correct result.

