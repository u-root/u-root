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
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

var (
	flags struct {
		all     bool
		nSidTty bool
		x       bool
		aux     bool
	}
	cmd     = "ps [-Aaex] [aux]"
	eUID    = os.Geteuid()
	mainPID = os.Getpid()
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

// main process table of ps
// used to make more flexible
type ProcessTable struct {
	table    []*Process
	headers  []string // each column to print
	fields   []string // which fields of process to print, on order
	fstring  []string // formated strings
	maxwidth int      // DEPRECATED: reason -> remove terminal stuff
}

// to use on sort.Sort
func (pT ProcessTable) Len() int {
	return len(pT.table)
}

// to use on sort.Sort
func (pT ProcessTable) Less(i, j int) bool {
	a, _ := strconv.Atoi(pT.table[i].Pid)
	b, _ := strconv.Atoi(pT.table[j].Pid)
	return a < b
}

// to use on sort.Sort
func (pT ProcessTable) Swap(i, j int) {
	pT.table[i], pT.table[j] = pT.table[j], pT.table[i]
}

// Gived a pid, search for a process
// Returns nil if not found
func (pT ProcessTable) GetProcess(pid int) (found *Process) {
	for _, p := range pT.table {
		if p.Pid == strconv.Itoa(pid) {
			found = p
			break
		}
	}
	return
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

// fetch the most long string of a field of ProcessTable
// example: biggest len of string Pid of processes
func (pT ProcessTable) MaxLenght(field string) int {
	slice := make([]int, 0)
	for _, p := range pT.table {
		slice = append(slice, len(p.Search(field)))
	}

	return max(slice)
}

// Defined the each header
// Print them pT.headers
func (pT ProcessTable) PrintHeader() {
	var row string
	for index, field := range pT.headers {
		formated := pT.fstring[index]
		row += fmt.Sprintf(formated, field)
	}

	fmt.Printf("%v\n", row)
}

// Print an single processing for defined fields
// by ith-position on table slice of ProcessTable
func (pT ProcessTable) PrintProcess(index int) {
	var row string
	p := pT.table[index]
	for index, f := range pT.fields {
		field := p.Search(f)
		formated := pT.fstring[index]
		row += fmt.Sprintf(formated, field)

	}

	fmt.Printf("%v\n", row)
}

func (pT *ProcessTable) PrepareString() {
	var (
		fstring  []string
		formated string
		PID      = pT.MaxLenght("Pid")
		TTY      = pT.MaxLenght("Ctty")
		STAT     = 4 | pT.MaxLenght("State") // min : 4
		TIME     = pT.MaxLenght("Time")
		CMD      = pT.MaxLenght("Cmd")
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
func ps(pT ProcessTable) error {
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

	mProc := pT.GetProcess(mainPID)

	pT.PrepareString()
	pT.PrintHeader()
	for index, p := range pT.table {
		uid, err := p.GetUid()
		if err != nil {
			// It is extremely common for a directory to disappear from
			// /proc when a process terminates, so ignore those errors.
			if os.IsNotExist(err) {
				continue
			}
			return err
		}

		switch {
		case flags.nSidTty:
			// no session leaders and no unlinked terminals
			if p.Sid == p.Pid || p.Ctty == "?" {
				continue
			}

		case flags.x:
			// print only process with same eUID of caller
			if eUID != uid {
				continue
			}

		case flags.all:
			// pass, print all

		default:
			// default for no flags only same session
			// and same uid process
			if p.Sid != mProc.Sid || eUID != uid {
				continue
			}
		}

		pT.PrintProcess(index)
	}

	return nil

}

func main() {
	flag.Parse()
	pT := ProcessTable{}
	if err := pT.LoadTable(); err != nil {
		log.Fatal(err)
	}

	err := ps(pT)
	if err != nil {
		log.Fatal(err)
	}
}
