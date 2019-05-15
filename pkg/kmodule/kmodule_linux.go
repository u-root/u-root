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

	"github.com/ulikunitz/xz"
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
	DryRunCB       func(string)
	RootDir        string
	KVer           string
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
	if err := LoadModule(modPath, modParams, 0, opts); err != nil {
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
	if err := LoadModule(path, "", 0, opts); err != nil {
		return err
	}
	m[path].state = loaded

	return nil
}

// loadModule loads a file as a module. Because of all the various hacks -- er,
// improvements, made to module loading over 20 years, this gets messy.
// o finit_module is optional and only takes an fd
// o you can only use flags with finit_module
// o some distros have compressed modules
//   - can you uncompress those to a file?
//     what if /tmp is not there for some reason? What if the module you are
//     loading is tmpfs and you can't make temp files until you modprobe tmpfs?
//     [this can happen]
//   - Can you write the uncompressed data to a pipe and pass that fd?
//     You may not be able to write all the data to a pipe. 4k is a common size.
//     What if the module loader stats the pipe? It won't get the right answer.
//     What if the writer gets behind and kernel reads 0 bytes for some reason and assumes EOF?
//     What if a writer messes up due to bad compression and you block the kernel?
// It's a mess.
//
// With luck, this function covers all the cases.
// I put it back in one function as I found it a bit more readable straight line --
// it's less than a page.
// Note that, at present, we don't even set flags. Arguably, using FinitModule
// is a waste of time.
func LoadModule(path, modParams string, flags int, opts ProbeOpts) error {
	if opts.DryRunCB != nil {
		opts.DryRunCB(path)
		return nil
	}

	isXZ := strings.HasSuffix(path, ".xz")
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if !isXZ {
		// Use e instead of err to ensure visitors from the future don't see err
		// and think the := is a typo.
		if e := unix.FinitModule(int(f.Fd()), modParams, flags); e != unix.ENOSYS {
			return e
		}
	}

	if flags != 0 {
		return fmt.Errorf("Can not have non-zero flags (%#x) with init_module", flags)
	}

	var r = io.Reader(f)
	if isXZ {
		if r, err = xz.NewReader(f); err != nil {
			return err
		}
	}
	img, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	return Init(img, modParams)
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
