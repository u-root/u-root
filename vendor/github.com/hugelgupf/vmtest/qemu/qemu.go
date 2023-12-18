// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package qemu provides a Go API for starting QEMU VMs.
//
// qemu is mainly suitable for running QEMU-based integration tests.
//
// The environment variable `VMTEST_QEMU` overrides the path to QEMU and the
// first few arguments. For example:
//
//	VMTEST_QEMU='qemu-system-x86_64 -L . -m 4096 -enable-kvm'
//
// Other environment variables:
//
//	VMTEST_ARCH (used when Arch is empty or ArchUseEnvv is set)
//	VMTEST_KERNEL (used when Options.Kernel is empty)
//	VMTEST_INITRAMFS (used when Options.Initramfs is empty)
//	VMTEST_TIMEOUT (used when Options.VMTimeout is empty)
package qemu

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Netflix/go-expect"
	"golang.org/x/exp/slices"
	"golang.org/x/sync/errgroup"
)

// ErrKernelRequiredForArgs is returned when KernelArgs is populated but Kernel is empty.
var ErrKernelRequiredForArgs = errors.New("KernelArgs can only be used when Kernel is also specified due to how QEMU bootloader works")

// ErrNoArch is returned when neither Arch nor VMTEST_ARCH are set.
var ErrNoArch = errors.New("no guest architecture specified -- guest arch is required to decide some QEMU command-line arguments")

// ErrUnsupportedArch is returned when an unsupported guest architecture value is used.
var ErrUnsupportedArch = errors.New("unsupported guest architecture specified -- guest arch is required to decide some QEMU command-line arguments")

// Arch is the QEMU guest architecture.
type Arch string

// Architecture values are derived from GOARCH values.
const (
	// ArchUseEnvv will derive the architecture from the VMTEST_ARCH env var.
	ArchUseEnvv Arch = ""

	// ArchAMD64 is the x86 64bit architecture.
	ArchAMD64 Arch = "amd64"

	// ArchI386 is the x86 32bit architecture.
	ArchI386 Arch = "i386"

	// ArchArm64 is the aarch64 architecture.
	ArchArm64 Arch = "arm64"

	// ArchArm is the arm 32bit architecture.
	ArchArm Arch = "arm"
)

// SupportedArches are the supported guest architecture values.
var SupportedArches = []Arch{
	ArchAMD64,
	ArchI386,
	ArchArm64,
	ArchArm,
}

// GuestArch returns the Guest architecture under test. Either VMTEST_ARCH or
// runtime.GOARCH.
func GuestArch() Arch {
	if env := Arch(os.Getenv("VMTEST_ARCH")); slices.Contains(SupportedArches, env) {
		return env
	}
	return Arch(runtime.GOARCH)
}

// Valid returns whether the guest arch is a supported guest arch value.
func (g Arch) Valid() bool {
	return slices.Contains(SupportedArches, g)
}

// Arch returns the guest architecture.
func (g Arch) Arch() Arch {
	if g == ArchUseEnvv {
		g = GuestArch()
	}
	return g
}

// Fn is a QEMU configuration option supplied to Start or OptionsFor.
//
// Fns rely on a QEMU architecture already having been determined.
type Fn func(*IDAllocator, *Options) error

// WithQEMUCommand sets a QEMU command. It's expected to provide a QEMU binary
// and optionally some arguments.
//
// cmd may contain additional QEMU args, such as "qemu-system-x86_64 -enable-kvm -m 1G".
// They will be appended to the command-line.
func WithQEMUCommand(cmd string) Fn {
	return func(alloc *IDAllocator, opts *Options) error {
		opts.QEMUCommand = cmd
		return nil
	}
}

// WithKernel sets the path to the kernel binary.
func WithKernel(kernel string) Fn {
	return func(alloc *IDAllocator, opts *Options) error {
		opts.Kernel = kernel
		return nil
	}
}

// WithInitramfs sets the path to the initramfs.
func WithInitramfs(initramfs string) Fn {
	return func(alloc *IDAllocator, opts *Options) error {
		opts.Initramfs = initramfs
		return nil
	}
}

// WithAppendKernel appends kernel arguments.
func WithAppendKernel(args ...string) Fn {
	return func(alloc *IDAllocator, opts *Options) error {
		opts.AppendKernel(strings.Join(args, " "))
		return nil
	}
}

// WithSerialOutput writes serial output to w as well.
func WithSerialOutput(w ...io.WriteCloser) Fn {
	return func(alloc *IDAllocator, opts *Options) error {
		opts.SerialOutput = append(opts.SerialOutput, w...)
		return nil
	}
}

// WithVMTimeout is a timeout for the QEMU guest subprocess.
func WithVMTimeout(timeout time.Duration) Fn {
	return func(alloc *IDAllocator, opts *Options) error {
		opts.VMTimeout = timeout
		return nil
	}
}

// WithTask adds a goroutine running alongside the guest.
//
// Task goroutines are started right before the guest is started.
//
// VM.Wait waits for all tasks to complete before returning an error. Errors
// produced by tasks are returned by VM.Wait.
//
// A task is expected to exit either when ctx is canceled or when the QEMU
// subprocess exits. When the context is canceled, the QEMU subprocess is
// expected to exit as well, and when the QEMU subprocess exits, the context is
// canceled.
func WithTask(t ...Task) Fn {
	return func(alloc *IDAllocator, opts *Options) error {
		opts.Tasks = append(opts.Tasks, t...)
		return nil
	}
}

// OptionsFor evaluates the given config functions and returns an Options object.
func OptionsFor(arch Arch, fns ...Fn) (*Options, error) {
	var vmTimeout time.Duration
	if d := os.Getenv("VMTEST_TIMEOUT"); len(d) > 0 {
		var err error
		vmTimeout, err = time.ParseDuration(d)
		if err != nil {
			return nil, fmt.Errorf("invalid VMTEST_TIMEOUT value: %w", err)
		}
	}

	o := &Options{
		QEMUCommand: os.Getenv("VMTEST_QEMU"),
		Kernel:      os.Getenv("VMTEST_KERNEL"),
		Initramfs:   os.Getenv("VMTEST_INITRAMFS"),
		VMTimeout:   vmTimeout,
		// Disable graphics by default.
		QEMUArgs: []string{"-nographic"},
	}

	if err := o.setArch(arch.Arch()); err != nil {
		return nil, err
	}

	alloc := NewIDAllocator()
	for _, f := range fns {
		if err := f(alloc, o); err != nil {
			return nil, err
		}
	}
	return o, nil
}

// Start starts a VM with the given configuration.
//
// SerialOutput will be relayed only if VM.Wait is also called some time after
// the VM starts.
func Start(arch Arch, fns ...Fn) (*VM, error) {
	return StartContext(context.Background(), arch, fns...)
}

// StartContext starts a VM with the given configuration and with the given context.
//
// When the context is done, the QEMU subprocess will be killed and all
// associated goroutines cleaned up as long as VM.Wait() is called.
func StartContext(ctx context.Context, arch Arch, fns ...Fn) (*VM, error) {
	o, err := OptionsFor(arch, fns...)
	if err != nil {
		return nil, err
	}
	return o.Start(ctx)
}

// Options are VM start-up parameters.
type Options struct {
	// arch is the QEMU architecture used.
	//
	// Some device decisions are made based on the architecture.
	// If empty, VMTEST_QEMU_ARCH env var will be used.
	arch Arch

	// QEMUCommand is QEMU binary to invoke and some additonal args.
	//
	// If empty, the VMTEST_QEMU env var will be used.
	QEMUCommand string

	// Path to the kernel to boot.
	//
	// If empty, VMTEST_KERNEL env var will be used.
	Kernel string

	// Path to the initramfs.
	//
	// If empty, VMTEST_INITRAMFS env var will be used.
	Initramfs string

	// Extra kernel command-line arguments.
	KernelArgs string

	// Where to send serial output.
	SerialOutput []io.WriteCloser

	// Tasks are goroutines running alongside the guest.
	//
	// Task goroutines are started right before the guest is started.
	//
	// A task is expected to exit either when ctx is canceled or when the
	// QEMU subprocess exits. When the context is canceled, the QEMU
	// subprocess is expected to exit as well, and when the QEMU subprocess
	// exits, the context is canceled.
	Tasks []Task

	// Additional QEMU cmdline arguments.
	QEMUArgs []string

	// VMTimeout is a timeout for the QEMU subprocess.
	VMTimeout time.Duration

	// ExtraFiles are extra files passed to QEMU on start.
	ExtraFiles []*os.File
}

// AddFile adds the file to the QEMU process and returns the FD it will be in
// the child process.
func (o *Options) AddFile(f *os.File) int {
	o.ExtraFiles = append(o.ExtraFiles, f)

	// 0, 1, 2 used for stdin/out/err.
	return len(o.ExtraFiles) + 2
}

// Task is a task running alongside the guest.
//
// A task is expected to exit either when ctx is canceled or when the QEMU
// subprocess exits.
type Task func(ctx context.Context, n *Notifications) error

// WaitVMStarted waits until the VM has started before starting t, or never
// starts t if context is canceled before the VM is started.
func WaitVMStarted(t Task) Task {
	return func(ctx context.Context, n *Notifications) error {
		// Wait until VM starts or exit if it never does.
		select {
		case <-n.VMStarted:
		case <-ctx.Done():
			return nil
		}
		return t(ctx, n)
	}
}

// Notifications gives tasks the option to wait for certain VM events.
//
// Tasks must not be required to listen on notifications; there must be no
// blocking channel I/O.
type Notifications struct {
	// VMStarted will be closed when the VM is started.
	VMStarted chan struct{}

	// VMExited will receive exactly 1 event when the VM exits and then be closed.
	VMExited chan error
}

func newNotifications() *Notifications {
	return &Notifications{
		VMStarted: make(chan struct{}),
		VMExited:  make(chan error, 1),
	}
}

// Arch returns the guest architecture.
func (o *Options) Arch() Arch {
	return o.arch
}

// Start starts a QEMU VM.
func (o *Options) Start(ctx context.Context) (*VM, error) {
	cmdline, err := o.Cmdline()
	if err != nil {
		return nil, err
	}

	c, err := expect.NewConsole()
	if err != nil {
		return nil, err
	}

	var cancel context.CancelFunc
	if o.VMTimeout != 0 {
		ctx, cancel = context.WithTimeout(ctx, o.VMTimeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}
	vm := &VM{
		Options: o,
		Console: c,
		cmdline: cmdline,
		cancel:  cancel,
	}
	for _, task := range o.Tasks {
		// Capture the var... Go stuff.
		task := task
		n := newNotifications()
		vm.taskWG.Go(func() error {
			return task(ctx, n)
		})
		vm.notifs = append(vm.notifs, n)
	}

	writers := []io.Writer{c.Tty()}
	for _, serial := range o.SerialOutput {
		writers = append(writers, serial)
	}
	cmd := exec.CommandContext(ctx, cmdline[0], cmdline[1:]...)
	cmd.Stdin = c.Tty()
	cmd.Stdout = io.MultiWriter(writers...)
	cmd.Stderr = io.MultiWriter(writers...)
	cmd.ExtraFiles = o.ExtraFiles
	if err := cmd.Start(); err != nil {
		// Cancel tasks.
		cancel()
		// Wait for tasks to exit. Some day we'll report their errors
		// with errors.Join.
		_ = vm.taskWG.Wait()
		return nil, err
	}
	vm.notifs.vmStarted()
	vm.cmd = cmd
	vm.wait = make(chan struct{})

	// A goroutine to wait on exit, as we need to close Console.Tty() to
	// unblock any waiting Expect calls.
	go func() {
		err := vm.cmd.Wait()
		vm.notifs.vmExited(err)

		// Close the pts end of the tty to unblock any potential
		// readers on ptm (i.e. Expect calls).
		//
		// Don't call vm.Console.Close() as that also closes the ptm,
		// which a blocking Expect call may still expect to read from.
		vm.Console.Tty().Close()
		vm.waitMu.Lock()
		vm.waitErr = err
		vm.waitMu.Unlock()
		close(vm.wait)
	}()
	return vm, nil
}

func (o *Options) setArch(arch Arch) error {
	if len(arch) == 0 {
		return ErrNoArch
	}
	if !arch.Valid() {
		return fmt.Errorf("%w: %s", ErrUnsupportedArch, arch)
	}
	o.arch = arch
	return nil
}

// AppendKernel appends to kernel args.
func (o *Options) AppendKernel(s ...string) {
	if len(s) == 0 {
		return
	}
	t := strings.Join(s, " ")
	if len(o.KernelArgs) == 0 {
		o.KernelArgs = t
	} else {
		o.KernelArgs += " " + t
	}
}

// AppendQEMU appends args to the QEMU command line.
func (o *Options) AppendQEMU(s ...string) {
	o.QEMUArgs = append(o.QEMUArgs, s...)
}

// Cmdline returns the command line arguments used to start QEMU. These
// arguments are derived from the given QEMU struct.
func (o *Options) Cmdline() ([]string, error) {
	var args []string

	// QEMU binary + initial args (may have been supplied via VMTEST_QEMU).
	args = append(args, strings.Fields(o.QEMUCommand)...)

	// Add user / configured args.
	args = append(args, o.QEMUArgs...)

	if len(o.Kernel) > 0 {
		args = append(args, "-kernel", o.Kernel)
		if len(o.KernelArgs) != 0 {
			args = append(args, "-append", o.KernelArgs)
		}
	} else if len(o.KernelArgs) != 0 {
		return nil, ErrKernelRequiredForArgs
	}

	if len(o.Initramfs) != 0 {
		args = append(args, "-initrd", o.Initramfs)
	}

	return args, nil
}

// VM is a running QEMU virtual machine.
type VM struct {
	// Console provides in/output to the QEMU subprocess.
	Console *expect.Console

	// Options are the options that were used to start the VM.
	//
	// They are not used once the VM is started.
	Options *Options

	// cmd is the QEMU subprocess.
	cmd *exec.Cmd

	// The cmdline that the QEMU subprocess was started with.
	cmdline []string

	// State related to tasks.
	taskWG errgroup.Group
	notifs notifications
	cancel func()

	wait       chan struct{}
	waitMu     sync.Mutex
	waitErr    error
	waitCalled atomic.Bool
}

// Cmdline is the command-line the VM was started with.
func (v *VM) Cmdline() []string {
	// Maybe return a copy?
	return v.cmdline
}

// Kill kills the QEMU subprocess.
//
// Callers are still responsible for calling VM.Wait after calling kill to
// clean up task goroutines and to get remaining serial console output.
func (v *VM) Kill() error {
	return v.cmd.Process.Kill()
}

// Signal signals the QEMU subprocess.
//
// Callers are still responsible for calling VM.Wait if the subprocess exits
// due to this signal to clean up task goroutines and to get remaining serial
// console output.
func (v *VM) Signal(sig os.Signal) error {
	return v.cmd.Process.Signal(sig)
}

// Waited returns whether Wait has been called on VM.
func (v *VM) Waited() bool {
	return v.waitCalled.Load()
}

// Wait waits for the VM to exit and expects EOF from the expect console.
func (v *VM) Wait() error {
	v.waitCalled.Store(true)

	<-v.wait

	v.waitMu.Lock()
	err := v.waitErr
	v.waitMu.Unlock()

	// Close everything but the pts (which was already closed).
	v.Console.Close()
	for _, w := range v.Options.SerialOutput {
		w.Close()
	}

	v.cancel()
	// Wait for all tasks to exit.
	if werr := v.taskWG.Wait(); werr != nil && err == nil {
		err = werr
	}
	return err
}

// CmdlineQuoted quotes any of QEMU's command line arguments containing a space
// so it is easy to copy-n-paste into a shell for debugging.
func (v *VM) CmdlineQuoted() string {
	args := make([]string, len(v.cmdline))
	for i, arg := range v.cmdline {
		if strings.ContainsAny(arg, " \t\n") {
			args[i] = fmt.Sprintf("'%s'", arg)
		} else {
			args[i] = arg
		}
	}
	return strings.Join(args, " ")
}

type notifications []*Notifications

func (n notifications) vmStarted() {
	for _, m := range n {
		close(m.VMStarted)
	}
}

func (n notifications) vmExited(err error) {
	for _, m := range n {
		m.VMExited <- err
		close(m.VMExited)
	}
}
