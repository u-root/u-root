package eval

import (
	"errors"
	"os"
	"os/exec"
)

// Command and process control.

var ErrNotInSameGroup = errors.New("not in the same process group")

func init() {
	addBuiltinFns(map[string]interface{}{
		// Command resolution
		"external":        external,
		"has-external":    hasExternal,
		"search-external": searchExternal,

		// Process control
		"fg":   fg,
		"exec": execFn,
		"exit": exit,
	})
}

func external(cmd string) ExternalCmd {
	return ExternalCmd{cmd}
}

func hasExternal(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func searchExternal(cmd string) (string, error) {
	return exec.LookPath(cmd)
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
