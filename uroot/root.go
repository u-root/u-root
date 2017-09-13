// Copyright 2014-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux

// package uroot contains various functions that might be needed more than
// one place.
package uroot

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
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
	return fmt.Sprintf("dir :%q: mode %o", d.Name, d.Mode)
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
	return fmt.Sprintf("file %q", f.Name)
}

type Symlink struct {
	Target   string
	Linkpath string
}

func (s Symlink) Create() error {
	os.Remove(s.Target)
	return os.Symlink(s.Linkpath, s.Target)
}

func (s Symlink) String() string {
	return fmt.Sprintf("symlink %q -> %q", s.Target, s.Linkpath)
}

type Link struct {
	Oldpath string
	Newpath string
}

func (s Link) Create() error {
	os.Remove(s.Newpath)
	return os.Link(s.Oldpath, s.Newpath)
}

func (s Link) String() string {
	return fmt.Sprintf("link %q -> %q", s.Oldpath, s.Newpath)
}

type Dev struct {
	Name    string
	Mode    uint32
	Dev     int
	Howmany int
}

func (d Dev) Create() error {
	os.Remove(d.Name)
	return syscall.Mknod(d.Name, d.Mode, d.Dev)
}

func (d Dev) String() string {
	return fmt.Sprintf("dev :%q: mode: %#o: magic: %v", d.Name, d.Mode, d.Dev)
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
	Profile string
	Envs    []string
	env     = map[string]string{
		"LD_LIBRARY_PATH": "/usr/local/lib",
		"GOROOT":          "/go",
		"GOPATH":          "/",
		"GOBIN":           "/ubin",
		"CGO_ENABLED":     "0",
	}

	namespace = []Creator{
		Dir{Name: "/proc", Mode: os.FileMode(0555)},
		Dir{Name: "/sys", Mode: os.FileMode(0555)},
		Dir{Name: "/buildbin", Mode: os.FileMode(0777)},
		Dir{Name: "/ubin", Mode: os.FileMode(0777)},
		Dir{Name: "/tmp", Mode: os.FileMode(0777)},
		Dir{Name: "/env", Mode: os.FileMode(0777)},
		Dir{Name: "/etc", Mode: os.FileMode(0777)},
		Dir{Name: "/tcz", Mode: os.FileMode(0777)},
		Dir{Name: "/dev", Mode: os.FileMode(0777)},
		Dir{Name: "/lib", Mode: os.FileMode(0777)},
		Dir{Name: "/usr/lib", Mode: os.FileMode(0777)},
		Dir{Name: "/go/pkg/linux_amd64", Mode: os.FileMode(0777)},
		// chicken and egg: these need to be there before you start and hence
		// built into the initial initramfs.
		//{Name: "/dev/null", Mode: uint32(syscall.S_IFCHR) | 0666, dev: 0x0103},
		//{Name: "/dev/console", Mode: uint32(syscall.S_IFCHR) | 0666, dev: 0x0501},
		Dev{Name: "/dev/tty", Mode: uint32(syscall.S_IFCHR) | 0666, Dev: 0x0500},
		Dev{Name: "/dev/urandom", Mode: uint32(syscall.S_IFCHR) | 0444, Dev: 0x0109},
		Dev{Name: "/dev/port", Mode: uint32(syscall.S_IFCHR) | 0640, Dev: 0x0104},
		Mount{Source: "proc", Target: "/proc", FSType: "proc", Flags: syscall.MS_MGC_VAL, Opts: ""},
		Mount{Source: "sys", Target: "/sys", FSType: "sysfs", Flags: syscall.MS_MGC_VAL, Opts: ""},
		Mount{Source: "cgroup", Target: "/sys/fs/cgroup", FSType: "tmpfs", Flags: syscall.MS_MGC_VAL, Opts: ""},
		Mount{Source: "none", Target: "/tmp", FSType: "tmpfs", Flags: syscall.MS_MGC_VAL, Opts: ""},
		// Kernel must be compiled with CONFIG_DEVTMPFS, otherwise
		// default to contents of Dev.cpio.
		Mount{Source: "none", Target: "/dev", FSType: "devtmpfs", Flags: syscall.MS_MGC_VAL},
		Dir{Name: "/dev/pts", Mode: os.FileMode(0777)},
		Mount{Source: "none", Target: "/dev/pts", FSType: "devpts", Flags: syscall.MS_MGC_VAL, Opts: "newinstance,ptmxmode=666,gid=5,mode=620"},
		Symlink{Linkpath: "/dev/pts/ptmx", Target: "/dev/ptmx"},
		File{Name: "/etc/resolv.conf", Contents: `nameserver 8.8.8.8`, Mode: os.FileMode(0644)},
		Dir{Name: "/sys/fs/cgroup/memory", Mode: os.FileMode(0555)},
		Dir{Name: "/sys/fs/cgroup/freezer", Mode: os.FileMode(0555)},
		Dir{Name: "/sys/fs/cgroup/devices", Mode: os.FileMode(0555)},
		Dir{Name: "/sys/fs/cgroup/cpu,cpuacct", Mode: os.FileMode(0555)},
		Symlink{Linkpath: "/sys/fs/cgroup/cpu,cpuacct", Target: "/sys/fs/cgroup/cpu"},
		Symlink{Linkpath: "/sys/fs/cgroup/cpu,cpuacct", Target: "/sys/fs/cgroup/cpuacct"},
		Mount{Source: "cgroup", Target: "/sys/fs/cgroup/memory", FSType: "cgroup", Flags: syscall.MS_MGC_VAL, Opts: "memory"},
		Mount{Source: "cgroup", Target: "/sys/fs/cgroup/freezer", FSType: "cgroup", Flags: syscall.MS_MGC_VAL, Opts: "freezer"},
		Mount{Source: "cgroup", Target: "/sys/fs/cgroup/devices", FSType: "cgroup", Flags: syscall.MS_MGC_VAL, Opts: "devices"},
		Mount{Source: "cgroup", Target: "/sys/fs/cgroup/cpu,cpuacct", FSType: "cgroup", Flags: syscall.MS_MGC_VAL, Opts: "cpu,cpuacct"},
	}
)

// build the root file system.
func Rootfs() {
	// Pick some reasonable values in the (unlikely!) even that Uname fails.
	uname := "linux"
	mach := "amd64"
	// There are three possible places for go:
	// The first is in /go/bin/$OS_$ARCH
	// The second is in /go/bin [why they still use this path is anyone's guess]
	// The third is in /go/pkg/tool/$OS_$ARCH
	if u, err := Uname(); err != nil {
		log.Printf("uroot.Utsname fails: %v, so assume %v_%v\n", err, uname, mach)
	} else {
		// Sadly, go and the OS disagree on many things.
		uname = strings.ToLower(u.Sysname)
		mach = strings.ToLower(u.Machine)
		// Yes, we really have to do this stupid thing.
		if mach[0:3] == "arm" {
			mach = "arm"
		}
		if mach == "x86_64" {
			mach = "amd64"
		}
	}
	goPath := fmt.Sprintf("/go/bin/%s_%s:/go/bin:/go/pkg/tool/%s_%s", uname, mach, uname, mach)
	env["PATH"] = fmt.Sprintf("%v:%v:%v:%v", goPath, PATHHEAD, PATHMID, PATHTAIL)

	for k, v := range env {
		os.Setenv(k, v)
		Envs = append(Envs, k+"="+v)
	}

	// Some systems wipe out all the environment variables we so carefully craft.
	// There is a way out -- we can put them into /etc/profile.d/uroot if we want.
	// The PATH variable has to change, however.
	env["PATH"] = fmt.Sprintf("%v:%v:%v:%v", goPath, PATHHEAD, "$PATH", PATHTAIL)
	for k, v := range env {
		Profile += "export " + k + "=" + v + "\n"
	}
	// The IFS lets us force a rehash every time we type a command, so that when we
	// build uroot commands we don't keep rebuilding them.
	Profile += "IFS=`hash -r`\n"
	// IF the profile is used, THEN when the user logs in they will need a private
	// tmpfs. There's no good way to do this on linux. The closest we can get for now
	// is to mount a tmpfs of /go/pkg/%s_%s :-(
	// Same applies to ubin. Each user should have their own.
	Profile += fmt.Sprintf("sudo mount -t tmpfs none /go/pkg/%s_%s\n", uname, mach)
	Profile += fmt.Sprintf("sudo mount -t tmpfs none /ubin\n")
	Profile += fmt.Sprintf("sudo mount -t tmpfs none /pkg\n")

	for _, c := range namespace {
		if err := c.Create(); err != nil {
			log.Printf("Error creating %s: %v", c, err)
		} else {
			log.Printf("Created %v", c)
		}
	}

	// only in case of emergency.
	if false {
		if err := filepath.Walk("/", func(name string, fi os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf(" WALK FAIL%v: %v\n", name, err)
				// That's ok, sometimes things are not there.
				return nil
			}
			fmt.Printf("%v\n", name)
			return nil
		}); err != nil {
			log.Printf("WALK fails %v\n", err)
		}
	}
}
