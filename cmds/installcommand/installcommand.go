// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// Install command from a go source file.
//
// Synopsis:
//     SYMLINK [ARGS...]
//     installcommand [INSTALLCOMMAND_ARGS...] COMMAND [ARGS...]
//
// Description:
//     u-root commands are lazily compiled. Uncompiled commands in the /bin
//     directory are symbolic links to installcommand. When executed through
//     the symbolic link, installcommand will build the command from source and
//     exec it.
//
//     The second form allows commands to be installed and exec'ed without a
//     symbolic link. In this form additional arguments such as `-v` and
//     `-ludicrous` can be passed into installcommand.
//
// Options:
//     -lowpri:    the scheduler priority to lowered before starting
//     -ludicrous: print out ALL the output from the go build commands
//     -noexec:    do not execute the process after building
//     -noforce:   do not build if a file already exists at the destination
//     -v:         print all build commands
import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"

	"github.com/u-root/u-root/uroot"
)

var (
	urpath    = "/go/bin:/ubin:/buildbin:/usr/local/bin:"
	lowpri    = flag.Bool("lowpri", false, "the scheduler priority is lowered before starting")
	ludicrous = flag.Bool("ludicrous", false, "print out ALL the output from the go build commands")
	noexec    = flag.Bool("noexec", false, "do not execute the process after building")
	noforce   = flag.Bool("noforce", false, "do not build if a file already exists at the destination")
	verbose   = flag.Bool("v", false, "print all build commands")
	debug     = func(string, ...interface{}) {}
)

type form struct {
	// Name of the command, ex: "ls"
	cmdName string
	// Args passed to the command, ex: {"-l", "-R"}
	cmdArgs []string
	// Args intended for installcommand
	lowPri    bool
	ludicrous bool
	noExec    bool
	noForce   bool
	verbose   bool
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: installcommand [INSTALLCOMMAND_ARGS...] COMMAND [ARGS...]\n")
	os.Exit(2)
}

// Parse the command line to determine the form.
func parseCommandLine() form {
	// First form:
	//     SYMLINK [ARGS...]
	if !strings.HasSuffix(os.Args[0], "installcommand") {
		return form{
			cmdName: path.Base(os.Args[0]),
			cmdArgs: os.Args[1:],
		}
	}

	// Second form:
	//     installcommand [INSTALLCOMMAND_ARGS...] COMMAND [ARGS...]
	flag.Parse()
	if flag.NArg() < 1 {
		log.Println("Second form requires a COMMAND argument")
		usage()
	}
	return form{
		cmdName:   flag.Arg(0),
		cmdArgs:   flag.Args()[1:],
		lowPri:    *lowpri,
		ludicrous: *ludicrous,
		noExec:    *noexec,
		noForce:   *noforce,
		verbose:   *verbose,
	}
}

func main() {
	form := parseCommandLine()

	if form.lowPri {
		if err := syscall.Setpriority(syscall.PRIO_PGRP, 0, 20); err != nil {
			log.Printf("Cannot set low priority: %v", err)
		}
	}

	a := []string{"install"}
	if form.verbose {
		debug = log.Printf
		a = append(a, "-x")
	}

	debug("Command name: %v\n", form.cmdName)
	destDir := "/ubin"
	destFile := path.Join(destDir, form.cmdName)

	// Optionally skip if already built.
	if form.noForce {
		if _, err := os.Stat(destFile); err == nil {
			os.Exit(0)
		}
	}

	cmd := exec.Command("go", append(a, path.Join(uroot.CmdsPath, form.cmdName))...)

	// Set GOGC if unset. The best value is determined empirically and
	// depends on the machine and Go version. For the workload of compiling
	// a small Go program, values larger than the default perform better.
	// See: /scripts/build_perf.sh
	if _, ok := os.LookupEnv("GOGC"); !ok {
		cmd.Env = append(os.Environ(), "GOGC=400")
	}

	cmd.Dir = "/"

	debug("Run %v", cmd)
	out, err := cmd.CombinedOutput()
	debug("installcommand: go build returned")

	if err != nil {
		p := os.Getenv("PATH")
		log.Fatalf("installcommand: trying to build {cmdName: %v, PATH %s, err %v, out %s}", form.cmdName, p, err, out)
	}

	if *ludicrous {
		debug(string(out))
	}

	if !form.noExec {
		if os.Getenv("INSTALLCOMMAND_NOFORK") == "1" {
			err = syscall.Exec(destFile, append([]string{form.cmdName}, form.cmdArgs...), os.Environ())
			// Regardless of whether err is nil, if Exec returns at all, it failed
			// at its job. Print an error and then let's see if a normal run can succeed.
			log.Printf("Failed to exec %s: %v", form.cmdName, err)
		}

		cmd = exec.Command(destFile)
		cmd.Args = append([]string{form.cmdName}, form.cmdArgs...)
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			exitErr, ok := err.(*exec.ExitError)
			if !ok {
				log.Fatal(err)
			}
			exitWithStatus(exitErr)
		}
	}
}
