// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// usage: invoke this with a list of directories. For each directory, it will
// run `CGO_ENABLED=0 GOARCH=amd64 GOOS=linux tinygo build -tags tinygo.enable`
// then attempt to fix-up the build tags by either adding or removing an
// go:build expression `(!tinygo || tinygo.enable)`

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const goBuild = "//go:build "
const constraint = "!tinygo || tinygo.enable"

type BuildStatus struct {
	passing  []string
	failing  []string
	excluded []string
}

// return trimmed output of "tinygo version"
func tinygoVersion(tinygo *string) string {
	out, err := exec.Command(*tinygo, "version").CombinedOutput()
	if nil != err {
		log.Fatalf("version: %v", err)
	}
	return strings.TrimSpace(string(out))
}

// check (via `go build -n`) if a given directory would have been skipped
// due to build constraints (e.g. cmds/core/bind only builds for plan9)
func isExcluded(dir string) bool {

	// to lazy to dynamically pull tags from `tinygo info`
	tags := []string{
		"tinygo",
		"tinygo.enable",
		"purego",
		"math_big_pure_go",
		"gc.precise",
		"scheduler.tasks",
		"serial.none",
	}
	c := exec.Command("go", "build",
		"-n",
		"-tags", strings.Join(tags, ","),
	)
	c.Env = append(os.Environ(), "GOOS=linux", "CGO_ENABLED=0", "GOARCH=amd64")
	c.Dir = dir
	out, _ := c.CombinedOutput()
	return strings.Contains(string(out), "build constraints exclude all Go files in")
}

// "tinygo build" in directory 'dir'
func build(tinygo *string, dir string) (err error) {
	c := exec.Command(*tinygo, "build", "-tags", "tinygo.enable")
	c.Dir = dir
	c.Stdout, c.Stderr = os.Stdout, os.Stderr
	c.Env = append(os.Environ(), "GOOS=linux", "CGO_ENABLED=0", "GOARCH=amd64")
	return c.Run()
}

// "tinygo build" in each of directories 'dirs'
func buildDirs(tinygo *string, dirs []string) (status BuildStatus, err error) {

	for _, dir := range dirs {
		log.Printf("Building %s\n", dir)
		if berr := build(tinygo, dir); berr != nil {
			if isExcluded(dir) {
				log.Printf("%v EXCLUDED\n", dir)
				status.excluded = append(status.excluded, dir)
				continue
			}
			log.Printf("%v FAILED %v\n", dir, berr)
			status.failing = append(status.failing, dir)
		} else {
			status.passing = append(status.passing, dir)
		}
	}
	return
}

func fixupConstraints(dir string, builds bool) (err error) {

	p := printer.Config{Mode: printer.UseSpaces | printer.TabIndent, Tabwidth: 8}

	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		log.Fatal(err)
	}
nextFile:
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

		goBuildPresent := false

	done:
		// modify existing //go:build line
		for _, cg := range f.Comments {
			for _, c := range cg.List {
				if !strings.HasPrefix(c.Text, goBuild) {
					continue
				}
				goBuildPresent = true

				contains := strings.Contains(c.Text, constraint)

				if (builds && !contains) || (!builds && contains) {
					log.Printf("Skipped, constraint up-to-date: %s\n", file)
					continue nextFile
				}

				if builds {
					re := regexp.MustCompile(`\(?\s*!tinygo\s+\|\|\s+tinygo.enable\s*\)?(\s+\&\&)?`)
					c.Text = re.ReplaceAllString(c.Text, "")
					log.Printf("Stripping build constraint %v\n", file)

					// handle potentially now-empty build constraint
					re = regexp.MustCompile(`^\s*//go:build\s*$`)
					if re.MatchString(c.Text) {
						c.Text = "//"
					}
				} else {
					c.Text = goBuild + "(" + constraint + ") && (" + c.Text[len(goBuild):] + ")"
				}
				break done
			}
		}

		if !builds && !goBuildPresent {
			// no //go:build line found: insert one
			var cg ast.CommentGroup
			cg.List = append(cg.List, &ast.Comment{Text: goBuild + constraint})

			if len(f.Comments) > 0 {
				// insert //go:build after first comment
				// group, assumed copyright. Doesn't seem
				// quite right but seems to work.
				cg.List[0].Slash = f.Comments[0].List[0].Slash + 1
				f.Comments = append([]*ast.CommentGroup{f.Comments[0], &cg}, f.Comments[1:]...)
			} else {
				// prepend //go:build
				f.Comments = append([]*ast.CommentGroup{&cg}, f.Comments...)
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
	return
}

func writeMarkdown(file *os.File, pathMD *string, tinygo *string, status BuildStatus) (err error) {

	fmt.Fprintf(file, "---\n\n")
	fmt.Fprintf(file, "DO NOT EDIT.\n\n")
	fmt.Fprintf(file, "Generated via `go run tools/tinygoize/main.go`\n\n")
	fmt.Fprintf(file, "%v\n\n", tinygoVersion(tinygo))
	fmt.Fprintf(file, "---\n\n")

	fmt.Fprintf(file, `# Status of u-root + tinygo
This document aims to track the process of enabling all u-root commands
to be built using tinygo. It will be updated as more commands can be built via:

    u-root> go run tools/tinygoize/main.go cmds/{core,exp,extra}/*

Commands that cannot be built with tinygo have a \"(!tinygo || tinygo.enable)\"
build constraint. Specify the "tinygo.enable" build tag to attempt to build
them.

    tinygo build -tags tinygo.enable cmds/core/ls

The list below is the result of building each command for Linux, x86_64.

The necessary additions to tinygo will be tracked in
[#2979](https://github.com/u-root/u-root/issues/2979).

---

## Commands Build Status
`)

	linkText := func(dir string) string {
		relPath, _ := filepath.Rel(filepath.Dir(*pathMD), dir)
		return fmt.Sprintf("[%v](%v)", dir, relPath)
	}

	processSet := func(header string, dirs []string) {

		fmt.Fprintf(file, "\n### %v (%v commands)\n", header, len(dirs))
		sort.Strings(dirs)

		if len(dirs) == 0 {
			fmt.Fprintf(file, "NONE\n")
		}
		for _, dir := range dirs {
			fmt.Fprintf(file, " - %v\n", linkText(dir))
		}

	}

	processSet("EXCLUDED", status.excluded)
	processSet("FAILING", status.failing)
	processSet("PASSING", status.passing)

	return
}

func main() {
	pathMD := flag.String("o", "tools/tinygoize/README.md", "Output file for markdown summary, '-' or '' for STDOUT")
	tinygo := flag.String("tinygo", "tinygo", "Path to tinygo")

	flag.Parse()

	var err error
	file := os.Stdout
	if len(*pathMD) > 0 && *pathMD != "-" {
		file, err = os.Create(*pathMD)
		if err != nil {
			fmt.Printf("Error creating opening file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
	}

	// generate list of commands that pass / fail / are excluded
	status, err := buildDirs(tinygo, flag.Args())
	if nil != err {
		log.Fatal(err)
	}

	// fix-up constraints in failing files
	for _, f := range status.failing {
		err = fixupConstraints(f, false)
		if nil != err {
			log.Fatal(err)
		}
	}

	// fix-up constraints in passing files
	for _, f := range status.passing {
		err = fixupConstraints(f, true)
		if nil != err {
			log.Fatal(err)
		}
	}

	// write markdown output
	err = writeMarkdown(file, pathMD, tinygo, status)
	if nil != err {
		log.Fatal(err)
	}
}
