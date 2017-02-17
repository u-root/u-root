// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// Install command from a go source file.
//
// Synopsis:
//     SYMLINK [ARGS...]
//     installcommand [-v] [-ludicrous] COMMAND [ARGS...]
//
// Description:
//     u-root commands are lazily compiled. Uncompiled commands in the /bin
//     directory are symbolic links to installcommand. When executed through
//     the symbolic link, installcommand will build the command from source and
//     exec it.
//
//     The second form allows commands to be installed and exec'ed without a
//     symbolic link. In this form the debug arguments `-v` and `-ludicrous`
//     can be passed into installcommand.
//
// Options:
//     -v: print all build commands
//     -ludicrous: print out ALL the output from the go build commands
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
	verbose   = flag.Bool("v", false, "print all build commands")
	ludicrous = flag.Bool("ludicrous", false, "print out ALL the output from the go build commands")
	debug     = func(string, ...interface{}) {}
	useExec   = flag.Bool("exec", false, "Use a direct exec system call instead of cmd.Run for the child")
)

type form struct {
	// Name of the command, ex: "ls"
	cmdName string
	// Args passed to the command, ex: {"-l", "-R"}
	cmdArgs []string
	// Args intended for installcommand
	verbose   bool
	ludicrous bool
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: installcommand [-v] [-ludicrous] COMMAND [ARGS...]\n")
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
	//     installcommand [-v] [-ludicrous] COMMAND [ARGS...]
	flag.Parse()
	if flag.NArg() < 1 {
		log.Println("Second form requires a COMMAND argument")
		usage()
	}
	return form{
		cmdName:   flag.Arg(0),
		cmdArgs:   flag.Args()[1:],
		verbose:   *verbose,
		ludicrous: *ludicrous,
	}
}

func main() {
	form := parseCommandLine()

	a := []string{"install"}
	if form.verbose {
		debug = log.Printf
		a = append(a, "-x")
	}

	debug("Command name: %v\n", form.cmdName)
	destDir := "/ubin"
	destFile := path.Join(destDir, form.cmdName)

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

	if *useExec {
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
		log.Fatal(err)
	}
}
