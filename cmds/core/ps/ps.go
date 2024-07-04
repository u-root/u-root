// Copyright 2013-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9
// +build !plan9

// Print process information.
//
// Synopsis:
//
//	ps [-Aaex] [aux]
//
// Description:
//
//	ps reads the /proc filesystem and prints nice things about what it
//	finds.  /proc in linux has grown by a process of Evilution, so it's
//	messy.
//
// Options:
//
//	 -A: select all processes. Identical to -e.
//	 -e: select all processes. Identical to -A.
//	 -x: BSD-Like style, with STAT Column and long CommandLine
//	 -a: print all process except whose are session leaders or unlinked with terminal
//	aux: see every process on the system using BSD syntax
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"

	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

// Flags
var (
	all     bool
	every   bool
	x       bool
	nSidTty bool
	aux     = false
)

var (
	psUsage = "ps: ps [flags] [aux]"
	eUID    = os.Geteuid()
)

func usage() {
	defUsage := flag.Usage
	flag.Usage = func() {
		os.Args[0] = psUsage
		defUsage()
	}
	flag.Usage()
}

// ProcessTable holds all the information needed for ps
type ProcessTable struct {
	table   []*Process
	mProc   *Process
	headers []string // each column to print
	fields  []string // which fields of process to print, on order
	fstring []string // formated strings
}

// NewProcessTable creates an empty process table
func NewProcessTable() *ProcessTable {
	return &ProcessTable{}
}

// Len returns the number of processes in the ProcessTable.
func (pT ProcessTable) Len() int {
	return len(pT.table)
}

// to use on sort.Sort
func (pT ProcessTable) Less(i, j int) bool {
	return pT.table[i].Pidno < pT.table[j].Pidno
}

// to use on sort.Sort
func (pT ProcessTable) Swap(i, j int) {
	pT.table[i], pT.table[j] = pT.table[j], pT.table[i]
}

// Return the biggest value in a slice of ints.
func max(slice []int) int {
	max := slice[0]
	for _, value := range slice {
		if value > max {
			max = value
		}
	}
	return max
}

// MaxLength returns the longest string of a field of ProcessTable
func (pT ProcessTable) MaxLength(field string) int {
	slice := make([]int, 0)
	for _, p := range pT.table {
		slice = append(slice, len(p.Search(field)))
	}

	return max(slice)
}

// PrintHeader prints the header for ps, with correct spacing.
func (pT ProcessTable) PrintHeader(w io.Writer) {
	var row string
	for index, field := range pT.headers {
		formated := pT.fstring[index]
		row += fmt.Sprintf(formated, field)
	}

	fmt.Fprintf(w, "%v\n", row)
}

// PrintProcess prints information about one process.
func (pT ProcessTable) PrintProcess(index int, w io.Writer) {
	var row string
	p := pT.table[index]
	for index, f := range pT.fields {
		field := p.Search(f)
		formated := pT.fstring[index]
		row += fmt.Sprintf(formated, field)

	}

	fmt.Fprintf(w, "%v\n", row)
}

// PrepareString figures out how to lay out a process table print
func (pT *ProcessTable) PrepareString() {
	var (
		fstring  []string
		formated string
		PID      = pT.MaxLength("Pid")
		TTY      = pT.MaxLength("Ctty")
		STAT     = 4 | pT.MaxLength("State") // min : 4
		TIME     = pT.MaxLength("Time")
		CMD      = pT.MaxLength("Cmd")
	)
	for _, f := range pT.headers {
		switch f {
		case "PID":
			formated = fmt.Sprintf("%%%dv ", PID)
		case "TTY":
			formated = fmt.Sprintf("%%-%dv    ", TTY)
		case "STAT":
			formated = fmt.Sprintf("%%-%dv    ", STAT)
		case "TIME":
			formated = fmt.Sprintf("%%%dv ", TIME)
		case "CMD":
			formated = fmt.Sprintf("%%-%dv ", CMD)
		}
		fstring = append(fstring, formated)
	}

	pT.fstring = fstring
}

// For now, just read /proc/pid/stat and dump its brains.
func ps(w io.Writer, args ...string) error {
	// The original ps was designed before many flag conventions existed.
	// It had switches not needing a -. Try to emulate that.
	// It's pretty awful, however :-)
	for _, a := range args {
		switch a {
		case "aux":
			all, every, aux = true, true, true
		default:
			usage()
			return nil
		}
	}
	pT := NewProcessTable()
	if err := pT.LoadTable(); err != nil {
		return err
	}

	if pT.Len() == 0 {
		return nil
	}
	// sorting ProcessTable by PID
	sort.Sort(pT)

	switch {
	case aux:
		pT.headers = []string{"PID", "PGRP", "SID", "TTY", "STAT", "TIME", "COMMAND"}
		pT.fields = []string{"Pid", "Pgrp", "Sid", "Ctty", "State", "Time", "Cmd"}
	case x:
		pT.headers = []string{"PID", "TTY", "STAT", "TIME", "COMMAND"}
		pT.fields = []string{"Pid", "Ctty", "State", "Time", "Cmd"}
	default:
		pT.headers = []string{"PID", "TTY", "TIME", "CMD"}
		pT.fields = []string{"Pid", "Ctty", "Time", "Cmd"}
	}

	pT.PrepareString()
	pT.PrintHeader(w)
	for index, p := range pT.table {
		switch {
		case nSidTty:
			// no session leaders and no unlinked terminals
			if p.Sid == p.Pid || p.Ctty == "?" {
				continue
			}

		case x:
			// print only process with same eUID of caller
			if eUID != p.uid {
				continue
			}

		case all || every:
			// pass, print all

		default:
			// default for no flags only same session
			// and same uid process
			if p.Sid != pT.mProc.Sid || eUID != p.uid {
				continue
			}
		}

		pT.PrintProcess(index, w)
	}

	return nil
}

func main() {
	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	f.BoolVar(&all, "All", false, "Select all processes.  Identical to -e.")
	f.BoolVar(&all, "A", false, "Select all processes.  Identical to -e. (shorthand)")

	f.BoolVar(&every, "every", false, "Select all processes.  Identical to -A.")
	f.BoolVar(&every, "e", false, "Select all processes.  Identical to -A. (shorthand)")

	f.BoolVar(&x, "bsd", false, "BSD-Like style, with STAT Column and long CommandLine")
	f.BoolVar(&x, "x", false, "BSD-Like style, with STAT Column and long CommandLine (shorthand)")

	f.BoolVar(&nSidTty, "anSIDTTY", false, "Print all process except whose are session leaders or unlinked with terminal")
	f.BoolVar(&nSidTty, "a", false, "Print all process except whose are session leaders or unlinked with terminal (shorthand)")

	f.Parse(unixflag.OSArgsToGoArgs())
	if err := ps(os.Stdout, f.Args()...); err != nil {
		log.Fatal(err)
	}
}
