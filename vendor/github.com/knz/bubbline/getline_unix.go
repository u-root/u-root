//go:build !windows
// +build !windows

package bubbline

import (
	"os"

	"golang.org/x/sys/unix"
)

var stopSignals = []os.Signal{unix.SIGINT, unix.SIGTERM}
