// Copyright 2018-2023 the u-root Authors. All rights reserved
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
//	-t string: write each string preceded by its byte offset from the start of the file (d decimal, o octal, x hexadecimal)
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

var (
	errInvalidFormatArgument = fmt.Errorf("invalid argument to option -t")
	errInvalidMinLength      = fmt.Errorf("invalid minimum string length -n")
)

type cmd struct {
	stdin  io.Reader
	stdout io.Writer
	params
	args   []string
	offset int
}

type params struct {
	t string
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
		return fmt.Errorf("%w: %d", errInvalidMinLength, c.n)
	}
	if c.t != "" && c.t != "d" && c.t != "o" && c.t != "x" {
		return fmt.Errorf("%w: %s", errInvalidFormatArgument, c.t)
	}
	if len(c.args) == 0 {
		rb := bufio.NewReader(c.stdin)
		if err := c.stringsIO(rb, c.stdout); err != nil {
			return err
		}
	}
	for _, file := range c.args {
		c.offset = 0
		if err := c.stringsFile(file, c.stdout); err != nil {
			return err
		}
	}
	return nil
}

func asciiIsPrint(char byte) bool {
	return char >= 32 && char <= 126
}

func (c *cmd) offsetValue(l int) string {
	offset := c.offset - l
	switch c.t {
	case "d":
		return fmt.Sprintf("%d ", offset)
	case "o":
		return fmt.Sprintf("%o ", offset)
	case "x":
		return fmt.Sprintf("%x ", offset)
	default:
		panic("t param parsed wrong")
	}
}

func (c *cmd) stringsIO(r *bufio.Reader, w io.Writer) error {
	var o []byte
	for {
		b, err := r.ReadByte()
		if err == io.EOF {
			if len(o) >= c.n {
				if c.t != "" {
					w.Write([]byte(c.offsetValue(len(o))))
				}
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
				if c.t != "" {
					w.Write([]byte(c.offsetValue(len(o))))
				}
				w.Write(o)
				w.Write([]byte{'\n'})
			}
			o = o[:0]
			c.offset++
			continue
		}
		// Prevent the buffer from growing indefinitely.
		if len(o) >= c.n+1024 {
			w.Write(o[:1024])
			o = o[1024:]
		}
		o = append(o, b)
		c.offset++
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
	var n int
	var t string
	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	f.IntVar(&n, "n", 4, "the minimum string length")
	f.StringVar(&t, "t", "", "write each string preceded by its byte offset from the start of the file (d decimal, o octal, x hexadecimal)")
	f.Parse(unixflag.OSArgsToGoArgs())
	if err := command(os.Stdin, os.Stdout, params{n: n, t: t}, f.Args()).run(); err != nil {
		log.Fatalf("strings: %v", err)
	}
}
