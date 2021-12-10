// Copyright 2015-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/cpio"
	"github.com/u-root/u-root/pkg/golang"
	"github.com/u-root/u-root/pkg/testutil"
	"github.com/u-root/u-root/pkg/uio"
	itest "github.com/u-root/u-root/pkg/uroot/initramfs/test"
)

var twocmds = []string{
	"github.com/u-root/u-root/cmds/core/ls",
	"github.com/u-root/u-root/cmds/core/init",
}

func xTestDCE(t *testing.T) {
	delFiles := false
	f, _ := buildIt(
		t,
		[]string{
			"-build=bb", "-no-strip",
			"world",
			"-github.com/u-root/u-root/cmds/exp/builtin",
			"-github.com/u-root/u-root/cmds/exp/run",
			"github.com/u-root/u-root/pkg/uroot/test/foo",
		},
		nil,
		nil)
	defer func() {
		if delFiles {
			os.RemoveAll(f.Name())
		}
	}()
	st, _ := f.Stat()
	t.Logf("Built %s, size %d", f.Name(), st.Size())
	cmd := golang.Default().GoCmd("tool", "nm", f.Name())
	nmOutput, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to run nm: %s %s", err, nmOutput)
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
		t.Logf("%s", sym)
	}
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
	cmd := golang.Default().GoCmd("tool", "nm", tf.Name())
	nmOutput, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run nm: %s %s", err, nmOutput)
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

// Turn this off until we make the full switch to modules.
func testUrootCmdline(t *testing.T) {
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
			name: "uinitcmd",
			args: []string{"-build=bb", "-uinitcmd=echo foobar fuzz", "-defaultsh=", "github.com/u-root/u-root/cmds/core/init", "github.com/u-root/u-root/cmds/core/echo"},
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
			name: "fix usage of an absolute path",
			args: []string{"-nocmd", fmt.Sprintf("-files=%s:/bin", sampledir)},
			err:  nil,
			validators: []itest.ArchiveValidator{
				itest.HasFile{"/bin/foo"},
				itest.HasFile{"/bin/bar"},
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
			name: "make sure dead code gets eliminated",
			args: []string{
				// Build the world + test symbols, unstripped.
				// Change default shell to ash/rush until a solution for upstream elvish is found
				"-defaultsh=/bin/ash", "-build=bb", "-no-strip", "world", "github.com/u-root/u-root/pkg/uroot/test/foo",
				// These are known to disable DCE and need to be exluded.
				// The reason is https://github.com/golang/go/issues/36021 and is fixed in Go 1.15,
				// so these can be removed once we no longer support Go < 1.15.
				"-github.com/u-root/u-root/cmds/exp/builtin",
				"-github.com/u-root/u-root/cmds/exp/run",
				// TODO(MDr164): Find a solution for upstream elvish
				"-github.com/u-root/u-root/cmds/core/elvish",
			},
			err: nil,
			validators: []itest.ArchiveValidator{
				noDeadCode{Path: "bbin/bb"},
			},
		},
		{
			name: "hosted bb mode",
			args: append([]string{"-build=bb", "-base=/dev/null", "-defaultsh=", "-initcmd="}, twocmds...),
		},
		{
			name: "AMD64 bb build",
			env:  []string{"GOARCH=amd64"},
			args: []string{"-build=bb", "all"},
		},
		{
			name: "MIPS bb build",
			env:  []string{"GOARCH=mips"},
			args: []string{"-build=bb", "all"},
		},
		{
			name: "MIPSLE bb build",
			env:  []string{"GOARCH=mipsle"},
			args: []string{"-build=bb", "all"},
		},
		{
			name: "MIPS64 bb build",
			env:  []string{"GOARCH=mips64"},
			args: []string{"-build=bb", "all"},
		},
		{
			name: "MIPS64LE bb build",
			env:  []string{"GOARCH=mips64le"},
			args: []string{"-build=bb", "all"},
		},
		{
			name: "ARM7 bb build",
			env:  []string{"GOARCH=arm", "GOARM=7"},
			args: []string{"-build=bb", "all"},
		},
		{
			name: "ARM64 bb build",
			env:  []string{"GOARCH=arm64"},
			args: []string{"-build=bb", "all"},
		},
		{
			name: "386 (32 bit) bb build",
			env:  []string{"GOARCH=386"},
			args: []string{"-build=bb", "all"},
		},
		{
			name: "Power 64bit bb build",
			env:  []string{"GOARCH=ppc64le"},
			args: []string{"-build=bb", "all"},
		},
		{
			name: "RISCV 64bit bb build",
			env:  []string{"GOARCH=riscv64"},
			args: []string{"-build=bb", "all"},
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
	f, err := os.CreateTemp("", "u-root-")
	if err != nil {
		t.Fatal(err)
	}
	// Use the u-root command outside of the $GOPATH tree to make sure it
	// still works.
	arg := append([]string{"-o", f.Name()}, args...)
	c := testutil.Command(t, arg...)
	t.Logf("Commandline: %v", arg)
	c.Env = append(c.Env, env...)
	c.Dir = t.TempDir()
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
