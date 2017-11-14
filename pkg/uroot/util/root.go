// Copyright 2014-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux

// Package util contains various u-root utility functions.
package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"syscall"
)

const (
	// Not all these paths may be populated or even exist but OTOH they might.
	PATHHEAD = "/ubin"
	PATHMID  = "/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/bin:/usr/local/sbin"
	PATHTAIL = "/buildbin"
	CmdsPath = "github.com/u-root/u-root/cmds"
)

type Creator interface {
	Create() error
	fmt.Stringer
}

type Dir struct {
	Name string
	Mode os.FileMode
}

func (d Dir) Create() error {
	return os.MkdirAll(d.Name, d.Mode)
}

func (d Dir) String() string {
	return fmt.Sprintf("dir %q (mode %#o)", d.Name, d.Mode)
}

type File struct {
	Name     string
	Contents string
	Mode     os.FileMode
}

func (f File) Create() error {
	return ioutil.WriteFile(f.Name, []byte(f.Contents), f.Mode)
}

func (f File) String() string {
	return fmt.Sprintf("file %q (mode %#o)", f.Name, f.Mode)
}

type Symlink struct {
	Target  string
	NewPath string
}

func (s Symlink) Create() error {
	os.Remove(s.NewPath)
	return os.Symlink(s.Target, s.NewPath)
}

func (s Symlink) String() string {
	return fmt.Sprintf("symlink %q -> %q", s.NewPath, s.Target)
}

type Link struct {
	OldPath string
	NewPath string
}

func (s Link) Create() error {
	os.Remove(s.NewPath)
	return os.Link(s.OldPath, s.NewPath)
}

func (s Link) String() string {
	return fmt.Sprintf("link %q -> %q", s.NewPath, s.OldPath)
}

type Dev struct {
	Name string
	Mode uint32
	Dev  int
}

func (d Dev) Create() error {
	os.Remove(d.Name)
	return syscall.Mknod(d.Name, d.Mode, d.Dev)
}

func (d Dev) String() string {
	return fmt.Sprintf("dev %q (mode %#o; magic %d)", d.Name, d.Mode, d.Dev)
}

type Mount struct {
	Source string
	Target string
	FSType string
	Flags  uintptr
	Opts   string
}

func (m Mount) Create() error {
	return syscall.Mount(m.Source, m.Target, m.FSType, m.Flags, m.Opts)
}

func (m Mount) String() string {
	return fmt.Sprintf("mount -t %q -o %s %q %q flags %#x", m.FSType, m.Opts, m.Source, m.Target, m.Flags)
}

var (
	namespace = []Creator{
		Dir{Name: "/buildbin", Mode: 0777},
		Dir{Name: "/ubin", Mode: 0777},
		Dir{Name: "/tmp", Mode: 0777},
		Dir{Name: "/env", Mode: 0777},
		Dir{Name: "/tcz", Mode: 0777},
		Dir{Name: "/lib", Mode: 0777},
		Dir{Name: "/usr/lib", Mode: 0777},
		Dir{Name: "/go/pkg/linux_amd64", Mode: 0777},

		Dir{Name: "/etc", Mode: 0777},
		File{Name: "/etc/resolv.conf", Contents: `nameserver 8.8.8.8`, Mode: 0644},

		Dir{Name: "/proc", Mode: 0555},
		Mount{Target: "/proc", FSType: "proc"},
		Mount{Target: "/tmp", FSType: "tmpfs"},

		Dir{Name: "/dev", Mode: 0777},
		Dev{Name: "/dev/tty", Mode: syscall.S_IFCHR | 0666, Dev: 0x0500},
		Dev{Name: "/dev/urandom", Mode: syscall.S_IFCHR | 0444, Dev: 0x0109},
		Dev{Name: "/dev/port", Mode: syscall.S_IFCHR | 0640, Dev: 0x0104},

		// Kernel must be compiled with CONFIG_DEVTMPFS.
		// Note that things kind of work even if this mount fails.
		// TODO: move the Dir commands above below this line?
		Mount{Target: "/dev", FSType: "devtmpfs"},

		Dir{Name: "/dev/pts", Mode: 0777},
		Mount{Target: "/dev/pts", FSType: "devpts", Opts: "newinstance,ptmxmode=666,gid=5,mode=620"},
		Symlink{NewPath: "/dev/ptmx", Target: "/dev/pts/ptmx"},

		// Note: shm is required at least for Chrome. If you don't mount
		// it chrome throws a bogus "out of memory" error, not the more
		// useful "I can't open /dev/shm/whatever". SAD!
		Dir{Name: "/dev/shm", Mode: 0777},
		Mount{Source: "tmpfs", Target: "/dev/shm", FSType: "tmpfs"},

		Dir{Name: "/sys", Mode: 0555},
		Mount{Source: "sys", Target: "/sys", FSType: "sysfs"},
		Mount{Source: "cgroup", Target: "/sys/fs/cgroup", FSType: "tmpfs"},
		Dir{Name: "/sys/fs/cgroup/memory", Mode: 0555},
		Dir{Name: "/sys/fs/cgroup/freezer", Mode: 0555},
		Dir{Name: "/sys/fs/cgroup/devices", Mode: 0555},
		Dir{Name: "/sys/fs/cgroup/cpu,cpuacct", Mode: 0555},
		Symlink{NewPath: "/sys/fs/cgroup/cpu", Target: "/sys/fs/cgroup/cpu,cpuacct"},
		Symlink{NewPath: "/sys/fs/cgroup/cpuacct", Target: "/sys/fs/cgroup/cpu,cpuacct"},
		Mount{Source: "cgroup", Target: "/sys/fs/cgroup/memory", FSType: "cgroup", Opts: "memory"},
		Mount{Source: "cgroup", Target: "/sys/fs/cgroup/freezer", FSType: "cgroup", Opts: "freezer"},
		Mount{Source: "cgroup", Target: "/sys/fs/cgroup/devices", FSType: "cgroup", Opts: "devices"},
		Mount{Source: "cgroup", Target: "/sys/fs/cgroup/cpu,cpuacct", FSType: "cgroup", Opts: "cpu,cpuacct"},
	}

	Env = map[string]string{
		"LD_LIBRARY_PATH": "/usr/local/lib",
		"GOROOT":          "/go",
		"GOPATH":          "/",
		"GOBIN":           "/ubin",
		"CGO_ENABLED":     "0",
	}
)

func GoBin() string {
	return fmt.Sprintf("/go/bin/%s_%s:/go/bin:/go/pkg/tool/%s_%s", runtime.GOOS, runtime.GOARCH, runtime.GOOS, runtime.GOARCH)
}

// build the root file system.
func Rootfs() {
	Env["PATH"] = fmt.Sprintf("%v:%v:%v:%v", GoBin(), PATHHEAD, PATHMID, PATHTAIL)
	for k, v := range Env {
		os.Setenv(k, v)
	}

	for _, c := range namespace {
		if err := c.Create(); err != nil {
			log.Printf("Error creating %s: %v", c, err)
		} else {
			log.Printf("Created %v", c)
		}
	}
}
