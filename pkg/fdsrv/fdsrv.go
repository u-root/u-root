// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Serves a file descriptor over an AF_UNIX socket when presented with a nonce.
//
// You must pass the socket path and nonce to the client via some out-of-band
// mechanism, such as gRPC or a bash script.
//
// Notes:
// - Uses the unix domain socket abstract namespace
// - Picks its own path in the abstract namespace for the socket.
// - Shared FDs are essentially duped, and they point to the same struct file:
// they share offsets and whatnot.
//
// Options:
// - WithServeOnce: serve once and shuts down (default is forever)
// - WithTimeout: cancel itself after a timeout (default none)
//
// Usage Server:
//
//	fds, err := NewServer(fd_to_share, "some_nonce", WithServeOnce())
//	var s path = fds.UDSPath()
//
//	// Pass path and some_nonce to the client via an out of band mechanism
//
//	fds.Serve(); // Blocks until the server is done
//	fds.Close()
//
// Usage Client:
//
//	sfd, err := GetSharedFD("uds_path", "some_nonce")
package fdsrv

import (
	"errors"
	"io"
	"net"
	"os"
	"syscall"
	"time"
)

var (
	ErrTruncatedWrite   = errors.New("truncated write")
	ErrEmptyNonce       = errors.New("nonce must not be empty")
	ErrMissingSCM       = errors.New("missing socket control message")
	ErrNotOneUnixRights = errors.New("expected exactly one unix rights")
)

type Server struct {
	dupedFD   int
	nonce     string
	listener  *net.UnixListener
	timeout   time.Duration
	serveOnce bool
}

// Serves the fd, returns true if successful, err for a server error.
// "false, nil" means the client was wrong, not the server.
func (fds *Server) handleConnection(uc *net.UnixConn) (bool, error) {
	defer uc.Close()

	buf := make([]byte, 4096)
	n, err := uc.Read(buf)
	if err != nil {
		return false, err
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
		return false, err
	}
	if goodn != len(good) || oobn != len(oob) {
		return false, ErrTruncatedWrite
	}
	return true, nil
}

// NewServer creates a server.  Close() it when you're done.
func NewServer(fd int, nonce string, options ...func(*Server) error) (*Server, error) {
	var err error
	fds := &Server{}

	if len(nonce) == 0 {
		return nil, ErrEmptyNonce
	}
	fds.nonce = nonce

	for _, op := range options {
		if err := op(fds); err != nil {
			return nil, err
		}
	}

	// An empty addr tells Linux to "autobind" to an available path in the
	// abstract unix domain socket namespace
	ua, err := net.ResolveUnixAddr("unix", "")
	if err != nil {
		return nil, err
	}
	fds.listener, err = net.ListenUnix("unix", ua)
	if err != nil {
		return nil, err
	}

	// Caller could close the file while we are running.  Keep our own copy.
	fds.dupedFD, err = syscall.Dup(int(fd))
	if err != nil {
		fds.listener.Close()
		return nil, err
	}

	return fds, nil
}

// WithTimeOut adds a timeout option to NewServer
func WithTimeout(timeout time.Duration) func(*Server) error {
	return func(fds *Server) error {
		fds.timeout = timeout
		return nil
	}
}

// WithServeOnce sets the "serve once and exit" option to NewServer
func WithServeOnce() func(*Server) error {
	return func(fds *Server) error {
		fds.serveOnce = true
		return nil
	}
}

// UDSPath returns the Unix Domain Socket path the server is listening on
func (fds *Server) UDSPath() string {
	return fds.listener.Addr().String()
}

// Close closes the server
func (fds *Server) Close() {
	fds.listener.Close()
	syscall.Close(fds.dupedFD)
}

// Serve serves the FD
func (fds *Server) Serve() error {
	var deadline time.Time
	if fds.timeout != 0 {
		deadline = time.Now().Add(fds.timeout)
	}
	fds.listener.SetDeadline(deadline)
	for {
		conn, err := fds.listener.AcceptUnix()
		// Clean up after ourselves, since we are initiating our own
		// closure through the timeout.
		if os.IsTimeout(err) {
			fds.Close()
			return err
		} else if errors.Is(err, net.ErrClosed) {
			return nil
		} else if err != nil {
			return err
		}
		conn.SetDeadline(deadline)
		succeeded, err := fds.handleConnection(conn)
		if err != nil {
			return err
		}
		if succeeded && fds.serveOnce {
			break
		}
	}
	return nil
}

// GetSharedFD gets an FD served at udsPath with nonce
func GetSharedFD(udsPath, nonce string) (int, error) {
	// If you don't send at least a byte, the server won't recvmsg.  This
	// is a Linux UDS SOCK_STREAM thing.
	if len(nonce) == 0 {
		return 0, ErrEmptyNonce
	}

	ua, err := net.ResolveUnixAddr("unix", udsPath)
	if err != nil {
		return 0, err
	}
	uc, err := net.DialUnix("unix", nil, ua)
	if err != nil {
		return 0, err
	}

	n, err := uc.Write([]byte(nonce))
	if err != nil {
		return 0, err
	}
	if n != len(nonce) {
		return 0, ErrTruncatedWrite
	}

	oob := make([]byte, 1024)
	_, oobn, _, _, err := uc.ReadMsgUnix(nil, oob)
	if err != nil {
		return 0, err
	}
	scm, err := syscall.ParseSocketControlMessage(oob[:oobn])
	if err != nil {
		return 0, err
	}
	if len(scm) != 1 {
		return 0, ErrMissingSCM
	}
	urs, err := syscall.ParseUnixRights(&scm[0])
	if err != nil {
		return 0, err
	}
	if len(urs) != 1 {
		return 0, ErrNotOneUnixRights
	}
	return urs[0], nil
}
