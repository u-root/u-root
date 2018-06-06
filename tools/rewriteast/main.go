// Copyright 2015-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// rewriteast rewrites a Go command to be bb-compatible.
package main

import (
	"flag"
	"go/importer"
	"log"

	"github.com/u-root/u-root/pkg/golang"
	"github.com/u-root/u-root/pkg/uroot"
)

var (
	pkg    = flag.String("pkg", "", "Go package path to rewrite")
	dest   = flag.String("dest", "", "Destination directory for rewritten package files")
	bbPath = flag.String("bbpath", "", "bb package Go import path")
)

func main() {
	flag.Parse()

	env := golang.Default()
	importer := importer.For("source", nil)

	if err := uroot.RewritePackage(env, *pkg, *dest, *bbPath, importer); err != nil {
		log.Fatalf("failed to rewrite package %v: %v", *pkg, err)
	}
}
