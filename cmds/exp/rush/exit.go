package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

var errExitUsage = errors.New("usage: exit <|0-255|>")

// Add builtins to the shell
func init() {
	_ = addBuiltIn("exit", exit)
}

// exit command: exit the shell with a given exit code
func exit(c *Command) error {
	code := 0 // Default exit code is 0

	if len(c.argv) > 1 {
		return errExitUsage
	}
	if len(c.argv) == 1 {
		// Attempt to convert argument to integer
		var err error
		code, err = strconv.Atoi(c.argv[0])
		if err != nil || code < 0 || code > 255 {
			return fmt.Errorf("exit: invalid exit code")
		}
	}

	// Exit with the specified code
	os.Exit(code)
	return nil
}
