// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/cpio"
)

// A version query is whatever follows the last "@", because a Go package path
// cannot contain one and a path may already carry a major version suffix.
func TestParsePatternSplitsOnTheLastAt(t *testing.T) {
	for _, tt := range []struct {
		in    string
		path  string
		query string
	}{
		{"github.com/u-root/u-root/cmds/core/...", "github.com/u-root/u-root/cmds/core/...", ""},
		{"github.com/u-root/u-root/cmds/core/ls@v0.16.0", "github.com/u-root/u-root/cmds/core/ls", "v0.16.0"},
		{"example.com/m@main", "example.com/m", "main"},
		// The major version suffix is part of the path, not a query, so
		// splitting on the first "@" would take the wrong half.
		{"gopkg.in/yaml.v2@v2.4.0", "gopkg.in/yaml.v2", "v2.4.0"},
		{"  example.com/m@latest  ", "example.com/m", "latest"},
	} {
		got := parsePattern(tt.in)
		if got.path != tt.path || got.query != tt.query {
			t.Errorf("parsePattern(%q) = {%q, %q}, want {%q, %q}", tt.in, got.path, got.query, tt.path, tt.query)
		}
	}
}

// fetchSpec puts the pattern back together for `go get`, and must not leave a
// bare "@" when no version was asked for.
func TestFetchSpecOmitsAnEmptyQuery(t *testing.T) {
	if got := parsePattern("example.com/m").fetchSpec(); got != "example.com/m" {
		t.Errorf("fetchSpec with no query = %q, want %q", got, "example.com/m")
	}
	if got := parsePattern("example.com/m@v1.2.3").fetchSpec(); got != "example.com/m@v1.2.3" {
		t.Errorf("fetchSpec with a query = %q, want %q", got, "example.com/m@v1.2.3")
	}
}

// A symlink in the archive is relative, and relative symlinks resolve from
// their parent directory rather than from the root. Getting this wrong points
// bin/sh at /bin/bbin/gosh, which exists nowhere.
func TestSymlinkTargetsAreRelativeToTheirParent(t *testing.T) {
	for _, tt := range []struct {
		source, cmd, want string
	}{
		{"init", "init", "bbin/init"},
		{"bin/sh", "gosh", "../bbin/gosh"},
		{"bin/defaultsh", "gosh", "../bbin/gosh"},
	} {
		rec, err := symlinkTo(tt.source, tt.cmd)
		if err != nil {
			t.Fatalf("symlinkTo(%q, %q): %v", tt.source, tt.cmd, err)
		}
		got, err := readTarget(rec)
		if err != nil {
			t.Fatal(err)
		}
		if got != tt.want {
			t.Errorf("symlinkTo(%q, %q) targets %q, want %q", tt.source, tt.cmd, got, tt.want)
		}
	}
}

// readTarget returns the target a symlink record points at.
func readTarget(rec cpio.Record) (string, error) {
	b := make([]byte, rec.Info.FileSize)
	if _, err := rec.ReaderAt.ReadAt(b, 0); err != nil {
		return "", err
	}
	return string(b), nil
}

// Every command gets a symlink in bbin pointing at the one binary, which is
// what makes a busybox dispatch on argv[0].
func TestArchiveLinksEveryCommandAtTheBusybox(t *testing.T) {
	bbPath := filepath.Join(t.TempDir(), "bb")
	if err := os.WriteFile(bbPath, []byte("not really a binary"), 0o755); err != nil {
		t.Fatal(err)
	}

	recs, err := archive(bbPath, []string{
		"github.com/u-root/u-root/cmds/core/init",
		"github.com/u-root/u-root/cmds/core/gosh",
		"github.com/u-root/u-root/cmds/core/ls",
	})
	if err != nil {
		t.Fatalf("archive: %v", err)
	}

	names := map[string]cpio.Record{}
	for _, r := range recs {
		names[r.Name] = r
	}
	for _, want := range []string{
		"bbin/bb", "bbin/init", "bbin/gosh", "bbin/ls",
		"init", "bin/sh", "bin/defaultsh",
		"dev/console", "dev/null", "etc/resolv.conf",
	} {
		if _, ok := names[want]; !ok {
			t.Errorf("archive is missing %q", want)
		}
	}
	for _, cmd := range []string{"init", "gosh", "ls"} {
		got, err := readTarget(names["bbin/"+cmd])
		if err != nil {
			t.Fatal(err)
		}
		if got != "bb" {
			t.Errorf("bbin/%s targets %q, want %q", cmd, got, "bb")
		}
	}
}

// Linking /init to a command that was not built produces an initramfs that
// cannot boot, and the kernel's complaint names neither the flag nor the
// command. Refuse it while the names are still in hand.
func TestArchiveRefusesASymlinkToAMissingCommand(t *testing.T) {
	bbPath := filepath.Join(t.TempDir(), "bb")
	if err := os.WriteFile(bbPath, []byte("not really a binary"), 0o755); err != nil {
		t.Fatal(err)
	}

	// gosh is present, so only the init symlink can fail here.
	_, err := archive(bbPath, []string{"github.com/u-root/u-root/cmds/core/gosh"})
	if err == nil {
		t.Fatal("archive with no init command succeeded, want an error")
	}
	if !strings.Contains(err.Error(), "init") {
		t.Errorf("error = %q, want it to name the missing command", err)
	}
}

// The skeleton carries the device nodes init expects to find already there,
// with the major and minor numbers Linux assigns them.
func TestSkeletonCarriesTheExpectedDeviceNodes(t *testing.T) {
	want := map[string][2]uint64{
		"dev/console": {5, 1},
		"dev/tty":     {5, 0},
		"dev/null":    {1, 3},
		"dev/port":    {1, 4},
		"dev/urandom": {1, 9},
	}
	got := map[string][2]uint64{}
	for _, r := range skeleton() {
		if strings.HasPrefix(r.Name, "dev/") {
			got[r.Name] = [2]uint64{r.Info.Rmajor, r.Info.Rminor}
		}
	}
	for name, mm := range want {
		if got[name] != mm {
			t.Errorf("%s is %v, want %v", name, got[name], mm)
		}
	}
	if len(got) != len(want) {
		t.Errorf("skeleton has %d device nodes, want %d", len(got), len(want))
	}
}

// goCmd must neutralise UROOT_SOURCE. findpkg reads it and turns a u-root
// package path into a file system lookup, which would quietly build a local
// checkout rather than what the module cache holds -- the one thing this tool
// exists to avoid.
func TestGoCmdClearsTheEnvironmentThatWouldUseLocalSource(t *testing.T) {
	t.Setenv("UROOT_SOURCE", "/somewhere/with/a/checkout")
	t.Setenv("GOPRIVATE", "example.com/*")
	t.Setenv("GONOPROXY", "example.com/*")

	last := map[string]string{}
	for _, kv := range goCmd(t.TempDir(), "list").Env {
		if k, v, ok := strings.Cut(kv, "="); ok {
			// exec uses the last value for a repeated key.
			last[k] = v
		}
	}
	for _, k := range []string{"UROOT_SOURCE", "GOPRIVATE", "GONOPROXY"} {
		if last[k] != "" {
			t.Errorf("%s = %q, want it cleared", k, last[k])
		}
	}
	if last["GOVCS"] != "*:off" {
		t.Errorf("GOVCS = %q, want %q", last["GOVCS"], "*:off")
	}
	if last["GOWORK"] != "off" {
		t.Errorf("GOWORK = %q, want %q", last["GOWORK"], "off")
	}
}
