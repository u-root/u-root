//+build linux

package netlink

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

// A netNS is a handle that can manipulate network namespaces.
//
// Operations performed on a netNS must use runtime.LockOSThread before
// manipulating any network namespaces.
type netNS struct {
	f *os.File
}

// threadNetNS constructs a netNS using the network namespace of the calling
// thread. If the namespace is not the default namespace, runtime.LockOSThread
// should be invoked first.
func threadNetNS() (*netNS, error) {
	f, err := os.Open(fmt.Sprintf("/proc/self/task/%d/ns/net", unix.Gettid()))
	if err != nil {
		return nil, err
	}

	return &netNS{f: f}, nil
}

// Close releases the handle to a network namespace.
func (n *netNS) Close() error { return n.f.Close() }

// FD returns a file descriptor which represents the network namespace.
func (n *netNS) FD() int { return int(n.f.Fd()) }

// Restore restores the original network namespace for the calling thread.
func (n *netNS) Restore() error { return n.Set(n.FD()) }

// Set sets a new network namespace for the current thread using fd.
func (n *netNS) Set(fd int) error {
	return os.NewSyscallError("setns", unix.Setns(fd, unix.CLONE_NEWNET))
}
