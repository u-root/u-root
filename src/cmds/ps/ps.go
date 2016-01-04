// Copyright 2013 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// ps reads the /proc and prints out nice things about what it finds.
// /proc in linux has grown by a process of Evolution, so it's messy.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	allProc = "^[0-9]+$"
	proc    = "/proc"
	USER_HZ = 100
)

var (
	flags struct {
		all bool
		x   bool
	}
	cmd  = "ps [-Aex]"
	eUID = os.Geteuid()
)

func usage() {
	fmt.Fprintln(os.Stderr, "Usage:", cmd)
	flag.PrintDefaults()
	os.Exit(1)
}

func init() {
	flag.BoolVar(&flags.all, "A", false, "Select all processes.  Identical to -e.")
	flag.BoolVar(&flags.all, "e", false, "Select all processes.  Identical to -A.")
	flag.BoolVar(&flags.x, "x", false, "BSD-Like style, with Stat Column and long CommandLine")
	flag.Parse()
	flag.Usage = usage
}

// Type set up to use a struct can iterate using reflect
type Field interface{}

func createField(name string) Field {
	return name
}

func fieldName(f Field) string {
	return fmt.Sprintf("%v", f)
}

func fieldLen(f Field) int {
	return len(fieldName(f))
}

// table content of stat file defined by:
// https://www.kernel.org/doc/Documentation/filesystems/proc.txt (2009)
// Section (ctrl + f) : Table 1-4: Contents of the stat files (as of 2.6.30-rc7)
type Process struct {
	// process id
	Pid Field
	// filename of the executable
	Tcomm Field
	// state (R is running, S is sleeping, D is sleeping in an
	// uninterruptible wait, Z is zombie, T is traced or stopped)
	State Field
	// process id of the parent process
	Ppid Field
	// pgrp of the process
	Pgrp Field
	// session id
	Sid Field
	// tty the process uses
	TTYNr Field
	// pgrp of the tty
	TTYPgrp Field
	// task flags
	Flags Field
	// number of minor faults
	MinFlt Field
	// number of minor faults with child's
	CminFlt Field
	// number of major faults
	MajFlt Field
	// number of major faults with child's
	CmajFlt Field
	// user mode jiffies
	Utime Field
	// kernel mode jiffies
	Stime Field
	// user mode jiffies with child's
	Cutime Field
	// kernel mode jiffies with child's
	Cstime Field
	// priority level
	Priority Field
	// nice level
	Nice Field
	// number of threads
	NumThreads Field
	// (obsolete, always 0)
	ItRealValue Field
	// time the process started after system boot
	StartTime Field
	// virtual memory size
	Vsize Field
	// resident set memory size
	Rss Field
	// current limit in bytes on the rss
	Rsslim Field
	// address above which program text can run
	StartCode Field
	// address below which program text can run
	EndCode Field
	// address of the start of the main process stack
	StartStack Field
	// current value of ESP
	Esp Field
	// current value of EIP
	Eip Field
	// bitmap of pending signals
	Pending Field
	// bitmap of blocked signals
	Blocked Field
	// bitmap of ignored signals
	Sigign Field
	// bitmap of caught signals
	Sigcatch Field
	// place holder, used to be the wchan address, use /proc/PID/wchan
	Wchan Field
	// ignored
	Zero1 Field
	// ignored
	Zero2 Field
	// signal to send to parent thread on exit
	ExitSignal Field
	// which CPU the task is scheduled on
	TaskCpu Field
	// realtime priority
	RtPriority Field
	// scheduling policy (man sched_setscheduler)
	Policy Field
	// time spent waiting for block IO
	BlkioTicks Field
	// guest time of the task in jiffies
	Gtime Field
	// guest time of the task children in jiffies
	Cgtime Field
	// address above which program data+bss is placed
	StartData Field
	// address below which program data+bss is placed
	EndData Field
	// address above which program heap can be expanded with brk()
	StartBrk Field
	// address above which program command line is placed
	ArgStart Field
	// address below which program command line is placed
	ArgEnd Field
	// address above which program environment is placed
	EnvStart Field
	// address below which program environment is placed
	EnvEnd Field
	// the thread's exit_code in the form reported by the waitpid system call
	ExitCode Field // end of table (from stat)
	Time     Field // extra member (don't parsed from stat)
	Ctty     Field // extra member (don't parsed from stat)
}

// Parse all content of Stat to a Process Struct
// by gived the pid
func (p *Process) ReadStat(pid string) error {
	b, err := ioutil.ReadFile(path.Join(proc, pid, "stat"))

	if err != nil {
		return err
	}

	fields := strings.Split(string(b), " ")

	// set struct fields from stat file data
	v := reflect.ValueOf(p).Elem()
	for i := 0; i < len(fields); i++ {
		setTo := createField(fields[i])
		fieldVal := v.Field(i)
		fieldVal.Set(reflect.ValueOf(setTo))
	}

	// get time space total
	p.Time = createField(p.GetTime())
	// current tty
	p.Ctty = createField(p.GetCtty())

	if flags.x {
		// breaks when cmdline is very big
		// implement limit width screen for each row to use that
		if false {
			cmdline, err := p.LongCmdLine()
			if err != nil {
				return err
			}
			p.Tcomm = createField(cmdline)
		}

	} else {
		// remove parenthesis of filename of executable
		cmd := fieldName(p.Tcomm)
		p.Tcomm = createField(cmd[1 : len(cmd)-1])
	}

	return nil
}

// ctty returns the ctty or "?" if none can be found.
// TODO: an right way to get ctty by p.TTYNr and p.TTYPgrp
func (p Process) GetCtty() string {
	if tty, err := os.Readlink(path.Join(proc, fieldName(p.Pid), "fd/2")); err != nil {
		return "?"
	} else if p.TTYPgrp != "-1" {
		if len(tty) > 5 && tty[:5] == "/dev/" {
			tty = tty[5:]
		}
		return tty
	}
	return "?"
}

// Get a named field of Process type
// e.g.: p.GetField("Pid") => '1'
func (p Process) GetField(field string) string {
	v := reflect.ValueOf(&p).Elem()
	return fieldName(v.FieldByName(field))
}

func (p Process) GetUid() (int, error) {
	b, err := ioutil.ReadFile(path.Join(proc, fieldName(p.Pid), "status"))

	var uid int
	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Uid") {
			fields := strings.Split(line, "\t")
			uid, err = strconv.Atoi(fields[1])
			break
		}
	}

	return uid, err

}

// change p.Tcomm to long command line with args
func (p Process) LongCmdLine() (string, error) {
	b, err := ioutil.ReadFile(path.Join(proc, fieldName(p.Pid), "cmdline"))

	if err != nil {
		return "", err
	}

	return string(b[:]), nil
}

// Get total time Process formated hh:mm:ss
func (p Process) GetTime() string {
	fields := []string{
		"Utime",
		"Stime",
		"Cutime",
		"Cstime",
	}

	jiffies := 0
	for _, field := range fields {
		value, _ := strconv.Atoi(p.GetField(field))
		jiffies += value
	}

	totalSeconds := jiffies / USER_HZ
	seconds := int(totalSeconds % 60)
	minutes := int((totalSeconds / 60) % 60)
	hours := totalSeconds / 360

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

// main process table of ps
type ProcessTable []Process

func (pT ProcessTable) Len() int {
	return len(pT)
}

// to use on sort.Sort
func (pT ProcessTable) Less(i, j int) bool {
	a, _ := strconv.Atoi(fieldName(pT[i].Pid))
	b, _ := strconv.Atoi(fieldName(pT[j].Pid))
	return a < b
}

// to use on sort.Sort
func (pT ProcessTable) Swap(i, j int) {
	pT[i], pT[j] = pT[j], pT[i]
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
	for _, p := range pT {
		slice = append(slice, fieldLen(p.GetField(field)))
	}

	return max(slice)
}

var processTable = ProcessTable{}

// used to make a pretty output like ps-gnu
type ColumnsWidth struct {
	PID  int
	TTY  int
	STAT int
	TIME int
	CMD  int
}

// get all the max width lenghts from processTable
func (cW *ColumnsWidth) Init(pT ProcessTable) {
	cW.PID = pT.MaxLenght("Pid")
	cW.TTY = pT.MaxLenght("Ctty")
	cW.STAT = pT.MaxLenght("State")
	cW.TIME = pT.MaxLenght("Time")
	cW.CMD = pT.MaxLenght("Tcomm")
}

// For now, just read /proc/pid/stat and dump its brains.
// TODO: a nice clean way to turn /proc/pid/stat into a struct. (trying now)
// There has to be a way.
func ps(pT ProcessTable) error {
	// sorting ProcessTable by PID
	sort.Sort(pT)
	// update command line
	cW := &ColumnsWidth{}
	cW.Init(pT)

	// Header
	var formatedRow string
	if flags.x {
		formatedRow = "%-*v %-*v\t%-*v\t%-*v %-*v\n"
		fmt.Printf(
			formatedRow,
			cW.PID, "PID",
			cW.TTY, "TTY",
			cW.STAT, "STAT",
			cW.TIME, "TIME",
			cW.CMD, "CMD",
		)
	} else {
		formatedRow = "%-*v %-*v\t%-*v %-*v\n"
		fmt.Printf(
			formatedRow,
			cW.PID, "PID",
			cW.TTY, "TTY",
			cW.TIME, "TIME",
			cW.CMD, "CMD",
		)
	}

	for _, p := range pT {
		uid, err := p.GetUid()
		if err != nil {
			return err
		}
		row := fmt.Sprintf(
			formatedRow,
			cW.PID, p.Pid,
			cW.TTY, p.Ctty,
			cW.TIME, p.Time,
			cW.CMD, p.Tcomm,
		)

		if flags.x {
			row = fmt.Sprintf(
				formatedRow,
				cW.PID, p.Pid,
				cW.TTY, p.Ctty,
				cW.STAT, p.State,
				cW.TIME, p.Time,
				cW.CMD, p.Tcomm,
			)
		}

		if flags.all == false && eUID != uid {
			continue
		}

		fmt.Printf(row)
	}

	return nil

}

func main() {
	nf := allProc
	flag.Parse()
	pf := regexp.MustCompile(nf)

	filepath.Walk(proc, func(name string, fi os.FileInfo, err error) error {
		if err != nil {
			log.Printf("%v: %v\n", name, err)
			return err
		}
		if name == proc {
			return nil
		}

		if pf.Match([]byte(fi.Name())) {
			p := &Process{}
			if err := p.ReadStat(fi.Name()); err != nil {
				log.Print(err)
				return err
			}
			processTable = append(processTable, *p)
		}

		return filepath.SkipDir
	})

	err := ps(processTable)
	if err != nil {
		log.Fatal(err)
	}
}
