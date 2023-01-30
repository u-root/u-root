//go:build windows
// +build windows

package bubbline

import "os"

var stopSignals = []os.Signal{os.Interrupt}
