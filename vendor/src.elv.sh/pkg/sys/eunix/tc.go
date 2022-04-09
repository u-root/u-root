//go:build !windows && !plan9
// +build !windows,!plan9

package eunix

import (
	"golang.org/x/sys/unix"
)

// Tcsetpgrp sets the terminal foreground process group.
func Tcsetpgrp(fd int, pid int) error {
	return unix.IoctlSetPointerInt(fd, unix.TIOCSPGRP, pid)
}
