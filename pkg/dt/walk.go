// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dt

import "fmt"

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
	return pq.p.String(), nil
}

// AsBytes returns the PropertyWalk value as a []byte.
func (pq *PropertyWalk) AsBytes() ([]byte, error) {
	if pq.err != nil {
		return nil, pq.err
	}
	return pq.p.Value, nil
}
