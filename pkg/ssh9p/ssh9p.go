// SPDX-License-Identifier: MIT
// Copyright 2026 Google LLC
//
// Package ssh9p provides functions to serve and mount 9P servers over SSH.

//go:build linux

package ssh9p

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/hugelgupf/p9/p9"
	"github.com/u-root/u-root/pkg/sshstream"
	"golang.org/x/crypto/ssh"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sys/unix"
)

// Server wraps a p9 Server with optional SSH support.
type Server struct {
	*p9.Server
	cfg *ssh.ServerConfig
}

// ServerOpt is an option function.
type ServerOpt func(s *Server)

// WithSSHServer uses SSH for the 9P server's transport.
func WithSSHServer(cfg *ssh.ServerConfig) ServerOpt {
	return func(s *Server) {
		s.cfg = cfg
	}
}

// NewServer wraps a p9.Server.
func NewServer(p9s *p9.Server, opts ...ServerOpt) *Server {
	s := &Server{Server: p9s}

	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Serve handles requests from the listener.  Will close the listener when the
// context is done.
func (s *Server) Serve(ctx context.Context, l net.Listener) error {
	var wg sync.WaitGroup
	defer wg.Wait()

	stopf := context.AfterFunc(ctx, func() {
		l.Close()
	})
	defer stopf()

	for {
		c, err := l.Accept()
		if err != nil {
			return fmt.Errorf("accept: %w", err)
		}

		wg.Go(func() {
			// The func is a closure, so c will be whatever the
			// current version of c is, i.e. it could be nc from
			// below.
			stopf := context.AfterFunc(ctx, func() {
				c.Close()
			})
			defer func() {
				// Abort the AfterFunc (which might never run)
				// and close it ourselves.
				if stopf() {
					c.Close()
				}
			}()
			if s.cfg != nil {
				nc, err := sshstream.NewServer(c, s.cfg)
				if err != nil {
					log.Printf("ssh9p NewServer: %v", err)
					return

				}
				c = nc
			}
			if err := s.Handle(c, c); err != nil {
				log.Printf("ssh9p Handle: %v", err)
			}
		})
	}
}

type mountOpts struct {
	unixFlags      int
	msize          int
	fastTCPTimeout bool
	cache          string
	cfg            *ssh.ClientConfig
}

// MountOpt is an option function.
type MountOpt func(*mountOpts)

// WithUnixFlags sets mount flags, e.g. unix.MS_NOSUID.
func WithUnixFlags(flags int) MountOpt {
	return func(m *mountOpts) {
		m.unixFlags = flags
	}
}

// WithMsize sets the 9P msize (max message size in bytes) parameter.
func WithMsize(msize int) MountOpt {
	return func(m *mountOpts) {
		m.msize = msize
	}
}

// WithFastTCPTimeout tells the kernel to quickly kill the underlying TCP
// connection, such that Mount9P() quickly tears down if the remote machine is
// unresponsive.
func WithFastTCPTimeout(enable bool) MountOpt {
	return func(m *mountOpts) {
		m.fastTCPTimeout = enable
	}
}

// WithSSHClient uses ssh for the net.Conn.  In this case, the caller must keep
// the io.Closer open to maintain the mount.
func WithSSHClient(cfg *ssh.ClientConfig) MountOpt {
	return func(m *mountOpts) {
		m.cfg = cfg
	}
}

// WithCache sets the cache option.  See
// https://docs.kernel.org/filesystems/9p.html for valid values.
func WithCache(cache string) MountOpt {
	return func(m *mountOpts) {
		m.cache = cache
	}
}

// fdFor9PTransportMount will take a conn and return an FD that can be passed to
// the kernel for a 9P transport mount.  Cancel the context when the mount is
// done.
//
// The kernel takes FDs that it can do raw reads and writes on, expecting those
// FDs to speak 9P.  If we're using a TCP socket, we can just hand the socket FD
// to the kernel.
//
// But if we're using ssh (or more broadly, not a net.Conn on a raw socket with
// a File() method), we need to translate ssh to 9P in userspace and hand the
// plaintext end of the pipe to the kernel.  This is why Mount9P needs something
// kept alive in userspace.
func fdFor9PTransportMount(ctx context.Context, eg *errgroup.Group, c net.Conn) (*os.File, error) {
	if tcpConn, ok := c.(*net.TCPConn); ok {
		tcpFile, err := tcpConn.File()
		if err == nil {
			return tcpFile, nil
		}
	}
	fds, err := unix.Socketpair(unix.AF_UNIX, syscall.SOCK_STREAM, 0)
	if err != nil {
		return nil, fmt.Errorf("socketpair: %w", err)
	}
	userside := os.NewFile(uintptr(fds[0]), "userside")
	kernelside := os.NewFile(uintptr(fds[1]), "kernelside")

	copyError := func(who string, w io.Writer, r io.Reader) error {
		_, err := io.Copy(w, r)
		// eg.Wait() returns the first err returned from its
		// goroutines.  If the ctx was cancelled from higher, we want
		// that error, not the error we get from being torn down.
		if err := ctx.Err(); err != nil {
			return err
		}
		if err != nil {
			return err
		}
		// Return some error to trigger the end of the eg.
		return errors.New("closed 9p " + who)
	}
	eg.Go(func() error { return copyError("write to conn", c, userside) })
	eg.Go(func() error { return copyError("write to userside", userside, c) })

	eg.Go(func() error {
		<-ctx.Done()
		// The ctx can be canelled from higher up for from one of the
		// copy goroutines erroring out.  Either way, our job is to make
		// sure the goroutines (if any remain) exit, and we do that by
		// making their io.Copy exit.  And the main thing to do there is
		// to make their Reads return.  (Also want the writes to return,
		// but the reads is where they are often blocked).
		//
		// Closing c (a net.Conn) will unblock any pending reads and
		// writes.  However, closing userside will *not* unblock any
		// reads.  It will unblock reads on the *other* side of the
		// socketpair (aka, kernelside), but not any reads of userside.
		// To handle that, we need to shutdown(2).
		if err := syscall.Shutdown(int(userside.Fd()), syscall.SHUT_RDWR); err != nil {

			log.Println("Failed to shutdown userside:", err)
		}
		userside.Close()
		c.Close()

		return ctx.Err()
	})

	return kernelside, nil
}

// Mount9P mounts a 9p server connection at mountPoint.
//
// If necessary (e.g. using ssh), this will block to maintain the mount.  If so:
// - You can tear it down by cancelling the context.
// - It will always return a non-nil error.
// - It will unmount the mountpoint when its done.  (Since we know the mount is
// inoperable).
func Mount9P(ctx context.Context, mntReady chan<- bool, c net.Conn, mountPoint string, opts ...MountOpt) error {
	// We're responsible for closing the mntReady chan, even on error.
	defer close(mntReady)

	m := &mountOpts{
		msize: 64 * 1024,
		cache: "none",
	}
	for _, opt := range opts {
		opt(m)
	}

	if m.fastTCPTimeout {
		// Could consider setting something for retrans too.
		if tcpConn, ok := c.(*net.TCPConn); ok {
			if err := tcpConn.SetKeepAliveConfig(net.KeepAliveConfig{
				Enable:   true,
				Idle:     5 * time.Second,
				Interval: 1 * time.Second,
				Count:    5,
			}); err != nil {
				log.Println("Failed to set TCP KeepAlive, continuing:", err)
			}
		}
	}

	if m.cfg != nil {
		nc, err := sshstream.NewClient(c, m.cfg)
		if err != nil {
			return err
		}
		c = nc
	}

	// If any goroutine of the eg errors, that ctx will be cancelled.  But
	// we could also error out on our own, for example the unix.Mount()
	// command.  That's why we set up the defer cancel() first.  (Note the
	// cancellable context is a parent to the eg's ctx, and cancellations
	// propagate downward).
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)

	tfd, err := fdFor9PTransportMount(ctx, eg, c)
	if err != nil {
		return err
	}
	// We can close the FD after we give it to the kernel.  (internal dup).
	defer tfd.Close()

	mntstr := fmt.Sprintf("version=9p2000.L,noxattr,cache=%s,trans=fd,rfdno=%d,wfdno=%d,debug=0,msize=%d", m.cache, tfd.Fd(), tfd.Fd(), m.msize)
	if err := unix.Mount("ssh9p", mountPoint, "9p", uintptr(m.unixFlags), mntstr); err != nil {
		return fmt.Errorf("9P mount: %w", err)
	}

	mntReady <- true

	// If we didn't spawn any goroutines in the eg, perhaps due to being a
	// raw TCP connection, eg.Wait() will just return immediately.  If there
	// are any post-mount errors, e.g. the remote server shut down, we'll
	// exit with an error.  Similarly if our original context is cancelled,
	// we'll return with whatever error came with that cancellation.
	err = eg.Wait()
	if err != nil {
		// Since we had goroutines maintaining the mount (e.g. for SSH),
		// and they are gone, we might as well tear down the mount too.
		if err := unix.Unmount(mountPoint, unix.MNT_FORCE|unix.MNT_DETACH); err != nil {
			log.Println("Could not unmount, continuing:", err)
		}
	}
	return err
}
