// Copyright 2015-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/u-root/gobusybox/src/pkg/golang"
)

var makebb = flag.String("makebb", "", "makebb binary path")

func skipIfUnsupported(t *testing.T, goVersion string, unsupportedGoVersions []string) {
	for _, unsupportedGoVersion := range unsupportedGoVersions {
		// Contains so that go1.13 also matches e.g. go1.13.1
		if strings.Contains(goVersion, unsupportedGoVersion) {
			t.Skipf("Version %s is unsupported for this test (unsupported versions: %v)", goVersion, unsupportedGoVersion)
		}
	}
}

func TestMakeBB(t *testing.T) {
	if *makebb == "" {
		t.Fatalf("Path to makebb is not set")
	}
	mkbb, _ := filepath.Abs(*makebb)
	wd, _ := os.Getwd()

	goVersion, err := golang.Default().Version()
	if err != nil {
		t.Fatalf("Could not determine Go version: %v", err)
	}

	for _, tt := range []struct {
		testname string
		// file paths to commands to compile
		cmds []string
		// Working directory
		wd string
		// extra args to makebb
		extraArgs []string
		// command name -> expected output
		want map[string]string
		// Go versions for which this test should be skipped.
		unsupportedGoVersions []string
	}{
		{
			testname:              "goembed",
			cmds:                  []string{"."},
			wd:                    filepath.Join(wd, "goembed"),
			want:                  map[string]string{"goembed": "hello\n"},
			unsupportedGoVersions: []string{"go1.15"},
		},
		{
			testname: "12-fancy-cmd",
			cmds:     []string{"."},
			wd:       filepath.Join(wd, "12-fancy-cmd"),
			want:     map[string]string{"12-fancy-cmd": "12-fancy-cmd\n"},
		},
		{
			testname:  "injectldvar",
			cmds:      []string{"."},
			wd:        filepath.Join(wd, "injectldvar"),
			extraArgs: []string{"-go-extra-args=-ldflags", "-go-extra-args=-X 'github.com/u-root/gobusybox/test/injectldvar.Something=Hello World'"},
			want:      map[string]string{"injectldvar": "Hello World\n"},
		},
		{
			testname: "implicitimport",
			cmds:     []string{"./cmd/loghello"},
			wd:       filepath.Join(wd, "implicitimport"),
			want:     map[string]string{"loghello": "Log Hello\n"},
		},
		/*{
			testname: "nested-modules",
			cmds:     []string{"./cmd/dmesg", "./cmd/strace", "./nestedmod/cmd/p9ufs"},
			wd:       filepath.Join(wd, "nested"),
		},*/
		{
			testname: "cross-module-deps",
			wd:       filepath.Join(wd, "normaldeps/mod1"),
			cmds:     []string{"./cmd/helloworld", "./cmd/getppid"},
			want: map[string]string{
				"helloworld": "test/normaldeps/mod2/hello: test/normaldeps/mod2/v2/hello\n",
				"getppid":    fmt.Sprintf("%d\n", os.Getpid()),
			},
		},
		{
			testname: "import-name-conflict",
			wd:       filepath.Join(wd, "nameconflict"),
			cmds:     []string{"./cmd/nameconflict"},
		},
		{
			testname: "diamond-module-dependency",
			wd:       filepath.Join(wd, "diamonddep/mod1"),
			cmds:     []string{"./cmd/hellowithdep", "./cmd/helloworld"},
			want: map[string]string{
				"hellowithdep": "test/diamonddep/mod1/hello: test/diamonddep/mod1/hello\n" +
					"test/diamonddep/mod2/hello: test/diamonddep/mod2/hello\n" +
					"test/diamonddep/mod2/exthello: test/diamonddep/mod2/exthello: test/diamonddep/mod1/hello and test/diamonddep/mod3/hello\n",
				"helloworld": "hello world\n",
			},
		},
	} {
		t.Run(tt.testname, func(t *testing.T) {
			skipIfUnsupported(t, goVersion, tt.unsupportedGoVersions)

			dir, err := ioutil.TempDir("", tt.testname)
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(dir)

			var goEnv string
			t.Run(goEnv, func(t *testing.T) {
				binary := filepath.Join(dir, "bb")

				// Build the bb.
				t.Logf("Run: %s %s -o %s %v %s", goEnv, mkbb, binary, strings.Join(tt.extraArgs, " "), strings.Join(tt.cmds, " "))
				args := append([]string{"-o", binary}, tt.extraArgs...)
				cmd := exec.Command(mkbb, append(args, tt.cmds...)...)
				cmd.Dir = tt.wd
				cmd.Env = append(os.Environ(), goEnv)
				out, err := cmd.CombinedOutput()
				if err != nil {
					t.Logf("makebb: %s", string(out))
					t.Fatalf("cmd: %v", err)
				}

				// There are some builds for which we
				// don't check the output since it's
				// unpredictable. We at least want to
				// check the binary exists.
				if _, err := os.Stat(binary); err != nil {
					t.Fatalf("Busybox binary does not exist: %v", err)
				}

				// Make sure that the bb contains all
				// the commands it's supposed to by
				// invoking them and checking their
				// output.
				for cmdName, want := range tt.want {
					t.Logf("Run: %s %s", binary, cmdName)
					out, err = exec.Command(binary, cmdName).CombinedOutput()
					if err != nil {
						t.Fatalf("cmd: %v", err)
					}
					if got := string(out); got != want {
						t.Errorf("Output of %s = %v, want %v", cmdName, got, want)
					}
				}
			})

			// Make sure that bb is reproducible.
			binaryOn, err := ioutil.ReadFile(filepath.Join(dir, "bb-on"))
			if err != nil {
				t.Errorf("bb binary for GO111MODULE=on does not exist: %v", err)
			}
			binaryAuto, err := ioutil.ReadFile(filepath.Join(dir, "bb-auto"))
			if err != nil {
				t.Errorf("bb binary for GO111MODULE=auto does not exist: %v", err)
			}
			if !bytes.Equal(binaryOn, binaryAuto) {
				t.Errorf("bb not reproducible")
			}
		})
	}
}

func TestBBSymlink(t *testing.T) {
	if *makebb == "" {
		t.Fatalf("Path to makebb is not set")
	}
	mkbb, _ := filepath.Abs(*makebb)

	dir := t.TempDir()

	cmdPath := "./cmd/loghello"
	want := "Log Hello\n"

	goEnv := "GO111MODULE=on"
	binary := filepath.Join(dir, "bb")

	// Build the bb.
	t.Logf("Run: (cd ./implicitimport && %s %s -o %s %s)", goEnv, *makebb, binary, cmdPath)

	cmd := exec.Command(mkbb, "-o", binary, cmdPath)
	cmd.Dir = "./implicitimport"
	cmd.Env = append(os.Environ(), goEnv)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("makebb: %s", string(out))
		t.Fatalf("cmd: %v", err)
	}

	t.Run("symlink", func(t *testing.T) {
		link := filepath.Join(dir, "loghello")
		// Test that symlinking works.
		if err := os.Symlink(binary, link); err != nil {
			t.Fatal(err)
		}

		t.Logf("Run: %s", link)
		out, err = exec.Command(link).CombinedOutput()
		if err != nil {
			t.Fatalf("cmd: %v", err)
		}
		if got := string(out); got != want {
			t.Errorf("Output of %s = %v, want %v", link, got, want)
		}
	})

	t.Run("argv1", func(t *testing.T) {
		cmdName := "loghello"
		t.Logf("Run: %s %s", binary, cmdName)
		out, err = exec.Command(binary, cmdName).CombinedOutput()
		if err != nil {
			t.Fatalf("cmd: %v", err)
		}
		if got := string(out); got != want {
			t.Errorf("Output of %s %s = %v, want %v", binary, cmdName, got, want)
		}
	})
}
