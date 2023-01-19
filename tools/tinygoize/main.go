// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// usage: invoke this with a list of directories.
// For each directory, it will try a tinygobuild with
// CGO_ENABLED=0, GOARCH=amd64, and GOOS=linux.
// If the tinygo build fails, it will use ast package
// to rewrite //go:build lines as follows:
// the line starts as //go:build expr
// it is rewritten to //go:build !tinygo && (expr)
// When the file is written, the expression seems
// to be simplified.

package main

import (
	"bytes"
	"flag"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const goBuild = "//go:build "

func main() {
	flag.Parse()

	p := printer.Config{Mode: printer.UseSpaces | printer.TabIndent, Tabwidth: 8}
	for _, d := range flag.Args() {
		c := exec.Command("tinygo", "build")
		c.Dir = d
		c.Stdout, c.Stderr = os.Stdout, os.Stderr
		c.Env = append(os.Environ(), "GOOS=linux", "CGO_ENABLED=0", "GOARCH=amd64")
		if err := c.Run(); err == nil {
			continue
		}
		files, err := filepath.Glob(filepath.Join(d, "*"))
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			if !strings.HasSuffix(file, ".go") {
				continue
			}
			log.Printf("Process %s", file)
			b, err := os.ReadFile(file)
			if err != nil {
				log.Fatal(err)
			}
			fset := token.NewFileSet() // positions are relative to fset
			f, err := parser.ParseFile(fset, file, string(b), parser.ParseComments|parser.SkipObjectResolution)
			if err != nil {
				log.Fatalf("parsing\n%v\n:%v", string(b), err)
			}
		done:
			for _, cg := range f.Comments {
				for _, c := range cg.List {
					if !strings.HasPrefix(c.Text, goBuild) {
						continue
					}
					c.Text = goBuild + "!tinygo && (" + c.Text[len(goBuild):] + ")"
					break done
				}
			}
			// Complete source file.
			var buf bytes.Buffer
			if err = p.Fprint(&buf, fset, f); err != nil {
				log.Fatalf("Printing:%v", err)
			}
			if err := os.WriteFile(file, buf.Bytes(), 0o644); err != nil {
				log.Fatal(err)
			}
		}
	}
}
