// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/core/find"
)

func prepareDirLayout(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	_, err := os.Create(filepath.Join(tmpDir, "file1"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Create(filepath.Join(tmpDir, "file2"))
	if err != nil {
		t.Fatal(err)
	}
	err = os.Mkdir(filepath.Join(tmpDir, "dir1"), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Create(filepath.Join(tmpDir, "dir1/file1"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Create(filepath.Join(tmpDir, "dir1/file2"))
	if err != nil {
		t.Fatal(err)
	}
	err = os.Mkdir(filepath.Join(tmpDir, "dir2"), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Create(filepath.Join(tmpDir, "dir2/file1"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Create(filepath.Join(tmpDir, "dir2/file3"))
	if err != nil {
		t.Fatal(err)
	}
	return tmpDir
}

func TestFind(t *testing.T) {
	tmpDir := prepareDirLayout(t)

	tests := []struct {
		wantStdout string
		wantErr    bool
		args       []string
	}{
		{
			wantStdout: filepath.Join(tmpDir, "file1") + "\n",
			args:       []string{filepath.Join(tmpDir, "file1")},
		},
		{
			wantStdout: filepath.Join(tmpDir, "dir1") + "\n" + filepath.Join(tmpDir, "dir1/file1") + "\n" + filepath.Join(tmpDir, "dir1/file2") + "\n",
			args:       []string{filepath.Join(tmpDir, "dir1")},
		},
		{
			wantStdout: filepath.Join(tmpDir, "dir1/file1") + "\n" + filepath.Join(tmpDir, "dir2/file1") + "\n" + filepath.Join(tmpDir, "file1") + "\n",
			args:       []string{"-name", "file1", tmpDir},
		},
		{
			wantStdout: tmpDir + "\n" + filepath.Join(tmpDir, "dir1") + "\n" + filepath.Join(tmpDir, "dir2") + "\n",
			args:       []string{"-type", "d", tmpDir},
		},
		{
			wantStdout: tmpDir + "\n" + filepath.Join(tmpDir, "dir1") + "\n" + filepath.Join(tmpDir, "dir2") + "\n",
			args:       []string{"-type", "directory", tmpDir},
		},
		{
			wantStdout: filepath.Join(tmpDir, "dir1/file1") + "\n" + filepath.Join(tmpDir, "dir1/file2") + "\n" + filepath.Join(tmpDir, "dir2/file1") + "\n" + filepath.Join(tmpDir, "dir2/file3") + "\n" + filepath.Join(tmpDir, "file1") + "\n" + filepath.Join(tmpDir, "file2") + "\n",
			args:       []string{"-type", "f", tmpDir},
		},
		{
			wantStdout: filepath.Join(tmpDir, "dir1/file1") + "\n" + filepath.Join(tmpDir, "dir1/file2") + "\n" + filepath.Join(tmpDir, "dir2/file1") + "\n" + filepath.Join(tmpDir, "dir2/file3") + "\n" + filepath.Join(tmpDir, "file1") + "\n" + filepath.Join(tmpDir, "file2") + "\n",
			args:       []string{"-type", "file", tmpDir},
		},
		{
			args:    []string{"-type", "notvalid", tmpDir},
			wantErr: true,
		},
		{
			wantStdout: filepath.Join(tmpDir, "file1") + "\n",
			args:       []string{"-mode", "420", filepath.Join(tmpDir, "file1")}, // 420 decimal = 644 octal
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			cmd := find.New()
			cmd.SetIO(nil, &stdout, &stderr)
			err := cmd.Run(tt.args...)
			if tt.wantErr {
				if err == nil {
					t.Fatal("want error got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("want nil got %v", err)
			}

			resStdout := stdout.String()
			if resStdout != tt.wantStdout {
				t.Errorf("args: %v\nwant\n %s, got\n %s", tt.args, tt.wantStdout, resStdout)
			}
		})
	}
}

func TestFindLong(t *testing.T) {
	tmpDir := prepareDirLayout(t)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := find.New()
	cmd.SetIO(nil, &stdout, &stderr)
	err := cmd.Run("-l", filepath.Join(tmpDir, "file1"))
	if err != nil {
		t.Fatal(err)
	}

	res := strings.TrimSpace(stdout.String())

	if !strings.HasPrefix(res, "-rw-r--r--") {
		t.Errorf("want prefix: -rw-r--r--, got prefix: %s", res[:10])
	}
	if !strings.HasSuffix(res, "file1") {
		t.Errorf("want suffix: file1, got suffix: %s", res[len(res)-5:])
	}
}
