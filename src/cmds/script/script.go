// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Script takes the arg list, does minimal rewriting, builds it and runs it
package main

import (
       "flag"
	"log"
	"os"
	"os/exec"

	"golang.org/x/tools/imports"
)

func main() {
     opts := imports.Options{
     Fragment: true,
     AllErrors: true,
     Comments: true,
     TabIndent: true,
     TabWidth: 8,
     }
	flag.Parse()
	a := "func main()"
	for _, v := range flag.Args() {
	    a = a + v
	    }
	   log.Printf("'%v'", a)
	goCode, err := imports.Process("commandline", []byte(a), &opts)
	if err != nil {
	   log.Fatalf("bad parse: '%v': %v", a, err)
	   }
	   log.Printf("%v", a)

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

}
