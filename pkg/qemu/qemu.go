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
	"net"
	"os"
	"regexp"
	"strings"
	"sync/atomic"
	"time"

	expect "github.com/google/goexpect"
)

// DefaultTimeout for `Expect` and `ExpectRE` functions.
var DefaultTimeout = 7 * time.Second

// TimeoutMultiplier increases all timeouts proportionally. Useful when running
// QEMU on a slow machine.
var TimeoutMultiplier = 2.0

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
	cmdline := cmdline(o)

	gExpect, _, err := expect.SpawnWithArgs(cmdline, -1,
		expect.Tee(o.SerialOutput),
		expect.CheckDuration(2*time.Millisecond))
	if err != nil {
		return nil, err
	}
	return &VM{
		Options: o,
		cmdline: cmdline,
		gExpect: gExpect,
	}, nil
}

// Device is a QEMU device to expose to a VM.
type Device interface {
	// Cmdline returns arguments to append to the QEMU command line for this device.
	Cmdline() []string
}

// Network is a Device that can connect multiple QEMU VMs to each other.
//
// Network uses the QEMU socket mechanism to connect multiple VMs with a simple
// TCP socket.
type Network struct {
	port uint16

	// numVMs must be atomically accessed so VMs can be started in parallel
	// in goroutines.
	numVMs uint32
}

func NewNetwork() *Network {
	return &Network{
		port: 1234,
	}
}

// Cmdline implements Device.
func (n *Network) Cmdline() []string {
	if n == nil {
		return nil
	}

	newNum := atomic.AddUint32(&n.numVMs, 1)
	num := newNum - 1

	// MAC for the virtualized NIC.
	//
	// This is from the range of locally administered address ranges.
	mac := net.HardwareAddr{0x0e, 0x00, 0x00, 0x00, 0x00, byte(num)}

	args := []string{"-net", fmt.Sprintf("nic,macaddr=%s", mac)}
	if num != 0 {
		args = append(args, "-net", fmt.Sprintf("socket,connect=:%d", n.port))
	} else {
		args = append(args, "-net", fmt.Sprintf("socket,listen=:%d", n.port))
	}
	return args
}

// ReadOnlyDirectory is a Device that exposes a directory as a /dev/sda1
// readonly vfat partition in the VM.
type ReadOnlyDirectory struct {
	// Dir is the directory to expose as a read-only vfat partition.
	Dir string
}

// Cmdline implements Device.
func (rod ReadOnlyDirectory) Cmdline() []string {
	if len(rod.Dir) == 0 {
		return nil
	}

	// Expose the temp directory to QEMU as /dev/sda1
	return []string{
		"-drive", fmt.Sprintf("file=fat:ro:%s,if=none,id=tmpdir", rod.Dir),
		"-device", "ich9-ahci,id=ahci",
		"-device", "ide-drive,drive=tmpdir,bus=ahci.0",
	}
}

// VirtioRandom exposes a PCI random number generator Device to the QEMU VM.
type VirtioRandom struct{}

// Cmdline implements Device.
func (VirtioRandom) Cmdline() []string {
	return []string{"-device", "virtio-rng-pci"}
}

// ArbitraryArgs allows users to add arbitrary arguments to the QEMU command
// line.
type ArbitraryArgs []string

// Cmdline implements Device.
func (aa ArbitraryArgs) Cmdline() []string {
	return aa
}

// cmdline returns the command line arguments used to start QEMU. These
// arguments are derived from the given QEMU struct.
func cmdline(o *Options) []string {
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
	if len(o.Kernel) != 0 {
		args = append(args, "-kernel", o.Kernel)
		cmdline := "console=ttyS0 earlyprintk=ttyS0"
		if len(o.KernelArgs) != 0 {
			cmdline += " " + o.KernelArgs
		}
		args = append(args, "-append", cmdline)
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
	return args
}

// VM is a running QEMU virtual machine.
type VM struct {
	Options *Options
	cmdline []string
	gExpect *expect.GExpect
}

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
