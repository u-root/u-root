// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

// nohup – invoke a utility immune to hangups.
//
// Synopsis:
//
//	nohup <command> [args...]
package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"golang.org/x/term"
)

var (
	errUsage  = fmt.Errorf("nohup <command> [args...]")
	errStart  = fmt.Errorf("failed to start")
	errFinish = fmt.Errorf("finished with error")
)

func main() {
	if err := run(os.Args); err != nil {
		if errors.Is(err, errUsage) {
			fmt.Fprintf(os.Stderr, "Usage: %v\n", errUsage)
			os.Exit(127)
		}
		log.Fatalf("nohup: %v", err)
	}
}

func run(args []string) error {
	if len(args) < 2 {
		return errUsage
	}

	signal.Ignore(syscall.SIGHUP)

	cmdName := args[1]
	cmdArgs := args[2:]

	cmd := exec.Command(cmdName, cmdArgs...)

	stdoutIsTerminal := term.IsTerminal(int(os.Stdout.Fd()))
	stderrIsTerminal := term.IsTerminal(int(os.Stderr.Fd()))

	if stdoutIsTerminal {
		outputFile, err := os.OpenFile("nohup.out", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		if err != nil {
			return fmt.Errorf("error opening file: %w", err)
		}
		defer outputFile.Close()

		cmd.Stdout = outputFile

		if stderrIsTerminal {
			cmd.Stderr = outputFile
		} else {
			cmd.Stderr = os.Stderr
		}
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("%s: %w: %w", cmdName, errStart, err)
	}

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("%s: %w: %w", cmdName, errFinish, err)
	}

	return nil
}
