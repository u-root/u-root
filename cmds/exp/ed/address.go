// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// address.go - contains methods for FileBuffer for line address resolution
package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// regex helpers
func reGroup(s string) string {
	return "(" + s + ")"
}

func reOr(s ...string) string {
	return strings.Join(s, "|")
}

func reOpt(s string) string {
	return s + "?"
}

func reStart(s string) string {
	return "^" + s
}

// addr regex strings
var (
	reWhitespace   = "(\\s)"
	reSingleSymbol = "([.$])"
	reNumber       = "([0-9]+)"
	reOffset       = reGroup("([-+])" + reOpt(reNumber))
	reMark         = "'([a-z])"
	reRE           = "(\\/((?:\\\\/|[^\\/])*)\\/|\\?((?:\\\\?|[^\\?])*)\\?)"
	reSingle       = reGroup(reOr(reSingleSymbol, reNumber, reOffset, reMark, reRE))
	reOff          = "(?:\\s+" + reGroup(reOr(reNumber, reOffset)) + ")"
)

// addr compiled regexes
var (
	rxWhitespace   = regexp.MustCompile(reStart(reWhitespace))
	rxSingleSymbol = regexp.MustCompile(reStart(reSingleSymbol))
	rxNumber       = regexp.MustCompile(reStart(reNumber))
	rxOffset       = regexp.MustCompile(reStart(reOffset))
	rxMark         = regexp.MustCompile(reStart(reMark))
	rxRE           = regexp.MustCompile(reStart(reRE))
	rxSingle       = regexp.MustCompile(reStart(reSingle))
	rxOff          = regexp.MustCompile(reStart(reOff))
)

// ResolveOffset resolves an offset to an addr
func (f *FileBuffer) ResolveOffset(cmd string) (offset, cmdOffset int, e error) {
	ms := rxOff.FindStringSubmatch(cmd)
	// 0: full
	// 1: without whitespace
	// 2: num
	// 3: offset full
	// 4: offset +/-
	// 5: offset num
	if len(ms) == 0 {
		return
	}
	n := 1
	cmdOffset = len(ms[0])
	switch {
	case len(ms[2]) > 0:
		// num
		if n, e = strconv.Atoi(ms[2]); e != nil {
			return
		}
		offset = n
	case len(ms[3]) > 0:
		// offset
		if len(ms[5]) > 0 {
			if n, e = strconv.Atoi(ms[5]); e != nil {
				return
			}
		}
		switch ms[4][0] {
		case '+':
			offset = n
		case '-':
			offset = -n
		}
	}
	return
}

// ResolveAddr resolves a command address from a cmd string
// - makes no attempt to verify that the resulting addr is valid
func (f *FileBuffer) ResolveAddr(cmd string) (line, cmdOffset int, e error) {
	line = f.GetAddr()
	cmdOffset = 0
	m := rxSingle.FindString(cmd)
	if len(m) == 0 {
		// no match
		return
	}
	cmdOffset = len(m)
	switch {
	case rxSingleSymbol.MatchString(m):
		// no need to rematch; these are all single char
		switch m[0] {
		case '.':
			// current
		case '$':
			// last
			line = f.Len() - 1
		}
	case rxNumber.MatchString(m):
		var n int
		ns := rxNumber.FindString(m)
		if n, e = strconv.Atoi(ns); e != nil {
			return
		}
		line = n - 1
	case rxOffset.MatchString(m):
		n := 1
		ns := rxOffset.FindStringSubmatch(m)
		if len(ns[3]) > 0 {
			if n, e = strconv.Atoi(ns[3]); e != nil {
				return
			}
		}
		switch ns[2][0] {
		case '+':
			line = f.GetAddr() + n
		case '-':
			line = f.GetAddr() - n
		}
	case rxMark.MatchString(m):
		c := m[1] // len should already be verified by regexp
		line, e = buffer.GetMark(c)
	case rxRE.MatchString(m):
		r := rxRE.FindAllStringSubmatch(m, -1)
		// 0: full
		// 1: regexp w/ delim
		// 2: regexp
		if len(r) < 1 || len(r[0]) < 3 {
			e = fmt.Errorf("invalid regexp: %s", m)
			return
		}
		restr := r[0][2]
		sign := 1
		if r[0][0][0] == '?' {
			sign = -1
			restr = r[0][3]
		}
		var re *regexp.Regexp
		if re, e = regexp.Compile(restr); e != nil {
			e = fmt.Errorf("invalid regexp: %w", e)
			return
		}
		var c int
		for i := 0; i < f.Len(); i++ {
			c = (sign*i + f.GetAddr() + f.Len()) % f.Len()
			if re.MatchString(f.GetMust(c, false)) {
				line = c
				return
			}
		}
		e = fmt.Errorf("regexp not found: %s", restr)
	}
	if e != nil {
		return
	}

	for {
		var off, cmdoff int
		if off, cmdoff, e = f.ResolveOffset(cmd[cmdOffset:]); e != nil {
			return
		}
		if cmdoff > 0 {
			// we got an offset
			line += off
			cmdOffset += cmdoff
		} else {
			// no more offsets
			break
		}
	}
	return
}

// ResolveAddrs resolves all addrs at the begining of a line
// - makes no attempt to verify that the resulting addrs are valid
// - will always return at least one addr as long as there isn't an error
// - if an error is reached, return value behavior is undefined
func (f *FileBuffer) ResolveAddrs(cmd string) (lines []int, cmdOffset int, e error) {
	var line, off int

Loop:
	for cmdOffset < len(cmd) {
		cmdOffset += wsOffset(cmd[cmdOffset:])
		if line, off, e = f.ResolveAddr(cmd[cmdOffset:]); e != nil {
			return
		}
		lines = append(lines, line)
		cmdOffset += off
		cmdOffset += wsOffset(cmd[cmdOffset:])
		if len(cmd)-1 <= cmdOffset {
			return
		}
		switch cmd[cmdOffset] { // do we have more addrs?
		case ',':
			cmdOffset++
		case ';':
			// we're  the left side of a ; set the current addr
			if e = f.SetAddr(line); e != nil {
				return
			}
			cmdOffset++
		case '%':
			lines = append(lines, 0, f.Len()-1)
			cmdOffset++
			cmdOffset += wsOffset(cmd[cmdOffset:])
			return
		default:
			break Loop
		}
	}
	return
}

// wsOffset is a helper to find the offset to skip whitespace
func wsOffset(cmd string) (o int) {
	o = 0
	ws := rxWhitespace.FindStringIndex(cmd)
	if ws != nil {
		o = ws[1]
	}
	return
}

/*
 * The following three functions, AddrValue, AddrRange, AddrRangeOrLine all get the specified
 * type of address from an already parsed address.
 */

// AddrValue gets the resolved single-line address
func (f *FileBuffer) AddrValue(addrs []int) (r int, e error) {
	if len(addrs) == 0 {
		e = ErrINV
		return
	}
	r = addrs[len(addrs)-1]
	if f.OOB(r) {
		e = ErrOOB
	}
	return
}

// AddrRange gets and address range, fails if we don't have a range specified
func (f *FileBuffer) AddrRange(addrs []int) (r [2]int, e error) {
	switch len(addrs) {
	case 0:
		e = ErrINV
		return
	case 1:
		r[0] = f.GetAddr()
		r[1] = addrs[0]
	default:
		r[0] = addrs[len(addrs)-2]
		r[1] = addrs[len(addrs)-1]
	}
	if f.OOB(r[0]) || f.OOB(r[1]) {
		e = ErrOOB
	}
	if r[0] > r[1] {
		e = fmt.Errorf("address out of order")
	}
	return
}

// AddrRangeOrLine returns a range or single line if no range could be resolved
func (f *FileBuffer) AddrRangeOrLine(addrs []int) (r [2]int, e error) {
	if len(addrs) > 1 {
		// delete a range
		if r, e = buffer.AddrRange(addrs); e != nil {
			return
		}
	} else {
		// delete a line
		if r[0], e = buffer.AddrValue(addrs); e != nil {
			return
		}
		r[1] = r[0]
	}
	return
}
