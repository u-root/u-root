// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package mkuimage defines mkuimage flags and creation function.
package mkuimage

import (
	"flag"
	"fmt"
	"os"

	"github.com/u-root/gobusybox/src/pkg/golang"
	"github.com/u-root/gobusybox/src/pkg/uflag"
	"github.com/u-root/mkuimage/uimage"
	"github.com/u-root/mkuimage/uimage/builder"
	"github.com/u-root/mkuimage/uimage/templates"
)

type optionalStringVar struct {
	s **string
}

// Set implements flag.Value.Set.
func (o optionalStringVar) Set(s string) error {
	*o.s = &s
	return nil
}

func (o *optionalStringVar) String() string {
	if o == nil || o.s == nil || *(o.s) == nil {
		return ""
	}
	return **(o.s)
}

// CommandFlags are flags related to Go commands to be built by mkuimage.
type CommandFlags struct {
	NoCommands bool
	Builder    string
	ShellBang  bool
	Mod        golang.ModBehavior
	BuildTags  []string
	BuildOpts  *golang.BuildOpts
}

// RegisterFlags registers flags related to Go commands being built.
func (c *CommandFlags) RegisterFlags(f *flag.FlagSet) {
	f.StringVar(&c.Builder, "build", c.Builder, "uimage command build format (e.g. bb/gbb or binary).")
	f.BoolVar(&c.NoCommands, "nocmd", c.NoCommands, "Build no Go commands; initramfs only")
	f.BoolVar(&c.ShellBang, "shellbang", c.ShellBang, "Use #! instead of symlinks for busybox")
	if c.BuildOpts == nil {
		c.BuildOpts = &golang.BuildOpts{}
	}
	c.BuildOpts.RegisterFlags(f)
	// Register an alias for -go-no-strip for backwards compatibility.
	f.BoolVar(&c.BuildOpts.NoStrip, "no-strip", false, "Build unstripped binaries")

	// Flags for golang.Environ.
	f.StringVar((*string)(&c.Mod), "go-mod", string(c.Mod), "Value of -mod to go commands (allowed: (empty), vendor, mod, readonly)")
	// Register an alias for -go-build-tags for backwards compatibility.
	f.Var((*uflag.Strings)(&c.BuildTags), "tags", "Go build tags -- repeat the flag for multiple values")
	f.Var((*uflag.Strings)(&c.BuildTags), "go-build-tags", "Go build tags -- repeat the flag for multiple values")
}

// Modifiers turns the flag values into uimage modifiers.
func (c *CommandFlags) Modifiers(packages ...string) ([]uimage.Modifier, error) {
	if c.NoCommands {
		// Later modifiers may still add packages, so let's set the right environment.
		return []uimage.Modifier{
			uimage.WithEnv(golang.WithBuildTag(c.BuildTags...), func(e *golang.Environ) {
				e.Mod = c.Mod
			}),
		}, nil
	}

	switch c.Builder {
	case "bb", "gbb":
		return []uimage.Modifier{
			uimage.WithEnv(golang.WithBuildTag(c.BuildTags...), func(e *golang.Environ) {
				e.Mod = c.Mod
			}),
			uimage.WithBusyboxCommands(packages...),
			uimage.WithShellBang(c.ShellBang),
			uimage.WithBusyboxBuildOpts(c.BuildOpts),
		}, nil
	case "binary":
		return []uimage.Modifier{
			uimage.WithEnv(golang.WithBuildTag(c.BuildTags...), func(e *golang.Environ) {
				e.Mod = c.Mod
			}),
			uimage.WithCommands(c.BuildOpts, builder.Binary, packages...),
		}, nil
	default:
		return nil, fmt.Errorf("%w: could not find binary builder format %q", os.ErrInvalid, c.Builder)
	}
}

// String can be used to fill in values for Init, Uinit, and Shell.
func String(s string) *string {
	return &s
}

// Flags are mkuimage command-line flags.
type Flags struct {
	TempDir     *string
	KeepTempDir bool

	Init  *string
	Uinit *string
	Shell *string

	Files []string

	BaseArchive     string
	ArchiveFormat   string
	OutputFile      string
	UseExistingInit bool

	Commands CommandFlags
}

// Modifiers return uimage modifiers created from the flags.
func (f *Flags) Modifiers(packages ...string) ([]uimage.Modifier, error) {
	m := []uimage.Modifier{
		uimage.WithFiles(f.Files...),
		uimage.WithExistingInit(f.UseExistingInit),
	}
	if f.TempDir != nil {
		m = append(m, uimage.WithTempDir(*f.TempDir))
	}
	if f.BaseArchive != "" {
		// ArchiveFormat does not determine this, as only CPIO is supported.
		m = append(m, uimage.WithBaseFile(f.BaseArchive))
	}
	if f.Init != nil {
		m = append(m, uimage.WithInit(*f.Init))
	}
	if f.Uinit != nil {
		m = append(m, uimage.WithUinitCommand(*f.Uinit))
	}
	if f.Shell != nil {
		m = append(m, uimage.WithShell(*f.Shell))
	}
	switch f.ArchiveFormat {
	case "cpio":
		m = append(m, uimage.WithCPIOOutput(f.OutputFile))
	case "dir":
		m = append(m, uimage.WithOutputDir(f.OutputFile))
	default:
		return nil, fmt.Errorf("%w: could not find output format %q", os.ErrInvalid, f.ArchiveFormat)
	}
	more, err := f.Commands.Modifiers(packages...)
	if err != nil {
		return nil, err
	}
	return append(m, more...), nil
}

// RegisterFlags registers flags.
func (f *Flags) RegisterFlags(fs *flag.FlagSet) {
	fs.Var(&optionalStringVar{&f.TempDir}, "tmp-dir", "Temporary directory to build binary and archive in. Deleted after build if --keep-tmp-dir is not set.")
	fs.BoolVar(&f.KeepTempDir, "keep-tmp-dir", f.KeepTempDir, "Keep temporary directory after build")

	fs.Var(&optionalStringVar{&f.Init}, "initcmd", "Symlink target for /init. Can be an absolute path or a Go command name. Use initcmd=\"\" if you don't want the symlink.")
	fs.Var(&optionalStringVar{&f.Uinit}, "uinitcmd", "Symlink target and arguments for /bin/uinit. Can be an absolute path or a Go command name, followed by command-line args. Use uinitcmd=\"\" if you don't want the symlink. E.g. -uinitcmd=\"echo foobar\"")
	fs.Var(&optionalStringVar{&f.Shell}, "defaultsh", "Default shell. Can be an absolute path or a Go command name. Use defaultsh=\"\" if you don't want the symlink.")

	fs.Var((*uflag.Strings)(&f.Files), "files", "Additional files, directories, and binaries (with their ldd dependencies) to add to archive. Can be specified multiple times.")

	fs.StringVar(&f.BaseArchive, "base", f.BaseArchive, "Base archive to add files to. By default, this is a couple of directories like /bin, /etc, etc. Has a default internally supplied set of files; use base=/dev/null if you don't want any base files.")
	fs.StringVar(&f.ArchiveFormat, "format", f.ArchiveFormat, "Archival input (for -base) and output (for -o) format.")
	fs.StringVar(&f.OutputFile, "o", f.OutputFile, "Path to output initramfs file.")
	fs.BoolVar(&f.UseExistingInit, "useinit", f.UseExistingInit, "Use existing init from base archive (only if --base was specified).")
	fs.BoolVar(&f.UseExistingInit, "use-init", f.UseExistingInit, "Use existing init from base archive (only if --base was specified).")

	f.Commands.RegisterFlags(fs)
}

// TemplateFlags are flags for uimage config templates.
type TemplateFlags struct {
	File   string
	Config string
}

// RegisterFlags registers template flags.
func (tc *TemplateFlags) RegisterFlags(f *flag.FlagSet) {
	f.StringVar(&tc.Config, "config", "", "Config to pick from templates")
	f.StringVar(&tc.File, "config-file", "", "Config file to read from (default: finds .mkuimage.yaml in cwd or parents)")
}

// Get turns template flags into templates.
func (tc *TemplateFlags) Get() (*templates.Templates, error) {
	if tc.File != "" {
		return templates.TemplateFromFile(tc.File)
	}

	tpl, err := templates.Template()
	// Only complain about not finding a template if user requested a templated config.
	if err != nil && tc.Config != "" {
		return nil, err
	}
	return tpl, nil
}
