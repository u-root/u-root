// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var errDiskCrashed = errors.New("disk crashed")

func TestTsort(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		stdin      io.Reader
		wantStdout string
		wantStderr string
		wantErr    error
	}{
		{
			name:  "stdin: empty",
			stdin: strings.NewReader(""),
		},
		{
			name:  "stdin: space",
			stdin: strings.NewReader(" "),
		},
		{
			name:  `stdin: \n`,
			stdin: strings.NewReader("\n"),
		},
		{
			name:  `stdin: \r`,
			stdin: strings.NewReader("\r"),
		},
		{
			name:  `stdin: \t`,
			stdin: strings.NewReader("\t"),
		},
		{
			name:  `stdin: \v`,
			stdin: strings.NewReader("\v"),
		},
		{
			name:  `stdin: \f`,
			stdin: strings.NewReader("\f"),
		},
		{
			name:       "stdin: one node: a",
			stdin:      strings.NewReader("a a"),
			wantStdout: "a\n",
		},
		{
			name:       "stdin: one node: b",
			stdin:      strings.NewReader("b b"),
			wantStdout: "b\n",
		},
		{
			name:       "stdin: one edge: a b",
			stdin:      strings.NewReader("a b"),
			wantStdout: "a\nb\n",
		},
		{
			name:       "stdin: one edge: b a",
			stdin:      strings.NewReader("b a"),
			wantStdout: "b\na\n",
		},
		{
			name:       "stdin: one edge: a A",
			stdin:      strings.NewReader("a A"),
			wantStdout: "a\nA\n",
		},
		{
			name:       "stdin: one edge: a  b",
			stdin:      strings.NewReader("a  b"),
			wantStdout: "a\nb\n",
		},
		{
			name:    "stdin: odd data count: 1",
			stdin:   strings.NewReader("a"),
			wantErr: errOddDataCount,
		},
		{
			name:    "stdin: odd data count: 3",
			stdin:   strings.NewReader("a b c"),
			wantErr: errOddDataCount,
		},
		{
			name:    "stdin: odd data count: 5",
			stdin:   strings.NewReader("a b c d e"),
			wantErr: errOddDataCount,
		},
		{
			name:       `stdin: one edge: a\nb`,
			stdin:      strings.NewReader("a\nb"),
			wantStdout: "a\nb\n",
		},
		{
			name:       `stdin: one edge: a\n\nb`,
			stdin:      strings.NewReader("a\n\nb"),
			wantStdout: "a\nb\n",
		},
		{
			name:       `stdin: one edge: \na\nb`,
			stdin:      strings.NewReader("\na\nb"),
			wantStdout: "a\nb\n",
		},
		{
			name:       `stdin: one edge: a\nb\n`,
			stdin:      strings.NewReader("a\nb\n"),
			wantStdout: "a\nb\n",
		},
		{
			name:       `stdin: one edge: a\rb`,
			stdin:      strings.NewReader("a\rb"),
			wantStdout: "a\nb\n",
		},
		{
			name:       `stdin: one edge: a\tb`,
			stdin:      strings.NewReader("a\tb"),
			wantStdout: "a\nb\n",
		},
		{
			name:       `stdin: one edge: a\vb`,
			stdin:      strings.NewReader("a\vb"),
			wantStdout: "a\nb\n",
		},
		{
			name:       `stdin: one edge: a\fb`,
			stdin:      strings.NewReader("a\fb"),
			wantStdout: "a\nb\n",
		},
		{
			name:    "stdin: error-returning stdin",
			stdin:   &errorReturningReader{err: errDiskCrashed},
			wantErr: errDiskCrashed,
		},
		{
			name:       "file: line: a b b c c d",
			args:       []string{tempFile(t, "a b b c c d")},
			wantStdout: "a\nb\nc\nd\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout := new(strings.Builder)
			stderr := new(strings.Builder)

			gotErr := run(tt.stdin, stdout, stderr, tt.args...)

			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf(`run() gotErr = %q, want %q`, gotErr, tt.wantErr)
			}

			gotStderr := stderr.String()
			if gotStderr != tt.wantStderr {
				t.Errorf(
					"run() gotStderr = %q, want %q",
					gotStderr,
					tt.wantStderr)
			}

			gotStdout := stdout.String()
			if gotStdout != tt.wantStdout {
				t.Errorf(
					"run() gotStdout = %q, want %q",
					gotStdout,
					tt.wantStdout)
			}
		})
	}

	t.Run("file: non-existent", func(t *testing.T) {
		var stdin io.Reader
		stdout := new(strings.Builder)
		stderr := new(strings.Builder)

		gotErr := run(stdin, stdout, stderr, "non-existent-file")

		if gotErr == nil || !strings.Contains(gotErr.Error(), "non-existent-file") {
			t.Errorf(`run() gotErr = %q, want <nil>`, gotErr)
		}

		gotStderr := stderr.String()
		if len(gotStderr) > 0 {
			t.Errorf(`run() gotStderr = %q, want empty string`, gotStderr)
		}

		gotStdout := stdout.String()
		if len(gotStdout) > 0 {
			t.Errorf("run() gotStdout = %q, want empty string", gotStdout)
		}
	})

	dagTests := []struct {
		name string
		g    string
	}{
		{
			name: "line: a b b c c d",
			g:    "a b b c c d",
		},
		{
			name: "line: a b c d b c",
			g:    "a b c d b c",
		},
		{
			name: "line: b c a b c d",
			g:    "b c a b c d",
		},
		{
			name: "line: b c c d a b",
			g:    "b c c d a b",
		},
		{
			name: "line: c d a b b c",
			g:    "c d a b b c",
		},
		{
			name: "line: c d b c a b",
			g:    "c d b c a b",
		},
		{
			name: "diamond: a b a c b d c d",
			g:    "a b a c b d c d",
		},
		{
			//    a     b      c
			//   / \   / \     |
			//  /   \ /   \    |
			// d     e     f   g
			//       |\   /
			//       | \ /
			//       h  i
			name: "dag: a d a e b e b f e h e i f i c g",
			g:    "a d a e b e b f e h e i f i c g",
		},
	}
	for _, tt := range dagTests {
		t.Run(fmt.Sprintf("stdin: %s", tt.name), func(t *testing.T) {
			stdin := strings.NewReader(tt.g)
			stdout := new(strings.Builder)
			stderr := new(strings.Builder)

			gotErr := run(stdin, stdout, stderr)

			if gotErr != nil {
				t.Errorf(`run() gotErr = %q, want <nil>`, gotErr)
			}

			gotStderr := stderr.String()
			if gotStderr != "" {
				t.Errorf(`run() gotStderr = %q, want ""`, gotStderr)
			}

			checkValidTopologicalOrdering(t, tt.g, stdout)
		})
	}

	cycleTests := []struct {
		name            string
		g               string
		wantStdoutAnyOf []string
		wantStderrAnyOf []string
	}{
		{
			name:            "stdin: cycle: a b b a",
			g:               "a b b a",
			wantStdoutAnyOf: abInAnyRotation(),
			wantStderrAnyOf: cycleABInAnyRotation(),
		},
		{
			name:            "stdin: cycle: b c c d d b",
			g:               "b c c d d b",
			wantStdoutAnyOf: bcdInAnyRotation(),
			wantStderrAnyOf: cycleBCDInAnyRotation(),
		},
		{
			name:            "stdin: two cycles: a b b a c d d c",
			g:               "a b b a c d d c",
			wantStdoutAnyOf: abcdInAnyOrder(),
			wantStderrAnyOf: cyclesABAndCDInAnyOrderAndRotation(),
		},
		{
			name:            "stdin: orphan node then cycle: d d a b b c c a",
			g:               "d d a b b c c a",
			wantStdoutAnyOf: abcdInAnyOrder(),
			wantStderrAnyOf: cycleABCInAnyRotation(),
		},
		{
			name:            "stdin: cycle then orphan node: a b b c c a d d",
			g:               "a b b c c a d d",
			wantStdoutAnyOf: abcdInAnyOrder(),
			wantStderrAnyOf: cycleABCInAnyRotation(),
		},
		{
			name:            "stdin: two connected cycles: a b b a a c c a",
			g:               "a b b a a c c a",
			wantStdoutAnyOf: abcInAnyOrder(),
			wantStderrAnyOf: cyclesABAndACInAnyOrderAndRotation(),
		},
		{
			name:            "stdin: cycle with duplicate edges: a b a b b a",
			g:               "a b a b b a",
			wantStdoutAnyOf: abInAnyRotation(),
			wantStderrAnyOf: cycleABInAnyRotation(),
		},
	}
	for _, tt := range cycleTests {
		t.Run(tt.name, func(t *testing.T) {
			stdin := strings.NewReader(tt.g)
			stdout := new(strings.Builder)
			stderr := new(strings.Builder)

			gotErr := run(stdin, stdout, stderr)

			if !errors.Is(gotErr, errNonFatal) {
				t.Errorf(`run() gotErr = %q, want %q`, gotErr, errNonFatal)
			}

			gotStderr := stderr.String()
			if !slices.Contains(tt.wantStderrAnyOf, gotStderr) {
				t.Errorf(
					"run() gotStderr = %q, want any of %q",
					gotStderr,
					tt.wantStderrAnyOf)
			}

			gotStdout := stdout.String()
			if !slices.Contains(tt.wantStdoutAnyOf, gotStdout) {
				t.Errorf(
					"run() gotStdout = %q, want any of %q",
					gotStdout,
					tt.wantStdoutAnyOf)
			}
		})
	}

	t.Run("stdin: diamond and cycle: a b a c b d c d b e e f f b", func(t *testing.T) {
		stdin := strings.NewReader("a b a c b d c d b e e f f b")
		stdout := new(strings.Builder)
		stderr := new(strings.Builder)

		gotErr := run(stdin, stdout, stderr)

		if !errors.Is(gotErr, errNonFatal) {
			t.Errorf(`run() gotErr = %q, want %q`, gotErr, errNonFatal)
		}

		gotStderr := stderr.String()
		wantStderrAnyOf := cycleBEFInAnyRotation()
		if !slices.Contains(wantStderrAnyOf, gotStderr) {
			t.Errorf(
				"run() gotStderr = %q, want any of %q",
				gotStderr,
				wantStderrAnyOf)
		}

		gotStdout := stdout.String()
		fields := strings.Fields(gotStdout)
		if len(fields) != 6 {
			t.Errorf(
				`topological ordering invalid: want 6 elements, got %d elements: %q`,
				len(fields), gotStdout)
		}
		for _, e := range []edge{{"a", "b"}, {"a", "c"}, {"b", "d"}, {"c", "d"}} {
			if slices.Index(fields, e.source) >= slices.Index(fields, e.target) {
				t.Errorf(
					`gotStdout %q: topological ordering invalid: %v is not before %v`,
					gotStdout, e.source, e.target)
			}
		}
		if !slices.Contains(fields, "e") {
			t.Errorf(
				`gotStdout %q: topological ordering invalid: did not contain "e"`,
				gotStdout)
		}
		if !slices.Contains(fields, "f") {
			t.Errorf(
				`gotStdout %q: topological ordering invalid: did not contain "f"`,
				gotStdout)
		}
	})
}

var acyclicGraph = func() string {
	var result strings.Builder
	rnd := rand.New(rand.NewSource(1))
	n := 10_000
	for range 100 * n {
		x := rnd.Intn(n + 1)
		y := rnd.Intn(n + 1)
		_, _ = fmt.Fprintln(&result, min(x, y), max(x, y))
	}
	return result.String()
}()

var cyclicGraph = func() string {
	var result strings.Builder
	rnd := rand.New(rand.NewSource(1))
	n := 200
	for range 100 * n {
		x := rnd.Intn(n + 1)
		y := rnd.Intn(n + 1)
		_, _ = fmt.Fprintln(&result, x, y)
	}
	return result.String()
}()

func BenchmarkTsortAcyclicGraph(b *testing.B) {
	for b.Loop() {
		err := run(strings.NewReader(acyclicGraph), io.Discard, io.Discard)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

func BenchmarkTsortCyclicGraph(b *testing.B) {
	for b.Loop() {
		err := run(strings.NewReader(cyclicGraph), io.Discard, io.Discard)
		if err != nil && !errors.Is(err, errNonFatal) {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

func tempFile(t *testing.T, contents string) (file string) {
	dir := t.TempDir()
	n := filepath.Join(dir, "file")
	if err := os.WriteFile(n, []byte(contents), 0o666); err != nil {
		t.Fatalf("temp file not created: %v", err)
	}
	return n
}

func abInAnyRotation() []string {
	return []string{"a\nb\n", "b\na\n"}
}

func abcInAnyOrder() []string {
	return []string{
		"a\nb\nc\n",
		"a\nc\nb\n",
		"b\na\nc\n",
		"b\nc\na\n",
		"c\na\nb\n",
		"c\nb\na\n",
	}
}

func bcdInAnyRotation() []string {
	return []string{
		"b\nc\nd\n",
		"c\nd\nb\n",
		"d\nb\nc\n",
	}
}

func abcdInAnyOrder() []string {
	return []string{
		"a\nb\nc\nd\n",
		"a\nb\nd\nc\n",
		"a\nc\nb\nd\n",
		"a\nc\nd\nb\n",
		"a\nd\nb\nc\n",
		"a\nd\nc\nb\n",
		// ...
		"b\na\nc\nd\n",
		"b\na\nd\nc\n",
		"b\nc\na\nd\n",
		"b\nc\nd\na\n",
		"b\nd\na\nc\n",
		"b\nd\nc\na\n",
		// ...
		"c\na\nb\nd\n",
		"c\na\nd\nb\n",
		"c\nb\na\nd\n",
		"c\nb\nd\na\n",
		"c\nd\na\nb\n",
		"c\nd\nb\na\n",
		// ...
		"d\na\nb\nc\n",
		"d\na\nc\nb\n",
		"d\nb\na\nc\n",
		"d\nb\nc\na\n",
		"d\nc\na\nb\n",
		"d\nc\nb\na\n",
	}
}

func cycleABInAnyRotation() []string {
	return []string{
		"tsort: cycle in data\ntsort: a\ntsort: b\n",
		"tsort: cycle in data\ntsort: b\ntsort: a\n",
	}
}

func cycleABCInAnyRotation() []string {
	return []string{
		"tsort: cycle in data\ntsort: a\ntsort: b\ntsort: c\n",
		"tsort: cycle in data\ntsort: b\ntsort: c\ntsort: a\n",
		"tsort: cycle in data\ntsort: c\ntsort: a\ntsort: b\n",
	}
}

func cycleBCDInAnyRotation() []string {
	return []string{
		"tsort: cycle in data\ntsort: b\ntsort: c\ntsort: d\n",
		"tsort: cycle in data\ntsort: c\ntsort: d\ntsort: b\n",
		"tsort: cycle in data\ntsort: d\ntsort: b\ntsort: c\n",
	}
}

func cycleBEFInAnyRotation() []string {
	return []string{
		"tsort: cycle in data\ntsort: b\ntsort: e\ntsort: f\n",
		"tsort: cycle in data\ntsort: e\ntsort: f\ntsort: b\n",
		"tsort: cycle in data\ntsort: f\ntsort: b\ntsort: e\n",
	}
}

func cyclesABAndCDInAnyOrderAndRotation() []string {
	return []string{
		"tsort: cycle in data\ntsort: a\ntsort: b\ntsort: cycle in data\ntsort: c\ntsort: d\n",
		"tsort: cycle in data\ntsort: a\ntsort: b\ntsort: cycle in data\ntsort: d\ntsort: c\n",
		"tsort: cycle in data\ntsort: b\ntsort: a\ntsort: cycle in data\ntsort: c\ntsort: d\n",
		"tsort: cycle in data\ntsort: b\ntsort: a\ntsort: cycle in data\ntsort: d\ntsort: c\n",
		"tsort: cycle in data\ntsort: c\ntsort: d\ntsort: cycle in data\ntsort: a\ntsort: b\n",
		"tsort: cycle in data\ntsort: c\ntsort: d\ntsort: cycle in data\ntsort: b\ntsort: a\n",
		"tsort: cycle in data\ntsort: d\ntsort: c\ntsort: cycle in data\ntsort: a\ntsort: b\n",
		"tsort: cycle in data\ntsort: d\ntsort: c\ntsort: cycle in data\ntsort: b\ntsort: a\n",
	}
}

func cyclesABAndACInAnyOrderAndRotation() []string {
	return []string{
		"tsort: cycle in data\ntsort: a\ntsort: b\ntsort: cycle in data\ntsort: a\ntsort: c\n",
		"tsort: cycle in data\ntsort: a\ntsort: b\ntsort: cycle in data\ntsort: c\ntsort: a\n",
		"tsort: cycle in data\ntsort: b\ntsort: a\ntsort: cycle in data\ntsort: a\ntsort: c\n",
		"tsort: cycle in data\ntsort: b\ntsort: a\ntsort: cycle in data\ntsort: c\ntsort: a\n",
		"tsort: cycle in data\ntsort: a\ntsort: c\ntsort: cycle in data\ntsort: a\ntsort: b\n",
		"tsort: cycle in data\ntsort: a\ntsort: c\ntsort: cycle in data\ntsort: b\ntsort: a\n",
		"tsort: cycle in data\ntsort: c\ntsort: a\ntsort: cycle in data\ntsort: a\ntsort: b\n",
		"tsort: cycle in data\ntsort: c\ntsort: a\ntsort: cycle in data\ntsort: b\ntsort: a\n",
	}
}

type errorReturningReader struct {
	err error
}

func (e *errorReturningReader) Read(_ []byte) (n int, err error) {
	return 0, e.err
}

func checkValidTopologicalOrdering(
	t *testing.T,
	graph string,
	topologicalOrdering fmt.Stringer,
) {
	graphNodes := nodes(graph)
	topoNodes := strings.Fields(topologicalOrdering.String())

	if diff := orderInsensitiveDiff(graphNodes, topoNodes); diff != "" {
		t.Errorf(
			"topological ordering mismatch (-graphNodes +topoNodes):\n%s",
			diff)
	}

	if hasDuplicates(topoNodes) {
		t.Errorf("topological ordering has duplicates: %q", topoNodes)
	}

	for _, e := range edges(graph) {
		if slices.Index(topoNodes, e.source) >=
			slices.Index(topoNodes, e.target) {
			t.Errorf(
				"topological ordering invalid: %q is not before %q",
				e.source, e.target)
		}
	}
}

func nodes(graph string) []string {
	fields := strings.Fields(graph)
	s := makeSet()

	var result []string
	for _, value := range fields {
		if !s.has(value) {
			s.add(value)
			result = append(result, value)
		}
	}

	return result
}

type edge struct {
	source string
	target string
}

func edges(graph string) []edge {
	var result []edge
	nodes := strings.Fields(graph)
	for i := 0; i < len(nodes); i += 2 {
		result = append(result, edge{source: nodes[i], target: nodes[i+1]})
	}
	return result
}

func orderInsensitiveDiff(a []string, b []string) string {
	return cmp.Diff(
		a, b, cmpopts.SortSlices(func(x, y string) bool { return x < y }))
}

func hasDuplicates(values []string) bool {
	s := makeSet()
	for _, value := range values {
		if s.has(value) {
			return true
		}
		s.add(value)
	}
	return false
}
