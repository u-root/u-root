// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux

// Package fdsrv serves a file descriptor over a Linux Unix domain socket (UDS)
// when presented with a nonce.
//
// You must pass the socket path and nonce to the client via some out-of-band
// mechanism, such as gRPC or a bash script.
//
// Notes:
//   - Uses the Linux UDS abstract namespace.
//   - Picks its own path in the abstract namespace for the socket.
//   - Shared FDs are essentially duped, and they point to the same struct file:
//     they share offsets and whatnot.
//
// Options:
//   - WithServeOnce: serve once and shut down (default is forever).
//
// Usage Server:
//
//	fds, err := NewServer(fdToShare, "some_nonce", WithServeOnce())
//	path := fds.UDSPath()
//
//	// Pass path and some_nonce to the client via an out-of-band mechanism.
//
//	fds.Serve(ctx) // Blocks until the server is done or ctx is canceled.
//	fds.Close()
//
// Usage Client:
//
//	sfd, err := GetSharedFD(ctx, udsPath, "some_nonce")
package fdsrv

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"
	"syscall"

	"golang.org/x/sys/unix"
)

var (
	// ErrTruncatedWrite indicates a short write on the UDS.
	ErrTruncatedWrite = errors.New("truncated write")
	// ErrEmptyNonce indicates that an empty nonce was provided.
	ErrEmptyNonce = errors.New("nonce must not be empty")
	// ErrMissingSCM indicates that expected socket control message was missing.
	ErrMissingSCM = errors.New("missing socket control message")
	// ErrNotOneUnixRights indicates that SCM_RIGHTS did not contain one FD.
	ErrNotOneUnixRights = errors.New("expected exactly one unix rights")
)

// Server serves a file descriptor over a Linux UDS in the abstract namespace.
type Server struct {
	dupedFD   int
	nonce     string
	listener  *net.UnixListener
	serveOnce bool
	closeOnce sync.Once
	closeErr  error
}

// Serves the fd, returns true if successful, err for a server error.
// "false, nil" means the client was wrong, not the server.
func (fds *Server) handleConnection(ctx context.Context, uc *net.UnixConn) (bool, error) {
	defer uc.Close()
	stop := context.AfterFunc(ctx, func() {
		uc.Close()
	})
	defer stop()

	buf := make([]byte, len(fds.nonce)+1)
	n, err := uc.Read(buf)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return false, ctxErr
		}
		return false, nil
	}
	query := string(buf[:n])
	if query != fds.nonce {
		io.WriteString(uc, "BAD NONCE")
		return false, nil
	}
	oob := syscall.UnixRights(fds.dupedFD)
	good := []byte("GOOD NONCE")
	goodn, oobn, err := uc.WriteMsgUnix(good, oob, nil)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return false, ctxErr
		}
		return false, err
	}
	if goodn != len(good) || oobn != len(oob) {
		return false, ErrTruncatedWrite
	}
	return true, nil
}

// NewServer creates a Server that shares fd when presented with nonce. Call
// Close when done.
func NewServer(fd int, nonce string, options ...func(*Server) error) (*Server, error) {
	if len(nonce) == 0 {
		return nil, ErrEmptyNonce
	}
	fds := &Server{
		nonce: nonce,
	}

	for _, op := range options {
		if err := op(fds); err != nil {
			return nil, err
		}
	}

	// An empty addr tells Linux to "autobind" to an available path in the
	// UDS abstract namespace.
	var err error
	fds.listener, err = net.ListenUnix("unix", &net.UnixAddr{Net: "unix"})
	if err != nil {
		return nil, err
	}

	// Caller could close the file while we are running. Keep our own copy
	// with FD_CLOEXEC set so child processes do not inherit it.
	r1, err := unix.FcntlInt(uintptr(fd), unix.F_DUPFD_CLOEXEC, 0)
	if err != nil {
		fds.listener.Close()
		return nil, err
	}
	fds.dupedFD = r1

	return fds, nil
}

// WithServeOnce sets the "serve once and exit" option on NewServer.
func WithServeOnce() func(*Server) error {
	return func(fds *Server) error {
		fds.serveOnce = true
		return nil
	}
}

// UDSPath returns the UDS path the server is listening on.
func (fds *Server) UDSPath() string {
	return fds.listener.Addr().String()
}

// Close closes the server. It is safe to call Close multiple times.
func (fds *Server) Close() error {
	fds.closeOnce.Do(func() {
		err1 := fds.listener.Close()
		err2 := syscall.Close(fds.dupedFD)
		fds.closeErr = errors.Join(err1, err2)
	})
	return fds.closeErr
}

// Serve serves the FD until the server finishes, Close is called, or ctx is
// done.
func (fds *Server) Serve(ctx context.Context) error {
	defer fds.Close()

	if err := ctx.Err(); err != nil {
		return err
	}
	stop := context.AfterFunc(ctx, func() {
		fds.Close()
	})
	defer stop()

	for {
		conn, err := fds.listener.AcceptUnix()
		if err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return ctxErr
			}
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			return err
		}
		succeeded, err := fds.handleConnection(ctx, conn)
		if err != nil {
			return err
		}
		if succeeded && fds.serveOnce {
			return nil
		}
	}
}

// GetSharedFD connects to udsPath, presents nonce, and returns the shared FD.
func GetSharedFD(ctx context.Context, udsPath, nonce string) (int, error) {
	// If you don't send at least a byte, the server won't recvmsg. This is
	// a Linux UDS SOCK_STREAM thing.
	if len(nonce) == 0 {
		return -1, ErrEmptyNonce
	}

	var d net.Dialer
	conn, err := d.DialContext(ctx, "unix", udsPath)
	if err != nil {
		return -1, err
	}
	defer conn.Close()
	uc := conn.(*net.UnixConn)

	stop := context.AfterFunc(ctx, func() {
		uc.Close()
	})
	defer stop()

	n, err := uc.Write([]byte(nonce))
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return -1, ctxErr
		}
		return -1, err
	}
	if n != len(nonce) {
		return -1, ErrTruncatedWrite
	}

	msg := make([]byte, 32)
	// Allocate space for up to 2 FDs (8 bytes) so that if a peer sends more
	// than 1 FD on a 32-bit architecture, the kernel does not truncate the
	// control message down to 1 FD.
	oob := make([]byte, syscall.CmsgSpace(8))
	_, oobn, _, _, err := uc.ReadMsgUnix(msg, oob)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return -1, ctxErr
		}
		return -1, err
	}
	scm, err := syscall.ParseSocketControlMessage(oob[:oobn])
	if err != nil {
		return -1, err
	}
	if len(scm) != 1 {
		return -1, ErrMissingSCM
	}
	urs, err := syscall.ParseUnixRights(&scm[0])
	if err != nil {
		return -1, err
	}
	if len(urs) != 1 {
		for _, fd := range urs {
			syscall.Close(fd)
		}
		return -1, ErrNotOneUnixRights
	}
	syscall.CloseOnExec(urs[0])
	return urs[0], nil
}
