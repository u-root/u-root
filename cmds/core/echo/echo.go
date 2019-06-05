// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Echo writes its arguments separated by blanks and terminated by a newline on
// the standard output.
//
// Synopsis:
//     echo [-e] [-n] [-E] [STRING]...
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type flags struct {
	noNewline, interpretEscapes bool
}

func escapeString(s string) (string, error) {
	if len(s) < 1 {
		return "", nil
	}

	s = strings.Split(s, "\\c")[0]
	s = strings.Replace(s, "\\0", "\\", -1)

	// Quote the string and scan it through %q to interpret backslash escapes
	s = fmt.Sprintf("\"%s\"", s)
	_, err := fmt.Sscanf(s, "%q", &s)
	if err != nil {
		return "", err
	}

	return s, nil
}

func echo(f flags, w io.Writer, s ...string) error {
	var err error
	line := strings.Join(s, " ")
	if f.interpretEscapes {
		line, err = escapeString(line)
		if err != nil {
			return err
		}

	}

	format := "%s"
	if !f.noNewline {
		format += "\n"
	}
	_, err = fmt.Fprintf(w, format, line)

	return err
}

func init() {
	defUsage := flag.Usage
	flag.Usage = func() {
		defUsage()
		fmt.Println(`
  If -e is in effect, the following sequences are recognized:
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
    \0NNN  byte with octal value NNN (1 to 3 digits)
    \xHH   byte with hexadecimal value HH (1 to 2 digits)`)
	}
}

func main() {
	var (
		f flags
		E bool
	)
	flag.BoolVar(&f.noNewline, "n", false, "suppress newline")
	flag.BoolVar(&f.interpretEscapes, "e", true, "enable interpretation of backslash escapes (default)")
	flag.BoolVar(&E, "E", false, "disable interpretation of backslash escapes")
	flag.Parse()
	if E {
		f.interpretEscapes = false
	}

	err := echo(f, os.Stdout, flag.Args()...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}
}
