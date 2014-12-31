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
	"syscall"

	"golang.org/x/tools/imports"
)

type mount struct {
	source string
	target string
	fstype string
	flags  uintptr
	opts   string
}

var (
	startPart = "package main\n"
	initPart  = "func init() {\n	addBuiltIn(\"%s\", b)\n}\nfunc b(c*Command) error {\nvar err error\n"
	//	endPart = "\n}\n)\n}\n"
	endPart   = "\nreturn err\n}\n"
	namespace = []mount{
		{source: "tmpfs", target: "/src/cmds/sh", fstype: "tmpfs", flags: syscall.MS_MGC_VAL, opts: ""},
		{source: "tmpfs", target: "/bin", fstype: "tmpfs", flags: syscall.MS_MGC_VAL, opts: ""},
	}
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
	bName := path.Join("/src/cmds/sh", a[0]+".go")
	filemap := make(map[string][]byte)
	fmt.Printf("filemap %v\n", filemap)
	filemap[bName] = fullCode

	// processed code, read in shell files.
	globs, err := filepath.Glob("/src/cmds/sh/*.go")
	if err != nil {
		log.Fatal(err)
	}
	for _, i := range globs {
		if b, err := ioutil.ReadFile(i); err != nil {
			log.Fatal(err)
		} else {
			if _, ok := filemap[i]; ok {
				log.Fatal("%v exists", i)
			}
			filemap[i] = b
		}
	}

	log.Printf("%v", a)

	log.Print(fullCode)
	if b, err := ioutil.ReadFile("/proc/mounts"); err == nil {
		fmt.Printf("m %v\n", b)
	}
	// we'd like to do this here, but it seems it doesn't end
	// up applying to all procs in this group, leading to confusion. 
	// sometimes they get the private mount, sometimes not. So we had
	// to hack it in the shell. 
	// FIXME
	if false {
		if err := syscall.Unshare(syscall.CLONE_NEWNS); err != nil {
			log.Fatal(err)
		}
	}
	if b, err := ioutil.ReadFile("/proc/mounts"); err == nil {
		fmt.Printf("m %v\n", b)
	}

	// We are rewriting the shell. We need to create a new binary, i.e.
	// rewrite the one in /bin. Sadly, there is no way to say "mount THIS bin
	// before THAT bin". There will be ca. 3.18 and we might as well wait for
	// that to become common. For now, we essentially erase /bin but mounting
	// a tmpfs over it.
	// This would be infinitely easier with a true union file system. Oh well.
	for _, m := range namespace {
		if err := syscall.Mount(m.source, m.target, m.fstype, m.flags, m.opts); err != nil {
			log.Printf("Mount :%s: on :%s: type :%s: flags %x: %v\n", m.source, m.target, m.fstype, m.flags, m.opts, err)
		}
	}
	log.Printf("filemap: %v", filemap)
	// write the new /src/cmds/sh
	for i, v := range filemap {
		if err = ioutil.WriteFile(i, v, 0600); err != nil {
			log.Fatal(err)
		}
	}

	// the big fun: just run it. The Right Things Happen.
	cmd := exec.Command("/buildbin/sh")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	// TODO: figure out why we get EPERM when we use this.
	//cmd.SysProcAttr = &syscall.SysProcAttr{Setctty: true, Setsid: true,}
	log.Printf("Run %v", cmd)
	err = cmd.Run()
	if err != nil {
		log.Printf("%v\n", err)
	}
	// Unshare doesn't work in a sane way due to a Go issue?
	for _, m := range namespace {
		if err := syscall.Unmount(m.target, syscall.MNT_FORCE); err != nil {
			log.Printf("Umount :%s: %v\n", m.target, err)
		}
	}
	log.Printf("init: /bin/sh returned!\n")
}
