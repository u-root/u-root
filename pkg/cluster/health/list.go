// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// vitalsigns prints out vital statistics about a node as JSON.
// It currently uses only the jaypipes/ghw package, but
// more may come later.
package health

import (
	"os/exec"
	"strings"
)

var V = func(string, ...any) {}

// NewNodeList creates a NodeList with an empty List. The List
// is populated by running cmd and args.
func NewNodeList(cmd string, args ...string) *NodeList {
	return &NodeList{cmd: cmd, args: args}
}

// Run creates a list of nodes.
// The node names are assumed to be a list that can be split
// by strings.Fields.
func (n *NodeList) Run() error {
	l, err := exec.Command(n.cmd, n.args...).CombinedOutput()
	if err != nil {
		return err
	}
	// Do not sort. Most systems will return sort order,
	// and it does not matter if it is sorted or not.
	n.List = strings.Fields(string(l))
	return nil
}
