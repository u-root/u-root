package uroot

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/cpio"
	"golang.org/x/sys/unix"
)

func TestArchiveFilesAddFile(t *testing.T) {
	for i, tt := range []struct {
		af          ArchiveFiles
		src         string
		dest        string
		result      ArchiveFiles
		errContains string
	}{
		{
			af:   NewArchiveFiles(),
			src:  "/foo/bar",
			dest: "bar/foo",

			result: ArchiveFiles{
				Files: map[string]string{
					"bar/foo": "/foo/bar",
				},
				Records: map[string]cpio.Record{},
			},
		},
		{
			af: ArchiveFiles{
				Files: map[string]string{
					"bar/foo": "/some/other/place",
				},
			},
			src:  "/foo/bar",
			dest: "bar/foo",
			result: ArchiveFiles{
				Files: map[string]string{
					"bar/foo": "/some/other/place",
				},
			},
			errContains: "already exists in archive",
		},
		{
			af: ArchiveFiles{
				Records: map[string]cpio.Record{
					"bar/foo": cpio.Symlink("bar/foo", "/some/other/place"),
				},
			},
			src:  "/foo/bar",
			dest: "bar/foo",
			result: ArchiveFiles{
				Records: map[string]cpio.Record{
					"bar/foo": cpio.Symlink("bar/foo", "/some/other/place"),
				},
			},
			errContains: "already exists in archive",
		},
		{
			af: ArchiveFiles{
				Files: map[string]string{
					"bar/foo": "/foo/bar",
				},
			},
			src:  "/foo/bar",
			dest: "bar/foo",
			result: ArchiveFiles{
				Files: map[string]string{
					"bar/foo": "/foo/bar",
				},
			},
		},
		{
			src:         "/foo/bar",
			dest:        "/bar/foo",
			errContains: "must not be absolute",
		},
		{
			src:         "foo/bar",
			dest:        "bar/foo",
			errContains: "must be absolute",
		},
	} {
		t.Run(fmt.Sprintf("Test %02d", i), func(t *testing.T) {
			err := tt.af.AddFile(tt.src, tt.dest)
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

func TestArchiveFilesAddRecord(t *testing.T) {
	for i, tt := range []struct {
		af     ArchiveFiles
		record cpio.Record

		result      ArchiveFiles
		errContains string
	}{
		{
			af:     NewArchiveFiles(),
			record: cpio.Symlink("bar/foo", ""),
			result: ArchiveFiles{
				Files: map[string]string{},
				Records: map[string]cpio.Record{
					"bar/foo": cpio.Symlink("bar/foo", ""),
				},
			},
		},
		{
			af: ArchiveFiles{
				Files: map[string]string{
					"bar/foo": "/some/other/place",
				},
			},
			record: cpio.Symlink("bar/foo", ""),
			result: ArchiveFiles{
				Files: map[string]string{
					"bar/foo": "/some/other/place",
				},
			},
			errContains: "already exists in archive",
		},
		{
			af: ArchiveFiles{
				Records: map[string]cpio.Record{
					"bar/foo": cpio.Symlink("bar/foo", "/some/other/place"),
				},
			},
			record: cpio.Symlink("bar/foo", ""),
			result: ArchiveFiles{
				Records: map[string]cpio.Record{
					"bar/foo": cpio.Symlink("bar/foo", "/some/other/place"),
				},
			},
			errContains: "already exists in archive",
		},
		{
			af: ArchiveFiles{
				Records: map[string]cpio.Record{
					"bar/foo": cpio.Symlink("bar/foo", "/some/other/place"),
				},
			},
			record: cpio.Symlink("bar/foo", "/some/other/place"),
			result: ArchiveFiles{
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

func TestArchiveFilesfillInParent(t *testing.T) {
	for i, tt := range []struct {
		af     ArchiveFiles
		result ArchiveFiles
	}{
		{
			af: ArchiveFiles{
				Records: map[string]cpio.Record{
					"foo/bar": cpio.Directory("foo/bar", 0777),
				},
			},
			result: ArchiveFiles{
				Records: map[string]cpio.Record{
					"foo/bar": cpio.Directory("foo/bar", 0777),
					"foo":     cpio.Directory("foo", 0755),
				},
			},
		},
		{
			af: ArchiveFiles{
				Files: map[string]string{
					"baz/baz/baz": "/somewhere",
				},
				Records: map[string]cpio.Record{
					"foo/bar": cpio.Directory("foo/bar", 0777),
				},
			},
			result: ArchiveFiles{
				Files: map[string]string{
					"baz/baz/baz": "/somewhere",
				},
				Records: map[string]cpio.Record{
					"foo/bar": cpio.Directory("foo/bar", 0777),
					"foo":     cpio.Directory("foo", 0755),
					"baz":     cpio.Directory("baz", 0755),
					"baz/baz": cpio.Directory("baz/baz", 0755),
				},
			},
		},
		{
			af:     ArchiveFiles{},
			result: ArchiveFiles{},
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
	return cpio.ReadCloserEqual(r1.ReadCloser, r2.ReadCloser)
}

func TestWriteFile(t *testing.T) {
	unix.Umask(0)

	for i, tt := range []struct {
		ma   *MockArchiver
		src  func() string
		dest string
		err  error
		want Records
	}{
		{
			ma: &MockArchiver{
				Records: make(Records),
			},
			src: func() string {
				f, err := ioutil.TempFile("", "foo")
				if err != nil {
					panic(err)
				}
				n := f.Name()
				f.Close()
				return n
			},
			dest: "foo/whatever",
			want: Records{
				"foo/whatever": cpio.Record{
					Info: cpio.Info{
						Name:  "foo/whatever",
						Mode:  unix.S_IFREG | 0600,
						UID:   uint64(os.Geteuid()),
						GID:   uint64(os.Getegid()),
						NLink: 1,
						Major: 253,
						Minor: 1,
					},
				},
			},
		},
		{
			ma: &MockArchiver{
				Records: make(Records),
			},
			src: func() string {
				f, err := ioutil.TempDir("", "foo")
				if err != nil {
					panic(err)
				}
				if err := ioutil.WriteFile(filepath.Join(f, "bla"), []byte("foo"), 0644); err != nil {
					panic(err)
				}
				if err := ioutil.WriteFile(filepath.Join(f, "bla2"), []byte("foo2"), 0644); err != nil {
					panic(err)
				}
				return f
			},
			dest: "etc",
			want: Records{
				"etc": cpio.Record{
					Info: cpio.Info{
						Name: "etc",
						Mode: unix.S_IFDIR | 0700,
					},
				},
				"etc/bla": cpio.Record{
					Info: cpio.Info{
						Name: "etc/bla",
						Mode: unix.S_IFREG | 0644,
					},
					ReadCloser: ioutil.NopCloser(bytes.NewReader([]byte("foo"))),
				},
				"etc/bla2": cpio.Record{
					Info: cpio.Info{
						Name: "etc/bla2",
						Mode: unix.S_IFREG | 0644,
					},
					ReadCloser: ioutil.NopCloser(bytes.NewReader([]byte("foo2"))),
				},
			},
		},
	} {
		t.Run(fmt.Sprintf("Test %02d", i), func(t *testing.T) {
			src := tt.src()
			defer os.RemoveAll(src)
			if err := WriteFile(tt.ma, src, tt.dest); err != tt.err {
				t.Errorf("WriteFile() = %v, want %v", err, tt.err)
			}
			if !RecordsEqual(tt.ma.Records, tt.want, sameNameModeContent) {
				t.Errorf("WriteFile() = %v, want %v", tt.ma.Records, tt.want)
			}
		})
	}
}

func TestArchiveOptsWrite(t *testing.T) {
	for i, tt := range []struct {
		desc string
		opts *ArchiveOpts
		ma   *MockArchiver
		want Records
		err  error
	}{
		{
			desc: "no conflicts, just records",
			opts: &ArchiveOpts{
				ArchiveFiles: ArchiveFiles{
					Records: map[string]cpio.Record{
						"foo": cpio.Symlink("foo", "elsewhere"),
					},
				},
				DefaultRecords: []cpio.Record{
					cpio.Directory("etc", 0777),
					cpio.Directory("etc/nginx", 0777),
				},
			},
			ma: &MockArchiver{
				Records: make(Records),
			},
			want: Records{
				"foo":       cpio.Symlink("foo", "elsewhere"),
				"etc":       cpio.Directory("etc", 0777),
				"etc/nginx": cpio.Directory("etc/nginx", 0777),
			},
		},
		{
			desc: "default already exists",
			opts: &ArchiveOpts{
				ArchiveFiles: ArchiveFiles{
					Records: map[string]cpio.Record{
						"etc": cpio.Symlink("etc", "whatever"),
					},
				},
				DefaultRecords: []cpio.Record{
					cpio.Directory("etc", 0777),
				},
			},
			ma: &MockArchiver{
				Records: make(Records),
			},
			want: Records{
				"etc": cpio.Symlink("etc", "whatever"),
			},
		},
		{
			desc: "no conflicts, missing parent automatically created",
			opts: &ArchiveOpts{
				ArchiveFiles: ArchiveFiles{
					Records: map[string]cpio.Record{
						"foo/bar/baz": cpio.Symlink("foo/bar/baz", "elsewhere"),
					},
				},
			},
			ma: &MockArchiver{
				Records: make(Records),
			},
			want: Records{
				"foo":         cpio.Directory("foo", 0755),
				"foo/bar":     cpio.Directory("foo/bar", 0755),
				"foo/bar/baz": cpio.Symlink("foo/bar/baz", "elsewhere"),
			},
		},
		{
			desc: "parent only automatically created if not already exists",
			opts: &ArchiveOpts{
				ArchiveFiles: ArchiveFiles{
					Records: map[string]cpio.Record{
						"foo/bar":     cpio.Directory("foo/bar", 0444),
						"foo/bar/baz": cpio.Symlink("foo/bar/baz", "elsewhere"),
					},
				},
			},
			ma: &MockArchiver{
				Records: make(Records),
			},
			want: Records{
				"foo":         cpio.Directory("foo", 0755),
				"foo/bar":     cpio.Directory("foo/bar", 0444),
				"foo/bar/baz": cpio.Symlink("foo/bar/baz", "elsewhere"),
			},
		},
		{
			desc: "base archive",
			opts: &ArchiveOpts{
				ArchiveFiles: ArchiveFiles{
					Records: map[string]cpio.Record{
						"foo/bar": cpio.Symlink("foo/bar", "elsewhere"),
						"exists":  cpio.Directory("exists", 0777),
					},
				},
			},
			ma: &MockArchiver{
				Records: make(Records),
				BaseArchive: []cpio.Record{
					cpio.Directory("etc", 0755),
					cpio.Directory("foo", 0444),
					cpio.Directory("exists", 0),
				},
			},
			want: Records{
				"etc":     cpio.Directory("etc", 0755),
				"exists":  cpio.Directory("exists", 0777),
				"foo":     cpio.Directory("foo", 0444),
				"foo/bar": cpio.Symlink("foo/bar", "elsewhere"),
			},
		},
		{
			desc: "base archive with init, no user init",
			opts: &ArchiveOpts{
				ArchiveFiles: ArchiveFiles{
					Records: map[string]cpio.Record{},
				},
			},
			ma: &MockArchiver{
				Records: make(Records),
				BaseArchive: []cpio.Record{
					cpio.StaticFile("init", "boo", 0555),
				},
			},
			want: Records{
				"init": cpio.StaticFile("init", "boo", 0555),
			},
		},
		{
			desc: "base archive with init and user init",
			opts: &ArchiveOpts{
				ArchiveFiles: ArchiveFiles{
					Records: map[string]cpio.Record{
						"init": cpio.StaticFile("init", "bar", 0444),
					},
				},
			},
			ma: &MockArchiver{
				Records: make(Records),
				BaseArchive: []cpio.Record{
					cpio.StaticFile("init", "boo", 0555),
				},
			},
			want: Records{
				"init":  cpio.StaticFile("init", "bar", 0444),
				"inito": cpio.StaticFile("inito", "boo", 0555),
			},
		},
		{
			desc: "base archive with init, use existing init",
			opts: &ArchiveOpts{
				ArchiveFiles: ArchiveFiles{
					Records: map[string]cpio.Record{},
				},
				UseExistingInit: true,
			},
			ma: &MockArchiver{
				Records: make(Records),
				BaseArchive: []cpio.Record{
					cpio.StaticFile("init", "boo", 0555),
				},
			},
			want: Records{
				"init": cpio.StaticFile("init", "boo", 0555),
			},
		},
		{
			desc: "base archive with init and user init, use existing init",
			opts: &ArchiveOpts{
				ArchiveFiles: ArchiveFiles{
					Records: map[string]cpio.Record{
						"init": cpio.StaticFile("init", "huh", 0111),
					},
				},
				UseExistingInit: true,
			},
			ma: &MockArchiver{
				Records: make(Records),
				BaseArchive: []cpio.Record{
					cpio.StaticFile("init", "boo", 0555),
				},
			},
			want: Records{
				"init":  cpio.StaticFile("init", "boo", 0555),
				"inito": cpio.StaticFile("inito", "huh", 0111),
			},
		},
	} {
		t.Run(fmt.Sprintf("Test %02d (%s)", i, tt.desc), func(t *testing.T) {
			tt.opts.BaseArchive = tt.ma
			tt.opts.OutputFile = tt.ma

			if err := tt.opts.Write(); err != tt.err {
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
