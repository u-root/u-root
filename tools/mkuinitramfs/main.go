// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command mkuinitramfs builds a u-root initramfs from the Go module cache,
// without a source checkout.
//
//	mkuinitramfs -o initramfs.cpio github.com/u-root/u-root/cmds/core/...@v0.16.0
//
// It writes a module whose only purpose is to require what the patterns name,
// which is what turns a package pattern into the module that provides it.
// That module stands in for the checkout: it holds no source, only the go.mod
// and go.sum recording what the patterns resolved to, which is enough for the
// go command to list the packages and build them. The source stays in the
// module cache.
//
// The archive is assembled with u-root's own pkg/cpio rather than a copy of
// it, so the fixes that land there reach this tool.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
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
	name := filepath.Base(args[0])
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	o := &opts{}
	fs.StringVar(&o.output, "o", "initramfs.cpio", "Path to write the initramfs to")
	fs.StringVar(&o.initCmd, "initcmd", "init", "Command to symlink /init to, or empty for none")
	fs.StringVar(&o.defaultSh, "defaultsh", "gosh", "Command to symlink /bin/sh to, or empty for none")
	fs.StringVar(&o.nameserver, "nameserver", "8.8.8.8", "Nameserver to write into /etc/resolv.conf")
	fs.BoolVar(&o.keepTmp, "keep-tmp", false, "Keep the fetching module and generated source")
	if err := fs.Parse(args[1:]); err != nil {
		return nil, err
	}
	if fs.NArg() == 0 {
		return nil, fmt.Errorf("usage: %s [flags] <pattern>[@version]...\n\ne.g. %s github.com/u-root/u-root/cmds/core/...@latest", name, name)
	}
	for _, arg := range fs.Args() {
		o.patterns = append(o.patterns, parsePattern(arg))
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
var gmt0 = "TZif2\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x04\x00\x00\x00\x00\x00\x00GMT\x00\x00\x00TZif2\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x01\x00\x00\x00\x04\xf8\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00GMT\x00\x00\x00\nGMT0\n"

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
func parsePattern(arg string) pattern {
	p := pattern{path: strings.TrimSpace(arg)}
	if i := strings.LastIndex(p.path, "@"); i >= 0 {
		p.path, p.query = p.path[:i], p.path[i+1:]
	}
	return p
}

// fetchSpec is the argument `go get` is given for this pattern.
func (p pattern) fetchSpec() string {
	if p.query == "" {
		return p.path
	}
	return p.path + "@" + p.query
}

// goCmd returns a go command to run against the fetching module.
//
// Source comes from the module proxy, never from a version control tool.
// GOVCS is what establishes that: it refuses the fetch whatever the proxy is,
// and names the rule rather than a missing tool. GOPRIVATE and GONOPROXY are
// cleared because both name paths that bypass the proxy and both are read from
// the caller's environment. GOPROXY is deliberately left alone, so a caller
// with a mirror or an offline cache keeps it.
//
// exec uses the last value for a repeated key, so these override the caller's.
func goCmd(dir, gocmd string, args ...string) *exec.Cmd {
	cmd := exec.Command("go", append([]string{gocmd}, args...)...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GOPRIVATE=",
		"GONOPROXY=",
		"GOVCS=*:off",
	)
	return cmd
}

// fixEnvironment sets, for this process, what every go command has to agree
// on. gobusybox runs its own lookups and compiles with
// append(os.Environ(), ...), so anything set only on the commands here would
// not reach them, and the two halves could disagree about which source they
// are looking at.
//
//   - GOWORK: a go.work above the temporary directory would otherwise decide
//     the build list for a module that is not one of its members.
//   - UROOT_SOURCE: findpkg rewrites a u-root package path into a file system
//     lookup when it is set, which would quietly build a local checkout.
//   - CGO_ENABLED: the busybox is built without cgo, so discovery has to run
//     without it too, or a cgo-only command is listed, never built, and left
//     as an archive symlink pointing at nothing.
func fixEnvironment() error {
	for k, v := range map[string]string{
		"GOWORK":       "off",
		"UROOT_SOURCE": "",
		"CGO_ENABLED":  "0",
	} {
		if err := os.Setenv(k, v); err != nil {
			return fmt.Errorf("setting %s: %w", k, err)
		}
	}
	return nil
}

// fetch writes a module requiring everything the patterns name, leaving the
// source in the module cache and the chosen versions in its go.mod.
func fetch(l *log.Logger, modDir string, patterns []pattern) error {
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
	if out, err := goCmd(modDir, "get", specs...).CombinedOutput(); err != nil {
		return fmt.Errorf("`go get %s`: %w: %s", strings.Join(specs, " "), err, out)
	}
	return nil
}

// commands lists the import paths of the main packages the patterns match.
//
// -e keeps a pattern that matches some unbuildable package from failing the
// whole listing, matching the build, which skips what it cannot build rather
// than refusing to build anything.
func commands(modDir string, patterns []pattern) ([]string, error) {
	var out []string
	seen := map[string]bool{}
	for _, p := range patterns {
		cmd := goCmd(modDir, "list", "-e", "-f", `{{if eq .Name "main"}}{{.ImportPath}}{{end}}`, p.path)
		b, err := cmd.Output()
		if err != nil {
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
	return out, nil
}

// skeleton is the directory tree and device nodes a u-root initramfs has
// always carried, which init expects to find already there.
func (o *opts) skeleton() []cpio.Record {
	return []cpio.Record{
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
		cpio.StaticFile("etc/resolv.conf", "nameserver "+o.nameserver+"\n", 0o644),
		cpio.StaticFile("etc/localtime", string(gmt0), 0o644),
	}
}

// symlinkTo makes a relative symlink at source pointing at the command named
// by cmd in bbin.
//
// Relative symlinks are evaluated from their parent directory, so the target
// is computed against that rather than against the root.
func symlinkTo(source, cmd string) (cpio.Record, error) {
	target, err := filepath.Rel(filepath.Join("/", filepath.Dir(source)), filepath.Join("/bbin", cmd))
	if err != nil {
		return cpio.Record{}, err
	}
	return cpio.Symlink(source, target), nil
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
	recs = append(recs, bb)

	// One symlink per command, all pointing at the single binary, which is
	// what makes a busybox dispatch by argv[0].
	for _, c := range cmds {
		name := path.Base(c)
		if name == "bb" {
			continue
		}
		recs = append(recs, cpio.Symlink(filepath.Join("bbin", name), "bb"))
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

	if err := fixEnvironment(); err != nil {
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

	if err := fetch(l, modDir, o.patterns); err != nil {
		return fmt.Errorf("fetching source: %w", err)
	}

	cmds, err := commands(modDir, o.patterns)
	if err != nil {
		return fmt.Errorf("listing commands: %w", err)
	}
	l.Printf("Building %d command(s) into one binary", len(cmds))

	bbPath := filepath.Join(genDir, "bb")
	env := golang.Default(golang.WithWorkingDir(modDir))
	env.Context.CgoEnabled = false
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

	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()

	w := cpio.NewDedupWriter(cpio.Newc.Writer(f))
	if err := cpio.WriteRecordsAndDirs(w, recs); err != nil {
		return fmt.Errorf("writing the initramfs: %w", err)
	}
	if err := cpio.WriteTrailer(w); err != nil {
		return fmt.Errorf("writing the trailer: %w", err)
	}

	st, err := f.Stat()
	if err != nil {
		return err
	}
	l.Printf("Wrote %s (%d commands, %d bytes)", out, len(cmds), st.Size())
	return nil
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
		l.Fatal(err)
	}
}
