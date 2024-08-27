// Copyright 2016-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

const (
	defaultGlob = "/proc"
	userHZ      = 100
)

var (
	psglob string
	// by convention, the first element of the path is "/proc"
	// This allows us to point to any place as our "/proc"
	procdir = "/proc"
)

// Process contains both kernel-dependent and kernel-independent information.
type Process struct {
	process
	status  string
	cmdline string
	stat    string
	Pidno   int // process id #
	uid     int
}

// table content of stat file defined by:
// https://www.kernel.org/doc/Documentation/filesystems/proc.txt (2009)
// Section (ctrl + f) : Table 1-4: Contents of the stat files (as of 2.6.30-rc7)
type process struct {
	Pid         string // process id name
	Cmd         string // filename of the executable
	State       string // state (R is running, S is sleeping, D is sleeping in an uninterruptible wait, Z is zombie, T is traced or stopped)
	Ppid        string // process id of the parent process
	Pgrp        string // pgrp of the process
	Sid         string // session id
	TTYNr       string // tty the process uses
	TTYPgrp     string // pgrp of the tty
	Flags       string // task flags
	MinFlt      string // number of minor faults
	CminFlt     string // number of minor faults with child's
	MajFlt      string // number of major faults
	CmajFlt     string // number of major faults with child's
	Utime       string // user mode jiffies
	Stime       string // kernel mode jiffies
	Cutime      string // user mode jiffies with child's
	Cstime      string // kernel mode jiffies with child's
	Priority    string // priority level
	Nice        string // nice level
	NumThreads  string // number of threads
	ItRealValue string // (obsolete, always 0)
	StartTime   string // time the process started after system boot
	Vsize       string // virtual memory size
	Rss         string // resident set memory size
	Rsslim      string // current limit in bytes on the rss
	StartCode   string // address above which program text can run
	EndCode     string // address below which program text can run
	StartStack  string // address of the start of the main process stack
	Esp         string // current value of ESP
	Eip         string // current value of EIP
	Pending     string // bitmap of pending signals
	Blocked     string // bitmap of blocked signals
	Sigign      string // bitmap of ignored signals
	Sigcatch    string // bitmap of caught signals
	Wchan       string // place holder, used to be the wchan address, use /proc/PID/wchan
	Zero1       string // ignored
	Zero2       string // ignored
	ExitSignal  string // signal to send to parent thread on exit
	TaskCPU     string // which CPU the task is scheduled on
	RtPriority  string // realtime priority
	Policy      string // scheduling policy (man sched_setscheduler)
	BlkioTicks  string // time spent waiting for block IO
	Gtime       string // guest time of the task in jiffies
	Cgtime      string // guest time of the task children in jiffies
	StartData   string // address above which program data+bss is placed
	EndData     string // address below which program data+bss is placed
	StartBrk    string // address above which program heap can be expanded with brk()
	ArgStart    string // address above which program command line is placed
	ArgEnd      string // address below which program command line is placed
	EnvStart    string // address above which program environment is placed
	EnvEnd      string // address below which program environment is placed
	ExitCode    string // the thread's exit_code in the form reported by the waitpid system call (end of stat)
	Ctty        string // extra member (don't parsed from stat)
	Time        string // extra member (don't parsed from stat)
}

// Parse all content of stat to a Process Struct
// by gived the pid (linux)
func (p *Process) readStat(s string) error {
	fields := strings.Split(s, " ")
	// set struct fields from stat file data
	v := reflect.ValueOf(&p.process).Elem()
	for i := 0; i < len(fields); i++ {
		fieldVal := v.Field(i)
		fieldVal.Set(reflect.ValueOf(fields[i]))
	}

	p.Time = p.getTime()
	p.Ctty = p.getCtty()
	p.Cmd = strings.TrimSuffix(strings.TrimPrefix(p.Cmd, "("), ")")
	if x && p.cmdline != "" {
		p.Cmd = p.cmdline
	}

	return nil
}

// Parse data from various strings in the Process struct
func (p *Process) Parse() error {
	err := p.readStat(p.stat)
	if err != nil {
		return err
	}
	if p.uid, err = p.GetUID(); err != nil {
		return err
	}
	return nil
}

// ctty returns the ctty or "?" if none can be found.
// TODO: an right way to get ctty by p.TTYNr and p.TTYPgrp
func (p process) getCtty() string {
	if tty, err := os.Readlink(filepath.Join(procdir, p.Pid, "fd/0")); err != nil {
		return "?"
	} else if p.TTYPgrp != "-1" {
		if len(tty) > 5 && tty[:5] == "/dev/" {
			tty = tty[5:]
		}
		return tty
	}
	return "?"
}

// Get a named field of stat type
// e.g.: p.getField("Pid") => '1'
func (p *process) getField(field string) string {
	v := reflect.ValueOf(p).Elem()
	return fmt.Sprintf("%v", v.FieldByName(field))
}

// Search for attributes about the process
func (p *Process) Search(field string) string {
	return p.process.getField(field)
}

// GetUID gets the UID of the process from the status string
func (p Process) GetUID() (int, error) {
	lines := strings.Split(p.status, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Uid") {
			fields := strings.Split(line, "\t")
			return strconv.Atoi(fields[1])
		}
	}

	return -1, fmt.Errorf("no Uid string in %s", p.status)
}

// Get total time stat formated hh:mm:ss
func (p process) getTime() string {
	utime, _ := strconv.Atoi(p.Utime)
	stime, _ := strconv.Atoi(p.Stime)
	jiffies := utime + stime

	tsecs := jiffies / userHZ
	secs := tsecs % 60
	mins := (tsecs / 60) % 60
	hrs := tsecs / 3600

	return fmt.Sprintf("%02d:%02d:%02d", hrs, mins, secs)
}

func getAllGlobNames() []string {
	psglob = os.Getenv("UROOT_PSPATH")
	if psglob == "" {
		// The reason we glob with stat, even though
		// we strip it off later, is it is a cheap way
		// to ensure we're getting a process directory
		// and not some other weird thing in /proc.
		psglob = defaultGlob
	}
	l := filepath.SplitList(psglob)
	if len(l) > 0 {
		procdir = l[0]
	}
	return l
}

// Create a set of stat file names from an array of globs
func getAllStatNames(globs []string) ([]string, error) {
	var list []string
	for _, g := range globs {
		l, err := filepath.Glob(filepath.Join(g, "[0-9]*/stat"))
		if err != nil {
			log.Printf("Glob err on %s: %v", g, err)
			continue
		}
		list = append(list, l...)
	}
	if len(list) == 0 {
		return nil, fmt.Errorf("no files found in %q; check if proc is mounted", psglob)
	}
	return list, nil
}

func file(s string) (string, error) {
	b, err := os.ReadFile(s)
	return string(b), err
}

func (pT *ProcessTable) doTable(statFileNames []string) error {
	var err error
	for _, stat := range statFileNames {
		p := &Process{}

		// log.Printf("Check %s", stat)
		// ps is a snapshot in time of /proc. Hence we want to grab
		// all the files we need in as close to an instant in time as
		// we can.
		// Read the files. It may have vanished or we may not have
		// access; we do not consider those to be errors.
		// if *any* of the files are not there, just skip this pid.
		p.stat, err = file(stat)
		if err != nil {
			continue
		}
		d := filepath.Dir(stat)
		pid := filepath.Base(d)
		pidno, err := strconv.Atoi(pid)
		if err != nil {
			return fmt.Errorf("last element of %v is not a number", pid)
		}
		p.status, err = file(filepath.Join(d, "status"))
		if err != nil {
			continue
		}
		if x {
			p.cmdline, err = file(filepath.Join(d, "cmdline"))
			if err != nil {
				continue
			}
		}
		// if filepath.Base is *not* proc, then use it, else
		// it's just the directory containing the pid.
		proot := filepath.Dir(d)
		// log.Printf("procdir %v d %v proot %v", procdir, d, proot)
		if proot != procdir {
			pid = filepath.Join(filepath.Base(proot), pid)
		}
		p.Pidno = pidno
		if err := p.Parse(); err != nil {
			return err
		}
		p.Pid = pid
		// log.Printf("stat is %v p is %v", stat,p)
		if p.Pidno == os.Getpid() {
			pT.mProc = p
		}
		pT.table = append(pT.table, p)
	}
	// if mProc is nil, something is really wrong.
	if pT.mProc == nil && len(pT.table) > 0 {
		pT.mProc = pT.table[0]
	}
	return nil
}

// LoadTable creates a ProcessTable containing stats on all processes.
// We use UROOT_PSPATH if set, else the default glob
// of /proc/[0-9]*/stat.
// We want to allow ps to run against the standard /proc but also
// proc mounted over a network in, e.g., /netproc/host/pid/...
// (i.e. we mount node:/proc on /netproc/node)
// The question then becomes what to store for the pid.
// For /proc, it's easy: strip the first directory component.
// For additional directories, e.g. /netproc/host/[0-9]*/stat,
// we can follow the same rule: strip the first component.
// We will do that for now and see if it works; if not we'll
// need more complex processing for UROOT_PSPATH.
func (pT *ProcessTable) LoadTable() error {
	g := getAllGlobNames()
	n, err := getAllStatNames(g)
	if err != nil {
		return err
	}
	return pT.doTable(n)
}
