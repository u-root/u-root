// Copyright 2009-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build freebsd || linux

package ldd

import (
	"debug/elf"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func parseinterp(input string) ([]string, error) {
	var names []string
	for _, p := range strings.Split(input, "\n") {
		f := strings.Fields(p)
		if len(f) < 3 {
			continue
		}
		if f[1] != "=>" || len(f[2]) == 0 {
			continue
		}
		if f[0] == f[2] {
			continue
		}
		// If the third part is a memory address instead
		// of a file system path, the entry should be skipped.
		// For example: linux-vdso.so.1 => (0x00007ffe4972d000)
		if f[1] == "=>" && string(f[2][0]) == "(" {
			continue
		}
		names = append(names, f[2])
	}
	return names, nil
}

// runinterp runs the interpreter with the --list switch
// and the file as an argument. For each returned line
// it looks for => as the second field, indicating a
// real .so (as opposed to the .vdso or a string like
// 'not a dynamic executable'.
func runinterp(interp, file string) ([]string, error) {
	o, err := exec.Command(interp, "--list", file).Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("%w: %s", err, ee.Stderr)
		}
		return nil, err
	}
	return parseinterp(string(o))
}

func GetInterp(file string) (string, error) {
	r, err := os.Open(file)
	if err != nil {
		return "fail", err
	}
	defer r.Close()
	f, err := elf.NewFile(r)
	if err != nil {
		return "", nil
	}

	s := f.Section(".interp")
	var interp string
	if s != nil {
		// If there is an interpreter section, it should be
		// an error if we can't read it.
		i, err := s.Data()
		if err != nil {
			return "fail", err
		}

		// .interp section is file name + \0 character.
		interp := strings.TrimRight(string(i), "\000")

		// Ignore #! interpreters
		if strings.HasPrefix(interp, "#!") {
			return "", nil
		}
		return interp, nil
	}

	if interp == "" {
		if f.Type != elf.ET_DYN || f.Class == elf.ELFCLASSNONE {
			return "", nil
		}
		bit64 := true
		if f.Class != elf.ELFCLASS64 {
			bit64 = false
		}

		// This is a shared library. Turns out you can run an
		// interpreter with --list and this shared library as an
		// argument. What interpreter do we use? Well, there's no way to
		// know. You have to guess.  I'm not sure why they could not
		// just put an interp section in .so's but maybe that would
		// cause trouble somewhere else.
		interp, err = LdSo(bit64)
		if err != nil {
			return "fail", err
		}
	}
	return interp, nil
}

// follow returns all paths and any files they recursively point to through
// symlinks.
func follow(paths ...string) ([]string, error) {
	seen := make(map[string]struct{})

	for _, path := range paths {
		if err := followInternal(path, seen); err != nil {
			return nil, err
		}
	}

	deps := make([]string, 0, len(seen))
	for s := range seen {
		deps = append(deps, s)
	}
	return deps, nil
}

func followInternal(path string, seen map[string]struct{}) error {
	for {
		if _, ok := seen[path]; ok {
			return nil
		}
		i, err := os.Lstat(path)
		if err != nil {
			return err
		}

		seen[path] = struct{}{}
		if i.Mode().IsRegular() {
			return nil
		}

		// If it's a symlink, read works; if not, it fails.
		// We can skip testing the type, since we still have to
		// handle any error if it's a link.
		next, err := os.Readlink(path)
		if err != nil {
			return err
		}

		// A relative link has to be interpreted relative to the file's
		// parent's path.
		if !filepath.IsAbs(next) {
			next = filepath.Join(filepath.Dir(path), next)
		}
		path = next
	}
}

// List returns a list of all library dependencies for a set of files.
//
// If a file has no dependencies, that is not an error. The only possible error
// is if a file does not exist, or it says it has an interpreter but we can't
// read it, or we are not able to run its interpreter.
//
// It's not an error for a file to not be an ELF.
func List(names ...string) ([]string, error) {
	list := make(map[string]struct{})
	interps := make(map[string]struct{})
	for _, n := range names {
		interp, err := GetInterp(n)
		if err != nil {
			return nil, err
		}
		if interp == "" {
			continue
		}
		interps[interp] = struct{}{}

		// Run the interpreter to get dependencies.
		sonames, err := runinterp(interp, n)
		if err != nil {
			return nil, err
		}
		for _, name := range sonames {
			list[name] = struct{}{}
		}
	}

	libs := make([]string, 0, len(list)+len(interps))

	// People expect to see the interps first.
	for s := range interps {
		libs = append(libs, s)
	}
	for s := range list {
		libs = append(libs, s)
	}
	return libs, nil
}

// FList returns a list of all library dependencies for a set of files,
// including following symlinks.
//
// If a file has no dependencies, that is not an error. The only possible error
// is if a file does not exist, or it says it has an interpreter but we can't
// read it, or we are not able to run its interpreter.
//
// It's not an error for a file to not be an ELF.
func FList(names ...string) ([]string, error) {
	deps, err := List(names...)
	if err != nil {
		return nil, err
	}
	return follow(deps...)
}
