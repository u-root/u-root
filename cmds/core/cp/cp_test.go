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

	"github.com/u-root/u-root/pkg/core/cp"
	"github.com/u-root/uio/uio"
)

const (
	maxSizeFile = 1000
	maxDirDepth = 5
	maxFiles    = 5
)

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
	for i := range maxFiles {
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
		name    string
		args    []string
		input   string
		wantErr error
	}{
		{
			name: "NoFlags-Success-",
			args: []string{file1.Name(), filepath.Join(tmpDir, "destination")},
		},
		{
			name:  "AskYes-Success-",
			args:  []string{"-i", file1.Name(), filepath.Join(tmpDir, "destination")},
			input: "yes\n",
		},
		{
			name:  "AskNo-Skip-",
			args:  []string{"-i", file1.Name(), filepath.Join(tmpDir, "destination")},
			input: "no\n",
		},
		{
			name: "Verbose",
			args: []string{"-v", file1.Name(), filepath.Join(tmpDir, "destination")},
		},
		{
			name: "SameFile-NoFlags",
			args: []string{file1.Name(), file1.Name()},
		},
		{
			name:    "NoFlags-Fail-SrcNotExist",
			args:    []string{"src", filepath.Join(tmpDir, "destination")},
			wantErr: fs.ErrNotExist,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			cmd := cp.New()
			var stdout, stderr bytes.Buffer
			cmd.SetIO(strings.NewReader(tt.input), &stdout, &stderr)

			err := cmd.Run(tt.args...)
			if tt.wantErr != nil {
				if err == nil || !errors.Is(err, tt.wantErr) {
					t.Errorf("Run() = %v, want error %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("Run() = %v, want nil", err)
			}
		})
	}
}

// TestCpSrcDirectory tests copying source to destination without recursive
// cmd-line equivalent: cp ~/dir ~/dir2
func TestCpSrcDirectory(t *testing.T) {
	tempDir := t.TempDir()
	tempDirTwo := t.TempDir()

	cmd := cp.New()
	var stdout, stderr bytes.Buffer
	var stdin bytes.Buffer
	cmd.SetIO(&stdin, &stdout, &stderr)

	err := cmd.Run(tempDir, tempDirTwo)
	if err != nil {
		t.Fatalf("Run() = %v, want nil", err)
	}

	outString := fmt.Sprintf("cp: -r not specified, omitting directory %s", tempDir)
	capturedString := stderr.String()
	if !strings.Contains(capturedString, outString) {
		t.Fatal("strings.Contains(capturedString, outString) = false, not true")
	}
}

// TestCpRecursive tests the recursive mode copy
// Copy dir hierarchies src-dir to dst-dir
// whose src-dir and dst-dir already exists
// cmd-line equivalent: $ cp -R src-dir/ dst-dir/
func TestCpRecursive(t *testing.T) {
	tempDir := t.TempDir()

	srcDir := filepath.Join(tempDir, "src")
	if err := os.Mkdir(srcDir, 0o755); err != nil {
		t.Fatalf("os.Mkdir(srcDir, 0o755) = %v, want nil", err)
	}
	dstDir := filepath.Join(tempDir, "dst-exists")
	if err := os.Mkdir(dstDir, 0o755); err != nil {
		t.Fatalf("os.Mkdir(dstDir, 0o755) = %v, want nil", err)
	}

	if err := createFilesTree(srcDir, maxDirDepth, 0); err != nil {
		t.Fatalf("createFilesTree(srcDir, maxDirDepth, 0) = %v, want nil", err)
	}

	t.Run("existing-dst-dir", func(t *testing.T) {
		cmd := cp.New()
		var stdout, stderr bytes.Buffer
		var stdin bytes.Buffer
		cmd.SetIO(&stdin, &stdout, &stderr)

		err := cmd.Run("-r", srcDir, dstDir)
		if err != nil {
			t.Fatalf("Run() = %v, want nil", err)
		}

		// Because dstDir already existed, a new dir was created inside it.
		realDestination := filepath.Join(dstDir, filepath.Base(srcDir))
		if err := IsEqualTree(cp.Options{}, srcDir, realDestination); err != nil {
			t.Fatalf("IsEqualTree() = %v, want nil", err)
		}
	})

	t.Run("non-existing-dst-dir", func(t *testing.T) {
		cmd := cp.New()
		var stdout, stderr bytes.Buffer
		var stdin bytes.Buffer
		cmd.SetIO(&stdin, &stdout, &stderr)

		notExistDstDir := filepath.Join(tempDir, "dst-does-not-exist")
		err := cmd.Run("-r", srcDir, notExistDstDir)
		if err != nil {
			t.Fatalf("Run() = %v, want nil", err)
		}

		if err := IsEqualTree(cp.Options{}, srcDir, notExistDstDir); err != nil {
			t.Fatalf("IsEqualTree() = %v, want nil", err)
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
