// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integration

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/u-root/cpu/client"
	"github.com/u-root/cpu/vm"
)

// TestCPUAMD64 tests both general and specific things. The specific parts are the io and cmos commands.
// It being cheaper to use a single generated initramfs, we use the full u-root for several tests.
func TestCPUKexecAMD64(t *testing.T) {
	d := t.TempDir()
	i, err := vm.New("linux", "amd64")
	if !errors.Is(err, nil) {
		t.Fatalf("Testing kernel=linux arch=amd64: got %v, want nil", err)
	}

	// Cancel before wg.Wait(), so goroutine can exit.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// TODO: have a one-time helper that builds a full u-root image once,
	// that all tests can use.
	// TODO: for all the tests, we need start only one VM. Even for kexec,
	// since it just starts a new kernel, and we can have that kernel use
	// the initramfs that runs cpud.
	n, err := i.Uroot(d)
	if err != nil {
		t.Skipf("skipping this test as we have no uroot command")
	}

	c, err := i.CommandContext(ctx, d, n)
	if err != nil {
		t.Fatalf("starting VM: got %v, want nil", err)
	}
	c.Stdout, c.Stderr = os.Stdout, os.Stderr

	if err := i.StartVM(c); err != nil {
		t.Fatalf("starting VM: got %v, want nil", err)
	}

	cpuTmp := "/tmp/cpu"

	for _, tt := range []struct {
		cmd  string
		args []string
		ok   bool
		wait bool
	}{
		{cmd: "/bbin/kexec", args: []string{"-l", "-d", "-i", filepath.Join(cpuTmp, "initramfs"), "--loadsyscall", "-reuse-cmdline", filepath.Join(cpuTmp, "kernel")}, ok: true, wait: true},
		{cmd: "/bbin/kexec", args: []string{"-e", "-d"}, ok: true, wait: false},
	} {
		cpu, err := i.CPUCommand(tt.cmd, tt.args...)
		if err != nil {
			t.Errorf("CPUCommand: got %v, want nil", err)
			continue
		}
		client.SetVerbose(t.Logf)

		if tt.wait {
			b, err := cpu.CombinedOutput()
			if err == nil != tt.ok {
				t.Errorf("%s %s: got %v, want %v", tt.cmd, tt.args, err == nil != tt.ok, err == nil == tt.ok)
			}
			t.Logf("%q", string(b))
		} else {
			if err := cpu.Start(); err != nil {
				t.Errorf("%s %s: got %v, want %v", tt.cmd, tt.args, err == nil != tt.ok, err == nil == tt.ok)
			}
		}
	}

	t.Logf("Delay")
	time.Sleep(10 * time.Second)
	t.Logf("Try to cpu to guest")

	for _, tt := range []struct {
		cmd  string
		args []string
		ok   bool
	}{
		{cmd: "/bbin/date", args: []string{}, ok: true},
	} {
		cpu, err := i.CPUCommand(tt.cmd, tt.args...)
		if err != nil {
			t.Errorf("CPUCommand: got %v, want nil", err)
			continue
		}
		client.SetVerbose(t.Logf)

		b, err := cpu.CombinedOutput()
		if err == nil != tt.ok {
			t.Errorf("%s %s: got %v, want %v", tt.cmd, tt.args, err == nil != tt.ok, err == nil == tt.ok)
		}
		t.Logf("%q", string(b))
	}

}
