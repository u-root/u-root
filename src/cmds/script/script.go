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
)

func main() {
	flag.Parse()
	goCode := "package main\n"
	a := flag.Args()
	// Simple programs are just bits of code for main ...
	if a[0] == "{" {
		goCode = goCode + "\nfunc main() {"
		for _, v := range a[1:] {
			if v == "}" {
				break
			}
			goCode = goCode + v
		}
	} else {
		for _, v := range a {
			if v == "{" {
				goCode = goCode + "\nfunc main() {\n"
				continue
			}
			if v == "}" {
				break
			}
			goCode = goCode + v
		}
	}
	goCode = goCode + "\n}\n"

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
