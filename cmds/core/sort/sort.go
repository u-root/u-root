// Copyright 2017-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Sort lines.
//
// Synopsis:
//
//	sort [OPTIONS]... [INPUT]...
//
// Description:
//
//	Sort copies lines from the input to the output, sorting them in the
//	process. This does nothing fancy (no multi-threading, compression,
//	optiminzations, ...); it simply uses Go's sort.Sort function.
//
// Options:
//
//	-r:      Reverse the result of comparisons
//	-C:      Check that the single input file is ordered. No warnings.
//	-o FILE: Specify the name of an output file to be used instead of the standard output.
package main

import (
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"sort"
	"strings"
)

var (
	reverse    = flag.Bool("r", false, "Reverse the result of comparisons.")
	ordered    = flag.Bool("C", false, "Check that the single input file is ordered. No warnings.")
	outputFile = flag.String("o", "", "Specify the name of an output file to be used instead of the standard output.")
)

var errNotOrdered = errors.New("not ordered")

type params struct {
	reverse    bool
	ordered    bool
	outputFile string
}

type cmd struct {
	stdin  io.ReadCloser
	stdout io.Writer
	stderr io.Writer
	params params
	args   []string
}

func command(stdin io.ReadCloser, stdout, stderr io.Writer, p params, args []string) *cmd {
	return &cmd{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
		params: p,
		args:   args,
	}
}

func (c *cmd) run() error {
	// Input files
	from := []io.ReadCloser{}
	for _, v := range c.args {
		f, err := os.Open(v)
		if err != nil {
			return err
		}
		defer f.Close()
		from = append(from, f)
	}

	if len(c.args) == 0 {
		from = append(from, c.stdin)
	}

	// Read unicode string from input
	fileContents := []string{}
	for _, f := range from {
		bytes, err := io.ReadAll(f)
		if err != nil {
			return err
		}
		s := string(bytes)
		fileContents = append(fileContents, s)
		// Add a separator between files if the file is not newline
		// terminated. Prevents concatenating lines between files.
		if len(s) > 0 && s[len(s)-1] != '\n' {
			fileContents = append(fileContents, "\n")
		}
	}

	s := strings.Join(fileContents, "")
	if len(s) > 0 && s[len(s)-1] == '\n' {
		s = s[:len(s)-1] // remove newline terminator
	}

	if c.params.ordered {
		lines := strings.Split(s, "\n")
		if sort.IsSorted(sort.StringSlice(lines)) {
			return nil
		}
		return errNotOrdered
	}

	if err := c.writeOutput(c.stdout, c.sortAlgorithm(s)); err != nil {
		return err
	}
	return nil
}

func (c *cmd) sortAlgorithm(s string) string {
	if len(s) == 0 {
		return "" // edge case mimics coreutils
	}
	lines := strings.Split(s, "\n")
	if c.params.reverse {
		sort.Sort(sort.Reverse(sort.StringSlice(lines)))
	} else {
		sort.Strings(lines)
	}
	return strings.Join(lines, "\n") + "\n" // append newline terminator
}

func (c *cmd) writeOutput(w io.Writer, s string) error {
	to := w
	if c.params.outputFile != "" {
		f, err := os.Create(c.params.outputFile)
		if err != nil {
			return err
		}
		defer f.Close()
		to = f
	}

	_, err := to.Write([]byte(s))
	return err
}

func main() {
	flag.Parse()
	p := params{*reverse, *ordered, *outputFile}
	if err := command(os.Stdin, os.Stdout, os.Stderr, p, flag.Args()).run(); err != nil {
		if err == errNotOrdered {
			os.Exit(1)
		}
		log.Fatal(err)
	}
}
