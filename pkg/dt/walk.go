// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dt

import (
	"fmt"
)

// NodeWalk is used to contain state for walking
// the FDT, such as an error. A Walk with a non-nil
// error value can not proceed. Many walks will start
// with a root, but it is possible to Walk to a node,
// make a copy of the Walk, and in that way do multiple
// Walks from that one node. This is very similar to how
// 9p clients walk 9p servers.
type NodeWalk struct {
	n   *Node
	err error
}

// AsString returns the NodeWalk Name and error as a string.
func (nq *NodeWalk) AsString() (string, error) {
	if nq.err != nil {
		return "", nq.err
	}
	return nq.n.Name, nil
}

// ListChildNodes returns a string array with the Names of each child Node
func (nq *NodeWalk) ListChildNodes() ([]string, error) {
	if nq.err != nil {
		return nil, nq.err
	}

	cs := make([]string, len(nq.n.Children))
	for i := range nq.n.Children {
		cs[i] = nq.n.Children[i].Name
	}
	return cs, nil
}

// Root returns the Root node from an FDT to start the walk.
func (fdt *FDT) Root() *NodeWalk {
	return &NodeWalk{n: fdt.RootNode}
}

// Walk walks from a node to a named Node, returning a NodeWalk.
func (nq *NodeWalk) Walk(name string) *NodeWalk {
	if nq.err != nil {
		return nq
	}
	for _, n := range nq.n.Children {
		if n.Name == name {
			return &NodeWalk{n: n}
		}
	}
	return &NodeWalk{err: fmt.Errorf("cannot find node name %q", name)}
}

// Find returns a Node given a matching function starting at the current
// NodeWalk.
func (nq *NodeWalk) Find(f func(*Node) bool) (*Node, error) {
	if nq.err != nil {
		return nil, nq.err
	}
	if matching, ok := nq.n.Find(f); ok {
		return matching, nil
	}
	return nil, fmt.Errorf("cannot find node with matching pattern")
}

// FindAll returns all Nodes given a matching function starting at the current
// NodeWalk.
func (nq *NodeWalk) FindAll(f func(*Node) bool) ([]*Node, error) {
	if nq.err != nil {
		return nil, nq.err
	}
	if matching, ok := nq.n.FindAll(f); ok {
		return matching, nil
	}
	return nil, fmt.Errorf("cannot find nodes with matching pattern")
}

// Property walks from a Node to a Property of that Node, returning a PropertyWalk.
func (nq *NodeWalk) Property(name string) *PropertyWalk {
	if nq.err != nil {
		return &PropertyWalk{err: nq.err}
	}
	for _, p := range nq.n.Properties {
		if p.Name == name {
			return &PropertyWalk{p: &p}
		}
	}
	return &PropertyWalk{err: fmt.Errorf("cannot find property name %q", name)}
}

// PropertyWalk contains the state from a Walk
// to a Property.
type PropertyWalk struct {
	p   *Property
	err error
}

// AsU64 returns the PropertyWalk value as a uint64.
func (pq *PropertyWalk) AsU64() (uint64, error) {
	if pq.err != nil {
		return 0, pq.err
	}
	return pq.p.AsU64()
}

// AsString returns the PropertyWalk value as a string.
func (pq *PropertyWalk) AsString() (string, error) {
	if pq.err != nil {
		return "", pq.err
	}
	return pq.p.AsString()
}

// AsBytes returns the PropertyWalk value as a []byte.
func (pq *PropertyWalk) AsBytes() ([]byte, error) {
	if pq.err != nil {
		return nil, pq.err
	}
	return pq.p.Value, nil
}
