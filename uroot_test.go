// Copyright 2015-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	gbbgolang "github.com/u-root/gobusybox/src/pkg/golang"
	"github.com/u-root/u-root/pkg/cpio"
	"github.com/u-root/u-root/pkg/testutil"
	itest "github.com/u-root/u-root/pkg/uroot/initramfs/test"
	"github.com/u-root/uio/uio"
)

var twocmds = []string{
	"github.com/u-root/u-root/cmds/core/ls",
	"github.com/u-root/u-root/cmds/core/init",
}

type noDeadCode struct {
	Path string
}

func (v noDeadCode) Validate(a *cpio.Archive) error {
	// 1. Extract BB binary into a temporary file.
	delFiles := true
	bbRecord, ok := a.Get(v.Path)
	if !ok {
		return fmt.Errorf("archive does not contain %s, but should", v.Path)
	}
	tf, err := os.CreateTemp("", "u-root-temp-bb-")
	if err != nil {
		return err
	}
	bbData, _ := uio.ReadAll(bbRecord)
	tf.Write(bbData)
	tf.Close()
	defer func() {
		if delFiles {
			os.RemoveAll(tf.Name())
		}
	}()
	// 2. Run "go nm" on it and build symbol table.
	cmd := gbbgolang.Default().GoCmd("tool", "nm", tf.Name())
	nmOutput, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run nm: %w %s", err, nmOutput)
	}
	symScanner := bufio.NewScanner(bytes.NewBuffer(nmOutput))
	syms := map[string]bool{}
	for symScanner.Scan() {
		line := symScanner.Text()
		parts := strings.Split(line, " ")
		if len(parts) == 0 {
			continue
		}
		sym := parts[len(parts)-1]
		syms[sym] = true
	}
	// 3. Check for presence and absence of particular symbols.
	if !syms["github.com/u-root/u-root/pkg/uroot/test/bar.Bar.UsedInterfaceMethod"] {
		// Sanity check of the test itself: this method must be in the binary.
		return fmt.Errorf("expected symbol not found, something is wrong with the build")
	}
	if syms["github.com/u-root/u-root/pkg/uroot/test/bar.Bar.UnusedNonInterfaceMethod"] {
		// Sanity check of the test itself: this method must be in the binary.
		delFiles = false
		return fmt.Errorf(
			"Unused non-interface method has not been eliminated, dead code elimination is not working properly.\n"+
				"The most likely reason is use of reflect.Value.Method or .MethodByName somewhere "+
				"(could be a command or vendor dependency, apologies for not being more precise here).\n"+
				"See https://golang.org/src/cmd/link/internal/ld/deadcode.go for explanation.\n"+
				"%s contains the resulting binary.\n", tf.Name())
	}
	return nil
}

func TestUrootCmdline(t *testing.T) {
	samplef, err := os.CreateTemp("", "u-root-test-")
	if err != nil {
		t.Fatal(err)
	}
	samplef.Close()
	defer os.RemoveAll(samplef.Name())
	sampledir := t.TempDir()
	if err = os.WriteFile(filepath.Join(sampledir, "foo"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(sampledir, "bar"), nil, 0o644); err != nil {
		t.Fatal(err)
	}

	type testCase struct {
		name       string
		env        []string
		args       []string
		err        error
		validators []itest.ArchiveValidator
	}

	noCmdTests := []testCase{
		{
			name: "include one extra file",
			args: []string{"-nocmd", "-defaultsh=", "-initcmd=", "-files=/bin/bash"},
			env:  []string{"GO111MODULE=off"},
			err:  nil,
			validators: []itest.ArchiveValidator{
				itest.HasFile{"bin/bash"},
			},
		},
		{
			name: "fix usage of an absolute path",
			args: []string{"-nocmd", "-defaultsh=", "-initcmd=", fmt.Sprintf("-files=%s:/bin", sampledir)},
			env:  []string{"GO111MODULE=off"},
			err:  nil,
			validators: []itest.ArchiveValidator{
				itest.HasFile{"/bin/foo"},
				itest.HasFile{"/bin/bar"},
			},
		},
		{
			name: "include multiple extra files",
			args: []string{"-nocmd", "-defaultsh=", "-initcmd=", "-files=/bin/bash", "-files=/bin/ls", fmt.Sprintf("-files=%s", samplef.Name())},
			env:  []string{"GO111MODULE=off"},
			validators: []itest.ArchiveValidator{
				itest.HasFile{"bin/bash"},
				itest.HasFile{"bin/ls"},
				itest.HasFile{samplef.Name()},
			},
		},
		{
			name: "include one extra file with rename",
			args: []string{"-nocmd", "-defaultsh=", "-initcmd=", "-files=/bin/bash:bin/bush"},
			env:  []string{"GO111MODULE=off"},
			validators: []itest.ArchiveValidator{
				itest.HasFile{"bin/bush"},
			},
		},
		{
			name: "supplied file can be uinit",
			args: []string{"-nocmd", "-defaultsh=", "-initcmd=", "-files=/bin/bash:bin/bash", "-uinitcmd=/bin/bash"},
			env:  []string{"GO111MODULE=off"},
			validators: []itest.ArchiveValidator{
				itest.HasFile{"bin/bash"},
				itest.HasRecord{cpio.Symlink("bin/uinit", "bash")},
			},
		},
	}

	bareTests := []testCase{
		{
			name: "uinitcmd",
			args: []string{"-uinitcmd=echo foobar fuzz", "-defaultsh=", "github.com/u-root/u-root/cmds/core/init", "github.com/u-root/u-root/cmds/core/echo"},
			err:  nil,
			validators: []itest.ArchiveValidator{
				itest.HasRecord{cpio.Symlink("bin/uinit", "../bbin/echo")},
				itest.HasContent{
					Path:    "etc/uinit.flags",
					Content: "\"foobar\"\n\"fuzz\"",
				},
			},
		},
		{
			name: "dead_code_elimination",
			args: []string{
				// Build the world + test symbols, unstripped.
				"-no-strip", "world", "github.com/u-root/u-root/pkg/uroot/test/foo",
			},
			err: nil,
			validators: []itest.ArchiveValidator{
				noDeadCode{Path: "bbin/bb"},
			},
		},
		{
			name: "hosted mode",
			args: append([]string{"-base=/dev/null", "-defaultsh=", "-initcmd="}, twocmds...),
		},
		{
			name: "AMD64 build",
			env:  []string{"GOARCH=amd64"},
			args: []string{"all"},
		},
		{
			name: "MIPS build",
			env:  []string{"GOARCH=mips"},
			args: []string{"all"},
		},
		{
			name: "MIPSLE build",
			env:  []string{"GOARCH=mipsle"},
			args: []string{"all"},
		},
		{
			name: "MIPS64 build",
			env:  []string{"GOARCH=mips64"},
			args: []string{"all"},
		},
		{
			name: "MIPS64LE build",
			env:  []string{"GOARCH=mips64le"},
			args: []string{"all"},
		},
		{
			name: "ARM7 build",
			env:  []string{"GOARCH=arm", "GOARM=7"},
			args: []string{"all"},
		},
		{
			name: "ARM64 build",
			env:  []string{"GOARCH=arm64"},
			args: []string{"all"},
		},
		{
			name: "386 (32 bit) build",
			env:  []string{"GOARCH=386"},
			args: []string{"all"},
		},
		{
			name: "Power 64bit build",
			env:  []string{"GOARCH=ppc64le"},
			args: []string{"all"},
		},
		{
			name: "RISCV 64bit build",
			env:  []string{"GOARCH=riscv64"},
			args: []string{"all"},
		},
	}
	var bbTests []testCase
	for _, test := range bareTests {
		gbbTest := test
		gbbTest.name = gbbTest.name + " gbb-gomodule"
		gbbTest.args = append([]string{"-build=gbb"}, gbbTest.args...)
		gbbTest.env = append(gbbTest.env, "GO111MODULE=on")

		bbTests = append(bbTests, gbbTest)
	}

	for _, tt := range append(noCmdTests, bbTests...) {
		t.Run(tt.name, func(t *testing.T) {
			delFiles := true
			var (
				f1, f2     *os.File
				sum1, sum2 []byte
				errs       [2]error
				wg         = &sync.WaitGroup{}
				remove     []string
			)

			wg.Add(2)
			go func() {
				defer wg.Done()
				f1, sum1, err = buildIt(t, tt.args, tt.env, tt.err)
				if err != nil {
					errs[0] = err
					return
				}

				a, err := itest.ReadArchive(f1.Name())
				if err != nil {
					errs[0] = err
					return
				}

				remove = append(remove, f1.Name())
				for _, v := range tt.validators {
					if err := v.Validate(a); err != nil {
						t.Errorf("validator failed: %v / archive:\n%s", err, a)
					}
				}
			}()

			go func() {
				defer wg.Done()
				var err error
				f2, sum2, err = buildIt(t, tt.args, tt.env, tt.err)
				if err != nil {
					errs[1] = err
					return
				}
				remove = append(remove, f2.Name())
			}()

			wg.Wait()
			defer func() {
				if delFiles {
					for _, n := range remove {
						os.RemoveAll(n)
					}
				}
			}()
			if errs[0] != nil {
				t.Error(errs[0])
				return
			}
			if errs[1] != nil {
				t.Error(errs[1])
				return
			}
			if !bytes.Equal(sum1, sum2) {
				delFiles = false
				t.Errorf("not reproducible, hashes don't match")
				t.Errorf("env: %v args: %v", tt.env, tt.args)
				t.Errorf("file1: %v file2: %v", f1.Name(), f2.Name())
			}
		})
	}
}

func buildIt(t *testing.T, args, env []string, want error) (*os.File, []byte, error) {
	f, err := os.CreateTemp("", "u-root-")
	if err != nil {
		return nil, nil, err
	}
	// Use the u-root command outside of the $GOPATH tree to make sure it
	// still works.
	arg := append([]string{"-o", f.Name()}, args...)
	c := testutil.Command(t, arg...)
	t.Logf("Commandline: %v u-root %v", strings.Join(env, " "), strings.Join(arg, " "))
	c.Env = append(c.Env, env...)
	if out, err := c.CombinedOutput(); err != want {
		return nil, nil, fmt.Errorf("Error: %w\nOutput:\n%s", err, out)
	} else if err != nil {
		h1 := sha256.New()
		if _, err := io.Copy(h1, f); err != nil {
			return nil, nil, err
		}
		return f, h1.Sum(nil), nil
	}
	return f, nil, nil
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}

func TestCheckArgs(t *testing.T) {
	for _, tt := range []struct {
		name string
		args []string
		err  error
	}{
		{"-files is only arg", []string{"-files"}, errEmptyFilesArg},
		{"-files followed by -files", []string{"-files", "-files"}, errEmptyFilesArg},
		{"-files followed by any other switch", []string{"-files", "-abc"}, errEmptyFilesArg},
		{"no args", []string{}, nil},
		{"u-root alone", []string{"u-root"}, nil},
		{"u-root with -files and other args", []string{"u-root", "-files", "/bin/bash", "core"}, nil},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkArgs(tt.args...); !errors.Is(err, tt.err) {
				t.Errorf("%q: got %v, want %v", tt.args, err, tt.err)
			}
		})
	}
}
