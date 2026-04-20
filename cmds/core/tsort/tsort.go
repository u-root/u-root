// Copyright 2012-2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// tsort writes to standard output a totally ordered list of items consistent
// with a partial ordering of items contained in the input. The standard input
// will be used if no file is specified.
//
// The input is a sequence of pairs of items, separated by <blank> characters.
// Pairs of different items (e.g., "a b") indicate ordering. Pairs of identical
// items (e.g., "c c") indicate presence, but not ordering.
//
// Synopsis:
//
//	tsort [FILE]
//
// Example:
//
//	tsort <<EOF
//	a b c c d e
//	g g
//	f g e f
//	h h
//	EOF
//
// produces an output like:
//
//	a
//	b
//	c
//	d
//	e
//	f
//	g
//	h
//
// which is one valid total ordering, but this is not guaranteed, it could
// equally be:
//
//	h
//	a
//	c
//	d
//	b
//	e
//	f
//	g
//
// or any other ordering where the following holds true:
//
//	- a is before b
//	- d is before e
//	- f is before g
//	- e is before f
//	- c is anywhere
//	- h is anywhere

package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"
)

var (
	errNonFatal     = errors.New("non-fatal")
	errOddDataCount = errors.New("odd data count")
)

func run(
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
	args ...string,
) error {
	var err error
	in := io.NopCloser(stdin)
	if len(args) >= 1 {
		in, err = os.Open(args[0])
		if err != nil {
			return err
		}
	}
	defer in.Close()

	var buf strings.Builder
	if _, err = io.Copy(&buf, in); err != nil {
		return err
	}

	g := newGraph()
	if err = parseInto(buf.String(), g); err != nil {
		return err
	}

	topologicalOrdering(
		g,
		func(node string) {
			fmt.Fprintf(stdout, "%v\n", node)
		},
		func(cycle []string) {
			fmt.Fprintf(stderr, "tsort: %v\n", "cycle in data")
			for _, node := range cycle {
				fmt.Fprintf(stderr, "tsort: %v\n", node)
			}
			err = errNonFatal
		})
	return err
}

func parseInto(buf string, g *graph) error {
	fields := strings.Fields(buf)
	var i int
	var odd bool

	next := func() (string, bool) {
		if i == len(fields) {
			return "", false
		}
		odd = !odd
		result := fields[i]
		i++
		return result, true
	}

	for {
		a, ok := next()
		if !ok {
			break
		}

		b, ok := next()
		if !ok {
			break
		}

		if a == b {
			g.addNode(a)
		} else {
			g.putEdge(a, b)
		}
	}

	if odd {
		return errOddDataCount
	}

	return nil
}

func topologicalOrdering(
	g *graph,
	f func(node string),
	cycles func(cycle []string),
) {
	// Variant of Kahn's algorithm that returns an ordering even for graphs
	// with cycles.
	roots := rootsOf(g)
	for g.nodeCount() != 0 {
		var next string
		next, roots = dequeueBreakingCycleIfNeeded(roots, g, cycles)
		f(next)
		for succ := range g.successors(next) {
			g.removeEdge(next, succ)
			if g.inDegree(succ) == 0 {
				roots.enqueue(succ)
			}
		}
		g.removeNode(next)
	}
}

func rootsOf(g *graph) queue {
	result := queue{}
	for node := range g.nodeToData {
		if g.inDegree(node) == 0 {
			result.enqueue(node)
		}
	}
	return result
}

func dequeueBreakingCycleIfNeeded(
	roots queue,
	g *graph,
	cycles func(cycle []string),
) (string, queue) {
	for {
		if next, ok := roots.dequeue(); ok {
			return next, roots
		}

		// The graph still has at least one node left, but there are no more
		// roots in the queue, so at least one cycle is present.
		//
		// Breaking a cycle has a chance of producing a new root in the graph,
		// so this loop repeatedly finds and breaks cycles until a new root
		// is found, which is immediately enqueued. This allows the greater
		// topological ordering algorithm to continue.
		cycle := findCycle(g)
		start, end := cycle[0], cycle[len(cycle)-1]
		g.removeEdge(end, start)
		cycles(cycle)
		if g.inDegree(start) == 0 {
			roots.enqueue(start)
		}
	}
}

func findCycle(g *graph) []string {
	var stack []string
	visited := makeSet()

	popStack := func() string {
		var result string
		result, stack = stack[len(stack)-1], stack[:len(stack)-1]
		return result
	}

	var cycle []string
	var dfs func() bool
	dfs = func() bool {
		for succ := range g.successors(top(stack)) {
			if visited.has(succ) {
				// cycle found
				cycle = append(cycle, popStack())
				for top(cycle) != succ {
					cycle = append(cycle, popStack())
				}
				slices.Reverse(cycle)
				return true
			}

			stack = append(stack, succ)
			visited.add(succ)
			if dfs() {
				return true
			}
		}

		visited.remove(popStack())
		return false
	}

	for node := range g.nodes() {
		if !visited.has(node) {
			stack = []string{node}
			visited.add(node)
			if dfs() {
				return cycle
			}
		}
	}

	panic("unreachable")
}

func top(s []string) string {
	return s[len(s)-1]
}

func main() {
	err := run(os.Stdin, os.Stdout, os.Stderr, os.Args[1:]...)
	if errors.Is(err, errNonFatal) {
		// All non-fatal warnings have been printed already, so just exit.
		os.Exit(1)
	}
	if err != nil {
		log.Fatalf("tsort: %v", err)
	}
}
