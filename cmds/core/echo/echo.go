// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// echo writes its arguments separated by blanks and terminated by a newline on
// the standard output.
//
// Synopsis:
//
//	echo [-e] [-n] [-E] [STRING]...
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/u-root/u-root/pkg/uroot/util"
)

var usage = `echo:
  If -e is in effect, the following sequences are recognized:
    \\     backslash
    \a     alert (BEL)
    \b     backspace
    \c     produce no further output
    \f     form feed
    \n     new line
    \r     carriage return
    \t     horizontal tab
    \v     vertical tab
    \0NNN  byte with octal value NNN (1 to 3 digits)
    \xHH   byte with hexadecimal value HH (1 to 2 digits)`

var (
	noNewline                 = flag.Bool("n", false, "suppress newline")
	interpretEscapes          = flag.Bool("e", true, "enable interpretation of backslash escapes (default)")
	interpretBackslashEscapes = flag.Bool("E", false, "disable interpretation of backslash escapes")
)

func escapeString(s string) (string, error) {
	if len(s) < 1 {
		return "", nil
	}

	s = strings.Split(s, "\\c")[0]
	s = strings.ReplaceAll(s, "\\0", "\\")

	// Quote the string and scan it through %q to interpret backslash escapes
	s = fmt.Sprintf("\"%s\"", s)
	_, err := fmt.Sscanf(s, "%q", &s)
	if err != nil {
		return "", err
	}
	return s, nil
}

func echo(w io.Writer, noNewline, escape, backslash bool, s ...string) error {
	var err error
	if backslash {
		escape = false
	}
	line := strings.Join(s, " ")
	if escape {
		line, err = escapeString(line)
		if err != nil {
			return err
		}
	}
	format := "%s"
	if !noNewline {
		format += "\n"
	}
	_, err = fmt.Fprintf(w, format, line)

	return err
}

func init() {
	flag.Usage = util.Usage(flag.Usage, usage)
}

func main() {
	flag.Parse()
	if err := echo(os.Stdout, *noNewline, *interpretEscapes, *interpretBackslashEscapes, flag.Args()...); err != nil {
		log.Fatalf("%v", err)
	}
}
