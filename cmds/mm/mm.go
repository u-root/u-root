// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/u-root/u-root/uroot"
)

var (
	dest      = flag.String("d", "/u", "destination directory")
	namespace = []uroot.Creator{
		uroot.Dir{Name: "proc", Mode: os.FileMode(0555)},
		uroot.Dir{Name: "sys", Mode: os.FileMode(0555)},
		uroot.Dir{Name: "buildbin", Mode: os.FileMode(0777)},
		uroot.Dir{Name: "ubin", Mode: os.FileMode(0777)},
		uroot.Dir{Name: "tmp", Mode: os.FileMode(0777)},
		uroot.Dir{Name: "env", Mode: os.FileMode(0777)},
		uroot.Dir{Name: "etc", Mode: os.FileMode(0777)},
		uroot.Dir{Name: "tcz", Mode: os.FileMode(0777)},
		uroot.Dir{Name: "dev", Mode: os.FileMode(0777)},
		uroot.Dir{Name: "dev/pts", Mode: os.FileMode(0777)},
		uroot.Dir{Name: "lib", Mode: os.FileMode(0777)},
		uroot.Dir{Name: "usr/lib", Mode: os.FileMode(0777)},
		uroot.Dir{Name: "go/pkg/linux_amd64", Mode: os.FileMode(0777)},
		uroot.Link{Oldpath: "/init", Newpath: "init"},
		//uroot.Dir{Name: "dev/null", Mode: uint32(syscall.S_IFCHR) | 0666, dev: 0x0103},
		//uroot.Dir{Name: "dev/console", Mode: uint32(syscall.S_IFCHR) | 0666, dev: 0x0501},
		//uroot.Dev{Name: "dev/tty", Mode: uint32(syscall.S_IFCHR) | 0666, dev: 0x0500},
		//uroot.Dev{Name: "dev/urandom", Mode: uint32(syscall.S_IFCHR) | 0444, dev: 0x0109},
		//mount{source: "proc", target: "proc", fstype: "proc", flags: syscall.MS_MGC_VAL, opts: ""},
		//mount{source: "sys", target: "sys", fstype: "sysfs", flags: syscall.MS_MGC_VAL, opts: ""},
		//// Kernel must be compiled with CONFIG_DEVTMPFS for this to work.
		//mount{source: "none", target: "dev", fstype: "devtmpfs", flags: syscall.MS_MGC_VAL},
		uroot.Mount{Source: "none", Target: "dev/pts", FSType: "devpts", Flags: syscall.MS_MGC_VAL, Opts: "newinstance,ptmxmode=666,gid=5,mode=620"},
		uroot.Symlink{Linkpath: "/dev/pts/ptmx", Target: "dev/ptmx"},
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
	// We won't do wholesale removal. That's up to you if this fails. */
	if err := os.Chdir(*dest); err == nil {
		log.Printf("Directory exists, skipping namespace setup")
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
