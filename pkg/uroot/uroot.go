// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uroot

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/u-root/u-root/pkg/cpio"
	"github.com/u-root/u-root/pkg/golang"
	"github.com/u-root/u-root/pkg/ldd"
)

// These constants are used in DefaultRamfs.
const (
	// This is the literal timezone file for GMT-0. Given that we have no
	// idea where we will be running, GMT seems a reasonable guess. If it
	// matters, setup code should download and change this to something
	// else.
	gmt0 = "TZif2\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x04\x00\x00\x00\x00\x00\x00GMT\x00\x00\x00TZif2\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x01\x00\x00\x00\x04\xf8\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00GMT\x00\x00\x00\nGMT0\n"

	nameserver = "nameserver 8.8.8.8\n"
)

var (
	builders = map[string]Build{
		"source": SourceBuild,
		"bb":     BBBuild,
		"binary": BinaryBuild,
	}
	archivers = map[string]Archiver{
		"cpio": CPIOArchiver{
			Format: "newc",
		},
		"dir": DirArchiver{},
	}
)

// DefaultRamfs are files that are contained in all u-root initramfs archives
// by default.
var DefaultRamfs = []cpio.Record{
	cpio.Directory("tcz", 0755),
	cpio.Directory("etc", 0755),
	cpio.Directory("dev", 0755),
	cpio.Directory("ubin", 0755),
	cpio.Directory("usr", 0755),
	cpio.Directory("usr/lib", 0755),
	cpio.Directory("lib64", 0755),
	cpio.Directory("bin", 0755),
	cpio.CharDev("dev/console", 0600, 5, 1),
	cpio.CharDev("dev/tty", 0666, 5, 0),
	cpio.CharDev("dev/null", 0666, 1, 3),
	cpio.CharDev("dev/port", 0640, 1, 4),
	cpio.CharDev("dev/urandom", 0666, 1, 9),
	cpio.StaticFile("etc/resolv.conf", nameserver, 0644),
	cpio.StaticFile("etc/localtime", gmt0, 0644),
}

// Opts are the arguments to CreateInitramfs.
type Opts struct {
	// Env is the build environment (OS, arch, etc).
	Env golang.Environ

	// Builder is the build format.
	//
	// This can currently be "source" or "bb".
	Builder Build

	// Archiver is the initramfs archival format.
	//
	// Only "cpio" is currently supported.
	Archiver Archiver

	// Packages are the Go packages to add to the archive.
	//
	// Currently allowed formats:
	//   Go package imports; e.g. github.com/u-root/u-root/cmds/ls
	//   Paths to Go package directories; e.g. $GOPATH/src/github.com/u-root/u-root/cmds/ls
	//   Globs of paths to Go package directories; e.g. ./cmds/*
	Packages []string

	// ExtraFiles are files to add to the archive in addition to the Go
	// packages.
	//
	// Shared library dependencies will automatically also be added to the
	// archive using ldd.
	ExtraFiles []string

	// TempDir is a temporary directory for the builder to store files in.
	TempDir string

	// OutputFile is the archive output file.
	OutputFile ArchiveWriter

	// BaseArchive is an existing initramfs to include in the resulting
	// initramfs.
	BaseArchive ArchiveReader

	// UseExistingInit determines whether the existing init from
	// BaseArchive should be used.
	//
	// If this is false, the "init" from BaseArchive will be renamed to
	// "inito".
	UseExistingInit bool
}

// CreateInitramfs creates an initramfs built to `opts`' specifications.
func CreateInitramfs(opts Opts) error {
	if _, err := os.Stat(opts.TempDir); os.IsNotExist(err) {
		return fmt.Errorf("temp dir %q must exist: %v", opts.TempDir, err)
	}
	if opts.OutputFile == nil {
		return fmt.Errorf("must give output file")
	}

	var importPaths []string
	// Resolve file system paths to package import paths.
	for _, pkg := range opts.Packages {
		matches, err := filepath.Glob(pkg)
		if len(matches) == 0 || err != nil {
			if _, perr := opts.Env.Package(pkg); perr != nil {
				return fmt.Errorf("%q is neither package or path/glob: %v / %v", pkg, err, perr)
			}
			importPaths = append(importPaths, pkg)
		}

		for _, match := range matches {
			p, err := opts.Env.PackageByPath(match)
			if err != nil {
				log.Printf("Skipping package %q: %v", match, err)
			} else {
				importPaths = append(importPaths, p.ImportPath)
			}
		}
	}

	builderTmpDir, err := ioutil.TempDir(opts.TempDir, "builder")
	if err != nil {
		return err
	}

	// Build the packages.
	bOpts := BuildOpts{
		Env:      opts.Env,
		Packages: importPaths,
		TempDir:  builderTmpDir,
	}
	files, err := opts.Builder(bOpts)
	if err != nil {
		return fmt.Errorf("error building %#v: %v", bOpts, err)
	}

	// Open the target initramfs file.
	archive := ArchiveOpts{
		ArchiveFiles:    files,
		OutputFile:      opts.OutputFile,
		BaseArchive:     opts.BaseArchive,
		UseExistingInit: opts.UseExistingInit,
		DefaultRecords:  DefaultRamfs,
	}

	// Add files from command line.
	for _, file := range opts.ExtraFiles {
		var src, dst string
		parts := strings.SplitN(file, ":", 2)
		if len(parts) == 2 {
			// treat the entry with the new src:dst syntax
			src = parts[0]
			dst = parts[1]
		} else {
			// plain old syntax
			src = file
			dst = file
		}
		src, err := filepath.Abs(src)
		if err != nil {
			return fmt.Errorf("couldn't find absolute path for %q: %v", src, err)
		}
		if err := archive.AddFile(src, dst); err != nil {
			return fmt.Errorf("couldn't add %q to archive: %v", file, err)
		}

		// Pull dependencies in the case of binaries. If `path` is not
		// a binary, `libs` will just be empty.
		libs, err := ldd.List([]string{src})
		if err != nil {
			return fmt.Errorf("couldn't list ldd dependencies for %q: %v", file, err)
		}
		for _, lib := range libs {
			if err := archive.AddFile(lib, lib[1:]); err != nil {
				return fmt.Errorf("couldn't add %q to archive: %v", lib, err)
			}
		}
	}

	// Finally, write the archive.
	if err := archive.Write(); err != nil {
		return fmt.Errorf("error archiving: %v", err)
	}
	return nil
}

// BuildOpts are arguments to the Build function.
type BuildOpts struct {
	// Env is the Go environment to use to compile and link packages.
	Env golang.Environ

	// Packages are the Go package import paths to compile.
	//
	// Builders need not support resolving packages by path.
	//
	// E.g. cmd/go or github.com/u-root/u-root/cmds/ls.
	Packages []string

	// TempDir is a temporary directory where the compilation mode compiled
	// binaries can be placed.
	//
	// TempDir should contain no files.
	TempDir string
}

// Build uses the given options to build Go packages and returns a list of
// files to be included in an initramfs archive.
type Build func(BuildOpts) (ArchiveFiles, error)

// ArchiveOpts are the options for building the initramfs archive.
type ArchiveOpts struct {
	// ArchiveFiles are the files to be included.
	//
	// Files in ArchiveFiles generally have priority over files in
	// DefaultRecords or BaseArchive.
	ArchiveFiles

	// DefaultRecords is a set of files to be included in the initramfs.
	DefaultRecords []cpio.Record

	// OutputFile is the file to write to.
	OutputFile ArchiveWriter

	// BaseArchive is an existing archive to add files to.
	//
	// BaseArchive may be nil.
	BaseArchive ArchiveReader

	// UseExistingInit determines whether the init from BaseArchive is used
	// or not, if BaseArchive is specified.
	//
	// If this is false, the "init" file in BaseArchive will be renamed
	// "inito" in the output archive.
	UseExistingInit bool
}

// Archiver is an archive format that builds an archive using a given set of
// files.
type Archiver interface {
	// OpenWriter opens an archive writer at `path`.
	//
	// If `path` is unspecified, implementations may choose an arbitrary
	// default location, potentially based on `goos` and `goarch`.
	OpenWriter(path, goos, goarch string) (ArchiveWriter, error)

	// Reader returns an ArchiveReader wrapper using the given io.Reader.
	Reader(io.ReaderAt) ArchiveReader
}

// ArchiveWriter is an object that files can be written to.
type ArchiveWriter interface {
	// WriteRecord writes the given file record.
	WriteRecord(cpio.Record) error

	// Finish finishes the archive.
	Finish() error
}

// ArchiveReader is an object that files can be read from.
type ArchiveReader interface {
	// ReadRecord reads a file record.
	ReadRecord() (cpio.Record, error)
}

// GetBuilder returns the Build function for the named build mode.
func GetBuilder(name string) (Build, error) {
	build, ok := builders[name]
	if !ok {
		return nil, fmt.Errorf("couldn't find builder %q", name)
	}
	return build, nil
}

// GetArchiver returns the archive mode for the named archive.
func GetArchiver(name string) (Archiver, error) {
	archiver, ok := archivers[name]
	if !ok {
		return nil, fmt.Errorf("couldn't find archival format %q", name)
	}
	return archiver, nil
}

// DefaultPackageImports returns a list of default u-root packages to include.
func DefaultPackageImports(env golang.Environ) ([]string, error) {
	// Find u-root directory.
	urootPkg, err := env.Package("github.com/u-root/u-root")
	if err != nil {
		return nil, fmt.Errorf("Couldn't find u-root src directory: %v", err)
	}

	matches, err := filepath.Glob(filepath.Join(urootPkg.Dir, "cmds/*"))
	if err != nil {
		return nil, fmt.Errorf("couldn't find u-root cmds: %v", err)
	}
	pkgs := make([]string, 0, len(matches))
	for _, match := range matches {
		pkg, err := env.PackageByPath(match)
		if err == nil {
			pkgs = append(pkgs, pkg.ImportPath)
		}
	}
	return pkgs, nil
}
