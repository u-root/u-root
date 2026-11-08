// SPDX-License-Identifier: MIT

// Copyright © 2012 Peter Harris
//
// Permission is hereby granted, free of charge, to any person obtaining a
// copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice (including the next
// paragraph) shall be included in all copies or substantial portions of the
// Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL
// THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

package liner

import (
	"github.com/mattn/go-runewidth"
	"github.com/rivo/uniseg"
)

// glyphs splits s into the sequences of runes a terminal draws as one
// character each, by the extended grapheme cluster rules of UAX #29.
//
// A glyph is often one rune and need not be. A letter followed by a combining
// accent is one, and so are the emoji built out of several runes: a flag is a
// pair of regional indicators, a thumb carries its skin tone in a second rune,
// and a family is four people joined by zero-width joiners.
func glyphs(s []rune) [][]rune {
	var out [][]rune
	g := uniseg.NewGraphemes(string(s))
	for g.Next() {
		out = append(out, []rune(g.Str()))
	}
	return out
}

// countGlyphs considers zero-width characters to be zero glyphs wide,
// and members of Chinese, Japanese, and Korean scripts to be 2 glyphs wide.
func countGlyphs(s []rune) int {
	n := 0
	for _, r := range s {
		// speed up the common case
		if r < 127 {
			n++
			continue
		}

		n += runewidth.RuneWidth(r)
	}
	return n
}

func countMultiLineGlyphs(s []rune, columns int, start int) int {
	n := start
	for _, r := range s {
		if r < 127 {
			n++
			continue
		}
		switch runewidth.RuneWidth(r) {
		case 0:
		case 1:
			n++
		case 2:
			n += 2
			// no room for a 2-glyphs-wide char in the ending
			// so skip a column and display it at the beginning
			if n%columns == 1 {
				n++
			}
		}
	}
	return n
}

func getPrefixGlyphs(s []rune, num int) []rune {
	g := glyphs(s)
	p := 0
	for n := 0; n < num && n < len(g); n++ {
		p += len(g[n])
	}
	return s[:p]
}

func getSuffixGlyphs(s []rune, num int) []rune {
	g := glyphs(s)
	p := len(s)
	for n := 0; n < num && n < len(g); n++ {
		p -= len(g[len(g)-1-n])
	}
	return s[p:]
}
