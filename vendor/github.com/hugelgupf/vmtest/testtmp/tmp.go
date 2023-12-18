// Copyright 2023 the Go Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package testtmp provides a temporary directory for tests that is only
// removed if the test passes.
//
// The directories are also retained if --keep-temp-dir is passed to the test.
package testtmp

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"unicode"
	"unicode/utf8"
)

var (
	keepTempDir = flag.Bool("keep-temp-dir", false, "Keep temporary directory after test, even if test passed")
)

var (
	mu       sync.Mutex
	tempDirs = map[string]string{}
	tempIdx  = map[string]int{}
)

// TempDir creates a temporary directory that is only cleaned up if the test
// passes.
//
// Each call to TempDir creates a new directory.
//
// If the test fails or if --keep-temp-dir is set, it will not be removed.
func TempDir(t testing.TB) string {
	mu.Lock()
	rootDir, ok := tempDirs[t.Name()]
	var rootErr error
	if !ok {
		// Drop unusual characters (such as path separators or
		// characters interacting with globs) from the directory name to
		// avoid surprising os.MkdirTemp behavior.
		//
		// Stolen from Go testing.T.TempDir()
		mapper := func(r rune) rune {
			if r < utf8.RuneSelf {
				const allowed = "!#$%&()+,-.=@^_{}~ "
				if '0' <= r && r <= '9' ||
					'a' <= r && r <= 'z' ||
					'A' <= r && r <= 'Z' {
					return r
				}
				if strings.ContainsRune(allowed, r) {
					return r
				}
			} else if unicode.IsLetter(r) || unicode.IsNumber(r) {
				return r
			}
			return -1
		}
		pattern := strings.Map(mapper, t.Name())

		rootDir, rootErr = os.MkdirTemp("", pattern)
		if rootErr == nil {
			tempDirs[t.Name()] = rootDir
			t.Cleanup(func() {
				switch {
				case t.Failed():
					t.Logf("Keeping temp dir due to test failure: %s", rootDir)

				case *keepTempDir:
					t.Logf("Keeping temp dir as requested by --keep-temp-dir: %s", rootDir)

				default:
					if err := os.RemoveAll(rootDir); err != nil {
						t.Errorf("Failed to remove temporary directory %s: %v", rootDir, err)
					}
					// Delete map keys for repeated test cases.
					mu.Lock()
					delete(tempDirs, t.Name())
					delete(tempIdx, t.Name())
					mu.Unlock()
				}
			})
		}
	}

	idx := tempIdx[t.Name()]
	if rootErr == nil {
		tempIdx[t.Name()]++
	}
	// Unlock before calling any Fatal function.
	mu.Unlock()

	if rootErr != nil {
		t.Fatalf("Failed to create temp dir for %s: %v", t.Name(), rootErr)
	}

	dir := filepath.Join(rootDir, fmt.Sprintf("%03d", idx))
	if err := os.Mkdir(dir, 0777); err != nil {
		t.Fatalf("Mkdir: %v", err)
	}
	return dir
}
