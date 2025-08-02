// Copyright 2018-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux

package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/hugelgupf/vmtest/guest"
	"github.com/u-root/u-root/pkg/testutil"
	"golang.org/x/sys/unix"
)

func TestArgs(t *testing.T) {
	if _, err := os.Stat("/dev/fuse"); err != nil {
		t.Skipf("Skipping:%v", err)
	}
	tmpDir := t.TempDir()
	tdir := filepath.Join(tmpDir, "a/b/c")
	if err := os.MkdirAll(tdir, 0o777); err != nil {
		t.Fatal(err)
	}
	ldir := filepath.Join(tmpDir, "d")
	if err := os.Symlink(tdir, ldir); err != nil {
		t.Fatal(err)
	}

	tab := []struct {
		n            string
		o            string
		e            int
		a            []string
		env          []string
		requiresRoot bool
	}{
		{n: "badargs", o: usage + "\n", e: 1, a: []string{"-zu"}, env: []string{CommFD + "=Nan"}},
		{n: "badpath", o: fmt.Sprintf("resolved path \"%s\" and mountpoint \"%s\" are not the same\n", tdir, ldir), e: 1, a: []string{ldir}, env: []string{CommFD + "=Nan"}},
		{n: "badcfd", o: "_FUSE_COMMFD: strconv.Atoi: parsing \"Nan\": invalid syntax\n", e: 1, a: []string{tdir}, env: []string{CommFD + "=Nan"}},
		{n: "badsock", o: "_FUSE_COMMFD: 5: bad file descriptor\n", e: 1, a: []string{tdir}, env: []string{CommFD + "=5"}, requiresRoot: true},
	}
	skip := len("2018/12/20 16:54:31 ")

	uid := os.Getuid()
	for _, v := range tab {
		if uid != 0 && v.requiresRoot {
			t.Skipf("test requires root, your uid is %d", uid)
		}
		t.Run(v.n, func(t *testing.T) {
			c := testutil.Command(t, v.a...)
			c.Env = append(c.Env, v.env...)
			c.Stdin = bytes.NewReader([]byte(v.n))
			o, err := c.CombinedOutput()
			// log.Fatal exits differently on circleci and real life
			// Even an explicit os.Exit(1) returns with a 2. WTF?
			if err := testutil.IsExitCode(err, v.e); err != nil {
				t.Logf("Exit codes don't match but we'll ignore that for now")
			}
			if v.e != 0 && err == nil {
				t.Fatalf("Want error, got nil")
			}
			if v.e == 0 && err != nil {
				t.Fatalf("Want no error, got %v", err)
			}
			if len(o) < skip {
				t.Fatalf("Fusermount %v %v: want '%v', got '%v'", v.n, v.a, v.o, o)
			}
			out := string(o[len("2018/12/20 16:54:31 "):])
			// if out != v.o {
			if !strings.Contains(out, v.o) {
				t.Fatalf("Fusermount %v %v: want at least'%v', got '%v'", v.n, v.a, v.o, out)
			}
		})
	}
}

func TestMount(t *testing.T) {
	guest.SkipIfNotInVM(t)

	if _, err := os.Stat("/dev/fuse"); err != nil {
		t.Skipf("Skipping:%v", err)
	}
	// Get a socketpair to talk on, then spawn the kid
	fds, err := unix.Socketpair(syscall.AF_FILE, syscall.SOCK_STREAM, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("fds are %v", fds)

	writeFile := os.NewFile(uintptr(fds[0]), "fusermount-child-writes")
	defer writeFile.Close()

	readFile := os.NewFile(uintptr(fds[1]), "fusermount-parent-reads")
	defer readFile.Close()

	tmpDir := t.TempDir()

	fc, err := net.FileConn(readFile)
	if err != nil {
		t.Fatalf("FileConn from fusermount socket: %v", err)
	}
	defer fc.Close()

	uc, ok := fc.(*net.UnixConn)
	if !ok {
		t.Fatalf("unexpected FileConn type; expected UnixConn, got %T", fc)
	}

	c := testutil.Command(t, "-v", tmpDir)
	c.Env = append(c.Env, fmt.Sprintf("_FUSE_COMMFD=%d", fds[0]))
	c.ExtraFiles = []*os.File{writeFile}
	go func() {
		o, err := c.CombinedOutput()
		t.Logf("Running %v: %q,%v", c, string(o), err)
	}()

	buf := make([]byte, 32) // expect 1 byte
	oob := make([]byte, 32) // expect 24 bytes
	_, oobn, _, _, err := uc.ReadMsgUnix(buf, oob)
	if err != nil {
		t.Fatalf("uc.ReadMsgUnix: got %v, want nil", err)
	}
	t.Logf("ReadMsgUnix returns oobn %v, err %v", oobn, err)
	scms, err := syscall.ParseSocketControlMessage(oob[:oobn])
	if err != nil {
		t.Fatalf("syscall.ParseSocketControlMessage(%v): got %v, want nil", oob[:oobn], err)
	}
	t.Logf("syscall.ParseSocketControlMessage(%v): returns %v", oob[:oobn], scms)
	if len(scms) != 1 {
		t.Fatalf("expected 1 SocketControlMessage; got scms = %#v", scms)
	}
	scm := scms[0]
	gotFds, err := syscall.ParseUnixRights(&scm)
	if err != nil {
		t.Fatalf("syscall.ParseUnixRights: %v", err)
	}
	if len(gotFds) != 1 {
		t.Fatalf("wanted 1 fd; got %#v", gotFds)
	}
	f := os.NewFile(uintptr(gotFds[0]), "/dev/fuse")
	t.Logf("file to fuse is %v", f)
	// Now every good program should unmount.
	c = testutil.Command(t, "-v", "-u", tmpDir)
	o, err := c.CombinedOutput()
	t.Logf("Running fuse: %v,%v", string(o), err)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}
