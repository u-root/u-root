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
			"cmds/*",
		},
		// Core should be things you don't want to live without.
		"core": {
			"cmds/ansi",
			"cmds/boot",
			"cmds/cat",
			"cmds/cbmem",
			"cmds/chmod",
			"cmds/chroot",
			"cmds/cmp",
			"cmds/console",
			"cmds/cp",
			"cmds/cpio",
			"cmds/date",
			"cmds/dd",
			"cmds/df",
			"cmds/dhclient",
			"cmds/dirname",
			"cmds/dmesg",
			"cmds/echo",
			"cmds/false",
			"cmds/field",
			"cmds/find",
			"cmds/free",
			"cmds/freq",
			"cmds/gpgv",
			"cmds/gpt",
			"cmds/grep",
			"cmds/gzip",
			"cmds/hexdump",
			"cmds/hostname",
			"cmds/id",
			"cmds/init",
			"cmds/insmod",
			"cmds/installcommand",
			"cmds/io",
			"cmds/ip",
			"cmds/kexec",
			"cmds/kill",
			"cmds/lddfiles",
			"cmds/ln",
			"cmds/losetup",
			"cmds/ls",
			"cmds/lsmod",
			"cmds/mkdir",
			"cmds/mkfifo",
			"cmds/mknod",
			"cmds/modprobe",
			"cmds/mount",
			"cmds/msr",
			"cmds/mv",
			"cmds/netcat",
			"cmds/ntpdate",
			"cmds/pci",
			"cmds/ping",
			"cmds/printenv",
			"cmds/ps",
			"cmds/pwd",
			"cmds/pxeboot",
			"cmds/readlink",
			"cmds/rm",
			"cmds/rmmod",
			"cmds/rsdp",
			"cmds/rush",
			"cmds/seq",
			"cmds/shutdown",
			"cmds/sleep",
			"cmds/sort",
			"cmds/stty",
			"cmds/switch_root",
			"cmds/sync",
			"cmds/tail",
			"cmds/tee",
			"cmds/true",
			"cmds/truncate",
			"cmds/umount",
			"cmds/uname",
			"cmds/uniq",
			"cmds/unshare",
			"cmds/validate",
			"cmds/vboot",
			"cmds/wc",
			"cmds/wget",
			"cmds/which",
		},
		// Minimal should be things you can't live without.
		"minimal": {
			"cmds/cat",
			"cmds/chmod",
			"cmds/cmp",
			"cmds/console",
			"cmds/cp",
			"cmds/date",
			"cmds/dd",
			"cmds/df",
			"cmds/dhclient",
			"cmds/dmesg",
			"cmds/echo",
			"cmds/find",
			"cmds/free",
			"cmds/gpgv",
			"cmds/grep",
			"cmds/gzip",
			"cmds/hostname",
			"cmds/id",
			"cmds/init",
			"cmds/insmod",
			"cmds/io",
			"cmds/ip",
			"cmds/kexec",
			"cmds/kill",
			"cmds/ln",
			"cmds/losetup",
			"cmds/ls",
			"cmds/lsmod",
			"cmds/mkdir",
			"cmds/mknod",
			"cmds/modprobe",
			"cmds/mount",
			"cmds/msr",
			"cmds/mv",
			"cmds/pci",
			"cmds/ping",
			"cmds/printenv",
			"cmds/ps",
			"cmds/pwd",
			"cmds/readlink",
			"cmds/rm",
			"cmds/rmmod",
			"cmds/rush",
			"cmds/seq",
			"cmds/shutdown",
			"cmds/sleep",
			"cmds/sync",
			"cmds/tail",
			"cmds/tee",
			"cmds/truncate",
			"cmds/umount",
			"cmds/uname",
			"cmds/unshare",
			"cmds/wc",
			"cmds/wget",
			"cmds/which",
		},
		// coreboot-app minimal environment
		"coreboot-app": {
			"cmds/cbmem",
			"cmds/chroot",
			"cmds/insmod",
			"cmds/modprobe",
			"cmds/rmmod",
			"cmds/sshd",
			"cmds/switch_root",
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
