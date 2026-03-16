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

//go:build windows || linux || darwin || openbsd || freebsd || netbsd || solaris
// +build windows linux darwin openbsd freebsd netbsd solaris

package liner

import "testing"

type testItem struct {
	list   []string
	prefix string
}

func TestPrefix(t *testing.T) {
	list := []testItem{
		{[]string{"food", "foot"}, "foo"},
		{[]string{"foo", "foot"}, "foo"},
		{[]string{"food", "foo"}, "foo"},
		{[]string{"food", "foe", "foot"}, "fo"},
		{[]string{"food", "foot", "barbeque"}, ""},
		{[]string{"cafeteria", "café"}, "caf"},
		{[]string{"cafe", "café"}, "caf"},
		{[]string{"cafè", "café"}, "caf"},
		{[]string{"cafés", "café"}, "café"},
		{[]string{"áéíóú", "áéíóú"}, "áéíóú"},
		{[]string{"éclairs", "éclairs"}, "éclairs"},
		{[]string{"éclairs are the best", "éclairs are great", "éclairs"}, "éclairs"},
		{[]string{"éclair", "éclairs"}, "éclair"},
		{[]string{"éclairs", "éclair"}, "éclair"},
		{[]string{"éclair", "élan"}, "é"},
	}

	for _, test := range list {
		lcp := longestCommonPrefix(test.list)
		if lcp != test.prefix {
			t.Errorf("%s != %s for %+v", lcp, test.prefix, test.list)
		}
	}
}
