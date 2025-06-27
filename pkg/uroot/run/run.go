// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package run

import (
	"io"
	"os"
)

type Params struct {
	Wd             string
	Env            []string
	Stdin          io.Reader
	Stdout, Stderr io.Writer
}

func DefaultParams() Params {
	wd, _ := os.Getwd()
	return Params{
		Env:    os.Environ(),
		Wd:     wd,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

type RunMain func(env Params, args []string) int
