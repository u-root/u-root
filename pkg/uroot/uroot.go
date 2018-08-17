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
	"path"
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
	builders = map[string]Builder{
		"source": SourceBuilder,
		"bb":     BBBuilder,
		"binary": BinaryBuilder,
	}
	archivers = map[string]Archiver{
		"cpio": CPIOArchiver{
			RecordFormat: cpio.Newc,
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
	cpio.Directory("tmp", 0777),
	cpio.Directory("ubin", 0755),
	cpio.Directory("usr", 0755),
	cpio.Directory("usr/lib", 0755),
	cpio.Directory("var/log", 0777),
	cpio.Directory("lib64", 0755),
	cpio.CharDev("dev/console", 0600, 5, 1),
	cpio.CharDev("dev/tty", 0666, 5, 0),
	cpio.CharDev("dev/null", 0666, 1, 3),
	cpio.CharDev("dev/port", 0640, 1, 4),
	cpio.CharDev("dev/urandom", 0666, 1, 9),
	cpio.StaticFile("etc/resolv.conf", nameserver, 0644),
	cpio.StaticFile("etc/localtime", gmt0, 0644),
}

// Commands specifies a list of packages to build with a specific builder.
type Commands struct {
	// Builder is the build format.
	Builder Builder

	// Packages are the Go packages to add to the archive.
	//
	// Currently allowed formats:
	//   Go package imports; e.g. github.com/u-root/u-root/cmds/ls
	//   Paths to Go package directories; e.g. $GOPATH/src/github.com/u-root/u-root/cmds/ls
	//   Globs of paths to Go package directories; e.g. ./cmds/*
	Packages []string

	// BinaryDir is the directory in which the resulting binaries are
	// placed inside the initramfs.
	BinaryDir string
}

// TargetDir returns the binary directory for these Commands.
func (c Commands) TargetDir() string {
	if len(c.BinaryDir) != 0 {
		return c.BinaryDir
	}
	return c.Builder.DefaultBinaryDir
}

// Opts are the arguments to CreateInitramfs.
type Opts struct {
	// Env is the build environment (OS, arch, etc).
	Env golang.Environ

	// Commands specify packages to build using a specific builder.
	Commands []Commands

	// TempDir is a temporary directory for builders to store files in.
	TempDir string

	// ExtraFiles are files to add to the archive in addition to the Go
	// packages.
	//
	// Shared library dependencies will automatically also be added to the
	// archive using ldd.
	ExtraFiles []string

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

	// InitCmd is the name of a command to link /init to.
	//
	// This can be an absolute path or the name of a command included in
	// Commands.
	//
	// If this is empty, no init symlink will be created.
	InitCmd string

	// DefaultShell is the default shell to start after init.
	//
	// This can be an absolute path or the name of a command included in
	// Commands.
	//
	// If this is empty, there will be no default shell.
	DefaultShell string
}

// resolvePackagePath finds import paths for a single import path or directory string
func resolvePackagePath(env golang.Environ, pkg string) ([]string, error) {
	// Search the current working directory, as well GOROOT and GOPATHs
	prefixes := append([]string{""}, env.SrcDirs()...)
	// Resolve file system paths to package import paths.
	for _, prefix := range prefixes {
		path := filepath.Join(prefix, pkg)
		matches, err := filepath.Glob(path)
		if len(matches) == 0 || err != nil {
			continue
		}

		var importPaths []string
		for _, match := range matches {
			p, err := env.PackageByPath(match)
			if err != nil {
				log.Printf("Skipping package %q: %v", match, err)
			} else if p.ImportPath == "." {
				// TODO: I do not completely understand why
				// this is triggered. This is only an issue
				// while this function is run inside the
				// process of a "go test".
				importPaths = append(importPaths, pkg)
			} else {
				importPaths = append(importPaths, p.ImportPath)
			}
		}
		return importPaths, nil
	}

	// No file import paths found. Check if pkg still resolves as a package name.
	if _, err := env.Package(pkg); err != nil {
		return nil, fmt.Errorf("%q is neither package or path/glob: %v", pkg, err)
	}
	return []string{pkg}, nil
}

func resolveCommandOrPath(cmd string, cmds []Commands) (string, error) {
	if filepath.IsAbs(cmd) {
		return cmd, nil
	}

	for _, c := range cmds {
		for _, p := range c.Packages {
			// Figure out which build mode the shell is in, and symlink to
			// that build modee
			if name := path.Base(p); name == cmd {
				return path.Join("/", c.TargetDir(), cmd), nil
			}
		}
	}

	return "", fmt.Errorf("command or path %q not included in u-root build", cmd)
}

// ResolvePackagePaths takes a list of Go package import paths and directories
// and turns them into exclusively import paths.
//
// Currently allowed formats:
//   Go package imports; e.g. github.com/u-root/u-root/cmds/ls
//   Paths to Go package directories; e.g. $GOPATH/src/github.com/u-root/u-root/cmds/ls
//   Globs of package imports, e.g. github.com/u-root/u-root/cmds/*
//   Globs of paths to Go package directories; e.g. ./cmds/*
func ResolvePackagePaths(env golang.Environ, pkgs []string) ([]string, error) {
	var importPaths []string
	for _, pkg := range pkgs {
		paths, err := resolvePackagePath(env, pkg)
		if err != nil {
			return nil, err
		}
		importPaths = append(importPaths, paths...)
	}
	return importPaths, nil
}

// ParseExtraFiles adds files from the extraFiles list to the archive, as
// parsed from the following formats:
//
// - hostPath:archivePath adds the file from hostPath at the relative archivePath in the archive.
// - justAPath is added to the archive under justAPath.
//
// ParseExtraFiles will also add ldd-listed dependencies if lddDeps is true.
func ParseExtraFiles(archive *ArchiveFiles, extraFiles []string, lddDeps bool) error {
	var err error
	// Add files from command line.
	for _, file := range extraFiles {
		var src, dst string
		parts := strings.SplitN(file, ":", 2)
		if len(parts) == 2 {
			// treat the entry with the new src:dst syntax
			src = filepath.Clean(parts[0])
			dst = filepath.Clean(parts[1])
		} else {
			// plain old syntax
			// filepath.Clean interprets an empty string as CWD for no good reason.
			if len(file) == 0 {
				continue
			}
			src = filepath.Clean(file)
			dst = src
			if filepath.IsAbs(dst) {
				dst, err = filepath.Rel("/", dst)
				if err != nil {
					return fmt.Errorf("cannot make path relative to /: %v: %v", dst, err)
				}
			}
		}
		src, err := filepath.Abs(src)
		if err != nil {
			return fmt.Errorf("couldn't find absolute path for %q: %v", src, err)
		}
		if err := archive.AddFile(src, dst); err != nil {
			return fmt.Errorf("couldn't add %q to archive: %v", file, err)
		}

		if lddDeps {
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
	}
	return nil
}

// CreateInitramfs creates an initramfs built to `opts`' specifications.
func CreateInitramfs(opts Opts) error {
	if _, err := os.Stat(opts.TempDir); err != nil {
		return fmt.Errorf("temp dir %q must exist and be accessible: %v", opts.TempDir, err)
	}
	if opts.OutputFile == nil {
		return fmt.Errorf("must give output file")
	}

	files := NewArchiveFiles()

	// Expand commands.
	for index, cmds := range opts.Commands {
		importPaths, err := ResolvePackagePaths(opts.Env, cmds.Packages)
		if err != nil {
			return err
		}
		opts.Commands[index].Packages = importPaths
	}

	// Add each build mode's commands to the archive.
	for _, cmds := range opts.Commands {
		builderTmpDir, err := ioutil.TempDir(opts.TempDir, "builder")
		if err != nil {
			return err
		}

		// Build packages.
		bOpts := BuildOpts{
			Env:       opts.Env,
			Packages:  cmds.Packages,
			TempDir:   builderTmpDir,
			BinaryDir: cmds.TargetDir(),
		}
		if err := cmds.Builder.Build(files, bOpts); err != nil {
			return fmt.Errorf("error building %#v: %v", bOpts, err)
		}
	}

	// Open the target initramfs file.
	archive := ArchiveOpts{
		ArchiveFiles:    files,
		OutputFile:      opts.OutputFile,
		BaseArchive:     opts.BaseArchive,
		UseExistingInit: opts.UseExistingInit,
		DefaultRecords:  DefaultRamfs,
	}

	if len(opts.DefaultShell) > 0 {
		if target, err := resolveCommandOrPath(opts.DefaultShell, opts.Commands); err != nil {
			log.Printf("No default shell: %v", err)
		} else if err := archive.AddRecord(cpio.Symlink("bin/defaultsh", target)); err != nil {
			return err
		}
	}

	if len(opts.InitCmd) > 0 {
		if target, err := resolveCommandOrPath(opts.InitCmd, opts.Commands); err != nil {
			return fmt.Errorf("could not find init: %v", err)
		} else if err := archive.AddRecord(cpio.Symlink("init", target)); err != nil {
			return err
		}
	}

	if err := ParseExtraFiles(archive.ArchiveFiles, opts.ExtraFiles, true); err != nil {
		return err
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

	// BinaryDir is the directory that built binaries are placed in in the
	// initramfs.
	//
	// BinaryDir must be specified.
	BinaryDir string
}

// Builder uses the given options to build Go packages and adds its files to be
// included in the initramfs to the given ArchiveFiles.
type Builder struct {
	// Build uses the given options to build Go packages and adds its files to be
	// included in the initramfs to the given ArchiveFiles.
	Build func(*ArchiveFiles, BuildOpts) error

	// DefaultBinaryDir is the default binary directory to place built
	// binaries in.
	DefaultBinaryDir string
}

// ArchiveOpts are the options for building the initramfs archive.
type ArchiveOpts struct {
	// ArchiveFiles are the files to be included.
	//
	// Files in ArchiveFiles generally have priority over files in
	// DefaultRecords or BaseArchive.
	*ArchiveFiles

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
	cpio.RecordWriter

	// Finish finishes the archive.
	Finish() error
}

// ArchiveReader is an object that files can be read from.
type ArchiveReader cpio.RecordReader

// GetBuilder returns the Build function for the named build mode.
func GetBuilder(name string) (Builder, error) {
	build, ok := builders[name]
	if !ok {
		return Builder{}, fmt.Errorf("couldn't find builder %q", name)
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
