// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package lockfile coordinates process-based file locking.
//
// This package is designed to aid concurrency issues between different
// processes.
//
// Sample usage:
//
//	lf := lockfile.New("/var/apt/apt.lock")
//	// Blocks and waits if /var/apt/apt.lock already exists.
//	if err := lf.Lock(); err != nil {
//		log.Fatal(err)
//	}
//
//	defer lf.MustUnlock()
//
//	// Do something in /var/apt/??
//
// Two concurrent invocations of this program will compete to create
// /var/apt/apt.lock, and then make the other wait until /var/apt/apt.lock
// disappears.
//
// If the lock holding process disappears without removing the lock, another
// process using this library will detect that and remove the lock.
//
// If some other entity removes the lockfile erroneously, the lock holder's
// call to Unlock() will return an ErrRogueDeletion.
package lockfile

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
)

var (
	// ErrRogueDeletion means the lock file was removed by someone
	// other than the lock holder.
	ErrRogueDeletion = errors.New("cannot unlock lockfile owned by another process")

	// ErrBusy means the lock is being held by another living process.
	ErrBusy = errors.New("file is locked by another process")

	// ErrInvalidPID means the lock file contains an incompatible syntax.
	ErrInvalidPID = errors.New("lockfile points to file with invalid content")

	// ErrProcessDead means the lock is held by a process that does not exist.
	ErrProcessDead = errors.New("lockfile points to invalid PID")
)

var (
	errUnlocked = errors.New("file is unlocked")
)

// Lockfile is a process-based file lock.
type Lockfile struct {
	// path is the file whose existense is the lock.
	path string

	// pid is our PID.
	//
	// This mostly exists for testing.
	pid int
}

// New returns a new lock file at the given path.
func New(path string) *Lockfile {
	return &Lockfile{
		path: path,
		pid:  os.Getpid(),
	}
}

func (l *Lockfile) pidfile() (string, error) {
	dir, base := filepath.Split(l.path)
	pidfile, err := ioutil.TempFile(dir, fmt.Sprintf("%s-", base))
	if err != nil {
		return "", err
	}
	defer pidfile.Close()

	if _, err := io.WriteString(pidfile, fmt.Sprintf("%d", l.pid)); err != nil {
		if err := os.Remove(pidfile.Name()); err != nil {
			log.Fatalf("Lockfile could not remove %q: %v", pidfile.Name(), err)
		}
		return "", err
	}
	return pidfile.Name(), nil
}

// TryLock attempts to create the lock file, or returns ErrBusy if a valid lock
// file already exists.
//
// If the lock file is detected not to be valid, it is removed and replaced
// with our lock file.
func (l *Lockfile) TryLock() error {
	pidpath, err := l.pidfile()
	if err != nil {
		return err
	}

	if err := l.lockWith(pidpath); err != nil {
		if err := os.Remove(pidpath); err != nil {
			log.Fatalf("Lockfile could not remove %q: %v", pidpath, err)
		}
		return err
	}
	return nil
}

// Lock blocks until it can create a valid lock file.
//
// If a valid lock file already exists, it waits for the file to be deleted or
// until the associated process dies.
//
// If an invalid lock file exists, it will be deleted and we'll retry locking.
func (l *Lockfile) Lock() error {
	pidpath, err := l.pidfile()
	if err != nil {
		return err
	}

	// Spin, oh, spin.
	if err := l.lockWith(pidpath); err == ErrBusy {
		return l.lockWith(pidpath)
	} else if err != nil {
		if err := os.Remove(pidpath); err != nil {
			log.Fatalf("Lockfile could not remove %q: %v", pidpath, err)
		}
		return err
	}
	return nil
}

func (l *Lockfile) checkLockfile() error {
	owningPid, err := ioutil.ReadFile(l.path)
	if os.IsNotExist(err) {
		return errUnlocked
	} else if err != nil {
		return err
	}

	if len(owningPid) == 0 {
		return ErrInvalidPID
	}

	pid, err := strconv.Atoi(string(owningPid))
	if err != nil || pid <= 0 {
		return ErrInvalidPID
	}

	p, err := os.FindProcess(pid)
	if err != nil {
		return ErrProcessDead
	}

	if err := p.Signal(syscall.Signal(0)); err != nil {
		return ErrProcessDead
	}

	if pid == l.pid {
		return nil
	}
	return ErrBusy
}

// Unlock attempts to delete the lock file.
//
// If we are not the lock holder, and the lock holder is an existing process,
// ErrRogueDeletion will be returned.
//
// If we are not the lock holder, and there is no valid lock holder (process
// died, invalid lock file syntax), ErrRogueDeletion will be returned.
func (l *Lockfile) Unlock() error {
	switch err := l.checkLockfile(); err {
	case nil:
		// Nuke the symlink and its target.
		target, err := os.Readlink(l.path)
		if os.IsNotExist(err) {
			return ErrRogueDeletion
		} else if err != nil {
			// The symlink is somehow screwed up. Just nuke it.
			if err := os.Remove(l.path); err != nil {
				log.Fatalf("Lockfile could not remove %q: %v", l.path, err)
			}
			return err
		}

		absTarget := resolveSymlinkTarget(l.path, target)
		if err := os.Remove(absTarget); os.IsNotExist(err) {
			return ErrRogueDeletion
		} else if err != nil {
			return err
		}

		if err := os.Remove(l.path); os.IsNotExist(err) {
			return ErrRogueDeletion
		} else if err != nil {
			return err
		}
		return nil

	case ErrInvalidPID, ErrProcessDead, errUnlocked, ErrBusy:
		return ErrRogueDeletion

	default:
		return err
	}
}

// MustUnlock panics if the call to Unlock fails.
func (l *Lockfile) MustUnlock() {
	if err := l.Unlock(); err != nil {
		log.Fatalf("could not unlock %q: %v", l.path, err)
	}
}

// resolveSymlinkTarget returns an absolute path for a given symlink target.
//
// Symlinks targets' "working directory" is the symlink's parent.
//
// Said another way, a symlink is always resolved relative to the symlink's
// parent.
//
// E.g.
// /foo/bar -> ./zoo resolves to the absolute path /foo/zoo
func resolveSymlinkTarget(symlink, target string) string {
	if filepath.IsAbs(target) {
		return target
	}

	return filepath.Join(filepath.Dir(symlink), target)
}

func (l *Lockfile) lockWith(pidpath string) error {
	switch err := os.Symlink(pidpath, l.path); {
	case err == nil:
		return nil

	case !os.IsExist(err):
		// Some kind of system error.
		return err

	default:
		// Symlink already exists.
		switch err := l.checkLockfile(); err {
		case errUnlocked:
			return l.lockWith(pidpath)

		case ErrInvalidPID, ErrProcessDead:
			// Nuke the symlink and its target.
			target, err := os.Readlink(l.path)
			if os.IsNotExist(err) {
				return l.lockWith(pidpath)
			} else if err != nil {
				// File might not be a symlink at all?
				// Leave it alone, in case it's someone's
				// legitimate file.
				return err
			}

			absTarget := resolveSymlinkTarget(l.path, target)
			// If it doesn't exist anymore, whatever. The symlink's
			// existence is the actual lock.
			if err := os.Remove(absTarget); !os.IsNotExist(err) && err != nil {
				return err
			}

			if err := os.Remove(l.path); os.IsNotExist(err) {
				return l.lockWith(pidpath)
			} else if err != nil {
				return err
			}

			// Retry making the symlink.
			return l.lockWith(pidpath)

		case ErrBusy, nil:
			return err

		default:
			return err
		}
	}
}
