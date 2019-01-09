// Copyright 2009-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// ldd returns all the library dependencies
// of a list of file names.
// The way this is done on GNU-based systems
// is interesting. For each ELF, one finds the
// .interp section. If there is no interpreter
// there's not much to do. If there is an interpreter,
// we run it with the --list option and the file as an argument.
// We need to parse the output.
// For all lines with =>  as the 2nd field, we take the
// 3rd field as a dependency. The field may be a symlink.
// Rather than stat the link and do other such fooling around,
// we can do a readlink on it; if it fails, we just need to add
// that file name; if it succeeds, we need to add that file name
// and repeat with the next link in the chain. We can let the
// kernel do the work of figuring what to do if and when we hit EMLINK.
package ldd

import (
	"debug/elf"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type FileInfo struct {
	FullName string
	os.FileInfo
}

// Follow starts at a pathname and adds it
// to a map if it is not there.
// If the pathname is a symlink, indicated by the Readlink
// succeeding, links repeats and continues
// for as long as the name is not found in the map.
func follow(l string, names map[string]*FileInfo) error {
	for {
		if names[l] != nil {
			return nil
		}
		i, err := os.Lstat(l)
		if err != nil {
			return fmt.Errorf("%v", err)
		}

		names[l] = &FileInfo{FullName: l, FileInfo: i}
		if i.Mode().IsRegular() {
			return nil
		}
		// If it's a symlink, the read works; if not, it fails.
		// we can skip testing the type, since we still have to
		// handle any error if it's a link.
		next, err := os.Readlink(l)
		if err != nil {
			return err
		}
		// It may be a relative link, so we need to
		// make it abs.
		if filepath.IsAbs(next) {
			l = next
			continue
		}
		l = filepath.Join(filepath.Dir(l), next)
	}
}

// runinterp runs the interpreter with the --list switch
// and the file as an argument. For each returned line
// it looks for => as the second field, indicating a
// real .so (as opposed to the .vdso or a string like
// 'not a dynamic executable'.
func runinterp(interp, file string) ([]string, error) {
	var names []string
	o, err := exec.Command(interp, "--list", file).Output()
	if err != nil {
		return nil, err
	}
	for _, p := range strings.Split(string(o), "\n") {
		f := strings.Split(p, " ")
		if len(f) < 3 {
			continue
		}
		if f[1] != "=>" || len(f[2]) == 0 {
			continue
		}
		names = append(names, f[2])
	}
	return names, nil
}

// Ldd returns a list of all library dependencies for a
// set of files, suitable for feeding into (e.g.) a cpio
// program. If a file has no dependencies, that is not an
// error. The only possible error is if a file does not
// exist, or it says it has an interpreter but we can't read
// it, or we are not able to run its interpreter.
// It's not an error for a file to not be an ELF, as
// this function should be convenient and the list might
// include non-ELF executables (a.out format, scripts)
func Ldd(names []string) ([]*FileInfo, error) {
	var (
		list    = make(map[string]*FileInfo)
		interps = make(map[string]*FileInfo)
		libs    []*FileInfo
	)
	for _, n := range names {
		if err := follow(n, list); err != nil {
			return nil, err
		}
	}
	for _, n := range names {
		r, err := os.Open(n)
		if err != nil {
			return nil, err
		}
		defer r.Close()
		f, err := elf.NewFile(r)
		if err != nil {
			continue
		}
		s := f.Section(".interp")
		var interp string
		if s != nil {
			// If there is an interpreter section, it should be
			// an error if we can't read it.
			i, err := s.Data()
			if err != nil {
				return nil, err
			}
			// Ignore #! interpreters
			if len(i) > 1 && i[0] == '#' && i[1] == '!' {
				continue
			}
			// annoyingly, s.Data() seems to return the null at the end and,
			// weirdly, that seems to confuse the kernel. Truncate it.
			interp = string(i[:len(i)-1])
		}
		if interp == "" {
			if f.Type != elf.ET_DYN || f.Class == elf.ELFCLASSNONE {
				continue
			}
			bit64 := true
			if f.Class != elf.ELFCLASS64 {
				bit64 = false
			}

			// This is a shared library. Turns out you can run an interpreter with
			// --list and this shared library as an argument. What interpreter
			// do we use? Well, there's no way to know. You have to guess.
			// I'm not sure why they could not just put an interp section in
			// .so's but maybe that would cause trouble somewhere else.
			interp, err = LdSo(bit64)
			if err != nil {
				return nil, err
			}
		}
		// We could just append the interp but people
		// expect to see that first.
		if interps[interp] == nil {
			err := follow(interp, interps)
			if err != nil {
				return nil, err
			}
		}
		// oh boy. Now to run the interp and get more names.
		n, err := runinterp(interp, n)
		if err != nil {
			return nil, err
		}
		for i := range n {
			if err := follow(n[i], list); err != nil {
				log.Fatalf("ldd: %v", err)
			}
		}
	}

	for i := range interps {
		libs = append(libs, interps[i])
	}

	for i := range list {
		libs = append(libs, list[i])
	}

	return libs, nil
}

func List(names []string) ([]string, error) {
	var list []string
	l, err := Ldd(names)
	if err != nil {
		return nil, err
	}
	for i := range l {
		list = append(list, l[i].FullName)
	}
	return list, nil
}
