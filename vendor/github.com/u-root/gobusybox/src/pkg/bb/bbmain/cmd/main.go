// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package main is the busybox main.go template.
package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/u-root/gobusybox/src/pkg/bb/bbmain"
	// There MUST NOT be any other dependencies here.
	//
	// It is preferred to copy minimal code necessary into this file, as
	// dependency management for this main file is... hard.
)

// AbsSymlink returns an absolute path for the link from a file to a target.
func AbsSymlink(originalFile, target string) string {
	if !filepath.IsAbs(originalFile) {
		var err error
		originalFile, err = filepath.Abs(originalFile)
		if err != nil {
			// This should not happen on Unix systems, or you're
			// already royally screwed.
			log.Fatalf("could not determine absolute path for %v: %v", originalFile, err)
		}
	}
	// Relative symlinks are resolved relative to the original file's
	// parent directory.
	//
	// E.g. /bin/defaultsh -> ../bbin/elvish
	if !filepath.IsAbs(target) {
		return filepath.Join(filepath.Dir(originalFile), target)
	}
	return target
}

// IsTargetSymlink returns true if a target of a symlink is also a symlink.
func IsTargetSymlink(originalFile, target string) bool {
	s, err := os.Lstat(AbsSymlink(originalFile, target))
	if err != nil {
		return false
	}
	return (s.Mode() & os.ModeSymlink) == os.ModeSymlink
}

// ResolveUntilLastSymlink resolves until the last symlink.
//
// This is needed when we have a chain of symlinks and want the last
// symlink, not the file pointed to (which is why we don't use
// filepath.EvalSymlinks)
//
// I.e.
//
// /foo/bar -> ../baz/foo
// /baz/foo -> bla
//
// ResolveUntilLastSymlink(/foo/bar) returns /baz/foo.
func ResolveUntilLastSymlink(p string) string {
	for target, err := os.Readlink(p); err == nil && IsTargetSymlink(p, target); target, err = os.Readlink(p) {
		p = AbsSymlink(p, target)
	}
	return p
}

func run() {
	// For shellbang files, the os.Arg[0] for different kernels is not always
	// consistent. In the case of most Unix, it is /bbin/bb; for Plan 9, it is the
	// base of the shellbang file path.
	// so, e.g, if /bbin/date is this:
	// #!/bbin/bb #!/bbin/date
	// on Plan 9, argv is
	// date #!/bbin/date /bbin/date
	// On most Unix, argv is
	// /bbin/bb #!/bbin/date /bbin/date
	// You can not just use os.Args[0] on Plan 9, because the first bbmain.Run below
	// will succeed, with two additional arguments, and the other two arguments will confuse it:
	// term% /bbin/date
	// Usage of date [-u] [-d FILE] [format] ...
	//
	// However:
	// u-root shellbang files have one argument for the interpreter: a string beginning
	// with #!, as above:
	// #!/bbin/bb #!/bbin/date
	// os.Args[2], for all kernels, is the full path or base of the command.
	// I.e., if we are pretty sure it is a shellbang, we can use os.Args[2].
	// So, if os.Args[1] is a shellbang, discard the first two arguments.
	// This was all explained several years ago in the first shellbang commit,
	// but that was discarded at some point. Please don't discard that again.
	//
	// Also, N.B.: the runtime.GOOS == "plan9" test will cause this
	// code to be elided by the toolchain. So don't worry about
	// some comparison happening here for all non-Plan 9 kernels.
	// This code won't make it into the final binary if it is not built for Plan 9.
	if runtime.GOOS == "plan9" && len(os.Args) > 2 && strings.HasPrefix(os.Args[1], "#!") {
		os.Args = os.Args[2:]
	}

	name := filepath.Base(os.Args[0])
	err := bbmain.Run(name)

	// This test should not run on Plan 9.
	if runtime.GOOS != "plan9" && errors.Is(err, bbmain.ErrNotRegistered) {
		if len(os.Args) > 1 {
			os.Args = os.Args[1:]
			err = bbmain.Run(filepath.Base(os.Args[0]))
		}
	}
	if errors.Is(err, bbmain.ErrNotRegistered) {
		log.SetFlags(0)
		log.Printf("Failed to run command: %v", err)

		log.Printf("Supported commands are:")
		for _, cmd := range bbmain.ListCmds() {
			log.Printf(" - %s", cmd)
		}
		os.Exit(1)
	} else if err != nil {
		log.SetFlags(0)
		log.Fatalf("Failed to run command: %v", err)
	}
}

func main() {
	os.Args[0] = ResolveUntilLastSymlink(os.Args[0])

	run()
}
