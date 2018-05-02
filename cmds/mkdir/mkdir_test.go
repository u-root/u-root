package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"syscall"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

var (
	umaskDefault = 022
)

func run(c *exec.Cmd) (string, string, error) {
	var o, e bytes.Buffer
	c.Stdout, c.Stderr = &o, &e
	err := c.Run()
	return o.String(), e.String(), err
}

func findFile(dir string, filename string) (os.FileInfo, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.Name() == filename {
			return file, nil
		}
	}
	return nil, nil
}

func TestMkdirErrors(t *testing.T) {
	syscall.Umask(umaskDefault)

	// Error Tests
	for _, test := range []struct {
		name string
		args []string
		want string
	}{
		{
			name: "No Arg Error",
			args: nil,
			want: "Usage",
		},
		{
			name: "Perm Mode Bits over 7 Error",
			args: []string{"-m=7778", "foo"},
			want: `invalid mode '7778'`,
		},
		{
			name: "More than 4 Perm Mode Bits Error",
			args: []string{"-m=11111", "foo"},
			want: `invalid mode '11111'`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			tmpDir, err := ioutil.TempDir("", "mkdir")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tmpDir)

			c := testutil.Command(t, test.args...)
			c.Dir = tmpDir
			_, stderr, err := run(c)
			if err == nil || !strings.Contains(stderr, test.want) {
				t.Fatalf("stderr is %v, want %v", stderr, test.want)
			}

			f, err := findFile(tmpDir, "foo")
			if err != nil {
				t.Fatalf("error finding file \"foo\": %v", err)
			}
			if f != nil {
				t.Errorf("wantected not to find file, but found file \"foo\"")
			}
		})
	}
}

func TestMkdirRegular(t *testing.T) {
	syscall.Umask(umaskDefault)

	for _, test := range []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "Create 1 Directory",
			args: []string{"foo"},
			want: []string{"foo"},
		},
		{
			name: "Create 2 Directories",
			args: []string{"foo", "bar"},
			want: []string{"foo", "bar"},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			tmpDir, err := ioutil.TempDir("", "mkdir")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tmpDir)

			c := testutil.Command(t, test.args...)
			c.Dir = tmpDir
			if err := c.Run(); err != nil {
				t.Fatalf("Error running mkdir: %v", err)
			}
			for _, dirName := range test.want {
				f, err := findFile(tmpDir, dirName)
				if err != nil {
					t.Fatal(err)
				}
				if f == nil {
					t.Errorf("file %q not found in dir", dirName)
				}
			}
		})
	}
}

func TestMkdirPermission(t *testing.T) {
	syscall.Umask(umaskDefault)

	for _, test := range []struct {
		name     string
		args     []string
		perm     os.FileMode
		dirNames []string
	}{
		{
			name:     "Default Perm",
			args:     []string{"foo"},
			perm:     0755 | os.ModeDir,
			dirNames: []string{"foo"},
		},
		{
			name:     "Custom Perm in Octal Form",
			args:     []string{"-m=0777", "foo"},
			perm:     0777 | os.ModeDir,
			dirNames: []string{"foo"},
		},
		{
			name:     "Custom Perm not in Octal Form",
			args:     []string{"-m=777", "foo"},
			perm:     0777 | os.ModeDir,
			dirNames: []string{"foo"},
		},
		{
			name:     "Custom Perm with Sticky Bit",
			args:     []string{"-m=1777", "foo"},
			perm:     0777 | os.ModeDir | os.ModeSticky,
			dirNames: []string{"foo"},
		},
		{
			name:     "Custom Perm with SGID Bit",
			args:     []string{"-m=2777", "foo"},
			perm:     0777 | os.ModeDir | os.ModeSetgid,
			dirNames: []string{"foo"},
		},
		{
			name:     "Custom Perm with SUID Bit",
			args:     []string{"-m=4777", "foo"},
			perm:     0777 | os.ModeDir | os.ModeSetuid,
			dirNames: []string{"foo"},
		},
		{
			name:     "Custom Perm with Sticky Bit and SUID Bit",
			args:     []string{"-m=5777", "foo"},
			perm:     0777 | os.ModeDir | os.ModeSticky | os.ModeSetuid,
			dirNames: []string{"foo"},
		},
		{
			name:     "Custom Perm for 2 Directories",
			args:     []string{"-m=5777", "foo", "bar"},
			perm:     0777 | os.ModeDir | os.ModeSticky | os.ModeSetuid,
			dirNames: []string{"foo", "bar"},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			tmpDir, err := ioutil.TempDir("", "mkdir")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tmpDir)

			c := testutil.Command(t, test.args...)
			c.Dir = tmpDir
			if err := c.Run(); err != nil {
				t.Fatalf("mkdir exited with error: %v", err)
			}

			for _, dirName := range test.dirNames {
				f, err := findFile(tmpDir, dirName)
				if err != nil {
					t.Fatal(err)
				}
				if f == nil {
					t.Errorf("could not find file %q", dirName)
				}
				if f != nil && !reflect.DeepEqual(f.Mode(), test.perm) {
					t.Errorf("file %q has mode %v, want %v", dirName, f.Mode(), test.perm)
				}
			}
		})
	}
}

func TestMain(m *testing.M) {
	if testutil.CallMain() {
		main()
		os.Exit(0)
	}

	os.Exit(m.Run())
}
