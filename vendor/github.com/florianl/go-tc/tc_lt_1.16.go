//go:build !go1.16
// +build !go1.16

package tc

import "log"

// devNull satisfies io.Writer, in case *log.Logger is not provided
type devNull struct{}

func (devNull) Write(p []byte) (int, error) {
	return 0, nil
}

func setDummyLogger() *log.Logger {
	return log.New(new(devNull), "", 0)
}
