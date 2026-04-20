// Copyright 2012-2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"slices"
	"strings"
	"testing"
)

func TestGraph(t *testing.T) {
	testAllNodes(t, fixtureGraph())
	testNodeCount(t, fixtureGraph())
	testSuccessors(t, fixtureGraph())
	testInDegree(t, fixtureGraph())
	testRemoveNode(t, fixtureGraph())
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
	// ...where edges are pointed downwards
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

func testAllNodes(t *testing.T, g *graph) {
	got := slices.Collect(g.nodes())
	want := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	if diff := orderInsensitiveDiff(got, want); diff != "" {
		t.Fatalf(
			"allNodes mismatch (-actual +expected):\n%s",
			diff)
	}
}

func testNodeCount(t *testing.T, g *graph) {
	if got, want := g.nodeCount(), 10; got != want {
		t.Errorf("g.nodeCount(): got %d, want %d", got, want)
	}

	g.addNode("k")
	if got, want := g.nodeCount(), 11; got != want {
		t.Errorf("g.nodeCount(): got %d, want %d", got, want)
	}
}

func testSuccessors(t *testing.T, g *graph) {
	if diff := orderInsensitiveDiff(slices.Collect(g.successors("a")), []string{"d", "e"}); diff != "" {
		t.Errorf(
			"set mismatch (-g.successors(\"a\") +expected):\n%s",
			diff)
	}
	if diff := orderInsensitiveDiff(slices.Collect(g.successors("b")), []string{"e", "f", "i"}); diff != "" {
		t.Errorf(
			"set mismatch (-g.successors(\"b\") +expected):\n%s",
			diff)
	}
	if diff := orderInsensitiveDiff(slices.Collect(g.successors("e")), []string{"h", "i"}); diff != "" {
		t.Errorf(
			"set mismatch (-g.successors(\"e\") +expected):\n%s",
			diff)
	}
	if diff := orderInsensitiveDiff(slices.Collect(g.successors("f")), []string{"i"}); diff != "" {
		t.Errorf(
			"set mismatch (-g.successors(\"f\") +expected):\n%s",
			diff)
	}
	if diff := orderInsensitiveDiff(slices.Collect(g.successors("c")), []string{"g"}); diff != "" {
		t.Errorf(
			"set mismatch (-g.successors(\"c\") +expected):\n%s",
			diff)
	}
	if got := slices.Collect(g.successors("h")); len(got) > 0 {
		t.Errorf(`g.successors("h"): want empty, got %v`, got)
	}
	if got := slices.Collect(g.successors("i")); len(got) > 0 {
		t.Errorf(`g.successors("i"): want empty, got %v`, got)
	}
	if got := slices.Collect(g.successors("j")); len(got) > 0 {
		t.Errorf(`g.successors("j"): want empty, got %v`, got)
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

func testRemoveNode(t *testing.T, g *graph) {
	caughtPanic := catchPanic(func() { g.removeNode("absent-node") })
	if caughtPanic == nil ||
		!strings.Contains(caughtPanic.Error(), "node is not in graph") {
		t.Errorf(
			`g.removeNode("absent-node"): want panic with message "node is not in graph", got %#v`,
			caughtPanic)
	}

	g.removeNode("j")
	if diff := orderInsensitiveDiff(
		slices.Collect(g.nodes()),
		[]string{"a", "b", "c", "d", "e", "f", "g", "h", "i"},
	); diff != "" {
		t.Fatalf("g.removeNode(\"j\"): nodes mismatch (-got +want):\n%s", diff)
	}

	g.removeNode("c")
	if diff := orderInsensitiveDiff(
		slices.Collect(g.nodes()),
		[]string{"a", "b", "d", "e", "f", "g", "h", "i"},
	); diff != "" {
		t.Errorf("g.removeNode(\"c\"): nodes mismatch (-got +want):\n%s", diff)
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
	if diff := orderInsensitiveDiff(slices.Collect(g.successors("b")), []string{"f", "i"}); diff != "" {
		t.Errorf(
			"set mismatch (-g.successors(\"b\") +expected):\n%s",
			diff)
	}
	if diff := orderInsensitiveDiff(slices.Collect(g.successors("e")), []string{"h", "i"}); diff != "" {
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

func catchPanic(f func()) (caughtPanic error) {
	defer func() {
		if e := recover(); e != nil {
			caughtPanic = fmt.Errorf("%v", e)
		}
	}()

	f()
	return
}
