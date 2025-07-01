// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package uimage creates root file systems from Go programs.
//
// uimage will appropriately compile the Go programs, create symlinks for their
// names, and assemble an initramfs with additional files as specified.
package uimage

import (
	"debug/elf"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/hugelgupf/go-shlex"
	"github.com/u-root/gobusybox/src/pkg/bb/findpkg"
	"github.com/u-root/gobusybox/src/pkg/golang"
	"github.com/u-root/mkuimage/cpio"
	"github.com/u-root/mkuimage/fileflag"
	"github.com/u-root/mkuimage/ldd"
	"github.com/u-root/mkuimage/uimage/builder"
	"github.com/u-root/mkuimage/uimage/initramfs"
	"github.com/u-root/uio/llog"
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

// DefaultRamfs returns a cpio.Archive for the target OS.
// If an OS is not known it will return a reasonable u-root specific
// default.
func DefaultRamfs() *cpio.Archive {
	switch golang.Default().GOOS {
	case "linux":
		a, _ := cpio.ArchiveFromRecords([]cpio.Record{
			cpio.Directory("bin", 0o755),
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
			cpio.StaticFile("etc/resolv.conf", nameserver, 0o644),
			cpio.StaticFile("etc/localtime", gmt0, 0o644),
		})
		return a
	default:
		a, _ := cpio.ArchiveFromRecords([]cpio.Record{
			cpio.Directory("ubin", 0o755),
			cpio.Directory("bbin", 0o755),
		})
		return a
	}
}

// Commands specifies a list of Golang packages to build with a builder, e.g.
// in busybox mode, source mode, or binary mode.
//
// See Builder for an explanation of build modes.
type Commands struct {
	// Builder is the Go compiler mode.
	Builder builder.Builder

	// Packages are the Go commands to include (compiled or otherwise) and
	// add to the archive.
	//
	// Currently allowed formats:
	//
	//   - package imports; e.g. github.com/u-root/u-root/cmds/ls
	//   - globs of package imports; e.g. github.com/u-root/u-root/cmds/*
	//   - paths to package directories; e.g. $GOPATH/src/github.com/u-root/u-root/cmds/ls
	//   - globs of paths to package directories; e.g. ./cmds/*
	//
	// Directories may be relative or absolute, with or without globs.
	// Globs are resolved using filepath.Glob.
	Packages []string

	// BinaryDir is the directory in which the resulting binaries are
	// placed inside the initramfs.
	//
	// BinaryDir may be empty, in which case Builder.DefaultBinaryDir()
	// will be used.
	BinaryDir string

	// Build options for building go binaries. Ultimate this holds all the
	// args that end up being passed to `go build`.
	BuildOpts *golang.BuildOpts
}

// TargetDir returns the initramfs binary directory for these Commands.
func (c Commands) TargetDir() string {
	if len(c.BinaryDir) != 0 {
		return c.BinaryDir
	}
	return c.Builder.DefaultBinaryDir()
}

// Opts are the arguments to CreateInitramfs.
//
// Opts contains everything that influences initramfs creation such as the Go
// build environment.
type Opts struct {
	// Env is the Golang build environment (GOOS, GOARCH, etc).
	//
	// If nil, golang.Default is used.
	Env *golang.Environ

	// Commands specify packages to build using a specific builder.
	//
	// E.g. the following will build 'ls' and 'ip' in busybox mode, but
	// 'cd' and 'cat' as separate binaries. 'cd', 'cat', 'bb', and symlinks
	// from 'ls' and 'ip' will be added to the final initramfs.
	//
	//   []Commands{
	//     Commands{
	//       Builder: builder.Busybox,
	//       Packages: []string{
	//         "github.com/u-root/u-root/cmds/ls",
	//         "github.com/u-root/u-root/cmds/ip",
	//       },
	//     },
	//     Commands{
	//       Builder: builder.Binary,
	//       Packages: []string{
	//         "github.com/u-root/u-root/cmds/cd",
	//         "github.com/u-root/u-root/cmds/cat",
	//       },
	//     },
	//   }
	Commands []Commands

	// UrootSource is the filesystem path to the locally checked out
	// u-root source tree. This is needed to resolve templates or
	// import paths of u-root commands.
	UrootSource string

	// TempDir is a temporary directory for builders to store files in.
	TempDir string

	// ExtraFiles are files to add to the archive in addition to the Go
	// packages.
	//
	// Shared library dependencies will automatically also be added to the
	// archive using ldd, unless SkipLDD (below) is true.
	//
	// The following formats are allowed in the list:
	//
	//   - "/home/chrisko/foo:root/bar" adds the file from absolute path
	//     /home/chrisko/foo on the host at the relative root/bar in the
	//     archive.
	//   - "/home/foo" is equivalent to "/home/foo:home/foo".
	ExtraFiles []string

	// Symlinks to create in the archive. File path in archive -> target
	//
	// Target can be the name of a command. If not, it will be created as given.
	Symlinks map[string]string

	// If true, do not use ldd to pick up dependencies from local machine for
	// ExtraFiles. Useful if you have all deps revision controlled and wish to
	// ensure builds are repeatable, and/or if the local machine's binaries use
	// instructions unavailable on the emulated cpu.
	//
	// If you turn this on but do not manually list all deps, affected binaries
	// will misbehave.
	SkipLDD bool

	// OutputFile is the archive output file.
	OutputFile initramfs.WriteOpener

	// BaseArchive is an existing initramfs to include in the resulting
	// initramfs.
	BaseArchive initramfs.ReadOpener

	// UseExistingInit determines whether the existing init from
	// BaseArchive should be used.
	//
	// If this is false, the "init" from BaseArchive will be renamed to
	// "inito" (init-original).
	UseExistingInit bool

	// InitCmd is the name of a command to link /init to.
	//
	// This can be an absolute path or the name of a command included in
	// Commands.
	//
	// If this is empty, no init symlink will be created, but a user may
	// still specify a command called init or include an /init file.
	InitCmd string

	// UinitCmd is the name of a command to link /bin/uinit to.
	//
	// This can be an absolute path or the name of a command included in
	// Commands.
	//
	// The u-root init will always attempt to fork/exec a uinit program,
	// and append arguments from both the kernel command-line
	// (uroot.uinitargs) as well as specified in UinitArgs.
	//
	// If this is empty, no uinit symlink will be created, but a user may
	// still specify a command called uinit or include a /bin/uinit file.
	UinitCmd string

	// UinitArgs are the arguments passed to /bin/uinit.
	UinitArgs []string

	// DefaultShell is the default shell to start after init.
	//
	// This can be an absolute path or the name of a command included in
	// Commands.
	//
	// This must be specified to have a default shell.
	DefaultShell string
}

// Modifier modifies uimage options.
type Modifier func(*Opts) error

// OptionsFor will creates Opts from the given modifiers.
func OptionsFor(mods ...Modifier) (*Opts, error) {
	o := &Opts{
		Env: golang.Default(),
	}
	if err := o.Apply(mods...); err != nil {
		return nil, err
	}
	return o, nil
}

// Create creates an initramfs from the given options o.
func (o *Opts) Create(l *llog.Logger) error {
	return CreateInitramfs(l, *o)
}

// Apply modifies o with the given modifiers.
func (o *Opts) Apply(mods ...Modifier) error {
	for _, mod := range mods {
		if mod != nil {
			if err := mod(o); err != nil {
				return err
			}
		}
	}
	return nil
}

// WithSkipLDD sets SkipLDD to true. If true, initramfs creation skips using
// ldd to pick up dependencies from the local file system when resolving
// ExtraFiles.
//
// Useful if you have all deps revision controlled and wish to ensure builds
// are repeatable, and/or if the local machine's binaries use instructions
// unavailable on the emulated CPU.
//
// If you turn this on but do not manually list all deps, affected binaries
// will misbehave.
func WithSkipLDD() Modifier {
	return func(o *Opts) error {
		o.SkipLDD = true
		return nil
	}
}

// WithReplaceEnv replaces the Go build environment.
func WithReplaceEnv(env *golang.Environ) Modifier {
	return func(o *Opts) error {
		o.Env = env
		return nil
	}
}

// WithEnv alters the Go build environment (e.g. build tags, GOARCH, GOOS env vars).
func WithEnv(gopts ...golang.Opt) Modifier {
	return func(o *Opts) error {
		if o.Env == nil {
			o.Env = golang.Default(gopts...)
		} else {
			o.Env.Apply(gopts...)
		}
		return nil
	}
}

// WithSymlink adds a symlink to the archive.
//
// Target can be the name of a command. If not, it will be created as given.
func WithSymlink(file string, target string) Modifier {
	return func(o *Opts) error {
		if o.Symlinks == nil {
			o.Symlinks = make(map[string]string)
		}
		if other, ok := o.Symlinks[file]; ok {
			return fmt.Errorf("%w: cannot add symlink for %q as %q, already points to %q", os.ErrExist, file, target, other)
		}
		o.Symlinks[file] = target
		return nil
	}
}

// WithFiles adds files to the archive.
//
// Shared library dependencies will automatically also be added to the archive
// using ldd, unless WithSkipLDD is set.
//
// The following formats are allowed in the list:
//
//   - "/home/chrisko/foo:root/bar" adds the file from absolute path
//     /home/chrisko/foo on the host at the relative root/bar in the archive.
//   - "/home/foo" is equivalent to "/home/foo:home/foo".
//   - "uroot_test.go" is equivalent to "uroot_test.go:uroot_test.go".
func WithFiles(file ...string) Modifier {
	return func(o *Opts) error {
		o.ExtraFiles = append(o.ExtraFiles, file...)
		return nil
	}
}

// WithCommands adds Go commands to compile and add to the archive.
//
// b is the method of building -- as a busybox or a binary.
//
// Currently allowed formats for cmd:
//
//   - package imports; e.g. github.com/u-root/u-root/cmds/ls
//   - globs of package imports; e.g. github.com/u-root/u-root/cmds/*
//   - paths to package directories; e.g. $GOPATH/src/github.com/u-root/u-root/cmds/ls
//   - globs of paths to package directories; e.g. ./cmds/*
//
// Directories may be relative or absolute, with or without globs.
// Globs are resolved using filepath.Glob.
func WithCommands(buildOpts *golang.BuildOpts, b builder.Builder, cmd ...string) Modifier {
	return func(o *Opts) error {
		o.AddCommands(Commands{
			Builder:   b,
			Packages:  cmd,
			BuildOpts: buildOpts,
		})
		return nil
	}
}

// WithBusyboxCommands adds Go commands to compile in a busybox and add to the
// archive.
//
// If there were already busybox commands added to the archive, the given cmd
// will be merged with them.
//
// Allowed formats for cmd are documented in [WithCommands].
func WithBusyboxCommands(cmd ...string) Modifier {
	return func(o *Opts) error {
		o.AddBusyboxCommands(cmd...)
		return nil
	}
}

// WithShellBang directs the busybox builder to use #! instead of symlinks.
func WithShellBang(b bool) Modifier {
	return func(o *Opts) error {
		for i, cmd := range o.Commands {
			if _, ok := cmd.Builder.(*builder.GBBBuilder); ok {
				// Make a copy, because the same object may
				// have been used in other builds.
				o.Commands[i].Builder = &builder.GBBBuilder{
					ShellBang: b,
				}
				return nil
			}
		}

		// Otherwise, add an empty builder with no packages.
		// AddBusyboxCommands/WithBusyboxCommands will append to this.
		//
		// Yeah, it's a hack, sue me.
		o.Commands = append(o.Commands, Commands{
			Builder: &builder.GBBBuilder{ShellBang: b},
		})
		return nil
	}
}

// WithBusyboxBuildOpts directs the busybox builder to use the given build opts.
//
// Overrides any previously defined build options.
func WithBusyboxBuildOpts(g *golang.BuildOpts) Modifier {
	return func(o *Opts) error {
		for i, cmd := range o.Commands {
			if _, ok := cmd.Builder.(*builder.GBBBuilder); ok {
				o.Commands[i].BuildOpts = g
				return nil
			}
		}

		// Otherwise, add an empty builder with no packages.
		// AddBusyboxCommands/WithBusyboxCommands will append to this.
		//
		// Yeah, it's a hack, sue me.
		o.Commands = append(o.Commands, Commands{
			Builder:   &builder.GBBBuilder{},
			BuildOpts: g,
		})
		return nil
	}
}

// WithBinaryCommands adds Go commands to compile as individual binaries and
// add to the archive.
//
// Allowed formats for cmd are documented in [WithCommands].
func WithBinaryCommands(cmd ...string) Modifier {
	return WithCommands(nil, builder.Binary, cmd...)
}

// WithBinaryCommandsOpts adds Go commands to compile as individual binaries
// and add to the archive.
//
// Allowed formats for cmd are documented in [WithCommands].
func WithBinaryCommandsOpts(gbOpts *golang.BuildOpts, cmd ...string) Modifier {
	return WithCommands(gbOpts, builder.Binary, cmd...)
}

// WithCoveredCommands adds Go commands to compile as individual binaries with
// -cover and -covermode=atomic for integration test coverage.
//
// Allowed formats for cmd are documented in [WithCommands].
func WithCoveredCommands(cmd ...string) Modifier {
	return WithCommands(&golang.BuildOpts{ExtraArgs: []string{"-cover", "-covermode=atomic"}}, builder.Binary, cmd...)
}

// WithOutput sets the archive output file.
func WithOutput(w initramfs.WriteOpener) Modifier {
	return func(o *Opts) error {
		o.OutputFile = w
		return nil
	}
}

// WithExistingInit sets whether an existing init from BaseArchive should remain the init.
//
// If not, it will be renamed inito.
func WithExistingInit(use bool) Modifier {
	return func(o *Opts) error {
		o.UseExistingInit = use
		return nil
	}
}

// WithCPIOOutput sets the archive output file to be a CPIO created at the given path.
func WithCPIOOutput(path string) Modifier {
	if path == "" {
		return nil
	}
	return WithOutput(&initramfs.CPIOFile{Path: path})
}

// WithOutputDir sets the archive output to be in the given directory.
func WithOutputDir(path string) Modifier {
	return WithOutput(&initramfs.Dir{Path: path})
}

// WithBase is an existing initramfs to include in the resulting initramfs.
func WithBase(base initramfs.ReadOpener) Modifier {
	return func(o *Opts) error {
		o.BaseArchive = base
		return nil
	}
}

// WithBaseFile is an existing initramfs read from a CPIO file at the given
// path to include in the resulting initramfs.
func WithBaseFile(path string) Modifier {
	if path == "" {
		return nil
	}
	return WithBase(&initramfs.CPIOFile{Path: path})
}

// WithBaseArchive is an existing initramfs to include in the resulting initramfs.
func WithBaseArchive(archive *cpio.Archive) Modifier {
	return WithBase(&initramfs.Archive{Archive: archive})
}

// WithUinitCommand is command to link to /bin/uinit with args.
//
// cmd will be tokenized by a very basic shlex.Split.
//
// This can be an absolute path or the name of a command included in
// Commands.
//
// The u-root init will always attempt to fork/exec a uinit program,
// and append arguments from both the kernel command-line
// (uroot.uinitargs) as well as those specified in cmd.
//
// If this is empty, no uinit symlink will be created, but a user may
// still specify a command called uinit or include a /bin/uinit file.
func WithUinitCommand(cmd string) Modifier {
	return func(opts *Opts) error {
		args := shlex.Split(cmd)
		if len(args) > 0 {
			opts.UinitCmd = args[0]
		} else {
			opts.UinitCmd = ""
		}
		if len(args) > 1 {
			opts.UinitArgs = args[1:]
		} else {
			opts.UinitArgs = nil
		}
		return nil
	}
}

// WithUinit is command to link to /bin/uinit with args.
//
// This can be an absolute path or the name of a command included in
// Commands.
//
// The u-root init will always attempt to fork/exec a uinit program,
// and append arguments from both the kernel command-line
// (uroot.uinitargs) as well as those specified in cmd.
func WithUinit(arg0 string, args ...string) Modifier {
	return func(opts *Opts) error {
		opts.UinitCmd = arg0
		opts.UinitArgs = args
		return nil
	}
}

// WithInit sets the name of a command to link /init to.
//
// This can be an absolute path or the name of a command included in
// Commands.
func WithInit(arg0 string) Modifier {
	return func(opts *Opts) error {
		opts.InitCmd = arg0
		return nil
	}
}

// WithShell sets the default shell to start after init, which is a symlink
// from /bin/sh.
//
// This can be an absolute path or the name of a command included in
// Commands.
func WithShell(arg0 string) Modifier {
	return func(opts *Opts) error {
		opts.DefaultShell = arg0
		return nil
	}
}

// WithTempDir sets a temporary directory to use for building commands.
func WithTempDir(dir string) Modifier {
	return func(o *Opts) error {
		o.TempDir = dir
		return nil
	}
}

// Create creates an initramfs from mods specifications.
func Create(l *llog.Logger, mods ...Modifier) error {
	o, err := OptionsFor(mods...)
	if err != nil {
		return err
	}
	return o.Create(l)
}

// CreateInitramfs creates an initramfs built to opts' specifications.
func CreateInitramfs(l *llog.Logger, opts Opts) error {
	if _, err := os.Stat(opts.TempDir); os.IsNotExist(err) {
		return fmt.Errorf("temp dir %q must exist: %w", opts.TempDir, err)
	}
	if opts.OutputFile == nil {
		return fmt.Errorf("must give output file")
	}

	env := golang.Default()
	if opts.Env != nil {
		env = opts.Env
	}
	files := initramfs.NewFiles()

	lookupEnv := findpkg.DefaultEnv()
	if opts.UrootSource != "" {
		lookupEnv.URootSource = opts.UrootSource
	}

	// Expand commands.
	for index, cmds := range opts.Commands {
		if len(cmds.Packages) == 0 {
			continue
		}
		paths, err := findpkg.ResolveGlobs(l.AtLevel(slog.LevelInfo), env, lookupEnv, cmds.Packages)
		if err != nil {
			return fmt.Errorf("%w: %w", errResolvePackage, err)
		}
		opts.Commands[index].Packages = paths
	}

	// Add each build mode's commands to the archive.
	for _, cmds := range opts.Commands {
		if len(cmds.Packages) == 0 {
			continue
		}
		builderTmpDir, err := os.MkdirTemp(opts.TempDir, "builder")
		if err != nil {
			return err
		}
		buildOpts := cmds.BuildOpts
		if buildOpts == nil {
			buildOpts = &golang.BuildOpts{}
		}

		// Build packages.
		bOpts := builder.Opts{
			Env:       env,
			BuildOpts: buildOpts,
			Packages:  cmds.Packages,
			TempDir:   builderTmpDir,
			BinaryDir: cmds.TargetDir(),
		}
		if err := cmds.Builder.Build(l, files, bOpts); err != nil {
			return fmt.Errorf("error building: %w", err)
		}
	}

	// Open the target initramfs file.
	archive := &initramfs.Opts{
		Files:           files,
		OutputFile:      opts.OutputFile,
		BaseArchive:     opts.BaseArchive,
		UseExistingInit: opts.UseExistingInit,
	}
	if err := ParseExtraFiles(l, archive.Files, opts.ExtraFiles, !opts.SkipLDD); err != nil {
		return err
	}
	if err := opts.addSymlinkTo(l, archive, opts.UinitCmd, "bin/uinit"); err != nil {
		return fmt.Errorf("%w: %w", err, errUinitSymlink)
	}
	if len(opts.UinitArgs) > 0 {
		if err := archive.AddRecord(cpio.StaticFile("etc/uinit.flags", fileflag.ArgvToFile(opts.UinitArgs), 0o444)); err != nil {
			return fmt.Errorf("%w: %w", err, errUinitArgs)
		}
	}
	if err := opts.addSymlinkTo(l, archive, opts.InitCmd, "init"); err != nil {
		return fmt.Errorf("%w: %w", err, errInitSymlink)
	}
	if err := opts.addSymlinkTo(l, archive, opts.DefaultShell, "bin/sh"); err != nil {
		return fmt.Errorf("%w: %w", err, errDefaultshSymlink)
	}
	if err := opts.addSymlinkTo(l, archive, opts.DefaultShell, "bin/defaultsh"); err != nil {
		return fmt.Errorf("%w: %w", err, errDefaultshSymlink)
	}
	for p, target := range opts.Symlinks {
		p = path.Clean(p)
		if len(p) >= 1 && p[0] == '/' {
			p = p[1:]
		}
		if err := opts.addSymlinkTo(l, archive, target, p); err != nil {
			return fmt.Errorf("%w: could not add additional symlink", err)
		}
	}

	// Finally, write the archive.
	if err := initramfs.Write(archive); err != nil {
		return fmt.Errorf("error archiving: %w", err)
	}
	return nil
}

var (
	errResolvePackage   = errors.New("failed to resolve package paths")
	errInitSymlink      = errors.New("specify -initcmd=\"\" to ignore this error and build without an init (or, did you specify a list, and are you missing github.com/u-root/u-root/cmds/core/init?)")
	errUinitSymlink     = errors.New("specify -uinitcmd=\"\" to ignore this error and build without a uinit")
	errDefaultshSymlink = errors.New("specify -defaultsh=\"\" to ignore this error and build without a shell")
	errSymlink          = errors.New("could not create symlink")
	errUinitArgs        = errors.New("could not add uinit arguments")
)

func (o *Opts) addSymlinkTo(l *llog.Logger, archive *initramfs.Opts, command string, source string) error {
	if len(command) == 0 {
		return nil
	}

	target, err := resolveCommandOrPath(command, o.Commands)
	if err != nil {
		if o.Commands != nil {
			return fmt.Errorf("%w from %q to %q: %w", errSymlink, source, command, err)
		}
		l.Errorf("Could not create symlink from %q to %q: %v", source, command, err)
		return nil
	}

	// Make a relative symlink from /source -> target
	//
	// E.g. bin/defaultsh -> target, so you need to
	// filepath.Rel(/bin, target) since relative symlinks are
	// evaluated from their PARENT directory.
	relTarget, err := filepath.Rel(filepath.Join("/", filepath.Dir(source)), target)
	if err != nil {
		return err
	}

	if err := archive.AddRecord(cpio.Symlink(source, relTarget)); err != nil {
		return fmt.Errorf("failed to add symlink %s -> %s to initramfs: %w", source, relTarget, err)
	}
	return nil
}

func resolveCommandOrPath(cmd string, cmds []Commands) (string, error) {
	if strings.ContainsRune(cmd, filepath.Separator) {
		return cmd, nil
	}

	// Each build mode has its own binary dir (/bbin or /bin or /ubin).
	//
	// Figure out which build mode the shell is in, and symlink to that
	// build mode.
	for _, c := range cmds {
		for _, p := range c.Packages {
			if name := path.Base(p); name == cmd {
				return path.Join("/", c.TargetDir(), cmd), nil
			}
		}
	}

	return "", fmt.Errorf("command or path %q not included in u-root build", cmd)
}

// ParseExtraFiles adds files from the extraFiles list to the archive.
//
// The following formats are allowed in the extraFiles list:
//
//   - "/home/chrisko/foo:root/bar" adds the file from absolute path
//     /home/chrisko/foo on the host at the relative root/bar in the
//     archive.
//   - "/home/foo" is equivalent to "/home/foo:home/foo".
//
// ParseExtraFiles will also add ldd-listed dependencies if lddDeps is true.
func ParseExtraFiles(l *llog.Logger, archive *initramfs.Files, extraFiles []string, lddDeps bool) error {
	var err error
	// Add files from command line.
	for _, file := range extraFiles {
		var src, dst string
		parts := strings.SplitN(file, ":", 2)
		if len(parts) == 2 {
			if len(parts[0]) == 0 {
				return fmt.Errorf("%w: invalid extra files %q", os.ErrInvalid, file)
			}
			if len(parts[1]) == 0 {
				return fmt.Errorf("%w: invalid extra files %q", os.ErrInvalid, file)
			}
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
			return fmt.Errorf("couldn't find absolute path for %q: %w", src, err)
		}
		if err := archive.AddFileNoFollow(src, dst); err != nil {
			return fmt.Errorf("couldn't add %q to archive: %w", file, err)
		}

		if lddDeps {
			// Users are frequently naming directories now, not just files.
			// Hence we must use walk here, not just check the one file.
			if err := filepath.Walk(src, func(name string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() {
					return nil
				}
				// Try to open it as an ELF. If that fails, we can skip the ldd
				// step. The file will still be included from above.
				f, err := elf.Open(name)
				if err != nil {
					return nil //nolint:nilerr
				}
				if err = f.Close(); err != nil {
					l.Warnf("Closing ELF file %q: %v", name, err)
				}
				// Pull dependencies in the case of binaries. If `path` is not
				// a binary, `libs` will just be empty.
				libs, err := ldd.FList(name)
				if err != nil {
					return fmt.Errorf("WARNING: couldn't add ldd dependencies for %q: %v", name, err)
				}
				for _, lib := range libs {
					if err := archive.AddFileNoFollow(lib, lib[1:]); err != nil {
						l.Warnf("WARNING: couldn't add ldd dependencies for %q: %v", lib, err)
					}
				}
				return nil
			}); err != nil {
				l.Warnf("Getting dependencies for %q: %v", src, err)
			}
		}
	}
	return nil
}

// AddCommands adds commands to the build.
func (o *Opts) AddCommands(c ...Commands) {
	o.Commands = append(o.Commands, c...)
}

// AddBusyboxCommands adds Go commands to the busybox build.
func (o *Opts) AddBusyboxCommands(pkgs ...string) {
	for i, cmds := range o.Commands {
		if _, ok := cmds.Builder.(*builder.GBBBuilder); ok {
			o.Commands[i].Packages = append(o.Commands[i].Packages, pkgs...)
			return
		}
	}

	// Not found? Add first busybox.
	o.AddCommands(BusyboxCmds(pkgs...)...)
}

// BinaryCmds returns a list of Commands with cmds built as a busybox.
func BinaryCmds(cmds ...string) []Commands {
	if len(cmds) == 0 {
		return nil
	}
	return []Commands{
		{
			Builder:  builder.Binary,
			Packages: cmds,
		},
	}
}

// BusyboxCmds returns a list of Commands with cmds built as a busybox.
func BusyboxCmds(cmds ...string) []Commands {
	if len(cmds) == 0 {
		return nil
	}
	return []Commands{
		{
			Builder:  builder.Busybox,
			Packages: cmds,
		},
	}
}
