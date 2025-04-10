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
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/sys/unix"
)

// Flags to finit_module(2) / FileInit.
const (
	// Ignore symbol version hashes.
	MODULE_INIT_IGNORE_MODVERSIONS = 0x1

	// Ignore kernel version magic.
	MODULE_INIT_IGNORE_VERMAGIC = 0x2
)

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

// Init loads the kernel module given by image with the given options.
func Init(image []byte, opts string) error {
	return unix.InitModule(image, opts)
}

// Wrapper for the compression readers
func CompressionReader(file *os.File) (reader io.Reader, err error) {
	return compressionReader(file)
}

// FileInit loads the kernel module contained by `f` with the given opts and
// flags. Uncompresses modules with a .xz and .gz suffix before loading.
//
// FileInit falls back to init_module(2) via Init when the finit_module(2)
// syscall is not available and when loading compressed modules.
func FileInit(f *os.File, opts string, flags uintptr) error {
	var r io.Reader
	var err error

	if r, err = CompressionReader(f); err != nil {
		return err
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
	return Init(img, opts)
}

// Delete removes a kernel module.
func Delete(name string, flags uintptr) error {
	return unix.DeleteModule(name, int(flags))
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
		return fmt.Errorf("could not generate dependency map %w", err)
	}

	modPath, err := findModPath(name, deps)
	if err != nil {
		return fmt.Errorf("could not find module path %q: %w", name, err)
	}

	dep := deps[modPath]

	if dep.state == builtin || dep.state == loaded {
		return nil
	}

	dep.state = loading
	for _, d := range dep.deps {
		if err := loadDeps(d, deps, opts); err != nil {
			return err
		}
	}
	return loadModule(modPath, modParams, opts)
}

func checkBuiltin(moduleDir string, deps depMap) error {
	f, err := os.Open(filepath.Join(moduleDir, "modules.builtin"))
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return fmt.Errorf("could not open builtin file: %w", err)
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

func genDeps(opts ProbeOpts) (depMap, error) {
	deps := make(depMap)
	rel := opts.KVer

	if rel == "" {
		var u unix.Utsname
		if err := unix.Uname(&u); err != nil {
			return nil, fmt.Errorf("could not get release (uname -r): %w", err)
		}
		rel = unix.ByteSliceToString(u.Release[:])
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
		return nil, fmt.Errorf("could not open dependency file: %w", err)
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

	if err = checkBuiltin(moduleDir, deps); err != nil {
		return nil, err
	}

	if !opts.IgnoreProcMods {
		fm, err := os.Open("/proc/modules")
		if err == nil {
			defer fm.Close()
			genLoadedMods(fm, deps)
		}
	}

	return deps, nil
}

func findModPath(name string, m depMap) (string, error) {
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

func loadDeps(path string, m depMap, opts ProbeOpts) error {
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
		return err
	}
	defer f.Close()

	if err := FileInit(f, modParams, 0); err != nil && err != unix.EEXIST {
		return err
	}

	return nil
}

func genLoadedMods(r io.Reader, deps depMap) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		arr := strings.Split(scanner.Text(), " ")
		name := arr[0]
		modPath, err := findModPath(name, deps)
		if err != nil {
			return fmt.Errorf("could not find module path %q: %w", name, err)
		}
		if deps[modPath] == nil {
			deps[modPath] = new(dependency)
		}
		deps[modPath].state = loaded
	}
	return scanner.Err()
}
