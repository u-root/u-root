package eval

import (
	"errors"
	"os"
)

// Command and process control.

var ErrNotInSameGroup = errors.New("not in the same process group")

func init() {
	addBuiltinFns(map[string]interface{}{
		// Process control
		"fg":   fg,
		"exec": execFn,
		"exit": exit,
	})
}

func exit(fm *Frame, codes ...int) error {
	code := 0
	switch len(codes) {
	case 0:
	case 1:
		code = codes[0]
	default:
		return ErrArgs
	}

	preExit(fm)
	os.Exit(code)
	// Does not return
	panic("os.Exit returned")
}

func preExit(fm *Frame) {
}

var errNotSupportedOnWindows = errors.New("not supported on Windows")

func notSupportedOnWindows() error {
	return errNotSupportedOnWindows
}
