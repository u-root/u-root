// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The `field` command reads newline-separated lines of data from either
// the standard input or the specified files. It splits those lines into
// a list of fields, separated by a specifiable regular expression. It
// then prints all or a subset of those fields to the standard output.
//
// The list of output fields is specified using a grammar given in the
// parsing code, below.
//
// Options '-F' and '-O' control the input and output separators,
// respectively. The NUL character can be used as an output separator if
// the '-0' is given. The '-e' and '-E' characters contol whether empty
// fields are collapsed in the input; '-e' unconditionally preserves such
// fields, '-E' discards them. If neither is specified, a heuristic is
// applied to guess: if the input specifier is more than one character in
// length, we discard empty fields, otherwise we preserve them.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"unicode"
	"unicode/utf8"
)

type frange struct {
	begin int
	end   int
}

const (
	lastField = 0x7FFFFFFF
	cmd       = "field [ -E | -e ] [ -F regexp ] [ -0 | -O delimiter ] <field list> [file...]"
)

var flags struct {
	nuloutsep     bool
	preserveEmpty bool
	discardEmpty  bool
	insep         string
	outsep        string
}

func init() {
	defUsage := flag.Usage
	flag.Usage = func() {
		os.Args[0] = cmd
		defUsage()
	}
	flag.BoolVar(&flags.nuloutsep, "0", false, "use the NUL character ('\\0') as output separator")
	flag.BoolVar(&flags.preserveEmpty, "e", false, "preseve empty input fields")
	flag.BoolVar(&flags.discardEmpty, "E", false, "discard empty input fields")
	flag.StringVar(&flags.insep, "F", "[ \t\v\r]+", "Input separator characters (regular expression)")
	flag.StringVar(&flags.outsep, "O", " ", "Output separater (string)")
}

func main() {
	flag.Parse()

	fstate := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { fstate[f.Name] = true })
	if fstate["e"] && fstate["E"] {
		fatal("flag conflict: -e and -E are mutually exclusive")
	}
	if fstate["0"] && fstate["O"] {
		fatal("flag conflict: -O and -0 are mutually exclusive")
	}

	collapse := shouldcollapse(flags.insep)
	delim, err := regexp.Compile(flags.insep)
	if err != nil {
		fatal("Delimiter regexp failed to parse: %v", err)
	}

	if flag.NArg() == 0 {
		fatal("Range specifier missing")
	}
	rv := parseranges(flag.Arg(0))

	if flag.NArg() == 1 {
		process(os.Stdin, rv, delim, flags.outsep, collapse)
		return
	}
	for i := 1; i < flag.NArg(); i++ {
		filename := flag.Arg(i)
		if filename == "-" {
			process(os.Stdin, rv, delim, flags.outsep, collapse)
			continue
		}
		file, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot open file %q: %v\n", filename, err)
			continue
		}
		process(file, rv, delim, flags.outsep, collapse)
		file.Close()
	}
}

func shouldcollapse(s string) bool {
	if flags.preserveEmpty {
		return false
	}
	if flags.discardEmpty {
		return true
	}
	l := utf8.RuneCountInString(s)
	r, _ := utf8.DecodeRuneInString(s)
	return l > 1 && (l != 2 || r != '\\')
}

// The field selection syntax is:
//
// ranges := range [[delim] range]
// range := field | NUM '-' [field]
// field := NUM | NF
// delim := ws+ | '|' | ','
// ws := c such that `isspace(c)` is true.
// NF := 'NF' | 'N'
// (Numbers can be negative)

func parseranges(input string) []frange {
	var rs []frange
	lex := &lexer{input: input}
	if input == "" {
		fatal("Empty field range")
	}
	lex.next()
	for {
		if lex.peektype() == tokSpace {
			lex.next()
		}
		r := parserange(lex)
		rs = append(rs, r)
		typ := lex.peektype()
		if typ == tokEOF {
			break
		}
		if !isdelim(typ) {
			fatal("Syntax error in field list, tok = %s", lex.peektok())
		}
		lex.next()
	}
	return rs
}

func parserange(lex *lexer) frange {
	r := frange{begin: lastField, end: lastField}
	if lex.peektype() == tokEOF {
		fatal("EOF at start of range")
	}
	fnum, typ := parsefield(lex)
	r.begin = fnum
	r.end = fnum
	if typ == tokNF {
		return r
	}
	typ = lex.peektype()
	if typ != tokDash {
		return r
	}
	lex.next()
	r.end = lastField
	typ = lex.peektype()
	if typ != tokEOF && !isdelim(typ) {
		r.end, _ = parsefield(lex)
	}
	return r
}

func parsefield(lex *lexer) (int, toktype) {
	typ := lex.peektype()
	if typ == tokNF {
		lex.next()
		return lastField, tokNF
	}
	return parsenum(lex), tokNum
}

func parsenum(lex *lexer) int {
	tok, typ := lex.next()
	if typ == tokEOF {
		fatal("EOF in number parser")
	}
	if typ == tokNum {
		num, _ := strconv.Atoi(tok)
		return num
	}
	if typ != tokDash {
		fatal("number parser error: unexpected token '%v'", tok)
	}
	tok, typ = lex.next()
	if typ == tokEOF {
		fatal("negative number parse error: unexpected EOF")
	}
	if typ != tokNum {
		fatal("number parser error: bad lexical token '%v'", tok)
	}
	num, _ := strconv.Atoi(tok)
	return -num
}

func isdelim(typ toktype) bool {
	return typ == tokComma || typ == tokPipe || typ == tokSpace
}

type toktype int

const (
	tokError toktype = iota
	tokEOF
	tokComma
	tokPipe
	tokDash
	tokNum
	tokSpace
	tokNF

	eof = -1
)

type lexer struct {
	input string
	tok   string
	typ   toktype
	start int
	pos   int
	width int
}

func (lex *lexer) peek() (string, toktype) {
	return lex.tok, lex.typ
}

func (lex *lexer) peektype() toktype {
	return lex.typ
}

func (lex *lexer) peektok() string {
	return lex.tok
}

func (lex *lexer) next() (string, toktype) {
	tok, typ := lex.peek()
	lex.tok, lex.typ = lex.scan()
	return tok, typ
}

func (lex *lexer) scan() (string, toktype) {
	switch r := lex.nextrune(); {
	case r == eof:
		return "", tokEOF
	case r == ',':
		return lex.token(), tokComma
	case r == '|':
		return lex.token(), tokPipe
	case r == '-':
		return lex.token(), tokDash
	case r == 'N':
		lex.consume()
		r = lex.nextrune()
		if r == 'F' {
			lex.consume()
		}
		lex.ignore()
		return lex.token(), tokNF
	case unicode.IsDigit(r):
		for r := lex.nextrune(); unicode.IsDigit(r); r = lex.nextrune() {
			lex.consume()
		}
		lex.ignore()
		return lex.token(), tokNum
	case unicode.IsSpace(r):
		for r := lex.nextrune(); unicode.IsSpace(r); r = lex.nextrune() {
			lex.consume()
		}
		lex.ignore()
		return lex.token(), tokSpace
	default:
		fatal("Lexical error at character '%v'", r)
	}
	return "", tokError
}

func (lex *lexer) nextrune() (r rune) {
	if lex.pos >= len(lex.input) {
		lex.width = 0
		return eof
	}
	r, lex.width = utf8.DecodeRuneInString(lex.input[lex.pos:])
	return r
}

func (lex *lexer) consume() {
	lex.pos += lex.width
	lex.width = 0
}

func (lex *lexer) ignore() {
	lex.width = 0
}

func (lex *lexer) token() string {
	lex.consume()
	tok := lex.input[lex.start:lex.pos]
	lex.start = lex.pos
	return tok
}

func process(file *os.File, rv []frange, delim *regexp.Regexp, outsep string, collapse bool) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		prefix := ""
		printed := false
		line := scanner.Text()
		fields := split(line, delim, collapse)
		for _, r := range rv {
			begin, end := r.begin, r.end
			switch {
			case begin == 0:
				pprefix(prefix)
				prefix = outsep
				fmt.Print(line)
				printed = true
			case begin == lastField:
				begin = len(fields) - 1
			case begin < 0:
				begin += len(fields)
			default:
				begin--
			}
			if end < 0 {
				end += len(fields) + 1
			}
			if begin < 0 || end < 0 || end < begin || len(fields) < begin {
				continue
			}
			for i := begin; i < end && i < len(fields); i++ {
				pprefix(prefix)
				prefix = outsep
				fmt.Print(fields[i])
				printed = true
			}
		}
		if printed || !collapse {
			fmt.Println()
		}
	}
	err := scanner.Err()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func split(s string, delim *regexp.Regexp, collapse bool) []string {
	sv := delim.Split(s, -1)
	if !collapse {
		return sv
	}
	rv := []string{}
	for _, s := range sv {
		if s != "" {
			rv = append(rv, s)
		}
	}
	return rv
}

func pprefix(prefix string) {
	if prefix == "" {
		return
	}
	if flags.nuloutsep {
		fmt.Print("\x00")
	} else {
		fmt.Print(prefix)
	}
}

func fatal(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", a...)
	flag.Usage()
	os.Exit(1)
}
