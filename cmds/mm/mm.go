// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"os"
	"os/exec"

	"github.com/u-root/u-root/pkg/uroot/util"
)

var (
	dest      = flag.String("d", "/u", "destination directory")
	namespace = []util.Creator{
		util.Dir{Name: "proc", Mode: 0555},
		util.Dir{Name: "sys", Mode: 0555},
		util.Dir{Name: "buildbin", Mode: 0777},
		util.Dir{Name: "ubin", Mode: 0777},
		util.Dir{Name: "tmp", Mode: 0777},
		util.Dir{Name: "env", Mode: 0777},
		util.Dir{Name: "etc", Mode: 0777},
		util.Dir{Name: "tcz", Mode: 0777},
		util.Dir{Name: "dev", Mode: 0777},
		util.Dir{Name: "dev/pts", Mode: 0777},
		util.Dir{Name: "lib", Mode: 0777},
		util.Dir{Name: "usr/lib", Mode: 0777},
		util.Dir{Name: "go/pkg/linux_amd64", Mode: 0777},
		util.Link{NewPath: "init", OldPath: "/init"},
		util.Mount{Target: "dev/pts", FSType: "devpts", Opts: "newinstance,ptmxmode=666,gid=5,mode=620"},
		util.Symlink{NewPath: "dev/ptmx", Target: "/dev/pts/ptmx"},
	}
	commands = []*exec.Cmd{
		exec.Command("minimega", "-e", "vm", "config", "filesystem"),
		exec.Command("minimega", "-e", "vm", "config", "snapshot", "false"),
		exec.Command("minimega", "-e", "vm", "launch", "container", "uroot"),
		exec.Command("minimega", "-e", "vm", "start", "all"),
	}
)

func main() {
	flag.Parse()
	// We won't do wholesale removal. That's up to you if this fails.
	if err := os.Chdir(*dest); err == nil {
		log.Printf("Directory exists, skipping namespace setup")
	} else if !os.IsNotExist(err) {
		log.Fatalf("Couldn't chdir(%q): %v", *dest, err)
	} else {
		if err := os.Mkdir(*dest, 0777); err != nil {
			log.Fatalf("Can't mkdir: %v", err)
		}
		if err := os.Chdir(*dest); err != nil {
			log.Fatalf("Can't cd: %v", err)
		}
		for _, c := range namespace {
			if err := c.Create(); err != nil {
				log.Fatalf("Error creating %s: %v", c, err)
			}
		}
	}
	log.Printf("now you need to run stuff")
	commands[0].Args = append(commands[0].Args, *dest)
	for _, c := range commands {
		o, err := c.CombinedOutput()
		log.Printf("%v says %v", c, o)
		if err != nil {
			log.Fatalf("Fails with %v", err)
		}
	}
}
