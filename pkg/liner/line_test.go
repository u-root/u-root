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
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestAppend(t *testing.T) {
	var s State
	s.AppendHistory("foo")
	s.AppendHistory("bar")

	var out bytes.Buffer
	num, err := s.WriteHistory(&out)
	if err != nil {
		t.Fatal("Unexpected error writing history", err)
	}
	if num != 2 {
		t.Fatalf("Expected 2 history entries, got %d", num)
	}

	s.AppendHistory("baz")
	num, err = s.WriteHistory(&out)
	if err != nil {
		t.Fatal("Unexpected error writing history", err)
	}
	if num != 3 {
		t.Fatalf("Expected 3 history entries, got %d", num)
	}

	s.AppendHistory("baz")
	num, err = s.WriteHistory(&out)
	if err != nil {
		t.Fatal("Unexpected error writing history", err)
	}
	if num != 3 {
		t.Fatalf("Expected 3 history entries after duplicate append, got %d", num)
	}

	s.AppendHistory("baz")

}

func TestHistory(t *testing.T) {
	input := `foo
bar
baz
quux
dingle`

	var s State
	num, err := s.ReadHistory(strings.NewReader(input))
	if err != nil {
		t.Fatal("Unexpected error reading history", err)
	}
	if num != 5 {
		t.Fatal("Wrong number of history entries read")
	}

	var out bytes.Buffer
	num, err = s.WriteHistory(&out)
	if err != nil {
		t.Fatal("Unexpected error writing history", err)
	}
	if num != 5 {
		t.Fatal("Wrong number of history entries written")
	}
	if strings.TrimSpace(out.String()) != input {
		t.Fatal("Round-trip failure")
	}

	// clear the history and re-write
	s.ClearHistory()
	num, err = s.WriteHistory(&out)
	if err != nil {
		t.Fatal("Unexpected error writing history", err)
	}
	if num != 0 {
		t.Fatal("Wrong number of history entries written, expected none")
	}
	// Test reading with a trailing newline present
	var s2 State
	num, err = s2.ReadHistory(&out)
	if err != nil {
		t.Fatal("Unexpected error reading history the 2nd time", err)
	}
	if num != 5 {
		t.Fatal("Wrong number of history entries read the 2nd time")
	}

	num, err = s.ReadHistory(strings.NewReader(input + "\n\xff"))
	if err == nil {
		t.Fatal("Unexpected success reading corrupted history", err)
	}
	if num != 5 {
		t.Fatal("Wrong number of history entries read the 3rd time")
	}
}

func TestColumns(t *testing.T) {
	list := []string{"foo", "food", "This entry is quite a bit longer than the typical entry"}

	output := []struct {
		width, columns, rows, maxWidth int
	}{
		{80, 1, 3, len(list[2]) + 1},
		{120, 2, 2, len(list[2]) + 1},
		{800, 14, 1, 0},
		{8, 1, 3, 7},
	}

	for i, o := range output {
		col, row, maxwidth := calculateColumns(o.width, list)
		if col != o.columns {
			t.Fatalf("Wrong number of columns, %d != %d, in TestColumns %d\n", col, o.columns, i)
		}
		if row != o.rows {
			t.Fatalf("Wrong number of rows, %d != %d, in TestColumns %d\n", row, o.rows, i)
		}
		if maxwidth != o.maxWidth {
			t.Fatalf("Wrong column width, %d != %d, in TestColumns %d\n", maxwidth, o.maxWidth, i)
		}
	}
}

// This example demonstrates a way to retrieve the current
// history buffer without using a file.
func ExampleState_WriteHistory() {
	var s State
	s.AppendHistory("foo")
	s.AppendHistory("bar")

	buf := new(bytes.Buffer)
	_, err := s.WriteHistory(buf)
	if err == nil {
		history := strings.Split(strings.TrimSpace(buf.String()), "\n")
		for i, line := range history {
			fmt.Println("History entry", i, ":", line)
		}
	}
	// Output:
	// History entry 0 : foo
	// History entry 1 : bar
}
