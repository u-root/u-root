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

//go:build !windows

package liner

import (
	"bufio"
	"bytes"
	"testing"
)

func (s *State) expectRune(t *testing.T, r rune) {
	item, err := s.readNext()
	if err != nil {
		t.Fatalf("Expected rune '%c', got error %s\n", r, err)
	}
	if v, ok := item.(rune); !ok {
		t.Fatalf("Expected rune '%c', got non-rune %v\n", r, v)
	} else {
		if v != r {
			t.Fatalf("Expected rune '%c', got rune '%c'\n", r, v)
		}
	}
}

func (s *State) expectAction(t *testing.T, a action) {
	item, err := s.readNext()
	if err != nil {
		t.Fatalf("Expected Action %d, got error %s\n", a, err)
	}
	if v, ok := item.(action); !ok {
		t.Fatalf("Expected Action %d, got non-Action %v\n", a, v)
	} else {
		if v != a {
			t.Fatalf("Expected Action %d, got Action %d\n", a, v)
		}
	}
}

func TestTypes(t *testing.T) {
	input := []byte{'A', 27, 'B', 27, 91, 68, 27, '[', '1', ';', '5', 'D', 'e'}
	var s State
	s.r = bufio.NewReader(bytes.NewBuffer(input))

	next := make(chan nexter)
	go func() {
		for {
			var n nexter
			n.r, _, n.err = s.r.ReadRune()
			next <- n
		}
	}()
	s.next = next

	s.expectRune(t, 'A')
	s.expectRune(t, 27)
	s.expectRune(t, 'B')
	s.expectAction(t, left)
	s.expectAction(t, wordLeft)

	s.expectRune(t, 'e')
}
