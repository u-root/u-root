// Copyright 2014-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package libinit creates the environment and root file system for u-root.
package libinit

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/u-root/u-root/pkg/cmdline"
	"github.com/u-root/u-root/pkg/cp"
	"github.com/u-root/u-root/pkg/kmodule"
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
	return unix.Mknod(d.Name, d.Mode, d.Dev)
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
	return unix.Mount(m.Source, m.Target, m.FSType, m.Flags, m.Opts)
}

func (m mount) String() string {
	return fmt.Sprintf("mount -t %q -o %s %q %q flags %#x", m.FSType, m.Opts, m.Source, m.Target, m.Flags)
}

type cpdir struct {
	Source string
	Target string
}

func (c cpdir) create() error {
	return cp.CopyTree(c.Source, c.Target)
}

func (c cpdir) String() string {
	return fmt.Sprintf("cp -a %q %q", c.Source, c.Target)
}

var (
	// These have to be created / mounted first, so that the logging works correctly.
	preNamespace = []creator{
		dir{Name: "/dev", Mode: 0o777},

		// Kernel must be compiled with CONFIG_DEVTMPFS.
		mount{Source: "devtmpfs", Target: "/dev", FSType: "devtmpfs"},
	}
	namespace = []creator{
		dir{Name: "/buildbin", Mode: 0o777},
		dir{Name: "/ubin", Mode: 0o777},
		dir{Name: "/tmp", Mode: 0o777},
		dir{Name: "/env", Mode: 0o777},
		dir{Name: "/tcz", Mode: 0o777},
		dir{Name: "/lib", Mode: 0o777},
		dir{Name: "/usr/lib", Mode: 0o777},
		dir{Name: "/var/log", Mode: 0o777},
		dir{Name: "/go/pkg/linux_amd64", Mode: 0o777},

		dir{Name: "/etc", Mode: 0o777},

		dir{Name: "/proc", Mode: 0o555},
		mount{Source: "proc", Target: "/proc", FSType: "proc"},
		mount{Source: "tmpfs", Target: "/tmp", FSType: "tmpfs"},

		dev{Name: "/dev/tty", Mode: unix.S_IFCHR | 0o666, Dev: 0x0500},
		dev{Name: "/dev/urandom", Mode: unix.S_IFCHR | 0o444, Dev: 0x0109},
		dev{Name: "/dev/port", Mode: unix.S_IFCHR | 0o640, Dev: 0x0104},

		dir{Name: "/dev/pts", Mode: 0o777},
		mount{Source: "devpts", Target: "/dev/pts", FSType: "devpts", Opts: "newinstance,ptmxmode=666,gid=5,mode=620"},
		// Note: if we mount /dev/pts with "newinstance", we *must* make "/dev/ptmx" a symlink to "/dev/pts/ptmx"
		symlink{NewPath: "/dev/ptmx", Target: "/dev/pts/ptmx"},
		// Note: shm is required at least for Chrome. If you don't mount
		// it chrome throws a bogus "out of memory" error, not the more
		// useful "I can't open /dev/shm/whatever". SAD!
		dir{Name: "/dev/shm", Mode: 0o777},
		mount{Source: "tmpfs", Target: "/dev/shm", FSType: "tmpfs"},

		dir{Name: "/sys", Mode: 0o555},
		mount{Source: "sysfs", Target: "/sys", FSType: "sysfs"},
		mount{Source: "securityfs", Target: "/sys/kernel/security", FSType: "securityfs"},
		mount{Source: "efivarfs", Target: "/sys/firmware/efi/efivars", FSType: "efivarfs"},

		cpdir{Source: "/etc", Target: "/tmp/etc"},
		mount{Source: "/tmp/etc", Target: "/etc", FSType: "tmpfs", Flags: unix.MS_BIND},
	}

	// cgroups are optional for most u-root users, especially
	// LinuxBoot/NERF. Some users use u-root for container stuff.
	cgroupsnamespace = []creator{
		mount{Source: "cgroup", Target: "/sys/fs/cgroup", FSType: "tmpfs"},
		dir{Name: "/sys/fs/cgroup/memory", Mode: 0o555},
		dir{Name: "/sys/fs/cgroup/freezer", Mode: 0o555},
		dir{Name: "/sys/fs/cgroup/devices", Mode: 0o555},
		dir{Name: "/sys/fs/cgroup/cpu,cpuacct", Mode: 0o555},
		dir{Name: "/sys/fs/cgroup/blkio", Mode: 0o555},
		dir{Name: "/sys/fs/cgroup/cpuset", Mode: 0o555},
		dir{Name: "/sys/fs/cgroup/pids", Mode: 0o555},
		dir{Name: "/sys/fs/cgroup/net_cls,net_prio", Mode: 0o555},
		dir{Name: "/sys/fs/cgroup/hugetlb", Mode: 0o555},
		dir{Name: "/sys/fs/cgroup/perf_event", Mode: 0o555},
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

// InitModuleLoader wraps the resources we need for early module loading
type InitModuleLoader struct {
	Cmdline      *cmdline.CmdLine
	Prober       func(name string, modParameters string) error
	ExcludedMods map[string]bool
}

func (i *InitModuleLoader) IsExcluded(mod string) bool {
	return i.ExcludedMods[mod]
}

func (i *InitModuleLoader) LoadModule(mod string) error {
	flags := i.Cmdline.FlagsForModule(mod)
	if err := i.Prober(mod, flags); err != nil {
		return fmt.Errorf("failed to load module: %s", err)
	}
	return nil
}

func NewInitModuleLoader() *InitModuleLoader {
	return &InitModuleLoader{
		Cmdline: cmdline.NewCmdLine(),
		Prober:  kmodule.Probe,
		ExcludedMods: map[string]bool{
			"idpf":     true,
			"idpf_imc": true,
		},
	}
}

// InstallAllModules installs kernel modules form the following locations in order:
// - .ko files from /lib/modules
// - modules found in .conf files from /lib/modules-load.d/
// - modules found in the cmdline argument modules_load= separated by ,
// Useful for modules that need to be loaded for boot (ie a network
// driver needed for netboot). It skips over blacklisted modules in
// excludedMods.
func InstallAllModules() error {
	loader := NewInitModuleLoader()
	modulePattern := "/lib/modules/*.ko"
	if err := InstallModulesFromDir(modulePattern, loader); err != nil {
		return err
	}
	var allModules []string
	moduleConfPattern := "/lib/modules-load.d/*.conf"
	modules, err := GetModulesFromConf(moduleConfPattern)
	if err != nil {
		return err
	}
	allModules = append(allModules, modules...)
	modules, err = GetModulesFromCmdline(loader)
	if err != nil {
		return err
	}
	allModules = append(allModules, modules...)
	InstallModules(loader, allModules)
	return nil
}

// InstallModules installs the passed modules using the InitModuleLoader
func InstallModules(m *InitModuleLoader, modules []string) {
	for _, moduleName := range modules {
		if m.IsExcluded(moduleName) {
			log.Printf("Skipping module %q", moduleName)
			continue
		}
		if err := m.LoadModule(moduleName); err != nil {
			log.Printf("InstallModulesFromModulesLoad: can't install %q: %v", moduleName, err)
		}
	}
}

// InstallModulesFromDir installs kernel modules (.ko files) from /lib/modules that
// match the given pattern, skipping those in the exclude list.
func InstallModulesFromDir(pattern string, loader *InitModuleLoader) error {
	files, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no modules found matching '%s'", pattern)
	}

	for _, filename := range files {
		f, err := os.Open(filename)
		if err != nil {
			log.Printf("InstallModules: can't open %q: %v", filename, err)
			continue
		}
		defer f.Close()
		// Module flags are passed to the command line in the from modulename.flag=val
		// And must be passed to FileInit as flag=val to be installed properly
		moduleName := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
		if loader.IsExcluded(moduleName) {
			log.Printf("Skipping module %q", moduleName)
			continue
		}

		flags := cmdline.FlagsForModule(moduleName)
		if err = kmodule.FileInit(f, flags, 0); err != nil {
			log.Printf("InstallModules: can't install %q: %v", filename, err)
		}
	}

	return nil
}

func readModules(f *os.File) []string {
	scanner := bufio.NewScanner(f)
	modules := []string{}
	for scanner.Scan() {
		i := scanner.Text()
		i = strings.TrimSpace(i)
		if i == "" || strings.HasPrefix(i, "#") {
			continue
		}
		modules = append(modules, i)
	}
	if err := scanner.Err(); err != nil {
		log.Println("error on reading:", err)
	}
	return modules
}

// GetModulesFromConf finds kernel modules from .conf files in /lib/modules-load.d/
func GetModulesFromConf(pattern string) ([]string, error) {
	var ret []string
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	for _, filename := range files {
		f, err := os.Open(filename)
		if err != nil {
			log.Printf("InstallModulesFromModulesLoad: can't open %q: %v", filename, err)
			continue
		}
		defer f.Close()
		modules := readModules(f)
		ret = append(ret, modules...)
	}
	return ret, nil
}

// GetModulesFromCmdline finds kernel modules from the modules_load kernel parameter
func GetModulesFromCmdline(m *InitModuleLoader) ([]string, error) {
	var ret []string
	modules, present := m.Cmdline.Flag("modules_load")
	if !present {
		return nil, nil
	}

	for _, moduleName := range strings.Split(modules, ",") {
		moduleName = strings.TrimSpace(moduleName)
		ret = append(ret, moduleName)
	}
	return ret, nil
}
