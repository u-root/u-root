// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (!tinygo || tinygo.enable) && !plan9 && goshsmall && !goshliner

package main

import (
	"io"

	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

func runInteractive(runner *interp.Runner, parser *syntax.Parser, stdout, stderr io.Writer) error {
	return errNotImplemented
}
