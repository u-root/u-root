// Copyright 2017-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package shlex is a POSIX command-line argument parser.
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

func isWhitespace(b byte) bool {
	return b == '\t' || b == '\n' || b == '\v' ||
		b == '\f' || b == '\r' || b == ' '
}

type quote uint8

const (
	unquoted quote = iota
	escape
	singleQuote
	doubleQuote
	doubleQuoteEscape
	comment
)

// Argv splits a command line according to usual simple shell rules.
//
// Argv was written from the spec of Grub quoting at
// https://www.gnu.org/software/grub/manual/grub/grub.html#Quoting
// except that the escaping of newline is not supported
func Argv(s string) []string {
	ret := []string{}
	var token []rune

	var context quote
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
			case '$', '"', '\\', '\n': // or newline
			default:
				token = append(token, '\\')
			}

			context = doubleQuote

		case comment:
			// should end on newline

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
