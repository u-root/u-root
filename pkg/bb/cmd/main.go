// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/u-root/u-root/pkg/bb"
)

func run() {
	name := filepath.Base(os.Args[0])
	if err := bb.Run(name); err != nil {
		log.Fatalf("%s: %v", name, err)
	}
}

func isTargetSymlink(originalFile, target string) bool {
	s, err := os.Lstat(absSymlink(originalFile, target))
	if err != nil {
		return false
	}
	return (s.Mode() & os.ModeSymlink) == os.ModeSymlink
}

func absSymlink(originalFile, target string) string {
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

// resolveUntilLastSymlink resolves until the last symlink, e.g.
func resolveUntilLastSymlink(p string) string {
	for target, err := os.Readlink(p); err == nil && isTargetSymlink(p, target); target, err = os.Readlink(p) {
		p = absSymlink(p, target)
	}
	return p
}

func main() {
	os.Args[0] = resolveUntilLastSymlink(os.Args[0])

	run()
}

func init() {
	m := func() {
		// Use argv[1] as the name.
		os.Args = os.Args[1:]
		run()
	}
	bb.Register("bb", bb.Noop, m)
	bb.RegisterDefault(bb.Noop, m)
}
