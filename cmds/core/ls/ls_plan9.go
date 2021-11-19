// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build plan9
// +build plan9

package main

import (
	flag "github.com/spf13/pflag"
)

var (
	final = flag.BoolP("print-last", "p", false, "Print only the final path element of each file name")
)
