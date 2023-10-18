// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// NAME
//
//	printf - format and print data
//
// SYNOPSIS
//
//	printf FORMAT [ARGUMENT]...
//
// DESCRIPTION
//
//				 Print ARGUMENT(s) according to FORMAT:
//
//	       FORMAT controls the output as in C printf.  Interpreted sequences are:
//
//	       \"     double quote
//
//	       \\     backslash
//
//	       \a     alert (BEL)
//
//	       \b     backspace
//
//	       \c     produce no further output
//
//	       \e     escape
//
//	       \f     form feed
//
//	       \n     new line
//
//	       \r     carriage return
//
//	       \t     horizontal tab
//
//	       \v     vertical tab
//
//	       \NNN   byte with octal value NNN (1 to 3 digits)
//
//	       \xHH   byte with hexadecimal value HH (1 to 2 digits)
//
//	       \uHHHH Unicode (ISO/IEC 10646) character with hex value HHHH (4 digits)
//
//	       \UHHHHHHHH
//	              Unicode character with hex value HHHHHHHH (8 digits)
//
//	       %%     a single %
//
//	       %b     ARGUMENT as a string with '\' escapes interpreted, except that octal escapes are of the form \0 or \0NNN
//
//	       %q     ARGUMENT  is  printed  in  a format that can be reused as shell input, escaping non-printable characters with the proposed POSIX $'' syntax.
//
//	       %{diouxXfeEgGcs} the C format specifications, with ARGUMENTs converted to proper type first.  Variable widths are handled

package main

import (
	"bytes"
	"io"
	"os"
)

type printf struct {
	Args   []string
	Stdout io.Writer
	Stderr io.Writer

	format    string
	arguments []string
}

func NewPrinter(stdout, stderr io.Writer, args []string) *printf {
	return &printf{
		Stdout: stdout,
		Stderr: stderr,
		Args:   args,
	}
}

func (c *printf) exec(w *bytes.Buffer) (err error) {
	return interpret(w, c.format, c.arguments, false, true)
}

func (c *printf) run() {
	if len(c.Args) < 1 {
		c.Stderr.Write([]byte("printf: not enough arguments\n"))
		return
	}
	c.format = c.Args[0]
	if len(c.Args) > 1 {
		c.arguments = c.Args[1:]
	}
	w := new(bytes.Buffer)
	err := c.exec(w)
	if err != nil {
		c.Stderr.Write([]byte("printf: " + err.Error() + "\n"))
		return
	}
	// flush on success
	w.WriteTo(c.Stdout)
}

func run() {
	cmd := NewPrinter(os.Stdout, os.Stderr, os.Args[1:])
	cmd.run()
	return
}

func main() {
	run()
}
