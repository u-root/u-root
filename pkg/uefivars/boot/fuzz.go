// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: BSD-3-Clause
//

//go:build gofuzz

package boot

import (
	"log"
	"os"
)

/*
go get github.com/dvyukov/go-fuzz/go-fuzz
go get github.com/dvyukov/go-fuzz/go-fuzz-build

go-fuzz-build -func FuzzParseFilePathList github.com/u-root/u-root/pkg/uefivars/boot
go-fuzz -bin=./boot-fuzz.zip -workdir=fuzz
...
*/

var tmpdir = "/tmp/fuzz-resolve-workdir"

func init() {
	// divert logging - greatly increases exec speed
	null, err := os.OpenFile("/dev/null", os.O_WRONLY, 0o200)
	if err != nil {
		panic(err)
	}
	log.SetOutput(null)
	log.SetFlags(0)

	err = os.MkdirAll(tmpdir, 0o755)
	if err != nil {
		panic(err)
	}
}

func FuzzParseFilePathList(data []byte) int {
	list, err := ParseFilePathList(data)
	if err != nil {
		return 0
	}
	_ = list.String()
	for _, p := range list {
		r, err := p.Resolver()
		if err != nil {
			continue
		}
		_ = r.String()
		_ = r.BlockInfo()
		_, _ = r.Resolve(tmpdir)
	}
	return 1
}
