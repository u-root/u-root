// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"golang.org/x/tools/imports"
)

const (
	rushPath = "/src/github.com/u-root/u-root/cmds/rush"
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
	initPart  = "func init() {\n	addBuiltIn(\"%s\", %s)\n}\nfunc %s(c*Command) error {\nvar err error\n"
	//	endPart = "\n}\n)\n}\n"
	endPart   = "\nreturn err\n}\n"
	namespace = []mount{
		{source: "tmpfs", target: rushPath, fstype: "tmpfs", flags: syscall.MS_MGC_VAL, opts: ""},
		{source: "tmpfs", target: "/ubin", fstype: "tmpfs", flags: syscall.MS_MGC_VAL, opts: ""},
	}
	debug = flag.Bool("d", false, "Print debug info")
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
	a := flag.Args()
	if len(a) < 2 || len(a)%2 != 0 {
		log.Fatalf("Usage: builtin <command> <code> [<command> <code>]*")
	}
	filemap := make(map[string][]byte)
	for ; len(a) > 0; a = a[2:] {
		goCode := startPart
		// Simple programs are just bits of code for main ...
		if a[1][0] == '{' {
			goCode = goCode + fmt.Sprintf(initPart, a[0], a[0], a[0])
			goCode = goCode + a[1][1:len(a[1])-1]
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
		if *debug {
			log.Printf("\n---------------------\n%v\n------------------------\n", goCode)
		}
		fullCode, err := imports.Process("commandline", []byte(goCode), &opts)
		if err != nil {
			log.Fatalf("bad parse: '%v': %v", goCode, err)
		}
		if *debug {
			log.Printf("\n----FULLCODE---------\n%v\n------FULLCODE----------\n", string(fullCode))
		}
		bName := filepath.Join(rushPath, a[0]+".go")
		filemap[bName] = fullCode
	}

	// processed code, read in shell files.
	globs, err := filepath.Glob(rushPath + "/*.go")
	if err != nil {
		log.Fatal(err)
	}
	for _, i := range globs {
		if b, err := ioutil.ReadFile(i); err != nil {
			log.Fatal(err)
		} else {
			if _, ok := filemap[i]; ok {
				log.Fatalf("%v exists", i)
			}
			filemap[i] = b
		}
	}

	if b, err := ioutil.ReadFile("/proc/mounts"); err == nil && false {
		log.Printf("m %v\n", string(b))
	}
	// we'd like to do this here, but it seems it doesn't end
	// up applying to all procs in this group, leading to confusion.
	// sometimes they get the private mount, sometimes not.
	// It's a fundamental limit in the go runtime.
	// So we hack it in the shell.
	// There is no FIXME
	if false {
		if err := syscall.Unshare(syscall.CLONE_NEWNS); err != nil {
			log.Fatal(err)
		}
	}
	if *debug {
		if b, err := ioutil.ReadFile("/proc/mounts"); err == nil {
			log.Printf("Reading /proc/mount:m %v\n", b)
		}
	}

	// We are rewriting the shell. We need to create a new binary, i.e.
	// rewrite the one in /ubin. Sadly, there is no way to say "mount THIS bin
	// before THAT bin". There will be ca. 3.18 and we might as well wait for
	// that to become common. For now, we essentially erase /ubin but mounting
	// a tmpfs over it.
	// This would be infinitely easier with a true union file system. Oh well.
	for _, m := range namespace {
		if err := syscall.Mount(m.source, m.target, m.fstype, m.flags, m.opts); err != nil {
			log.Printf("Mount :%s: on :%s: type :%s: flags %x: opts %v: %v\n", m.source, m.target, m.fstype, m.flags, m.opts, err)
		}
	}
	// write the new rushPath
	for i, v := range filemap {
		if err = ioutil.WriteFile(i, v, 0600); err != nil {
			log.Fatal(err)
		}
	}

	// the big fun: just run it. The Right Things Happen.
	cmd := exec.Command("/buildbin/rush")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	// TODO: figure out why we get EPERM when we use this.
	//cmd.SysProcAttr = &syscall.SysProcAttr{Setctty: true, Setsid: true,}
	if *debug {
		log.Printf("Run %v", cmd)
	}
	if err := cmd.Run(); err != nil {
		log.Printf("%v\n", err)
	}
	// Unshare doesn't work in a sane way due to a Go issue?
	for _, m := range namespace {
		if err := syscall.Unmount(m.target, syscall.MNT_FORCE); err != nil {
			log.Printf("Umount :%s: %v\n", m.target, err)
		}
	}
	log.Printf("builtin: /ubin/rush returned!\n")
}
