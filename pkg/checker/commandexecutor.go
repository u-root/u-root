// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package checker

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// DefaultShell is used by EmergencyShell
var DefaultShell = "elvish"

func init() {
	registerCheckFun(CommandExecutor)
	registerCheckFun(EmergencyShell)
}

// CommandExecutor returns a check that runs the provided command and arguments.
func CommandExecutor(args CheckArgs) error {
	cmd, cmdArgs, err := commandExecutorParseArgs(args)
	if err != nil {
		return err
	}

	command := exec.Command(cmd, cmdArgs...)
	command.Stdin, command.Stdout, command.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := command.Run(); err != nil {
		return err
	}
	return nil
}

func commandExecutorParseArgs(args CheckArgs) (string, []string, error) {
	if args["cmd"] == "" {
		return "", nil, fmt.Errorf("argument 'cmd' is required")
	}

	cmdArgs := make([]string, 0)
	for i := 1; i <= 256; i++ {
		argName := fmt.Sprintf("arg%d", i)
		arg := args[argName]
		if arg == "" {
			break
		}

		cmdArgs = append(cmdArgs, arg)
	}

	return args["cmd"], cmdArgs, nil
}

// EmergencyShell is a remediation that prints the given banner, and then calls
// an emergency shell.
func EmergencyShell(args CheckArgs) error {
	log.Print(green("Running emergency shell: %s", DefaultShell))
	if args["banner"] != "" {
		log.Print(args["banner"])
	}
	return CommandExecutor(CheckArgs{
		"cmd": DefaultShell,
	})
}
