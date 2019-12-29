// +build linux,go1.13

package netlink

import (
	"errors"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/sys/unix"
)

func TestOpErrorUnwrapLinux(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		target error
		ok     bool
	}{
		{
			name:   "ENOBUFS",
			err:    unix.ENOBUFS,
			target: os.ErrNotExist,
		},
		{
			name: "OpError ENOBUFS",
			err: &OpError{
				Op:  "receive",
				Err: unix.ENOBUFS,
			},
			target: os.ErrNotExist,
		},
		{
			name: "OpError os.SyscallError ENOBUFS",
			err: &OpError{
				Op:  "receive",
				Err: os.NewSyscallError("recvmsg", unix.ENOBUFS),
			},
			target: os.ErrNotExist,
		},
		{
			name:   "ENOENT",
			err:    unix.ENOENT,
			target: os.ErrNotExist,
			ok:     true,
		},
		{
			name: "OpError ENOENT",
			err: &OpError{
				Op:  "receive",
				Err: unix.ENOENT,
			},
			target: os.ErrNotExist,
			ok:     true,
		},
		{
			name: "OpError os.SyscallError ENOENT",
			err: &OpError{
				Op:  "receive",
				Err: os.NewSyscallError("recvmsg", unix.ENOENT),
			},
			target: os.ErrNotExist,
			ok:     true,
		},
		{
			name: "OpError os.SyscallError EEXIST",
			err: &OpError{
				Op:  "receive",
				Err: os.NewSyscallError("recvmsg", unix.EEXIST),
			},
			target: os.ErrExist,
			ok:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := errors.Is(tt.err, tt.target)
			if diff := cmp.Diff(tt.ok, got); diff != "" {
				t.Fatalf("unexpected result (-want +got):\n%s", diff)
			}
		})
	}
}
