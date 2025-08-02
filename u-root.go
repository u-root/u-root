// Copyright 2015-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command u-root builds CPIO archives with the given files and Go commands.
package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"runtime"

	"github.com/dustin/go-humanize"
	"github.com/u-root/gobusybox/src/pkg/golang"
	"github.com/u-root/mkuimage/uimage"
	"github.com/u-root/mkuimage/uimage/mkuimage"
	"github.com/u-root/uio/llog"
)

var errEmptyFilesArg = errors.New("empty argument to -files")

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
		return fmt.Errorf("last argument is -files:%w", errEmptyFilesArg)
	}

	// We know the last arg is not -files; scan the arguments for -files
	// followed by a switch.
	for i := 0; i < len(args)-1; i++ {
		if args[i] == "-files" && args[i+1][0] == '-' {
			return fmt.Errorf("-files argument %d is followed by a switch: %w", i, errEmptyFilesArg)
		}
	}

	return nil
}

func main() {
	log.SetFlags(log.Ltime)
	if err := checkArgs(os.Args...); err != nil {
		log.Fatal(err)
	}

	env := golang.Default(golang.DisableCGO())
	f := &mkuimage.Flags{
		Commands:      mkuimage.CommandFlags{Builder: "bb"},
		ArchiveFormat: "cpio",
		OutputFile:    defaultFile(env),
	}
	f.RegisterFlags(flag.CommandLine)

	l := llog.Default()
	l.RegisterVerboseFlag(flag.CommandLine, "v", slog.LevelDebug)

	tf := &mkuimage.TemplateFlags{}
	tf.RegisterFlags(flag.CommandLine)
	flag.Parse()

	// Set defaults.
	m := []uimage.Modifier{
		uimage.WithReplaceEnv(env),
		uimage.WithBaseArchive(uimage.DefaultRamfs()),
		uimage.WithCPIOOutput(defaultFile(env)),
		uimage.WithInit("init"),
		uimage.WithShellBang(runtime.GOOS == "plan9" || os.Getenv("GOOS") == "plan9"),
	}
	if golang.Default().GOOS != "plan9" {
		m = append(m, uimage.WithShell("gosh"))
	}

	pkgs := flag.Args()
	// Only add default packages if no config template was given.
	//
	// Otherwise, the template can't erase the default packages and all
	// templates would be forced to use cmds/core/*.
	if len(pkgs) == 0 && tf.Config == "" {
		pkgs = []string{"github.com/u-root/u-root/cmds/core/*"}
	}
	if err := mkuimage.CreateUimage(l, m, tf, f, pkgs); err != nil {
		l.Errorf("mkuimage error: %v", err)
		os.Exit(1)
	}

	if stat, err := os.Stat(f.OutputFile); err == nil && f.ArchiveFormat == "cpio" {
		l.Infof("Successfully built %q (size %d bytes -- %s).", f.OutputFile, stat.Size(), humanize.IBytes(uint64(stat.Size())))
	}
}

func defaultFile(env *golang.Environ) string {
	if len(env.GOOS) == 0 || len(env.GOARCH) == 0 {
		return "/tmp/initramfs.cpio"
	}
	return fmt.Sprintf("/tmp/initramfs.%s_%s.cpio", env.GOOS, env.GOARCH)
}
