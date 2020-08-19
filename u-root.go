// Copyright 2015-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/u-root/u-root/pkg/golang"
	"github.com/u-root/u-root/pkg/shlex"
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

// Flags for u-root builder.
var (
	build, format, tmpDir, base, outputPath *string
	uinitCmd, initCmd                       *string
	defaultShell                            *string
	useExistingInit                         *bool
	fourbins                                *bool
	noCommands                              *bool
	extraFiles                              multiFlag
	noStrip                                 *bool
)

func init() {
	var sh string
	switch golang.Default().GOOS {
	case "plan9":
		sh = ""
	default:
		sh = "elvish"
	}

	fourbins = flag.Bool("fourbins", false, "build installcommand on boot, no ahead of time, so we have only four binares")
	build = flag.String("build", "bb", "u-root build format (e.g. bb or source).")
	format = flag.String("format", "cpio", "Archival format.")

	tmpDir = flag.String("tmpdir", "", "Temporary directory to put binaries in.")

	base = flag.String("base", "", "Base archive to add files to. By default, this is a couple of directories like /bin, /etc, etc. u-root has a default internally supplied set of files; use base=/dev/null if you don't want any base files.")
	useExistingInit = flag.Bool("useinit", false, "Use existing init from base archive (only if --base was specified).")
	outputPath = flag.String("o", "", "Path to output initramfs file.")

	initCmd = flag.String("initcmd", "init", "Symlink target for /init. Can be an absolute path or a u-root command name. Use initcmd=\"\" if you don't want the symlink.")
	uinitCmd = flag.String("uinitcmd", "", "Symlink target and arguments for /bin/uinit. Can be an absolute path or a u-root command name. Use uinitcmd=\"\" if you don't want the symlink. E.g. -uinitcmd=\"echo foobar\"")
	defaultShell = flag.String("defaultsh", sh, "Default shell. Can be an absolute path or a u-root command name. Use defaultsh=\"\" if you don't want the symlink.")

	noCommands = flag.Bool("nocmd", false, "Build no Go commands; initramfs only")

	flag.Var(&extraFiles, "files", "Additional files, directories, and binaries (with their ldd dependencies) to add to archive. Can be speficified multiple times.")

	noStrip = flag.Bool("no-strip", false, "Build unstripped binaries")
}

func main() {
	flag.Parse()

	// Main is in a separate functions so defers run on return.
	if err := Main(); err != nil {
		log.Fatal(err)
	}
	log.Printf("Successfully wrote initramfs.")
}

var recommendedVersions = []string{
	"go1.13",
	"go1.14",
}

func isRecommendedVersion(v string) bool {
	for _, r := range recommendedVersions {
		if strings.HasPrefix(v, r) {
			return true
		}
	}
	return false
}

// Main is a separate function so defers are run on return, which they wouldn't
// on exit.
func Main() error {
	env := golang.Default()
	if *fourbins && env.GOROOT == "" {
		log.Fatalf("You have to set GOROOT for fourbins to work")
	}
	if env.CgoEnabled {
		log.Printf("Disabling CGO for u-root...")
		env.CgoEnabled = false
	}
	log.Printf("Build environment: %s", env)
	if env.GOOS != "linux" {
		log.Printf("GOOS is not linux. Did you mean to set GOOS=linux?")
	}

	v, err := env.Version()
	if err != nil {
		log.Printf("Could not get environment's Go version, using runtime's version: %v", err)
		v = runtime.Version()
	}
	if !isRecommendedVersion(v) {
		log.Printf(`WARNING: You are not using one of the recommended Go versions (have = %s, recommended = %v).
			Some packages may not compile.
			Go to https://golang.org/doc/install to find out how to install a newer version of Go,
			or use https://godoc.org/golang.org/dl/%s to install an additional version of Go.`,
			v, recommendedVersions, recommendedVersions[0])
	}

	archiver, err := initramfs.GetArchiver(*format)
	if err != nil {
		return err
	}

	logger := log.New(os.Stderr, "", log.LstdFlags)
	// Open the target initramfs file.
	w, err := archiver.OpenWriter(logger, *outputPath, env.GOOS, env.GOARCH)
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
		tempDir, err = ioutil.TempDir("", "u-root")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tempDir)
	} else if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		if err := os.MkdirAll(tempDir, 0755); err != nil {
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
		case "bb":
			b = builder.BBBuilder{}
		case "binary":
			b = builder.BinaryBuilder{}
		case "source":
			b = builder.SourceBuilder{
				FourBins: *fourbins,
			}
		default:
			return fmt.Errorf("could not find builder %q", *build)
		}

		// Resolve globs into package imports.
		//
		// Currently allowed formats:
		//   Go package imports; e.g. github.com/u-root/u-root/cmds/ls (must be in $GOPATH)
		//   Paths to Go package directories; e.g. $GOPATH/src/github.com/u-root/u-root/cmds/*
		var pkgs []string
		for _, a := range flag.Args() {
			p, ok := templates[a]
			if !ok {
				pkgs = append(pkgs, a)
				continue
			}
			pkgs = append(pkgs, p...)
		}
		if len(pkgs) == 0 {
			pkgs = []string{"github.com/u-root/u-root/cmds/core/*"}
		}

		if *fourbins && *build == "source" {
			initCommand = "/go/bin/go"
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
		TempDir:         tempDir,
		ExtraFiles:      extraFiles,
		OutputFile:      w,
		BaseArchive:     baseFile,
		UseExistingInit: *useExistingInit,
		InitCmd:         initCommand,
		DefaultShell:    *defaultShell,
		NoStrip:         *noStrip,
	}
	uinitArgs := shlex.Argv(*uinitCmd)
	if len(uinitArgs) > 0 {
		opts.UinitCmd = uinitArgs[0]
	}
	if len(uinitArgs) > 1 {
		opts.UinitArgs = uinitArgs[1:]
	}
	return uroot.CreateInitramfs(logger, opts)
}
