// Copyright 2015-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	gbbgolang "github.com/u-root/gobusybox/src/pkg/golang"
	"github.com/u-root/u-root/pkg/shlex"
	"github.com/u-root/u-root/pkg/ulog"
	"github.com/u-root/u-root/pkg/uroot"
	"github.com/u-root/u-root/pkg/uroot/builder"
	"github.com/u-root/u-root/pkg/uroot/initramfs"
)

// multiFlag is used for flags that support multiple invocations, e.g. -files
type multiFlag []string

func (m *multiFlag) String() string {
	return fmt.Sprint(*m)
}

func (m *multiFlag) Set(value string) error {
	*m = append(*m, value)
	return nil
}

// errors from the u-root command
var (
	ErrEmptyFilesArg = errors.New("Empty argument to -files")
)

// Flags for u-root builder.
var (
	build, format, tmpDir, base, outputPath *string
	uinitCmd, initCmd                       *string
	defaultShell                            *string
	useExistingInit                         *bool
	noCommands                              *bool
	extraFiles                              multiFlag
	statsOutputPath                         *string
	statsLabel                              *string
	shellbang                               *bool
	tags                                    *string
	// For the new gobusybox support
	usegobusybox *bool
	genDir       *string
	// For the new "filepath only" logic
	urootSourceDir *string
)

func init() {
	var sh string
	switch gbbgolang.Default().GOOS {
	case "plan9":
		sh = ""
	default:
		sh = "elvish"
	}

	build = flag.String("build", "gbb", "u-root build format (e.g. bb/gbb or binary).")
	format = flag.String("format", "cpio", "Archival format.")

	tmpDir = flag.String("tmpdir", "", "Temporary directory to put binaries in.")

	base = flag.String("base", "", "Base archive to add files to. By default, this is a couple of directories like /bin, /etc, etc. u-root has a default internally supplied set of files; use base=/dev/null if you don't want any base files.")
	useExistingInit = flag.Bool("useinit", false, "Use existing init from base archive (only if --base was specified).")
	outputPath = flag.String("o", "", "Path to output initramfs file.")

	initCmd = flag.String("initcmd", "init", "Symlink target for /init. Can be an absolute path or a u-root command name. Use initcmd=\"\" if you don't want the symlink.")
	uinitCmd = flag.String("uinitcmd", "", "Symlink target and arguments for /bin/uinit. Can be an absolute path or a u-root command name. Use uinitcmd=\"\" if you don't want the symlink. E.g. -uinitcmd=\"echo foobar\"")
	defaultShell = flag.String("defaultsh", sh, "Default shell. Can be an absolute path or a u-root command name. Use defaultsh=\"\" if you don't want the symlink.")

	noCommands = flag.Bool("nocmd", false, "Build no Go commands; initramfs only")

	flag.Var(&extraFiles, "files", "Additional files, directories, and binaries (with their ldd dependencies) to add to archive. Can be specified multiple times.")

	shellbang = flag.Bool("shellbang", false, "Use #! instead of symlinks for busybox")

	statsOutputPath = flag.String("stats-output-path", "", "Write build stats to this file (JSON)")
	statsLabel = flag.String("stats-label", "", "Use this statsLabel when writing stats")

	tags = flag.String("tags", "", "Comma separated list of build tags")

	// Flags for the gobusybox, which we hope to move to, since it works with modules.
	genDir = flag.String("gen-dir", "", "Directory to generate source in")

	// Flag for the new filepath only mode. This will be required to find the u-root commands and make templates work
	// In almost every case, "." is fine.
	urootSourceDir = flag.String("uroot-source", ".", "Path to the locally checked out u-root source tree in case commands from there are desired.")
}

type buildStats struct {
	Label      string  `json:"label,omitempty"`
	Time       int64   `json:"time"`
	Duration   float64 `json:"duration"`
	OutputSize int64   `json:"output_size"`
}

func writeBuildStats(stats buildStats, path string) error {
	var allStats []buildStats
	if data, err := os.ReadFile(*statsOutputPath); err == nil {
		json.Unmarshal(data, &allStats)
	}
	found := false
	for i, s := range allStats {
		if s.Label == stats.Label {
			allStats[i] = stats
			found = true
			break
		}
	}
	if !found {
		allStats = append(allStats, stats)
		sort.Slice(allStats, func(i, j int) bool {
			return strings.Compare(allStats[i].Label, allStats[j].Label) == -1
		})
	}
	data, err := json.MarshalIndent(allStats, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(*statsOutputPath, data, 0o644); err != nil {
		return err
	}
	return nil
}

func generateLabel(env *gbbgolang.Environ) string {
	var baseCmds []string
	if len(flag.Args()) > 0 {
		// Use the last component of the name to keep the label short
		for _, e := range flag.Args() {
			baseCmds = append(baseCmds, path.Base(e))
		}
	} else {
		baseCmds = []string{"core"}
	}
	return fmt.Sprintf("%s-%s-%s-%s", *build, env.GOOS, env.GOARCH, strings.Join(baseCmds, "_"))
}

// checkArgs checks for common mistakes that cause confusion.
//  1. -files as the last argument
//  2. -files followed by any switch, indicating a shell expansion problem
//     This is usually caused by Makfiles structured as follows
//     u-root -files `which ethtool` -files `which bash`
//     if ethtool is not installed, the expansion yields
//     u-root -files -files `which bash`
//     and the rather confusing error message
//     16:14:51 Skipping /usr/bin/bash because it is not a directory
//     which, in practice, nobody understands
func checkArgs(args ...string) error {
	if len(args) == 0 {
		return nil
	}

	if args[len(args)-1] == "-files" {
		return fmt.Errorf("last argument is -files:%w", ErrEmptyFilesArg)
	}

	// We know the last arg is not -files; scan the arguments for -files
	// followed by a switch.
	for i := 0; i < len(args)-1; i++ {
		if args[i] == "-files" && args[i+1][0] == '-' {
			return fmt.Errorf("-files argument %d is followed by a switch: %w", i, ErrEmptyFilesArg)
		}
	}

	return nil
}

func main() {
	if err := checkArgs(os.Args...); err != nil {
		log.Fatal(err)
	}
	gbbOpts := &gbbgolang.BuildOpts{}
	gbbOpts.RegisterFlags(flag.CommandLine)

	l := log.New(os.Stderr, "", log.Ltime)

	// Register an alias for -go-no-strip for backwards compatibility.
	flag.CommandLine.BoolVar(&gbbOpts.NoStrip, "no-strip", false, "Build unstripped binaries")
	flag.Parse()

	if usrc := os.Getenv("UROOT_SOURCE"); usrc != "" && *urootSourceDir == "" {
		*urootSourceDir = usrc
	}

	env := gbbgolang.Default()
	env.BuildTags = strings.Split(*tags, ",")
	if env.CgoEnabled {
		l.Printf("Disabling CGO for u-root...")
		env.CgoEnabled = false
	}
	l.Printf("Build environment: %s", env)
	if env.GOOS != "linux" {
		l.Printf("GOOS is not linux. Did you mean to set GOOS=linux?")
	}

	start := time.Now()

	// Main is in a separate functions so defers run on return.
	if err := Main(l, env, gbbOpts); err != nil {
		l.Fatalf("Build error: %v", err)
	}

	elapsed := time.Now().Sub(start)

	stats := buildStats{
		Label:    *statsLabel,
		Time:     start.Unix(),
		Duration: float64(elapsed.Milliseconds()) / 1000,
	}
	if stats.Label == "" {
		stats.Label = generateLabel(env)
	}
	if stat, err := os.Stat(*outputPath); err == nil && stat.ModTime().After(start) {
		l.Printf("Successfully built %q (size %d).", *outputPath, stat.Size())
		stats.OutputSize = stat.Size()
		if *statsOutputPath != "" {
			if err := writeBuildStats(stats, *statsOutputPath); err == nil {
				l.Printf("Wrote stats to %q (label %q)", *statsOutputPath, stats.Label)
			} else {
				l.Printf("Failed to write stats to %s: %v", *statsOutputPath, err)
			}
		}
	}
}

var recommendedVersions = []string{
	"go1.19",
	"go1.20",
}

func isRecommendedVersion(v string) bool {
	for _, r := range recommendedVersions {
		if strings.HasPrefix(v, r) {
			return true
		}
	}
	return false
}

func canFindSource(dir string) error {
	d := filepath.Join(dir, "cmds", "core")
	if _, err := os.Stat(d); err != nil {
		return fmt.Errorf("can not build u-root in %q:%w (-uroot-source may be incorrect or not set)", *urootSourceDir, os.ErrNotExist)
	}
	return nil
}

// Main is a separate function so defers are run on return, which they wouldn't
// on exit.
func Main(l ulog.Logger, env *gbbgolang.Environ, buildOpts *gbbgolang.BuildOpts) error {
	v, err := env.Version()
	if err != nil {
		l.Printf("Could not get environment's Go version, using runtime's version: %v", err)
		v = runtime.Version()
	}
	if !isRecommendedVersion(v) {
		l.Printf(`WARNING: You are not using one of the recommended Go versions (have = %s, recommended = %v).
			Some packages may not compile.
			Go to https://golang.org/doc/install to find out how to install a newer version of Go,
			or use https://godoc.org/golang.org/dl/%s to install an additional version of Go.`,
			v, recommendedVersions, recommendedVersions[0])
	}

	archiver, err := initramfs.GetArchiver(*format)
	if err != nil {
		return err
	}

	// Open the target initramfs file.
	if *outputPath == "" {
		if len(env.GOOS) == 0 && len(env.GOARCH) == 0 {
			return fmt.Errorf("passed no path, GOOS, and GOARCH to CPIOArchiver.OpenWriter")
		}
		*outputPath = fmt.Sprintf("/tmp/initramfs.%s_%s.cpio", env.GOOS, env.GOARCH)
	}
	w, err := archiver.OpenWriter(l, *outputPath)
	if err != nil {
		return err
	}

	var baseFile initramfs.Reader
	if *base != "" {
		bf, err := os.Open(*base)
		if err != nil {
			return err
		}
		defer bf.Close()
		baseFile = archiver.Reader(bf)
	} else {
		baseFile = uroot.DefaultRamfs().Reader()
	}

	tempDir := *tmpDir
	if tempDir == "" {
		var err error
		tempDir, err = os.MkdirTemp("", "u-root")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tempDir)
	} else if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		if err := os.MkdirAll(tempDir, 0o755); err != nil {
			return fmt.Errorf("temporary directory %q did not exist; tried to mkdir but failed: %v", tempDir, err)
		}
	}

	var (
		c           []uroot.Commands
		initCommand = *initCmd
	)
	if !*noCommands {
		var b builder.Builder
		switch *build {
		case "bb", "gbb":
			l.Printf("NOTE: building with the new gobusybox; to get the old behavior check out commit 8b790de")
			b = builder.GBBBuilder{ShellBang: *shellbang}
		case "binary":
			b = builder.BinaryBuilder{}
		case "source":
			return fmt.Errorf("source mode has been deprecated")
		default:
			return fmt.Errorf("could not find builder %q", *build)
		}

		// Resolve globs into package imports.
		//
		// Currently allowed format:
		//   Paths to Go package directories; e.g. $GOPATH/src/github.com/u-root/u-root/cmds/*
		//   u-root templates; e.g. all, core, minimal (requires uroot-source be valid)
		//   Import paths of u-root commands; e.g. github.com/u-root/u-root/cmds/* (requires uroot-source)
		var pkgs []string
		for _, a := range flag.Args() {
			p, ok := templates[a]
			if !ok {
				if !validateArg(a) {
					l.Printf("%q is not a valid path, allowed are only existing relative or absolute file paths!", a)
					continue
				}
				pkgs = append(pkgs, a)
				continue
			}
			pkgs = append(pkgs, p...)
		}
		if len(pkgs) == 0 {
			pkgs = []string{"github.com/u-root/u-root/cmds/core/*"}
		}

		// The command-line tool only allows specifying one build mode
		// right now.
		c = append(c, uroot.Commands{
			Builder:  b,
			Packages: pkgs,
		})
	}

	opts := uroot.Opts{
		Env:             env,
		Commands:        c,
		UrootSource:     *urootSourceDir,
		TempDir:         tempDir,
		ExtraFiles:      extraFiles,
		OutputFile:      w,
		BaseArchive:     baseFile,
		UseExistingInit: *useExistingInit,
		InitCmd:         initCommand,
		DefaultShell:    *defaultShell,
		BuildOpts:       buildOpts,
	}
	uinitArgs := shlex.Argv(*uinitCmd)
	if len(uinitArgs) > 0 {
		opts.UinitCmd = uinitArgs[0]
	}
	if len(uinitArgs) > 1 {
		opts.UinitArgs = uinitArgs[1:]
	}
	return uroot.CreateInitramfs(l, opts)
}

func validateArg(arg string) bool {
	// Do the simple thing first: stat the path.
	// This saves incorrect diagnostics when the
	// path is a perfectly valid relative path.
	if _, err := os.Stat(arg); err == nil {
		return true
	}
	if !checkPrefix(arg) {
		paths, err := filepath.Glob(arg)
		if err != nil {
			return false
		}
		for _, path := range paths {
			if !checkPrefix(path) {
				return false
			}
		}
	}

	return true
}

func checkPrefix(arg string) bool {
	prefixes := []string{".", "/", "-", "cmds", "github.com/u-root/u-root"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(arg, prefix) {
			return true
		}
	}

	return false
}
