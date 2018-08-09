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
	"strconv"
	"syscall"

	"github.com/u-root/u-root/pkg/cmdline"
)

const (
	// Not all these paths may be populated or even exist but OTOH they might.
	PATHHEAD = "/ubin"
	PATHMID  = "/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/bin:/usr/local/sbin"
	PATHTAIL = "/buildbin:/bbin"
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
		Dir{Name: "/var/log", Mode: 0777},
		Dir{Name: "/go/pkg/linux_amd64", Mode: 0777},

		Dir{Name: "/etc", Mode: 0777},

		Dir{Name: "/proc", Mode: 0555},
		Mount{Source: "proc", Target: "/proc", FSType: "proc"},
		Mount{Source: "tmpfs", Target: "/tmp", FSType: "tmpfs"},

		Dir{Name: "/dev", Mode: 0777},
		Dev{Name: "/dev/tty", Mode: syscall.S_IFCHR | 0666, Dev: 0x0500},
		Dev{Name: "/dev/urandom", Mode: syscall.S_IFCHR | 0444, Dev: 0x0109},
		Dev{Name: "/dev/port", Mode: syscall.S_IFCHR | 0640, Dev: 0x0104},

		// Kernel must be compiled with CONFIG_DEVTMPFS.
		// Note that things kind of work even if this mount fails.
		// TODO: move the Dir commands above below this line?
		Mount{Source: "devtmpfs", Target: "/dev", FSType: "devtmpfs"},

		Dir{Name: "/dev/pts", Mode: 0777},
		Mount{Source: "devpts", Target: "/dev/pts", FSType: "devpts", Opts: "newinstance,ptmxmode=666,gid=5,mode=620"},
		Dev{Name: "/dev/ptmx", Mode: syscall.S_IFCHR | 0666, Dev: 0x0502},
		// Note: shm is required at least for Chrome. If you don't mount
		// it chrome throws a bogus "out of memory" error, not the more
		// useful "I can't open /dev/shm/whatever". SAD!
		Dir{Name: "/dev/shm", Mode: 0777},
		Mount{Source: "tmpfs", Target: "/dev/shm", FSType: "tmpfs"},

		Dir{Name: "/sys", Mode: 0555},
		Mount{Source: "sysfs", Target: "/sys", FSType: "sysfs"},
	}
	cgroupsnamespace = []Creator{
		Mount{Source: "cgroup", Target: "/sys/fs/cgroup", FSType: "tmpfs"},
		Dir{Name: "/sys/fs/cgroup/memory", Mode: 0555},
		Dir{Name: "/sys/fs/cgroup/freezer", Mode: 0555},
		Dir{Name: "/sys/fs/cgroup/devices", Mode: 0555},
		Dir{Name: "/sys/fs/cgroup/cpu,cpuacct", Mode: 0555},
		Dir{Name: "/sys/fs/cgroup/blkio", Mode: 0555},
		Dir{Name: "/sys/fs/cgroup/cpuset", Mode: 0555},
		Dir{Name: "/sys/fs/cgroup/pids", Mode: 0555},
		Dir{Name: "/sys/fs/cgroup/net_cls,net_prio", Mode: 0555},
		Dir{Name: "/sys/fs/cgroup/hugetlb", Mode: 0555},
		Dir{Name: "/sys/fs/cgroup/perf_event", Mode: 0555},
		Symlink{NewPath: "/sys/fs/cgroup/cpu", Target: "/sys/fs/cgroup/cpu,cpuacct"},
		Symlink{NewPath: "/sys/fs/cgroup/cpuacct", Target: "/sys/fs/cgroup/cpu,cpuacct"},
		Symlink{NewPath: "/sys/fs/cgroup/net_cls", Target: "/sys/fs/cgroup/net_cls,net_prio"},
		Symlink{NewPath: "/sys/fs/cgroup/net_prio", Target: "/sys/fs/cgroup/net_cls,net_prio"},
		Mount{Source: "cgroup", Target: "/sys/fs/cgroup/memory", FSType: "cgroup", Opts: "memory"},
		Mount{Source: "cgroup", Target: "/sys/fs/cgroup/freezer", FSType: "cgroup", Opts: "freezer"},
		Mount{Source: "cgroup", Target: "/sys/fs/cgroup/devices", FSType: "cgroup", Opts: "devices"},
		Mount{Source: "cgroup", Target: "/sys/fs/cgroup/cpu,cpuacct", FSType: "cgroup", Opts: "cpu,cpuacct"},
		Mount{Source: "cgroup", Target: "/sys/fs/cgroup/blkio", FSType: "cgroup", Opts: "blkio"},
		Mount{Source: "cgroup", Target: "/sys/fs/cgroup/cpuset", FSType: "cgroup", Opts: "cpuset"},
		Mount{Source: "cgroup", Target: "/sys/fs/cgroup/pids", FSType: "cgroup", Opts: "pids"},
		Mount{Source: "cgroup", Target: "/sys/fs/cgroup/net_cls,net_prio", FSType: "cgroup", Opts: "net_cls,net_prio"},
		Mount{Source: "cgroup", Target: "/sys/fs/cgroup/hugetlb", FSType: "cgroup", Opts: "hugetlb"},
		Mount{Source: "cgroup", Target: "/sys/fs/cgroup/perf_event", FSType: "cgroup", Opts: "perf_event"},
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

func create(namespace []Creator) {
	for _, c := range namespace {
		if err := c.Create(); err != nil {
			log.Printf("Error creating %s: %v", c, err)
		} else {
			log.Printf("Created %v", c)
		}
	}
}

// build the root file system.
func Rootfs() {
	Env["PATH"] = fmt.Sprintf("%v:%v:%v:%v", GoBin(), PATHHEAD, PATHMID, PATHTAIL)
	for k, v := range Env {
		os.Setenv(k, v)
	}
	create(namespace)

	// systemd gets upset when it discovers something has already setup cgroups
	// We have to do this after the base namespace is created, so we have /proc
	initFlags := cmdline.GetInitFlagMap()
	systemd, present := initFlags["systemd"]
	systemdEnabled, boolErr := strconv.ParseBool(systemd)
	if !present || boolErr != nil || systemdEnabled == false {
		create(cgroupsnamespace)
	}

}
