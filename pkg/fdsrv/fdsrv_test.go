// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fdsrv

import (
	"os"
	"syscall"
	"testing"
	"time"
)

// Returns an fd of the read side of a pipe that has the value 'x' written in
// it.
func xPipe(t *testing.T) int {
	var p [2]int
	if err := syscall.Pipe(p[:]); err != nil {
		t.Fatal("pipe:", err)
	}
	if _, err := syscall.Write(p[1], []byte("x")); err != nil {
		t.Fatal("write:", err)
	}
	if err := syscall.Close(p[1]); err != nil {
		t.Fatal("close:", err)
	}
	return p[0]
}

// Returns an *Server, serving an xPipe with "some_nonce"
func allocPipeFDs(t *testing.T, options ...func(*Server) error) *Server {
	fd := xPipe(t)
	fds, err := NewServer(fd, "some_nonce", options...)
	if err != nil {
		t.Fatal("alloc:", err)
	}
	if err := syscall.Close(fd); err != nil {
		t.Fatal("close:", err)
	}
	return fds
}

// Read a string from an fd
func readString(t *testing.T, fd int) string {
	buf := make([]byte, 128)
	n, err := syscall.Read(fd, buf)
	if err != nil {
		t.Fatal("read:", err)
	}
	return string(buf[:n])
}

// Gets a shared fd, makes sure we can read "x" from it
func testSharedOK(t *testing.T, udspath, nonce string) {
	sfd, err := GetSharedFD(udspath, nonce)
	if err != nil {
		t.Error("getsharedfd:", err)
	}
	got := readString(t, sfd)
	if got != "x" {
		t.Errorf("expected x, got %s", got)
	}
	if err := syscall.Close(sfd); err != nil {
		t.Error("close:", err)
	}
}

func TestPassFD(t *testing.T) {
	fds := allocPipeFDs(t, WithServeOnce())

	serveErr := make(chan error)
	go func() {
		serveErr <- fds.Serve()
	}()

	testSharedOK(t, fds.UDSPath(), "some_nonce")

	fds.Close()
	if err := <-serveErr; err != nil {
		t.Errorf("Serve: %v", err)
	}
}

func TestBadNonce(t *testing.T) {
	fds := allocPipeFDs(t, WithServeOnce())

	serveErr := make(chan error)
	go func() {
		serveErr <- fds.Serve()
	}()

	sfd, err := GetSharedFD(fds.UDSPath(), "bad_nonce")
	if err == nil {
		t.Errorf("should have failed, but got sfd %d", sfd)
	}
	fds.Close()
	if err := <-serveErr; err != nil {
		t.Errorf("Serve: %v", err)
	}
}

func TestBadSubsetNonce(t *testing.T) {
	fds := allocPipeFDs(t, WithServeOnce())

	serveErr := make(chan error)
	go func() {
		serveErr <- fds.Serve()
	}()

	sfd, err := GetSharedFD(fds.UDSPath(), "some_non")
	if err == nil {
		t.Errorf("should have failed, but got sfd %d", sfd)
	}
	fds.Close()
	if err := <-serveErr; err != nil {
		t.Errorf("Serve: %v", err)
	}
}

func TestBadEmptyNonce(t *testing.T) {
	fds := allocPipeFDs(t, WithServeOnce())

	serveErr := make(chan error)
	go func() {
		serveErr <- fds.Serve()
	}()

	sfd, err := GetSharedFD(fds.UDSPath(), "")
	if err == nil {
		t.Errorf("should have failed, but got sfd %d", sfd)
	}

	fds.Close()
	if err := <-serveErr; err != nil {
		t.Errorf("Serve: %v", err)
	}
}

func TestEmptyNonce(t *testing.T) {
	fds, err := NewServer(0, "")
	if err == nil {
		t.Error("should have failed to alloc")
		fds.Close()
	}
}

// Might flake, based on timing
func TestTimeoutDoesntFire(t *testing.T) {
	fds := allocPipeFDs(t, WithServeOnce(), WithTimeout(time.Second))

	serveErr := make(chan error)
	go func() {
		serveErr <- fds.Serve()
	}()

	testSharedOK(t, fds.UDSPath(), "some_nonce")

	fds.Close()
	if err := <-serveErr; err != nil {
		t.Errorf("Serve: %v", err)
	}
}

// Might flake, based on timing
func TestTimeoutFires(t *testing.T) {
	fds := allocPipeFDs(t, WithServeOnce(), WithTimeout(100*time.Millisecond))

	serveErr := make(chan error)
	go func() {
		serveErr <- fds.Serve()
	}()

	time.Sleep(time.Millisecond * 1000)

	sfd, err := GetSharedFD(fds.UDSPath(), "some_nonce")
	if err == nil {
		t.Errorf("should have timed out, but got sfd %d", sfd)
	}
	fds.Close()
	if err := <-serveErr; err != nil && !os.IsTimeout(err) {
		t.Errorf("Serve: %v", err)
	}
}

func TestWaitTimeout(t *testing.T) {
	fds := allocPipeFDs(t, WithServeOnce(),
		WithTimeout(time.Millisecond*10))

	err := fds.Serve()
	if err == nil || !os.IsTimeout(err) {
		t.Error("expected timeout:", err)
	}
	fds.Close()
}

func TestMultiServe(t *testing.T) {
	fds := allocPipeFDs(t, WithTimeout(time.Second*5))

	serveErr := make(chan error)
	go func() {
		serveErr <- fds.Serve()
	}()

	testSharedOK(t, fds.UDSPath(), "some_nonce")
	// The second reader won't see 'x', the pipe was already drained
	sfd, err := GetSharedFD(fds.UDSPath(), "some_nonce")
	if err != nil {
		t.Error("getsharedfd:", err)
	}
	if err := syscall.Close(sfd); err != nil {
		t.Error("close:", err)
	}

	fds.Close()
	if err := <-serveErr; err != nil && !os.IsTimeout(err) {
		t.Errorf("Serve: %v", err)
	}
}
