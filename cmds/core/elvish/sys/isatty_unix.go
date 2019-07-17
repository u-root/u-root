// +build !windows,!plan9

package sys

import (
	"github.com/u-root/u-root/pkg/termios"
	"os"
)

// IsATTY returns true if the given file is a terminal.
func IsATTY(file *os.File) bool {
	_, err := termios.GetTermios(file.Fd())
	return err == nil
}
