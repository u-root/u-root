// Copyright 2021-2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9
// +build !plan9

package main

import (
	"errors"
	"os"

	"src.elv.sh/pkg/buildinfo"
	"src.elv.sh/pkg/prog"
	"src.elv.sh/pkg/shell"
)

var ErrNotSupported = errors.New("daemon mode is not supported in this build")

func main() {
	os.Exit(run(os.Stdin, os.Stdout, os.Stderr, os.Args))
}

func run(stdin, stdout, stderr *os.File, args []string) int {
	return prog.Run([3]*os.File{stdin, stdout, stderr}, args, prog.Composite(buildinfo.Program, daemonStub{}, shell.Program{ActivateDaemon: nil}))
}

type daemonStub struct{}

func (daemonStub) Run(fds [3]*os.File, f *prog.Flags, args []string) error {
	if f.Daemon {
		return ErrNotSupported
	}
	return prog.ErrNotSuitable
}
