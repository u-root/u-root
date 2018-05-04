package ldd

import (
	"os"
	"os/exec"
)

const ldso = "/libexec/ld-elf*.so.*"

// lddOutput runs the interpreter and returns its output.
func lddOutput(interp, file string) ([]byte, error) {
	cmd := exec.Command(interp, file)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "LD_TRACE_LOADED_OBJECTS=1")
	return cmd.Output()
}
