// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

func newGraph() *graph {
	return &graph{
		nodes:            make(set),
		nodeToInDegree:   newMultiset(),
		nodeToSuccessors: make(map[string]set),
	}
}

type graph struct {
	nodes            set
	nodeToInDegree   multiset
	nodeToSuccessors map[string]set
}

func (g *graph) addNode(node string) {
	g.nodes.add(node)
}

func (g *graph) putEdge(source, target string) {
	g.addNode(source)
	g.addNode(target)

	successors, ok := g.nodeToSuccessors[source]
	if !ok {
		successors = make(set)
		g.nodeToSuccessors[source] = successors
	}
	if !successors.has(target) {
		successors.add(target)
		g.nodeToInDegree.add(target, 1)
	}
}

func (g *graph) successors(node string) set {
	if !g.nodes.has(node) {
		panic("node is not in graph")
	}

	return g.nodeToSuccessors[node]
}

func (g *graph) removeEdge(source, target string) {
	if !g.nodes.has(source) {
		panic("source node is not in graph")
	}
	if !g.nodes.has(target) {
		panic("target node is not in graph")
	}

	successors := g.nodeToSuccessors[source]

	delete(successors, target)
	if len(successors) == 0 {
		delete(g.nodeToSuccessors, source)
	}

	g.nodeToInDegree.removeOne(target)
}

func (g *graph) inDegree(node string) int {
	return g.nodeToInDegree.count(node)
}
