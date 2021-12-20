// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// created by Manoel Vilela <manoel_vilela@engineer.com>

package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/cp"
	"github.com/u-root/u-root/pkg/uio"
	"golang.org/x/sys/unix"
)

const (
	maxSizeFile = 1000
	maxDirDepth = 5
	maxFiles    = 5
)

// resetFlags is used to reset the cp flags to default
func resetFlags() {
	flags.recursive = false
	flags.ask = false
	flags.force = false
	flags.verbose = false
	flags.noFollowSymlinks = false
}

// randomFile create a random file with random content
func randomFile(fpath, prefix string) (*os.File, error) {
	f, err := os.CreateTemp(fpath, prefix)
	if err != nil {
		return nil, err
	}
	// generate random content for files
	bytes := []byte{}
	for i := 0; i < rand.Intn(maxSizeFile); i++ {
		bytes = append(bytes, byte(i))
	}
	f.Write(bytes)

	return f, nil
}

// createFilesTree create a random files tree
func createFilesTree(root string, maxDepth, depth int) error {
	// create more one dir if don't achieve the maxDepth
	if depth < maxDepth {
		newDir, err := os.MkdirTemp(root, fmt.Sprintf("cpdir_%d_", depth))
		if err != nil {
			return err
		}

		if err = createFilesTree(newDir, maxDepth, depth+1); err != nil {
			return err
		}
	}
	// generate random files
	for i := 0; i < maxFiles; i++ {
		f, err := randomFile(root, fmt.Sprintf("cpfile_%d_", i))
		if err != nil {
			return err
		}
		f.Close()
	}

	return nil
}

func TestRunSimple(t *testing.T) {
	tmpDir := t.TempDir()
	file1, err := randomFile(tmpDir, "src-")
	if err != nil {
		t.Errorf("failed to create tmp dir: %q", err)
	}

	for _, tt := range []struct {
		name     string
		recursiv bool
		ask      bool
		force    bool
		verbose  bool
		noFollow bool
		args     []string
		wantErr  error
	}{
		{
			name: "NoFlags-Success-",
			args: []string{file1.Name(), filepath.Join(tmpDir, "destination")},
		},
		{
			name:    "Verbose",
			args:    []string{file1.Name(), filepath.Join(tmpDir, "destination")},
			verbose: true,
			wantErr: cp.ErrSkip,
		},
		{
			name:    "SameFile-NoFlags",
			args:    []string{file1.Name(), file1.Name()},
			wantErr: cp.ErrSkip,
		},
		{
			name:    "NoFlags-Fail-SrcNotExist",
			args:    []string{"src", filepath.Join(tmpDir, "destination")},
			wantErr: fs.ErrNotExist,
		},
		{
			name:    "NoFlags-Fail-DstcNotExist",
			args:    []string{file1.Name(), "dst"},
			wantErr: fs.ErrNotExist,
		},
		{
			name:    "NoFlags-Success-",
			args:    []string{file1.Name(), "dst", "src"},
			wantErr: unix.ENOTDIR,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			if err := run(tt.args, tt.recursiv, tt.ask, tt.force, tt.verbose, tt.noFollow, &out); err != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("%q failed. Error Got: %q, Error Want:%q", tt.name, err.Error(), tt.wantErr)
				}
				return
			}

			if err := IsEqualTree(cp.Default, tt.args[0], tt.args[1]); err != nil {
				t.Errorf("%q failed at: IsEqualTree failed with: %q", tt.name, err)
			}
		})

		t.Run(tt.name+"PreCallBack", func(t *testing.T) {
			var out bytes.Buffer
			f := setupPreCallback(tt.recursiv, tt.ask, tt.force, &out)
			srcfi, err := os.Stat(tt.args[0])
			// If the src file does not exist, there is no point in continue, but it is not an error so the say.
			// Also we catch that error in the previous test
			if errors.Is(err, fs.ErrNotExist) {
				return
			}
			if err := f(tt.args[0], tt.args[1], srcfi); errors.Is(err, cp.ErrSkip) {
				t.Log(err)
			}
		})

		t.Run(tt.name+"PostCallBack", func(t *testing.T) {
			var out bytes.Buffer
			f := setupPostCallback(tt.verbose, &out)
			f(tt.args[0], tt.args[1])
			if tt.verbose {
				if out.String() != fmt.Sprintf("%q -> %q\n", tt.args[0], tt.args[1]) {
					t.Errorf("PostCallback failed. Want: %q, Got: %q", fmt.Sprintf("%q -> %q\n", tt.args[0], tt.args[1]), out.String())
				}
			}
		})
	}
}

// TestCpSrcDirectory tests copying source to destination without recursive
// cmd-line equivalent: cp ~/dir ~/dir2
func TestCpSrcDirectory(t *testing.T) {
	flags.recursive = false
	defer resetFlags()

	tempDir := t.TempDir()
	tempDirTwo := t.TempDir()

	// capture log output to verify
	var logBytes bytes.Buffer

	if err := run([]string{tempDir, tempDirTwo}, false, false, false, false, false, &logBytes); err != nil {
		t.Fatalf("copy(%q -> %q) = %v, want nil", tempDir, tempDirTwo, err)
	}

	outString := fmt.Sprintf("cp: -r not specified, omitting directory %s", tempDir)
	capturedString := logBytes.String()
	if !strings.Contains(capturedString, outString) {
		t.Fatalf("copy(%q -> %q) = %v, want %v", tempDir, tempDirTwo, capturedString, outString)
	}
}

// TestCpRecursive tests the recursive mode copy
// Copy dir hierarchies src-dir to dst-dir
// whose src-dir and dst-dir already exists
// cmd-line equivalent: $ cp -R src-dir/ dst-dir/
func TestCpRecursive(t *testing.T) {
	flags.recursive = true
	defer resetFlags()

	tempDir := t.TempDir()

	srcDir := filepath.Join(tempDir, "src")
	if err := os.Mkdir(srcDir, 0o755); err != nil {
		t.Fatal(err)
	}
	dstDir := filepath.Join(tempDir, "dst-exists")
	if err := os.Mkdir(dstDir, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := createFilesTree(srcDir, maxDirDepth, 0); err != nil {
		t.Fatalf("cannot create files tree on directory %q: %v", srcDir, err)
	}

	t.Run("existing-dst-dir", func(t *testing.T) {
		var out bytes.Buffer
		if err := run([]string{srcDir, dstDir}, true, false, false, false, false, &out); err != nil {
			t.Fatalf("cp(%q -> %q) = %v, want nil", srcDir, dstDir, err)
		}
		// Because dstDir already existed, a new dir was created inside it.
		realDestination := filepath.Join(dstDir, filepath.Base(srcDir))
		if err := IsEqualTree(cp.Default, srcDir, realDestination); err != nil {
			t.Fatalf("copy(%q -> %q): file trees not equal: %v", srcDir, realDestination, err)
		}
	})

	t.Run("non-existing-dst-dir", func(t *testing.T) {
		var out bytes.Buffer
		notExistDstDir := filepath.Join(tempDir, "dst-does-not-exist")
		if err := run([]string{srcDir, notExistDstDir}, true, false, false, false, false, &out); err != nil {
			t.Fatalf("cp(%q -> %q) = %v, want nil", srcDir, notExistDstDir, err)
		}

		if err := IsEqualTree(cp.Default, srcDir, notExistDstDir); err != nil {
			t.Fatalf("copy(%q -> %q): file trees not equal: %v", srcDir, notExistDstDir, err)
		}
	})
}

// Other test to verify the CopyRecursive
// whose dir$n and dst-dir already exists
// cmd-line equivalent: $ cp -R dir1/ dir2/ dir3/ dst-dir/
//
// dst-dir will content dir{1, 3}
// $ dst-dir/
// ..	dir1/
// ..	dir2/
// ..   dir3/
func TestCpRecursiveMultiple(t *testing.T) {
	flags.recursive = true
	defer resetFlags()
	tempDir := t.TempDir()

	dstTest := filepath.Join(tempDir, "destination")
	if err := os.Mkdir(dstTest, 0o755); err != nil {
		t.Fatalf("Failed on build directory %v: %v", dstTest, err)
	}

	// create multiple random directories sources
	srcDirs := []string{}
	for i := 0; i < maxDirDepth; i++ {
		srcTest := t.TempDir()

		if err := createFilesTree(srcTest, maxDirDepth, 0); err != nil {
			t.Fatalf("cannot create files tree on directory %v: %v", srcTest, err)
		}

		srcDirs = append(srcDirs, srcTest)

	}
	var out bytes.Buffer
	args := srcDirs
	args = append(args, dstTest)
	if err := run(args, true, false, false, false, false, &out); err != nil {
		t.Fatalf("cp %q exit with error: %v", args, err)
	}
	// Make sure we can do it twice.
	flags.force = true
	if err := run(args, true, false, false, false, false, &out); err != nil {
		t.Fatalf("cp %q exit with error: %v", args, err)
	}
	for _, src := range srcDirs {
		_, srcFile := filepath.Split(src)

		dst := filepath.Join(dstTest, srcFile)
		if err := IsEqualTree(cp.Default, src, dst); err != nil {
			t.Fatalf("The copy %q -> %q failed: %v", src, dst, err)
		}
	}
}

// using -P don't follow symlinks, create other symlink
// cmd-line equivalent: $ cp -P symlink symlink-copy
func TestCpSymlink(t *testing.T) {
	tempDir := t.TempDir()

	f, err := randomFile(tempDir, "src-")
	if err != nil {
		t.Fatalf("cannot create a random file: %v", err)
	}
	defer f.Close()

	srcFpath := f.Name()
	srcFname := filepath.Base(srcFpath)

	newName := filepath.Join(tempDir, srcFname+"_link")
	if err := os.Symlink(srcFname, newName); err != nil {
		t.Fatalf("cannot create a link %q with target %q: %v", newName, srcFname, err)
	}

	t.Run("no-follow-symlink", func(t *testing.T) {
		defer resetFlags()
		var out bytes.Buffer
		flags.noFollowSymlinks = true

		dst := filepath.Join(tempDir, "dst-no-follow")
		if err := run([]string{newName, dst}, false, false, false, false, true, &out); err != nil {
			t.Fatalf("cp(%q -> %q) = %v, want nil", newName, dst, err)
		}
		if err := IsEqualTree(cp.NoFollowSymlinks, newName, dst); err != nil {
			t.Fatalf("The copy %q -> %q failed: %v", newName, dst, err)
		}
	})

	t.Run("follow-symlink", func(t *testing.T) {
		defer resetFlags()
		var out bytes.Buffer
		flags.noFollowSymlinks = false

		dst := filepath.Join(tempDir, "dst-follow")
		if err := run([]string{newName, dst}, false, false, false, false, false, &out); err != nil {
			t.Fatalf("cp(%q -> %q) = %v, want nil", newName, dst, err)
		}
		if err := IsEqualTree(cp.Default, newName, dst); err != nil {
			t.Fatalf("The copy %q -> %q failed: %v", newName, dst, err)
		}
	})
}

// isEqualFile compare two files by checksum
func isEqualFile(fpath1, fpath2 string) error {
	file1, err := os.Open(fpath1)
	if err != nil {
		return err
	}
	defer file1.Close()
	file2, err := os.Open(fpath2)
	if err != nil {
		return err
	}
	defer file2.Close()

	if !uio.ReaderAtEqual(file1, file2) {
		return fmt.Errorf("%q and %q do not have equal content", fpath1, fpath2)
	}
	return nil
}

func readDirNames(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var basenames []string
	for _, entry := range entries {
		basenames = append(basenames, entry.Name())
	}
	return basenames, nil
}

func stat(o cp.Options, path string) (os.FileInfo, error) {
	if o.NoFollowSymlinks {
		return os.Lstat(path)
	}
	return os.Stat(path)
}

// IsEqualTree compare the content in the file trees in src and dst paths
func IsEqualTree(o cp.Options, src, dst string) error {
	srcInfo, err := stat(o, src)
	if err != nil {
		return err
	}
	dstInfo, err := stat(o, dst)
	if err != nil {
		return err
	}
	if sm, dm := srcInfo.Mode()&os.ModeType, dstInfo.Mode()&os.ModeType; sm != dm {
		return fmt.Errorf("mismatched mode: %q has mode %s while %q has mode %s", src, sm, dst, dm)
	}

	switch {
	case srcInfo.Mode().IsDir():
		srcEntries, err := readDirNames(src)
		if err != nil {
			return err
		}
		dstEntries, err := readDirNames(dst)
		if err != nil {
			return err
		}
		// os.ReadDir guarantees these are sorted.
		if !reflect.DeepEqual(srcEntries, dstEntries) {
			return fmt.Errorf("directory contents did not match:\n%q had %v\n%q had %v", src, srcEntries, dst, dstEntries)
		}
		for _, basename := range srcEntries {
			if err := IsEqualTree(o, filepath.Join(src, basename), filepath.Join(dst, basename)); err != nil {
				return err
			}
		}
		return nil

	case srcInfo.Mode().IsRegular():
		return isEqualFile(src, dst)

	case srcInfo.Mode()&os.ModeSymlink == os.ModeSymlink:
		srcTarget, err := os.Readlink(src)
		if err != nil {
			return err
		}
		dstTarget, err := os.Readlink(dst)
		if err != nil {
			return err
		}
		if srcTarget != dstTarget {
			return fmt.Errorf("target mismatch: symlink %q had target %q, while %q had target %q", src, srcTarget, dst, dstTarget)
		}
		return nil

	default:
		return fmt.Errorf("unsupported mode: %s", srcInfo.Mode())
	}
}
