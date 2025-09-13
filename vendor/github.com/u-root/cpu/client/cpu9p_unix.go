// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !windows && !plan9
// +build !windows,!plan9

package client

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/hugelgupf/p9/p9"
	"golang.org/x/sys/unix"
)

// SetAttr implements p9.File.SetAttr.
func (l *CPU9P) SetAttr(mask p9.SetAttrMask, attr p9.SetAttr) error {
	var err error
	// Any or ALL can be set.
	// A setattr could include things to set,
	// and a permission value that makes setting those
	// things impossible. Therefore, do these
	// permission-y things last:
	// Permissions
	// GID
	// UID
	// Since changing, e.g., UID or GID might make
	// changing permissions impossible.
	//
	// The test actually caught this ...

	if mask.Size {
		if e := unix.Truncate(l.path, int64(attr.Size)); e != nil {
			err = errors.Join(err, fmt.Errorf("truncate:%w", err))
		}
	}
	if mask.ATime || mask.MTime {
		atime, mtime := time.Now(), time.Now()
		if mask.ATimeNotSystemTime {
			atime = time.Unix(int64(attr.ATimeSeconds), int64(attr.ATimeNanoSeconds))
		}
		if mask.MTimeNotSystemTime {
			mtime = time.Unix(int64(attr.MTimeSeconds), int64(attr.MTimeNanoSeconds))
		}
		if e := os.Chtimes(l.path, atime, mtime); e != nil {
			err = errors.Join(err, e)
		}
	}

	if mask.CTime {
		// The Linux client sets CTime. I did not even know that was allowed.
		// if e := errors.New("Can not set CTime on Unix"); e != nil { err = errors.Join(e)}
		verbose("mask.CTime is set by client; ignoring")
	}
	if mask.Permissions {
		perm := uint32(attr.Permissions)
		if e := unix.Chmod(l.path, perm); e != nil {
			err = errors.Join(err, fmt.Errorf("%q:%o:%w", l.path, perm, err))
		}
	}

	if mask.GID {
		if e := unix.Chown(l.path, -1, int(attr.GID)); e != nil {
			err = errors.Join(err, e)
		}
	}
	if mask.UID {
		if e := unix.Chown(l.path, int(attr.UID), -1); e != nil {
			err = errors.Join(err, e)
		}
	}
	return err
}

// Lock implements p9.File.Lock.
func (l *CPU9P) Lock(pid int, locktype p9.LockType, flags p9.LockFlags, start, length uint64, client string) (p9.LockStatus, error) {
	var cmd int
	switch flags {
	case p9.LockFlagsBlock:
		cmd = unix.F_SETLKW
	case p9.LockFlagsReclaim:
		return p9.LockStatusError, unix.ENOSYS
	default:
		cmd = unix.F_SETLK
	}
	var t int16
	switch locktype {
	case p9.ReadLock:
		t = unix.F_RDLCK
	case p9.WriteLock:
		t = unix.F_WRLCK
	case p9.Unlock:
		t = unix.F_UNLCK
	default:
		return p9.LockStatusError, unix.ENOSYS
	}
	lk := &unix.Flock_t{
		Type:   t,
		Whence: unix.SEEK_SET,
		Start:  int64(start),
		Len:    int64(length),
	}
	if err := unix.FcntlFlock(l.file.Fd(), cmd, lk); err != nil {
		if errors.Is(err, unix.EAGAIN) {
			return p9.LockStatusBlocked, nil
		}
		return p9.LockStatusError, err
	}
	return p9.LockStatusOK, nil
}
