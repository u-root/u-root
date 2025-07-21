// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This builds a tinygo image that does not need
// fork/exec.

package main

import (
	"flag"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	list   = flag.String("l", "cmds/exp/rush", "comma-separated list of u-root commands to build")
	flash  = flag.Bool("flash", false, "whether to flash as well as build")
	target = flag.String("target", "microbit-v2", "target name")
	tinygo = flag.Bool("tinygo", false, "use tinygo, not go")
	code   = `// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build tinygo

package bbrush

import bbmain "bb.u-root.com/bb/pkg/bbmain"
import "log"

func runone(c*Command) error {
   // put a recover here for the panic at some point.
   os.Args = append([]string{c.cmd}, c.argv...)
   flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.PanicOnError)
   log.Printf("run %v", c)
   return bbmain.Run(c.cmd)
}

`
)

func main() {
	flag.Parse()

	dir, err := os.MkdirTemp("", "tiny")
	if err != nil {
		log.Fatal(err)
	}

	c := exec.Command("./u-root", append([]string{"-tags", "tinygo", "-tmpdir", dir}, strings.Split(*list, ",")...)...)
	c.Stdout, c.Stderr = os.Stdout, os.Stderr
	if err := c.Run(); err != nil {
		// An error is expected, for now.
		log.Printf("Running u-root: %v", err)
	}

	build, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		log.Fatal(err)
	}
	if len(build) != 1 {
		log.Fatalf("can not find unique builddir from %q, got %q", dir, build)
	}

	rushdir := filepath.Join(build[0], "src/github.com/u-root/u-root/cmds/exp/rush/")

	// Fixup rush.
	if err := os.WriteFile(filepath.Join(rushdir, "rush_tinygo.go"), []byte(code), 0o644); err != nil {
		log.Fatal(err)
	}

	// Walk the tree, find all go files, in each file, replace os.Exit
	// with //.Exit. Sleazy.
	// This is the wrong way to do this, should use AST of course, as we do it with
	// the busybox. I was in a hurry.
	err = filepath.WalkDir(build[0], func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		b, errRead := os.ReadFile(path)
		if errRead != nil {
			return errRead
		}
		s := strings.ReplaceAll(string(b), "os.Exit", "//.Exit")
		log.Printf("replaced os.Exit in %q, output %s", path, s)
		if err := os.WriteFile(path, []byte(s), 0o644); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Walkdir %q: %v", build[0], err)
	}

	// The rewrite step may result in os no longer being needed in some
	// source. run imports.
	c = exec.Command("goimports", "-w", ".")
	c.Dir = build[0]
	c.Stdout, c.Stderr = os.Stdout, os.Stderr
	log.Printf("Now run imports: %v in %q", c, c.Dir)
	if err := c.Run(); err != nil {
		log.Fatalf("Running go: %v, %v", c, err)
	}

	c = exec.Command("go", "build", "-tags", "tinygo")
	if *tinygo {
		action := "build"
		if *flash {
			action = "flash"
		}
		c = exec.Command("tinygo", action, "-target", *target)
	}
	c.Dir = filepath.Join(build[0], "src/bb.u-root.com/bb")
	c.Stdout, c.Stderr = os.Stdout, os.Stderr
	log.Printf("Now compile it: %v", c)
	if err := c.Run(); err != nil {
		log.Fatalf("Running go: %v, %v", c, err)
	}

	bin := filepath.Join(c.Dir, "bb")
	fi, err := os.Stat(bin)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("The binary is: %q, info %v", bin, fi)
}
