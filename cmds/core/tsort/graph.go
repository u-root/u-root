// Copyright 2012-2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"iter"
	"slices"
)

// nodeID is an efficient handle for a node that was added to a graph. nodeID
// values are contiguous, in strictly increasing order and within the range of
// [0..len(graph.nodeCount())].
//
// To get the original value associated with a nodeID, call
// graph.valueFor(nodeID).
type nodeID int32

func newGraph() *graph {
	return &graph{
		nodeToID:             make(map[string]nodeID),
		idToNode:             make([]string, 0),
		nodeIDToSuccessorIDs: make([][]nodeID, 0),
	}
}

type graph struct {
	nodeToID             map[string]nodeID
	idToNode             []string
	nodeIDToSuccessorIDs [][]nodeID
}

func (g *graph) addNode(node string) {
	_ = g.addNodeInternal(node)
}

func (g *graph) addNodeInternal(node string) nodeID {
	if id, ok := g.nodeToID[node]; ok {
		return id
	}

	id := nodeID(len(g.idToNode))
	g.nodeToID[node] = id
	g.idToNode = append(g.idToNode, node)

	g.nodeIDToSuccessorIDs = append(g.nodeIDToSuccessorIDs, nil)

	return id
}

func (g *graph) putEdge(source, target string) {
	sourceID := g.addNodeInternal(source)
	targetID := g.addNodeInternal(target)

	succs := g.nodeIDToSuccessorIDs[sourceID]
	if !slices.Contains(succs, targetID) {
		g.nodeIDToSuccessorIDs[sourceID] = append(succs, targetID)
	}
}

func (g *graph) valueFor(id nodeID) string {
	return g.idToNode[id]
}

func (g *graph) nodeCount() int {
	return len(g.nodeIDToSuccessorIDs)
}

func (g *graph) nodeIDs() iter.Seq[nodeID] {
	return func(yield func(nodeID) bool) {
		for id := range len(g.idToNode) {
			if !yield(nodeID(id)) {
				return
			}
		}
	}
}

func (g *graph) successorIDs(id nodeID) []nodeID {
	return g.nodeIDToSuccessorIDs[id]
}
