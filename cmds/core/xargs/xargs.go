// Copyright 2013-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if err := run(os.Stdin, os.Stdout, os.Stderr, os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func run(stdin io.Reader, stdout, stderr io.Writer, args []string) error {
	if len(args) == 0 {
		args = append(args, "echo")
	}

	scanner := bufio.NewScanner(stdin)
	for scanner.Scan() {
		sp := strings.Fields(scanner.Text())
		args = append(args, sp...)
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return cmd.Run()
}
