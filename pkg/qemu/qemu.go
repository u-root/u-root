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
	"io"
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

// QEMU is filled and pass to `Start()`.
type QEMU struct {
	// Path to the bzImage kernel
	Kernel string

	// Path to the initramfs
	InitRAMFS string

	// Extra kernel arguments
	KernelArgs string

	// Extra QEMU arguments
	ExtraArgs []string

	// Where to send serial output.
	SerialOutput io.WriteCloser

	gExpect *expect.GExpect
}

// CmdLine returns the command line arguments used to start QEMU. These
// arguments are derived from the given QEMU struct.
func (q *QEMU) CmdLine() []string {
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
	if q.Kernel != "" {
		args = append(args, "-kernel", q.Kernel)
		args = append(args, "-append", "console=ttyS0 earlyprintk=ttyS0")
		if q.KernelArgs != "" {
			args[len(args)-1] += " " + q.KernelArgs
		}
	}
	if q.InitRAMFS != "" {
		args = append(args, "-initrd", q.InitRAMFS)
	}

	return append(args, q.ExtraArgs...)
}

// CmdLineQuoted quotes any of QEMU's command line arguments containing a space
// so it is easy to copy-n-paste into a shell for debugging.
func (q *QEMU) CmdLineQuoted() (cmdline string) {
	args := q.CmdLine()
	for i, v := range q.CmdLine() {
		if strings.ContainsAny(v, " \t\n") {
			args[i] = "'" + v + "'"
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
	q.gExpect, _, err = expect.SpawnWithArgs(q.CmdLine(), -1,
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
	return q.ExpectTimeout(search, DefaultTimeout)
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
