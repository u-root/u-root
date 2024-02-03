// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package initramfs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/cpio"
	"github.com/u-root/uio/uio"
)

func TestFilesAddFileNoFollow(t *testing.T) {
	regularFile, err := os.CreateTemp("", "archive-files-add-file")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(regularFile.Name())

	dir := t.TempDir()
	dir2 := t.TempDir()

	os.Create(filepath.Join(dir, "foo2"))
	os.Symlink(filepath.Join(dir, "foo2"), filepath.Join(dir2, "foo3"))

	for i, tt := range []struct {
		name        string
		af          *Files
		src         string
		dest        string
		result      *Files
		errContains string
	}{
		{
			name: "just add a file",
			af:   NewFiles(),

			src:  regularFile.Name(),
			dest: "bar/foo",

			result: &Files{
				Files: map[string]string{
					"bar/foo": regularFile.Name(),
				},
				Records: map[string]cpio.Record{},
			},
		},
		{
			name: "add symlinked file, NOT following",
			af:   NewFiles(),
			src:  filepath.Join(dir2, "foo3"),
			dest: "bar/foo",
			result: &Files{
				Files: map[string]string{
					"bar/foo": filepath.Join(dir2, "foo3"),
				},
				Records: map[string]cpio.Record{},
			},
		},
	} {
		t.Run(fmt.Sprintf("Test %02d: %s", i, tt.name), func(t *testing.T) {
			err := tt.af.AddFileNoFollow(tt.src, tt.dest)
			if err != nil && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("Error is %v, does not contain %v", err, tt.errContains)
			}
			if err == nil && len(tt.errContains) > 0 {
				t.Errorf("Got no error, want %v", tt.errContains)
			}

			if tt.result != nil && !reflect.DeepEqual(tt.af, tt.result) {
				t.Errorf("got %v, want %v", tt.af, tt.result)
			}
		})
	}
}

func TestFilesAddFile(t *testing.T) {
	regularFile, err := os.CreateTemp("", "archive-files-add-file")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(regularFile.Name())

	dir := t.TempDir()
	dir2 := t.TempDir()
	dir3 := t.TempDir()

	os.Create(filepath.Join(dir, "foo"))
	os.Create(filepath.Join(dir, "foo2"))
	os.Symlink(filepath.Join(dir, "foo2"), filepath.Join(dir2, "foo3"))

	fooDir := filepath.Join(dir3, "fooDir")
	os.Mkdir(fooDir, os.ModePerm)
	symlinkToDir3 := filepath.Join(dir3, "fooSymDir/")
	os.Symlink(fooDir, symlinkToDir3)
	os.Create(filepath.Join(fooDir, "foo"))
	os.Create(filepath.Join(fooDir, "bar"))

	for i, tt := range []struct {
		name        string
		af          *Files
		src         string
		dest        string
		result      *Files
		errContains string
	}{
		{
			name: "just add a file",
			af:   NewFiles(),

			src:  regularFile.Name(),
			dest: "bar/foo",

			result: &Files{
				Files: map[string]string{
					"bar/foo": regularFile.Name(),
				},
				Records: map[string]cpio.Record{},
			},
		},
		{
			name: "add symlinked file, following",
			af:   NewFiles(),
			src:  filepath.Join(dir2, "foo3"),
			dest: "bar/foo",
			result: &Files{
				Files: map[string]string{
					"bar/foo": filepath.Join(dir, "foo2"),
				},
				Records: map[string]cpio.Record{},
			},
		},
		{
			name: "add symlinked directory, following",
			af:   NewFiles(),
			src:  symlinkToDir3,
			dest: "foo/",
			result: &Files{
				Files: map[string]string{
					"foo":     fooDir,
					"foo/foo": filepath.Join(fooDir, "foo"),
					"foo/bar": filepath.Join(fooDir, "bar"),
				},
				Records: map[string]cpio.Record{},
			},
		},
		{
			name: "add file that exists in Files",
			af: &Files{
				Files: map[string]string{
					"bar/foo": "/some/other/place",
				},
			},
			src:  regularFile.Name(),
			dest: "bar/foo",
			result: &Files{
				Files: map[string]string{
					"bar/foo": "/some/other/place",
				},
			},
			errContains: "already exists in archive",
		},
		{
			name: "add a file that exists in Records",
			af: &Files{
				Records: map[string]cpio.Record{
					"bar/foo": cpio.Symlink("bar/foo", "/some/other/place"),
				},
			},
			src:  regularFile.Name(),
			dest: "bar/foo",
			result: &Files{
				Records: map[string]cpio.Record{
					"bar/foo": cpio.Symlink("bar/foo", "/some/other/place"),
				},
			},
			errContains: "already exists in archive",
		},
		{
			name: "add a file that already exists in Files, but is the same one",
			af: &Files{
				Files: map[string]string{
					"bar/foo": regularFile.Name(),
				},
			},
			src:  regularFile.Name(),
			dest: "bar/foo",
			result: &Files{
				Files: map[string]string{
					"bar/foo": regularFile.Name(),
				},
			},
		},
		{
			name: "absolute destination paths are made relative",
			af: &Files{
				Files: map[string]string{},
			},
			src:  dir,
			dest: "/bar/foo",
			result: &Files{
				Files: map[string]string{
					"bar/foo":      dir,
					"bar/foo/foo":  filepath.Join(dir, "foo"),
					"bar/foo/foo2": filepath.Join(dir, "foo2"),
				},
			},
		},
		{
			name: "add a directory",
			af: &Files{
				Files: map[string]string{},
			},
			src:  dir,
			dest: "bar/foo",
			result: &Files{
				Files: map[string]string{
					"bar/foo":      dir,
					"bar/foo/foo":  filepath.Join(dir, "foo"),
					"bar/foo/foo2": filepath.Join(dir, "foo2"),
				},
			},
		},
		{
			name: "add a different directory to the same destination, no overlapping children",
			af: &Files{
				Files: map[string]string{
					"bar/foo":     "/some/place/real",
					"bar/foo/zed": "/some/place/real/zed",
				},
			},
			src:  dir,
			dest: "bar/foo",
			result: &Files{
				Files: map[string]string{
					"bar/foo":      dir,
					"bar/foo/foo":  filepath.Join(dir, "foo"),
					"bar/foo/foo2": filepath.Join(dir, "foo2"),
					"bar/foo/zed":  "/some/place/real/zed",
				},
			},
		},
		{
			name: "add a different directory to the same destination, overlapping children",
			af: &Files{
				Files: map[string]string{
					"bar/foo":      "/some/place/real",
					"bar/foo/foo2": "/some/place/real/zed",
				},
			},
			src:         dir,
			dest:        "bar/foo",
			errContains: "already exists in archive",
		},
	} {
		t.Run(fmt.Sprintf("Test %02d: %s", i, tt.name), func(t *testing.T) {
			err := tt.af.AddFile(tt.src, tt.dest)
			if err != nil && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("Error is %v, does not contain %v", err, tt.errContains)
			}
			if err == nil && len(tt.errContains) > 0 {
				t.Errorf("Got no error, want %v", tt.errContains)
			}

			if tt.result != nil && !reflect.DeepEqual(tt.af, tt.result) {
				t.Errorf("got %v, want %v", tt.af, tt.result)
			}
		})
	}
}

func TestFilesAddRecord(t *testing.T) {
	for i, tt := range []struct {
		af     *Files
		record cpio.Record

		result      *Files
		errContains string
	}{
		{
			af:     NewFiles(),
			record: cpio.Symlink("bar/foo", ""),
			result: &Files{
				Files: map[string]string{},
				Records: map[string]cpio.Record{
					"bar/foo": cpio.Symlink("bar/foo", ""),
				},
			},
		},
		{
			af: &Files{
				Files: map[string]string{
					"bar/foo": "/some/other/place",
				},
			},
			record: cpio.Symlink("bar/foo", ""),
			result: &Files{
				Files: map[string]string{
					"bar/foo": "/some/other/place",
				},
			},
			errContains: "already exists in archive",
		},
		{
			af: &Files{
				Records: map[string]cpio.Record{
					"bar/foo": cpio.Symlink("bar/foo", "/some/other/place"),
				},
			},
			record: cpio.Symlink("bar/foo", ""),
			result: &Files{
				Records: map[string]cpio.Record{
					"bar/foo": cpio.Symlink("bar/foo", "/some/other/place"),
				},
			},
			errContains: "already exists in archive",
		},
		{
			af: &Files{
				Records: map[string]cpio.Record{
					"bar/foo": cpio.Symlink("bar/foo", "/some/other/place"),
				},
			},
			record: cpio.Symlink("bar/foo", "/some/other/place"),
			result: &Files{
				Records: map[string]cpio.Record{
					"bar/foo": cpio.Symlink("bar/foo", "/some/other/place"),
				},
			},
		},
		{
			record:      cpio.Symlink("/bar/foo", ""),
			errContains: "must not be absolute",
		},
	} {
		t.Run(fmt.Sprintf("Test %02d", i), func(t *testing.T) {
			err := tt.af.AddRecord(tt.record)
			if err != nil && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("Error is %v, does not contain %v", err, tt.errContains)
			}
			if err == nil && len(tt.errContains) > 0 {
				t.Errorf("Got no error, want %v", tt.errContains)
			}

			if !reflect.DeepEqual(tt.af, tt.result) {
				t.Errorf("got %v, want %v", tt.af, tt.result)
			}
		})
	}
}

func TestFilesfillInParent(t *testing.T) {
	for i, tt := range []struct {
		af     *Files
		result *Files
	}{
		{
			af: &Files{
				Records: map[string]cpio.Record{
					"foo/bar": cpio.Directory("foo/bar", 0o777),
				},
			},
			result: &Files{
				Records: map[string]cpio.Record{
					"foo/bar": cpio.Directory("foo/bar", 0o777),
					"foo":     cpio.Directory("foo", 0o755),
				},
			},
		},
		{
			af: &Files{
				Files: map[string]string{
					"baz/baz/baz": "/somewhere",
				},
				Records: map[string]cpio.Record{
					"foo/bar": cpio.Directory("foo/bar", 0o777),
				},
			},
			result: &Files{
				Files: map[string]string{
					"baz/baz/baz": "/somewhere",
				},
				Records: map[string]cpio.Record{
					"foo/bar": cpio.Directory("foo/bar", 0o777),
					"foo":     cpio.Directory("foo", 0o755),
					"baz":     cpio.Directory("baz", 0o755),
					"baz/baz": cpio.Directory("baz/baz", 0o755),
				},
			},
		},
		{
			af:     &Files{},
			result: &Files{},
		},
	} {
		t.Run(fmt.Sprintf("Test %02d", i), func(t *testing.T) {
			tt.af.fillInParents()
			if !reflect.DeepEqual(tt.af, tt.result) {
				t.Errorf("got %v, want %v", tt.af, tt.result)
			}
		})
	}
}

type MockArchiver struct {
	Records      Records
	FinishCalled bool
	BaseArchive  []cpio.Record
}

func (ma *MockArchiver) WriteRecord(r cpio.Record) error {
	if _, ok := ma.Records[r.Name]; ok {
		return fmt.Errorf("file exists")
	}
	ma.Records[r.Name] = r
	return nil
}

func (ma *MockArchiver) Finish() error {
	ma.FinishCalled = true
	return nil
}

func (ma *MockArchiver) ReadRecord() (cpio.Record, error) {
	if len(ma.BaseArchive) > 0 {
		next := ma.BaseArchive[0]
		ma.BaseArchive = ma.BaseArchive[1:]
		return next, nil
	}
	return cpio.Record{}, io.EOF
}

type Records map[string]cpio.Record

func RecordsEqual(r1, r2 Records, recordEqual func(cpio.Record, cpio.Record) bool) bool {
	for name, s1 := range r1 {
		s2, ok := r2[name]
		if !ok {
			return false
		}
		if !recordEqual(s1, s2) {
			return false
		}
	}
	for name := range r2 {
		if _, ok := r1[name]; !ok {
			return false
		}
	}
	return true
}

func sameNameModeContent(r1 cpio.Record, r2 cpio.Record) bool {
	if r1.Name != r2.Name || r1.Mode != r2.Mode {
		return false
	}
	return uio.ReaderAtEqual(r1.ReaderAt, r2.ReaderAt)
}

func TestOptsWrite(t *testing.T) {
	for i, tt := range []struct {
		desc string
		opts *Opts
		ma   *MockArchiver
		want Records
		err  error
	}{
		{
			desc: "no conflicts, just records",
			opts: &Opts{
				Files: &Files{
					Records: map[string]cpio.Record{
						"foo": cpio.Symlink("foo", "elsewhere"),
					},
				},
			},
			ma: &MockArchiver{
				Records: make(Records),
				BaseArchive: []cpio.Record{
					cpio.Directory("etc", 0o777),
					cpio.Directory("etc/nginx", 0o777),
				},
			},
			want: Records{
				"foo":       cpio.Symlink("foo", "elsewhere"),
				"etc":       cpio.Directory("etc", 0o777),
				"etc/nginx": cpio.Directory("etc/nginx", 0o777),
			},
		},
		{
			desc: "default already exists",
			opts: &Opts{
				Files: &Files{
					Records: map[string]cpio.Record{
						"etc": cpio.Symlink("etc", "whatever"),
					},
				},
			},
			ma: &MockArchiver{
				Records: make(Records),
				BaseArchive: []cpio.Record{
					cpio.Directory("etc", 0o777),
				},
			},
			want: Records{
				"etc": cpio.Symlink("etc", "whatever"),
			},
		},
		{
			desc: "no conflicts, missing parent automatically created",
			opts: &Opts{
				Files: &Files{
					Records: map[string]cpio.Record{
						"foo/bar/baz": cpio.Symlink("foo/bar/baz", "elsewhere"),
					},
				},
			},
			ma: &MockArchiver{
				Records: make(Records),
			},
			want: Records{
				"foo":         cpio.Directory("foo", 0o755),
				"foo/bar":     cpio.Directory("foo/bar", 0o755),
				"foo/bar/baz": cpio.Symlink("foo/bar/baz", "elsewhere"),
			},
		},
		{
			desc: "parent only automatically created if not already exists",
			opts: &Opts{
				Files: &Files{
					Records: map[string]cpio.Record{
						"foo/bar":     cpio.Directory("foo/bar", 0o444),
						"foo/bar/baz": cpio.Symlink("foo/bar/baz", "elsewhere"),
					},
				},
			},
			ma: &MockArchiver{
				Records: make(Records),
			},
			want: Records{
				"foo":         cpio.Directory("foo", 0o755),
				"foo/bar":     cpio.Directory("foo/bar", 0o444),
				"foo/bar/baz": cpio.Symlink("foo/bar/baz", "elsewhere"),
			},
		},
		{
			desc: "base archive",
			opts: &Opts{
				Files: &Files{
					Records: map[string]cpio.Record{
						"foo/bar": cpio.Symlink("foo/bar", "elsewhere"),
						"exists":  cpio.Directory("exists", 0o777),
					},
				},
			},
			ma: &MockArchiver{
				Records: make(Records),
				BaseArchive: []cpio.Record{
					cpio.Directory("etc", 0o755),
					cpio.Directory("foo", 0o444),
					cpio.Directory("exists", 0),
				},
			},
			want: Records{
				"etc":     cpio.Directory("etc", 0o755),
				"exists":  cpio.Directory("exists", 0o777),
				"foo":     cpio.Directory("foo", 0o444),
				"foo/bar": cpio.Symlink("foo/bar", "elsewhere"),
			},
		},
		{
			desc: "base archive with init, no user init",
			opts: &Opts{
				Files: &Files{
					Records: map[string]cpio.Record{},
				},
			},
			ma: &MockArchiver{
				Records: make(Records),
				BaseArchive: []cpio.Record{
					cpio.StaticFile("init", "boo", 0o555),
				},
			},
			want: Records{
				"init": cpio.StaticFile("init", "boo", 0o555),
			},
		},
		{
			desc: "base archive with init and user init",
			opts: &Opts{
				Files: &Files{
					Records: map[string]cpio.Record{
						"init": cpio.StaticFile("init", "bar", 0o444),
					},
				},
			},
			ma: &MockArchiver{
				Records: make(Records),
				BaseArchive: []cpio.Record{
					cpio.StaticFile("init", "boo", 0o555),
				},
			},
			want: Records{
				"init":  cpio.StaticFile("init", "bar", 0o444),
				"inito": cpio.StaticFile("inito", "boo", 0o555),
			},
		},
		{
			desc: "base archive with init, use existing init",
			opts: &Opts{
				Files: &Files{
					Records: map[string]cpio.Record{},
				},
				UseExistingInit: true,
			},
			ma: &MockArchiver{
				Records: make(Records),
				BaseArchive: []cpio.Record{
					cpio.StaticFile("init", "boo", 0o555),
				},
			},
			want: Records{
				"init": cpio.StaticFile("init", "boo", 0o555),
			},
		},
		{
			desc: "base archive with init and user init, use existing init",
			opts: &Opts{
				Files: &Files{
					Records: map[string]cpio.Record{
						"init": cpio.StaticFile("init", "huh", 0o111),
					},
				},
				UseExistingInit: true,
			},
			ma: &MockArchiver{
				Records: make(Records),
				BaseArchive: []cpio.Record{
					cpio.StaticFile("init", "boo", 0o555),
				},
			},
			want: Records{
				"init":  cpio.StaticFile("init", "boo", 0o555),
				"inito": cpio.StaticFile("inito", "huh", 0o111),
			},
		},
	} {
		t.Run(fmt.Sprintf("Test %02d (%s)", i, tt.desc), func(t *testing.T) {
			tt.opts.BaseArchive = tt.ma
			tt.opts.OutputFile = tt.ma

			if err := Write(tt.opts); err != tt.err {
				t.Errorf("Write() = %v, want %v", err, tt.err)
			} else if err == nil && !tt.ma.FinishCalled {
				t.Errorf("Finish wasn't called on archive")
			}

			if !RecordsEqual(tt.ma.Records, tt.want, sameNameModeContent) {
				t.Errorf("Write() = %v, want %v", tt.ma.Records, tt.want)
			}
		})
	}
}
