// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ls

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Test listName func
func TestListName(t *testing.T) {
	// Create some directorys.
	tmpDir := t.TempDir()
	if err := os.Mkdir(filepath.Join(tmpDir, "d1"), 0777); err != nil {
		t.Fatalf("err in os.Mkdir: %v", err)
	}
	// Create some files.
	files := []string{"f1", "f2", "f3\nline 2", ".f4", "d1/f4"}
	for _, file := range files {
		if _, err := os.Create(filepath.Join(tmpDir, file)); err != nil {
			t.Fatalf("err in os.Create: %v", err)
		}
	}

	// Creating test table
	for _, tt := range []struct {
		name   string
		input  string
		want   string
		flag   cmd
		prefix bool
	}{
		{
			name:  "ls without arguments",
			input: tmpDir,
			want:  fmt.Sprintf("%s\n%s\n%s\n%s\n", "d1", "f1", "f2", "f3?line 2"),
			flag:  cmd{},
		},
		{
			name:  "ls osfi.IsDir() path, quoted = true",
			input: tmpDir,
			want:  fmt.Sprintf("\"%s\"\n\"%s\"\n\"%s\"\n\"%s\"\n", "d1", "f1", "f2", "f3\\nline 2"),
			flag: cmd{
				quoted: true,
			},
		},
		{
			name:  "ls osfi.IsDir() path, quoted = true, prefix = true ",
			input: tmpDir,
			want:  fmt.Sprintf("\"%s\":\n\"%s\"\n\"%s\"\n\"%s\"\n\"%s\"\n", tmpDir, "d1", "f1", "f2", "f3\\nline 2"),
			flag: cmd{
				quoted: true,
			},
			prefix: true,
		},
		{
			name:   "ls osfi.IsDir() path, quoted = false, prefix = true ",
			input:  tmpDir,
			want:   fmt.Sprintf("%s:\n%s\n%s\n%s\n%s\n", tmpDir, "d1", "f1", "f2", "f3?line 2"),
			prefix: true,
		},
		{
			name:  "ls recurse = true",
			input: tmpDir,
			want: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n", tmpDir, filepath.Join(tmpDir, ".f4"), filepath.Join(tmpDir, "d1"),
				filepath.Join(tmpDir, "d1/f4"), filepath.Join(tmpDir, "f1"), filepath.Join(tmpDir, "f2"), filepath.Join(tmpDir, "f3?line 2")),
			flag: cmd{
				all:     true,
				recurse: true,
			},
		},
		{
			name:  "ls directory = true",
			input: tmpDir,
			want:  fmt.Sprintf("%s\n", "tmpDir"),
			flag: cmd{
				directory: true,
			},
		},
		{
			name:  "ls classify = true",
			input: tmpDir,
			want:  fmt.Sprintf("%s\n%s\n%s\n%s\n", "d1/", "f1", "f2", "f3?line 2"),
			flag: cmd{
				classify: true,
			},
		},
		{
			name:  "file does not exist",
			input: "dir",
			flag:  cmd{},
		},
	} {
		// Running the tests
		t.Run(tt.name, func(t *testing.T) {
			// Write output in buffer.
			var buf bytes.Buffer
			var s Stringer = NameStringer{}
			if tt.flag.quoted {
				s = QuotedStringer{}
			}
			if tt.flag.long {
				s = LongStringer{Human: tt.flag.human, Name: s}
			}
			tt.flag.stdout = &buf
			if err := tt.flag.listName(s, tt.input, tt.prefix); err != nil {
				if buf.String() != tt.want {
					t.Errorf("listName() = '%v', want: '%v'", buf.String(), tt.want)
				}
			}
		})
	}
}

// Test list func
func TestRun(t *testing.T) {
	// Creating test table
	for _, tt := range []struct {
		name   string
		args   []string
		err    error
		prefix bool
	}{
		{
			name: "exist, input empty, quoted = true, long = true",
			args: []string{"ls", "-Ql"},
			err:  nil,
		},
		{
			name: "not exist, input empty, quoted = true, long = true",
			args: []string{"ls", "-Ql", "dir"},
			err:  os.ErrNotExist,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			wd, _ := os.Getwd()
			stderr := &bytes.Buffer{}
			if code := Command(io.Discard, stderr, wd).Run(tt.args...); code != 0 {
				stderrs := stderr.String()
				tterrs := fmt.Sprintln(tt.err)
				if tterrs != stderrs {
					t.Errorf("list() = %#v, want: %#v", stderrs, tterrs)
				}
			} else if tt.err != nil {
				t.Errorf("expected error \"%v\", got nil", tt.err)
			}
		})
	}
}

// Test indicator func
func TestIndicator(t *testing.T) {
	// Creating test table
	for _, test := range []struct {
		lsInfo FileInfo
		symbol string
	}{
		{
			FileInfo{
				Mode: os.ModeDir,
			},
			"/",
		},
		{
			FileInfo{
				Mode: os.ModeNamedPipe,
			},
			"|",
		},
		{
			FileInfo{
				Mode: os.ModeSymlink,
			},
			"@",
		},
		{
			FileInfo{
				Mode: os.ModeSocket,
			},
			"=",
		},
		{
			FileInfo{
				Mode: 0b110110100,
			},
			"",
		},
		{
			FileInfo{
				Mode: 0b111111101,
			},
			"*",
		},
	} {
		// Run tests
		got := indicator(test.lsInfo)
		if got != test.symbol {
			t.Errorf("for mode '%b' expected '%q', got '%q'", test.lsInfo.Mode, test.symbol, got)
		}
	}
}

// Make sure if perms fail in a dir, we still list the dir.
func TestPermHandling(t *testing.T) {
	d := t.TempDir()
	for _, v := range []string{"a", "c", "d"} {
		if err := os.Mkdir(filepath.Join(d, v), 0777); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Mkdir(filepath.Join(d, "b"), 0); err != nil {
		t.Fatal(err)
	}
	for _, v := range []string{"0", "1", "2"} {
		if err := os.Mkdir(filepath.Join(d, v), 0777); err != nil {
			t.Fatal(err)
		}
	}
	b := &bytes.Buffer{}
	var c = cmd{stdout: b}

	if err := c.listName(NameStringer{}, d, false); err != nil {
		t.Fatalf("listName(NameString{}, %q, w, false): %v != nil", d, err)
	}
	// the output varies very widely between kernels and Go versions :-(
	// Just look for 'permission denied' and more than 6 lines of output ...
	if !strings.Contains(b.String(), "0\n1\n2\na\nb\nc\nd\n") {
		t.Errorf("ls %q: output %q did not contain %q", d, b.String(), "0\n1\n2\na\nb\nc\nd\n")
	}
}
