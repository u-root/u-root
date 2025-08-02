// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syscallfilter

import (
	"bytes"
	"os"
	"strconv"
	"strings"
	"syscall"
	"testing"

	"github.com/u-root/u-root/pkg/strace"
)

// If we are being traced, none of this will work.
// Be conservative: if anything at all fails, just return
// true, i.e. assume we're being traced somehow..
// Yep, this is a kludge.
func traced() bool {
	b, err := os.ReadFile("/proc/self/status")
	if err != nil {
		return true
	}

	s := strings.SplitN(string(b), "TracerPid:\t", 2)
	// There should be two pieces, the bit before TracerPid:
	// and the bit after, included the value of TracerPid.
	// If we can't find this then we have no way to know.
	// Assume traced: somebody could be messing with us.
	// log.Printf("split %s", s)
	if len(s) < 2 {
		return true
	}

	tracerPid, err := strconv.Atoi(strings.Fields(s[1])[0])

	if err != nil || tracerPid != 0 {
		return true
	}

	return false
}

func TestEventMap(t *testing.T) {
	tests := []struct {
		v    string
		find string
		err  error
	}{
		{"NewChild,error,-1", "NewChild", nil},
		{"E.*timeof.*,error,-1", "Egettimeofday", nil},
	}
	for _, tt := range tests {
		c := Command("date")
		err := c.AddActions(tt.v)
		if err != nil && tt.err == nil {
			t.Errorf("%v.AddActions(%s): err %v != %v", c, tt.v, err, tt.err)
			continue
		}
		fe := findEvent(c.events, tt.find)
		if fe == nil {
			t.Errorf(`findEvent(%v, %v): is nil`, c, tt.find)
		}
	}
}

func TestNoActions(t *testing.T) {
	if traced() {
		t.Skipf("Skipping, we're being traced already")
	}

	c := Command("echo", "hi")
	var stdout, stderr bytes.Buffer
	c.Stdout, c.Stderr = &stdout, &stderr

	if err := c.Run(); err != nil {
		t.Fatalf(`%v.Run(), "echo", "hi"): %v != nil`, c, err)
	}
	t.Logf("stdout: %q, stderr: %q", stdout.String(), stderr.String())
	if stdout.String() != "hi\n" {
		t.Errorf("stdout: string is %q, not %q", stdout.String(), "hi\n")
	}
	if len(stderr.String()) != 0 {
		t.Errorf("stderr.String: got %q, want %q", stderr.String(), "")
	}
}

func TestNoErrorExit(t *testing.T) {
	if traced() {
		t.Skipf("Skipping, we're being traced already")
	}

	c := Command("date")
	var stdout, stderr bytes.Buffer
	c.Stdout, c.Stderr = &stdout, &stderr

	if err := c.Run(); err != nil {
		t.Fatalf(`%v.Run(): %v != nil`, c, err)
	}
	t.Logf("stdout: %q, stderr: %q", stdout.String(), stderr.String())
	if len(stdout.String()) == 0 {
		t.Errorf("stdout.String: got %q, want output", "")
	}
	if len(stderr.String()) != 0 {
		t.Errorf("stderr.String: got %q, want %q", stderr.String(), "")
	}
}

func TestErrorExitLog(t *testing.T) {
	if traced() {
		t.Skipf("Skipping, we're being traced already")
	}

	filter := "E.*open.*,error,-2"
	c := Command("date")
	var stdout, stderr bytes.Buffer
	c.Stdout, c.Stderr = &stdout, &stderr
	err := c.AddActions(filter)
	if err != nil {
		t.Fatalf("addactions(%s): err %v != %v", filter, err, nil)
	}
	var b bytes.Buffer
	c.Log = &b
	if err := c.Run(); err == nil {
		t.Fatalf(`%v.Run(): nil != %v`, c, "exit status 1")
	}
	t.Logf("stdout: %q, stderr: %q", stdout.String(), stderr.String())
	if b.Len() == 0 {
		t.Errorf("%v.Run(): Log is zero, not > 0", c)
	}
	t.Logf("Log: %v", b.String())
}

func TestErrorExitLogAll(t *testing.T) {
	if traced() {
		t.Skipf("Skipping, we're being traced already")
	}

	filter := "E.*open.*,error,-2"
	c := Command("date")
	var stdout, stderr bytes.Buffer
	c.Stdout, c.Stderr = &stdout, &stderr
	if err := c.AddActions(filter); err != nil {
		t.Fatalf("addactions(%s): err %v != %v", filter, err, nil)
	}
	var b bytes.Buffer
	c.Log = &b
	if err := c.Run(); err == nil {
		t.Fatalf(`%v.Run(): nil != %v`, c, "exit status 1")
	}
	t.Logf("stdout: %q, stderr: %q", stdout.String(), stderr.String())
	if b.Len() == 0 {
		t.Errorf("%v.Run(): Log is zero, not > 0", c)
	}
	t.Logf("Log: %v", b.String())
	// OK, now run it again, with LogAllActions set; the new log should be longer.
	c = Command("date")
	c.Stdout, c.Stderr = &stdout, &stderr
	if err := c.AddActions("E.*read.*,log,0", "E.*exit.*,error,-2"); err != nil {
		t.Fatalf("addactions(%s): err %v != %v", filter, err, nil)
	}
	var b2 bytes.Buffer
	c.Log = &b2
	if err := c.Run(); err == nil {
		t.Fatalf(`%v.Run(): nil != %v`, c, "exit status 1")
	}
	t.Logf("stdout: %q, stderr: %q", stdout.String(), stderr.String())
	if b2.Len() == 0 {
		t.Errorf("%v.Run(): Log is zero bytes, not > 0", c)
	}
	if b.Len() >= b2.Len() {
		t.Errorf("%v.Run() with LogAllRecords: Log is %d bytes, not > %d", c, b2.Len(), b.Len())
	}
	t.Logf("Log: %v", b2.String())
}

// This is a simple example of how you might test the reboot command
// without being root and without using qemu.
// It is disabled because it's just a bit too dangerous.
func testRebootOK(t *testing.T) {
	filter := ".*reboot,error,0"
	c := Command("reboot", "-f", "-f")
	var stdout, stderr bytes.Buffer
	c.Stdout, c.Stderr = &stdout, &stderr
	err := c.AddActions(filter)
	if err != nil {
		t.Fatalf("addactions(%s): err %v != %v", filter, err, nil)
	}

	if err := c.Run(); err == nil {
		t.Fatalf(`%v.Run(): nil != %v`, c, "exit status 1")
	}
	t.Logf("stdout: %q, stderr: %q", stdout.String(), stderr.String())
}

func testRebootFail(t *testing.T) {
	if traced() {
		t.Skipf("Skipping, we're being traced already")
	}

	filter := ".*reboot,error,-1"
	c := Command("reboot", "-f", "-f")
	var stdout, stderr bytes.Buffer
	c.Stdout, c.Stderr = &stdout, &stderr
	err := c.AddActions(filter)
	if err != nil {
		t.Fatalf("addactions(%s): err %v != %v", filter, err, nil)
	}

	if err := c.Run(); err == nil {
		t.Fatalf(`%v.Run(): nil != %v`, c, "exit status 1")
	}
	t.Logf("stdout: %q, stderr: %q", stdout.String(), stderr.String())
}

func TestEventName(t *testing.T) {
	if traced() {
		t.Skipf("Skipping, we're being traced already")
	}

	for _, test := range []struct {
		r *strace.TraceRecord
		n string
	}{
		{r: &strace.TraceRecord{Event: strace.SyscallEnter, Syscall: &strace.SyscallEvent{Sysno: syscall.SYS_READ}}, n: "Eread"},
		{r: &strace.TraceRecord{Event: strace.SyscallExit, Syscall: &strace.SyscallEvent{Sysno: syscall.SYS_READ}}, n: "Xread"},
		{r: &strace.TraceRecord{Event: strace.SyscallEnter, Syscall: &strace.SyscallEvent{Sysno: 0xabcd}}, n: "Eabcd"},
		{r: &strace.TraceRecord{Event: strace.SyscallExit, Syscall: &strace.SyscallEvent{Sysno: 0xbcde}}, n: "Xbcde"},
		{r: &strace.TraceRecord{Event: strace.SignalExit}, n: "SignalExit"},
		{r: &strace.TraceRecord{Event: strace.Exit}, n: "Exit"},
		{r: &strace.TraceRecord{Event: strace.SignalStop}, n: "SignalStop"},
		{r: &strace.TraceRecord{Event: strace.NewChild}, n: "NewChild"},
	} {
		n := eventName(test.r)
		// There's a problem testing the system call name: linux system call numbers to names
		// vary, widely, across architectures. But read at 0 seems safe, so ...
		if test.n != n {
			t.Errorf("eventName(%v): %v != %v ", test, n, test.n)
		}
	}
}
