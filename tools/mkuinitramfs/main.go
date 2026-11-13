// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/u-root/gobusybox/src/pkg/bb"
	"github.com/u-root/gobusybox/src/pkg/bb/findpkg"
	"github.com/u-root/gobusybox/src/pkg/golang"
	"github.com/u-root/u-root/pkg/cpio"
)

// opts is what the command line asks for.
type opts struct {
	output     string
	initCmd    string
	defaultSh  string
	nameserver string
	keepTmp    bool
	patterns   []pattern
}

// parseArgs turns an argument list into opts, using a FlagSet of its own so
// that a test can call run without touching package state.
func parseArgs(args []string) (*opts, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("no arguments, not even a command name")
	}
	name := filepath.Base(args[0])
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	o := &opts{}
	fs.StringVar(&o.output, "o", "initramfs.cpio", "Path to write the initramfs to")
	fs.StringVar(&o.initCmd, "initcmd", "init", "Command to symlink /init to, or empty for none")
	fs.StringVar(&o.defaultSh, "defaultsh", "gosh", "Command to symlink /bin/sh to, or empty for none")
	fs.StringVar(&o.nameserver, "nameserver", "8.8.8.8", "Nameserver to write into /etc/resolv.conf, or empty for no resolv.conf")
	fs.BoolVar(&o.keepTmp, "keep-tmp", false, "Keep the fetching module and generated source")
	if err := fs.Parse(args[1:]); err != nil {
		return nil, err
	}
	if fs.NArg() == 0 {
		return nil, fmt.Errorf("usage: %s [flags] <pattern>[@version]...\n\ne.g. %s github.com/u-root/u-root/cmds/core/...@latest", name, name)
	}
	for _, arg := range fs.Args() {
		p, err := parsePattern(arg)
		if err != nil {
			return nil, err
		}
		o.patterns = append(o.patterns, p)
	}
	return o, nil
}

// synthModulePath is the module written to require the requested patterns.
//
// Nothing imports it and it is never published, so the path only has to be one
// no input could also declare.
const synthModulePath = "mkuinitramfs.u-root.com/fetch"

// gmt0 is a minimal TZif blob, so that formatting a time does not depend on
// the host having zoneinfo. Taken verbatim from mkuimage rather than retyped:
// an earlier hand-written copy of it did not parse, and nothing would have
// said so at runtime, because Go falls back to UTC in silence.
const gmt0 = "TZif2\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x04\x00\x00\x00\x00\x00\x00GMT\x00\x00\x00TZif2\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x01\x00\x00\x00\x04\xf8\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00GMT\x00\x00\x00\nGMT0\n"

// pattern is one command-line argument, split into the part naming packages
// and the part naming a version.
type pattern struct {
	path  string
	query string
}

// parsePattern splits an argument into a package pattern and a version query.
//
// A Go package path cannot contain "@", so everything after the last one is
// the query. Splitting on the last rather than the first matters for a path
// that already carries a major version suffix, such as gopkg.in/yaml.v2.
//
// An argument opening with "-" is refused rather than passed on: go get would
// read it as a flag of its own, so a typo would silently become an option to
// a subprocess the user did not know was there.
func parsePattern(arg string) (pattern, error) {
	p := pattern{path: strings.TrimSpace(arg)}
	if i := strings.LastIndex(p.path, "@"); i >= 0 {
		p.path, p.query = p.path[:i], p.path[i+1:]
		if p.query == "" {
			return pattern{}, fmt.Errorf("%q has an empty version query after \"@\"", arg)
		}
	}
	if p.path == "" {
		return pattern{}, fmt.Errorf("%q names no package", arg)
	}
	if strings.HasPrefix(p.path, "-") {
		return pattern{}, fmt.Errorf("%q starts with \"-\": a pattern is a package path, not a flag", arg)
	}
	return p, nil
}

// fetchSpec is the argument `go get` is given for this pattern.
func (p pattern) fetchSpec() string {
	if p.query == "" {
		return p.path
	}
	return p.path + "@" + p.query
}

// buildEnv is the one Go environment every phase of a build shares.
//
// Fetching, listing and compiling must agree about which go they are using
// and what it is looking at, and the way to make them agree is to give them
// one Environ rather than to set variables on this process. gobusybox builds
// both its lookups and its compile as append(os.Environ(), env.Env()...), so
// what is set here reaches the half this file does not run.
//
// Using it is also what stops the two halves running different go commands.
// Environ runs $GOROOT/bin/go, and GOROOT is baked into this binary when it
// is built, so a mkuinitramfs installed the way doc.go recommends --
// GOTOOLCHAIN=go1.N go install -- has a GOROOT naming that toolchain.
// exec.Command("go") would instead fetch and list with whatever go is on
// PATH, which is the older one the install was meant to get away from.
//
//   - GOWORK: a go.work above the temporary directory would otherwise decide
//     the build list for a module that is not one of its members.
//   - CGO_ENABLED: the busybox is built without cgo, so discovery has to run
//     without it too, or a cgo-only command is listed, never built, and left
//     as an archive symlink pointing at nothing.
//   - GO111MODULE: the fetching module is a module, and GO111MODULE=off in
//     the caller's environment otherwise fails the fetch outright with
//     "modules disabled by GO111MODULE=off".
//   - Mod: gobusybox passes this to the compile as a -mod argument, where it
//     beats anything in GOFLAGS. goCmd does the same job by other means.
//
// UROOT_SOURCE needs nothing here: it is read only by findpkg.DefaultEnv,
// and run passes an empty findpkg.Env instead.
func buildEnv(modDir string) *golang.Environ {
	return golang.Default(
		golang.WithWorkingDir(modDir),
		golang.DisableCGO(),
		golang.WithGOWORK("off"),
		golang.WithGO111MODULE("on"),
		golang.WithMod(golang.ModOnly),
	)
}

// goCmd returns a go command to run against the fetching module.
//
// Source comes from the module proxy, never from a version control tool.
// GOVCS is what establishes that: it refuses the fetch whatever the proxy is,
// and names the rule rather than a missing tool. GONOPROXY=none is what sends
// every module through the proxy, including the ones GOPRIVATE marks private.
// GOPROXY is deliberately left alone, so a caller with a mirror or an offline
// cache keeps it.
//
// GOPRIVATE is left alone too, because it is the default for GONOSUMDB as
// well as for GONOPROXY: clearing it would send a private module path to
// sum.golang.org, which both leaks the name and fails the fetch. Overriding
// GONOPROXY alone is the split `go help private` documents.
//
// exec uses the last value for a repeated key, so these override the caller's.
func goCmd(env *golang.Environ, gocmd string, args ...string) *exec.Cmd {
	cmd := env.GoCmd(gocmd, args...)
	cmd.Env = append(cmd.Env,
		"GONOPROXY=none",
		"GOVCS=*:off",
		"GOFLAGS="+withoutModFlag(os.Getenv("GOFLAGS")),
	)
	return cmd
}

// preferOwnToolchain puts this binary's own Go toolchain first on PATH.
//
// Everything this file runs goes through the Environ and so through
// $GOROOT/bin/go already. gobusybox's package lookup does not: it is
// go/packages underneath, which runs exec.Command("go"), and exec resolves a
// bare name against the PATH of *this* process rather than the one it is
// about to pass to the child. No Environ can reach that, so the only way to
// make the lookup agree with the fetch, the listing and the compile is to
// make the toolchain's go the first one found.
//
// Prepended rather than replacing PATH: Environ.Env() replaces it outright
// for its own children, which leaves them without /usr/bin, and nothing here
// wants to do that to this process.
//
// It also makes the command work when the caller has no go at all, which an
// installed binary often does -- $GOPATH/bin is on PATH far more often than a
// Go toolchain is.
func preferOwnToolchain(env *golang.Environ) error {
	if env.GOROOT == "" {
		return nil
	}
	bin := filepath.Join(env.GOROOT, "bin")
	if path := os.Getenv("PATH"); path != "" {
		bin += string(os.PathListSeparator) + path
	}
	if err := os.Setenv("PATH", bin); err != nil {
		return fmt.Errorf("putting %s first on PATH: %w", bin, err)
	}
	return nil
}

// withoutModFlag returns goflags with any -mod setting removed.
//
// The compile is already immune to one: gobusybox appends -mod to the go
// build arguments, and an argument beats GOFLAGS. These commands have no such
// cover, and go get does not accept a -mod argument at all, so the
// environment is the only place to take it out. Left in, GOFLAGS=-mod=vendor
// fails the listing with "inconsistent vendoring" over a vendor directory the
// fetching module was never going to have.
//
// Only -mod is dropped. The rest of GOFLAGS stays, because -tags decides
// which files a package is built from and so which packages are commands at
// all: stripping it here would let the listing and the compile disagree about
// what is being built.
//
// Splitting on spaces is how the go command itself reads GOFLAGS; it has no
// quoting.
func withoutModFlag(goflags string) string {
	var keep []string
	for f := range strings.FieldsSeq(goflags) {
		if strings.HasPrefix(f, "-mod=") || strings.HasPrefix(f, "--mod=") {
			continue
		}
		keep = append(keep, f)
	}
	return strings.Join(keep, " ")
}

// fetch writes a module requiring everything the patterns name, leaving the
// source in the module cache and the chosen versions in its go.mod.
func fetch(l *log.Logger, env *golang.Environ, modDir string, patterns []pattern) error {
	gomod := fmt.Sprintf("module %s\n", synthModulePath)
	if err := os.WriteFile(filepath.Join(modDir, "go.mod"), []byte(gomod), 0o644); err != nil {
		return fmt.Errorf("writing go.mod: %w", err)
	}

	specs := make([]string, 0, len(patterns))
	for _, p := range patterns {
		specs = append(specs, p.fetchSpec())
	}

	// One invocation rather than one per pattern, so the go command
	// resolves all of them against a single build list.
	l.Printf("Fetching %s", strings.Join(specs, " "))
	if out, err := goCmd(env, "get", specs...).CombinedOutput(); err != nil {
		return fmt.Errorf("`go get %s`: %w: %s", strings.Join(specs, " "), err, out)
	}
	return nil
}

// commands lists the import paths of the main packages the patterns match.
//
// -e keeps a pattern that matches some unbuildable package from failing the
// whole listing, matching the build, which skips what it cannot build rather
// than refusing to build anything.
func commands(env *golang.Environ, patterns []pattern) ([]string, error) {
	var out []string
	seen := map[string]bool{}
	for _, p := range patterns {
		cmd := goCmd(env, "list", "-e", "-f", `{{if eq .Name "main"}}{{.ImportPath}}{{end}}`, p.path)
		b, err := cmd.Output()
		if err != nil {
			// Output puts the diagnosis in ExitError.Stderr, which %w
			// does not print: without this the whole report of a
			// failed listing is "exit status 1".
			var ee *exec.ExitError
			if errors.As(err, &ee) && len(ee.Stderr) > 0 {
				return nil, fmt.Errorf("`go list %s`: %w: %s", p.path, err, ee.Stderr)
			}
			return nil, fmt.Errorf("`go list %s`: %w", p.path, err)
		}
		for line := range strings.SplitSeq(string(b), "\n") {
			if line = strings.TrimSpace(line); line != "" && !seen[line] {
				seen[line] = true
				out = append(out, line)
			}
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no Go commands match %v", patterns)
	}
	// Sorted, so that the symlink order in the archive depends on which
	// commands were found rather than on the order the patterns were typed.
	slices.Sort(out)
	return out, nil
}

// skeleton is the directory tree and device nodes a u-root initramfs has
// always carried, which init expects to find already there.
//
// Every parent directory any later record needs is in here, which is what
// lets the archive be written without filling parents in. Nothing names ".":
// a record for it would chmod the root the kernel unpacks into.
func (o *opts) skeleton() []cpio.Record {
	recs := []cpio.Record{
		cpio.Directory("bin", 0o755),
		cpio.Directory("bbin", 0o755),
		cpio.Directory("dev", 0o755),
		cpio.Directory("env", 0o755),
		cpio.Directory("etc", 0o755),
		cpio.Directory("lib64", 0o755),
		cpio.Directory("proc", 0o755),
		cpio.Directory("sys", 0o755),
		cpio.Directory("tcz", 0o755),
		cpio.Directory("tmp", 0o777),
		cpio.Directory("ubin", 0o755),
		cpio.Directory("usr", 0o755),
		cpio.Directory("usr/lib", 0o755),
		cpio.Directory("var", 0o755),
		cpio.Directory("var/log", 0o777),
		cpio.Directory("var/lock", 0o777),
		cpio.CharDev("dev/console", 0o600, 5, 1),
		cpio.CharDev("dev/tty", 0o666, 5, 0),
		cpio.CharDev("dev/null", 0o666, 1, 3),
		cpio.CharDev("dev/port", 0o640, 1, 4),
		cpio.CharDev("dev/urandom", 0o666, 1, 9),
		cpio.StaticFile("etc/localtime", gmt0, 0o644),
	}
	// An empty -nameserver means no resolv.conf, rather than one naming a
	// nameserver called "".
	if o.nameserver != "" {
		recs = append(recs, cpio.StaticFile("etc/resolv.conf", "nameserver "+o.nameserver+"\n", 0o644))
	}
	return recs
}

// symlinkTo makes a relative symlink at source pointing at the command named
// by cmd in bbin.
//
// Relative symlinks are evaluated from their parent directory, so the target
// is computed against that rather than against the root.
//
// filepath.Rel is the only way to compute this and it answers in the host's
// separators, so the result is put back into slashes: a cpio holds Unix paths
// whatever built it, and a Windows host would otherwise write "..\bbin\gosh".
func symlinkTo(source, cmd string) (cpio.Record, error) {
	target, err := filepath.Rel(path.Join("/", path.Dir(source)), path.Join("/bbin", cmd))
	if err != nil {
		return cpio.Record{}, err
	}
	return cpio.Symlink(source, filepath.ToSlash(target)), nil
}

// archive assembles every record the initramfs holds, in the order they are
// written.
func (o *opts) archive(bbPath string, cmds []string) ([]cpio.Record, error) {
	recs := o.skeleton()

	bb, err := cpio.NewRecorder().GetRecord(bbPath)
	if err != nil {
		return nil, fmt.Errorf("recording the busybox: %w", err)
	}
	bb.Name = "bbin/bb"
	// GetRecord takes the mode from the host file, which the linker created
	// under the caller's umask, and MakeReproducible does not normalize it.
	// Left alone, a umask of 0o027 ships a /bbin/bb no ordinary user in the
	// guest can execute, and two hosts build different archives from the
	// same input.
	bb.Mode = cpio.S_IFREG | 0o755
	recs = append(recs, bb)

	// One symlink per command, all pointing at the single binary, which is
	// what makes a busybox dispatch by argv[0]. path, not filepath: these
	// are archive names, which are Unix paths whatever host wrote them.
	for _, c := range cmds {
		name := path.Base(c)
		if name == "bb" {
			continue
		}
		recs = append(recs, cpio.Symlink(path.Join("bbin", name), "bb"))
	}

	have := map[string]bool{}
	for _, c := range cmds {
		have[path.Base(c)] = true
	}
	for _, s := range []struct{ source, cmd string }{
		{"init", o.initCmd},
		{"bin/sh", o.defaultSh},
		{"bin/defaultsh", o.defaultSh},
	} {
		if s.cmd == "" {
			continue
		}
		if !have[s.cmd] {
			return nil, fmt.Errorf("cannot link /%s to %q: no such command was built", s.source, s.cmd)
		}
		rec, err := symlinkTo(s.source, s.cmd)
		if err != nil {
			return nil, err
		}
		recs = append(recs, rec)
	}
	return recs, nil
}

// run assembles an initramfs from the patterns specified in the argument list.
func run(l *log.Logger, args []string) error {
	o, err := parseArgs(args)
	if err != nil {
		return err
	}

	out, err := filepath.Abs(o.output)
	if err != nil {
		return err
	}

	modDir, err := os.MkdirTemp("", "mkuinitramfs-mod-")
	if err != nil {
		return err
	}
	defer o.cleanup(l, modDir)

	genDir, err := os.MkdirTemp("", "mkuinitramfs-src-")
	if err != nil {
		return err
	}
	defer o.cleanup(l, genDir)

	// One environment for the fetch, the listing and the compile, so that
	// all three agree about which go is running and what it is looking at.
	env := buildEnv(modDir)
	if err := preferOwnToolchain(env); err != nil {
		return err
	}

	if err := fetch(l, env, modDir, o.patterns); err != nil {
		return fmt.Errorf("fetching source: %w", err)
	}

	cmds, err := commands(env, o.patterns)
	if err != nil {
		return fmt.Errorf("listing commands: %w", err)
	}
	l.Printf("Building %d command(s) into one binary", len(cmds))

	bbPath := filepath.Join(genDir, "bb")
	if err := bb.BuildBusybox(l, &bb.Opts{
		Env: env,
		// Empty rather than DefaultEnv, so UROOT_SOURCE in the caller's
		// environment cannot turn a u-root package path into a file
		// system lookup and build local source instead.
		LookupEnv:    &findpkg.Env{},
		GenSrcDir:    filepath.Join(genDir, "src"),
		CommandPaths: cmds,
		BinaryPath:   bbPath,
		GoBuildOpts:  &golang.BuildOpts{},
	}); err != nil {
		return fmt.Errorf("building the busybox: %w", err)
	}

	recs, err := o.archive(bbPath, cmds)
	if err != nil {
		return fmt.Errorf("assembling the initramfs: %w", err)
	}
	cpio.MakeAllReproducible(recs)

	size, err := writeArchive(out, recs)
	if err != nil {
		return err
	}
	l.Printf("Wrote %s (%d commands, %d bytes)", out, len(cmds), size)
	return nil
}

// writeArchive writes recs to out and returns the size of the result.
//
// A failed write leaves a truncated archive behind, which looks like output
// and is not, so out is removed unless the whole of it succeeded.
func writeArchive(out string, recs []cpio.Record) (int64, error) {
	f, err := os.Create(out)
	if err != nil {
		return 0, err
	}
	size, err := writeRecords(f, recs)
	if err != nil {
		f.Close()
		os.Remove(out)
		return 0, err
	}
	// Close explicitly rather than deferring it: a write is first reported
	// here on a network file system, and claiming success for a truncated
	// archive is worse than failing.
	if err := f.Close(); err != nil {
		os.Remove(out)
		return 0, fmt.Errorf("closing %s: %w", out, err)
	}
	return size, nil
}

// writeRecords writes every record and the trailer to f.
//
// cpio.Newc.Writer is already a DedupWriter and archive emits every parent
// directory itself, so the records go out as assembled. WriteRecordsAndDirs
// would fill parents in at 0o777, and the parent of a top-level name is ".",
// so every archive would open by chmodding the root it unpacks into. Measured
// in a guest: the rootfs is 0o1777 by default and such a record leaves it
// 0o777, dropping the sticky bit that stops one user deleting another's files
// in a directory everyone can write.
func writeRecords(f *os.File, recs []cpio.Record) (int64, error) {
	w := cpio.Newc.Writer(f)
	if err := cpio.WriteRecords(w, recs); err != nil {
		return 0, fmt.Errorf("writing the initramfs: %w", err)
	}
	if err := cpio.WriteTrailer(w); err != nil {
		return 0, fmt.Errorf("writing the trailer: %w", err)
	}
	st, err := f.Stat()
	if err != nil {
		return 0, err
	}
	return st.Size(), nil
}

// cleanup removes a temporary directory, or says where it was left.
func (o *opts) cleanup(l *log.Logger, dir string) {
	if o.keepTmp {
		l.Printf("Keeping %s", dir)
		return
	}
	if err := os.RemoveAll(dir); err != nil {
		l.Printf("Could not remove %s: %v", dir, err)
	}
}

func main() {
	l := log.New(os.Stderr, "", log.Ltime)
	if err := run(l, os.Args); err != nil {
		// -h and -help are answered by the FlagSet printing usage, so
		// there is nothing left to report and nothing went wrong.
		if errors.Is(err, flag.ErrHelp) {
			return
		}
		l.Fatal(err)
	}
}
