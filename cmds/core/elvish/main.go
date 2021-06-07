// Copyright (c) elvish developers and contributors, All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found at https://github.com/elves/elvish/blob/master/LICENSE.

// Elvish is a cross-platform shell, supporting Linux, BSDs and Windows. It
// features an expressive programming language, with features like namespacing
// and anonymous functions, and a fully programmable user interface with
// friendly defaults. It is suitable for both interactive use and scripting.
package main

import (
	"os"

	"src.elv.sh/pkg/buildinfo"
	"src.elv.sh/pkg/daemon"
	"src.elv.sh/pkg/prog"
	"src.elv.sh/pkg/shell"
)

func main() {
	os.Exit(prog.Run(
		[3]*os.File{os.Stdin, os.Stdout, os.Stderr}, os.Args,
		buildinfo.Program, daemon.Program, shell.Program))
}
