// Copyright 2012-2024 the u-root Authors. All rights reserved
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
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
)

var errNonFatal = errors.New("non-fatal")
var errOddDataCount = errors.New("odd data count")

func run(
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
	args ...string,
) error {
	in := io.NopCloser(stdin)
	if len(args) >= 1 {
		var err error
		in, err = os.Open(args[0])
		if err != nil {
			return err
		}
	}
	defer in.Close()

	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanWords)

	g := newGraph()

	if err := parseInto(scanner, g); err != nil {
		return err
	}

	var err error
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

func parseInto(scanner *bufio.Scanner, g *graph) error {
	var odd bool

	next := func() (string, bool) {
		if !scanner.Scan() {
			return "", false
		}
		odd = !odd
		return scanner.Text(), true
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

	if err := scanner.Err(); err != nil {
		return err
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
	for {
		// Kahn's algorithm
		var result []string
		roots := rootsOf(g)
		nonRoots := nonRootsOf(g)
		for !roots.isEmpty() {
			next := roots.dequeue()
			result = append(result, next)
			for succ := range g.successors(next) {
				nonRoots.removeOne(succ)
				if !nonRoots.has(succ) {
					roots.enqueue(succ)
				}
			}
		}
		if nonRoots.isEmpty() {
			// No cycles left
			for _, value := range result {
				f(value)
			}
			break
		}

		// Break a cycle and try Kahn's algorithm again
		nonRoots.forEachUnique(func(next string) bool {
			cycle := cycleStartingAt(g, next)
			if len(cycle) == 0 {
				return true
			}

			g.removeEdge(cycle[len(cycle)-1], cycle[0])
			cycles(cycle)
			return false
		})
	}
}

func rootsOf(g *graph) queue {
	result := queue{}
	for node := range g.nodes {
		if g.inDegree(node) == 0 {
			result.enqueue(node)
		}
	}
	return result
}

func nonRootsOf(g *graph) multiset {
	result := newMultiset()
	for node := range g.nodes {
		if g.inDegree(node) > 0 {
			result.add(node, g.inDegree(node))
		}
	}
	return result
}

func cycleStartingAt(g *graph, node string) []string {
	stack := []string{node}
	inStack := make(set)
	inStack.add(node)
	popStack := func() string {
		var result string
		result, stack = stack[len(stack)-1], stack[:len(stack)-1]
		return result
	}

	var cycle []string
	var dfs func() bool
	dfs = func() bool {
		for succ := range g.successors(top(stack)) {
			if inStack.has(succ) {
				// cycle found
				cycle = append(cycle, popStack())
				for top(cycle) != succ {
					cycle = append(cycle, popStack())
				}
				slices.Reverse(cycle)
				return true
			}

			stack = append(stack, succ)
			inStack.add(succ)
			if dfs() {
				return true
			}
		}

		delete(inStack, popStack())
		return false
	}
	dfs()
	return cycle
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
