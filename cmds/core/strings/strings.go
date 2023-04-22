// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Strings finds printable strings.
//
// Synopsis:
//
//	strings OPTIONS [FILES]...
//
// Description:
//
//	Prints all sequences of `n` or more printable characters terminated by a
//	non-printable character (or EOF).
//
//	If no files are specified, read from stdin.
//
// Options:
//
//	-n number: the minimum string length (default is 4)
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	flag "github.com/spf13/pflag"
)

var n = flag.Int("n", 4, "the minimum string length")

type cmd struct {
	stdin  io.Reader
	stdout io.Writer
	params
	args []string
}

type params struct {
	n int
}

func command(stdin io.Reader, stdout io.Writer, p params, args []string) *cmd {
	return &cmd{
		stdin:  stdin,
		stdout: stdout,
		params: p,
		args:   args,
	}
}

func (c *cmd) run() error {
	if c.n < 1 {
		return fmt.Errorf("strings: invalid minimum string length %v", c.n)
	}
	if len(c.args) == 0 {
		rb := bufio.NewReader(c.stdin)
		if err := c.stringsIO(rb, c.stdout); err != nil {
			return err
		}
	}
	for _, file := range c.args {
		if err := c.stringsFile(file, c.stdout); err != nil {
			return err
		}
	}
	return nil
}

func asciiIsPrint(char byte) bool {
	return char >= 32 && char <= 126
}

func (c *cmd) stringsIO(r *bufio.Reader, w io.Writer) error {
	var o []byte
	for {
		b, err := r.ReadByte()
		if err == io.EOF {
			if len(o) >= c.n {
				w.Write(o)
				w.Write([]byte{'\n'})
			}
			return nil
		}
		if err != nil {
			return err
		}
		if !asciiIsPrint(b) {
			if len(o) >= c.n {
				w.Write(o)
				w.Write([]byte{'\n'})
			}
			o = o[:0]
			continue
		}
		// Prevent the buffer from growing indefinitely.
		if len(o) >= c.n+1024 {
			w.Write(o[:1024])
			o = o[1024:]
		}
		o = append(o, b)
	}
}

func (c *cmd) stringsFile(file string, w io.Writer) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	// Buffer reduces number of syscalls.
	rb := bufio.NewReader(f)
	return c.stringsIO(rb, w)
}

func main() {
	flag.Parse()
	if err := command(os.Stdin, os.Stdout, params{n: *n}, flag.Args()).run(); err != nil {
		log.Fatalf("strings: %v", err)
	}
}
