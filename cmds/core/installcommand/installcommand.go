// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// installcommand installs a command from Go source files.
//
// Synopsis:
//     SYMLINK [ARGS...]
//     installcommand [INSTALLCOMMAND_ARGS...] COMMAND [ARGS...]
//
// Description:
//     In u-root's source mode, uncompiled commands in the /bin directory are
//     symbolic links to installcommand. When executed through the symbolic
//     link, installcommand will build the command from source and exec it.
//
//     The second form allows commands to be installed and exec'ed without a
//     symbolic link. In this form additional arguments such as `-v` and
//     `-ludicrous` can be passed into installcommand.
//
// Options:
//     -lowpri:    the scheduler priority to lowered before starting
//     -exec:      build and exec the command
//     -force:     do not build if a file already exists at the destination
//     -v:         print all build commands
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/u-root/u-root/pkg/golang"
	"github.com/u-root/u-root/pkg/upath"
)

var (
	lowpri = flag.Bool("lowpri", false, "the scheduler priority is lowered before starting")
	exe    = flag.Bool("exec", true, "build AND execute the command")
	force  = flag.Bool("force", false, "build even if a file already exists at the destination")

	verbose = flag.Bool("v", false, "print all build commands")
	r       = upath.UrootPath
)

type form struct {
	// Name of the command, ex: "ls"
	cmdName string
	// Args passed to the command, ex: {"-l", "-R"}
	cmdArgs []string

	// Args intended for installcommand
	lowPri  bool
	exec    bool
	force   bool
	verbose bool
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
		// This is almost certain to be a symlink, and it's no harm
		// to check it.
		f := upath.ResolveUntilLastSymlink(os.Args[0])
		return form{
			cmdName: filepath.Base(f),
			cmdArgs: os.Args[1:],
			lowPri:  *lowpri,
			exec:    *exe,
			force:   *force,
			verbose: *verbose,
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
		cmdName: flag.Arg(0),
		cmdArgs: flag.Args()[1:],
		lowPri:  *lowpri,
		exec:    *exe,
		force:   *force,
		verbose: *verbose,
	}
}

// run runs the command with the information from form.
// Since run can potentially never return, since it can use Exec,
// it should never return in any other case. Hence, if all goes well
// at the end, we os.Exit(0)
func run(n string, form form) {
	cmd := exec.Command(n, form.cmdArgs...)
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
	os.Exit(0)
}

func main() {
	form := parseCommandLine()

	if form.lowPri {
		if err := syscall.Setpriority(syscall.PRIO_PROCESS, 0, 20); err != nil {
			log.Printf("Cannot set low priority: %v", err)
		}
	}

	destFile := filepath.Join(r("/ubin"), form.cmdName)

	// Is the command there? This covers a race condition
	// in that some other process may have caused it to be
	// built.
	if _, err := os.Stat(destFile); err == nil {
		if !form.exec {
			os.Exit(0)
		}
		run(destFile, form)
	}

	env := golang.Default()
	env.Context.GOROOT = r("/go")
	env.Context.GOPATH = r("/")
	env.Context.CgoEnabled = false

	var srcDir string
	err := filepath.Walk(r("/src"), func(p string, fi os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if fi.IsDir() && filepath.Base(p) == form.cmdName {
			// Make sure it's an actual Go command.
			pkg, err := env.PackageByPath(p)
			if err == nil && pkg.IsCommand() {
				srcDir = p
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	if len(srcDir) == 0 {
		log.Fatalf("Can not find source code for %q", form.cmdName)
	}

	if err := env.BuildDir(srcDir, destFile, golang.BuildOpts{}); err != nil {
		log.Fatalf("Couldn't compile %q: %v", form.cmdName, err)
	}

	if form.exec {
		run(destFile, form)
	}
}
