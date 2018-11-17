// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Print process information.
//
// Synopsis:
//     ps [-Aaex] [aux]
//
// Description:
//     ps reads the /proc filesystem and prints nice things about what it
//     finds.  /proc in linux has grown by a process of Evilution, so it's
//     messy.
//
// Options:
//     -A: select all processes. Identical to -e.
//     -e: select all processes. Identical to -A.
//     -x: BSD-Like style, with STAT Column and long CommandLine
//     -a: print all process except whose are session leaders or unlinked with terminal
//    aux: see every process on the system using BSD syntax
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
)

var (
	flags struct {
		all     bool
		nSidTty bool
		x       bool
		aux     bool
	}
	cmd  = "ps [-Aaex] [aux]"
	eUID = os.Geteuid()
)

func init() {
	defUsage := flag.Usage
	flag.Usage = func() {
		os.Args[0] = cmd
		defUsage()
	}
	flag.BoolVar(&flags.all, "A", false, "Select all processes.  Identical to -e.")
	flag.BoolVar(&flags.all, "e", false, "Select all processes.  Identical to -A.")
	flag.BoolVar(&flags.x, "x", false, "BSD-Like style, with STAT Column and long CommandLine")
	flag.BoolVar(&flags.nSidTty, "a", false, "Print all process except whose are session leaders or unlinked with terminal")

	if len(os.Args) > 1 {
		if isPermutation(os.Args[1], "aux") {
			flags.aux = true
			flags.all = true
		}
	}
}

// ProcessTable holds all the information needed for ps
type ProcessTable struct {
	table    []*Process
	mProc    *Process
	headers  []string // each column to print
	fields   []string // which fields of process to print, on order
	fstring  []string // formated strings
	maxwidth int      // DEPRECATED: reason -> remove terminal stuff
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

func isPermutation(check string, ref string) bool {
	if len(check) != len(ref) {
		return false
	}
	checkArray := strings.Split(check, "")
	refArray := strings.Split(ref, "")

	sort.Strings(checkArray)
	sort.Strings(refArray)

	for i := range check {
		if checkArray[i] != refArray[i] {
			return false
		}
	}
	return true
}

// For now, just read /proc/pid/stat and dump its brains.
// TODO: a nice clean way to turn /proc/pid/stat into a struct. (trying now)
// There has to be a way.
func ps(pT *ProcessTable, w io.Writer) error {
	if len(pT.table) == 0 {
		return nil
	}
	// sorting ProcessTable by PID
	sort.Sort(pT)

	switch {
	case flags.aux:
		pT.headers = []string{"PID", "PGRP", "SID", "TTY", "STAT", "TIME", "COMMAND"}
		pT.fields = []string{"Pid", "Pgrp", "Sid", "Ctty", "State", "Time", "Cmd"}
	case flags.x:
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
		case flags.nSidTty:
			// no session leaders and no unlinked terminals
			if p.Sid == p.Pid || p.Ctty == "?" {
				continue
			}

		case flags.x:
			// print only process with same eUID of caller
			if eUID != p.uid {
				continue
			}

		case flags.all:
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
	flag.Parse()
	pT := NewProcessTable()
	if err := pT.LoadTable(); err != nil {
		log.Fatal(err)
	}

	err := ps(pT, os.Stderr)
	if err != nil {
		log.Fatal(err)
	}
}
