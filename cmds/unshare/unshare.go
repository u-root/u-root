// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Disassociate parts of the process execution context.
//
// Synopsis:
//     unshare [OPTIONS] [PROGRAM [ARGS]...]
//
// Description:
//     Go applications use multiple processes, and the Go user level scheduler
//     schedules goroutines onto those processes. For this reason, it is not
//     possible to use syscall.Unshare. A goroutine can call `syscall.Unshare`
//     from process m and the scheduler can resume that goroutine in process n,
//     which has not had the unshare operation! This is a known problem with
//     any system call that modifies the name space or file system context of
//     only one process as opposed to the entire Go application, i.e. all of
//     its processes. Examples include chroot and unshare. There has been
//     lively discussion of this problem but no resolution as of yet. In sum:
//     it is not possible to use `syscall.Unshare` from Go with any reasonable
//     expectation of success.
//
//     If PROGRAM is not specified, unshare defaults to /ubin/rush.
//     The unsharing options are highly kernel dependent, but most kernels
//     support at least one of them.
package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
)

var fatal= log.Fatalf

func main() {
	flag.Parse()

	a := flag.Args()
	if len(a) == 0 {
		a = []string{"/ubin/rush", "rush"}
	}

	c := exec.Command(a[0], a[1:]...)
	c.Stdin, c.Stdout, c.Stderr = os.Stdin, os.Stdout, os.Stderr
	setSysProcAttr(c)

	if err := c.Run(); err != nil {
		fatal(err.Error())
	}
}
