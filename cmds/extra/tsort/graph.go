// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

func newGraph() *graph {
	return &graph{
		nodeToData: make(map[string]*nodeData),
	}
}

type graph struct {
	nodeToData map[string]*nodeData
}

type nodeData struct {
	inDegree   int
	successors set
}

func (g *graph) addNode(node string) {
	if _, ok := g.nodeToData[node]; !ok {
		g.nodeToData[node] = &nodeData{
			inDegree:   0,
			successors: makeSet(),
		}
	}
}

func (g *graph) putEdge(source, target string) {
	g.addNode(source)
	g.addNode(target)

	successors := g.nodeToData[source].successors
	if !successors.has(target) {
		successors.add(target)
		g.nodeToData[target].inDegree++
	}
}

func (g *graph) successors(node string) set {
	n, ok := g.nodeToData[node]
	if !ok {
		panic("node is not in graph")
	}

	return n.successors
}

func (g *graph) removeEdge(source, target string) {
	if _, ok := g.nodeToData[source]; !ok {
		panic("source node is not in graph")
	}
	if _, ok := g.nodeToData[target]; !ok {
		panic("target node is not in graph")
	}

	g.nodeToData[source].successors.remove(target)
	g.nodeToData[target].inDegree--
}

func (g *graph) inDegree(node string) int {
	n, ok := g.nodeToData[node]
	if !ok {
		return 0
	}
	return n.inDegree
}
