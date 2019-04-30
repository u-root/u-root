package checker

import (
	"os"
	"os/exec"
)

// CommandExecutor returns a check that runs the provided command and arguments.
func CommandExecutor(prog string, args ...string) Checker {
	return func() error {
		cmd := exec.Command(prog, args...)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
		return nil
	}
}
