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

	"github.com/u-root/u-root/pkg/golang"
	"github.com/u-root/u-root/pkg/uroot"
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
	initCmd                                 *string
	defaultShell                            *string
	useExistingInit                         *bool
	extraFiles                              multiFlag
	templates                               = map[string][]string{
		"all": {
			"github.com/u-root/u-root/cmds/*",
		},
		// Core should be things you don't want to live without.
		"core": {
			"github.com/u-root/u-root/cmds/ansi",
			"github.com/u-root/u-root/cmds/boot",
			"github.com/u-root/u-root/cmds/cat",
			"github.com/u-root/u-root/cmds/cbmem",
			"github.com/u-root/u-root/cmds/chmod",
			"github.com/u-root/u-root/cmds/chroot",
			"github.com/u-root/u-root/cmds/cmp",
			"github.com/u-root/u-root/cmds/console",
			"github.com/u-root/u-root/cmds/cp",
			"github.com/u-root/u-root/cmds/cpio",
			"github.com/u-root/u-root/cmds/date",
			"github.com/u-root/u-root/cmds/dd",
			"github.com/u-root/u-root/cmds/df",
			"github.com/u-root/u-root/cmds/dhclient",
			"github.com/u-root/u-root/cmds/dirname",
			"github.com/u-root/u-root/cmds/dmesg",
			"github.com/u-root/u-root/cmds/echo",
			"github.com/u-root/u-root/cmds/false",
			"github.com/u-root/u-root/cmds/field",
			"github.com/u-root/u-root/cmds/find",
			"github.com/u-root/u-root/cmds/free",
			"github.com/u-root/u-root/cmds/freq",
			"github.com/u-root/u-root/cmds/gpgv",
			"github.com/u-root/u-root/cmds/gpt",
			"github.com/u-root/u-root/cmds/grep",
			"github.com/u-root/u-root/cmds/gzip",
			"github.com/u-root/u-root/cmds/hexdump",
			"github.com/u-root/u-root/cmds/hostname",
			"github.com/u-root/u-root/cmds/id",
			"github.com/u-root/u-root/cmds/init",
			"github.com/u-root/u-root/cmds/insmod",
			"github.com/u-root/u-root/cmds/installcommand",
			"github.com/u-root/u-root/cmds/io",
			"github.com/u-root/u-root/cmds/ip",
			"github.com/u-root/u-root/cmds/kexec",
			"github.com/u-root/u-root/cmds/kill",
			"github.com/u-root/u-root/cmds/lddfiles",
			"github.com/u-root/u-root/cmds/ln",
			"github.com/u-root/u-root/cmds/losetup",
			"github.com/u-root/u-root/cmds/ls",
			"github.com/u-root/u-root/cmds/lsmod",
			"github.com/u-root/u-root/cmds/mkdir",
			"github.com/u-root/u-root/cmds/mkfifo",
			"github.com/u-root/u-root/cmds/mknod",
			"github.com/u-root/u-root/cmds/modprobe",
			"github.com/u-root/u-root/cmds/mount",
			"github.com/u-root/u-root/cmds/msr",
			"github.com/u-root/u-root/cmds/mv",
			"github.com/u-root/u-root/cmds/netcat",
			"github.com/u-root/u-root/cmds/ntpdate",
			"github.com/u-root/u-root/cmds/pci",
			"github.com/u-root/u-root/cmds/ping",
			"github.com/u-root/u-root/cmds/printenv",
			"github.com/u-root/u-root/cmds/ps",
			"github.com/u-root/u-root/cmds/pwd",
			"github.com/u-root/u-root/cmds/pxeboot",
			"github.com/u-root/u-root/cmds/readlink",
			"github.com/u-root/u-root/cmds/rm",
			"github.com/u-root/u-root/cmds/rmmod",
			"github.com/u-root/u-root/cmds/rsdp",
			"github.com/u-root/u-root/cmds/rush",
			"github.com/u-root/u-root/cmds/seq",
			"github.com/u-root/u-root/cmds/shutdown",
			"github.com/u-root/u-root/cmds/sleep",
			"github.com/u-root/u-root/cmds/sort",
			"github.com/u-root/u-root/cmds/stty",
			"github.com/u-root/u-root/cmds/switch_root",
			"github.com/u-root/u-root/cmds/sync",
			"github.com/u-root/u-root/cmds/tail",
			"github.com/u-root/u-root/cmds/tee",
			"github.com/u-root/u-root/cmds/true",
			"github.com/u-root/u-root/cmds/truncate",
			"github.com/u-root/u-root/cmds/umount",
			"github.com/u-root/u-root/cmds/uname",
			"github.com/u-root/u-root/cmds/uniq",
			"github.com/u-root/u-root/cmds/unshare",
			"github.com/u-root/u-root/cmds/validate",
			"github.com/u-root/u-root/cmds/vboot",
			"github.com/u-root/u-root/cmds/wc",
			"github.com/u-root/u-root/cmds/wget",
			"github.com/u-root/u-root/cmds/which",
		},
		// Minimal should be things you can't live without.
		"minimal": {
			"github.com/u-root/u-root/cmds/cat",
			"github.com/u-root/u-root/cmds/chmod",
			"github.com/u-root/u-root/cmds/cmp",
			"github.com/u-root/u-root/cmds/console",
			"github.com/u-root/u-root/cmds/cp",
			"github.com/u-root/u-root/cmds/date",
			"github.com/u-root/u-root/cmds/dd",
			"github.com/u-root/u-root/cmds/df",
			"github.com/u-root/u-root/cmds/dhclient",
			"github.com/u-root/u-root/cmds/dmesg",
			"github.com/u-root/u-root/cmds/echo",
			"github.com/u-root/u-root/cmds/find",
			"github.com/u-root/u-root/cmds/free",
			"github.com/u-root/u-root/cmds/gpgv",
			"github.com/u-root/u-root/cmds/grep",
			"github.com/u-root/u-root/cmds/gzip",
			"github.com/u-root/u-root/cmds/hostname",
			"github.com/u-root/u-root/cmds/id",
			"github.com/u-root/u-root/cmds/init",
			"github.com/u-root/u-root/cmds/insmod",
			"github.com/u-root/u-root/cmds/io",
			"github.com/u-root/u-root/cmds/ip",
			"github.com/u-root/u-root/cmds/kexec",
			"github.com/u-root/u-root/cmds/kill",
			"github.com/u-root/u-root/cmds/ln",
			"github.com/u-root/u-root/cmds/losetup",
			"github.com/u-root/u-root/cmds/ls",
			"github.com/u-root/u-root/cmds/lsmod",
			"github.com/u-root/u-root/cmds/mkdir",
			"github.com/u-root/u-root/cmds/mknod",
			"github.com/u-root/u-root/cmds/modprobe",
			"github.com/u-root/u-root/cmds/mount",
			"github.com/u-root/u-root/cmds/msr",
			"github.com/u-root/u-root/cmds/mv",
			"github.com/u-root/u-root/cmds/pci",
			"github.com/u-root/u-root/cmds/ping",
			"github.com/u-root/u-root/cmds/printenv",
			"github.com/u-root/u-root/cmds/ps",
			"github.com/u-root/u-root/cmds/pwd",
			"github.com/u-root/u-root/cmds/readlink",
			"github.com/u-root/u-root/cmds/rm",
			"github.com/u-root/u-root/cmds/rmmod",
			"github.com/u-root/u-root/cmds/rush",
			"github.com/u-root/u-root/cmds/seq",
			"github.com/u-root/u-root/cmds/shutdown",
			"github.com/u-root/u-root/cmds/sleep",
			"github.com/u-root/u-root/cmds/sync",
			"github.com/u-root/u-root/cmds/tail",
			"github.com/u-root/u-root/cmds/tee",
			"github.com/u-root/u-root/cmds/truncate",
			"github.com/u-root/u-root/cmds/umount",
			"github.com/u-root/u-root/cmds/uname",
			"github.com/u-root/u-root/cmds/unshare",
			"github.com/u-root/u-root/cmds/wc",
			"github.com/u-root/u-root/cmds/wget",
			"github.com/u-root/u-root/cmds/which",
		},
		// coreboot-app minimal environment
		"coreboot-app": {
			"github.com/u-root/u-root/cmds/cat",
			"github.com/u-root/u-root/cmds/cbmem",
			"github.com/u-root/u-root/cmds/chroot",
			"github.com/u-root/u-root/cmds/console",
			"github.com/u-root/u-root/cmds/cp",
			"github.com/u-root/u-root/cmds/dd",
			"github.com/u-root/u-root/cmds/dhclient",
			"github.com/u-root/u-root/cmds/dmesg",
			"github.com/u-root/u-root/cmds/find",
			"github.com/u-root/u-root/cmds/grep",
			"github.com/u-root/u-root/cmds/id",
			"github.com/u-root/u-root/cmds/init",
			"github.com/u-root/u-root/cmds/insmod",
			"github.com/u-root/u-root/cmds/ip",
			"github.com/u-root/u-root/cmds/kill",
			"github.com/u-root/u-root/cmds/ls",
			"github.com/u-root/u-root/cmds/modprobe",
			"github.com/u-root/u-root/cmds/mount",
			"github.com/u-root/u-root/cmds/pci",
			"github.com/u-root/u-root/cmds/ping",
			"github.com/u-root/u-root/cmds/ps",
			"github.com/u-root/u-root/cmds/pwd",
			"github.com/u-root/u-root/cmds/rm",
			"github.com/u-root/u-root/cmds/rmmod",
			"github.com/u-root/u-root/cmds/rush",
			"github.com/u-root/u-root/cmds/shutdown",
			"github.com/u-root/u-root/cmds/sshd",
			"github.com/u-root/u-root/cmds/switch_root",
			"github.com/u-root/u-root/cmds/tail",
			"github.com/u-root/u-root/cmds/tee",
			"github.com/u-root/u-root/cmds/uname",
			"github.com/u-root/u-root/cmds/wget",
		},
	}
)

func init() {
	build = flag.String("build", "source", "u-root build format (e.g. bb or source).")
	format = flag.String("format", "cpio", "Archival format.")

	tmpDir = flag.String("tmpdir", "", "Temporary directory to put binaries in.")

	base = flag.String("base", "", "Base archive to add files to.")
	useExistingInit = flag.Bool("useinit", false, "Use existing init from base archive (only if --base was specified).")
	outputPath = flag.String("o", "", "Path to output initramfs file.")

	initCmd = flag.String("initcmd", "init", "Symlink target for /init. Can be an absolute path or a u-root command name.")
	defaultShell = flag.String("defaultsh", "rush", "Default shell. Can be an absolute path or a u-root command name.")

	flag.Var(&extraFiles, "files", "Additional files, directories, and binaries (with their ldd dependencies) to add to archive. Can be speficified multiple times.")
}

func main() {
	flag.Parse()

	// Main is in a separate functions so defers run on return.
	if err := Main(); err != nil {
		log.Fatal(err)
	}
	log.Printf("Successfully wrote initramfs.")
}

// Main is a separate function so defers are run on return, which they wouldn't
// on exit.
func Main() error {
	env := golang.Default()
	if env.CgoEnabled {
		log.Printf("Disabling CGO for u-root...")
		env.CgoEnabled = false
	}
	log.Printf("Build environment: %s", env)
	if env.GOOS != "linux" {
		log.Printf("GOOS is not linux. Did you mean to set GOOS=linux?")
	}

	builder, err := uroot.GetBuilder(*build)
	if err != nil {
		return err
	}
	archiver, err := uroot.GetArchiver(*format)
	if err != nil {
		return err
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
		var err error
		pkgs, err = uroot.DefaultPackageImports(env)
		if err != nil {
			return err
		}
	}

	// Open the target initramfs file.
	w, err := archiver.OpenWriter(*outputPath, env.GOOS, env.GOARCH)
	if err != nil {
		return err
	}

	var baseFile uroot.ArchiveReader
	if *base != "" {
		bf, err := os.Open(*base)
		if err != nil {
			return err
		}
		defer bf.Close()
		baseFile = archiver.Reader(bf)
	}

	opts := uroot.Opts{
		Env: env,
		// The command-line tool only allows specifying one build mode
		// right now.
		Commands: []uroot.Commands{
			{
				Builder:  builder,
				Packages: pkgs,
			},
		},
		Archiver:        archiver,
		TempDir:         tempDir,
		ExtraFiles:      extraFiles,
		OutputFile:      w,
		BaseArchive:     baseFile,
		UseExistingInit: *useExistingInit,
		InitCmd:         *initCmd,
		DefaultShell:    *defaultShell,
	}
	return uroot.CreateInitramfs(opts)
}
