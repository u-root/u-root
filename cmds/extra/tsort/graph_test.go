// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGraph(t *testing.T) {
	testSuccessors(t, fixtureGraph())
	testInDegree(t, fixtureGraph())
	testRemoveEdge(t, fixtureGraph())
}

func fixtureGraph() *graph {
	//    a     b      c   j
	//   / \   /|\     |
	//  /   \ / | \    |
	// d     e  |  f   g
	//       |\ | /
	//       | \|/
	//       h  i
	g := newGraph()
	g.putEdge("a", "d")
	g.putEdge("a", "e")
	g.putEdge("b", "e")
	g.putEdge("b", "f")
	g.putEdge("b", "i")
	g.putEdge("b", "i")
	g.putEdge("e", "h")
	g.putEdge("e", "i")
	g.putEdge("f", "i")
	g.putEdge("c", "g")
	g.addNode("j")
	return g
}

func testSuccessors(t *testing.T, g *graph) {
	if diff := cmp.Diff(g.successors("a"), setOf("d", "e")); diff != "" {
		t.Errorf(
			"set mismatch (-g.successors(\"a\") +expected):\n%s",
			diff)
	}
	if diff := cmp.Diff(g.successors("b"), setOf("e", "f", "i")); diff != "" {
		t.Errorf(
			"set mismatch (-g.successors(\"b\") +expected):\n%s",
			diff)
	}
	if diff := cmp.Diff(g.successors("e"), setOf("h", "i")); diff != "" {
		t.Errorf(
			"set mismatch (-g.successors(\"e\") +expected):\n%s",
			diff)
	}
	if diff := cmp.Diff(g.successors("f"), setOf("i")); diff != "" {
		t.Errorf(
			"set mismatch (-g.successors(\"f\") +expected):\n%s",
			diff)
	}
	if diff := cmp.Diff(g.successors("c"), setOf("g")); diff != "" {
		t.Errorf(
			"set mismatch (-g.successors(\"c\") +expected):\n%s",
			diff)
	}
	if len(g.successors("h")) > 0 {
		t.Errorf(`g.successors("h"): want empty, got %v`, g.successors("h"))
	}
	if len(g.successors("i")) > 0 {
		t.Errorf(`g.successors("i"): want empty, got %v`, g.successors("i"))
	}
	if len(g.successors("j")) > 0 {
		t.Errorf(`g.successors("j"): want empty, got %v`, g.successors("j"))
	}
	caughtPanic := catchPanic(func() { g.successors("absent") })
	if caughtPanic == nil ||
		!strings.Contains(caughtPanic.Error(), "node is not in graph") {
		t.Errorf(
			`g.successors("absent"): want panic with message "node is not in graph", got %#v`,
			caughtPanic)
	}
}

func testInDegree(t *testing.T, g *graph) {
	if g.inDegree("a") != 0 {
		t.Errorf(`g.inDegree("a"): want 0, got %d`, g.inDegree("a"))
	}
	if g.inDegree("d") != 1 {
		t.Errorf(`g.inDegree("d"): want 1, got %d`, g.inDegree("d"))
	}
	if g.inDegree("e") != 2 {
		t.Errorf(`g.inDegree("e"): want 2, got %d`, g.inDegree("e"))
	}
	if g.inDegree("i") != 3 {
		t.Errorf(`g.inDegree("i"): want 3, got %d`, g.inDegree("e"))
	}
	if g.inDegree("absent-node") != 0 {
		t.Errorf(
			`g.inDegree("absent-node"): want 0, got %d`,
			g.inDegree("absent-node"))
	}
}

func testRemoveEdge(t *testing.T, g *graph) {
	caughtPanic := catchPanic(func() { g.removeEdge("absent-source-node", "a") })
	if caughtPanic == nil ||
		!strings.Contains(caughtPanic.Error(), "source node is not in graph") {
		t.Errorf(
			`g.removeEdge("absent-source-node", "a"): want panic with message "source node is not in graph", got %#v`,
			caughtPanic)
	}
	testSuccessors(t, g) // test that there were no changes

	caughtPanic = catchPanic(func() { g.removeEdge("a", "absent-target-node") })
	if caughtPanic == nil ||
		!strings.Contains(caughtPanic.Error(), "target node is not in graph") {
		t.Errorf(
			`g.removeEdge("absent-target-node", "a"): want panic with message "target node is not in graph", got %#v`,
			caughtPanic)
	}
	testSuccessors(t, g) // test that there were no changes

	g.removeEdge("b", "e")
	if diff := cmp.Diff(g.successors("b"), setOf("f", "i")); diff != "" {
		t.Errorf(
			"set mismatch (-g.successors(\"b\") +expected):\n%s",
			diff)
	}
	if diff := cmp.Diff(g.successors("e"), setOf("h", "i")); diff != "" {
		t.Errorf(
			"set mismatch (-g.successors(\"e\") +expected):\n%s",
			diff)
	}
	if g.inDegree("e") != 1 {
		t.Errorf(
			`g.removeEdge("b", "e"): want g.inDegree("e") to be 1, got %d`,
			g.inDegree("e"))
	}
}
