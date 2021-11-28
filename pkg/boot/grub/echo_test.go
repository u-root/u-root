// Copyright 2017-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grub

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/curl"
	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/mount/block"
)

var update = flag.Bool("run-bash", false, "run bash and update golden file")

// TestMain is used to warp a go utility in the same test binary
// when the environment variable BE_ECHO is set to 1, the binary will echo its
// parameters using the %#v format string, so the parameters are escaped and can
// be recovered.
func TestMain(m *testing.M) {
	if os.Getenv("BE_ECHO") == "1" {
		fmt.Printf("echo:%#v\n", os.Args[1:])
		return
	} // call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}

// TestHelperEcho tests the echo wrapper in TestMain
func TestHelperEcho(t *testing.T) {
	cmd := exec.Command(os.Args[0], "echothis")
	cmd.Env = append(os.Environ(), "BE_ECHO=1")
	out, err := cmd.Output()
	t.Logf("%q\n", out)
	if err != nil {
		t.Fatalf("process ran with err %v", err)
	}
	want := "echo:[]string{\"echothis\"}\n"
	if string(out) != want {
		t.Fatalf("wrong process output got `%s` want `%s`", out, want)
	}
}

// TestBashWrapper tests that the "./testdata/bash_wrapper.sh" works as expected
// bash_wrapper.sh is a script that replace the internal command echo with its
// first argument and source its second argument.
// The goal is to be able to run grub's tests scripts, see TestGrubTests
func TestBashWrapper(t *testing.T) {
	if !*update {
		t.Skip("use -run-bash flag to run this")
	}
	cmd := exec.Command("./testdata/bash_wrapper.sh", os.Args[0], "./testdata/test_bash_wrapper.sh")
	cmd.Env = append(os.Environ(), "BE_ECHO=1")
	out, err := cmd.Output()
	t.Logf("%q\n", out)
	if err != nil {
		t.Fatalf("process ran with err %v", err)
	}
	want := "echo:[]string{\"param1\", \"param2\"}\n"
	if string(out) != want {
		t.Fatalf("wrong process output got `%s` want `%s`", out, want)
	}
}

// TestGrubTests run tests imported from grub source to check our parser
// grub has tests in for of scripts that are run both by grub and bash, they
// mostly use echo and the test then compare the output of both runs.
// In our case we don't want to compare the output of echo, but get the token
// passed to echo. So we replace the echo command in bash with the wrapper (see
// above). We can then compare the bash output to our parser output.
// Also to avoid keeping the dependency on bash, the output are saved in the
// golden files. One must run the test with '-run-bash' to update the golden
// files in case new tests are added or the echo format is changed.
func TestGrubTests(t *testing.T) {
	files, err := filepath.Glob("testdata/*.in")
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		name := strings.TrimSuffix(filepath.Base(file), ".in")
		t.Run(name, func(t *testing.T) {
			golden := strings.TrimSuffix(file, ".in") + ".out"
			var out []byte
			if *update {
				cmd := exec.Command("./testdata/bash_wrapper.sh", os.Args[0], file)
				cmd.Env = append(os.Environ(), "BE_ECHO=1")
				out, err = cmd.Output()
				// t.Logf("%s\n", out)
				if err != nil {
					t.Fatalf("process ran with err %v", err)
				}
			} else {
				out, err = os.ReadFile(golden)
				if err != nil {
					t.Fatalf("error loading file `%s`, %v", golden, err)
				}
			}
			// parse with our parser and compare
			var b bytes.Buffer
			wd := &url.URL{
				Scheme: "file",
				Path:   "./testdata",
			}
			mountPool := &mount.Pool{}
			c := newParser(wd, block.BlockDevices{}, mountPool, curl.DefaultSchemes)
			c.W = &b

			script, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("error loading file `%s`, %v", file, err)
			}
			err = c.append(context.Background(), string(script))
			if err != nil {
				t.Fatalf("error parsing file `%s`, %v", file, err)
			}

			if b.String() != string(out) {
				t.Fatalf("wrong script parsing output got `%s` want `%s`", b.String(), string(out))
			}
			// update/create golden file on success
			if *update {
				err := os.WriteFile(golden, out, 0o644)
				if err != nil {
					t.Fatalf("error writing file `%s`, %v", file, err)
				}
			}
		})

	}
}
