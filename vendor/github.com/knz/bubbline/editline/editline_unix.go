//go:build !windows
// +build !windows

package editline

import "syscall"

var canSuspendProcess = true

func suspendProcess() {
	// Send SIGTSTP to the entire process group.
	_ = syscall.Kill(0, syscall.SIGTSTP)
}
