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
	"strconv"
	"testing"
)

func accent(in []rune) []rune {
	var out []rune
	for _, r := range in {
		out = append(out, r)
		out = append(out, '\u0301')
	}
	return out
}

type testCase struct {
	s      []rune
	glyphs int
}

var testCases = []testCase{
	{[]rune("query"), 5},
	{[]rune("私"), 2},
	{[]rune("hello『世界』"), 13},
}

func TestCountGlyphs(t *testing.T) {
	for _, testCase := range testCases {
		count := countGlyphs(testCase.s)
		if count != testCase.glyphs {
			t.Errorf("ASCII count incorrect. %d != %d", count, testCase.glyphs)
		}
		count = countGlyphs(accent(testCase.s))
		if count != testCase.glyphs {
			t.Errorf("Accent count incorrect. %d != %d", count, testCase.glyphs)
		}
	}
}

func compare(a, b []rune, name string, t *testing.T) {
	if len(a) != len(b) {
		t.Errorf(`"%s" != "%s" in %s"`, string(a), string(b), name)
		return
	}
	for i := range a {
		if a[i] != b[i] {
			t.Errorf(`"%s" != "%s" in %s"`, string(a), string(b), name)
			return
		}
	}
}

func TestPrefixGlyphs(t *testing.T) {
	for _, testCase := range testCases {
		for i := 0; i <= len(testCase.s); i++ {
			iter := strconv.Itoa(i)
			out := getPrefixGlyphs(testCase.s, i)
			compare(out, testCase.s[:i], "ascii prefix "+iter, t)
			out = getPrefixGlyphs(accent(testCase.s), i)
			compare(out, accent(testCase.s[:i]), "accent prefix "+iter, t)
		}
		out := getPrefixGlyphs(testCase.s, 999)
		compare(out, testCase.s, "ascii prefix overflow", t)
		out = getPrefixGlyphs(accent(testCase.s), 999)
		compare(out, accent(testCase.s), "accent prefix overflow", t)

		out = getPrefixGlyphs(testCase.s, -3)
		if len(out) != 0 {
			t.Error("ascii prefix negative")
		}
		out = getPrefixGlyphs(accent(testCase.s), -3)
		if len(out) != 0 {
			t.Error("accent prefix negative")
		}
	}
}

func TestSuffixGlyphs(t *testing.T) {
	for _, testCase := range testCases {
		for i := 0; i <= len(testCase.s); i++ {
			iter := strconv.Itoa(i)
			out := getSuffixGlyphs(testCase.s, i)
			compare(out, testCase.s[len(testCase.s)-i:], "ascii suffix "+iter, t)
			out = getSuffixGlyphs(accent(testCase.s), i)
			compare(out, accent(testCase.s[len(testCase.s)-i:]), "accent suffix "+iter, t)
		}
		out := getSuffixGlyphs(testCase.s, 999)
		compare(out, testCase.s, "ascii suffix overflow", t)
		out = getSuffixGlyphs(accent(testCase.s), 999)
		compare(out, accent(testCase.s), "accent suffix overflow", t)

		out = getSuffixGlyphs(testCase.s, -3)
		if len(out) != 0 {
			t.Error("ascii suffix negative")
		}
		out = getSuffixGlyphs(accent(testCase.s), -3)
		if len(out) != 0 {
			t.Error("accent suffix negative")
		}
	}
}

// Glyphs a terminal draws as one character out of several runes. An accent is
// the case the tests above already cover; these are the rest of them.
var multiRuneGlyphs = []struct {
	name string
	s    string
}{
	{"family", "\U0001F468‍\U0001F469‍\U0001F467‍\U0001F466"},
	{"rainbow flag", "\U0001F3F3️‍\U0001F308"},
	{"thumb with a skin tone", "\U0001F44D\U0001F3FD"},
	{"flag", "\U0001F1EF\U0001F1F5"},
}

// Moving over one glyph must move over all of it. Stopping partway leaves the
// cursor inside the character, and the runes either side of it then get
// edited separately.
func TestGlyphMovementKeepsGlyphsWhole(t *testing.T) {
	for _, g := range multiRuneGlyphs {
		s := []rune(g.s)
		if out := getPrefixGlyphs(s, 1); len(out) != len(s) {
			t.Errorf("%s: prefix of one glyph took %d of its %d runes",
				g.name, len(out), len(s))
		}
		if out := getSuffixGlyphs(s, 1); len(out) != len(s) {
			t.Errorf("%s: suffix of one glyph took %d of its %d runes",
				g.name, len(out), len(s))
		}

		// One glyph of "a<glyph>b" is the "a" at the front and the "b" at the
		// back, never a piece of the glyph between them.
		padded := append(append([]rune("a"), s...), 'b')
		if out := getPrefixGlyphs(padded, 1); len(out) != 1 {
			t.Errorf("%s: prefix of one glyph of %q took %d runes, wanted 1",
				g.name, string(padded), len(out))
		}
		if out := getSuffixGlyphs(padded, 1); len(out) != 1 {
			t.Errorf("%s: suffix of one glyph of %q took %d runes, wanted 1",
				g.name, string(padded), len(out))
		}
	}
}
