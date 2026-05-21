// Copyright 2012-2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand/v2"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"testing/iotest"
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
			name:       "stdin: duplicate edge: a b a b",
			stdin:      strings.NewReader("a b a b"),
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
			stdin:   iotest.ErrReader(errDiskCrashed),
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
				t.Errorf(`gotErr = %q, want %q`, gotErr, tt.wantErr)
			}

			gotStderr := stderr.String()
			if gotStderr != tt.wantStderr {
				t.Errorf(
					"gotStderr = %q, want %q",
					gotStderr,
					tt.wantStderr)
			}

			gotStdout := stdout.String()
			if gotStdout != tt.wantStdout {
				t.Errorf(
					"gotStdout = %q, want %q",
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
			t.Errorf(`gotErr = %q, want <nil>`, gotErr)
		}

		gotStderr := stderr.String()
		if len(gotStderr) > 0 {
			t.Errorf(`gotStderr = %q, want empty string`, gotStderr)
		}

		gotStdout := stdout.String()
		if len(gotStdout) > 0 {
			t.Errorf("gotStdout = %q, want empty string", gotStdout)
		}
	})

	directedAcyclicGraphTests := []struct {
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
			//    a
			//   / \
			//  b   c
			//   \ /
			//    d
			// ...where edges are pointing downwards
			name: "diamond: a b a c b d c d",
			g:    "a b a c b d c d",
		},
		{
			//    a     b      c  j
			//   / \   / \     |
			//  /   \ /   \    |
			// d     e     f   g
			//       |\   /
			//       | \ /
			//       h  i
			// ...where edges are pointing downwards
			name: "directed acyclic graph: a d a e b e b f e h e i f i c g j j",
			g:    "a d a e b e b f e h e i f i c g j j",
		},
	}
	for _, tt := range directedAcyclicGraphTests {
		t.Run(fmt.Sprintf("stdin: %s", tt.name), func(t *testing.T) {
			stdin := strings.NewReader(tt.g)
			stdout := new(strings.Builder)
			stderr := new(strings.Builder)

			gotErr := run(stdin, stdout, stderr)

			if gotErr != nil {
				t.Errorf(`gotErr = %q, want <nil>`, gotErr)
			}

			if gotStderr := stderr.String(); gotStderr != "" {
				t.Errorf(`gotStderr = %q, want ""`, gotStderr)
			}

			checkValidTopologicalOrdering(t, tt.g, stdout)
		})
	}

	// When cycles are detected, we make no guarantees about which order the nodes are returned,
	// and we do not guarantee that every cycle is reported. This allows for more performance
	// optimizations.
	cycleTests := []struct {
		name            string
		g               string
		wantStdoutAnyOf []string
		wantStderrAnyOf []string
	}{
		{
			name:            "stdin: cycle: a b b a",
			g:               "a b b a",
			wantStdoutAnyOf: abInAnyOrder(),
			wantStderrAnyOf: cycleABInAnyOrder(),
		},
		{
			name:            "stdin: cycle: b c c d d b",
			g:               "b c c d d b",
			wantStdoutAnyOf: bcdInAnyOrder(),
			wantStderrAnyOf: cycleBCDInAnyRotation(),
		},
		{
			name:            "stdin: two cycles: a b b a c d d c",
			g:               "a b b a c d d c",
			wantStdoutAnyOf: abcdInAnyOrder(),
			wantStderrAnyOf: cyclesABOrCDInAnyOrderAndRotation(),
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
			wantStderrAnyOf: cyclesABOrACInAnyOrderAndRotation(),
		},
		{
			name:            "stdin: cycle with duplicate edges: a b a b b a",
			g:               "a b a b b a",
			wantStdoutAnyOf: abInAnyOrder(),
			wantStderrAnyOf: cycleABInAnyOrder(),
		},
	}
	for _, tt := range cycleTests {
		t.Run(tt.name, func(t *testing.T) {
			stdin := strings.NewReader(tt.g)
			stdout := new(strings.Builder)
			stderr := new(strings.Builder)

			gotErr := run(stdin, stdout, stderr)

			if !errors.Is(gotErr, errNonFatal) {
				t.Errorf(`gotErr = %q, want %q`, gotErr, errNonFatal)
			}

			gotStderr := stderr.String()
			if !slices.Contains(tt.wantStderrAnyOf, gotStderr) {
				t.Errorf(
					"gotStderr = %q, want any of %q",
					gotStderr,
					tt.wantStderrAnyOf)
			}

			gotStdout := stdout.String()
			if !slices.Contains(tt.wantStdoutAnyOf, gotStdout) {
				t.Errorf(
					"gotStdout = %q, want any of %q",
					gotStdout,
					tt.wantStdoutAnyOf)
			}
		})
	}

	directedGraphWithCycleTests := []struct {
		name                     string
		g                        string
		wantStdoutToRespectEdges []edge
		wantStderrAnyOf          []string
	}{
		{
			//    a
			//   / \
			//  b   c
			//  |\  |
			//  | \ |
			//  |   d
			//  e
			//  |
			//  f--->b (cycle back)
			// ...where vertical edges are pointing downwards
			name:                     "stdin: diamond and cycle: a b a c b d c d b e e f f b",
			g:                        "a b a c b d c d b e e f f b",
			wantStdoutToRespectEdges: []edge{{"a", "b"}, {"a", "c"}, {"b", "d"}, {"c", "d"}},
			wantStderrAnyOf:          cycleBEFInAnyRotation(),
		},
		{
			//    a
			//   / \
			//  b   f
			//  |
			//  c--->a (cycle back)
			//  |\
			//  d e
			// ...where vertical edges are pointing downwards
			name:                     "stdin: directed graph with cycle: a b b c c a c d c e a f",
			g:                        "a b b c c a c d c e a f",
			wantStdoutToRespectEdges: []edge{{"a", "f"}, {"c", "d"}, {"c", "e"}},
			wantStderrAnyOf:          cycleABCInAnyRotation(),
		},
	}
	for _, tt := range directedGraphWithCycleTests {
		t.Run(tt.name, func(t *testing.T) {
			stdin := strings.NewReader(tt.g)
			stdout := new(strings.Builder)
			stderr := new(strings.Builder)

			gotErr := run(stdin, stdout, stderr)

			if !errors.Is(gotErr, errNonFatal) {
				t.Errorf("gotErr = %q, want %q", gotErr, errNonFatal)
			}

			gotStderr := stderr.String()
			if !slices.Contains(tt.wantStderrAnyOf, gotStderr) {
				t.Errorf(
					"gotStderr = %q, want any of %q",
					gotStderr,
					tt.wantStderrAnyOf)
			}

			checkValidSoftTopologicalOrdering(
				t,
				tt.g,
				stdout,
				tt.wantStdoutToRespectEdges,
			)
		})
	}
}

var (
	rnd = rand.New(rand.NewPCG(1, 1))
)

func randomDirectedAcyclicGraph(nodeCount uint16, edgeCountRatio float64) string {
	if edgeCountRatio < 0.0 || edgeCountRatio > 1.0 {
		panic(fmt.Sprintf(
			"edgeCountRatio %v must be between 0.0 and 1.0",
			edgeCountRatio,
		))
	}

	totalPossibleEdges := maxEdgesForDirectedAcyclicGraph(nodeCount)
	edgeCount := uint(math.Round(float64(totalPossibleEdges) * edgeCountRatio))

	// filled with `false` by default
	randomEdges := make([]bool, totalPossibleEdges)
	for i := range edgeCount {
		randomEdges[i] = true
	}
	rnd.Shuffle(len(randomEdges), func(i, j int) {
		randomEdges[i], randomEdges[j] = randomEdges[j], randomEdges[i]
	})

	result := new(strings.Builder)
	for i := range nodeCount {
		_, _ = fmt.Fprintln(result, i, i)
	}
	index := 0
	for i := uint16(0); i < nodeCount-1; i++ {
		for j := i + 1; j < nodeCount; j++ {
			if randomEdges[index] {
				_, _ = fmt.Fprintln(result, i, j)
			}
			index++
		}
	}
	return result.String()
}

// For any directed acyclic graph, the maximum number of edges is equal to (n * (n - 1) / 2),
// where n is the number of nodes in the graph.
func maxEdgesForDirectedAcyclicGraph(nodeCount uint16) uint {
	return uint(nodeCount) * (uint(nodeCount) - 1) / 2
}

func randomDirectedCyclicGraph(nodeCount uint) string {
	result := new(strings.Builder)
	// Produces a cyclic graph through a fixed RNG seed and sheer probability.
	for i := range nodeCount {
		_, _ = fmt.Fprintln(result, i, i)
	}
	for range 100 * nodeCount {
		x := rnd.UintN(nodeCount + 1)
		y := rnd.UintN(nodeCount + 1)
		_, _ = fmt.Fprintln(result, x, y)
	}
	return result.String()
}

func BenchmarkTsortAcyclicGraph(b *testing.B) {
	b.Skipf("Fix testutils before re-enabling this, so we can skip in a vm")
	benchmarkCases := []struct {
		name         string
		acyclicGraph string
	}{
		{
			name:         "small sparse directed acyclic graph",
			acyclicGraph: randomDirectedAcyclicGraph(10, 0.1),
		},
		{
			name:         "small half-total-edges directed acyclic graph",
			acyclicGraph: randomDirectedAcyclicGraph(10, 0.5),
		},
		{
			name:         "small edgeless directed acyclic graph",
			acyclicGraph: randomDirectedAcyclicGraph(10, 0.0),
		},
		{
			name:         "small tournament directed acyclic graph",
			acyclicGraph: randomDirectedAcyclicGraph(10, 1.0),
		},
		{
			name:         "medium sparse directed acyclic graph",
			acyclicGraph: randomDirectedAcyclicGraph(100, 0.1),
		},
		{
			name:         "medium half-total-edges directed acyclic graph",
			acyclicGraph: randomDirectedAcyclicGraph(100, 0.5),
		},
		{
			name:         "medium edgeless directed acyclic graph",
			acyclicGraph: randomDirectedAcyclicGraph(100, 0),
		},
		{
			name:         "medium tournament directed acyclic graph",
			acyclicGraph: randomDirectedAcyclicGraph(100, 1.0),
		},
		{
			name:         "large sparse directed acyclic graph",
			acyclicGraph: randomDirectedAcyclicGraph(1_000, 0.1),
		},
		{
			name:         "large half-total-edges directed acyclic graph",
			acyclicGraph: randomDirectedAcyclicGraph(1_000, 0.5),
		},
		{
			name:         "large edgeless directed acyclic graph",
			acyclicGraph: randomDirectedAcyclicGraph(1_000, 0.0),
		},
		{
			name:         "large tournament directed acyclic graph",
			acyclicGraph: randomDirectedAcyclicGraph(1_000, 1.0),
		},
	}
	for _, bc := range benchmarkCases {
		b.Run(bc.name, func(b *testing.B) {
			g := bc.acyclicGraph
			for b.Loop() {
				err := run(strings.NewReader(g), io.Discard, io.Discard)
				if err != nil {
					b.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

func BenchmarkTsortCyclicGraph(b *testing.B) {
	b.Skipf("Fix testutils before re-enabling this, so we can skip in a vm")
	benchmarkCases := []struct {
		name        string
		cyclicGraph string
	}{
		{
			name:        "small cyclic graph",
			cyclicGraph: randomDirectedCyclicGraph(10),
		},
		{
			name:        "medium cyclic graph",
			cyclicGraph: randomDirectedCyclicGraph(50),
		},
		{
			name:        "large cyclic graph",
			cyclicGraph: randomDirectedCyclicGraph(100),
		},
	}

	for _, bc := range benchmarkCases {
		b.Run(bc.name, func(b *testing.B) {
			g := bc.cyclicGraph
			for b.Loop() {
				err := run(strings.NewReader(g), io.Discard, io.Discard)
				if err != nil && !errors.Is(err, errNonFatal) {
					b.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

func BenchmarkStressTest(b *testing.B) {
	b.Skipf("Fix testutils before re-enabling this, so we can skip in a vm")

	// Stress test the implementation to make sure it can handle humongously
	// deep graphs without crashing.
	//
	// WARNING: Consumes ~2GB of memory whilst running.
	g := lineGraphFrom0To2000000()

	for b.Loop() {
		err := run(strings.NewReader(g), io.Discard, io.Discard)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

func lineGraphFrom0To2000000() string {
	result := new(strings.Builder)
	for i := range uint(2_000_000) {
		_, _ = fmt.Fprintln(result, i, i+1)
	}
	return result.String()
}

func tempFile(t *testing.T, contents string) (file string) {
	n := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(n, []byte(contents), 0o600); err != nil {
		t.Fatalf("temp file not created: %v", err)
	}
	return n
}

func abInAnyOrder() []string {
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

func bcdInAnyOrder() []string {
	return []string{
		"b\nc\nd\n",
		"b\nd\nc\n",
		"c\nb\nd\n",
		"c\nd\nb\n",
		"d\nb\nc\n",
		"d\nc\nb\n",
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

func cycleABInAnyOrder() []string {
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

func cyclesABOrCDInAnyOrderAndRotation() []string {
	return []string{
		"tsort: cycle in data\ntsort: a\ntsort: b\n",
		"tsort: cycle in data\ntsort: b\ntsort: a\n",
		"tsort: cycle in data\ntsort: c\ntsort: d\n",
		"tsort: cycle in data\ntsort: d\ntsort: c\n",
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

func cyclesABOrACInAnyOrderAndRotation() []string {
	return []string{
		"tsort: cycle in data\ntsort: a\ntsort: b\n",
		"tsort: cycle in data\ntsort: b\ntsort: a\n",
		"tsort: cycle in data\ntsort: a\ntsort: c\n",
		"tsort: cycle in data\ntsort: c\ntsort: a\n",
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

func checkValidTopologicalOrdering(
	t *testing.T,
	graph string,
	topologicalOrdering fmt.Stringer,
) {
	checkValidSoftTopologicalOrdering(t, graph, topologicalOrdering, edges(graph))
}

func checkValidSoftTopologicalOrdering(
	t *testing.T,
	graph string,
	topologicalOrdering fmt.Stringer,
	shouldRespectEdges []edge,
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

	positions := make(map[string]int, len(topoNodes))
	for i, node := range topoNodes {
		positions[node] = i
	}

	for _, e := range shouldRespectEdges {
		if positions[e.source] >= positions[e.target] {
			t.Errorf(
				"topological ordering invalid: %q is not before %q",
				e.source, e.target)
		}
	}
}

func nodes(graph string) []string {
	result := slices.Sorted(strings.FieldsSeq(graph))
	return slices.Compact(result)
}

type edge struct {
	source string
	target string
}

func edges(graph string) []edge {
	var result []edge
	ns := strings.Fields(graph)
	for i := 0; i < len(ns); i += 2 {
		if ns[i] == ns[i+1] {
			// skip self-loop edges like 'edge{"a", "a"}', because tsort
			// treats them as single nodes rather than proper edges.
			continue
		}
		result = append(result, edge{source: ns[i], target: ns[i+1]})
	}
	return result
}

func hasDuplicates(values []string) bool {
	s := make(map[string]struct{})
	for _, value := range values {
		if _, ok := s[value]; ok {
			return true
		}
		s[value] = struct{}{}
	}
	return false
}
