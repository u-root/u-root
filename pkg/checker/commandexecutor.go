package checker

import "os/exec"

// CommandExecutor returns a check that runs the provided command and arguments.
func CommandExecutor(prog string, args ...string) Checker {
	return func() error {
		cmd := exec.Command(prog, args...)
		if err := cmd.Start(); err != nil {
			return err
		}
		if err := cmd.Wait(); err != nil {
			return err
		}
		return nil
	}
}
