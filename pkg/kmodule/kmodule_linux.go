// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package kmodule interfaces with Linux kernel modules.
//
// kmodule allows loading and unloading kernel modules with dependencies, as
// well as locating them through probing.
package kmodule

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/klauspost/compress/zstd"
	"github.com/klauspost/pgzip"
	"github.com/ulikunitz/xz"
	"golang.org/x/sys/unix"
)

// Flags to finit_module(2) / FileInit.
const (
	// Ignore symbol version hashes.
	MODULE_INIT_IGNORE_MODVERSIONS = 0x1

	// Ignore kernel version magic.
	MODULE_INIT_IGNORE_VERMAGIC = 0x2

	modPath = "/proc/modules"
)

type LinuxLoader struct {
	modules string
	io.ReadCloser
}

// Init loads the kernel module given by image with the given options.
func (l *LinuxLoader) Init(image []byte, opts string) error {
	return unix.InitModule(image, opts)
}

// Read implements io.Reader
func (l *LinuxLoader) Read(b []byte) (int, error) {
	return l.ReadCloser.Read(b)
}

// FileInit loads the kernel module contained by `f` with the given opts and
// flags. Uncompresses modules with a .xz and .gz suffix before loading.
//
// FileInit falls back to init_module(2) via Init when the finit_module(2)
// syscall is not available and when loading compressed modules.
func (l *LinuxLoader) FileInit(f *os.File, opts string, flags uintptr) error {
	var r io.Reader
	var err error
	switch filepath.Ext(f.Name()) {
	case ".xz":
		if r, err = xz.NewReader(f); err != nil {
			return err
		}
	case ".gz":
		if r, err = pgzip.NewReader(f); err != nil {
			return err
		}
	case ".zst":
		if r, err = zstd.NewReader(f); err != nil {
			return err
		}
	}

	if r == nil {
		err := unix.FinitModule(int(f.Fd()), opts, int(flags))
		if err == unix.ENOSYS {
			if flags != 0 {
				return err
			}
			// Fall back to init_module(2).
			r = f
		} else {
			return err
		}
	}

	img, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	return l.Init(img, opts)
}

// FileInit loads the kernel module contained by `f` with the given opts and
// flags. Uncompresses modules with a .xz and .gz suffix before loading.
//
// FileInit falls back to init_module(2) via Init when the finit_module(2)
// syscall is not available and when loading compressed modules.
func FileInit(f *os.File, opts string, flags uintptr) error {
	l, err := New()
	if err != nil {
		return err
	}
	return l.FileInit(f, opts, flags)
}

// Delete removes a kernel module.
func (l *LinuxLoader) Delete(name string, flags uintptr) error {
	return unix.DeleteModule(name, int(flags))
}

type modState uint8

const (
	unloaded modState = iota
	loading
	loaded
	builtin
)

type dependency struct {
	state modState
	deps  []string
}

type depMap map[string]*dependency

// ProbeOpts contains optional parameters to Probe.
//
// An empty ProbeOpts{} should lead to the default behavior.
type ProbeOpts struct {
	DryRunCB       func(string)
	RootDir        string
	KVer           string
	IgnoreProcMods bool
}

// NewPath creates a *LinuxLoader using a given path.
func NewPath(n string) (*LinuxLoader, error) {
	f, err := os.Open(n)
	if err != nil {
		return nil, err
	}
	return &LinuxLoader{modules: n, ReadCloser: f}, nil
}

// New creates a *LinuxLoader from /proc/modules
func New() (*LinuxLoader, error) {
	return NewPath(modPath)
}

var _ Loader = &LinuxLoader{}

// ProbePath loads the given kernel module and its dependencies,
// using a provide path to a /proc/modules-like file.
func ProbePath(modules, name, modParams string) error {
	l, err := NewPath(modules)
	if err != nil {
		return err
	}
	return l.Probe(name, modParams)
}

// Probe loads the given kernel module and its dependencies.
// It is calls ProbeOptions with the default ProbeOpts.
func Probe(name, modParams string) error {
	return ProbePath(modPath, name, modParams)
}

// Probe loads the given kernel module and its dependencies.
// It is calls ProbeOptions with the default ProbeOpts.
func (l *LinuxLoader) Probe(name string, modParams string) error {
	return l.ProbeOptions(name, modParams, ProbeOpts{})
}

// ProbeOptions loads the given kernel module and its dependencies.
// This functions takes ProbeOpts.
func (l *LinuxLoader) ProbeOptions(name, modParams string, opts ProbeOpts) error {
	deps, err := l.genDeps(opts)
	if err != nil {
		return fmt.Errorf("could not generate dependency map %v", err)
	}

	modPath, err := l.findModPath(name, deps)
	if err != nil {
		return fmt.Errorf("could not find module path %q: %v", name, err)
	}

	dep := deps[modPath]

	if dep.state == builtin || dep.state == loaded {
		return nil
	}

	dep.state = loading
	for _, d := range dep.deps {
		if err := l.loadDeps(d, deps, opts); err != nil {
			return err
		}
	}
	return l.loadModule(modPath, modParams, opts)
}

// ProbeOptions loads the given kernel module and its dependencies.
// This functions takes ProbeOpts.
func ProbeOptions(name, modParams string, opts ProbeOpts) error {
	l, err := New()
	if err != nil {
		return err
	}
	return l.ProbeOptions(name, modParams, opts)
}

func (l *LinuxLoader) checkBuiltin(moduleDir string, deps depMap) error {
	f, err := os.Open(filepath.Join(moduleDir, "modules.builtin"))
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return fmt.Errorf("could not open builtin file: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		modPath := filepath.Join(moduleDir, strings.TrimSpace(txt))
		if deps[modPath] == nil {
			deps[modPath] = new(dependency)
		}
		deps[modPath].state = builtin
	}

	return scanner.Err()
}

func (l *LinuxLoader) genDeps(opts ProbeOpts) (depMap, error) {
	deps := make(depMap)
	rel := opts.KVer

	if rel == "" {
		var u unix.Utsname
		if err := unix.Uname(&u); err != nil {
			return nil, fmt.Errorf("could not get release (uname -r): %v", err)
		}
		rel = string(u.Release[:bytes.IndexByte(u.Release[:], 0)])
	}

	var moduleDir string
	for _, n := range []string{"/lib/modules", "/usr/lib/modules"} {
		moduleDir = filepath.Join(opts.RootDir, n, strings.TrimSpace(rel))
		if _, err := os.Stat(moduleDir); err == nil {
			break
		}
	}

	f, err := os.Open(filepath.Join(moduleDir, "modules.dep"))
	if err != nil {
		return nil, fmt.Errorf("could not open dependency file: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		nameDeps := strings.Split(txt, ":")
		modPath, modDeps := nameDeps[0], nameDeps[1]
		modPath = filepath.Join(moduleDir, strings.TrimSpace(modPath))

		var dependency dependency
		if len(modDeps) > 0 {
			for _, dep := range strings.Split(strings.TrimSpace(modDeps), " ") {
				dependency.deps = append(dependency.deps, filepath.Join(moduleDir, dep))
			}
		}
		deps[modPath] = &dependency
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if err = l.checkBuiltin(moduleDir, deps); err != nil {
		return nil, err
	}

	if !opts.IgnoreProcMods {
		fm, err := os.Open(modPath)
		if err == nil {
			defer fm.Close()
			if err := l.genLoadedMods(fm, deps); err != nil {
				return nil, err
			}
		}
	}

	return deps, nil
}

func (l *LinuxLoader) findModPath(name string, m depMap) (string, error) {
	// Kernel modules do not have any consistency with use of hyphens and underscores
	// matching from the module's name to the module's file path. Thus try matching
	// the provided name using either.
	nameH := strings.Replace(name, "_", "-", -1)
	nameU := strings.Replace(name, "-", "_", -1)

	for mp := range m {
		switch path.Base(mp) {
		case nameH + ".ko", nameH + ".ko.gz", nameH + ".ko.xz", nameH + ".ko.zst":
			return mp, nil
		case nameU + ".ko", nameU + ".ko.gz", nameU + ".ko.xz", nameU + ".ko.zst":
			return mp, nil
		}
	}

	return "", fmt.Errorf("could not find path for module %q", name)
}

func (l *LinuxLoader) loadDeps(path string, m depMap, opts ProbeOpts) error {
	dependency, ok := m[path]
	if !ok {
		return fmt.Errorf("could not find dependency %q", path)
	}

	if dependency.state == loading {
		return fmt.Errorf("circular dependency! %q already LOADING", path)
	} else if (dependency.state == loaded) || (dependency.state == builtin) {
		return nil
	}

	m[path].state = loading

	for _, dep := range dependency.deps {
		if err := l.loadDeps(dep, m, opts); err != nil {
			return err
		}
	}

	// done with dependencies, load module
	if err := l.loadModule(path, "", opts); err != nil {
		return err
	}
	m[path].state = loaded

	return nil
}

func (l *LinuxLoader) loadModule(path, modParams string, opts ProbeOpts) error {
	if opts.DryRunCB != nil {
		opts.DryRunCB(path)
		return nil
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := l.FileInit(f, modParams, 0); err != nil && err != unix.EEXIST {
		return err
	}

	return nil
}

func (l *LinuxLoader) genLoadedMods(r io.Reader, deps depMap) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		arr := strings.Split(scanner.Text(), " ")
		name := arr[0]
		modPath, err := l.findModPath(name, deps)
		if err != nil {
			return fmt.Errorf("could not find module path %q: %v", name, err)
		}
		if deps[modPath] == nil {
			deps[modPath] = new(dependency)
		}
		deps[modPath].state = loaded
	}
	return scanner.Err()
}

// Pretty prints the string from /proc/modules in a pretty format.
func Pretty(w io.Writer, s string) error {
	in := "Module\tSize\tUsed by\n"
	for i, line := range strings.Split(s, "\n") {
		if len(line) == 0 {
			continue
		}
		f := strings.Fields(line)
		if len(f) < 4 {
			in += fmt.Sprintf("Line %d is malformed: %q\n", i, line)
			continue
		}
		// Don't use fmt.Sprintf here; it turns tabs to spaces.
		in += f[0] + "\t" + f[1] + "\t" + f[2] + "\t" + f[3] + "\n"
	}
	tw := tabwriter.NewWriter(w, 0, 8, 0, '\t', 0)
	if _, err := tw.Write([]byte(in)); err != nil {
		return err
	}

	return tw.Flush()
}
