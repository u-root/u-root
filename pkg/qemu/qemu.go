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
	"errors"
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

func (n *Network) NewVM(q *QEMU) {
	num := n.numVMs
	n.numVMs++
	q.network = &networkState{
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

// QEMU is filled and pass to `Start()`.
type QEMU struct {
	// Path to the bzImage kernel
	Kernel string

	// Path to the initramfs.
	Initramfs string

	// Extra kernel arguments
	KernelArgs string

	SharedDir string

	network *networkState

	// Extra QEMU arguments
	ExtraArgs []string

	// Where to send serial output.
	SerialOutput io.WriteCloser

	Timeout time.Duration

	gExpect *expect.GExpect
}

// Cmdline returns the command line arguments used to start QEMU. These
// arguments are derived from the given QEMU struct.
func (q *QEMU) Cmdline() []string {
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
	// - q.KernelArgs: extra, optional kernel arguments
	if len(q.Kernel) != 0 {
		args = append(args, "-kernel", q.Kernel)
		cmdline := "console=ttyS0 earlyprintk=ttyS0"
		if len(q.KernelArgs) != 0 {
			cmdline += " " + q.KernelArgs
		}
		args = append(args, "-append", cmdline)
	}
	if len(q.Initramfs) != 0 {
		args = append(args, "-initrd", q.Initramfs)
	}

	if len(q.SharedDir) != 0 {
		// Expose the temp directory to QEMU as /dev/sda1
		args = append(args, "-drive", fmt.Sprintf("file=fat:ro:%s,if=none,id=tmpdir", q.SharedDir))
		args = append(args, "-device", "ich9-ahci,id=ahci")
		args = append(args, "-device", "ide-drive,drive=tmpdir,bus=ahci.0")
	}

	if q.network != nil {
		args = append(args, "-net", fmt.Sprintf("nic,macaddr=%s", q.network.mac))
		if q.network.connect {
			args = append(args, "-net", fmt.Sprintf("socket,connect=:%d", q.network.port))
		} else {
			args = append(args, "-net", fmt.Sprintf("socket,listen=:%d", q.network.port))
		}
	}

	if q.ExtraArgs != nil {
		args = append(args, q.ExtraArgs...)
	}
	return args
}

// CmdlineQuoted quotes any of QEMU's command line arguments containing a space
// so it is easy to copy-n-paste into a shell for debugging.
func (q *QEMU) CmdlineQuoted() (cmdline string) {
	args := q.Cmdline()
	for i, v := range q.Cmdline() {
		if strings.ContainsAny(v, " \t\n") {
			args[i] = fmt.Sprintf("'%s'", v)
		}
	}
	return strings.Join(args, " ")
}

// Start QEMU.
func (q *QEMU) Start() error {
	if q.gExpect != nil {
		return errors.New("QEMU already started")
	}
	var err error
	q.gExpect, _, err = expect.SpawnWithArgs(q.Cmdline(), -1,
		expect.Tee(q.SerialOutput),
		expect.CheckDuration(2*time.Millisecond))
	return err
}

// Close stops QEMU.
func (q *QEMU) Close() {
	q.gExpect.Close()
	q.gExpect = nil
}

// Send sends a string to QEMU's serial.
func (q *QEMU) Send(in string) {
	q.gExpect.Send(in)
}

// Expect returns an error if the given string is not found in QEMU's serial
// output within `DefaultTimeout`.
func (q *QEMU) Expect(search string) error {
	timeout := q.Timeout
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	return q.ExpectTimeout(search, timeout)
}

// ExpectTimeout returns an error if the given string is not found in QEMU's serial
// output within the given timeout.
func (q *QEMU) ExpectTimeout(search string, timeout time.Duration) error {
	_, err := q.ExpectRETimeout(regexp.MustCompile(regexp.QuoteMeta(search)), timeout)
	return err
}

// ExpectRE returns an error if the given regular expression is not found in
// QEMU's serial output within `DefaultTimeout`. The matched string is
// returned.
func (q *QEMU) ExpectRE(pattern *regexp.Regexp) (string, error) {
	return q.ExpectRETimeout(pattern, DefaultTimeout)
}

// ExpectRETimeout returns an error if the given regular expression is not
// found in QEMU's serial output within the given timeout. The matched string
// is returned.
func (q *QEMU) ExpectRETimeout(pattern *regexp.Regexp, timeout time.Duration) (string, error) {
	scaled := time.Duration(float64(timeout) * TimeoutMultiplier)
	str, _, err := q.gExpect.Expect(pattern, scaled)
	return str, err
}
