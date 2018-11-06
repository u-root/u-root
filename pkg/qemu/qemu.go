// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package qemu is suitable for running QEMU-based integration tests.
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
	"time"

	"github.com/google/goexpect"
)

// DefaultTimeout for `Expect` and `ExpectRE` functions.
var DefaultTimeout = 7 * time.Second

// TimeoutMultiplier increases all timeouts proportionally. Useful when running
// QEMU on a slow machine.
var TimeoutMultiplier = 2.0

type Network struct {
	port   uint16
	numVMs uint8
}

func NewNetwork() *Network {
	return &Network{
		port: 1234,
	}
}

func (n *Network) newVM() *networkState {
	num := n.numVMs
	n.numVMs++
	return &networkState{
		connect: num != 0,
		mac:     net.HardwareAddr{0x0e, 0x00, 0x00, 0x00, 0x00, byte(num)},
		port:    n.port,
	}
}

type networkState struct {
	// Whether to connect or listen.
	connect bool
	mac     net.HardwareAddr
	port    uint16
}

// Options is filled and pass to `Start()`.
type Options struct {
	// Path to the bzImage kernel
	Kernel string

	// Path to the initramfs.
	Initramfs string

	// Extra kernel arguments.
	KernelArgs string

	// SharedDir is a directory that will be mountable inside the VM as
	// /dev/sda1.
	SharedDir string

	// ExtraArgs are additional QEMU arguments.
	ExtraArgs []string

	// Where to send serial output.
	SerialOutput io.WriteCloser

	// Timeout is the expect timeout.
	Timeout time.Duration

	// Network to expose inside the VM.
	Network *Network
}

// cmdline returns the command line arguments used to start QEMU. These
// arguments are derived from the given QEMU struct.
func cmdline(o *Options, net *networkState) []string {
	// Read first few arguments for env.
	env := os.Getenv("UROOT_QEMU")
	if env == "" {
		env = "qemu" // default
	}
	args := strings.Fields(env)

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

	if len(o.SharedDir) != 0 {
		// Expose the temp directory to QEMU as /dev/sda1
		args = append(args, "-drive", fmt.Sprintf("file=fat:ro:%s,if=none,id=tmpdir", o.SharedDir))
		args = append(args, "-device", "ich9-ahci,id=ahci")
		args = append(args, "-device", "ide-drive,drive=tmpdir,bus=ahci.0")
	}

	if net != nil {
		args = append(args, "-net", fmt.Sprintf("nic,macaddr=%s", net.mac))
		if net.connect {
			args = append(args, "-net", fmt.Sprintf("socket,connect=:%d", net.port))
		} else {
			args = append(args, "-net", fmt.Sprintf("socket,listen=:%d", net.port))
		}
	}

	if o.ExtraArgs != nil {
		args = append(args, o.ExtraArgs...)
	}
	return args
}

// VM is a running QEMU virtual machine.
type VM struct {
	Options *Options
	cmdline []string
	network *networkState
	gExpect *expect.GExpect
}

// Start a QEMU VM.
func Start(o *Options) (*VM, error) {
	var net *networkState
	if o.Network != nil {
		net = o.Network.newVM()
	}

	cmdline := cmdline(o, net)

	gExpect, _, err := expect.SpawnWithArgs(cmdline, -1,
		expect.Tee(o.SerialOutput),
		expect.CheckDuration(2*time.Millisecond))
	if err != nil {
		return nil, err
	}
	return &VM{
		Options: o,
		cmdline: cmdline,
		network: net,
		gExpect: gExpect,
	}, nil
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
