// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/u-root/gobusybox/src/pkg/bb/findpkg"
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
		got, err := parsePattern(tt.in)
		if err != nil {
			t.Errorf("parsePattern(%q): %v", tt.in, err)
			continue
		}
		if got.path != tt.path || got.query != tt.query {
			t.Errorf("parsePattern(%q) = {%q, %q}, want {%q, %q}", tt.in, got.path, got.query, tt.path, tt.query)
		}
	}
}

// A pattern that is not a package path has to be refused here, because go get
// and go list would read a leading "-" as a flag of their own.
func TestParsePatternRefusesWhatIsNotAPackagePath(t *testing.T) {
	for _, in := range []string{
		"",
		"   ",
		// An empty query: `go get pkg@` is not what anyone means.
		"example.com/m@",
		// A flag, which go get would act on rather than reject.
		"-tags=foo",
		"--",
	} {
		if got, err := parsePattern(in); err == nil {
			t.Errorf("parsePattern(%q) = %+v, want an error", in, got)
		}
	}
}

// fetchSpec puts the pattern back together for `go get`, and must not leave a
// bare "@" when no version was asked for.
func TestFetchSpecOmitsAnEmptyQuery(t *testing.T) {
	for _, tt := range []struct{ in, want string }{
		{"example.com/m", "example.com/m"},
		{"example.com/m@v1.2.3", "example.com/m@v1.2.3"},
	} {
		p, err := parsePattern(tt.in)
		if err != nil {
			t.Fatalf("parsePattern(%q): %v", tt.in, err)
		}
		if got := p.fetchSpec(); got != tt.want {
			t.Errorf("fetchSpec(%q) = %q, want %q", tt.in, got, tt.want)
		}
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

// goCmd must keep the fetch on the module proxy, without disturbing the
// caller's idea of which of its modules are private.
//
// GONOPROXY=none is what sends every module through the proxy. GOPRIVATE has
// to survive, because it is the default for GONOSUMDB as well: clearing it
// sends a private module path to sum.golang.org, which leaks the name and
// then fails the fetch with a 404. That split is what `go help private`
// documents, and `go env GONOSUMDB` reports GOPRIVATE's value when GONOSUMDB
// is unset.
func TestGoCmdKeepsTheFetchOnTheProxy(t *testing.T) {
	t.Setenv("GOPRIVATE", "example.com/*")
	t.Setenv("GONOPROXY", "example.com/*")

	last := map[string]string{}
	for _, kv := range goCmd(buildEnv(t.TempDir()), "list").Env {
		if k, v, ok := strings.Cut(kv, "="); ok {
			// exec uses the last value for a repeated key.
			last[k] = v
		}
	}
	for k, want := range map[string]string{
		"GONOPROXY": "none",
		"GOVCS":     "*:off",
		"GOPRIVATE": "example.com/*",
	} {
		if last[k] != want {
			t.Errorf("%s = %q, want %q", k, last[k], want)
		}
	}
}

// A failed `go list` has to say why it failed.
//
// Output puts the diagnosis in ExitError.Stderr, and %w prints none of it, so
// wrapping the error alone reports a bare "exit status 1" for every cause
// there is.
func TestCommandsReportsWhatGoListSaid(t *testing.T) {
	dir := t.TempDir()
	// A go.mod the go command cannot parse: module-level, so -e does not
	// tolerate it the way it tolerates an unbuildable package.
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module \x00bad\ngo 1.27\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := commands(buildEnv(dir), []pattern{{path: "all"}})
	if err == nil {
		t.Fatal("commands on an unparseable go.mod succeeded, want an error")
	}
	if !strings.Contains(err.Error(), "go.mod") {
		t.Errorf("error = %q, want it to carry what go list wrote to stderr", err)
	}
}

// The archive must name every parent directory before whatever sits in it,
// and must not name "." at all.
//
// Linux will not create intermediate directories for an initramfs, so a
// record whose parent never appears is silently dropped. The usual answer is
// cpio.WriteRecordsAndDirs, which fills parents in -- but it fills them at
// 0o777, and the parent of a top-level name such as "init" is ".", so every
// archive would open by chmodding the root it unpacks into, taking it from
// the kernel's 0o1777 to 0o777 and losing the sticky bit. The skeleton
// carries every parent instead, which is what lets the write be a plain
// WriteRecords; this test is what keeps that true.
func TestArchiveCarriesEveryParentAndNeverTheRoot(t *testing.T) {
	recs := readBackArchive(t, testOpts(), []string{"example.com/cmd/init", "example.com/cmd/gosh"})

	dirs := map[string]bool{}
	for _, r := range recs {
		if r.Name == "." {
			t.Errorf("archive names %q, at mode 0%o: nothing may chmod the root", r.Name, r.Mode&0o7777)
			continue
		}
		if parent := path.Dir(r.Name); parent != "." && !dirs[parent] {
			t.Errorf("%q appears before its parent %q", r.Name, parent)
		}
		if r.Mode&cpio.S_IFMT == cpio.S_IFDIR {
			dirs[r.Name] = true
		}
	}
}

// /bbin/bb must be 0755 whatever the host's umask is.
//
// GetRecord takes the mode from the file the linker wrote, so it arrives as
// 0777 &^ umask, and MakeReproducible normalizes everything about a record
// except its Mode. A umask of 0o027 would ship a busybox no ordinary user in
// the guest can execute, and two hosts would build different archives from
// identical input.
func TestTheBusyboxIsExecutableWhateverTheHostUmaskIs(t *testing.T) {
	old := syscall.Umask(0o027)
	defer syscall.Umask(old)

	dir := t.TempDir()
	bbPath := filepath.Join(dir, "bb")
	f, err := os.OpenFile(bbPath, os.O_CREATE|os.O_WRONLY, 0o777)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString("a busybox"); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	if fi, err := os.Stat(bbPath); err != nil {
		t.Fatal(err)
	} else if fi.Mode().Perm() == 0o755 {
		t.Fatalf("the umask did not take: %v is already the wanted mode", fi.Mode())
	}

	recs, err := testOpts().archive(bbPath, []string{"example.com/cmd/init", "example.com/cmd/gosh"})
	if err != nil {
		t.Fatal(err)
	}
	cpio.MakeAllReproducible(recs)
	for _, r := range recs {
		if r.Name == "bbin/bb" {
			if want := uint64(cpio.S_IFREG | 0o755); r.Mode != want {
				t.Errorf("bbin/bb mode = 0%o, want 0%o", r.Mode, want)
			}
			return
		}
	}
	t.Error("no bbin/bb record in the archive")
}

// A write that fails must not leave a truncated archive looking like output.
func TestAFailedWriteLeavesNoArchive(t *testing.T) {
	out := filepath.Join(t.TempDir(), "initramfs.cpio")
	// A record claiming more content than its reader holds fails mid-write,
	// after the file has been created and partly written.
	recs := []cpio.Record{
		cpio.Directory("bbin", 0o755),
		{
			ReaderAt: strings.NewReader("short"),
			Info:     cpio.Info{Name: "bbin/bb", Mode: cpio.S_IFREG | 0o755, FileSize: 1 << 20},
		},
	}
	if _, err := writeArchive(out, recs); err == nil {
		t.Fatal("writeArchive succeeded on a short record, want an error")
	}
	if _, err := os.Stat(out); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("os.Stat(%q) = %v, want the truncated archive to have been removed", out, err)
	}
}

// readBackArchive writes the archive the way run does and reads every record
// out of the result, so a test sees what a guest would.
func readBackArchive(t *testing.T, o *opts, cmds []string) []cpio.Record {
	t.Helper()
	dir := t.TempDir()
	bbPath := filepath.Join(dir, "bb")
	if err := os.WriteFile(bbPath, []byte("a busybox"), 0o755); err != nil {
		t.Fatal(err)
	}
	recs, err := o.archive(bbPath, cmds)
	if err != nil {
		t.Fatal(err)
	}
	cpio.MakeAllReproducible(recs)

	out := filepath.Join(dir, "initramfs.cpio")
	if _, err := writeArchive(out, recs); err != nil {
		t.Fatal(err)
	}
	f, err := os.Open(out)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	rr, err := cpio.Newc.NewFileReader(f)
	if err != nil {
		t.Fatal(err)
	}
	var got []cpio.Record
	for {
		r, err := rr.ReadRecord()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, r)
	}
	return got
}

// The settings every phase has to agree on live on one Environ, not on this
// process.
//
// gobusybox runs its own lookups and compiles with append(os.Environ(),
// env.Env()...), so an Environ reaches the half this file does not run, which
// is what makes one shared value possible at all. A stray GOWORK would pick a
// workspace over the fetched module, and CGO_ENABLED would let discovery list
// a cgo-only command the build then skips, leaving a symlink to nothing.
func TestOneEnvironCarriesWhatEveryPhaseMustAgreeOn(t *testing.T) {
	t.Setenv("GOWORK", "/somewhere/go.work")
	t.Setenv("CGO_ENABLED", "1")
	t.Setenv("GO111MODULE", "off")

	last := map[string]string{}
	for _, kv := range buildEnv(t.TempDir()).Env() {
		if k, v, ok := strings.Cut(kv, "="); ok {
			last[k] = v
		}
	}
	for k, want := range map[string]string{
		"GOWORK":      "off",
		"CGO_ENABLED": "0",
		"GO111MODULE": "on",
	} {
		if last[k] != want {
			t.Errorf("Env() has %s=%q, want %q", k, last[k], want)
		}
	}
}

// Every phase must run the same go binary.
//
// Environ runs $GOROOT/bin/go, and GOROOT is baked in when this binary is
// built, so installing with GOTOOLCHAIN=go1.N -- which doc.go recommends,
// because the command type-checks in process -- makes that the toolchain's
// own go. exec.Command("go") would fetch and list with whatever is on PATH,
// which is the older one the install was meant to avoid.
func TestEveryPhaseRunsTheSameGo(t *testing.T) {
	// A decoy first on PATH, so this cannot pass merely because the go on
	// PATH happens to be the same one. It usually is, which is what hides
	// the divergence until someone installs with GOTOOLCHAIN set.
	decoyDir := t.TempDir()
	decoy := filepath.Join(decoyDir, "go")
	if err := os.WriteFile(decoy, []byte("#!/bin/sh\nexit 1\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", decoyDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	env := buildEnv(t.TempDir())
	if env.GOROOT == "" {
		t.Fatal("GOROOT is empty, so there is no toolchain to prefer")
	}
	want := filepath.Join(env.GOROOT, "bin", "go")
	got := goCmd(env, "list").Path
	if got == decoy {
		t.Errorf("go command = %q, the first go on PATH, want the toolchain's %q", got, want)
	} else if got != want {
		t.Errorf("go command = %q, want %q", got, want)
	}
}

// The package lookup finds go through PATH, so the toolchain has to be first
// on it.
//
// Everything this file runs goes through the Environ and so through
// $GOROOT/bin/go. gobusybox's lookup is go/packages underneath, which runs
// exec.Command("go"), and exec resolves a bare name against this process's
// PATH rather than the environment it passes to the child -- so an Environ
// cannot reach it. Measured over a whole build with a logging go first on
// PATH: 8 invocations took it before this, 5 with the Environ alone, and 0
// with this as well.
//
// Prepended, not replacing: Environ.Env() replaces PATH outright for its own
// children, and doing that to this process would take /usr/bin away from
// everything else.
func TestTheToolchainGoesFirstOnPath(t *testing.T) {
	t.Setenv("PATH", "/somewhere/else")

	env := buildEnv(t.TempDir())
	if env.GOROOT == "" {
		t.Fatal("GOROOT is empty, so there is no toolchain to prefer")
	}
	if err := preferOwnToolchain(env); err != nil {
		t.Fatalf("preferOwnToolchain: %v", err)
	}
	want := filepath.Join(env.GOROOT, "bin")
	got := os.Getenv("PATH")
	first, rest, _ := strings.Cut(got, string(os.PathListSeparator))
	if first != want {
		t.Errorf("PATH starts with %q, want %q", first, want)
	}
	if rest != "/somewhere/else" {
		t.Errorf("PATH keeps %q, want the caller's %q", rest, "/somewhere/else")
	}
}

// UROOT_SOURCE is defused by the lookup environment rather than by unsetting
// it, which is why nothing sets it any more.
//
// findpkg reads it in DefaultEnv and nowhere else, and run passes an empty
// findpkg.Env, so the variable is never consulted. Were it consulted, a
// u-root package path would turn into a file system lookup and a local
// checkout would be built instead of the named version.
func TestUROOTSOURCEIsReadOnlyByTheDefaultLookupEnv(t *testing.T) {
	t.Setenv("UROOT_SOURCE", "/somewhere/with/a/checkout")
	if got := (&findpkg.Env{}).URootSource; got != "" {
		t.Errorf("the empty findpkg.Env has URootSource %q, want it empty", got)
	}
	if got := findpkg.DefaultEnv().URootSource; got == "" {
		t.Error("findpkg.DefaultEnv no longer reads UROOT_SOURCE, so run's empty Env may no longer be what defuses it")
	}
}

// GOFLAGS reaches these commands, and a -mod in it breaks them.
//
// The compile is covered: gobusybox appends -mod to the go build arguments
// and an argument beats GOFLAGS. go get does not accept a -mod argument at
// all, so the environment is the only place to take it out. Everything else
// in GOFLAGS stays, because -tags decides which files a package is built
// from, and so which packages are commands: dropping it here would let the
// listing and the compile disagree about what is being built.
func TestOnlyTheModFlagIsTakenOutOfGoflags(t *testing.T) {
	for _, tt := range []struct{ in, want string }{
		{"", ""},
		{"-mod=vendor", ""},
		{"--mod=vendor", ""},
		{"-tags=foo", "-tags=foo"},
		{"-tags=foo -mod=vendor -x", "-tags=foo -x"},
		{"-mod=readonly -tags=a,b", "-tags=a,b"},
		// Not a -mod setting, and must survive.
		{"-modfile=x.mod", "-modfile=x.mod"},
	} {
		if got := withoutModFlag(tt.in); got != tt.want {
			t.Errorf("withoutModFlag(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}

	t.Setenv("GOFLAGS", "-mod=vendor -tags=keepme")
	last := map[string]string{}
	for _, kv := range goCmd(buildEnv(t.TempDir()), "list").Env {
		if k, v, ok := strings.Cut(kv, "="); ok {
			last[k] = v
		}
	}
	if got := last["GOFLAGS"]; got != "-tags=keepme" {
		t.Errorf("GOFLAGS on the command = %q, want %q", got, "-tags=keepme")
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

// An empty -nameserver means no resolv.conf at all. Writing "nameserver \n"
// instead gives the guest a file naming a nameserver called "", which is
// worse than having no file: a resolver reads it and fails rather than
// falling back.
func TestAnEmptyNameserverWritesNoResolvConf(t *testing.T) {
	o, err := parseArgs([]string{"mkuinitramfs", "-nameserver=", "example.com/m/cmd"})
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range o.skeleton() {
		if r.Name == "etc/resolv.conf" {
			t.Error("an empty -nameserver still wrote etc/resolv.conf")
		}
	}
}

// Asking for help is not a failure.
//
// The FlagSet has already printed the usage by the time Parse returns
// flag.ErrHelp, so there is nothing left to report; main must exit 0 rather
// than log "flag: help requested" and exit 1.
func TestHelpIsNotAnError(t *testing.T) {
	for _, flagName := range []string{"-h", "-help"} {
		var buf bytes.Buffer
		err := run(log.New(&buf, "", 0), []string{"mkuinitramfs", flagName})
		if !errors.Is(err, flag.ErrHelp) {
			t.Errorf("run(%s) = %v, want flag.ErrHelp so main can exit 0", flagName, err)
		}
	}
}

// parseArgs must not index an argument list that has nothing in it.
func TestParseArgsRefusesAnEmptyArgumentList(t *testing.T) {
	if _, err := parseArgs(nil); err == nil {
		t.Error("parseArgs(nil) = nil, want an error rather than a panic")
	}
}
