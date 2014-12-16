// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Script takes the arg list, does minimal rewriting, builds it and runs it
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"golang.org/x/tools/imports"
)

var (
	startPart = "package main\n"
	initPart  = "func init() {\n	addBuiltIn(\"%s\", b)\n}\nfunc b(cmd string, s []string) error {\nvar err error\n"
	//	endPart = "\n}\n)\n}\n"
	endPart = "\nreturn err\n}\n"
)

func main() {
	opts := imports.Options{
		Fragment:  true,
		AllErrors: true,
		Comments:  true,
		TabIndent: true,
		TabWidth:  8,
	}
	flag.Parse()
	goCode := startPart
	a := flag.Args()
	if len(a) < 3 {
		log.Fatalf("Usage: builtin <command> <code>")
	}
	// Simple programs are just bits of code for main ...
	if a[1] == "{" {
		goCode = goCode + fmt.Sprintf(initPart, a[0])
		for _, v := range a[2:] {
			if v == "}" {
				break
			}
			goCode = goCode + v + "\n"
		}
	} else {
		for _, v := range a[1:] {
			if v == "{" {
				goCode = goCode + fmt.Sprintf(initPart, a[0])
				continue
			}
			// FIXME: should only look for last arg.
			if v == "}" {
				break
			}
			goCode = goCode + v + "\n"
		}
	}
	goCode = goCode + endPart
	log.Printf("%v", goCode)
	fullCode, err := imports.Process("commandline", []byte(goCode), &opts)
	if err != nil {
		log.Fatalf("bad parse: '%v': %v", goCode, err)
	}
	log.Printf("%v", a)

	log.Print(fullCode)

	d, err := ioutil.TempDir("", "builtin")
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(path.Join(d, a[0]+".go"), []byte(fullCode), 0666); err != nil {
		log.Fatal(err)
	}

	/* copy all of /src/cmds/sh/*.go to the directory. */
	globs, err := filepath.Glob("/src/cmds/sh/*.go")
	if err != nil {
		log.Fatal(err)
	}
	for _, i := range globs {
		if b, err := ioutil.ReadFile(i); err != nil {
			log.Fatal(err)
		} else {
			_, df := path.Split(i)
			f := path.Join(d, df)
			if err = ioutil.WriteFile(f, b, 0600); err != nil {
				log.Fatal(err)
			}
		}
	}

	os.Setenv("GOBIN", d)
	cmd := exec.Command("go", "build", "-x", ".")
	cmd.Dir = d

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	log.Printf("Install %v", a[0])
	if err = cmd.Run(); err != nil {
		log.Printf("%v\n", err)
	}

	// stupid, but hey ...
	_, execName := path.Split(d)
	execName = path.Join(d, execName)
	os.Setenv("GOBIN", "/bin")
	cmd = exec.Command(execName)
	cmd.Dir = d

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	log.Printf("Run %v", execName)
	if err := cmd.Run(); err != nil {
		log.Printf("%v\n", err)
	}

}
