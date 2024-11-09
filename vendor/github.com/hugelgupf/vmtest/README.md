## vmtest

[![Go Report Card](https://goreportcard.com/badge/github.com/hugelgupf/vmtest)](https://goreportcard.com/report/github.com/hugelgupf/vmtest)
[![GoDoc](https://godoc.org/github.com/hugelgupf/vmtest?status.svg)](https://godoc.org/github.com/hugelgupf/vmtest)

vmtest is a Go API for launching QEMU VMs.

* [The `qemu` package](https://pkg.go.dev/github.com/hugelgupf/vmtest/qemu)
  contains APIs for

    * launching QEMU processes
    * configuring QEMU devices (such as a shared 9P directory, networking,
      serial logging, etc)
    * running tasks (goroutines) bound to the VM process lifetime, and
    * using expect-scripting to check for outputs.

    * [`quimage`](https://pkg.go.dev/github.com/hugelgupf/vmtest/quimage)
      can be used to configure a Go-based u-root initramfs to be used as the
      root file system, including any arbitrary Go commands.
    * [`qnetwork`](https://pkg.go.dev/github.com/hugelgupf/vmtest/qemu/qnetwork)
      (WIP) is an API to QEMU network devices.
    * [`qevent`](https://pkg.go.dev/github.com/hugelgupf/vmtest/qemu/qevent)
      provides a JSON-over-virtio-serial channel from guest to host.
    * [`qcoverage`](https://pkg.go.dev/github.com/hugelgupf/vmtest/qemu/qcoverage)
      adds utilities to collect kernel & Go
      [`GOCOVERDIR`-based](https://go.dev/doc/build-cover) integration test
      coverage.

* [The `govmtest` package](https://pkg.go.dev/github.com/hugelgupf/vmtest/govmtest)
  (WIP) contains an API for running Go unit tests in the guest and collecting
  their test coverage as well as results.

* [The `scriptvm` package](https://pkg.go.dev/github.com/hugelgupf/vmtest/scriptvm)
  (WIP) contains an API for running shell scripts in the guest.

## Running Tests

The `qemu` API picks up the following values from env vars by default:

* `VMTEST_QEMU`: QEMU binary + arguments (e.g.
  `VMTEST_QEMU="qemu-system-x86_64 -enable-kvm"`)
* `VMTEST_KERNEL`: Kernel to boot.
* `VMTEST_ARCH`: Guest architecture (same as GOARCH values). Must match the QEMU
  binary supplied. If not supplied, defaults to `runtime.GOARCH`, i.e. it
  matches the host's GOARCH.
* `VMTEST_KERNEL_APPEND`: is added to kernel command-line arguments
* `VMTEST_QEMU_APPEND`: is added to QEMU command-line arguments.
* `VMTEST_TIMEOUT`: Timeout value (e.g. `1m20s` -- parsed by Go's
  `time.ParseDuration`).
* `VMTEST_INITRAMFS`: Initramfs to boot.

Most of these values can be overriden in the Go API, but typically only
`VMTEST_INITRAMFS` and `VMTEST_TIMEOUT` are. `VMTEST_KERNEL_APPEND` and
`VMTEST_QEMU_APPEND` are always additive.

The `runvmtest` tool automatically downloads `VMTEST_QEMU` and
`VMTEST_KERNEL` for use with tests based on a provided `VMTEST_ARCH`. E.g.

```sh
go install github.com/hugelgupf/vmtest/tools/runvmtest@latest

# See how it works:
runvmtest -- bash -c "echo \$VMTEST_KERNEL -- \$VMTEST_QEMU"

# Intended usage:
runvmtest -- go test -v ./tests/gohello

# Or run an Arm64 guest:
VMTEST_ARCH=arm64 runvmtest -- go test -v ./tests/gohello
```

You can also override one or both, which will just be passed through:

```sh
# Will only download VMTEST_KERNEL.
VMTEST_ARCH=arm64 VMTEST_QEMU="qemu-system-aarch64 -enable-kvm" runvmtest -- go test -v ./tests/gohello
```

To keep the artifacts around locally to reproduce the same test:

```s
runvmtest --keep-artifacts -- go test -v ./tests/gohello
```

The default kernel and QEMU supplied by `runvmtest` may of course not work well
for your tests. You can configure `runvmtest` to supply your own `VMTEST_KERNEL`
and `VMTEST_QEMU` -- but also any additional environment variables. See
[`runvmtest` configuration](#custom-runvmtest-configuration).

To build your own kernel or QEMU, check out
[images/kernel-arm64](./images/kernel-arm64) for building a kernel-image-only
Docker image, and [images/qemu](./images/qemu/Dockerfile) for how we build a
Docker image with just QEMU binaries and their dependencies.

## Writing Tests

### Architecture-independent test with shared directory

The recommendation for architecture-independent tests is to allow
architecture-dependent values like `VMTEST_QEMU` and `VMTEST_KERNEL` be provided
by `runvmtest` based on `VMTEST_ARCH`, while configuring everything else with
the Go API.

See e.g. [examples/shareddir](./examples/shareddir/vm_test.go)

Run it with:

```sh
$ cd examples/shareddir
$ VMTEST_ARCH=amd64 runvmtest -- go test -v
$ VMTEST_ARCH=arm64 runvmtest -- go test -v
$ VMTEST_ARCH=riscv64 runvmtest -- go test -v
$ VMTEST_ARCH=arm runvmtest -- go test -v
```

They all should pass.

Every test in this repository is written this way to test every feature on all
architectures.

### Example: Go unit tests in VM

See [tests/gobench](./tests/gobench/bench_test.go)

### Example: qemu API with u-root initramfs

```go
func TestStartVM(t *testing.T) {
    vm := qemu.StartT(
        t,
        "vm", // Prefix for t.Logf serial console lines.
        qemu.ArchUseEnvv,
        quimage.WithUimageT(t,
                uimage.WithInit("init"),
                uimage.WithUinit("shutdownafter", "--", "cat", "etc/thatfile"),
                uimage.WithBusyboxCommands(
                        "github.com/u-root/u-root/cmds/core/init",
                        "github.com/u-root/u-root/cmds/core/cat",
                        "github.com/hugelgupf/vmtest/vminit/shutdownafter",
                ),
                uimage.WithFiles("./testdata/foo:etc/thatfile"),
        ),

        // Other options...
    )
    // ...
}
```

### Example: qemu API

The qemu API can be used without `testing.TB`, and any of the environment
variables can be overridden in the Go API:

```go
func TestStartVM(t *testing.T) {
    vm, err := qemu.Start(
        // Or use qemu.ArchUseEnvv and set VMTEST_ARCH=amd64 (values like GOARCH)
        qemu.ArchAMD64,

        // Or omit and set VMTEST_QEMU="qemu-system-x86_64 -enable-kvm"
        qemu.WithQEMUCommand("qemu-system-x86_64 -enable-kvm"),

        // Or omit and set VMTEST_KERNEL=./foobar
        qemu.WithKernel("./foobar"),

        // Or omit and set VMTEST_INITRAMFS=./somewhere.cpio
        // Or use u-root initramfs builder in quimage package. See example below.
        qemu.WithInitramfs("./somewhere.cpio"),

        qemu.WithAppendKernel("console=ttyS0 earlyprintk=ttyS0"),
        qemu.LogSerialByLine(qemu.DefaultPrint("vm", t.Logf)),
    )
    if err != nil {
        t.Fatalf("Failed to start VM: %v", err)
    }
    t.Logf("cmdline: %#v", vm.CmdlineQuoted())

    if _, err := vm.Console.ExpectString("Kernel command line:"); err != nil {
        t.Errorf("Error expecting kernel command line string: %v", err)
    }

    if err := vm.Wait(); err != nil {
        t.Fatalf("Error waiting for VM to exit: %v", err)
    }
}
```

### Example: Tasks

```go
func TestStartVM(t *testing.T) {
    vm, err := qemu.Start(
        qemu.ArchUseEnvv,
        // Other config ...

        // Runs a goroutine alongside the QEMU process, which is canceled via
        // context when QEMU exits.
        qemu.WithTask(
            func(ctx context.Context, n *qemu.Notifications) error {
                // If this were an HTTP server or something not expected to exit
                // cleanly when the guest exits, probably want to ignore SIGKILL error.
                return exec.CommandContext(ctx, "sleep", "900").Run()
            },
        ),

        // Task that runs when the VM exits.
        qemu.WithTask(qemu.Cleanup(func() error {
            // Do something.
            return fmt.Errorf("this is returned by vm.Wait()")
        })),

        // Task that only runs when VM starts.
        qemu.WithTask(qemu.WaitVMStarted(...)),
    )
    // ...
}
```

## Custom runvmtest configuration

`runvmtest` tries to read a config from `.vmtest.yaml` in the current working
directory or any of its parents.

Given this is a Go-based test framework, the recommendation would be to place
`.vmtest.yaml` in the same directory as your `go.mod` so that the config is
available anywhere `go test` is for that module.

`runvmtest` can be configured to set up any number of environment variables.
Config format looks like this:

```
VMTEST_ARCH:
  ENV_VAR:
    container: <container name>
    template: "{{.somedir}} -foobar {{.somefile}}"
    files:
      somefile: <path in container to copy to a tmpfile>
    directories:
      somedir: <path in container to copy to a tmpdir>
```

Check out the example in
[tools/runvmtest/example-vmtest.yaml](./tools/runvmtest/example-vmtest.yaml).
