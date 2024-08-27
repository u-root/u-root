// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// clusterstats gathers information from a cluster.
// It uses three comands, specified as flags:
// -list: list the nodes from which to gather statistics
// -run: schedule a command to run
// -nodestat: command to run to get the data
//
//	This command *must* return the data as JSON.
//
// The JSON from the nodestats command is un-marshal'ed into
// a slice of structs, then marshal'ed back into JSON. Doing
// so provides an easy check that the JSON is valid.
//
// There is an easy test for basic operation, using only local resources:
// clusterstats -list 'echo localhost' -run ssh  -node nodestats
// This will create a list with just localhost, use ssh to run on that host
// and run the nodestats command.
// Slightly more extreme
// clusterstats -list 'echo nodestats' -run 'bash -c'  -node "-h"
// In this case, the list is the command, because run is a bash -c.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/u-root/u-root/pkg/cluster/health"
)

type cmd struct {
	c      exec.Cmd
	stdout io.Writer
	stderr io.Writer
	list   string
	sched  string
	node   string
	v      func(string, ...any)
}

func new() *cmd {
}

func (c *cmd) run() error {
	lc := strings.Fields(c.list)
	if len(lc) == 0 {
		log.Fatal("list command must have at least a ...command")
	}

	n := health.NewNodeList(lc[0], lc[1:]...)
	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
	rc := strings.Fields(c.sched)
	if len(rc) == 0 {
		log.Fatalf("to sched on %v hosts, sched command must have at least a ...command", n.List)
	}
	g := n.NewGather(rc[0], rc[1:]...)
	nc := strings.Fields(c.node)
	if len(nc) == 0 {
		log.Fatalf("to run gather %v on %v, node command must have at least a ...command", g, g)
	}
	all, err := g.Run(nc[0], nc[1:]...)
	if err != nil {
		log.Fatal(err)
	}
	j, err := json.MarshalIndent(all, "", "\t")
	if err != nil {
		log.Fatalf("marshaling all stats: %v", err)
	}
	fmt.Printf("%s\n", string(j))
	return nil
}

func main() {
	c := &cmd{stdout: os.Stdout, stderr: os.Stderr, v: func(string, ...any) {}}
	var verbose = flag.Bool("v", false, "print stuff")
	flag.StringVar(&c.list, "list", `sinfo -N -h -o %N`, "command to list nodes")
	flag.StringVar(&c.sched, "run", `srun -t 5 -i none -w`, "template of nodes to run on")
	flag.StringVar(&c.node, "node", "node", "node command to get stats")

	flag.Parse()
	if *verbose {
		c.v = log.Printf
		health.V = c.v
	}
	if err := c.run(); err != nil {
		log.Fatal(err)
	}

}
