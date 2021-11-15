// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"

	"src.elv.sh/pkg/buildinfo"
	"src.elv.sh/pkg/prog"
	"src.elv.sh/pkg/shell"
)

func main() {
	os.Exit(prog.Run([3]*os.File{os.Stdin, os.Stdout, os.Stderr}, os.Args, buildinfo.Program, daemonStub{}, shell.Program{}))
}

type daemonStub struct{}

func (daemonStub) ShouldRun(f *prog.Flags) bool {
	return f.Daemon
}

func (daemonStub) Run(fds [3]*os.File, f *prog.Flags, args []string) error {
	return fmt.Errorf("daemon mode not supported in this build")
}
