// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build tinygo && amd64 && !noasm

package main

// Intentional build failure if using tinygo without noasm tag for amd64
// required by github.com/klauspost/compress/flate.matchLen
//
// go/pkg/mod/github.com/klauspost/compress@v1.17.4/flate/level3.go:107: linker
// could not find symbol github.com/klauspost/compress/flate.matchLen
var _ = TINYGO_USE_TAG_NOASM
