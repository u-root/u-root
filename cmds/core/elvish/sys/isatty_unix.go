// +build !windows,!plan9

package sys

import (
	"os"
	"github.com/u-root/u-root/pkg/termios"
)

// IsATTY returns true if the given file is a terminal.
func IsATTY(file *os.File) bool {
	_, err := termios.GetTermios(file.Fd())
	return err == nil
}
