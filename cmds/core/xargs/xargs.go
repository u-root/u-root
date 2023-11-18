// Copyright 2013-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

const defaultMaxArgs = 5000

func main() {
	var maxNumber = flag.Int("n", defaultMaxArgs, "max number of arguments per command")
	var trace = flag.Bool("t", false, "enable trace mode, each command is written to stderr")
	flag.Parse()
	if err := run(os.Stdin, os.Stdout, os.Stderr, *maxNumber, *trace, flag.Args()...); err != nil {
		log.Fatal(err)
	}
}

func run(stdin io.Reader, stdout, stderr io.Writer, maxArgs int, trace bool, args ...string) error {
	if len(args) == 0 {
		args = append(args, "echo")
	}

	var xArgs []string
	scanner := bufio.NewScanner(stdin)
	for scanner.Scan() {
		sp := strings.Fields(scanner.Text())
		xArgs = append(xArgs, sp...)
	}

	argsLen := len(args)

	for i := 0; i < len(xArgs); i += maxArgs {
		m := len(xArgs)
		if i+maxArgs < m {
			m = i + maxArgs
		}
		args = append(args, xArgs[i:m]...)

		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdin = stdin
		cmd.Stdout = stdout
		cmd.Stderr = stderr

		if trace {
			fmt.Fprintf(stderr, "%s\n", strings.Join(args, " "))
		}

		if err := cmd.Run(); err != nil {
			return err
		}

		args = args[:argsLen]
	}

	return nil
}
