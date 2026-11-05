// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

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

	recs, err := testOpts().archive(bbPath, []string{
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
	_, err := testOpts().archive(bbPath, []string{"github.com/u-root/u-root/cmds/core/gosh"})
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
	for _, r := range testOpts().skeleton() {
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

// goCmd must keep the fetch on the module proxy. GOPRIVATE and GONOPROXY both
// name paths that bypass it and both are read from the caller's environment.
func TestGoCmdKeepsTheFetchOnTheProxy(t *testing.T) {
	t.Setenv("GOPRIVATE", "example.com/*")
	t.Setenv("GONOPROXY", "example.com/*")

	last := map[string]string{}
	for _, kv := range goCmd(t.TempDir(), "list").Env {
		if k, v, ok := strings.Cut(kv, "="); ok {
			// exec uses the last value for a repeated key.
			last[k] = v
		}
	}
	for _, k := range []string{"GOPRIVATE", "GONOPROXY"} {
		if last[k] != "" {
			t.Errorf("%s = %q, want it cleared", k, last[k])
		}
	}
	if last["GOVCS"] != "*:off" {
		t.Errorf("GOVCS = %q, want %q", last["GOVCS"], "*:off")
	}
}

// The settings every go command has to agree on live in the process
// environment, not on the commands this file builds.
//
// gobusybox runs its own lookups and compiles with append(os.Environ(), ...),
// so anything set only on the commands here would not reach them. UROOT_SOURCE
// would then make findpkg build a local checkout; a stray GOWORK would pick a
// workspace over the fetched module; and CGO_ENABLED would let discovery list
// a cgo-only command the build then skips, leaving a symlink to nothing.
func TestFixEnvironmentCoversEveryGoCommand(t *testing.T) {
	t.Setenv("UROOT_SOURCE", "/somewhere/with/a/checkout")
	t.Setenv("GOWORK", "/somewhere/go.work")
	t.Setenv("CGO_ENABLED", "1")

	if err := fixEnvironment(); err != nil {
		t.Fatalf("fixEnvironment: %v", err)
	}
	for k, want := range map[string]string{
		"UROOT_SOURCE": "",
		"GOWORK":       "off",
		"CGO_ENABLED":  "0",
	} {
		if got := os.Getenv(k); got != want {
			t.Errorf("%s = %q, want %q", k, got, want)
		}
	}
}

// The embedded zoneinfo must be something Go can actually read.
//
// An earlier hand-written copy of this blob did not parse, and nothing would
// have reported it: a guest whose /etc/localtime is malformed falls back to
// UTC in silence, so the archive looks right and the clock is merely wrong.
func TestTheEmbeddedZoneinfoParses(t *testing.T) {
	loc, err := time.LoadLocationFromTZData("GMT0", []byte(gmt0))
	if err != nil {
		t.Fatalf("LoadLocationFromTZData on the embedded blob: %v", err)
	}
	at := time.Date(2026, 6, 1, 12, 0, 0, 0, loc)
	if name, offset := at.Zone(); name != "GMT" || offset != 0 {
		t.Errorf("zone = %q offset = %d, want GMT and 0", name, offset)
	}
}

// testOpts is what the flags default to, which is what the archive tests
// above are describing.
func testOpts() *opts {
	o, err := parseArgs([]string{"mkuinitramfs", "example.com/m/cmd"})
	if err != nil {
		panic(err)
	}
	return o
}

// parseArgs owns the flags, so a test can ask for a configuration without
// touching package state or the real command line.
func TestParseArgsReadsTheFlagsAndPatterns(t *testing.T) {
	o, err := parseArgs([]string{
		"mkuinitramfs", "-o", "out.cpio", "-initcmd", "myinit", "-defaultsh", "",
		"-nameserver", "1.1.1.1", "-keep-tmp",
		"example.com/m/cmd/...@v1.2.3", "example.com/other/cmd",
	})
	if err != nil {
		t.Fatalf("parseArgs: %v", err)
	}
	if o.output != "out.cpio" || o.initCmd != "myinit" || o.defaultSh != "" ||
		o.nameserver != "1.1.1.1" || !o.keepTmp {
		t.Errorf("opts = %+v, want the flags as given", o)
	}
	if len(o.patterns) != 2 || o.patterns[0].query != "v1.2.3" || o.patterns[1].path != "example.com/other/cmd" {
		t.Errorf("patterns = %+v, want both, with the version split off the first", o.patterns)
	}
}

// Naming no pattern is a usage error rather than an empty initramfs.
func TestParseArgsRefusesNoPatterns(t *testing.T) {
	if _, err := parseArgs([]string{"mkuinitramfs", "-o", "out.cpio"}); err == nil {
		t.Error("parseArgs with no pattern succeeded, want a usage error")
	}
}

// run returns its failures rather than exiting, which is what lets a test
// reach it at all, and what keeps the temporary directories collectable.
func TestRunReturnsUsageRatherThanExiting(t *testing.T) {
	var buf bytes.Buffer
	err := run(log.New(&buf, "", 0), []string{"mkuinitramfs"})
	if err == nil {
		t.Fatal("run with no pattern returned nil, want a usage error")
	}
	if !strings.Contains(err.Error(), "usage:") {
		t.Errorf("err = %q, want it to say usage", err)
	}
}

// The nameserver flag reaches the file the guest reads.
func TestNameserverFlagReachesResolvConf(t *testing.T) {
	o, err := parseArgs([]string{"mkuinitramfs", "-nameserver", "9.9.9.9", "example.com/m/cmd"})
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range o.skeleton() {
		if r.Name != "etc/resolv.conf" {
			continue
		}
		b := make([]byte, r.Info.FileSize)
		if _, err := r.ReaderAt.ReadAt(b, 0); err != nil {
			t.Fatal(err)
		}
		if got, want := string(b), "nameserver 9.9.9.9\n"; got != want {
			t.Errorf("resolv.conf = %q, want %q", got, want)
		}
		return
	}
	t.Error("the skeleton has no etc/resolv.conf")
}
