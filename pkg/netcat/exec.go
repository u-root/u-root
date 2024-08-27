// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netcat

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
)

type ExecType int

const (
	EXEC_TYPE_NATIVE ExecType = iota
	EXEC_TYPE_SHELL
	EXEC_TYPE_LUA
	EXEC_TYPE_NONE // For faster case switching this is appended at the end
)

type Exec struct {
	Type    ExecType
	Command string
}

func ParseCommands(execs ...Exec) (Exec, error) {
	cmds := 0
	last_valid := -1

	for i, e := range execs {
		if e.Command == "" {
			continue
		}
		last_valid = i

		cmds++

	}

	// This is a recoverable error, we can just ignore the command
	if last_valid == -1 {
		return Exec{Type: EXEC_TYPE_NONE}, nil
	}

	if cmds > 1 {
		return Exec{}, fmt.Errorf("only one of --exec, --sh-exec, and --lua-exec is allowed")
	}

	return Exec{
		Type:    execs[last_valid].Type,
		Command: execs[last_valid].Command,
	}, nil
}

// Execute a given command on the host system
// stdout of the command is send to to the connection
// stderr of the command is displayed on stdout of the host
// The host process exits with the exit code of the command unless --keep-open is specified
func (n *Exec) Execute(stdin io.ReadWriter, stdout io.Writer, stderr io.Writer, eol []byte) error {
	var (
		cmd    *exec.Cmd
		buffer bytes.Buffer
	)

	if n.Command == "" {
		return fmt.Errorf("empty command")
	}

	switch n.Type {
	case EXEC_TYPE_NATIVE:
		commandParts := strings.Fields(n.Command)
		cmd = exec.Command(commandParts[0], commandParts[1:]...)
	case EXEC_TYPE_SHELL:
		cmd = exec.Command(DEFAULT_SHELL, "-c", n.Command)
	case EXEC_TYPE_LUA:
		return fmt.Errorf("not implemented")
	default:
		return fmt.Errorf("invalid exec type")
	}

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	scanner := bufio.NewScanner(stdin)
	for scanner.Scan() {
		buffer.WriteString(scanner.Text())
		buffer.Write(eol)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	cmd.Stdin = &buffer

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("exec start: %w", err)
	}

	// Wait waits for the command to exit and waits for any copying to stdin or copying from stdout or stderr to complete.
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("exec wait: %w", err)
	}

	return nil
}
