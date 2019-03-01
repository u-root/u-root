// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kmodule

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
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

// Init loads the kernel module given by image with the given options.
func Init(image []byte, opts string) error {
	return unix.InitModule(image, opts)
}

// FileInit loads the kernel module contained by `f` with the given opts and
// flags.
//
// FileInit falls back to Init when the finit_module(2) syscall is not available.
func FileInit(f *os.File, opts string, flags uintptr) error {
	err := unix.FinitModule(int(f.Fd()), opts, int(flags))
	if err == unix.ENOSYS {
		if flags != 0 {
			return err
		}

		// Fall back to regular init_module(2).
		img, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}
		return Init(img, opts)
	}

	return err
}

// Delete removes a kernel module.
func Delete(name string, flags uintptr) error {
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
	DryRunCB func(string)
	RootDir  string
	KVer     string
	IgnoreProcMods bool
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
		return fmt.Errorf("could not generate dependency map %v", err)
	}

	modPath, err := findModPath(name, deps)
	if err != nil {
		return fmt.Errorf("could not find module path %q: %v", name, err)
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
	if err := loadModule(modPath, modParams, opts); err != nil {
		return err
	}

	return nil
}

func checkBuiltin(moduleDir string, deps depMap) error {
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
		name := strings.Replace(arr[0], "_", "-", -1)
		modPath, err := findModPath(name, deps)
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
