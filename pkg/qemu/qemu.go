// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package qemu provides a Go API for starting QEMU VMs.
//
// qemu is mainly suitable for running QEMU-based integration tests.
//
// The environment variable `UROOT_QEMU` overrides the path to QEMU and the
// first few arguments (defaults to "qemu"). For example, I use:
//
//     UROOT_QEMU='qemu-system-x86_64 -L . -m 4096 -enable-kvm'
//
// For CI, this environment variable is set in `.circleci/images/integration/Dockerfile`.
package qemu

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	expect "github.com/google/goexpect"
)

// DefaultTimeout for `Expect` and `ExpectRE` functions.
var DefaultTimeout = 7 * time.Second

// TimeoutMultiplier increases all timeouts proportionally. Useful when running
// QEMU on a slow machine.
var TimeoutMultiplier = 1.0

func init() {
	if timeoutMultS := os.Getenv("UROOT_QEMU_TIMEOUT_X"); len(timeoutMultS) > 0 {
		t, err := strconv.ParseFloat(timeoutMultS, 64)
		if err == nil {
			TimeoutMultiplier = t
		}
	}
}

// Options are VM start-up parameters.
type Options struct {
	// QEMUPath is the path to the QEMU binary to invoke.
	//
	// If left unspecified, the UROOT_QEMU env var will be used.
	// If the env var is unspecified, "qemu" is the default.
	QEMUPath string

	// Path to the kernel to boot.
	Kernel string

	// Path to the initramfs.
	Initramfs string

	// Extra kernel command-line arguments.
	KernelArgs string

	// Where to send serial output.
	SerialOutput io.WriteCloser

	// Timeout is the expect timeout.
	Timeout time.Duration

	// Devices are devices to expose to the QEMU VM.
	Devices []Device
}

// Start starts a QEMU VM.
func (o *Options) Start() (*VM, error) {
	cmdline, err := o.Cmdline()
	if err != nil {
		return nil, err
	}

	gExpect, ch, err := expect.SpawnWithArgs(cmdline, -1,
		expect.Tee(o.SerialOutput),
		expect.CheckDuration(2*time.Millisecond))
	if err != nil {
		return nil, err
	}
	return &VM{
		Options: o,
		errCh:   ch,
		cmdline: cmdline,
		gExpect: gExpect,
	}, nil
}

// cmdline returns the command line arguments used to start QEMU. These
// arguments are derived from the given QEMU struct.
func (o *Options) Cmdline() ([]string, error) {
	var args []string
	if len(o.QEMUPath) > 0 {
		args = append(args, o.QEMUPath)
	} else {
		// Read first few arguments for env.
		env := os.Getenv("UROOT_QEMU")
		if env == "" {
			env = "qemu" // default
		}
		args = append(args, strings.Fields(env)...)
	}

	// Disable graphics because we are using serial.
	args = append(args, "-nographic")

	// Arguments passed to the kernel:
	//
	// - earlyprintk=ttyS0: print very early debug messages to the serial
	// - console=ttyS0: /dev/console points to /dev/ttyS0 (the serial port)
	// - o.KernelArgs: extra, optional kernel arguments
	// - args required by devices
	for _, dev := range o.Devices {
		if dev != nil {
			if a := dev.KArgs(); a != nil {
				o.KernelArgs += " " + strings.Join(a, " ")
			}
		}
	}
	if len(o.Kernel) != 0 {
		args = append(args, "-kernel", o.Kernel)
		if len(o.KernelArgs) != 0 {
			args = append(args, "-append", o.KernelArgs)
		}
	} else if len(o.KernelArgs) != 0 {
		err := fmt.Errorf("kernel args are required but cannot be added due to bootloader")
		return nil, err
	}
	if len(o.Initramfs) != 0 {
		args = append(args, "-initrd", o.Initramfs)
	}

	for _, dev := range o.Devices {
		if dev != nil {
			if c := dev.Cmdline(); c != nil {
				args = append(args, c...)
			}
		}
	}
	return args, nil
}

// VM is a running QEMU virtual machine.
type VM struct {
	Options *Options
	cmdline []string
	errCh   <-chan error
	gExpect *expect.GExpect
}

// Wait waits for the VM to exit.
func (v *VM) Wait() error {
	return <-v.errCh
}

// Cmdline is the command-line the VM was started with.
func (v *VM) Cmdline() []string {
	// Maybe return a copy?
	return v.cmdline
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

// Close stops QEMU.
func (v *VM) Close() {
	v.gExpect.Close()
	v.gExpect = nil
}

// Send sends a string to QEMU's serial.
func (v *VM) Send(in string) {
	v.gExpect.Send(in)
}

func (v *VM) TimeoutOr() time.Duration {
	if v.Options.Timeout == 0 {
		return DefaultTimeout
	}
	return v.Options.Timeout
}

// Expect returns an error if the given string is not found in vEMU's serial
// output within `DefaultTimeout`.
func (v *VM) Expect(search string) error {
	return v.ExpectTimeout(search, v.TimeoutOr())
}

// ExpectTimeout returns an error if the given string is not found in vEMU's serial
// output within the given timeout.
func (v *VM) ExpectTimeout(search string, timeout time.Duration) error {
	_, err := v.ExpectRETimeout(regexp.MustCompile(regexp.QuoteMeta(search)), timeout)
	return err
}

// ExpectRE returns an error if the given regular expression is not found in
// vEMU's serial output within `DefaultTimeout`. The matched string is
// returned.
func (v *VM) ExpectRE(pattern *regexp.Regexp) (string, error) {
	return v.ExpectRETimeout(pattern, v.TimeoutOr())
}

// ExpectRETimeout returns an error if the given regular expression is not
// found in vEMU's serial output within the given timeout. The matched string
// is returned.
func (v *VM) ExpectRETimeout(pattern *regexp.Regexp, timeout time.Duration) (string, error) {
	scaled := time.Duration(float64(timeout) * TimeoutMultiplier)
	str, _, err := v.gExpect.Expect(pattern, scaled)
	return str, err
}
