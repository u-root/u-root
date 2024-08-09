// Copyright 2009-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build freebsd || linux || darwin

// ldd returns library dependencies of an executable.
//
// The way this is done on GNU-based systems is interesting. For each ELF, one
// finds the .interp section. If there is no interpreter there's not much to
// do.
//
// If there is an interpreter, we run it with the --list option and the file as
// an argument. We need to parse the output. For all lines with =>  as the 2nd
// field, we take the 3rd field as a dependency.
//
// On many Unix kernels, the kernel ABI is stable. On OSX, the stability
// is held in the library interface; the kernel ABI is explicitly not
// stable. The ldd package on OSX will only return the files passed to it.
package ldd
