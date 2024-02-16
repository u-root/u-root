// Copyright 2017-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package shlex is a Unicode-supporting POSIX command-line argument parser.
//
// shlex will parse for example
//
//     start --append="foobar foobaz" --nogood 'food'
//
// into the appropriate argvs to start the command.
package shlex

import (
	"unicode"
)

type state uint8

const (
	unquoted state = iota
	escape
	singleQuote
	doubleQuote
	doubleQuoteEscape
	comment
)

// Split splits a command line according to Bash shell rules.
//
// Split is compatible with Bash quoting as described in
// https://www.gnu.org/software/bash/manual/html_node/Quoting.html
// and GRUB quoting as described in
// https://www.gnu.org/software/grub/manual/grub/grub.html#Quoting
//
// Split treats $, ", \, \n, and ` as special within double quotes, as does
// Bash. This is slightly different from GRUB, but Grub can live with it.
func Split(s string) []string {
	ret := []string{}
	var token []rune

	var context state
	lastWhiteSpace := true
	for _, r := range s {
		quotes := context != unquoted
		switch context {
		case unquoted:
			switch r {
			case '\\':
				context = escape
				// strip out the quote
				continue
			case '\'':
				context = singleQuote
				// strip out the quote
				continue
			case '"':
				context = doubleQuote
				// strip out the quote
				continue
			case '#':
				if lastWhiteSpace {
					context = comment
					// strip out the rest
					continue
				}
			}

		case escape:
			context = unquoted

		case singleQuote:
			if r == '\'' {
				context = unquoted
				// strip out the quote
				continue
			}

		case doubleQuote:
			switch r {
			case '\\':
				context = doubleQuoteEscape
				// strip out the quote
				continue
			case '"':
				context = unquoted
				// strip out the quote
				continue
			}

		case doubleQuoteEscape:
			// GNU Bash manual:
			//
			// The backslash retains its special meaning only when
			// followed by one of the following characters: ‘$’,
			// ‘`’, ‘"’, ‘\’, or newline. Within double quotes,
			// backslashes that are followed by one of these
			// characters are removed.
			switch r {
			case '$', '"', '\\', '\n', '`': // or newline
			default:
				token = append(token, '\\')
			}

			context = doubleQuote

		case comment:
			switch r {
			case '\n':
				context = unquoted
			}

			// strip out the rest
			continue
		}

		lastWhiteSpace = unicode.IsSpace(r)

		if !lastWhiteSpace || quotes {
			token = append(token, r)
		} else if len(token) > 0 {
			ret = append(ret, string(token))
			token = token[:0]
		}
	}

	if len(token) > 0 {
		ret = append(ret, string(token))
	}
	return ret
}
