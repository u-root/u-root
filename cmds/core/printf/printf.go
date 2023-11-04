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
	"os"

	"github.com/u-root/u-root/pkg/printf"
)

func run() {
	_, err := printf.Fprintf(os.Stdout, os.Args[1:]...)
	if err != nil {
		os.Stderr.Write([]byte(err.Error() + "\n"))
		return
	}
	return
}

func main() {
	run()
}
