// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kmodule

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

// Flags to finit_module(2) / FileInit.
const (
	// Ignore symbol version hashes.
	MODULE_INIT_IGNORE_MODVERSIONS = 0x1

	// Ignore kernel version magic.
	MODULE_INIT_IGNORE_VERMAGIC = 0x2
)

// SyscallError contains an error message as well as the actual syscall Errno
type SyscallError struct {
	Msg   string
	Errno syscall.Errno
}

func (s *SyscallError) Error() string {
	if s.Errno != 0 {
		return fmt.Sprintf("%s: %v", s.Msg, s.Errno)
	}
	return s.Msg
}

// Init loads the kernel module given by image with the given options.
func Init(image []byte, opts string) error {
	optsNull, err := unix.BytePtrFromString(opts)
	if err != nil {
		return &SyscallError{Msg: fmt.Sprintf("kmodule.Init: could not convert %q to C string: %v", opts, err)}
	}

	if _, _, e := unix.Syscall(unix.SYS_INIT_MODULE, uintptr(unsafe.Pointer(&image[0])), uintptr(len(image)), uintptr(unsafe.Pointer(optsNull))); e != 0 {
		return &SyscallError{
			Msg:   fmt.Sprintf("init_module(%v, %q) failed", image, opts),
			Errno: e,
		}
	}

	return nil
}

// FileInit loads the kernel module contained by `f` with the given opts and
// flags.
//
// FileInit falls back to Init when the finit_module(2) syscall is not available.
func FileInit(f *os.File, opts string, flags uintptr) error {
	optsNull, err := unix.BytePtrFromString(opts)
	if err != nil {
		return &SyscallError{Msg: fmt.Sprintf("kmodule.Init: could not convert %q to C string: %v", opts, err)}
	}

	if _, _, e := unix.Syscall(unix.SYS_FINIT_MODULE, f.Fd(), uintptr(unsafe.Pointer(optsNull)), flags); e == unix.ENOSYS {
		if flags != 0 {
			return &SyscallError{Msg: fmt.Sprintf("finit_module unavailable"), Errno: e}
		}

		// Fall back to regular init_module(2).
		img, err := ioutil.ReadAll(f)
		if err != nil {
			return &SyscallError{Msg: fmt.Sprintf("kmodule.FileInit: %v", err)}
		}
		return Init(img, opts)
	} else if e != 0 {
		return &SyscallError{
			Msg:   fmt.Sprintf("finit_module(%v, %q, %#x) failed", f, opts, flags),
			Errno: e,
		}
	}

	return nil
}

// Delete removes a kernel module.
func Delete(name string, flags uintptr) error {
	modnameptr, err := unix.BytePtrFromString(name)
	if err != nil {
		return &SyscallError{Msg: fmt.Sprintf("could not delete module %q: %v", name, err)}
	}

	if _, _, e := unix.Syscall(unix.SYS_DELETE_MODULE, uintptr(unsafe.Pointer(modnameptr)), flags, 0); e != 0 {
		return &SyscallError{Msg: fmt.Sprintf("could not delete module %q", name), Errno: e}
	}

	return nil
}

type modState uint8

const (
	unloaded modState = iota
	loading
	loaded
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
	DryRunCB func(string)
	RootDir  string
	KVer     string
}

// Probe loads the given kernel module and its dependencies.
// It is calls ProbeOptions with the default ProbeOpts.
func Probe(name string, modParams string) error {
	return ProbeOptions(name, modParams, ProbeOpts{})
}

// ProbeOptions loads the given kernel module and its dependencies.
// This functions takes ProbeOpts.
func ProbeOptions(name, modParams string, opts ProbeOpts) error {
	deps, err := genDeps(opts)
	if err != nil {
		return &SyscallError{Msg: fmt.Sprintf("could not generate dependency map %v", err)}
	}

	modPath, err := findModPath(name, deps)
	if err != nil {
		return &SyscallError{Msg: fmt.Sprintf("could not find module path %q: %v", name, err)}
	}

	if opts.DryRunCB == nil {
		// if the module is already loaded or does not have deps, or all of them are loaded
		// then this succeeds and we are done
		if err := loadModule(modPath, modParams, opts); err == nil {
			return nil
		}
		// okay, we have to try the hard way and load dependencies first.
	}

	deps[modPath].state = loading
	for _, d := range deps[modPath].deps {
		if err := loadDeps(d, deps, opts); err != nil {
			return err
		}
	}
	if err := loadModule(modPath, modParams, opts); err != nil {
		return err
	}
	// we don't care to set the state to loaded
	// deps[modPath].state = loaded
	return nil
}

func genDeps(opts ProbeOpts) (depMap, error) {
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

	return deps, nil
}

func findModPath(name string, m depMap) (string, error) {
	for mp := range m {
		if path.Base(mp) == name+".ko" {
			return mp, nil
		}
	}

	return "", fmt.Errorf("Could not find path for module %q", name)
}

func loadDeps(path string, m depMap, opts ProbeOpts) error {
	dependency, ok := m[path]
	if !ok {
		return &SyscallError{Msg: fmt.Sprintf("could not find dependency %q", path)}
	}

	if dependency.state == loading {
		return &SyscallError{Msg: fmt.Sprintf("circular dependency! %q already LOADING", path)}
	} else if dependency.state == loaded {
		return nil
	}

	m[path].state = loading

	for _, dep := range dependency.deps {
		if err := loadDeps(dep, m, opts); err != nil {
			return err
		}
	}

	// done with dependencies, load module
	if err := loadModule(path, "", opts); err != nil {
		return err
	}
	m[path].state = loaded

	return nil
}

func loadModule(path, modParams string, opts ProbeOpts) error {
	if opts.DryRunCB != nil {
		opts.DryRunCB(path)
		return nil
	}

	f, err := os.Open(path)
	if err != nil {
		return &SyscallError{Msg: fmt.Sprintf("could not open %q: %v", path, err)}
	}
	defer f.Close()

	if err := FileInit(f, modParams, 0); err != nil {
		if serr, ok := err.(*SyscallError); !ok || (ok && serr.Errno != unix.EEXIST) {
			return err
		}
	}

	return nil
}
