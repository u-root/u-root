## vmtest

[![CircleCI](https://circleci.com/gh/hugelgupf/vmtest.svg?style=svg)](https://circleci.com/gh/hugelgupf/vmtest)
[![Go Report Card](https://goreportcard.com/badge/github.com/hugelgupf/vmtest)](https://goreportcard.com/report/github.com/hugelgupf/vmtest)
[![GoDoc](https://godoc.org/github.com/hugelgupf/vmtest?status.svg)](https://godoc.org/github.com/hugelgupf/vmtest)

Fun stuff coming

### Example: qemu API

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
        // Or use u-root initramfs builder in uqemu package. See example below.
        qemu.WithInitramfs("./somewhere.cpio"),

        qemu.WithAppendKernel("console=ttyS0 earlyprintk=ttyS0"),
        qemu.LogSerialByLine(qemu.PrintLineWithPrefix("vm", t.Logf)),
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

### Example: qemu API with u-root initramfs

```go
func TestStartVM(t *testing.T) {
    l := &ulogtest.Logger{TB: t}
    initramfs := uroot.Opts{
        TempDir:   t.TempDir(),
        InitCmd:   "init",
        UinitCmd:  "cat",
        UinitArgs: []string{"etc/thatfile"},
        Commands: uroot.BusyBoxCmds(
            "github.com/u-root/u-root/cmds/core/init",
            "github.com/u-root/u-root/cmds/core/cat",
        ),
        ExtraFiles: []string{
            "./testdata/foo:etc/thatfile",
        },
    }
    vm, err := qemu.Start(
        qemu.ArchUseEnvv,
        uqemu.WithUrootInitramfs(l, initramfs, filepath.Join(t.TempDir(), "initramfs.cpio")),

        // Other options...
    )
    // ...
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
    )
    // ...
}
```

