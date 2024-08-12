// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package health is used for quick measuring of cluster health.
// Health for a node is provided in a Stat struct, containing
// everything we want to know. A lot of the data is provided by
// the github.com/jaypipes/ghw.HostInfo struct.
//
// To measure health requires invoking three different commands.
// We have found it convenient to structure this via three separate
// steps.
//
// First is the command to list the nodes.
// Second is to create the struct with the command to invoke the scheduler to run a command on a node.
// Third is to run the command using the struct
//
// These three steps result in running two commands: a lister and a gatherer.
//
// The lister operation is based on the NodeList struct.
// The zero value of a NodeList is a command, arguments, and empty list.
// NewNodeList takes a command and arguments.
// Run runs the command and populates the list. The List is exported and can
// be filled other ways. The command and args can be as simple as
// 'echo', '1', '2', '3' if that is desired; the List will then be '1', '2', and '3'.
//
// Gather runs a scheduler command to gather Stat structs.
// NewGather takes the scheduler command and arguments needed to *schedule*
// a process on a node. Run takes the command and arguments to be *run
// on each node*. Run produces a JSON-encoded []Stat.
//
// Typical commands for NodeList: sinfo -N -h -o %N
// Simpler commands for NodeList: echo ...names
// Typical commands for Gather: srun -t 5 -i none -w
// Simpler commands for Gather: cat ...file
package health
