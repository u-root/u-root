package ldd

import (
	"os/exec"
)

const ldso = "/lib*/ld-linux-*.so.*"

// lddOutput runs the interpreter and returns its output.
func lddOutput(interp, file string) ([]byte, error) {
	return exec.Command(interp, "--list", file).Output()
}
