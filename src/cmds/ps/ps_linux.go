// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	allProc = "^[0-9]+$"
	proc    = "/proc"
	USER_HZ = 100
)

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

// Walk from the proc files
// and parsing them
func (pT *ProcessTable) LoadTable() error {
	pf := regexp.MustCompile(allProc)
	err := filepath.Walk(proc, func(name string, fi os.FileInfo, err error) error {
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

	if err.Error() == "skip this directory" {
		return nil
	}

	return err
}
