// Copyright 2015-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/cpio"
	"github.com/u-root/u-root/pkg/testutil"
	itest "github.com/u-root/u-root/pkg/uroot/initramfs/test"
)

var twocmds = []string{
	"github.com/u-root/u-root/cmds/core/ls",
	"github.com/u-root/u-root/cmds/core/init",
}

var srcmds = []string{
	"github.com/u-root/u-root/cmds/core/ls",
	"github.com/u-root/u-root/cmds/core/init",
	"github.com/u-root/u-root/cmds/core/installcommand",
}

type buildSourceValidator struct {
	gopath string
	goroot string
	env    []string
}

func (b buildSourceValidator) Validate(a *cpio.Archive) error {
	dir, err := ioutil.TempDir("", "u-root-source-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	if err := os.Mkdir(filepath.Join(dir, "tmp"), 0755); err != nil {
		return err
	}

	// Unpack into dir.
	err = cpio.ForEachRecord(a.Reader(), func(r cpio.Record) error {
		return cpio.CreateFileInRoot(r, dir, false)
	})
	if err != nil {
		return err
	}

	goroot := filepath.Join(dir, b.goroot)
	gopath := filepath.Join(dir, b.gopath)
	// go build ./src/...
	c := exec.Command(filepath.Join(goroot, "bin/go"), "build", filepath.Join(gopath, "src/..."))
	c.Env = append(b.env,
		fmt.Sprintf("GOPATH=%s", gopath),
		fmt.Sprintf("GOCACHE=%s", filepath.Join(dir, "tmp")),
		fmt.Sprintf("GOROOT=%s", goroot))
	out, err := c.CombinedOutput()
	if err != nil {
		return fmt.Errorf("could not build go source %v; output\n%s", err, out)
	}
	return nil
}

func TestUrootCmdline(t *testing.T) {
	samplef, err := ioutil.TempFile("", "u-root-sample-")
	if err != nil {
		t.Fatal(err)
	}
	samplef.Close()
	defer os.RemoveAll(samplef.Name())

	for _, tt := range []struct {
		name       string
		env        []string
		args       []string
		err        error
		validators []itest.ArchiveValidator
	}{
		{
			name: "include one extra file",
			args: []string{"-nocmd", "-files=/bin/bash"},
			err:  nil,
			validators: []itest.ArchiveValidator{
				itest.HasFile{"bin/bash"},
			},
		},
		{
			name: "fix usage of an absolute path",
			args: []string{"-nocmd", "-files=/bin:/bin"},
			err:  nil,
			validators: []itest.ArchiveValidator{
				itest.HasFile{"bin/bash"},
			},
		},
		{
			name: "include multiple extra files",
			args: []string{"-nocmd", "-files=/bin/bash", "-files=/bin/ls", fmt.Sprintf("-files=%s", samplef.Name())},
			validators: []itest.ArchiveValidator{
				itest.HasFile{"bin/bash"},
				itest.HasFile{"bin/ls"},
				itest.HasFile{samplef.Name()},
			},
		},
		{
			name: "include one extra file with rename",
			args: []string{"-nocmd", "-files=/bin/bash:bin/bush"},
			validators: []itest.ArchiveValidator{
				itest.HasFile{"bin/bush"},
			},
		},
		{
			name: "hosted source mode",
			args: append([]string{"-build=source", "-base=/dev/null", "-defaultsh=", "-initcmd="}, srcmds...),
		},
		{
			name: "hosted bb mode",
			args: append([]string{"-build=bb", "-base=/dev/null", "-defaultsh=", "-initcmd="}, twocmds...),
		},
		{
			name: "AMD64 bb build",
			env:  []string{"GOARCH=amd64"},
			args: []string{"-build=bb"},
		},
		{
			name: "AMD64 source build",
			env:  []string{"GOARCH=amd64"},
			args: []string{"-build=source"},
			validators: []itest.ArchiveValidator{
				buildSourceValidator{
					goroot: "/go",
					gopath: ".",
					env:    []string{"GOARCH=amd64"},
				},
			},
		},
		{
			name: "ARM7 bb build",
			env:  []string{"GOARCH=arm", "GOARM=7"},
			args: []string{"-build=bb"},
		},
		{
			name: "ARM64 bb build",
			env:  []string{"GOARCH=arm64"},
			args: []string{"-build=bb"},
		},
		{
			name: "Power 64bit bb build",
			env:  []string{"GOARCH=ppc64le"},
			args: []string{"-build=bb"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			delFiles := true
			f, sum1 := buildIt(t, tt.args, tt.env, tt.err)
			defer func() {
				if delFiles {
					os.RemoveAll(f.Name())
				}
			}()

			a, err := itest.ReadArchive(f.Name())
			if err != nil {
				t.Fatal(err)
			}

			for _, v := range tt.validators {
				if err := v.Validate(a); err != nil {
					t.Errorf("validator failed: %v / archive:\n%s", err, a)
				}
			}

			f2, sum2 := buildIt(t, tt.args, tt.env, tt.err)
			defer func() {
				if delFiles {
					os.RemoveAll(f2.Name())
				}
			}()
			if !bytes.Equal(sum1, sum2) {
				delFiles = false
				t.Errorf("not reproducible, hashes don't match")
				t.Errorf("env: %v args: %v", tt.env, tt.args)
				t.Errorf("file1: %v file2: %v", f.Name(), f2.Name())
			}
		})
	}
}

func buildIt(t *testing.T, args, env []string, want error) (*os.File, []byte) {
	f, err := ioutil.TempFile("", "u-root-")
	if err != nil {
		t.Fatal(err)
	}

	arg := append([]string{"-o", f.Name()}, args...)
	c := testutil.Command(t, arg...)
	t.Logf("Commandline: %v", arg)
	c.Env = append(c.Env, env...)
	if out, err := c.CombinedOutput(); err != want {
		t.Fatalf("Error: %v\nOutput:\n%s", err, out)
	} else if err != nil {
		h1 := sha256.New()
		if _, err := io.Copy(h1, f); err != nil {
			t.Fatal()
		}
		return f, h1.Sum(nil)
	}
	return f, nil
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}
