// Copyright 2013 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
ps reads the /proc and prints out nice things about what it finds. 
/proc in linux has grown by a process of Evilution, so it's messy.
*/

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	allProc = "^[0-9][0-9]*$"
	proc = "/proc"
)

// ctty returns the ctty or "none" if none can be found.
func ctty(pid string) string {
	if tty, err := os.Readlink(path.Join(proc, pid, "fd/2")); err != nil {
		return "none"
	} else {
		if len(tty) > 5 && tty[:5] == "/dev/" {
			tty = tty[5:]
		}
		return tty
	}
}
// For now, just read /proc/pid/stat and dump its brains.
// TODO; a nice clean way to turn /proc/pid/stat into a struct.
// There has to be a way.
func ps(pid string) error {
	b, err := ioutil.ReadFile(path.Join(proc, pid, "stat"))
	if err != nil {
		return err
	}
	l := strings.Split(string(b), " ")
	// sum the times. But what's the divisor? 
	times := 0
	for _, v := range l[13:17] {
		t, err := strconv.Atoi(v)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
		times = times + t
	}
	// convert to minutes, assume USER HZ is 100, suck.
	times = times / (100 * 60)
	fmt.Printf("%v\t%v\t%v\t%v\n", l[0], ctty(pid), times, l[1][1:len(l[1])-1])
	return nil
}

func main() {
	nf := allProc
	flag.Parse()
	pf := regexp.MustCompile(nf)
	fmt.Printf("PID\tTTY\tTIME\tCMD\n")
	filepath.Walk(proc, func(name string, fi os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v: %v\n", name, err)
			return err
		}
		if name == proc {
			return nil
		}

		if pf.Match([]byte(fi.Name())) {
			if err := ps(fi.Name()); err != nil {
				fmt.Fprintf(os.Stderr, "%v", err)
			}
		}
		
		return filepath.SkipDir
	})
}
