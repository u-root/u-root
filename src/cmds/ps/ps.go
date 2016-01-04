// Copyright 2013 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// ps reads the /proc and prints out nice things about what it finds.
// /proc in linux has grown by a process of Evilution, so it's messy.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

var (
	flags struct {
		all    bool
		notSid bool
		x      bool
	}
	cmd     = "ps [-Aaex]"
	eUID    = os.Geteuid()
	mainPID = os.Getpid()
)

func usage() {
	fmt.Fprintln(os.Stderr, "Usage:", cmd)
	flag.PrintDefaults()
	os.Exit(1)
}

func init() {
	flag.BoolVar(&flags.all, "A", false, "Select all processes.  Identical to -e.")
	flag.BoolVar(&flags.all, "e", false, "Select all processes.  Identical to -A.")
	flag.BoolVar(&flags.x, "x", false, "BSD-Like style, with STAT Column and long CommandLine")
	flag.BoolVar(&flags.notSid, "a", false, "Print all process except whose are session leaders or unlinked with terminal")
	flag.Parse()
	flag.Usage = usage
}

// main process table of ps
// used to make more flexible
type ProcessTable struct {
	table   []Process
	headers []string // each column to print
	fields  []string // which fields of process to print, on order
	fstring []string // formated strings
}

// to use on sort.Sort
func (pT ProcessTable) Len() int {
	return len(pT.table)
}

// to use on sort.Sort
func (pT ProcessTable) Less(i, j int) bool {
	a, _ := strconv.Atoi(fieldString(pT.table[i].Pid))
	b, _ := strconv.Atoi(fieldString(pT.table[j].Pid))
	return a < b
}

// to use on sort.Sort
func (pT ProcessTable) Swap(i, j int) {
	pT.table[i], pT.table[j] = pT.table[j], pT.table[i]
}

// Gived a pid, search for a process
// Returns nil if not found
func (pT ProcessTable) GetProcess(pid int) (found Process) {
	for _, p := range pT.table {
		if fieldString(p.Pid) == strconv.Itoa(pid) {
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
		slice = append(slice, fieldLen(p.GetField(field)))
	}

	return max(slice)
}

// Defined the each header
// Print them pT.headers
func (pT ProcessTable) PrintHeader() {
	for index, field := range pT.headers {
		formated := pT.fstring[index]
		fmt.Printf(formated, field)
	}
	fmt.Printf("\n")
}

// Print an single processing for defined fields
// by ith-position on table slice of ProcessTable
func (pT ProcessTable) PrintProcess(index int) {
	var row string
	p := pT.table[index]
	for index, f := range pT.fields {
		field := p.GetField(f)
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
		STAT     = pT.MaxLenght("State")
		TIME     = pT.MaxLenght("Time")
		CMD      = pT.MaxLenght("Tcomm")
	)
	for _, f := range pT.headers {
		switch f {
		case "PID":
			formated = fmt.Sprintf("%%%dv ", PID)
		case "TTY":
			formated = fmt.Sprintf("%%-%dv\t ", TTY)
		case "STAT":
			formated = fmt.Sprintf("%%-%dv\t", STAT)
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
// TODO: a nice clean way to turn /proc/pid/stat into a struct. (trying now)
// There has to be a way.
func ps(pT ProcessTable) error {
	// sorting ProcessTable by PID
	sort.Sort(pT)

	if flags.x {
		pT.headers = []string{"PID", "TTY", "STAT", "TIME", "CMD"}
		pT.fields = []string{"Pid", "Ctty", "State", "Time", "Tcomm"}
	} else {
		pT.headers = []string{"PID", "TTY", "TIME", "CMD"}
		pT.fields = []string{"Pid", "Ctty", "Time", "Tcomm"}
	}

	mainProcess := pT.GetProcess(mainPID)

	pT.PrepareString()
	pT.PrintHeader()
	for index, p := range pT.table {
		if flags.notSid {
			if p.Sid == p.Pid { // ignore if is a session leader
				continue
			}
			if fieldString(p.TTYPgrp) == "-1" { // without any terminal linked
				continue
			}
		} else if flags.all == false { // default case
			uid, err := p.GetUid()
			if err != nil {
				return err
			}
			if p.Sid != mainProcess.Sid { // ignore is not same session
				continue
			}
			if eUID != uid {
				continue
			}
		}

		pT.PrintProcess(index)
	}

	return nil

}

func main() {
	pT := ProcessTable{}
	if err := pT.LoadTable(); err != nil {
		log.Fatal(err)
	}

	err := ps(pT)
	if err != nil {
		log.Fatal(err)
	}
}
