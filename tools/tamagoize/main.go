// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This builds a tinygo image that does not need
// fork/exec.

package main

import (
	"flag"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	// this used to include rush, but tamago has their own shell
	// should you want to bring rush back, you will need to grab the
	// support from tinygobb.
	list = flag.String("l", "", "comma-separated list of u-root commands to build")
	debug = flag.Bool("d", false, "enable debug prints")
	v = func(string, ...any) {}
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
		v("Running u-root: %v", err)
	}

	build, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		log.Fatal(err)
	}
	if len(build) != 1 {
		log.Fatalf("can not find unique builddir from %q, got %q", dir, build)
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
		v("replaced os.Exit in %q, output %s", path, s)
		if err := ioutil.WriteFile(path, []byte(s), 0644); err != nil {
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

	os.Exit(0)

	// this is where we would build tamago and maybe it would work.
	c = exec.Command("go", "build", "-tags", "tinygo")
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
