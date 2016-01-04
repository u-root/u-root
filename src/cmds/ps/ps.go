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

// Type set up to use a struct can iterate using reflect
type Field interface{}

func createField(name string) Field {
	return name
}

func fieldString(f Field) string {
	return fmt.Sprintf("%s", f)
}

func fieldLen(f Field) int {
	return len(fieldString(f))
}

// table content of stat file defined by:
// https://www.kernel.org/doc/Documentation/filesystems/proc.txt (2009)
// Section (ctrl + f) : Table 1-4: Contents of the stat files (as of 2.6.30-rc7)
type Process struct {
	Pid         Field // process id
	Tcomm       Field // filename of the executable
	State       Field // state (R is running, S is sleeping, D is sleeping in an uninterruptible wait, Z is zombie, T is traced or stopped)
	Ppid        Field // process id of the parent process
	Pgrp        Field // pgrp of the process
	Sid         Field // session id
	TTYNr       Field // tty the process uses
	TTYPgrp     Field // pgrp of the tty
	Flags       Field // task flags
	MinFlt      Field // number of minor faults
	CminFlt     Field // number of minor faults with child's
	MajFlt      Field // number of major faults
	CmajFlt     Field // number of major faults with child's
	Utime       Field // user mode jiffies
	Stime       Field // kernel mode jiffies
	Cutime      Field // user mode jiffies with child's
	Cstime      Field // kernel mode jiffies with child's
	Priority    Field // priority level
	Nice        Field // nice level
	NumThreads  Field // number of threads
	ItRealValue Field // (obsolete, always 0)
	StartTime   Field // time the process started after system boot
	Vsize       Field // virtual memory size
	Rss         Field // resident set memory size
	Rsslim      Field // current limit in bytes on the rss
	StartCode   Field // address above which program text can run
	EndCode     Field // address below which program text can run
	StartStack  Field // address of the start of the main process stack
	Esp         Field // current value of ESP
	Eip         Field // current value of EIP
	Pending     Field // bitmap of pending signals
	Blocked     Field // bitmap of blocked signals
	Sigign      Field // bitmap of ignored signals
	Sigcatch    Field // bitmap of caught signals
	Wchan       Field // place holder, used to be the wchan address, use /proc/PID/wchan
	Zero1       Field // ignored
	Zero2       Field // ignored
	ExitSignal  Field // signal to send to parent thread on exit
	TaskCpu     Field // which CPU the task is scheduled on
	RtPriority  Field // realtime priority
	Policy      Field // scheduling policy (man sched_setscheduler)
	BlkioTicks  Field // time spent waiting for block IO
	Gtime       Field // guest time of the task in jiffies
	Cgtime      Field // guest time of the task children in jiffies
	StartData   Field // address above which program data+bss is placed
	EndData     Field // address below which program data+bss is placed
	StartBrk    Field // address above which program heap can be expanded with brk()
	ArgStart    Field // address above which program command line is placed
	ArgEnd      Field // address below which program command line is placed
	EnvStart    Field // address above which program environment is placed
	EnvEnd      Field // address below which program environment is placed
	ExitCode    Field // the thread's exit_code in the form reported by the waitpid system call (end of stat)
	Time        Field // extra member (don't parsed from stat)
	Ctty        Field // extra member (don't parsed from stat)
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

	p.Time = createField(p.GetTime())
	p.Ctty = createField(p.GetCtty())
	cmd := fieldString(p.Tcomm)
	p.Tcomm = createField(cmd[1 : len(cmd)-1])
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
	}

	return nil
}

// ctty returns the ctty or "?" if none can be found.
// TODO: an right way to get ctty by p.TTYNr and p.TTYPgrp
func (p Process) GetCtty() string {
	if tty, err := os.Readlink(path.Join(proc, fieldString(p.Pid), "fd/0")); err != nil {
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
	return fieldString(v.FieldByName(field))
}

// read UID of process based on or
func (p Process) GetUid() (int, error) {
	b, err := ioutil.ReadFile(path.Join(proc, fieldString(p.Pid), "status"))

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
	b, err := ioutil.ReadFile(path.Join(proc, fieldString(p.Pid), "cmdline"))

	if err != nil {
		return "", err
	}

	return string(b), nil
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
		v, _ := strconv.Atoi(p.GetField(field))
		jiffies += v
	}

	tsecs := jiffies / USER_HZ
	secs := int(tsecs % 60)
	mins := int((tsecs / 60) % 60)
	hrs := tsecs / 360

	return fmt.Sprintf("%02d:%02d:%02d", hrs, mins, secs)
}

// main process table of ps
// used to make more flexible
type ProcessTable struct {
	table   []Process // the matrix of all process
	headers []string  // each column to print
	fields  []string  // which fields of process to print, on order
	fstring []string  // formated strings
}

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
func (pT ProcessTable) PrintHeader() {
	for index, field := range pT.headers {
		formated := pT.fstring[index]
		fmt.Printf(formated, field)
	}
	fmt.Printf("\n")
}

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

func (pT *ProcessTable) LoadTable() {
	pf := regexp.MustCompile(allProc)
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
			pT.table = append(pT.table, *p)
		}

		return filepath.SkipDir
	})
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
			if p.Sid == p.Pid {
				continue
			}
			if fieldString(p.TTYPgrp) == "-1" {
				continue
			}
		} else if flags.all == false {
			uid, err := p.GetUid()
			if err != nil {
				return err
			}
			if p.Sid != mainProcess.Sid {
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
	pT.LoadTable()

	err := ps(pT)
	if err != nil {
		log.Fatal(err)
	}
}
