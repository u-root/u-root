// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
/*
NAME
       printf - format and print data

SYNOPSIS
       printf FORMAT [ARGUMENT]...

DESCRIPTION
			 Print ARGUMENT(s) according to FORMAT:

       FORMAT controls the output as in C printf.  Interpreted sequences are:

       \"     double quote

       \\     backslash

       \a     alert (BEL)

       \b     backspace

       \c     produce no further output

       \e     escape

       \f     form feed

       \n     new line

       \r     carriage return

       \t     horizontal tab

       \v     vertical tab

       \NNN   byte with octal value NNN (1 to 3 digits)

       \xHH   byte with hexadecimal value HH (1 to 2 digits)

       \uHHHH Unicode (ISO/IEC 10646) character with hex value HHHH (4 digits)

       \UHHHHHHHH
              Unicode character with hex value HHHHHHHH (8 digits)

       %%     a single %

       %b     ARGUMENT as a string with '\' escapes interpreted, except that octal escapes are of the form \0 or \0NNN

       %q     ARGUMENT  is  printed  in  a format that can be reused as shell input, escaping non-printable characters with the proposed
              POSIX $'' syntax.

       and all C format specifications ending with one of diouxXfeEgGcs, with ARGUMENTs converted to proper type first.  Variable widths
       are handled.
*/
package main

import (
	"bytes"
	"io"
	"os"
)

type command struct {
	format string
	args   []string

	stdout io.Writer
	stderr io.Writer
}

func (c *command) execFormat(w *bytes.Buffer) (err error) {
	// the current implementation performs this in two passes
	return interpret(w, c.format, c.args, false, true)
}

func (c *command) exec(w *bytes.Buffer) (err error) {
	return c.execFormat(w)
}

func (c *command) run() {
	w := new(bytes.Buffer)
	err := c.exec(w)
	if err != nil {
		c.stderr.Write([]byte("printf: " + err.Error() + "\n"))
		return
	}
	// flush on success
	w.WriteTo(c.stdout)
}

func run() {
	cmd := &command{}
	cmd.stdout = os.Stdout
	cmd.stderr = os.Stderr
	if len(os.Args) < 1 {
		cmd.stderr.Write([]byte("printf: not enough arguments\n"))
		return
	}
	cmd.format = os.Args[1]
	if len(os.Args) > 2 {
		cmd.args = os.Args[2:]
	}
	cmd.run()
	return
}

func main() {
	run()
}
