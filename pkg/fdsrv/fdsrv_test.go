// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux

package fdsrv

import (
	"context"
	"errors"
	"net"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"

	"golang.org/x/sys/unix"
)

// Returns an fd of the read side of a pipe that has the value 'x' written in
// it.
func xPipe(t *testing.T) int {
	t.Helper()
	var p [2]int
	if err := syscall.Pipe2(p[:], syscall.O_CLOEXEC); err != nil {
		t.Fatalf("Pipe2: got %v, want nil", err)
	}
	if _, err := syscall.Write(p[1], []byte("x")); err != nil {
		t.Fatalf("Write: got %v, want nil", err)
	}
	if err := syscall.Close(p[1]); err != nil {
		t.Fatalf("Close: got %v, want nil", err)
	}
	return p[0]
}

// Returns a *Server, serving an xPipe with "some_nonce".
func allocPipeFDs(t *testing.T, options ...func(*Server) error) *Server {
	t.Helper()
	fd := xPipe(t)
	fds, err := NewServer(fd, "some_nonce", options...)
	if err != nil {
		t.Fatalf("NewServer: got %v, want nil", err)
	}
	if err := syscall.Close(fd); err != nil {
		t.Fatalf("Close: got %v, want nil", err)
	}
	return fds
}

// Read a string from an fd.
func readString(t *testing.T, fd int) string {
	t.Helper()
	buf := make([]byte, 128)
	n, err := syscall.Read(fd, buf)
	if err != nil {
		t.Fatalf("Read: got %v, want nil", err)
	}
	return string(buf[:n])
}

// Gets a shared fd, makes sure we can read "x" from it.
func testSharedOK(ctx context.Context, t *testing.T, udspath, nonce string) {
	t.Helper()
	sfd, err := GetSharedFD(ctx, udspath, nonce)
	if err != nil {
		t.Fatalf("GetSharedFD(%q, %q): got %v, want nil", udspath, nonce, err)
	}
	got := readString(t, sfd)
	if got != "x" {
		t.Errorf("readString(%d): got %q, want %q", sfd, got, "x")
	}
	if err := syscall.Close(sfd); err != nil {
		t.Errorf("Close(%d): got %v, want nil", sfd, err)
	}
}

func TestPassFD(t *testing.T) {
	ctx := t.Context()
	fds := allocPipeFDs(t, WithServeOnce())

	serveErr := make(chan error)
	go func() {
		serveErr <- fds.Serve(ctx)
	}()

	testSharedOK(ctx, t, fds.UDSPath(), "some_nonce")

	if err := fds.Close(); err != nil {
		t.Errorf("Close: got %v, want nil", err)
	}
	if err := <-serveErr; err != nil {
		t.Errorf("Serve: got %v, want nil", err)
	}
}

func TestBadNonce(t *testing.T) {
	ctx := t.Context()
	fds := allocPipeFDs(t, WithServeOnce())

	serveErr := make(chan error)
	go func() {
		serveErr <- fds.Serve(ctx)
	}()

	sfd, err := GetSharedFD(ctx, fds.UDSPath(), "bad_nonce")
	if err == nil {
		t.Errorf("GetSharedFD(bad_nonce): got (%d, nil), want (-1, non-nil error)", sfd)
	}
	// Server should remain alive after a bad nonce and still serve a valid
	// client.
	testSharedOK(ctx, t, fds.UDSPath(), "some_nonce")

	if err := fds.Close(); err != nil {
		t.Errorf("Close: got %v, want nil", err)
	}
	if err := <-serveErr; err != nil {
		t.Errorf("Serve: got %v, want nil", err)
	}
}

func TestBadSubsetNonce(t *testing.T) {
	ctx := t.Context()
	fds := allocPipeFDs(t, WithServeOnce())

	serveErr := make(chan error)
	go func() {
		serveErr <- fds.Serve(ctx)
	}()

	sfd, err := GetSharedFD(ctx, fds.UDSPath(), "some_non")
	if err == nil {
		t.Errorf("GetSharedFD(some_non): got (%d, nil), want (-1, non-nil error)", sfd)
	}
	sfd, err = GetSharedFD(ctx, fds.UDSPath(), "some_nonce_extra")
	if err == nil {
		t.Errorf("GetSharedFD(some_nonce_extra): got (%d, nil), want (-1, non-nil error)", sfd)
	}
	if err := fds.Close(); err != nil {
		t.Errorf("Close: got %v, want nil", err)
	}
	if err := <-serveErr; err != nil {
		t.Errorf("Serve: got %v, want nil", err)
	}
}

func TestBadEmptyNonce(t *testing.T) {
	ctx := t.Context()
	fds := allocPipeFDs(t, WithServeOnce())

	serveErr := make(chan error)
	go func() {
		serveErr <- fds.Serve(ctx)
	}()

	sfd, err := GetSharedFD(ctx, fds.UDSPath(), "")
	if !errors.Is(err, ErrEmptyNonce) {
		t.Errorf("GetSharedFD(\"\"): got (%d, %v), want (-1, %v)", sfd, err, ErrEmptyNonce)
	}

	// Also verify that a client connecting and closing without sending any
	// bytes does not tear down the server.
	var d net.Dialer
	conn, err := d.DialContext(ctx, "unix", fds.UDSPath())
	if err != nil {
		t.Fatalf("DialContext: got %v, want nil", err)
	}
	conn.Close()

	testSharedOK(ctx, t, fds.UDSPath(), "some_nonce")

	if err := fds.Close(); err != nil {
		t.Errorf("Close: got %v, want nil", err)
	}
	if err := <-serveErr; err != nil {
		t.Errorf("Serve: got %v, want nil", err)
	}
}

func TestNewServerErrors(t *testing.T) {
	if fds, err := NewServer(0, ""); !errors.Is(err, ErrEmptyNonce) {
		t.Errorf("NewServer(0, \"\"): got (%v, %v), want (nil, %v)", fds, err, ErrEmptyNonce)
	}
	if fds, err := NewServer(-1, "some_nonce"); err == nil {
		fds.Close()
		t.Errorf("NewServer(-1, \"some_nonce\"): got (%v, nil), want (nil, non-nil error)", fds)
	}
	optErr := errors.New("bad option")
	if fds, err := NewServer(0, "some_nonce", func(*Server) error { return optErr }); !errors.Is(err, optErr) {
		t.Errorf("NewServer with failing option: got (%v, %v), want (nil, %v)", fds, err, optErr)
	}
}

func TestTimeoutDoesntFire(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), time.Second)
	defer cancel()
	fds := allocPipeFDs(t, WithServeOnce())

	serveErr := make(chan error)
	go func() {
		serveErr <- fds.Serve(ctx)
	}()

	testSharedOK(ctx, t, fds.UDSPath(), "some_nonce")

	if err := fds.Close(); err != nil {
		t.Errorf("Close: got %v, want nil", err)
	}
	if err := <-serveErr; err != nil {
		t.Errorf("Serve: got %v, want nil", err)
	}
}

func TestTimeoutFires(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 50*time.Millisecond)
	defer cancel()
	fds := allocPipeFDs(t, WithServeOnce())

	serveErr := make(chan error)
	go func() {
		serveErr <- fds.Serve(ctx)
	}()

	if err := <-serveErr; !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Serve: got %v, want %v", err, context.DeadlineExceeded)
	}

	sfd, err := GetSharedFD(t.Context(), fds.UDSPath(), "some_nonce")
	if err == nil {
		t.Errorf("GetSharedFD after timeout: got (%d, nil), want (-1, non-nil error)", sfd)
	}
	if err := fds.Close(); err != nil {
		t.Errorf("Close: got %v, want nil", err)
	}
}

func TestWaitTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Millisecond)
	defer cancel()
	fds := allocPipeFDs(t, WithServeOnce())

	err := fds.Serve(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Serve: got %v, want %v", err, context.DeadlineExceeded)
	}
	if err := fds.Close(); err != nil {
		t.Errorf("Close: got %v, want nil", err)
	}
}

func TestCanceledContextAndStalledPeers(t *testing.T) {
	// Serve with an already-canceled context.
	canceledCtx, cancel := context.WithCancel(t.Context())
	cancel()
	fds := allocPipeFDs(t, WithServeOnce())
	if err := fds.Serve(canceledCtx); !errors.Is(err, context.Canceled) {
		t.Errorf("Serve(canceledCtx): got %v, want %v", err, context.Canceled)
	}

	// Client connects and stalls without writing nonce until server ctx
	// expires inside handleConnection.
	srvCtx, srvCancel := context.WithTimeout(t.Context(), 50*time.Millisecond)
	defer srvCancel()
	fds2 := allocPipeFDs(t, WithServeOnce())
	serveErr := make(chan error)
	go func() {
		serveErr <- fds2.Serve(srvCtx)
	}()
	var d net.Dialer
	stalledConn, err := d.DialContext(t.Context(), "unix", fds2.UDSPath())
	if err != nil {
		t.Fatalf("DialContext: got %v, want nil", err)
	}
	defer stalledConn.Close()
	if err := <-serveErr; !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Serve with stalled client: got %v, want %v", err, context.DeadlineExceeded)
	}

	// Client GetSharedFD times out while waiting for a stalled server, and
	// rejects a server sending multiple FDs in SCM_RIGHTS.
	ln, err := net.ListenUnix("unix", &net.UnixAddr{Net: "unix"})
	if err != nil {
		t.Fatalf("ListenUnix: got %v, want nil", err)
	}
	defer ln.Close()

	p := xPipe(t)
	defer syscall.Close(p)

	go func() {
		// First connection: read nonce and stall until client times out.
		if c1, err := ln.AcceptUnix(); err == nil {
			buf := make([]byte, 32)
			c1.Read(buf)
			time.Sleep(100 * time.Millisecond)
			c1.Close()
		}
		// Second connection: send two FDs in SCM_RIGHTS.
		if c2, err := ln.AcceptUnix(); err == nil {
			defer c2.Close()
			buf := make([]byte, 32)
			c2.Read(buf)
			oob := syscall.UnixRights(p, p)
			c2.WriteMsgUnix([]byte("TWO FDS"), oob, nil)
		}
	}()

	cliCtx, cliCancel := context.WithTimeout(t.Context(), 20*time.Millisecond)
	defer cliCancel()
	if sfd, err := GetSharedFD(cliCtx, ln.Addr().String(), "some_nonce"); !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("GetSharedFD(stalled server): got (%d, %v), want (-1, %v)", sfd, err, context.DeadlineExceeded)
	}

	if sfd, err := GetSharedFD(t.Context(), ln.Addr().String(), "some_nonce"); !errors.Is(err, ErrNotOneUnixRights) {
		t.Errorf("GetSharedFD(two FDs): got (%d, %v), want (-1, %v)", sfd, err, ErrNotOneUnixRights)
	}
}

func TestMultiServe(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	fds := allocPipeFDs(t)

	serveErr := make(chan error)
	go func() {
		serveErr <- fds.Serve(ctx)
	}()

	testSharedOK(ctx, t, fds.UDSPath(), "some_nonce")
	// The second reader won't see 'x'; the pipe was already drained.
	sfd, err := GetSharedFD(ctx, fds.UDSPath(), "some_nonce")
	if err != nil {
		t.Fatalf("GetSharedFD: got %v, want nil", err)
	}
	if err := syscall.Close(sfd); err != nil {
		t.Errorf("Close(%d): got %v, want nil", sfd, err)
	}

	if err := fds.Close(); err != nil {
		t.Errorf("Close: got %v, want nil", err)
	}
	if err := <-serveErr; err != nil {
		t.Errorf("Serve: got %v, want nil", err)
	}
}

func TestLongNonce(t *testing.T) {
	ctx := t.Context()
	longNonce := strings.Repeat("n", 8192)
	fd := xPipe(t)
	fds, err := NewServer(fd, longNonce, WithServeOnce())
	if err != nil {
		t.Fatalf("NewServer: got %v, want nil", err)
	}
	if err := syscall.Close(fd); err != nil {
		t.Fatalf("Close(%d): got %v, want nil", fd, err)
	}

	serveErr := make(chan error)
	go func() {
		serveErr <- fds.Serve(ctx)
	}()

	testSharedOK(ctx, t, fds.UDSPath(), longNonce)

	if err := <-serveErr; err != nil {
		t.Errorf("Serve: got %v, want nil", err)
	}
}

func TestCloexecAndNoFDLeak(t *testing.T) {
	countOpenFDs := func(t *testing.T) int {
		t.Helper()
		entries, err := os.ReadDir("/proc/self/fd")
		if err != nil {
			t.Fatalf("ReadDir(/proc/self/fd): got %v, want nil", err)
		}
		return len(entries)
	}

	before := countOpenFDs(t)

	ctx := t.Context()
	fds := allocPipeFDs(t, WithServeOnce())

	flags, err := unix.FcntlInt(uintptr(fds.dupedFD), unix.F_GETFD, 0)
	if err != nil {
		t.Fatalf("FcntlInt(F_GETFD) on dupedFD: got %v, want nil", err)
	}
	if flags&unix.FD_CLOEXEC == 0 {
		t.Errorf("dupedFD flags: got %#x, want FD_CLOEXEC (%#x) set", flags, unix.FD_CLOEXEC)
	}

	serveErr := make(chan error)
	go func() {
		serveErr <- fds.Serve(ctx)
	}()

	sfd, err := GetSharedFD(ctx, fds.UDSPath(), "some_nonce")
	if err != nil {
		t.Fatalf("GetSharedFD: got %v, want nil", err)
	}
	flags, err = unix.FcntlInt(uintptr(sfd), unix.F_GETFD, 0)
	if err != nil {
		t.Fatalf("FcntlInt(F_GETFD) on sfd: got %v, want nil", err)
	}
	if flags&unix.FD_CLOEXEC == 0 {
		t.Errorf("sfd flags: got %#x, want FD_CLOEXEC (%#x) set", flags, unix.FD_CLOEXEC)
	}
	if err := syscall.Close(sfd); err != nil {
		t.Fatalf("Close(%d): got %v, want nil", sfd, err)
	}
	if err := <-serveErr; err != nil {
		t.Fatalf("Serve: got %v, want nil", err)
	}
	if err := fds.Close(); err != nil {
		t.Fatalf("Close: got %v, want nil", err)
	}

	after := countOpenFDs(t)
	if after != before {
		t.Errorf("open FD count after test: got %d, want %d", after, before)
	}
}
