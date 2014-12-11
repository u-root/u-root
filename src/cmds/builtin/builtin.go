// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Script takes the arg list, does minimal rewriting, builds it and runs it
package main

import (
       "flag"
	"fmt"
	"log"
	"os"
_	"os/exec"
)
var (
	startPart = "package main\n"
	initPart = "func init() {\n	addBuiltIn(\"%s\", func %s(cmd string, s []string) error {\nvar err error\n"
	endPart = "\n}\n)\n}\n"
)

func main() {
	flag.Parse()
	goCode := startPart
	a := flag.Args()
	if len (a) < 3 {
		log.Fatalf("Usage: builtin <command> <code>")
	}
	// Simple programs are just bits of code for main ...
	if a[0] == "{" {
		goCode = goCode + fmt.Sprintf(initPart, a[0], a[0])
		for _, v := range a[1:] {
			if v == "}" {
				break
			}
			goCode = goCode + v + "\n"
		}
	} else {
		for _, v := range a[1:] {
			if v == "{" {
				goCode = goCode + fmt.Sprintf(initPart, a[0], a[0])
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
	log.Print(goCode)
	os.Exit(1)

/*
	f, err := TempFile("", "script%s.go")
	if err != nil {
		log.Fatalf("Script: opening TempFile: %v", err)
	}

	if _, err := f.Write([]byte(goCode)); err != nil {
		log.Fatalf("Script: Writing %v: %v", f, err)
	}
	if err := f.Close(); err != nil {
		log.Fatalf("Script: Closing %v: %v", f, err)
	}

	os.Setenv("GOBIN", "/tmp")
	cmd := exec.Command("go", "install", "-x", f.Name())
	//installenvs = append(envs, "GOBIN=/tmp")
	cmd.Dir = "/"

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	log.Printf("Install %v", f.Name())
	if err = cmd.Run(); err != nil {
		log.Printf("%v\n", err)
	}

	// stupid, but hey ...
	execName := f.Name()
	execName = execName[:len(execName)-3]
	cmd = exec.Command(execName)
	cmd.Dir = "/tmp"

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	log.Printf("Run %v", f.Name())
	if err := cmd.Run(); err != nil {
		log.Printf("%v\n", err)
	}
*/

}
