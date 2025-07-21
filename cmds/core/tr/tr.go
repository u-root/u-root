// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// tr - translate or delete characters

// Synopsis:
//     tr [OPTION]... SET1 [SET2]

// Description:
//     Translate, squeeze, and/or delete characters from standard input, writing
//     to standard output.
//
//     -d, --delete: delete characters in SET1, do not translate
//
// SETs  are  specified  as  strings of characters. Most represent themselves.
// Interpreted sequences are:
//     \\        backslash
//     \a        audible BEL
//     \b        backspace
//     \f        form feed
//     \n        new line
//     \r        return
//     \t        horizontal tab
//     \v        vertical tab
//     [:alnum:] all letters and digits
//     [:alpha:] all letters
//     [:digit:] all digits
//     [:graph:] all printable characters
//     [:cntrl:] all control characters
//     [:lower:] all lower case letters
//     [:upper:] all upper case letters
//     [:space:] all whitespaces

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"unicode"
)

var del bool

type cmd struct {
	stdin  io.Reader
	stdout io.Writer
	del    bool
	tr     *transformer
}

func command(in io.Reader, out io.Writer, args []string, del bool) (*cmd, error) {
	tr, err := parse(args, del)
	if err != nil {
		return nil, err
	}

	return &cmd{
		stdin:  in,
		stdout: out,
		del:    del,
		tr:     tr,
	}, nil
}

const name = "tr"

var escapeChars = map[rune]rune{
	'\\': '\\',
	'a':  '\a',
	'b':  '\b',
	'f':  '\f',
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
	'v':  '\v',
}

type Set string

const (
	ALPHA Set = "[:alpha:]"
	DIGIT Set = "[:digit:]"
	GRAPH Set = "[:graph:]"
	CNTRL Set = "[:cntrl:]"
	PUNCT Set = "[:punct:]"
	SPACE Set = "[:space:]"
	ALNUM Set = "[:alnum:]"
	LOWER Set = "[:lower:]"
	UPPER Set = "[:upper:]"
)

var sets = map[Set]func(r rune) bool{
	ALNUM: func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsDigit(r)
	},

	ALPHA: unicode.IsLetter,
	DIGIT: unicode.IsDigit,
	GRAPH: unicode.IsGraphic,
	CNTRL: unicode.IsControl,
	PUNCT: unicode.IsPunct,
	SPACE: unicode.IsSpace,
	LOWER: unicode.IsLower,
	UPPER: unicode.IsUpper,
}

type transformer struct {
	transform func(r rune) rune
}

func setToRune(s Set, outRune rune) *transformer {
	check := sets[s]
	return &transformer{
		transform: func(r rune) rune {
			if check(r) {
				return outRune
			}
			return r
		},
	}
}

func lowerToUpper() *transformer {
	return &transformer{
		transform: func(r rune) rune {
			return unicode.ToUpper(r)
		},
	}
}

func upperToLower() *transformer {
	return &transformer{
		transform: func(r rune) rune {
			return unicode.ToLower(r)
		},
	}
}

func runesToRunes(in []rune, out ...rune) *transformer {
	convs := make(map[rune]rune)
	l := len(out)
	for i, r := range in {
		ind := min(i, l-1)
		convs[r] = out[ind]
	}
	return &transformer{
		transform: func(r rune) rune {
			if outRune, ok := convs[r]; ok {
				return outRune
			}
			return r
		},
	}
}

func (c *cmd) run() error {
	in := bufio.NewReader(c.stdin)
	out := bufio.NewWriter(c.stdout)

	defer out.Flush()

	for {
		inRune, size, err := in.ReadRune()
		if inRune == unicode.ReplacementChar {
			// can skip error handling here, because
			// previous operation was in.ReadRune()
			in.UnreadRune()

			b, err := in.ReadByte()
			if err != nil {
				return fmt.Errorf("read error: %w", err)
			}

			if err := out.WriteByte(b); err != nil {
				return fmt.Errorf("write error: %w", err)
			}
		} else if size > 0 {
			if outRune := c.tr.transform(inRune); outRune != unicode.ReplacementChar {
				if _, err := out.WriteRune(outRune); err != nil {
					return fmt.Errorf("write error: %w", err)
				}
			}
		}

		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

func parse(args []string, del bool) (*transformer, error) {
	narg := len(args)

	switch {
	case narg == 0 || (narg == 1 && !del):
		return nil, fmt.Errorf("missing operand")
	case narg > 1 && del:
		return nil, fmt.Errorf("extra operand after %q", args[0])
	case narg > 2:
		return nil, fmt.Errorf("extra operand after %q", args[1])
	}

	set1 := Set(args[0])
	arg1, err := unescape(set1)
	if err != nil {
		return nil, err
	}

	var set2 Set
	if del {
		set2 = Set(unicode.ReplacementChar)
	} else {
		set2 = Set(args[1])
	}

	if set1 == LOWER && set2 == UPPER {
		return lowerToUpper(), nil
	}
	if set1 == UPPER && set2 == LOWER {
		return upperToLower(), nil
	}

	if (set2 == LOWER || set2 == UPPER) && (set1 != LOWER && set1 != UPPER) ||
		(set1 == LOWER && set2 == LOWER) || (set1 == UPPER && set2 == UPPER) {
		return nil, fmt.Errorf("misaligned [:upper:] and/or [:lower:] construct")
	}

	if _, ok := sets[set2]; ok {
		return nil, fmt.Errorf(`the only character classes that may appear in SET2 are 'upper' and 'lower'`)
	}

	arg2, err := unescape(set2)
	if err != nil {
		return nil, err
	}
	if len(arg2) == 0 {
		return nil, fmt.Errorf("SET2 must be non-empty")
	}
	if _, ok := sets[set1]; ok {
		return setToRune(set1, arg2[0]), nil
	}
	return runesToRunes(arg1, arg2...), nil
}

func unescape(s Set) ([]rune, error) {
	var out []rune
	var escape bool
	for _, r := range s {
		if escape {
			v, ok := escapeChars[r]
			if !ok {
				return nil, fmt.Errorf("unknown escape sequence '\\%c'", r)
			}
			out = append(out, v)
			escape = false
			continue
		}

		if r == '\\' {
			escape = true
			continue
		}

		out = append(out, r)
	}
	return out, nil
}

func main() {
	flag.BoolVar(&del, "delete", false, "delete characters in SET1, do not translate")
	flag.BoolVar(&del, "d", false, "delete characters in SET1, do not translate (shorthand)")

	flag.Parse()

	cmd, err := command(os.Stdin, os.Stdout, flag.Args(), del)
	if err != nil {
		log.Fatalf("%s: %v\n", name, err)
	}
	if err := cmd.run(); err != nil {
		log.Fatalf("%s: %v\n", name, err)
	}
}
