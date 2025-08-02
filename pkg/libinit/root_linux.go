// Copyright 2014-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package libinit creates the environment and root file system for u-root.
package libinit

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/u-root/u-root/pkg/cmdline"
	"github.com/u-root/u-root/pkg/core/cp"
	"github.com/u-root/u-root/pkg/kmodule"
	"github.com/u-root/u-root/pkg/ulog"
	"golang.org/x/sys/unix"
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

type Dev struct {
	Name string
	Mode uint32
	Dev  int
}

func (d Dev) Create() error {
	os.Remove(d.Name)
	return unix.Mknod(d.Name, d.Mode, d.Dev)
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
	return unix.Mount(m.Source, m.Target, m.FSType, m.Flags, m.Opts)
}

func (m Mount) String() string {
	return fmt.Sprintf("mount -t %q -o %s %q %q flags %#x", m.FSType, m.Opts, m.Source, m.Target, m.Flags)
}

type CpDir struct {
	Source string
	Target string
}

func (c CpDir) Create() error {
	copier := cp.Options{
		NoFollowSymlinks: true,
	}
	return copier.CopyTree(c.Source, c.Target)
}

func (c CpDir) String() string {
	return fmt.Sprintf("cp -a %q %q", c.Source, c.Target)
}

var (
	// These have to be created / mounted first, so that the logging works correctly.
	PreNamespace = []Creator{
		Dir{Name: "/dev", Mode: 0o777},

		// Kernel must be compiled with CONFIG_DEVTMPFS.
		Mount{Source: "devtmpfs", Target: "/dev", FSType: "devtmpfs"},
	}
	Namespace = []Creator{
		Dir{Name: "/buildbin", Mode: 0o777},
		Dir{Name: "/ubin", Mode: 0o777},
		Dir{Name: "/tmp", Mode: 0o777},
		Dir{Name: "/env", Mode: 0o777},
		Dir{Name: "/tcz", Mode: 0o777},
		Dir{Name: "/lib", Mode: 0o777},
		Dir{Name: "/usr/lib", Mode: 0o777},
		Dir{Name: "/var/log", Mode: 0o777},
		Dir{Name: "/go/pkg/linux_amd64", Mode: 0o777},

		Dir{Name: "/etc", Mode: 0o777},

		Dir{Name: "/proc", Mode: 0o555},
		Mount{Source: "proc", Target: "/proc", FSType: "proc"},
		Mount{Source: "tmpfs", Target: "/tmp", FSType: "tmpfs"},

		Dev{Name: "/dev/tty", Mode: unix.S_IFCHR | 0o666, Dev: 0x0500},
		Dev{Name: "/dev/urandom", Mode: unix.S_IFCHR | 0o444, Dev: 0x0109},
		Dev{Name: "/dev/port", Mode: unix.S_IFCHR | 0o640, Dev: 0x0104},
		Dev{Name: "/dev/ttyhvc0", Mode: unix.S_IFCHR | 0o666, Dev: 0xe500},

		Dir{Name: "/dev/pts", Mode: 0o777},
		Mount{Source: "devpts", Target: "/dev/pts", FSType: "devpts", Opts: "newinstance,ptmxmode=666,gid=5,mode=620"},
		// Note: if we mount /dev/pts with "newinstance", we *must* make "/dev/ptmx" a symlink to "/dev/pts/ptmx"
		Symlink{NewPath: "/dev/ptmx", Target: "/dev/pts/ptmx"},
		// Note: shm is required at least for Chrome. If you don't mount
		// it chrome throws a bogus "out of memory" error, not the more
		// useful "I can't open /dev/shm/whatever". SAD!
		Dir{Name: "/dev/shm", Mode: 0o777},
		Mount{Source: "tmpfs", Target: "/dev/shm", FSType: "tmpfs"},

		Dir{Name: "/sys", Mode: 0o555},
		Mount{Source: "sysfs", Target: "/sys", FSType: "sysfs"},
		Mount{Source: "securityfs", Target: "/sys/kernel/security", FSType: "securityfs"},
		Mount{Source: "efivarfs", Target: "/sys/firmware/efi/efivars", FSType: "efivarfs"},
		Mount{Source: "debugfs", Target: "/sys/kernel/debug", FSType: "debugfs"},

		CpDir{Source: "/etc", Target: "/tmp/etc"},
		Mount{Source: "/tmp/etc", Target: "/etc", FSType: "tmpfs", Flags: unix.MS_BIND},
	}

	// cgroups are optional for most u-root users, especially
	// LinuxBoot/NERF. Some users use u-root for container stuff.
	CgroupsNamespace = []Creator{
		Mount{Source: "cgroup", Target: "/sys/fs/cgroup", FSType: "tmpfs"},
		Dir{Name: "/sys/fs/cgroup/memory", Mode: 0o555},
		Dir{Name: "/sys/fs/cgroup/freezer", Mode: 0o555},
		Dir{Name: "/sys/fs/cgroup/devices", Mode: 0o555},
		Dir{Name: "/sys/fs/cgroup/cpu,cpuacct", Mode: 0o555},
		Dir{Name: "/sys/fs/cgroup/blkio", Mode: 0o555},
		Dir{Name: "/sys/fs/cgroup/cpuset", Mode: 0o555},
		Dir{Name: "/sys/fs/cgroup/pids", Mode: 0o555},
		Dir{Name: "/sys/fs/cgroup/net_cls,net_prio", Mode: 0o555},
		Dir{Name: "/sys/fs/cgroup/hugetlb", Mode: 0o555},
		Dir{Name: "/sys/fs/cgroup/perf_event", Mode: 0o555},
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
)

func goBin() string {
	return fmt.Sprintf("/go/bin/%s_%s:/go/bin:/go/pkg/tool/%s_%s", runtime.GOOS, runtime.GOARCH, runtime.GOOS, runtime.GOARCH)
}

func Create(namespace []Creator, optional bool) {
	// Clear umask bits so that we get stuff like ptmx right.
	m := unix.Umask(0)
	defer unix.Umask(m)
	for _, c := range namespace {
		if err := c.Create(); err != nil {
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
	Create(PreNamespace, false)
	ulog.KernelLog.Reinit()

	Create(Namespace, false)

	// systemd gets upset when it discovers something has already setup cgroups
	// We have to do this after the base namespace is created, so we have /proc
	initFlags := cmdline.GetInitFlagMap()
	systemd, present := initFlags["systemd"]
	systemdEnabled, boolErr := strconv.ParseBool(systemd)
	if !present || boolErr != nil || !systemdEnabled {
		Create(CgroupsNamespace, true)
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
		return fmt.Errorf("failed to load module: %w", err)
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
	if err := InstallModulesFromDir(modulePattern, loader); !errors.Is(err, ErrNoModulesFound) {
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

// ErrNoModulesFound is the error returned when InstallModulesFromDir does not
// find any valid modules in the path.
var ErrNoModulesFound = fmt.Errorf("no modules found")

// InstallModulesFromDir installs kernel modules (.ko files) from /lib/modules that
// match the given pattern, skipping those in the exclude list.
func InstallModulesFromDir(pattern string, loader *InitModuleLoader) error {
	files, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return ErrNoModulesFound
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

// OpenTTYDevices opens the TTY devices with the given names in /dev.
// It uses a best-effort approach, returning the devices that could be opened.
// If no devices could be opened, it returns an error.
func OpenTTYDevices(names []string) ([]io.Writer, error) {
	return openTTYDevices("/dev", names)
}

func openTTYDevices(prefix string, names []string) ([]io.Writer, error) {
	var err error
	devs := make([]io.Writer, 0, len(names))

	if len(names) == 0 {
		return devs, nil
	}

	for _, name := range names {
		d, e := os.OpenFile(filepath.Join(prefix, name), os.O_RDWR, 0)
		if e != nil {
			ulog.KernelLog.Printf("open TTY: %v", e)
			err = errors.Join(err, e)
			continue
		}
		devs = append(devs, d)
	}

	if len(devs) == 0 {
		return devs, err
	}

	return devs, nil
}
