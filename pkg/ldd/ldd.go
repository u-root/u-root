// Copyright 2009-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build freebsd || linux || darwin
// +build freebsd linux darwin

// ldd returns all the library dependencies of an executable.
//
// The way this is done on GNU-based systems is interesting. For each ELF, one
// finds the .interp section. If there is no interpreter
// there's not much to do.
//
// If there is an interpreter, we run it with the --list option and the file as
// an argument. We need to parse the output.
//
// For all lines with =>  as the 2nd field, we take the 3rd field as a
// dependency. The field may be a symlink.  Rather than stat the link and do
// other such fooling around, we can do a readlink on it; if it fails, we just
// need to add that file name; if it succeeds, we need to add that file name
// and repeat with the next link in the chain. We can let the kernel do the
// work of figuring what to do if and when we hit EMLINK.
//
// On many Unix kernels, the kernel ABI is stable. On OSX, the stability
// is held in the library interface; the kernel ABI is explicitly not
// stable. The ldd package on OSX will only return the files passed to it.
// It will continue to resolve symbolic links.
package ldd

type FileInfo struct {
	FullName string
	os.FileInfo
}

// Paths returns the dependency file paths of files in names.
func Paths(names ...string) ([]string, error) {
	l, err := Ldd(names...)
	if err != nil {
		return nil, err
	}

	var list []string
	for i := range l {
		list = append(list, l[i].FullName)
	}
	return list, nil
}
