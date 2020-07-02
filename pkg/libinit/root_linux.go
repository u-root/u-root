// Copyright 2014-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package libinit creates the environment and root file system for u-root.
package libinit

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"syscall"

	"github.com/u-root/u-root/pkg/cmdline"
	"github.com/u-root/u-root/pkg/ulog"
	"golang.org/x/sys/unix"
)

type creator interface {
	create() error
	fmt.Stringer
}

type dir struct {
	Name string
	Mode os.FileMode
}

func (d dir) create() error {
	return os.MkdirAll(d.Name, d.Mode)
}

func (d dir) String() string {
	return fmt.Sprintf("dir %q (mode %#o)", d.Name, d.Mode)
}

type symlink struct {
	Target  string
	NewPath string
}

func (s symlink) create() error {
	os.Remove(s.NewPath)
	return os.Symlink(s.Target, s.NewPath)
}

func (s symlink) String() string {
	return fmt.Sprintf("symlink %q -> %q", s.NewPath, s.Target)
}

type dev struct {
	Name string
	Mode uint32
	Dev  int
}

func (d dev) create() error {
	os.Remove(d.Name)
	return syscall.Mknod(d.Name, d.Mode, d.Dev)
}

func (d dev) String() string {
	return fmt.Sprintf("dev %q (mode %#o; magic %d)", d.Name, d.Mode, d.Dev)
}

type mount struct {
	Source string
	Target string
	FSType string
	Flags  uintptr
	Opts   string
}

func (m mount) create() error {
	return syscall.Mount(m.Source, m.Target, m.FSType, m.Flags, m.Opts)
}

func (m mount) String() string {
	return fmt.Sprintf("mount -t %q -o %s %q %q flags %#x", m.FSType, m.Opts, m.Source, m.Target, m.Flags)
}

var (
	// These have to be created / mounted first, so that the logging works correctly.
	preNamespace = []creator{
		dir{Name: "/dev", Mode: 0777},

		// Kernel must be compiled with CONFIG_DEVTMPFS.
		mount{Source: "devtmpfs", Target: "/dev", FSType: "devtmpfs"},
	}
	namespace = []creator{
		dir{Name: "/buildbin", Mode: 0777},
		dir{Name: "/ubin", Mode: 0777},
		dir{Name: "/tmp", Mode: 0777},
		dir{Name: "/env", Mode: 0777},
		dir{Name: "/tcz", Mode: 0777},
		dir{Name: "/lib", Mode: 0777},
		dir{Name: "/usr/lib", Mode: 0777},
		dir{Name: "/var/log", Mode: 0777},
		dir{Name: "/go/pkg/linux_amd64", Mode: 0777},

		dir{Name: "/etc", Mode: 0777},

		dir{Name: "/proc", Mode: 0555},
		mount{Source: "proc", Target: "/proc", FSType: "proc"},
		mount{Source: "tmpfs", Target: "/tmp", FSType: "tmpfs"},

		dev{Name: "/dev/tty", Mode: syscall.S_IFCHR | 0666, Dev: 0x0500},
		dev{Name: "/dev/urandom", Mode: syscall.S_IFCHR | 0444, Dev: 0x0109},
		dev{Name: "/dev/port", Mode: syscall.S_IFCHR | 0640, Dev: 0x0104},

		dir{Name: "/dev/pts", Mode: 0777},
		mount{Source: "devpts", Target: "/dev/pts", FSType: "devpts", Opts: "newinstance,ptmxmode=666,gid=5,mode=620"},
		// Note: if we mount /dev/pts with "newinstance", we *must* make "/dev/ptmx" a symlink to "/dev/pts/ptmx"
		symlink{NewPath: "/dev/ptmx", Target: "/dev/pts/ptmx"},
		// Note: shm is required at least for Chrome. If you don't mount
		// it chrome throws a bogus "out of memory" error, not the more
		// useful "I can't open /dev/shm/whatever". SAD!
		dir{Name: "/dev/shm", Mode: 0777},
		mount{Source: "tmpfs", Target: "/dev/shm", FSType: "tmpfs"},

		dir{Name: "/sys", Mode: 0555},
		mount{Source: "sysfs", Target: "/sys", FSType: "sysfs"},
		mount{Source: "securityfs", Target: "/sys/kernel/security", FSType: "securityfs"},
	}

	// cgroups are optional for most u-root users, especially
	// LinuxBoot/NERF. Some users use u-root for container stuff.
	cgroupsnamespace = []creator{
		mount{Source: "cgroup", Target: "/sys/fs/cgroup", FSType: "tmpfs"},
		dir{Name: "/sys/fs/cgroup/memory", Mode: 0555},
		dir{Name: "/sys/fs/cgroup/freezer", Mode: 0555},
		dir{Name: "/sys/fs/cgroup/devices", Mode: 0555},
		dir{Name: "/sys/fs/cgroup/cpu,cpuacct", Mode: 0555},
		dir{Name: "/sys/fs/cgroup/blkio", Mode: 0555},
		dir{Name: "/sys/fs/cgroup/cpuset", Mode: 0555},
		dir{Name: "/sys/fs/cgroup/pids", Mode: 0555},
		dir{Name: "/sys/fs/cgroup/net_cls,net_prio", Mode: 0555},
		dir{Name: "/sys/fs/cgroup/hugetlb", Mode: 0555},
		dir{Name: "/sys/fs/cgroup/perf_event", Mode: 0555},
		symlink{NewPath: "/sys/fs/cgroup/cpu", Target: "/sys/fs/cgroup/cpu,cpuacct"},
		symlink{NewPath: "/sys/fs/cgroup/cpuacct", Target: "/sys/fs/cgroup/cpu,cpuacct"},
		symlink{NewPath: "/sys/fs/cgroup/net_cls", Target: "/sys/fs/cgroup/net_cls,net_prio"},
		symlink{NewPath: "/sys/fs/cgroup/net_prio", Target: "/sys/fs/cgroup/net_cls,net_prio"},
		mount{Source: "cgroup", Target: "/sys/fs/cgroup/memory", FSType: "cgroup", Opts: "memory"},
		mount{Source: "cgroup", Target: "/sys/fs/cgroup/freezer", FSType: "cgroup", Opts: "freezer"},
		mount{Source: "cgroup", Target: "/sys/fs/cgroup/devices", FSType: "cgroup", Opts: "devices"},
		mount{Source: "cgroup", Target: "/sys/fs/cgroup/cpu,cpuacct", FSType: "cgroup", Opts: "cpu,cpuacct"},
		mount{Source: "cgroup", Target: "/sys/fs/cgroup/blkio", FSType: "cgroup", Opts: "blkio"},
		mount{Source: "cgroup", Target: "/sys/fs/cgroup/cpuset", FSType: "cgroup", Opts: "cpuset"},
		mount{Source: "cgroup", Target: "/sys/fs/cgroup/pids", FSType: "cgroup", Opts: "pids"},
		mount{Source: "cgroup", Target: "/sys/fs/cgroup/net_cls,net_prio", FSType: "cgroup", Opts: "net_cls,net_prio"},
		mount{Source: "cgroup", Target: "/sys/fs/cgroup/hugetlb", FSType: "cgroup", Opts: "hugetlb"},
		mount{Source: "cgroup", Target: "/sys/fs/cgroup/perf_event", FSType: "cgroup", Opts: "perf_event"},
	}
)

func goBin() string {
	return fmt.Sprintf("/go/bin/%s_%s:/go/bin:/go/pkg/tool/%s_%s", runtime.GOOS, runtime.GOARCH, runtime.GOOS, runtime.GOARCH)
}

func create(namespace []creator, optional bool) {
	// Clear umask bits so that we get stuff like ptmx right.
	m := unix.Umask(0)
	defer unix.Umask(m)
	for _, c := range namespace {
		if err := c.create(); err != nil {
			if optional {
				ulog.KernelLog.Printf("u-root init [optional]: warning creating %s: %v", c, err)
			} else {
				ulog.KernelLog.Printf("u-root init: error creating %s: %v", c, err)
			}
		}
	}
}

// SetEnv sets the default u-root environment.
func SetEnv() {
	env := map[string]string{
		"LD_LIBRARY_PATH": "/usr/local/lib",
		"GOROOT":          "/go",
		"GOPATH":          "/",
		"GOBIN":           "/ubin",
		"CGO_ENABLED":     "0",
		"USER":            "root",
	}

	// Not all these paths may be populated or even exist but OTOH they might.
	path := "/ubin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/bin:/usr/local/sbin:/buildbin:/bbin"

	env["PATH"] = fmt.Sprintf("%v:%v", goBin(), path)
	for k, v := range env {
		os.Setenv(k, v)
	}
}

// CreateRootfs creates the default u-root file system.
func CreateRootfs() {
	// Mount devtmpfs, then open /dev/kmsg with Reinit.
	create(preNamespace, false)
	ulog.KernelLog.Reinit()

	create(namespace, false)

	// systemd gets upset when it discovers something has already setup cgroups
	// We have to do this after the base namespace is created, so we have /proc
	initFlags := cmdline.GetInitFlagMap()
	systemd, present := initFlags["systemd"]
	systemdEnabled, boolErr := strconv.ParseBool(systemd)
	if !present || boolErr != nil || !systemdEnabled {
		create(cgroupsnamespace, true)
	}
}
