package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"syscall"
)

var (
	ipc     = flag.Bool("ipc", false, "Unshare the IPC namespace")
	mount   = flag.Bool("mount", false, "Unshared the mount namespace")
	pid     = flag.Bool("pid", false, "Unshared the pid namespace")
	net     = flag.Bool("net", false, "Unshared the net namespace")
	uts     = flag.Bool("uts", false, "Unshared the uts namespace")
	user    = flag.Bool("user", false, "Unshared the user namespace")
	maproot = flag.Bool("map-root-user", false, "map current uid to root. Not working")
)

func main() {
	flag.Parse()
	a := flag.Args()
	if len(a) == 0 {
		a = []string{"/bin/bash", "bash"}
	}
	c := exec.Command(a[0], a[1:]...)
	c.SysProcAttr = &syscall.SysProcAttr{}
	if *mount {
		c.SysProcAttr.Cloneflags |= syscall.CLONE_NEWNS
	}
	if *uts {
		c.SysProcAttr.Cloneflags |= syscall.CLONE_NEWUTS
	}
	if *ipc {
		c.SysProcAttr.Cloneflags |= syscall.CLONE_NEWIPC
	}
	if *net {
		c.SysProcAttr.Cloneflags |= syscall.CLONE_NEWNET
	}
	if *pid {
		c.SysProcAttr.Cloneflags |= syscall.CLONE_NEWPID
	}
	if *user {
		c.SysProcAttr.Cloneflags |= syscall.CLONE_NEWUSER
	}

	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		log.Printf(err.Error())
	}
}
