// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// vitalsigns prints out vital statistics about a node as JSON.
// It currently uses only the jaypipes/ghw package, but
// more may come later.
package health

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"golang.org/x/sync/errgroup"
)

// NewGather creates a Gather from a NodeList. The cmd and args
// are used to invoke a scheduler command (e.g. srun, sbatch) to
// invoke a command on each node in turn. The assumption is that the
// command pattern looks like
// cmd args... nodename command-for-node args-for-node
// It is possible that, over time, we will need something more flexible,
// but this works for most schedulers we know of.
func (n *NodeList) NewGather(cmd string, args ...string) *Gather {
	return &Gather{cmd: cmd, args: args, Nodes: n}
}

// Run gathers data from the cluster, running cmd and args on
// each node. How they are run on the node is determined by
// the cmd and args in the Gather struct. Almost every error that
// can occur will occur on a per-node basis; those errors are accumulated
// into the Stderr member of the Stat struct, not returned.
// The intent is that this will saved in storage, so including errors
// is important.
func (g *Gather) Run(cmd string, args ...string) ([]Stat, error) {
	eg := new(errgroup.Group)
	c := make(chan Stat, len(g.Nodes.List))
	V("gather: run node %v, %v via scheduler %v", cmd, args, g.cmd)
	for _, node := range g.Nodes.List {
		V("Run on %v", node)
		eg.Go(func() error {
			cmdargs := append(append(g.args, node, cmd), args...)
			// Note: we do not use errgroup timeout or anything,
			// since the command itself should have a timeout.
			V("run %v %v", g.cmd, cmdargs)
			stat, err := exec.Command(g.cmd, cmdargs...).CombinedOutput()
			if err != nil {
				V("err stat %v %v", string(stat), err)
				c <- Stat{Hostname: node, Err: fmt.Sprintf("%s %v", string(stat), err)}
				return nil
			}
			var s Stat
			if err := json.Unmarshal(stat, &s); err != nil {
				V("err %v", err)
				c <- Stat{Hostname: node, Err: err.Error()}
				return nil
			}
			V("push %v to chan", s)
			c <- s
			return nil
		})
	}

	eg.Wait()
	close(c)

	V("all returned")
	all := make([]Stat, 0, len(g.Nodes.List))
	for s := range c {
		all = append(all, s)
	}

	V("all done gather")
	return all, nil
}
